package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/hashicorp/go-getter"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/zeta-chain/node/zetaclient/config"
)

const zetaclientdBinaryName = "zetaclientd"

var defaultUpgradesDir = os.ExpandEnv("$HOME/.zetaclientd/upgrades")

// serializedWriter wraps an io.Writer and ensures that writes to it from multiple goroutines
// are serialized
type serializedWriter struct {
	upstream io.Writer
	lock     sync.Mutex
}

func (w *serializedWriter) Write(p []byte) (n int, err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	return w.upstream.Write(p)
}

func getLogger(cfg config.Config, out io.Writer) zerolog.Logger {
	var logger zerolog.Logger
	switch cfg.LogFormat {
	case "json":
		logger = zerolog.New(out).Level(zerolog.Level(cfg.LogLevel)).With().Timestamp().Logger()
	case "text":
		logger = zerolog.New(zerolog.ConsoleWriter{Out: out, TimeFormat: time.RFC3339}).
			Level(zerolog.Level(cfg.LogLevel)).
			With().
			Timestamp().
			Logger()
	default:
		logger = zerolog.New(zerolog.ConsoleWriter{Out: out, TimeFormat: time.RFC3339})
	}

	return logger
}

type zetaclientdSupervisor struct {
	zetacoredConn      *grpc.ClientConn
	reloadSignals      chan bool
	logger             zerolog.Logger
	upgradesDir        string
	upgradePlanName    string
	enableAutoDownload bool
}

func newZetaclientdSupervisor(
	zetaCoreURL string,
	logger zerolog.Logger,
	enableAutoDownload bool,
) (*zetaclientdSupervisor, error) {
	logger = logger.With().Str("module", "zetaclientdSupervisor").Logger()
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:9090", zetaCoreURL),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc dial: %w", err)
	}
	return &zetaclientdSupervisor{
		zetacoredConn:      conn,
		logger:             logger,
		reloadSignals:      make(chan bool, 1),
		upgradesDir:        defaultUpgradesDir,
		enableAutoDownload: enableAutoDownload,
	}, nil
}

func (s *zetaclientdSupervisor) Start(ctx context.Context) {
	go s.watchForVersionChanges(ctx)
	go s.handleCoreUpgradePlan(ctx)
}

func (s *zetaclientdSupervisor) WaitForReloadSignal(ctx context.Context) {
	select {
	case <-s.reloadSignals:
	case <-ctx.Done():
	}
}

func (s *zetaclientdSupervisor) dirForVersion(version string) string {
	return path.Join(s.upgradesDir, version)
}

func atomicSymlink(target, linkName string) error {
	linkNameTmp := linkName + ".tmp"
	_, err := os.Stat(target)
	if err != nil {
		return fmt.Errorf("stat target: %w", err)
	}
	err = os.Remove(linkNameTmp)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("remove old current tmp: %w", err)
	}
	err = os.Symlink(target, linkNameTmp)
	if err != nil {
		return fmt.Errorf("new symlink: %w", err)
	}
	err = os.Rename(linkNameTmp, linkName)
	if err != nil {
		return fmt.Errorf("rename symlink: %w", err)
	}
	return nil
}

func (s *zetaclientdSupervisor) watchForVersionChanges(ctx context.Context) {
	client := tmservice.NewServiceClient(s.zetacoredConn)
	prevVersion := ""
	for {
		select {
		case <-time.After(time.Second):
		case <-ctx.Done():
			return
		}
		res, err := client.GetNodeInfo(ctx, &tmservice.GetNodeInfoRequest{})
		if err != nil {
			s.logger.Warn().Err(err).Msg("get node info")
			continue
		}
		newVersion := res.ApplicationVersion.Version
		if prevVersion == "" {
			prevVersion = newVersion
		}
		if prevVersion == newVersion {
			continue
		}
		s.logger.Warn().Msgf("core version change (%s -> %s)", prevVersion, newVersion)

		prevVersion = newVersion

		// TODO: just use newVersion when #2135 is merged
		// even without #2135, the version will still change and trigger the update
		newVersionDir := s.dirForVersion(s.upgradePlanName)
		currentLinkPath := s.dirForVersion("current")

		err = atomicSymlink(newVersionDir, currentLinkPath)
		if err != nil {
			s.logger.Error().
				Err(err).
				Msgf("unable to update current symlink (%s -> %s)", newVersionDir, currentLinkPath)
			return
		}
		s.reloadSignals <- true
	}
}

func (s *zetaclientdSupervisor) handleCoreUpgradePlan(ctx context.Context) {
	client := upgradetypes.NewQueryClient(s.zetacoredConn)

	prevPlanName := ""
	for {
		// wait for either a second or context cancel
		select {
		case <-time.After(time.Second):
		case <-ctx.Done():
			return
		}

		resp, err := client.CurrentPlan(ctx, &upgradetypes.QueryCurrentPlanRequest{})
		if err != nil {
			s.logger.Warn().Err(err).Msg("get current upgrade plan")
			continue
		}
		if resp.Plan == nil {
			continue
		}
		plan := resp.Plan
		if prevPlanName == plan.Name {
			continue
		}
		s.logger.Warn().Msgf("got new upgrade plan (%s)", plan.Name)
		prevPlanName = plan.Name
		s.upgradePlanName = plan.Name

		if !s.enableAutoDownload {
			s.logger.Warn().Msg("skipping autodownload because of configuration")
			continue
		}
		err = s.downloadZetaclientd(ctx, plan)
		if err != nil {
			s.logger.Error().Err(err).Msg("downloadZetaclientd failed")
		}
	}
}

// UpgradeConfig is expected format for the info field to allow auto-download
// this structure is copied from cosmosvisor
type upgradeConfig struct {
	Binaries map[string]string `json:"binaries"`
}

func (s *zetaclientdSupervisor) downloadZetaclientd(ctx context.Context, plan *upgradetypes.Plan) error {
	if plan.Info == "" {
		return errors.New("upgrade info empty")
	}
	var cfg upgradeConfig
	err := json.Unmarshal([]byte(plan.Info), &cfg)
	if err != nil {
		return fmt.Errorf("unmarshal upgrade config: %w", err)
	}

	s.logger.Info().Msg("downloading zetaclientd")

	binKey := fmt.Sprintf("%s-%s/%s", zetaclientdBinaryName, runtime.GOOS, runtime.GOARCH)
	binURL, ok := cfg.Binaries[binKey]
	if !ok {
		return fmt.Errorf("no binary found for: %s", binKey)
	}
	upgradeDir := s.dirForVersion(plan.Name)
	err = os.MkdirAll(upgradeDir, 0o750)
	if err != nil {
		return fmt.Errorf("mkdir %s: %w", upgradeDir, err)
	}
	upgradePath := path.Join(upgradeDir, zetaclientdBinaryName)
	// TODO: retry?
	// GetFile should validate checksum so long as it was provided in the url
	err = getter.GetFile(upgradePath, binURL, getter.WithContext(ctx), getter.WithUmask(0o750))
	if err != nil {
		return fmt.Errorf("get file %s: %w", binURL, err)
	}

	// ensure binary is executable
	info, err := os.Stat(upgradePath)
	if err != nil {
		return fmt.Errorf("stat binary: %w", err)
	}
	newMode := info.Mode().Perm() | 0o111
	err = os.Chmod(upgradePath, newMode)
	if err != nil {
		return fmt.Errorf("chmod %s: %w", upgradePath, err)
	}
	return nil
}

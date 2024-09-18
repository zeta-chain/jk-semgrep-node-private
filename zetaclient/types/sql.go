package types

import (
	"gorm.io/gorm"
)

const (
	// LastBlockNumID is the identifier to access the last block number in the database
	LastBlockNumID = 0xBEEF

	// LastTxHashID is the identifier to access the last transaction hash in the database
	LastTxHashID = 0xBEF0
)

// LastBlockSQLType is a model for storing the last block number
type LastBlockSQLType struct {
	gorm.Model
	Num uint64
}

// LastTransactionSQLType is a model for storing the last transaction hash
type LastTransactionSQLType struct {
	gorm.Model
	Hash string
}

// ToLastBlockSQLType converts a last block number to a LastBlockSQLType
func ToLastBlockSQLType(lastBlock uint64) *LastBlockSQLType {
	return &LastBlockSQLType{
		Model: gorm.Model{ID: LastBlockNumID},
		Num:   lastBlock,
	}
}

// ToLastTxHashSQLType converts a last transaction hash to a LastTransactionSQLType
func ToLastTxHashSQLType(lastTx string) *LastTransactionSQLType {
	return &LastTransactionSQLType{
		Model: gorm.Model{ID: LastTxHashID},
		Hash:  lastTx,
	}
}

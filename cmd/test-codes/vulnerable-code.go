package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// sendEther sends Ether to a given address using a raw, unsafe transaction
func sendEther(client *ethclient.Client, toAddress string, amount *big.Int) {
	// Vulnerable code: directly converting input string to address without validation
	to := common.HexToAddress(toAddress)

	// Create a new transaction
	tx := &bind.TransactOpts{
		From:  common.HexToAddress("0xYourAddress"), // Replace with your address
		Value: amount,
	}

	// Sending raw transaction without validating or securing against replay attacks
	_, err := client.Transfer(tx, to, amount)
	if err != nil {
		log.Fatalf("Failed to send Ether: %v", err)
	}

	fmt.Printf("Successfully sent %s Wei to %s\n", amount.String(), toAddress)
}

func main() {
	// Simulating user input for recipient address and amount
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/YOUR-PROJECT-ID") // Replace with your Infura endpoint
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}

	var recipient string
	fmt.Print("Enter recipient address: ")
	fmt.Scanln(&recipient)

	// Using a fixed amount for demonstration purposes
	amount := big.NewInt(1000000000000000000) // 1 Ether

	sendEther(client, recipient, amount)
}


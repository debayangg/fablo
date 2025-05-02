package main

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing key-value pairs
type SmartContract struct {
	contractapi.Contract
}

// Set stores a new key-value pair
func (s *SmartContract) Set(ctx contractapi.TransactionContextInterface, key string, value string) error {
	return ctx.GetStub().PutState(key, []byte(value))
}

// Get retrieves the value for a given key
func (s *SmartContract) Get(ctx contractapi.TransactionContextInterface, key string) (string, error) {
	value, err := ctx.GetStub().GetState(key)
	if err != nil {
		return "", fmt.Errorf("failed to read key: %v", err)
	}
	if value == nil {
		return "", fmt.Errorf("key not found: %s", key)
	}
	return string(value), nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		panic("Error creating chaincode: " + err.Error())
	}
	if err := chaincode.Start(); err != nil {
		panic("Error starting chaincode: " + err.Error())
	}
}
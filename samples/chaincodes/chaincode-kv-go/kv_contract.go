package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Response is a generic response structure for chaincode operations
type Response struct {
	Success string `json:"success,omitempty"`
	Error   string `json:"error,omitempty"`
}

// KVContract implements the chaincode logic for key-value operations
type KVContract struct {
	contractapi.Contract
}

// InitLedger invoked during chaincode instantiation to initialize any data
func (c *KVContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	// No initialization required
	return nil
}

// Put stores the given value under the specified key
func (c *KVContract) Put(ctx contractapi.TransactionContextInterface, key string, value string) (*Response, error) {
	err := ctx.GetStub().PutState(key, []byte(value))
	if err != nil {
		return &Response{Error: err.Error()}, err
	}
	return &Response{Success: "OK"}, nil
}

// Get retrieves the value stored under the specified key
func (c *KVContract) Get(ctx contractapi.TransactionContextInterface, key string) (*Response, error) {
	buffer, err := ctx.GetStub().GetState(key)
	if err != nil {
		return &Response{Error: err.Error()}, err
	}
	if buffer == nil || len(buffer) == 0 {
		return &Response{Error: "NOT_FOUND"}, nil
	}
	return &Response{Success: string(buffer)}, nil
}

// PutPrivateMessage stores a transient message into a private data collection
func (c *KVContract) PutPrivateMessage(ctx contractapi.TransactionContextInterface, collection string) (*Response, error) {
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return &Response{Error: err.Error()}, err
	}
	message, exists := transientMap["message"]
	if !exists {
		return &Response{Error: "TRANSIENT_KEY_NOT_FOUND"}, fmt.Errorf("transient key 'message' not found")
	}
	err = ctx.GetStub().PutPrivateData(collection, "message", message)
	if err != nil {
		return &Response{Error: err.Error()}, err
	}
	return &Response{Success: "OK"}, nil
}

// GetPrivateMessage retrieves the private message from a collection
func (c *KVContract) GetPrivateMessage(ctx contractapi.TransactionContextInterface, collection string) (*Response, error) {
	message, err := ctx.GetStub().GetPrivateData(collection, "message")
	if err != nil {
		return &Response{Error: err.Error()}, err
	}
	if message == nil {
		return &Response{Error: "NOT_FOUND"}, nil
	}
	return &Response{Success: string(message)}, nil
}

// VerifyPrivateMessage checks that the hash of the transient message matches stored private data hash
func (c *KVContract) VerifyPrivateMessage(ctx contractapi.TransactionContextInterface, collection string) (*Response, error) {
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return &Response{Error: err.Error()}, err
	}
	messageBytes, exists := transientMap["message"]
	if !exists {
		return &Response{Error: "TRANSIENT_KEY_NOT_FOUND"}, fmt.Errorf("transient key 'message' not found")
	}
	// Compute SHA256 of the transient message
	hash := sha256.New()
	hash.Write(messageBytes)
	currentHash := hex.EncodeToString(hash.Sum(nil))

	// Retrieve stored private data hash
	hashBytes, err := ctx.GetStub().GetPrivateDataHash(collection, "message")
	if err != nil {
		return &Response{Error: err.Error()}, err
	}
	storedHash := hex.EncodeToString(hashBytes)

	// Compare hashes
	if storedHash != currentHash {
		return &Response{Error: "VERIFICATION_FAILED"}, nil
	}
	return &Response{Success: "OK"}, nil
}

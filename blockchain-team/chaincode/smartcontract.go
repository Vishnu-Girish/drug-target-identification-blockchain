package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing medical recommendation hashes
type SmartContract struct {
	contractapi.Contract
}

// MedicalRecord defines the structural asset stored in the World State
type MedicalRecord struct {
	RecordID  string `json:"recordID"`
	Hash      string `json:"hash"`
	Timestamp string `json:"timestamp"`
}

// HistoryQueryResult structure used for returning historical audit trails
type HistoryQueryResult struct {
	TxId      string         `json:"txId"`
	Value     *MedicalRecord `json:"value"`
	Timestamp string         `json:"timestamp"`
	IsDelete  bool           `json:"isDelete"`
}

// CreateRecord stores a new cryptographic recommendation hash on the ledger
func (s *SmartContract) CreateRecord(ctx contractapi.TransactionContextInterface, recordID string, hash string) error {
	// Check if the record already exists in the World State
	exists, err := s.RecordExists(ctx, recordID)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if exists {
		return fmt.Errorf("the medical record %s already exists", recordID)
	}

	// Create the asset structure
	record := MedicalRecord{
		RecordID:  recordID,
		Hash:      hash,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// Serialize the asset into a JSON byte array
	recordBytes, err := json.Marshal(record)
	if err != nil {
		return err
	}

	// Commit the byte array to the ledger World State using native PutState()
	return ctx.GetStub().PutState(recordID, recordBytes)
}

// ReadRecord retrieves a stored medical record hash from the World State
func (s *SmartContract) ReadRecord(ctx contractapi.TransactionContextInterface, recordID string) (*MedicalRecord, error) {
	// Retrieve bytes from World State using native GetState()
	recordBytes, err := ctx.GetStub().GetState(recordID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if recordBytes == nil {
		return nil, fmt.Errorf("the medical record %s does not exist", recordID)
	}

	// Unmarshal the JSON bytes back into the Go structural asset
	var record MedicalRecord
	err = json.Unmarshal(recordBytes, &record)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

// VerifyRecord cross-checks an incoming hash against the immutable ledger record
func (s *SmartContract) VerifyRecord(ctx contractapi.TransactionContextInterface, recordID string, incomingHash string) (string, error) {
	record, err := s.ReadRecord(ctx, recordID)
	if err != nil {
		return "", err
	}

	// Deterministic string comparison of hashes
	if record.Hash == incomingHash {
		return "VALID", nil
	}

	return "TAMPERED", nil
}

// GetHistory unrolls the complete ledger lifecycle log for a specific record ID
func (s *SmartContract) GetHistory(ctx contractapi.TransactionContextInterface, recordID string) ([]HistoryQueryResult, error) {
	// Fetch structural iterator using native GetHistoryForKey()
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(recordID)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var record MedicalRecord
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &record)
			if err != nil {
				return nil, err
			}
		}

		// Parse block commitment timestamp
		txTimestamp := time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).Format(time.RFC3339)

		historyRecord := HistoryQueryResult{
			TxId:      response.TxId,
			Value:     &record,
			Timestamp: txTimestamp,
			IsDelete:  response.IsDelete,
		}
		records = append(records, historyRecord)
	}

	return records, nil
}

// RecordExists helper function to verify if an asset key currently has a value assigned
func (s *SmartContract) RecordExists(ctx contractapi.TransactionContextInterface, recordID string) (bool, error) {
	recordBytes, err := ctx.GetStub().GetState(recordID)
	if err != nil {
		return false, err
	}
	return recordBytes != nil, nil
}

// Main execution entrypoint to boot up the chaincode binary
func main() {
	cc, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		panic(fmt.Sprintf("Error creating drug-toxicity chaincode: %v", err))
	}

	if err := cc.Start(); err != nil {
		panic(fmt.Sprintf("Error starting drug-toxicity chaincode: %v", err))
	}
}
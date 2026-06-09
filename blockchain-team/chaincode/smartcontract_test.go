package main

import (
	"encoding/json"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// FakeStub simulates the ledger database state mapping using a native Go map
type FakeStub struct {
	shim.ChaincodeStubInterface
	State map[string][]byte
}

func (f *FakeStub) PutState(key string, value []byte) error {
	f.State[key] = value
	return nil
}

func (f *FakeStub) GetState(key string) ([]byte, error) {
	return f.State[key], nil
}

func (f *FakeStub) GetTxTimestamp() (*timestamppb.Timestamp, error) {
	return timestamppb.Now(), nil
}

// FakeContext simulates the invocation transaction runtime wrapper
type FakeContext struct {
	contractapi.TransactionContextInterface
	stub *FakeStub
}

func (f *FakeContext) GetStub() shim.ChaincodeStubInterface {
	return f.stub
}

func (f *FakeContext) GetClientIdentity() cid.ClientIdentity {
	return nil
}

func TestChaincodeWorkflow(t *testing.T) {
	// Initialize our independent database simulator
	stub := &FakeStub{State: make(map[string][]byte)}
	ctx := &FakeContext{stub: stub}

	contract := SmartContract{}
	recordID := "R001"
	expectedHash := "ABC123XYZ"

	// 1. Validate Ingestion Engine
	t.Run("Test CreateRecord", func(t *testing.T) {
		err := contract.CreateRecord(ctx, recordID, expectedHash)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		var record MedicalRecord
		bytes := stub.State[recordID]
		json.Unmarshal(bytes, &record)

		if record.Hash != expectedHash {
			t.Errorf("Expected hash %s, got %s", expectedHash, record.Hash)
		}
	})

	// 2. Validate Retrieval Engine
	t.Run("Test ReadRecord", func(t *testing.T) {
		record, err := contract.ReadRecord(ctx, recordID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if record.Hash != expectedHash {
			t.Errorf("Expected hash %s, got %s", expectedHash, record.Hash)
		}
	})

	// 3. Validate Match Verification Rules
	t.Run("Test VerifyRecord - VALID", func(t *testing.T) {
		status, err := contract.VerifyRecord(ctx, recordID, expectedHash)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if status != "VALID" {
			t.Errorf("Expected status VALID, got %s", status)
		}
	})

	t.Run("Test VerifyRecord - TAMPERED", func(t *testing.T) {
		status, err := contract.VerifyRecord(ctx, recordID, "HACKED_HASH")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if status != "TAMPERED" {
			t.Errorf("Expected status TAMPERED, got %s", status)
		}
	})
}
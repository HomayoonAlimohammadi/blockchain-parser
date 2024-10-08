package storage

import (
	"testing"

	"github.com/HomayoonAlimohammadi/blockchain-parser/internal/parser"
)

func TestAddTransactionFor(t *testing.T) {
	store := NewInMemory()
	address := "test_address"
	txn := parser.Transaction{Data: "txn1"}

	err := store.AddTransactionFor(address, txn)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	transactions, err := store.GetTransactionsFor(address)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(transactions) != 1 {
		t.Fatalf("expected 1 transaction, got %d", len(transactions))
	}

	if transactions[0].Data != txn.Data {
		t.Fatalf("expected transaction ID %s, got %s", txn.Data, transactions[0].Data)
	}
}

func TestGetTransactionsFor(t *testing.T) {
	store := NewInMemory()
	address := "test_address"
	txn1 := parser.Transaction{Data: "txn1"}
	txn2 := parser.Transaction{Data: "txn2"}

	store.AddTransactionFor(address, txn1)
	store.AddTransactionFor(address, txn2)

	transactions, err := store.GetTransactionsFor(address)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(transactions) != 2 {
		t.Fatalf("expected 2 transactions, got %d", len(transactions))
	}

	if transactions[0].Data != txn1.Data || transactions[1].Data != txn2.Data {
		t.Fatalf("expected transaction IDs %s and %s, got %s and %s", txn1.Data, txn2.Data, transactions[0].Data, transactions[1].Data)
	}
}

func TestGetTransactionsForEmptyAddress(t *testing.T) {
	store := NewInMemory()
	address := "non_existent_address"

	transactions, err := store.GetTransactionsFor(address)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(transactions) != 0 {
		t.Fatalf("expected 0 transactions, got %d", len(transactions))
	}
}

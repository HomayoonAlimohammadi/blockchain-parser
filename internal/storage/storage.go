package storage

import (
	"sync"

	"github.com/HomayoonAlimohammadi/blockchain-parser/internal/parser"
)

// NewInMemory creates a new in-memory storage
func NewInMemory() *inMemory {
	return &inMemory{
		mu:            &sync.RWMutex{},
		addressToTxns: make(map[string][]parser.Transaction),
	}
}

type inMemory struct {
	mu            *sync.RWMutex
	addressToTxns map[string][]parser.Transaction
}

// AddTransactionFor adds a transaction for a given address
func (s *inMemory) AddTransactionFor(address string, txn parser.Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.addressToTxns == nil {
		s.addressToTxns = make(map[string][]parser.Transaction)
	}

	s.addressToTxns[address] = append(s.addressToTxns[address], txn)

	return nil
}

// GetTransactionsFor returns the transactions for a given address
func (s *inMemory) GetTransactionsFor(address string) ([]parser.Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.addressToTxns[address], nil
}

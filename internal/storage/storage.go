package storage

import (
	"fmt"
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

// inMemory is an in-memory storage
type inMemory struct {
	mu            *sync.RWMutex
	addressToTxns map[string][]parser.Transaction
	activeAddrs   map[string]struct{}
}

// AddTransactionFor adds a transaction for a given address
func (s *inMemory) AddTransactionFor(address string, txn parser.Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.addressToTxns == nil {
		s.addressToTxns = make(map[string][]parser.Transaction)
	}

	if err := s.AddActiveAddress(address); err != nil {
		return fmt.Errorf("failed to add active address %q: %w", address, err)
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

// AddActiveAddress adds an address to the active list
func (s *inMemory) AddActiveAddress(address string) error {
	if s.activeAddrs == nil {
		s.activeAddrs = make(map[string]struct{})
	}

	s.activeAddrs[address] = struct{}{}
	return nil
}

// GetActiveAddresses returns the set of active addresses
func (s *inMemory) GetActiveAddresses() (map[string]struct{}, error) {
	return s.activeAddrs, nil
}

// RemoveActiveAddress removes an address from the active list
func (s *inMemory) RemoveActiveAddress(address string) error {
	delete(s.activeAddrs, address)
	return nil
}

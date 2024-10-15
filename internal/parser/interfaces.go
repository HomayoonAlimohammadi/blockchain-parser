package parser

import "context"

// Parser interface for blockchain parsing
type Parser interface {
	// GetCurrentBlock returns the current block number
	GetCurrentBlock(context.Context) (int, error)
	// Subscribe adds an address to the observer
	Subscribe(ctx context.Context, address string) error
	// GetTransactions returns the list of inbound or outbound transactions for an address
	GetTransactions(address string) ([]Transaction, error)
}

// Storage interface for storing transactions
type Storage interface {
	// AddActiveAddress adds an address to the set of observed addresses
	AddActiveAddress(address string) error
	// GetActiveAddresses returns the set of addresses being observed
	GetActiveAddresses() (map[string]struct{}, error)
	// RemoveActiveAddress removes an address from the set of observed addresses
	RemoveActiveAddress(address string) error
	// GetTransactionsFor returns the transactions for a given address
	GetTransactionsFor(address string) ([]Transaction, error)
	// AddTransactionFor adds a transaction for a given address
	AddTransactionFor(address string, txn Transaction) error
}

// RPCCaller calls methods of eth JSON RPC
type RPCCaller interface {
	// Subscribe calls the eth_subscribe method
	Subscribe(ctx context.Context, address string) (<-chan Transaction, error)
	// BlockNumber calls the eth_blockNumber method
	BlockNumber(ctx context.Context) (string, error)
}

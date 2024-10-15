package parser

import (
	"context"
	"fmt"
	"strconv"

	"github.com/HomayoonAlimohammadi/blockchain-parser/pkg/log"
)

// Transaction structure
type Transaction struct {
	Address          string   `json:"address"`
	BlockHash        string   `json:"blockHash"`
	BlockNumber      string   `json:"blockNumber"`
	Data             string   `json:"data"`
	LogIndex         string   `json:"logIndex"`
	Topics           []string `json:"topics"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
}

// EthereumParser implements the Parser interface
type EthereumParser struct {
	rpcCaller RPCCaller
	storage   Storage
}

// NewEthereumParser creates a new parser
func NewEthereumParser(rpcCaller RPCCaller, storage Storage) *EthereumParser {
	return &EthereumParser{
		rpcCaller: rpcCaller,
		storage:   storage,
	}
}

// GetCurrentBlock returns the current block number
func (p *EthereumParser) GetCurrentBlock(ctx context.Context) (int, error) {
	blockHex, err := p.rpcCaller.BlockNumber(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to call eth_blockNumber: %w", err)
	}

	result, err := strconv.ParseInt(blockHex[2:], 16, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse block hex: %w", err)
	}

	return int(result), nil
}

// Subscribe adds an address to the subscribed list
func (p *EthereumParser) Subscribe(ctx context.Context, address string) error {
	if alreadySubscribed, err := p.isAlreadySubscribed(address); err != nil {
		return fmt.Errorf("failed to check if address %q is already subscribed: %w", address, err)
	} else if alreadySubscribed {
		return fmt.Errorf("address %q already subscribed", address)
	}

	resChan, err := p.rpcCaller.Subscribe(ctx, address)
	if err != nil {
		return fmt.Errorf("failed to subscribe to address %q: %w", address, err)
	}

	if err := p.storage.AddActiveAddress(address); err != nil {
		return fmt.Errorf("failed to add active address %q: %w", address, err)
	}

	watchCtx := context.Background()
	go p.watchForTransactions(watchCtx, resChan, address)

	return nil
}

// GetTransactions returns the transactions for a given address
func (p *EthereumParser) GetTransactions(address string) ([]Transaction, error) {
	txns, err := p.storage.GetTransactionsFor(address)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions for address %q: %w", address, err)
	}

	return txns, nil
}

// watchForTransactions watches for transactions and adds them to the storage
func (p *EthereumParser) watchForTransactions(ctx context.Context, resChan <-chan Transaction, address string) {
	defer func() {
		if err := p.storage.RemoveActiveAddress(address); err != nil {
			log.Error(err, "failed to remove active address", "address", address)
		}
	}()

	log.Info("watching for transactions...", "address", address)
	for {
		select {
		case <-ctx.Done():
			log.Error(ctx.Err(), "context done")
			return
		case txn, ok := <-resChan:
			if !ok {
				log.Info("response channel close")
				return
			}

			log.Info("got transaction", "txn", txn)
			if err := p.storage.AddTransactionFor(address, txn); err != nil {
				log.Error(err, "failed to add transaction for address", "address", address)
			}
		}
	}
}

// isAlreadySubscribed checks if an address is already subscribed
func (p *EthereumParser) isAlreadySubscribed(address string) (bool, error) {
	activeAddrs, err := p.storage.GetActiveAddresses()
	if err != nil {
		return false, fmt.Errorf("failed to get active addresses: %w", err)
	}

	_, ok := activeAddrs[address]
	return ok, nil
}

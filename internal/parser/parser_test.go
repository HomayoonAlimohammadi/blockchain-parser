package parser

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRPCCaller is a mock implementation of the RPCCaller interface
type MockRPCCaller struct {
	mock.Mock
}

func (m *MockRPCCaller) BlockNumber(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *MockRPCCaller) Subscribe(ctx context.Context, address string) (<-chan Transaction, error) {
	args := m.Called(ctx, address)
	return args.Get(0).(<-chan Transaction), args.Error(1)
}

// MockStorage is a mock implementation of the Storage interface
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) GetTransactionsFor(address string) ([]Transaction, error) {
	args := m.Called(address)
	return args.Get(0).([]Transaction), args.Error(1)
}

func (m *MockStorage) AddTransactionFor(address string, txn Transaction) error {
	args := m.Called(address, txn)
	return args.Error(0)
}

func TestGetCurrentBlock(t *testing.T) {
	ctx := context.Background()
	mockRPCCaller := new(MockRPCCaller)
	mockStorage := new(MockStorage)
	parser := NewEthereumParser(mockRPCCaller, mockStorage)

	mockRPCCaller.On("BlockNumber", ctx).Return("0x10", nil)

	blockNumber, err := parser.GetCurrentBlock(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 16, blockNumber)

	mockRPCCaller.AssertExpectations(t)
}

func TestGetCurrentBlock_Error(t *testing.T) {
	ctx := context.Background()
	mockRPCCaller := new(MockRPCCaller)
	mockStorage := new(MockStorage)
	parser := NewEthereumParser(mockRPCCaller, mockStorage)

	mockRPCCaller.On("BlockNumber", ctx).Return("", errors.New("rpc error"))

	blockNumber, err := parser.GetCurrentBlock(ctx)
	assert.Error(t, err)
	assert.Equal(t, 0, blockNumber)

	mockRPCCaller.AssertExpectations(t)
}

func TestSubscribe(t *testing.T) {
	ctx := context.Background()
	mockRPCCaller := new(MockRPCCaller)
	mockStorage := new(MockStorage)
	parser := NewEthereumParser(mockRPCCaller, mockStorage)

	resChan := make(chan Transaction)
	mockRPCCaller.On("Subscribe", ctx, "0xAddress").Return(resChan, nil)

	err := parser.Subscribe(ctx, "0xAddress")
	assert.NoError(t, err)

	mockRPCCaller.AssertExpectations(t)
}

func TestSubscribe_Error(t *testing.T) {
	ctx := context.Background()
	mockRPCCaller := new(MockRPCCaller)
	mockStorage := new(MockStorage)
	parser := NewEthereumParser(mockRPCCaller, mockStorage)

	mockRPCCaller.On("Subscribe", ctx, "0xAddress").Return(nil, errors.New("subscribe error"))

	err := parser.Subscribe(ctx, "0xAddress")
	assert.Error(t, err)

	mockRPCCaller.AssertExpectations(t)
}

func TestGetTransactions(t *testing.T) {
	mockRPCCaller := new(MockRPCCaller)
	mockStorage := new(MockStorage)
	parser := NewEthereumParser(mockRPCCaller, mockStorage)

	expectedTxns := []Transaction{
		{Address: "0xAddress1"},
		{Address: "0xAddress2"},
	}
	mockStorage.On("GetTransactionsFor", "0xAddress").Return(expectedTxns, nil)

	txns, err := parser.GetTransactions("0xAddress")
	assert.NoError(t, err)
	assert.Equal(t, expectedTxns, txns)

	mockStorage.AssertExpectations(t)
}

func TestGetTransactions_Error(t *testing.T) {
	mockRPCCaller := new(MockRPCCaller)
	mockStorage := new(MockStorage)
	parser := NewEthereumParser(mockRPCCaller, mockStorage)

	mockStorage.On("GetTransactionsFor", "0xAddress").Return(nil, errors.New("storage error"))

	txns, err := parser.GetTransactions("0xAddress")
	assert.Error(t, err)
	assert.Nil(t, txns)

	mockStorage.AssertExpectations(t)
}

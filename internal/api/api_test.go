package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/HomayoonAlimohammadi/blockchain-parser/internal/api"
	parserpkg "github.com/HomayoonAlimohammadi/blockchain-parser/internal/parser"
)

type MockParser struct {
	mock.Mock
}

func (m *MockParser) Subscribe(ctx context.Context, address string) error {
	args := m.Called(ctx, address)
	return args.Error(0)
}

func (m *MockParser) GetTransactions(address string) ([]parserpkg.Transaction, error) {
	args := m.Called(address)
	return args.Get(0).([]parserpkg.Transaction), args.Error(1)
}

func (m *MockParser) GetCurrentBlock(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func TestSubscribeHandler(t *testing.T) {
	mockParser := new(MockParser)
	apiInstance := api.NewAPI(mockParser)

	t.Run("MethodNotAllowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/subscribe", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(apiInstance.SubscribeHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	})

	t.Run("BadRequest", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/subscribe", bytes.NewBuffer([]byte("invalid json")))
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(apiInstance.SubscribeHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Success", func(t *testing.T) {
		mockParser.On("Subscribe", mock.Anything, "test-address").Return(nil)

		body := map[string]string{"address": "test-address"}
		bodyBytes, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/subscribe", bytes.NewBuffer(bodyBytes))
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(apiInstance.SubscribeHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		mockParser.AssertExpectations(t)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		mockParser.On("Subscribe", mock.Anything, "test-address").Return(fmt.Errorf("error"))

		body := map[string]string{"address": "test-address"}
		bodyBytes, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/subscribe", bytes.NewBuffer(bodyBytes))
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(apiInstance.SubscribeHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockParser.AssertExpectations(t)
	})
}

func TestGetTransactionsHandler(t *testing.T) {
	mockParser := new(MockParser)
	apiInstance := api.NewAPI(mockParser)

	t.Run("MethodNotAllowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/transactions", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(apiInstance.GetTransactionsHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	})

	t.Run("Success", func(t *testing.T) {
		mockTransactions := []parserpkg.Transaction{{Data: "tx1"}, {Data: "tx2"}}
		mockParser.On("GetTransactions", "test-address").Return(mockTransactions, nil)

		req, _ := http.NewRequest(http.MethodGet, "/transactions?address=test-address", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(apiInstance.GetTransactionsHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockParser.AssertExpectations(t)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		mockParser.On("GetTransactions", "test-address").Return(nil, fmt.Errorf("error"))

		req, _ := http.NewRequest(http.MethodGet, "/transactions?address=test-address", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(apiInstance.GetTransactionsHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockParser.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockParser.On("GetTransactions", "test-address").Return(nil, nil)

		req, _ := http.NewRequest(http.MethodGet, "/transactions?address=test-address", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(apiInstance.GetTransactionsHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		mockParser.AssertExpectations(t)
	})
}

func TestGetBlockNumberHandler(t *testing.T) {
	mockParser := new(MockParser)
	apiInstance := api.NewAPI(mockParser)

	t.Run("MethodNotAllowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/blocknumber", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(apiInstance.GetBlockNumberHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	})

	t.Run("Success", func(t *testing.T) {
		mockParser.On("GetCurrentBlock", mock.Anything).Return(12345, nil)

		req, _ := http.NewRequest(http.MethodGet, "/blocknumber", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(apiInstance.GetBlockNumberHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockParser.AssertExpectations(t)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		mockParser.On("GetCurrentBlock", mock.Anything).Return(0, fmt.Errorf("error"))

		req, _ := http.NewRequest(http.MethodGet, "/blocknumber", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(apiInstance.GetBlockNumberHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockParser.AssertExpectations(t)
	})
}

package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	parserpkg "github.com/HomayoonAlimohammadi/blockchain-parser/internal/parser"
)

type api struct {
	parser parserpkg.Parser
}

// NewAPI creates a new API instance
func NewAPI(parser parserpkg.Parser) *api {
	return &api{
		parser: parser,
	}
}

// SubscribeHandler handles address subscription
func (a *api) SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		JSONError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed: %s", r.Method), nil)
		return
	}

	var req struct {
		Address string `json:"address"`
	}

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, fmt.Errorf("failed to decode request: %w", err), nil)
		return
	}

	if err := a.parser.Subscribe(r.Context(), req.Address); err != nil {
		JSONError(w, http.StatusInternalServerError, fmt.Errorf("failed to subscribe to address: %w", err), nil)
		return
	}

	JSONResponse(w, http.StatusCreated, "Address subscribed", nil)
}

// GetTransactionsHandler returns transactions for a given address
func (a *api) GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		JSONError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed: %s", r.Method), nil)
		return
	}

	address := r.URL.Query().Get("address")
	transactions, err := a.parser.GetTransactions(address)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, fmt.Errorf("failed to get transactions: %w", err), nil)
		return
	}

	if len(transactions) == 0 {
		JSONError(w, http.StatusNotFound, fmt.Errorf("no transactions found"), nil)
		return
	}

	resp := map[string]any{
		"transactions": transactions,
	}
	JSONResponse(w, http.StatusOK, "Transactions for address", resp)
}

// GetBlockNumberHandler returns the current block number
func (a *api) GetBlockNumberHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		JSONError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed: %s", r.Method), nil)
		return
	}

	blockNumber, err := a.parser.GetCurrentBlock(r.Context())
	if err != nil {
		JSONError(w, http.StatusInternalServerError, fmt.Errorf("failed to get block number: %w", err), nil)
		return
	}

	resp := map[string]any{
		"blockNumber": blockNumber,
	}
	JSONResponse(w, http.StatusOK, "Current block number", resp)
}

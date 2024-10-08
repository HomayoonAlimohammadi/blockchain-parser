package eth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HomayoonAlimohammadi/blockchain-parser/internal/parser"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestRPCCaller_BlockNumber(t *testing.T) {
	expectedResult := "0x10"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := RPCResponse{
			Jsonrpc: "2.0",
			ID:      1,
			Result:  expectedResult,
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := server.Client()
	rpcCaller := NewRPCCaller(client, nil)

	result, err := rpcCaller.BlockNumber(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestRPCCaller_Subscribe(t *testing.T) {
	expectedTxn := parser.Transaction{
		Data: "0x123",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _ := websocket.Upgrade(w, r, nil, 1024, 1024)
		defer conn.Close()

		// Send ack message
		conn.WriteMessage(websocket.TextMessage, []byte(`{"jsonrpc":"2.0","id":1,"result":"0x1"}`))

		// Send a transaction message
		txnMsg, _ := json.Marshal(expectedTxn)
		conn.WriteMessage(websocket.TextMessage, txnMsg)
	}))
	defer server.Close()

	wsDialer := websocket.DefaultDialer

	rpcCaller := NewRPCCaller(nil, wsDialer)
	resChan, err := rpcCaller.Subscribe(context.Background(), "0xAddress")
	assert.NoError(t, err)

	txn := <-resChan
	assert.Equal(t, expectedTxn, txn)
}

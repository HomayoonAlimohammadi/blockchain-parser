package eth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"

	"github.com/HomayoonAlimohammadi/blockchain-parser/internal/parser"
	"github.com/HomayoonAlimohammadi/blockchain-parser/pkg/log"
	wspkg "github.com/HomayoonAlimohammadi/blockchain-parser/pkg/websocket"
)

// JSON-RPC request structure
type RPCRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	ID      int    `json:"id"`
}

// JSON-RPC response structure
type RPCResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  string `json:"result"`
}

// RPC caller structure
type rpcCaller struct {
	client   *http.Client
	wsDialer *websocket.Dialer
}

// NewRPCCaller creates a new RPC caller
func NewRPCCaller(client *http.Client, wsDialer *websocket.Dialer) *rpcCaller {
	return &rpcCaller{
		client:   client,
		wsDialer: wsDialer,
	}
}

// Subscribe calls eth_subscribe
func (c *rpcCaller) Subscribe(ctx context.Context, address string) (<-chan parser.Transaction, error) {
	u := url.URL{Scheme: webSocketScheme, Host: rpcHost, Path: "/"}

	conn, _, err := c.wsDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to dial websocket: %w", err)
	}

	req := RPCRequest{
		Jsonrpc: rpcVersion,
		Method:  subscribeMethod,
		Params: []any{
			"logs",
			map[string]string{
				"address": address,
			},
		},
	}

	err = conn.WriteJSON(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send subscription message: %w", err)
	}

	// First message is the ack with different format
	if _, _, err := conn.ReadMessage(); err != nil {
		return nil, fmt.Errorf("failed to read ack message: %w", err)
	}

	resChan := make(chan parser.Transaction, 999999)
	// We use new context because the parent context is cancelled when the request is done
	listenCtx := context.Background()
	go listenForTxn(listenCtx, conn, resChan)

	return resChan, nil
}

// listenForTxn listens for transactions
func listenForTxn(ctx context.Context, conn *websocket.Conn, resChan chan<- parser.Transaction) {
	defer func() {
		close(resChan)
		conn.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			log.Error(ctx.Err(), "context done")
			return
		default:
			_, message, err := conn.ReadMessage()
			if wspkg.IsCloseError(err) {
				log.Error(err, "connection closed")
				return
			} else if err != nil {
				log.Error(err, "failed to read message")
			}

			var txn parser.Transaction
			if err := json.Unmarshal(message, &txn); err != nil {
				log.Error(err, "failed to unmarshal message")
			}

			select {
			case resChan <- txn:
			default:
				log.Warn("transaction missed", "txn", txn)
			}
		}
	}
}

// BlockNumber calls eth_blockNumber
func (c *rpcCaller) BlockNumber(ctx context.Context) (string, error) {
	reqBody := RPCRequest{
		Jsonrpc: rpcVersion,
		Method:  blockNumberMethod,
	}

	jsonReq, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	url := fmt.Sprintf("https://%s", rpcHost)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonReq))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	var rpcResp RPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return rpcResp.Result, nil
}

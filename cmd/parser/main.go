package main

import (
	"net/http"

	"github.com/HomayoonAlimohammadi/blockchain-parser/internal/api"
	parserpkg "github.com/HomayoonAlimohammadi/blockchain-parser/internal/parser"
	storagepkg "github.com/HomayoonAlimohammadi/blockchain-parser/internal/storage"
	"github.com/HomayoonAlimohammadi/blockchain-parser/pkg/eth"
	"github.com/HomayoonAlimohammadi/blockchain-parser/pkg/log"
	"github.com/gorilla/websocket"
)

func main() {
	rpcCaller := eth.NewRPCCaller(http.DefaultClient, websocket.DefaultDialer)
	storage := storagepkg.NewInMemory()
	parser := parserpkg.NewEthereumParser(rpcCaller, storage)

	api := api.NewAPI(parser)
	http.HandleFunc("/subscribe", api.SubscribeHandler)
	http.HandleFunc("/transactions", api.GetTransactionsHandler)
	http.HandleFunc("/blocknumber", api.GetBlockNumberHandler)

	log.Info("starting to listen on :8080")
	log.Error(http.ListenAndServe(":8080", nil), "failed to listen and serve")
}

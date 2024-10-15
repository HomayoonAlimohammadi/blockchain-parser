# Usage

To use Blockchain Parser, run the server with the following command:

```bash
go build -o parser ./cmd/parser
./parser 
```

To subscribe to an address, run:

```bash
curl -X POST -d '{"address": "0x28C6c06298d514Db089934071355E5743bf21d60"}' http://localhost:8080/subscribe
```

* NOTE: `0x28C6c06298d514Db089934071355E5743bf21d60` is the "Binance 14" with over 20M transactions and more than 235k ETH.

To get the transactions for an address:

```bash
curl http://localhost:8080/transactions\?address\=0x28C6c06298d514Db089934071355E5743bf21d60
```

To get the current block number:

```bash
curl http://localhost:8080/blocknumber
```

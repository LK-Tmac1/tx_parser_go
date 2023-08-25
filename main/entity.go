package main

// Internal business domain entities

type Transaction struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Value            string `json:"value"`
}

type Block struct {
	Number       string        `json:"number"`
	Hash         string        `json:"hash"`
	Transactions []Transaction `json:"transactions"`
}

// API request or response entities

type CurrentBlockResponse struct {
	CurrentBlock int64 `json:"current_block"`
}

type SubscribeRequest struct {
	Address string `json:"address"`
}

type SubscribeResponse struct {
	Subscribed bool `json:"subscribed"`
}

type GetTransactionsRequest struct {
	Address string `json:"address"`
}

type GetTransactionsResponse struct {
	Transactions *[]Transaction `json:"transactions"`
}

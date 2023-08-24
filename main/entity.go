package main

type Transaction struct {
	Hash     string `json:"hash"`
	From     string `json:"from"`
	To       string `json:"to"`
	Value    string `json:"value"`
	GasPrice string `json:"gasPrice"`
}

type CurrentBlock struct {
	Result string `json:"result"`
}

type CurrentBlockResponse struct {
	CurrentBlock int64 `json:"current_block"`
}

type SubscribeRequest struct {
	Address string `json:"address"`
}

type GetTransactionsRequest struct {
	Address string `json:"address"`
}

type GetTransactionsResponse struct {
	InboundTransactions  []Transaction `json:"inbound_transactions"`
	OutboundTransactions []Transaction `json:"outbound_transactions"`
}

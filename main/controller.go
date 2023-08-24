package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	ethRPCURL  = "https://cloudflare-eth.com"
	jsonRpcVer = "2.0"
	id         = 1 // a dummy value
)

type Parser struct {
	subscriptions map[string]bool
}

func NewParser() *Parser {
	return &Parser{
		subscriptions: make(map[string]bool),
	}
}

func (p *Parser) GetCurrentBlock() (*int64, error) {
	var params []map[string]interface{}
	rpcResponse, err := callEthRpc("eth_blockNumber", params)
	if err != nil {
		return nil, err
	}

	var currentBlock CurrentBlock
	if err := json.NewDecoder(rpcResponse.Body).Decode(&currentBlock); err != nil {
		return nil, err
	}

	result, err := hexStrToInt(currentBlock.Result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func hexStrToInt(s string) (*int64, error) {
	s = strings.Replace(s, "0x", "", 1)
	value, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		return nil, err
	}

	return &value, nil
}

func (p *Parser) Subscribe(address string) {
	p.subscriptions[address] = true
	fmt.Println(address)
}

func callEthRpc(methodName string, params []map[string]interface{}) (*http.Response, error) {
	payload := map[string]interface{}{
		"jsonrpc": jsonRpcVer,
		"method":  methodName,
		"params":  params,
		"id":      id,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(ethRPCURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (p *Parser) GetTransactions(address string) ([]Transaction, error) {
	if _, ok := p.subscriptions[address]; !ok {
		return nil, nil
	}

	params := []map[string]interface{}{
		{
			"address":   address,
			"fromBlock": "latest",
			"toBlock":   "latest",
		},
	}
	response, err := callEthRpc("eth_getLogs", params)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var rpcResponse struct {
		Result []Transaction `json:"result"`
	}

	if err := json.NewDecoder(response.Body).Decode(&rpcResponse); err != nil {
		return nil, err
	}

	return rpcResponse.Result, nil
}

func test1() {
	parser := NewParser()

	address := "0x407d73d8a49eeb85d32cf465507dd71d507100c1"
	parser.Subscribe(address)

	// Fetch transactions for the subscribed address
	transactions, err := parser.GetTransactions(address)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Transactions:")
	for _, tx := range transactions {
		fmt.Printf("Hash: %s, From: %s, To: %s, Value: %s\n", tx.Hash, tx.From, tx.To, tx.Value)
	}
	fmt.Println("Transactions done")
}

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
	ethRpcURL  = "https://cloudflare-eth.com"
	jsonRpcVer = "2.0"
	id         = 1 // a dummy value
)

type Parser struct {
	storage Storage
}

func NewParser(s Storage) *Parser {
	return &Parser{
		storage: s,
	}
}

func (p Parser) GetCurrentBlock() (*int64, error) {
	var params []map[string]interface{}
	rpcResponse, err := callEthRpc("eth_blockNumber", params)
	if err != nil {
		return nil, err
	}

	var currentBlock CurrentBlock
	if err = json.NewDecoder(rpcResponse.Body).Decode(&currentBlock); err != nil {
		return nil, err
	}

	result := hexStrToInt(currentBlock.Result)
	return &result, nil
}

func hexStrToInt(s string) int64 {
	s = strings.Replace(s, "0x", "", 1)
	value, _ := strconv.ParseInt(s, 16, 64)
	return value
}

func hexStrToIntStr(s string) string {
	return strconv.FormatInt(hexStrToInt(s), 10)
}

func (p Parser) Subscribe(address string) bool {
	fmt.Println("subscribing address", address)
	return p.storage.Subscribe(address) == nil
}

func (p Parser) GetTransactions(address string) ([]Transaction, error) {
	if yes, err := p.storage.IsSubscribed(address); err != nil {
		return nil, err
	} else if !*yes {
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

	t := p.getTransaction(address)
	rpcResponse.Result = append(rpcResponse.Result, *t)

	return rpcResponse.Result, nil
}

func callEthRpc(methodName string, params interface{}) (*http.Response, error) {
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

	resp, err := http.Post(ethRpcURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (p Parser) getTransaction(address string) *Transaction {
	params := []string{address}
	response, err := callEthRpc("eth_getTransactionByHash", params)
	if err != nil {
		fmt.Println("eth_getTransactionByHash call error", err)
		return nil
	}
	defer response.Body.Close()

	var rpcResponse struct {
		Result *Transaction `json:"result"`
	}
	if err := json.NewDecoder(response.Body).Decode(&rpcResponse); err != nil {
		fmt.Println("eth_getTransactionByHash encoding error", err)
		return nil
	}

	fmt.Println("mapping", rpcResponse.Result)
	convertTransactionHex(rpcResponse.Result)
	return rpcResponse.Result
}

func convertTransactionHex(t *Transaction) {
	t.BlockHash = hexStrToIntStr(t.BlockHash)
	t.BlockNumber = hexStrToIntStr(t.BlockNumber)
	t.From = hexStrToIntStr(t.From)
	t.Gas = hexStrToIntStr(t.Gas)
	t.GasPrice = hexStrToIntStr(t.GasPrice)
	t.Hash = hexStrToIntStr(t.Hash)
	t.Input = hexStrToIntStr(t.Input)
	t.Nonce = hexStrToIntStr(t.Nonce)
	t.To = hexStrToIntStr(t.To)
	t.TransactionIndex = hexStrToIntStr(t.TransactionIndex)
	t.Value = hexStrToIntStr(t.Value)
}

func test1() {
	parser := NewParser(NewInMemoryStore())
	address := "0x88df016429689c079f3b2f6ad39fa052532c56795b733da78a91ebe6a713944b"

	t := parser.getTransaction(address)
	if t != nil {
		fmt.Println("transaction %+v", *t)
	} else {
		fmt.Println("transaction nil")
	}
}

func test2() {
	parser := NewParser(NewInMemoryStore())
	address := "0x407d73d8a49eeb85d32cf465507dd71d507100c1"

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

func main2() {
	test1()
}

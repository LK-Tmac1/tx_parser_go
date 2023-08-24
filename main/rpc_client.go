package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	ethRpcURL  = "https://cloudflare-eth.com"
	jsonRpcVer = "2.0"
)

type RpcClient struct{}

func NewRpcClient() RpcClient {
	return RpcClient{}
}

func (c RpcClient) GetCurrentBlock() (*string, error) {
	var params []map[string]interface{}
	rpcResponse, err := callEthRpcHelper("eth_blockNumber", params)
	if err != nil {
		return nil, err
	}

	var currentBlock CurrentBlock
	if err = json.NewDecoder(rpcResponse.Body).Decode(&currentBlock); err != nil {
		return nil, err
	}

	return &currentBlock.Result, nil
}

func (c RpcClient) GetLogs(address string) (*[]Log, error) {
	params := []map[string]interface{}{{"address": []string{address}}}
	response, err := callEthRpcHelper("eth_getLogs", params)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// then get all transactions by those logs
	var rpcResponse struct {
		Result []Log `json:"result"`
	}
	if err := json.NewDecoder(response.Body).Decode(&rpcResponse); err != nil {
		return nil, err
	}

	return &rpcResponse.Result, err
}

func (c RpcClient) GetTransactionByHash(txnHash string) (*Transaction, error) {
	params := []string{txnHash}
	response, err := callEthRpcHelper("eth_getTransactionByHash", params)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var rpcResponse struct {
		Result *Transaction `json:"result"`
	}
	if err := json.NewDecoder(response.Body).Decode(&rpcResponse); err != nil {
		fmt.Println("eth_getTransactionByHash encoding error", err)
		return nil, err
	}
	return rpcResponse.Result, err
}

func callEthRpcHelper(methodName string, params interface{}) (*http.Response, error) {
	fmt.Println("calling eth rpc with methodName and params:", methodName, params)
	payload := map[string]interface{}{
		"jsonrpc": jsonRpcVer,
		"method":  methodName,
		"params":  params,
		"id":      "1", // dummy value
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("calling eth rpc with methodName and params with error", methodName, params, err)
		return nil, err
	}

	resp, err := http.Post(ethRpcURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("calling eth rpc with methodName and params with error", methodName, params, err)
		return nil, err
	}
	return resp, nil
}

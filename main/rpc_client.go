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

func (c RpcClient) GetLatestBlock() (*string, error) {
	var params []map[string]interface{}
	rpcResponse, err := callEthRpcHelper("eth_blockNumber", params)
	if err != nil {
		return nil, err
	}

	var latestBlock struct {
		Result string `json:"result"`
	}
	if err = json.NewDecoder(rpcResponse.Body).Decode(&latestBlock); err != nil {
		return nil, err
	}

	return &latestBlock.Result, nil
}

func (c RpcClient) GetBlockByNumber(hex string) (*Block, error) {
	params := []interface{}{hex, true}
	rpcResponse, err := callEthRpcHelper("eth_getBlockByNumber", params)
	if err != nil {
		return nil, err
	}
	defer rpcResponse.Body.Close()

	var blockByNumber struct {
		Result *Block `json:"result"`
	}

	if err := json.NewDecoder(rpcResponse.Body).Decode(&blockByNumber); err != nil {
		fmt.Println("eth_getBlockByNumber encoding error", err)
		return nil, err
	}

	return blockByNumber.Result, nil
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

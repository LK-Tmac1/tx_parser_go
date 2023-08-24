package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Parser struct {
	storage   Storage
	rpcClient RpcClient
}

func NewParser(s Storage) *Parser {
	return &Parser{
		storage:   s,
		rpcClient: NewRpcClient(),
	}
}

func (p Parser) GetCurrentBlock() (*int64, error) {
	result, err := p.rpcClient.GetCurrentBlock()
	if err != nil {
		return nil, err
	}

	s := strings.Replace(*result, "0x", "", 1)
	value, err := strconv.ParseInt(s, 16, 64)
	return &value, err
}

func (p Parser) Subscribe(address string) bool {
	fmt.Println("subscribing address...", address)
	return p.storage.Subscribe(address) == nil
}

func (p Parser) GetTransactions(address string) ([]Transaction, error) {
	if yes, err := p.storage.IsSubscribed(address); err != nil {
		fmt.Println("error when checking if subscribed for address", address, err)
		return nil, err
	} else if !*yes {
		fmt.Println("address not subscribed yet", address)
		return nil, nil
	}

	// first get all logs of an address
	logs, err := p.rpcClient.GetLogs(address)
	if err != nil {
		return nil, err
	}

	// then get all transactions by those logs
	return p.getAllTransactionsByLogs(*logs)
}

func (p Parser) getAllTransactionsByLogs(logs []Log) ([]Transaction, error) {
	var results []Transaction
	for _, log := range logs {
		if t, err := p.getTransactionByLog(log); err == nil {
			results = append(results, *t)
		} else {
			return nil, err
		}
	}

	fmt.Println("all logs processed", logs)
	return results, nil
}

func (p Parser) getTransactionByLog(log Log) (*Transaction, error) {
	// first query storage for any previously persisted transactions
	txn, err := p.storage.GetTransaction(log.Address, log.TransactionHash)
	if err != nil || txn != nil {
		return txn, err
	}

	// not available on storage, so query eth rpc for transaction
	result, err := p.rpcClient.GetTransactionByHash(log.TransactionHash)
	if err != nil {
		return nil, err
	}

	if err = p.storage.SaveTransaction(log.Address, *result); err != nil {
		fmt.Println("failed to persist transaction hash/address with error", result.Hash, log.Address, err)
	}

	return result, nil
}

package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Parser struct {
	storage   Storage
	rpcClient RpcClient
	worker    *Worker
}

func NewParser(s Storage) *Parser {
	newParser := &Parser{
		storage:   s,
		rpcClient: NewRpcClient(),
	}
	newParser.worker = NewWorker(newParser.autoSyncTransactions)
	return newParser
}

func (p Parser) initBlock() error {
	// initialize the latest block as the current block since app starts
	latestBlock, err := p.rpcClient.GetLatestBlock()
	if err != nil {
		return err
	}
	return p.storage.SaveBlockNumber(*latestBlock)
}

func (p Parser) GetCurrentBlock() (*int64, error) {
	hexStr, err := p.storage.GetSavedBlockNumber()
	if err != nil {
		return nil, err
	}

	return hexStringToInt(*hexStr)
}

func (p Parser) Subscribe(address string) bool {
	fmt.Println("subscribing address...", address)
	return p.storage.Subscribe(address) == nil
}

func (p Parser) GetTransactions(address string) (*[]Transaction, error) {
	if yes, err := p.storage.IsSubscribed(address); err != nil {
		fmt.Println("error when checking if subscribed for address", address, err)
		return nil, err
	} else if !*yes {
		fmt.Println("address not subscribed yet", address)
		return nil, nil
	}

	transactions, err := p.storage.GetTransactions(address)
	if err != nil {
		return nil, err
	}

	return &transactions, nil
}

// below are functions for auto update of transactions of subscribed addresses

func (p Parser) initBackgroundUpdate(period time.Duration) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	p.worker.RunBackground(ctx, period)
}

func (p Parser) autoSyncTransactions() error {
	// first get the latest block number
	printAutoSyncLog("START")
	latestBlockHexStr, err := p.rpcClient.GetLatestBlock()
	if err != nil {
		return err
	}

	latestBlockInt, err := hexStringToInt(*latestBlockHexStr)
	if err != nil {
		return err
	}

	// then the last saved block number
	fmt.Println("autoSyncTransactions...GetSavedBlockNumber")
	lastSavedBlockHexStr, err := p.storage.GetSavedBlockNumber()
	if err != nil {
		return err
	}
	lastSavedBlockInt, err := hexStringToInt(*lastSavedBlockHexStr)
	if err != nil {
		return err
	}

	// then get and save all transactions start/end blocks
	if err = p.updateTransactionsBetweenBlocks(*lastSavedBlockInt, *latestBlockInt); err != nil {
		printAutoSyncLog("FAILED")
		return err
	}
	printAutoSyncLog("SUCCEED")
	return nil
}

func printAutoSyncLog(prefix string) {
	fmt.Println(prefix+"-------autoSyncTransactions ts=", time.Now().String())
}

func (p Parser) updateTransactionsBetweenBlocks(startBlock, endBlock int64) error {
	fmt.Println("updateTransactionsBetweenBlocks start/end=", startBlock, endBlock)
	for blockId := startBlock + 1; blockId <= endBlock; blockId++ {
		if err := p.getAndSaveTransactionsByBlock(blockId); err != nil {
			return err
		}
	}
	return nil
}

func (p Parser) getAndSaveTransactionsByBlock(blockId int64) error {
	blockHex := fmt.Sprintf("0x%x", blockId)
	block, err := p.rpcClient.GetBlockByNumber(blockHex)
	if err != nil {
		return err
	}
	if block == nil {
		fmt.Println("getAndSaveTransactionsByBlock block is nil")
		return nil
	}

	fmt.Println("getAndSaveTransactionsByBlock size of block transactions=", len(block.Transactions))
	for _, t := range block.Transactions {
		// save the transaction according to if the from/to address is subscribed
		if err = p.saveAddressIfSubscribed(t.From, t); err != nil {
			return err
		}
		if err = p.saveAddressIfSubscribed(t.To, t); err != nil {
			return err
		}
	}
	// update the current block
	return p.storage.SaveBlockNumber(blockHex)
}

func (p Parser) saveAddressIfSubscribed(address string, transaction Transaction) error {
	fmt.Println("saveAddressIfSubscribed address/transaction hash=", address, transaction.Hash)
	isSubscribed, err := p.storage.IsSubscribed(address)
	if err != nil {
		return err
	}
	if !*isSubscribed {
		fmt.Println("saveAddressIfSubscribed address not subscribed", address)
		return nil
	}

	fmt.Println("saveAddressIfSubscribed SaveTransaction address/hash", address, transaction.Hash)
	return p.storage.SaveTransaction(address, transaction)
}

func hexStringToInt(hexStr string) (*int64, error) {
	s := strings.Replace(hexStr, "0x", "", 1)
	value, err := strconv.ParseInt(s, 16, 64)
	return &value, err
}

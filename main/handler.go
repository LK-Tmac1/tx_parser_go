package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var parser *Parser

func GetCurrentBlock(w http.ResponseWriter, r *http.Request) {
	currentBlock, err := parser.GetCurrentBlock()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := CurrentBlockResponse{CurrentBlock: *currentBlock}
	writeSucceedResponse(w, response)
}

func Subscribe(w http.ResponseWriter, r *http.Request) {
	var data SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := SubscribeResponse{Subscribed: parser.Subscribe(data.Address)}
	writeSucceedResponse(w, response)
}

func GetTransactions(w http.ResponseWriter, r *http.Request) {
	var data GetTransactionsRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	transactions, err := parser.GetTransactions(data.Address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	response := GetTransactionsResponse{Transactions: transactions}
	writeSucceedResponse(w, response)
}

func writeSucceedResponse(w http.ResponseWriter, response interface{}) {
	fmt.Println("writeSucceedResponse", response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Fatal("error when encoding json result", err)
	}
}

func initHandler(port string, p *Parser) {
	parser = p
	http.HandleFunc("/getCurrentBlock", GetCurrentBlock)
	http.HandleFunc("/subscribe", Subscribe)
	http.HandleFunc("/getTransactions", GetTransactions)

	log.Fatal(http.ListenAndServe("localhost:"+port, nil))
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var parser = NewParser()

func GetCurrentBlock(w http.ResponseWriter, r *http.Request) {
	currentBlock, err := parser.GetCurrentBlock()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := CurrentBlockResponse{CurrentBlock: *currentBlock}
	writeSucceedResponse(w, result)
}

func Subscribe(w http.ResponseWriter, r *http.Request) {
	var data SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parser.Subscribe(data.Address)
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

	writeSucceedResponse(w, transactions)
}

func writeSucceedResponse(w http.ResponseWriter, results interface{}) {
	fmt.Println("writeSucceedResponse", results)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}

func initHandler(port string) {
	http.HandleFunc("/getCurrentBlock", GetCurrentBlock)
	http.HandleFunc("/subscribe", Subscribe)
	http.HandleFunc("/getTransactions", GetTransactions)

	fmt.Println("starting service on", port)
	log.Fatal(http.ListenAndServe("localhost:"+port, nil))
}

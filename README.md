# Tx Parser

### How to run the program
```
cd main

# there should be a "main" bin created
go build

# a new terminal would be opened
open main

# now let's test all those endpoints, assuming with a valid address

curl -X GET http://localhost:8080/getCurrentBlock

curl -X POST http://localhost:8080/subscribe --data '{"address": "0xYour address"}'

# ensure to subscribe the address first
curl -X POST http://localhost:8080/getTransactions --data '{"address": "0xYour address"}'
```

### Design
The app consists of these components, by following the classic MVC model:
1. Entity: defines internal domain entities plus external request/response entities
2. Storage: abstracts the methods of persistence layer, extensible to other storages besides in-memory e.g. DB
3. Handler: forwards different requests to API endpoints then to specific controller methods
4. Controller: implements the specific business logic for those endpoints
![image](https://github.com/LK-Tmac1/tx_parser_go/assets/7871066/c731eca0-bd9d-474d-93be-47efd2e39320)


### Implementations
On a high level, the implementations of those endpoints are similar in that each will:
1. Interact with storage if necessary
2. Call ETH JSON RPC to get related data
3. Transform the data and return as a response

The most complicated one is ***GetTransactions***:
1. First get the logs of the input address by ***eth_getLogs***, with transactions returned, if any
2. For each log/transaction, read from storage if there is any persisted data available 
3. If not, use the transaction hash value to get the details by ***eth_getTransactionByHash***
4. Persist the transaction details to storage, if any

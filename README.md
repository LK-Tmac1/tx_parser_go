# Tx Parser

### How to run the program
```
cd main
# there should be a "tx_parser_go" bin created
go build
# a new terminal would be opened
open tx_parser_go

# now test all those endpoints, assuming with valid address(es)
curl -X GET http://localhost:8080/getCurrentBlock

curl -X POST http://localhost:8080/subscribe --data '{"address": "0xYour address"}'

# ensure to subscribe the address first
curl -X POST http://localhost:8080/getTransactions --data '{"address": "0xYour address"}'
```

### Design
The app consists of these components:
1. Entity: defines internal domain entities plus external request/response entities
2. RPC client: calls the external ETH cloudflare proxy to read data
3. Storage: abstracts the methods of persistence layer, extensible to other storages besides in-memory e.g. DB
4. Handler: forwards different requests to API endpoints then to specific controller methods
5. Controller: implements the specific business logic for those endpoints plus the transaction auto sync job logic
6. Worker: a background job that automatically update the transactions of subscribed addresses

![image](https://github.com/LK-Tmac1/tx_parser_go/assets/7871066/c8a2b4c9-a8df-4e28-aa5e-55311604d22d)


### Implementations
On a high level, the implementations of those endpoints are straightforward and similar in that each will simply:
1. Read data from storage
2. Transform the data and return as a response

The complicated one is the worker job that automatically sync the transactions i.e. ***autoSyncTransactions***:
1. First get the latest block info via RPC client
2. Then find the range of block IDs between the last saved block ID and this latest block ID
3. For each block ID within this range, find the block info with transaction details via ***eth_getBlockByNumber***
4. For each detailed transaction, if its address is subscribed before, persist its info to storage

### Misc. Items
I added many logging to better debug the app, although it could be verbose :)

The major TODO left is the make storage thread-safe, since we have a background job updating the data periodically.

However, considering the complexity, this is skipped. 

Also, I didn't find some good sample addresses, so some edge cases might not be covered.

Some sample outputs:
```
curl -X POST http://localhost:8080/subscribe --data '{"address": "0x388c818ca8b9251b393131c08a736a67ccb19297"}' 
{"subscribed":true}

curl -X POST http://localhost:8080/getTransactions --data '{"address": "0x388c818ca8b9251b393131c08a736a67ccb19297"}'
{
    "transactions":[
        {
            "blockHash":"0x83e4e72c86bdc82547d1e970b944f7e00511846979501c2e647fb7d760450a70",
            "blockNumber":"0x1127d0f",
            "from":"0xbaf6dc2e647aeb6f510f9e318856a1bcd66c5e19",
            "gas":"0x565f",
            "gasPrice":"0x680cfee39",
            "hash":"0x263a06a12ee90071e4a6968d76f21b9b57f23f3b542c29d17ac8d4ea214b7e91",
            "input":"0x",
            "nonce":"0xb549",
            "to":"0x388c818ca8b9251b393131c08a736a67ccb19297",
            "transactionIndex":"0x77",
            "value":"0xa4c2f6e62f05c7"
        }
    ]
}
```

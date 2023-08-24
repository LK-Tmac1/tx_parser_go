package main

// Storage The data access layer but right now we only have in memory DAL, can be extensible to e.g. DB
type Storage interface {
	GetTransaction(string, string) (*Transaction, error)
	SaveTransaction(string, Transaction) error
	IsSubscribed(string) (*bool, error)
	Subscribe(string) error
}

var (
	// local var for in memory storage
	subscriptions      = make(map[string]bool)
	transactionsRecord = make(map[string]*TransactionsRecord)
)

type InMemoryStore struct{}

func NewInMemoryStore() InMemoryStore {
	return InMemoryStore{}
}

func (ims InMemoryStore) GetTransaction(address, txnHash string) (*Transaction, error) {
	result, ok := transactionsRecord[address]
	if !ok {
		return nil, nil
	}
	target := result.Transactions[txnHash]
	return &target, nil
}

func (ims InMemoryStore) SaveTransaction(address string, t Transaction) error {
	if _, ok := transactionsRecord[address]; !ok {
		transactionsRecord[address] = &TransactionsRecord{
			Address:      address,
			Transactions: make(map[string]Transaction),
		}
	}

	transactionsRecord[address].Transactions[t.Hash] = t
	return nil
}

func (ims InMemoryStore) IsSubscribed(address string) (*bool, error) {
	flag := subscriptions[address]
	return &flag, nil
}

func (ims InMemoryStore) Subscribe(address string) error {
	subscriptions[address] = true
	return nil
}

package main

// Storage The data access layer but right now we only have in memory DAL, can be extensible to e.g. DB
type Storage interface {
	GetTransactionsRecord(string) (*TransactionsRecord, error)
	SaveTransactions(string, Transaction) error
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

func (ims InMemoryStore) GetTransactionsRecord(address string) (*TransactionsRecord, error) {
	if result, ok := transactionsRecord[address]; ok {
		return result, nil
	}
	return nil, nil
}

func (ims InMemoryStore) SaveTransactions(address string, t Transaction) error {
	if _, ok := transactionsRecord[address]; !ok {
		transactionsRecord[address] = &TransactionsRecord{
			Address:      address,
			Transactions: make([]Transaction, 0),
		}
	}

	record := transactionsRecord[address]
	record.Transactions = append(record.Transactions, t)
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

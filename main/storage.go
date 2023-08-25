package main

// Storage The data access layer but right now we only have in memory DAL, can be extensible to e.g. DB
type Storage interface {
	GetTransactions(string) ([]Transaction, error)
	SaveTransaction(string, Transaction) error

	IsSubscribed(string) (*bool, error)
	Subscribe(string) error

	SaveBlockNumber(string) error
	GetSavedBlockNumber() (*string, error)
}

var (
	// local var for in memory storage
	// TODO: since we rely on a background update, need to make them as syncMap/atomic.Int64
	subscriptions         = make(map[string]bool)
	transactionsByAddress = make(map[string][]Transaction)
	currentBlockNumber    = ""
)

type InMemoryStore struct{}

func NewInMemoryStore() InMemoryStore {
	return InMemoryStore{}
}

func (ims InMemoryStore) GetTransactions(address string) ([]Transaction, error) {
	return transactionsByAddress[address], nil
}

func (ims InMemoryStore) SaveTransaction(address string, t Transaction) error {
	transactionsByAddress[address] = append(transactionsByAddress[address], t)
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

func (ims InMemoryStore) SaveBlockNumber(newNum string) error {
	currentBlockNumber = newNum
	return nil
}
func (ims InMemoryStore) GetSavedBlockNumber() (*string, error) {
	return &currentBlockNumber, nil
}

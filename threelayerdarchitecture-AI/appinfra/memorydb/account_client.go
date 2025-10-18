package memorydb

import (
    "errors"
    "fmt"
    "sync"
    "sync/atomic"
)

var (
	ErrAccountNotFound    = errors.New("memorydb: account not found")
	ErrTransactionNotFound = errors.New("memorydb: transaction not found")
)

type AccountRecord struct {
	ID       string
	UserID   string
	Balance  int64
	Currency string
}

type TransactionRecord struct {
	ID               string
	FromAccountID    string
	ToAccountID      string
	Amount           int64
	Currency         string
	OccurredAtUnix   int64
	CorrelationID    string
	InitiatingUserID string
}

type AccountClient struct {
	mu       sync.RWMutex
	accounts map[string]AccountRecord
}

func NewAccountClient(seed []AccountRecord) *AccountClient {
	client := &AccountClient{
		accounts: make(map[string]AccountRecord),
	}
	for _, record := range seed {
		client.accounts[record.ID] = record
	}
	return client
}

func (c *AccountClient) Get(id string) (AccountRecord, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	record, ok := c.accounts[id]
	if !ok {
		return AccountRecord{}, ErrAccountNotFound
	}
	return record, nil
}

func (c *AccountClient) Put(record AccountRecord) error {
	if record.ID == "" {
		return errors.New("memorydb: account id required")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.accounts[record.ID] = record
	return nil
}

type TransactionClient struct {
	mu           sync.RWMutex
	transactions map[string]TransactionRecord
	idGenerator  IDGenerator
}

type IDGenerator interface {
	NewID(prefix string) string
}

type defaultIDGenerator struct{}

var idCounter uint64

func (defaultIDGenerator) NewID(prefix string) string {
    next := atomic.AddUint64(&idCounter, 1)
    return fmt.Sprintf("%s-%d", prefix, next)
}

func NewTransactionClient(generator IDGenerator) *TransactionClient {
	if generator == nil {
		generator = defaultIDGenerator{}
	}
	return &TransactionClient{
		transactions: make(map[string]TransactionRecord),
		idGenerator:  generator,
	}
}

func (c *TransactionClient) Put(record TransactionRecord) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if record.ID == "" {
		record.ID = c.idGenerator.NewID("txn")
	}
	c.transactions[record.ID] = record
	return record.ID, nil
}

func (c *TransactionClient) Get(id string) (TransactionRecord, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	record, ok := c.transactions[id]
	if !ok {
		return TransactionRecord{}, ErrTransactionNotFound
	}
	return record, nil
}

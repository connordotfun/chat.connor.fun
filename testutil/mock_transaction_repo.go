package testutil

import "github.com/aaronaaeng/chat.connor.fun/db"

type MockTransactionalRepository struct {
	MockRepository
}

func (r *MockTransactionalRepository) CreateTransaction() db.Transaction {
	return &MockTransaction{MockRepository: r.MockRepository}
}

type MockTransaction struct {
	MockRepository
}

func (r *MockTransaction) Commit() error {
	return nil
}

func (r *MockTransaction) Rollback() error {
	return nil
}

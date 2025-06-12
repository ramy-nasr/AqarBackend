package application

import (
	"context"
	"transaction-backend/domain"
)

type TransactionRepository interface {
	Save(ctx context.Context, txn domain.Transaction) error
	GetAll(ctx context.Context) ([]domain.Transaction, error)
}

type TransactionBroadcaster interface {
	Broadcast(txn domain.Transaction)
}

type TransactionService struct {
	repo        TransactionRepository
	broadcaster TransactionBroadcaster
}

func NewTransactionService(r TransactionRepository, b TransactionBroadcaster) *TransactionService {
	return &TransactionService{repo: r, broadcaster: b}
}

func (s *TransactionService) HandleNewTransaction(ctx context.Context, txn domain.Transaction) error {
	if err := s.repo.Save(ctx, txn); err != nil {
		return err
	}
	s.broadcaster.Broadcast(txn)
	return nil
}

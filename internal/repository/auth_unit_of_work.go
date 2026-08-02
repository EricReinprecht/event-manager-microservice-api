package repository

import (
	"context"

	"github.com/reinp/event-platform/backend/internal/database"
	emailVerificationRepository "github.com/reinp/event-platform/backend/internal/repository/email_verification"
	userRepository "github.com/reinp/event-platform/backend/internal/repository/user"
)

type AuthTransactionRepositories struct {
	Users         *userRepository.UserRepository
	Verifications *emailVerificationRepository.EmailVerificationRepository
}

type AuthUnitOfWork struct {
	transactionManager *database.TransactionManager
}

func NewAuthUnitOfWork(transactionManager *database.TransactionManager) *AuthUnitOfWork {
	return &AuthUnitOfWork{transactionManager: transactionManager}
}

func (u *AuthUnitOfWork) Transaction(
	ctx context.Context,
	fn func(*AuthTransactionRepositories) error,
) error {
	return u.transactionManager.Transaction(ctx, func(tx database.DBExecutor) error {
		return fn(&AuthTransactionRepositories{
			Users:         userRepository.NewUserRepository(tx),
			Verifications: emailVerificationRepository.NewEmailVerificationRepository(tx),
		})
	})
}

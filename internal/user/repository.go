package user

import (
	"context"

	"github.com/ttagiyeva/entain/internal/model"
)

//go:generate mockgen -source ./repository.go -mock_names Repository=MockUserRepository -package mocks -destination mocks/userRepository.mock.gen.go

// Repository is a repository for users
type Repository interface {
	GetUser(ctx context.Context, id string) (*model.UserDao, error)
	UpdateUserBalance(ctx context.Context, user *model.UserDao) error
}

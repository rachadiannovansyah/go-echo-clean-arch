package usecase

import (
	"context"
	"time"

	"github.com/rachadiannovansyah/go-echo-clean-arch/domain"
)

type userUsecase struct {
	userRepo       domain.UserRepository
	contextTimeout time.Duration
}

// NewUserUsecase will create new an UserUsecase object representation of domain.UserUsecase interface
func NewUserUsecase(a domain.UserRepository, timeout time.Duration) domain.UserUsecase {
	return &userUsecase{
		userRepo:       a,
		contextTimeout: timeout,
	}
}

func (a *userUsecase) Fetch(c context.Context, cursor string, num int64) (res []domain.User, nextCursor string, err error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, nextCursor, err = a.userRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}

	return
}

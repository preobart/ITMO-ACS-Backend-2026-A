package me

import (
	"context"

	"restaurant-booking/internal/shared/userrepo"
)

type Usecase struct {
	repo *userrepo.Repo
}

func NewUsecase(repo *userrepo.Repo) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) Me(ctx context.Context, input Input) (Output, error) {
	user, err := u.repo.GetByID(ctx, input.UserID)
	if err != nil {
		return Output{}, err
	}
	return Output{User: user}, nil
}

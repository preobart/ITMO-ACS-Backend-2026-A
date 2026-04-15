package list

import (
	"context"

	"restaurant-booking/internal/shared/bookingrepo"
)

type Usecase struct {
	repo *bookingrepo.Repo
}

func NewUsecase(repo *bookingrepo.Repo) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) List(ctx context.Context, input Input) (Output, error) {
	items, err := u.repo.ListByUser(ctx, input.UserID)
	if err != nil {
		return Output{}, err
	}
	return Output{Bookings: items}, nil
}

package get

import (
	"context"

	"github.com/google/uuid"

	"restaurant-booking/internal/domain"
	"restaurant-booking/internal/shared/bookingrepo"
)

type Usecase struct {
	repo *bookingrepo.Repo
}

func NewUsecase(repo *bookingrepo.Repo) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) Get(ctx context.Context, input Input) (Output, error) {
	bid, err := uuid.Parse(input.BookingID)
	if err != nil {
		return Output{}, domain.ErrInvalidInput
	}
	b, err := u.repo.GetByUserAndID(ctx, input.UserID, bid)
	if err != nil {
		return Output{}, err
	}
	return Output{Booking: b}, nil
}

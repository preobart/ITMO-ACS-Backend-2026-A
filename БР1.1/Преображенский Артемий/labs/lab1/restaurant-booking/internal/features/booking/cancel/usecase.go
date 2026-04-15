package cancel

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

func (u *Usecase) Cancel(ctx context.Context, input Input) (Output, error) {
	bid, err := uuid.Parse(input.BookingID)
	if err != nil {
		return Output{}, domain.ErrInvalidInput
	}
	if err := u.repo.Cancel(ctx, input.UserID, bid); err != nil {
		return Output{}, err
	}
	return Output{OK: true}, nil
}

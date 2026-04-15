package availability

import (
	"context"
	"strings"

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

func (u *Usecase) Check(ctx context.Context, input Input) (Output, error) {
	rid, err := uuid.Parse(input.RestaurantID)
	if err != nil {
		return Output{}, domain.ErrInvalidInput
	}
	tid, err := uuid.Parse(input.TableID)
	if err != nil {
		return Output{}, domain.ErrInvalidInput
	}
	d := strings.TrimSpace(input.BookingDate)
	st := strings.TrimSpace(input.StartTime)
	et := strings.TrimSpace(input.EndTime)
	if d == "" || st == "" || et == "" {
		return Output{}, domain.ErrInvalidInput
	}
	_, err = u.repo.TableForRestaurant(ctx, tid, rid)
	if err != nil {
		return Output{}, err
	}
	overlap, err := u.repo.HasOverlap(ctx, tid, d, st, et)
	if err != nil {
		return Output{}, err
	}
	return Output{Available: !overlap}, nil
}

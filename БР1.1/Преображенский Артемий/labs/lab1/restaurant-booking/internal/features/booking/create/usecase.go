package create

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

func (u *Usecase) Create(ctx context.Context, input Input) (Output, error) {
	d := strings.TrimSpace(input.BookingDate)
	st := strings.TrimSpace(input.StartTime)
	et := strings.TrimSpace(input.EndTime)
	if d == "" || st == "" || et == "" {
		return Output{}, domain.ErrInvalidInput
	}
	if input.GuestsCount <= 0 {
		return Output{}, domain.ErrInvalidInput
	}
	if input.RestaurantID == uuid.Nil || input.TableID == uuid.Nil {
		return Output{}, domain.ErrInvalidInput
	}
	t, err := u.repo.TableForRestaurant(ctx, input.TableID, input.RestaurantID)
	if err != nil {
		return Output{}, err
	}
	if t.Seats < input.GuestsCount {
		return Output{}, domain.ErrUnavailable
	}
	overlap, err := u.repo.HasOverlap(ctx, input.TableID, d, st, et)
	if err != nil {
		return Output{}, err
	}
	if overlap {
		return Output{}, domain.ErrUnavailable
	}
	b := domain.Booking{
		UserID:       input.UserID,
		RestaurantID: input.RestaurantID,
		TableID:      input.TableID,
		GuestsCount:  input.GuestsCount,
		BookingDate:  d,
		StartTime:    st,
		EndTime:      et,
	}
	created, err := u.repo.Create(ctx, b)
	if err != nil {
		return Output{}, err
	}
	return Output{Booking: created}, nil
}

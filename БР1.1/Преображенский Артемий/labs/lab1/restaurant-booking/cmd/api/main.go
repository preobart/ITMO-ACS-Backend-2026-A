package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"restaurant-booking/config"
	"restaurant-booking/internal/adapter/postgres"
	httpcontroller "restaurant-booking/internal/controller/http"
	"restaurant-booking/internal/features/auth/login"
	"restaurant-booking/internal/features/auth/register"
	bookingcancel "restaurant-booking/internal/features/booking/cancel"
	bookingcreate "restaurant-booking/internal/features/booking/create"
	bookingget "restaurant-booking/internal/features/booking/get"
	bookinglist "restaurant-booking/internal/features/booking/list"
	bookingavailability "restaurant-booking/internal/features/booking/availability"
	menulist "restaurant-booking/internal/features/menu/list"
	reviewcreate "restaurant-booking/internal/features/review/create"
	reviewlist "restaurant-booking/internal/features/review/list"
	restaurantdelete "restaurant-booking/internal/features/restaurant/delete"
	restaurantget "restaurant-booking/internal/features/restaurant/get"
	restaurantlist "restaurant-booking/internal/features/restaurant/list"
	tablelist "restaurant-booking/internal/features/table/list"
	userme "restaurant-booking/internal/features/user/me"
	"restaurant-booking/internal/shared/bookingrepo"
	"restaurant-booking/internal/shared/userrepo"
	"restaurant-booking/pkg/jwt"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	if err := AppRun(context.Background(), cfg); err != nil {
		panic(err)
	}
}

func AppRun(ctx context.Context, cfg config.Config) error {
	pgPool, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		return fmt.Errorf("postgres.New: %w", err)
	}

	jwtDur, err := time.ParseDuration(cfg.JWTExpires)
	if err != nil {
		jwtDur = 24 * time.Hour
	}
	jwtCfg := jwt.Config{
		Secret:  []byte(cfg.JWTSecret),
		Expires: jwtDur,
	}

	userRepo := userrepo.New(pgPool)
	bookingRepo := bookingrepo.New(pgPool)

	restaurantListRepo := restaurantlist.NewPostgres(pgPool)
	restaurantListUsecase := restaurantlist.NewUsecase(restaurantListRepo)

	restaurantGetRepo := restaurantget.NewPostgres(pgPool)
	restaurantGetUsecase := restaurantget.NewUsecase(restaurantGetRepo)

	restaurantDeleteRepo := restaurantdelete.NewPostgres(pgPool)
	restaurantDeleteUsecase := restaurantdelete.NewUsecase(restaurantDeleteRepo)

	menuListRepo := menulist.NewPostgres(pgPool)
	menuListUsecase := menulist.NewUsecase(menuListRepo)

	reviewListRepo := reviewlist.NewPostgres(pgPool)
	reviewListUsecase := reviewlist.NewUsecase(reviewListRepo)

	reviewCreateRepo := reviewcreate.NewPostgres(pgPool)
	reviewCreateUsecase := reviewcreate.NewUsecase(reviewCreateRepo)

	tableListRepo := tablelist.NewPostgres(pgPool)
	tableListUsecase := tablelist.NewUsecase(tableListRepo)

	bookingAvailabilityUsecase := bookingavailability.NewUsecase(bookingRepo)
	bookingCreateUsecase := bookingcreate.NewUsecase(bookingRepo)
	bookingListUsecase := bookinglist.NewUsecase(bookingRepo)
	bookingGetUsecase := bookingget.NewUsecase(bookingRepo)
	bookingCancelUsecase := bookingcancel.NewUsecase(bookingRepo)

	registerUsecase := register.NewUsecase(userRepo, jwtCfg)
	loginUsecase := login.NewUsecase(userRepo, jwtCfg)
	meUsecase := userme.NewUsecase(userRepo)

	routes := httpcontroller.Routes{
		JWT: jwtCfg,
		Auth: httpcontroller.AuthRoutes{
			Register: register.HTTP(registerUsecase),
			Login:    login.HTTP(loginUsecase),
		},
		Me: httpcontroller.MeRoutes{
			Profile:       userme.HTTP(meUsecase),
			BookingList:   bookinglist.HTTP(bookingListUsecase),
			BookingGet:    bookingget.HTTP(bookingGetUsecase),
			BookingCancel: bookingcancel.HTTP(bookingCancelUsecase),
		},
		Restaurants: httpcontroller.RestaurantRoutes{
			List:         restaurantlist.HTTP(restaurantListUsecase),
			Get:          restaurantget.HTTP(restaurantGetUsecase),
			Delete:       restaurantdelete.HTTP(restaurantDeleteUsecase),
			Menu:         menulist.HTTP(menuListUsecase),
			ReviewsList:  reviewlist.HTTP(reviewListUsecase),
			ReviewCreate: reviewcreate.HTTP(reviewCreateUsecase),
			Tables:       tablelist.HTTP(tableListUsecase),
			Availability: bookingavailability.HTTP(bookingAvailabilityUsecase),
		},
		Bookings: httpcontroller.BookingRoutes{
			Create: bookingcreate.HTTP(bookingCreateUsecase),
		},
	}

	router := httpcontroller.Router(routes)

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		fmt.Printf("shutdown signal received: %s\n", sig.String())
	case err := <-errCh:
		pgPool.Close()
		return fmt.Errorf("http server: %w", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		pgPool.Close()
		return fmt.Errorf("server shutdown: %w", err)
	}

	pgPool.Close()

	return nil
}

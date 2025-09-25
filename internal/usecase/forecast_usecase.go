package usecase

import (
	"context"
	"time"

	"ect-borrow-be/internal/dto"
	"ect-borrow-be/internal/repository"
)

type ForecastUsecase interface {
	Get(ctx context.Context, date time.Time, q string) (dto.ForecastData, error)
}

type forecastUsecase struct {
	repo repository.ForecastRepository
}

func NewForecastUsecase(repo repository.ForecastRepository) ForecastUsecase {
	return &forecastUsecase{repo: repo}
}

func (u *forecastUsecase) Get(ctx context.Context, date time.Time, q string) (dto.ForecastData, error) {
	return u.repo.Forecast(ctx, date, q)
}

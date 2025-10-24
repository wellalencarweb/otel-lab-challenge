package location

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/wellalencarweb/otel-lab-challenge/internal/entities"
	"github.com/wellalencarweb/otel-lab-challenge/internal/pkg/customerrors"
	"github.com/wellalencarweb/otel-lab-challenge/internal/pkg/httpclient"
)

type FindByZipCodeUseCaseInterface interface {
	Execute(ctx context.Context, zipCode string) (*entities.Location, error)
}

type FindByZipCodeUseCase struct {
	HttpClient httpclient.HttpClientInterface
	Logger     zerolog.Logger
}

func NewFindByZipCodeUseCase(
	httpClient httpclient.HttpClientInterface,
	logger zerolog.Logger,
) *FindByZipCodeUseCase {
	return &FindByZipCodeUseCase{
		HttpClient: httpClient,
		Logger:     logger,
	}
}

func (uc *FindByZipCodeUseCase) Execute(ctx context.Context, zipCode string) (*entities.Location, error) {
	var location entities.Location

	uc.Logger.Info().Msgf("[FindByZipCode] Calling API with zipcode [%s]", zipCode)

	if err := uc.HttpClient.Get(ctx, fmt.Sprintf("/%s/json/", zipCode), &location); err != nil {
		if *err.StatusCode == http.StatusNotFound {
			return nil, &customerrors.NotFoundError{
				Err:     err.Error,
				Message: "can not find zipcode",
				Tags: map[string]interface{}{
					"zipCode": zipCode,
				},
			}
		}

		return nil, &customerrors.UnknownError{
			Err:     err.Error,
			Message: "Unknown error getting location",
			Tags: map[string]interface{}{
				"zipCode": zipCode,
			},
		}
	}

	uc.Logger.Debug().Msgf("[FindByZipCode] Got location [%+v]", location)

	return &location, nil
}

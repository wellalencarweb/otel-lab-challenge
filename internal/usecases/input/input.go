package input

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/wellalencarweb/otel-lab-challenge/internal/entities/dto"
	"github.com/wellalencarweb/otel-lab-challenge/internal/pkg/customerrors"
	"github.com/wellalencarweb/otel-lab-challenge/internal/pkg/httpclient"
)

type InputUseCaseInterface interface {
	Execute(ctx context.Context, input dto.InputUCInput) (*dto.GetTemperaturesByZipCodeOutput, error)
}

type InputUseCase struct {
	HttpClient httpclient.HttpClientInterface
	Logger     zerolog.Logger
}

func NewInputUseCase(
	httpClient httpclient.HttpClientInterface,
	logger zerolog.Logger,
) *InputUseCase {
	return &InputUseCase{
		HttpClient: httpClient,
		Logger:     logger,
	}
}

func (uc *InputUseCase) Execute(ctx context.Context, input dto.InputUCInput) (*dto.GetTemperaturesByZipCodeOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	uc.Logger.Info().Msgf("[Input] Calling Orchestrator API with zipcode [%s]", input.Zipcode)

	var response dto.GetTemperaturesByZipCodeOutput

	if err := uc.HttpClient.Get(ctx, fmt.Sprintf("/?zipcode=%s", input.Zipcode), &response); err != nil {
		if *err.StatusCode == http.StatusNotFound {
			return nil, &customerrors.NotFoundError{
				Err:     err.Error,
				Message: "can not find zipcode",
				Tags: map[string]interface{}{
					"zipCode": input.Zipcode,
				},
			}
		}

		return nil, &customerrors.UnknownError{
			Err:     err.Error,
			Message: "Unknown error getting location",
			Tags: map[string]interface{}{
				"zipCode": input.Zipcode,
			},
		}
	}

	uc.Logger.Debug().Msgf("[Input] Got data: %+v", response)

	return &response, nil
}

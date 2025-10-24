package handlers

import (
	"errors"
	"net/http"
	"regexp"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/wellalencarweb/otel-lab-challenge/internal/entities/dto"
	"github.com/wellalencarweb/otel-lab-challenge/internal/pkg/responsehandler"
	"github.com/wellalencarweb/otel-lab-challenge/internal/usecases/climate"
	"github.com/wellalencarweb/otel-lab-challenge/internal/usecases/location"
)

type WebClimateHandlerInterface interface {
	GetTemperaturesByZipCode(w http.ResponseWriter, r *http.Request)
}

type WebClimateHandler struct {
	ResponseHandler              responsehandler.WebResponseHandlerInterface
	FindLocationByZipCodeUseCase location.FindByZipCodeUseCaseInterface
	FindClimateByCityNameUseCase climate.FindByCityNameUseCaseInterface
	Tracer                       trace.Tracer
}

func NewWebClimateHandler(
	rh responsehandler.WebResponseHandlerInterface,
	findByZipCodeUC location.FindByZipCodeUseCaseInterface,
	findByCityNameUC climate.FindByCityNameUseCaseInterface,
	tracer trace.Tracer,
) *WebClimateHandler {
	return &WebClimateHandler{
		ResponseHandler:              rh,
		FindLocationByZipCodeUseCase: findByZipCodeUC,
		FindClimateByCityNameUseCase: findByCityNameUC,
		Tracer:                       tracer,
	}
}

func (h *WebClimateHandler) GetTemperaturesByZipCode(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), carrier)
	ctx, span := h.Tracer.Start(ctx, "climate")
	defer span.End()

	qs := r.URL.Query()
	zipStr := qs.Get("zipcode")

	if err := validateInput(zipStr); err != nil {
		span.SetStatus(codes.Error, "invalid zipcode")
		span.RecordError(err)
		span.End()

		h.ResponseHandler.RespondWithError(w, http.StatusUnprocessableEntity, err)
		return
	}

	zipCodeCtx, zipCodeSpan := h.Tracer.Start(ctx, "find-location-by-zipcode")
	location, err := h.FindLocationByZipCodeUseCase.Execute(zipCodeCtx, zipStr)
	if err != nil {
		zipCodeSpan.SetStatus(codes.Error, "error finding location by zipcode")
		zipCodeSpan.RecordError(err)
		zipCodeSpan.End()

		h.ResponseHandler.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}
	if location.City == "" {
		zipCodeSpan.SetStatus(codes.Error, "zipcode not found")
		zipCodeSpan.RecordError(err)
		zipCodeSpan.End()

		h.ResponseHandler.RespondWithError(w, http.StatusNotFound, errors.New("zipcode not found"))
		return
	}

	zipCodeSpan.End()

	climateCtx, climateSpan := h.Tracer.Start(ctx, "find-climate-by-city-name")
	climate, err := h.FindClimateByCityNameUseCase.Execute(climateCtx, location.City)
	if err != nil {
		climateSpan.SetStatus(codes.Error, "error finding climate by city name")
		climateSpan.RecordError(err)
		climateSpan.End()

		h.ResponseHandler.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	climateSpan.End()

	fahrenheit, kelvin := convertTemperature(climate.Current.TempC)

	h.ResponseHandler.Respond(w, http.StatusOK, dto.GetTemperaturesByZipCodeOutput{
		City:       location.City,
		Celcius:    float32(climate.Current.TempC),
		Fahrenheit: float32(fahrenheit),
		Kelvin:     float32(kelvin),
	})
}

func validateInput(zipcode string) error {
	if zipcode == "" {
		return errors.New("invalid zipcode")
	}

	matched, err := regexp.MatchString(`\b\d{5}[\-]?\d{3}\b`, zipcode)
	if !matched || err != nil {
		return errors.New("invalid zipcode")
	}

	return nil
}

func convertTemperature(celcius float64) (float64, float64) {
	fahrenheit := celcius*1.8 + 32
	kelvin := celcius + 273.15

	return fahrenheit, kelvin
}

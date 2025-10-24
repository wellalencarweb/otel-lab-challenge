package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/wellalencarweb/otel-lab-challenge/internal/entities"
	"github.com/wellalencarweb/otel-lab-challenge/internal/pkg/mocks"
	"github.com/wellalencarweb/otel-lab-challenge/internal/pkg/responsehandler"
	"go.opentelemetry.io/otel"
)

type ClimateHandlerTestSuite struct {
	suite.Suite
	FindLocationByZipCodeUseCaseMock *mocks.FindByZipCodeUseCaseMock
	FindClimateByCityNameUseCaseMock *mocks.FindByCityNameUseCaseMock
	ResponseHandler                  *responsehandler.WebResponseHandler
	WebClimateHandler                *WebClimateHandler
}

func TestFindByZipCodeUseCase(t *testing.T) {
	suite.Run(t, new(ClimateHandlerTestSuite))
}

func (s *ClimateHandlerTestSuite) SetupTest() {
	findLocationByZipCodeUseCaseMock := new(mocks.FindByZipCodeUseCaseMock)
	findClimateByCityNameUseCaseMock := new(mocks.FindByCityNameUseCaseMock)
	responseHandler := responsehandler.NewWebResponseHandler()
	tracer := otel.Tracer("climate-test")

	s.FindLocationByZipCodeUseCaseMock = findLocationByZipCodeUseCaseMock
	s.FindClimateByCityNameUseCaseMock = findClimateByCityNameUseCaseMock
	s.ResponseHandler = responseHandler

	s.WebClimateHandler = NewWebClimateHandler(
		responseHandler,
		findLocationByZipCodeUseCaseMock,
		findClimateByCityNameUseCaseMock,
		tracer,
	)
}

func (s *ClimateHandlerTestSuite) clearMocks() {
	s.FindLocationByZipCodeUseCaseMock.ExpectedCalls = nil
	s.FindClimateByCityNameUseCaseMock.ExpectedCalls = nil
}

func (s *ClimateHandlerTestSuite) TestGetTemperaturesByZipCode() {
	s.Run("should return temperatures by zipcode", func() {
		defer s.clearMocks()

		zipCode := "22021001"
		city := "Rio de Janeiro"

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/?zipcode=%s", zipCode), nil)
		w := httptest.NewRecorder()

		expectedLocation := entities.Location{
			City:    city,
			Zipcode: zipCode,
		}

		expectedClimate := entities.Climate{
			Current: entities.ClimateData{
				TempC: 30,
			},
		}

		s.FindLocationByZipCodeUseCaseMock.On("Execute", zipCode).Return(&expectedLocation, nil)
		s.FindClimateByCityNameUseCaseMock.On("Execute", city).Return(&expectedClimate, nil)

		s.WebClimateHandler.GetTemperaturesByZipCode(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, _ := io.ReadAll(res.Body)
		expectedResponse := "{\"temp_C\":30,\"temp_F\":86,\"temp_K\":303.15}"

		s.Equal(http.StatusOK, res.StatusCode)
		s.Equal(expectedResponse, strings.TrimSuffix(string(data), "\n"))
	})

	s.Run("should return error when zipcode is empty", func() {
		defer s.clearMocks()

		zipCode := "011530000"
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/?zipcode=%s", zipCode), nil)
		w := httptest.NewRecorder()

		s.WebClimateHandler.GetTemperaturesByZipCode(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, _ := io.ReadAll(res.Body)
		expectedResponse := "{\"message\":\"invalid zipcode\"}"

		s.Equal(http.StatusUnprocessableEntity, res.StatusCode)
		s.Equal(expectedResponse, strings.TrimSuffix(string(data), "\n"))
	})

	s.Run("should return error when zipcode is invalid", func() {
		defer s.clearMocks()

		req := httptest.NewRequest(http.MethodGet, "/?zipcode=", nil)
		w := httptest.NewRecorder()

		s.WebClimateHandler.GetTemperaturesByZipCode(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, _ := io.ReadAll(res.Body)
		expectedResponse := "{\"message\":\"invalid zipcode\"}"

		s.Equal(http.StatusUnprocessableEntity, res.StatusCode)
		s.Equal(expectedResponse, strings.TrimSuffix(string(data), "\n"))
	})

	s.Run("should return error when zipcode is not found", func() {
		defer s.clearMocks()

		zipCode := "22021001"

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/?zipcode=%s", zipCode), nil)
		w := httptest.NewRecorder()

		expectedLocation := entities.Location{
			City:    "",
			Zipcode: zipCode,
		}

		s.FindLocationByZipCodeUseCaseMock.On("Execute", zipCode).Return(&expectedLocation, nil)

		s.WebClimateHandler.GetTemperaturesByZipCode(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, _ := io.ReadAll(res.Body)
		expectedResponse := "{\"message\":\"zipcode not found\"}"

		s.Equal(http.StatusNotFound, res.StatusCode)
		s.Equal(expectedResponse, strings.TrimSuffix(string(data), "\n"))
	})
}

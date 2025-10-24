package climate

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"

	"github.com/wellalencarweb/otel-lab-challenge/internal/entities"
	"github.com/wellalencarweb/otel-lab-challenge/internal/pkg/httpclient"
	"github.com/wellalencarweb/otel-lab-challenge/internal/pkg/mocks"
)

const API_KEY = "any-api-key"

type FindByCityNameUseCaseTestSuite struct {
	suite.Suite
	HttpClientMock        *mocks.HttpClientMock
	FindByCityNameUseCase *FindByCityNameUseCase
}

func TestFindByCityNameUseCase(t *testing.T) {
	suite.Run(t, new(FindByCityNameUseCaseTestSuite))
}

func (s *FindByCityNameUseCaseTestSuite) SetupTest() {
	httpClientMock := new(mocks.HttpClientMock)

	s.HttpClientMock = httpClientMock
	s.FindByCityNameUseCase = NewFindByCityNameUseCase(httpClientMock, zerolog.Nop(), API_KEY)
}

func (s *FindByCityNameUseCaseTestSuite) clearMocks() {
	s.HttpClientMock.ExpectedCalls = nil
}

func (s *FindByCityNameUseCaseTestSuite) TestFindByCityNameUseCase() {
	s.Run("should return climate", func() {
		defer s.clearMocks()

		ctx := context.Background()
		city := "Rio de Janeiro"
		endpoint := fmt.Sprintf("/v1/current.json?key=%s&q=%s&aqi=no", API_KEY, url.QueryEscape(city))

		s.HttpClientMock.On("Get", ctx, endpoint, &entities.Climate{}).Return(nil)

		result, err := s.FindByCityNameUseCase.Execute(ctx, city)

		s.Nil(err)
		s.NotNil(result)
	})

	s.Run("should return error when http client returns error", func() {
		defer s.clearMocks()

		ctx := context.Background()
		city := "Rio de Janeiro"
		endpoint := fmt.Sprintf("/v1/current.json?key=%s&q=%s&aqi=no", API_KEY, url.QueryEscape(city))

		s.HttpClientMock.On("Get", ctx, endpoint, &entities.Climate{}).Return(&httpclient.HttpClientError{
			Error: fmt.Errorf("any-error"),
		})

		result, err := s.FindByCityNameUseCase.Execute(ctx, city)

		s.Error(err)
		s.Nil(result)
	})
}

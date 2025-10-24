package location

import (
	"context"
	"fmt"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"

	"github.com/wellalencarweb/otel-lab-challenge/internal/entities"
	"github.com/wellalencarweb/otel-lab-challenge/internal/pkg/httpclient"
	"github.com/wellalencarweb/otel-lab-challenge/internal/pkg/mocks"
)

const API_KEY = "any-api-key"

type FindByZipCodeUseCaseTestSuite struct {
	suite.Suite
	HttpClientMock       *mocks.HttpClientMock
	FindByZipCodeUseCase *FindByZipCodeUseCase
}

func TestFindByZipCodeUseCase(t *testing.T) {
	suite.Run(t, new(FindByZipCodeUseCaseTestSuite))
}

func (s *FindByZipCodeUseCaseTestSuite) SetupTest() {
	httpClientMock := new(mocks.HttpClientMock)

	s.HttpClientMock = httpClientMock
	s.FindByZipCodeUseCase = NewFindByZipCodeUseCase(httpClientMock, zerolog.Nop())
}

func (s *FindByZipCodeUseCaseTestSuite) clearMocks() {
	s.HttpClientMock.ExpectedCalls = nil
}

func (s *FindByZipCodeUseCaseTestSuite) TestFindByZipCodeUseCase() {
	s.Run("should return location", func() {
		defer s.clearMocks()

		ctx := context.Background()
		zipCode := "22021-001"
		endpoint := fmt.Sprintf("/%s/json/", zipCode)

		s.HttpClientMock.On("Get", endpoint, &entities.Location{}).Return(nil)

		result, err := s.FindByZipCodeUseCase.Execute(ctx, zipCode)

		s.Nil(err)
		s.NotNil(result)
	})

	s.Run("should return error when http client returns error", func() {
		defer s.clearMocks()

		ctx := context.Background()
		zipCode := "22021-001"
		endpoint := fmt.Sprintf("/%s/json/", zipCode)

		s.HttpClientMock.On("Get", endpoint, &entities.Location{}).Return(&httpclient.HttpClientError{
			Error: fmt.Errorf("any-error"),
		})

		result, err := s.FindByZipCodeUseCase.Execute(ctx, zipCode)

		s.Error(err)
		s.Nil(result)
	})
}

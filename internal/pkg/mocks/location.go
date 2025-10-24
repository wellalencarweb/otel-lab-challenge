package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/wellalencarweb/otel-lab-challenge/internal/entities"
)

type FindByZipCodeUseCaseMock struct {
	mock.Mock
}

func (m *FindByZipCodeUseCaseMock) Execute(ctx context.Context, city string) (*entities.Location, error) {
	args := m.Called(ctx, city)
	return args.Get(0).(*entities.Location), args.Error(1)
}

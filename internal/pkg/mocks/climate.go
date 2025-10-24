package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/wellalencarweb/otel-lab-challenge/internal/entities"
)

type FindByCityNameUseCaseMock struct {
	mock.Mock
}

func (m *FindByCityNameUseCaseMock) Execute(ctx context.Context, city string) (*entities.Climate, error) {
	args := m.Called(ctx, city)
	return args.Get(0).(*entities.Climate), args.Error(1)
}

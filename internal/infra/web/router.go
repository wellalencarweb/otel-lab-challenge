package web

import (
	"net/http"

	"github.com/wellalencarweb/otel-lab-challenge/internal/infra/web/handlers"
)

type WebRouterInterface interface {
	Build() []RouteHandler
}

type InputWebRouter struct {
	WebInputHandler handlers.WebInputHandlerInterface
}

type OrchestratorWebRouter struct {
	WebClimateHandler handlers.WebClimateHandlerInterface
}

func NewInputWebRouter(webInputHandler handlers.WebInputHandlerInterface) *InputWebRouter {
	return &InputWebRouter{
		WebInputHandler: webInputHandler,
	}
}

func NewOrchestratorWebRouter(webClimateHandler handlers.WebClimateHandlerInterface) *OrchestratorWebRouter {
	return &OrchestratorWebRouter{
		WebClimateHandler: webClimateHandler,
	}
}

func (wr *InputWebRouter) Build() []RouteHandler {
	return []RouteHandler{
		{
			Path:        "/",
			Method:      http.MethodPost,
			HandlerFunc: wr.WebInputHandler.Handle,
		},
	}
}

func (wr *OrchestratorWebRouter) Build() []RouteHandler {
	return []RouteHandler{
		{
			Path:        "/",
			Method:      http.MethodGet,
			HandlerFunc: wr.WebClimateHandler.GetTemperaturesByZipCode,
		},
	}
}

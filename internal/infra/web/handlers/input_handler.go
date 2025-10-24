package handlers

import (
	"encoding/json"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/wellalencarweb/otel-lab-challenge/internal/entities/dto"
	"github.com/wellalencarweb/otel-lab-challenge/internal/pkg/customerrors"
	"github.com/wellalencarweb/otel-lab-challenge/internal/pkg/responsehandler"
	"github.com/wellalencarweb/otel-lab-challenge/internal/usecases/input"
)

type WebInputHandlerInterface interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type WebInputHandler struct {
	ResponseHandler responsehandler.WebResponseHandlerInterface
	InputUseCase    input.InputUseCaseInterface
	Tracer          trace.Tracer
}

func NewWebInputHandler(
	rh responsehandler.WebResponseHandlerInterface,
	inputUC input.InputUseCaseInterface,
	tracer trace.Tracer,
) *WebInputHandler {
	return &WebInputHandler{
		ResponseHandler: rh,
		InputUseCase:    inputUC,
		Tracer:          tracer,
	}
}

func (h *WebInputHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var dto dto.InputUCInput

	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, span := h.Tracer.Start(ctx, "climate")
	defer span.End()

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		recordSpan(span, err, "error decoding request body")

		h.ResponseHandler.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(r.Header))

	input, err := h.InputUseCase.Execute(ctx, dto)
	if err != nil {
		switch err.(type) {
		case *customerrors.NotFoundError:
			recordSpan(span, err, "not found")

			h.ResponseHandler.RespondWithError(w, http.StatusNotFound, err)
			return
		case *customerrors.ValidationError:
			recordSpan(span, err, "invalid zipcode")

			h.ResponseHandler.RespondWithError(w, http.StatusUnprocessableEntity, err)
			return
		case *customerrors.UnknownError:
		default:
			recordSpan(span, err, "error getting location")

			h.ResponseHandler.RespondWithError(w, http.StatusInternalServerError, err)
			return
		}
	}

	h.ResponseHandler.Respond(w, http.StatusOK, input)
}

func recordSpan(span trace.Span, err error, description string) {
	span.SetStatus(codes.Error, description)
	span.RecordError(err)
	span.End()
}

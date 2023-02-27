package handlers

import (
	"context"

	"go.opentelemetry.io/otel/trace"

	"github.com/labstack/echo/v4"
)

type HealthHandler struct {
}

type healthResponse struct {
	Status string `json:"status"`
}

func (h HealthHandler) Handler(ctx context.Context, tracer trace.Tracer) echo.HandlerFunc {
	return func(c echo.Context) error {
		if tracer != nil {
			_, span := tracer.Start(ctx, "health")
			defer span.End()
		}
		response := &healthResponse{
			Status: "UP",
		}
		return c.JSON(200, response)
	}
}

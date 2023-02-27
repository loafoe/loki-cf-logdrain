package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	syslog "github.com/RackSec/srslog"
	"go.opentelemetry.io/otel/trace"

	v2syslog "github.com/influxdata/go-syslog/v2"
	"github.com/influxdata/go-syslog/v2/rfc5424"
	"github.com/labstack/echo/v4"
)

type SyslogHandler struct {
	debug  bool
	token  string
	writer *syslog.Writer
	parser v2syslog.Machine
}

func NewSyslogHandler(token, promtailAddr string) (*SyslogHandler, error) {
	if token == "" {
		return nil, fmt.Errorf("Missing TOKEN value")
	}
	handler := &SyslogHandler{}
	handler.token = token

	parser := rfc5424.NewParser()

	if os.Getenv("DEBUG") == "true" {
		handler.debug = true
	}
	writer, err := syslog.Dial("tcp", promtailAddr,
		syslog.LOG_WARNING|syslog.LOG_DAEMON, "loki-cf-logdrain")
	if err != nil {
		return nil, fmt.Errorf("promtail: %w", err)
	}
	writer.SetFramer(syslog.RFC5425MessageLengthFramer)
	writer.SetFormatter(RFC5424PassThroughFormatter)
	handler.writer = writer
	handler.parser = parser
	return handler, nil
}

func RFC5424PassThroughFormatter(p syslog.Priority, hostname, tag, content string) string {
	return content
}

func (h *SyslogHandler) Handler(ctx context.Context, tracer trace.Tracer) echo.HandlerFunc {
	return func(c echo.Context) error {
		if tracer != nil {
			_, span := tracer.Start(ctx, "syslog")
			defer span.End()
		}
		t := c.Param("token")
		if h.token != t {
			return c.String(http.StatusUnauthorized, "")
		}
		b, _ := io.ReadAll(c.Request().Body)
		syslogMessage, err := h.parser.Parse(b)
		if err != nil {
			return err
		}
		fmt.Printf("version=%d\n", syslogMessage.Version())
		_, _ = h.writer.Write(b)
		return c.String(http.StatusOK, "")
	}
}

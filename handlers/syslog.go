package handlers

import (
	"fmt"
	"io/ioutil"
	"log/syslog"
	"net/http"
	"os"

	v2syslog "github.com/influxdata/go-syslog/v2"
	"github.com/influxdata/go-syslog/v2/rfc5424"
	"github.com/labstack/echo-contrib/zipkintracing"
	"github.com/openzipkin/zipkin-go"

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
		syslog.LOG_WARNING|syslog.LOG_DAEMON, "lokiproxy")
	if err != nil {
		return nil, err
	}
	handler.writer = writer
	handler.parser = parser
	return handler, nil
}

func (h *SyslogHandler) Handler(tracer *zipkin.Tracer) echo.HandlerFunc {
	return func(c echo.Context) error {
		if tracer != nil {
			defer zipkintracing.TraceFunc(c, "syslog_handler", zipkintracing.DefaultSpanTags, tracer)()
		}
		t := c.Param("token")
		if h.token != t {
			return c.String(http.StatusUnauthorized, "")
		}
		b, _ := ioutil.ReadAll(c.Request().Body)
		go func() {
			if tracer != nil {
				span := zipkintracing.StartChildSpan(c, "push", tracer)
				defer span.Finish()
				traceID := span.Context().TraceID.String()
				fmt.Printf("handler=syslog traceID=%s\n", traceID)
			}
			syslogMessage, err := h.parser.Parse(b)
			if err != nil {
				fmt.Printf("error parsing message: %v\n", err)
				return
			}
			fmt.Printf("version=%d\n", syslogMessage.Version())
			_, _ = h.writer.Write(b)
		}()
		return c.String(http.StatusOK, "")
	}
}

package logging

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

func NewStructuredLogger(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLogger{logger})
}

type StructuredLogger struct {
	Logger *logrus.Logger
}

func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &StructuredLoggerEntry{Logger: logrus.NewEntry(l.Logger)}
	logFields := logrus.Fields{}

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields["request_id"] = reqID
	}

	logFields["http_method"] = r.Method
	logFields["remote_address"] = r.RemoteAddr
	logFields["uri"] = r.RequestURI
	logFields["request_body"] = copyRequestBody(r)

	entry.Logger = entry.Logger.WithFields(logFields)

	entry.Logger.Info("Request started")

	return entry
}

func copyRequestBody(r *http.Request) string {
	body, _ := io.ReadAll(r.Body)
	_ = r.Body.Close()
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	return string(body)
}

type StructuredLoggerEntry struct {
	Logger logrus.FieldLogger
}

const responseRounding = 100

func (e *StructuredLoggerEntry) Write(status, bytes int, _ http.Header, elapsed time.Duration, _ interface{}) {
	e.Logger = e.Logger.WithFields(logrus.Fields{
		"resp_status":       status,
		"resp_bytes_length": bytes,
		"resp_elapsed":      elapsed.Round(time.Millisecond / responseRounding).String(),
	})
	e.Logger.Info("Request completed")
}

func (e *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	e.Logger = e.Logger.WithFields(logrus.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}

func GetLogEntry(r *http.Request) logrus.FieldLogger {
	entry, ok := middleware.GetLogEntry(r).(*StructuredLoggerEntry)
	if !ok {
		panic("LogEntry isn't *StructuredLoggerEntry")
	}

	return entry.Logger
}

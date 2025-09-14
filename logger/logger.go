package logger

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

var levelNames = map[string]Level{
	"DEBUG": LevelDebug,
	"INFO":  LevelInfo,
	"WARN":  LevelWarn,
	"ERROR": LevelError,
}

// Logger is a simple JSON logger with trace/span support.
type Logger struct {
	out   io.Writer
	level Level
	mu    sync.Mutex
}

// default logger instance
var defaultLogger = &Logger{out: os.Stdout, level: LevelInfo}

// Init configures the default logger using environment variables (via Viper).
// Supported variables:
// - LOG_LEVEL: DEBUG|INFO|WARN|ERROR (default INFO)
// - LOG_OUTPUT: console|file (default console)
// - LOG_FILE: path to log file when LOG_OUTPUT=file (default app.log)
func Init() error {
	v := viper.New()
	v.SetDefault("LOG_LEVEL", "INFO")
	v.SetDefault("LOG_OUTPUT", "console")
	v.SetDefault("LOG_FILE", "app.log")

	v.SetConfigFile(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	_ = v.ReadInConfig() // best-effort
	v.AutomaticEnv()

	lvl := strings.ToUpper(strings.TrimSpace(v.GetString("LOG_LEVEL")))
	if l, ok := levelNames[lvl]; ok {
		defaultLogger.level = l
	}

	out := strings.ToLower(strings.TrimSpace(v.GetString("LOG_OUTPUT")))
	if out == "file" {
		path := v.GetString("LOG_FILE")
		f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		defaultLogger.out = f
	} else {
		defaultLogger.out = os.Stdout
	}
	return nil
}

// Context keys
type ctxKey int

const (
	traceKey ctxKey = iota
	spanKey
)

// WithTrace returns a new context with the provided trace id.
func WithTrace(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceKey, traceID)
}

// WithSpan returns a new context with the provided span id.
func WithSpan(ctx context.Context, spanID string) context.Context {
	return context.WithValue(ctx, spanKey, spanID)
}

// StartSpan creates a new span, inheriting or creating a trace id.
// It returns the new context, trace id, and span id.
func StartSpan(ctx context.Context) (context.Context, string, string) {
	tr := TraceID(ctx)
	if tr == "" {
		tr = newID()
	}
	sp := newID()
	ctx = WithTrace(ctx, tr)
	ctx = WithSpan(ctx, sp)
	return ctx, tr, sp
}

// SetTraceSpan is a convenience to set both IDs.
func SetTraceSpan(ctx context.Context, traceID, spanID string) context.Context {
	ctx = WithTrace(ctx, traceID)
	ctx = WithSpan(ctx, spanID)
	return ctx
}

// TraceID extracts the trace id from context.
func TraceID(ctx context.Context) string {
	if v := ctx.Value(traceKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// SpanID extracts the span id from context.
func SpanID(ctx context.Context) string {
	if v := ctx.Value(spanKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// Log writes a JSON log line at the specified level.
func (l *Logger) Log(ctx context.Context, level Level, msg string, fields map[string]any) {
	if level < l.level {
		return
	}
	// Ensure we always have trace and span IDs
	tr := TraceID(ctx)
	sp := SpanID(ctx)
	if tr == "" {
		tr = newID()
	}
	if sp == "" {
		sp = newID()
	}
	entry := map[string]any{
		"time":     time.Now().UTC().Format(time.RFC3339Nano),
		"level":    levelString(level),
		"message":  msg,
		"trace_id": tr,
		"span_id":  sp,
	}
	// Merge fields, converting error values to strings to avoid memory addresses
	if fields != nil {
		for k, v := range fields {
			switch val := v.(type) {
			case error:
				entry[k] = val.Error()
			default:
				entry[k] = v
			}
		}
	}
	// Allow explicit override via fields
	if fields != nil {
		if tid, ok := fields["trace_id"].(string); ok && tid != "" {
			entry["trace_id"] = tid
		}
		if sid, ok := fields["span_id"].(string); ok && sid != "" {
			entry["span_id"] = sid
		}
	}

	b, _ := json.Marshal(entry)
	l.mu.Lock()
	defer l.mu.Unlock()
	_, _ = l.out.Write(append(b, '\n'))
}

func levelString(l Level) string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	default:
		return "ERROR"
	}
}

// helper to generate 16-byte hex id
func newID() string {
	var b [16]byte
	_, err := rand.Read(b[:])
	if err != nil {
		// fallback to timestamp-based id
		return strings.ReplaceAll(time.Now().Format("20060102T150405.000000000"), ".", "")
	}
	return hex.EncodeToString(b[:])
}

// Exported convenience wrappers using default logger
func Debug(ctx context.Context, msg string, fields map[string]any) {
	defaultLogger.Log(ctx, LevelDebug, msg, fields)
}
func Info(ctx context.Context, msg string, fields map[string]any) {
	defaultLogger.Log(ctx, LevelInfo, msg, fields)
}
func Warn(ctx context.Context, msg string, fields map[string]any) {
	defaultLogger.Log(ctx, LevelWarn, msg, fields)
}
func Error(ctx context.Context, msg string, fields map[string]any) {
	defaultLogger.Log(ctx, LevelError, msg, fields)
}

// Must is a helper to log an error and exit for fatal scenarios.
func Must(ctx context.Context, err error, msg string, fields map[string]any) {
	if err == nil {
		return
	}
	defaultLogger.Log(ctx, LevelError, msg, merge(fields, map[string]any{"error": err.Error()}))
	os.Exit(1)
}

func merge(a, b map[string]any) map[string]any {
	if a == nil && b == nil {
		return nil
	}
	out := map[string]any{}
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		out[k] = v
	}
	return out
}

// GinMiddleware attaches trace and span IDs to the request context and response headers.
func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ctx, tr, sp := StartSpan(ctx)
		c.Request = c.Request.WithContext(ctx)
		c.Header("X-Trace-ID", tr)
		c.Header("X-Span-ID", sp)
		c.Next()
	}
}

// InjectIntoRequest is a helper to inject current trace/span into an outbound HTTP request.
func InjectIntoRequest(ctx context.Context, req *http.Request) *http.Request {
	tr := TraceID(ctx)
	sp := SpanID(ctx)
	if tr == "" || sp == "" {
		ctx, tr, sp = StartSpan(ctx)
		_ = ctx // context is updated but caller holds original; tracing ids used below
	}
	req.Header.Set("X-Trace-ID", tr)
	req.Header.Set("X-Span-ID", sp)
	return req
}

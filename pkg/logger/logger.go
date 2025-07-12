package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// Color codes for different log levels
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorYellow = "\033[33m"
	ColorGreen  = "\033[32m"
	ColorBlue   = "\033[34m"
	ColorGray   = "\033[37m"
	ColorBold   = "\033[1m"
)

// PrettyJSONHandler wraps slog.JSONHandler to provide pretty-printed JSON output
type PrettyJSONHandler struct {
	handler slog.Handler
	writer  io.Writer
}

// NewPrettyJSONHandler creates a new PrettyJSONHandler
func NewPrettyJSONHandler(w io.Writer, opts *slog.HandlerOptions) *PrettyJSONHandler {
	return &PrettyJSONHandler{
		handler: slog.NewJSONHandler(w, opts),
		writer:  w,
	}
}

// colorizeLevel adds color to the log level based on its importance
func colorizeLevel(level string) string {
	switch level {
	case "DEBUG":
		return ColorGray + level + ColorReset
	case "INFO":
		return ColorGreen + level + ColorReset
	case "WARN":
		return ColorYellow + level + ColorReset
	case "ERROR":
		return ColorRed + ColorBold + level + ColorReset
	case "FATAL", "PANIC":
		return ColorRed + ColorBold + level + ColorReset
	default:
		return level
	}
}

// Handle formats the log record as pretty-printed JSON with colored output
func (h *PrettyJSONHandler) Handle(ctx context.Context, r slog.Record) error {
	// Create a map to hold the log data
	logData := make(map[string]interface{})

	// Add basic fields (without colors for JSON)
	logData["time"] = r.Time.Format("2006-01-02T15:04:05.000Z07:00")
	logData["level"] = r.Level.String()
	logData["msg"] = r.Message

	// Add attributes
	r.Attrs(func(a slog.Attr) bool {
		logData[a.Key] = a.Value.Any()
		return true
	})

	// Marshal with indentation
	jsonData, err := json.MarshalIndent(logData, "", "  ")
	if err != nil {
		return err
	}

	// Convert to string and apply colors
	jsonStr := string(jsonData)

	// Replace the level field with colored version
	levelPattern := fmt.Sprintf(`"level": "%s"`, r.Level.String())
	coloredLevel := fmt.Sprintf(`"level": "%s"`, colorizeLevel(r.Level.String()))
	jsonStr = strings.Replace(jsonStr, levelPattern, coloredLevel, 1)

	// Write to output
	_, err = h.writer.Write([]byte(jsonStr + "\n"))
	return err
}

// Enabled delegates to the underlying handler
func (h *PrettyJSONHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// WithAttrs delegates to the underlying handler
func (h *PrettyJSONHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyJSONHandler{
		handler: h.handler.WithAttrs(attrs),
		writer:  h.writer,
	}
}

// WithGroup delegates to the underlying handler
func (h *PrettyJSONHandler) WithGroup(name string) slog.Handler {
	return &PrettyJSONHandler{
		handler: h.handler.WithGroup(name),
		writer:  h.writer,
	}
}

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(NewPrettyJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	default:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	return log
}

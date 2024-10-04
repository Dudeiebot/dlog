package dlog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/fatih/color"
)

const (
	LevelTrace = slog.LevelDebug - 4
	LevelFatal = slog.LevelError + 4
)

type HandlerOptions struct {
	slog.HandlerOptions
	TimeStr string
}

type prettyHandler struct {
	slog.Handler
	opts *HandlerOptions
	w    io.Writer
}

func NewPrettyHandler(out io.Writer, opts *HandlerOptions) slog.Handler {
	h := slog.NewTextHandler(out, &opts.HandlerOptions)
	return &prettyHandler{h, opts, out}
}

func (h *prettyHandler) Handle(ctx context.Context, r slog.Record) error {
	timeStr := r.Time.Format(h.opts.TimeStr)
	level := r.Level.String()

	switch r.Level {
	case LevelTrace:
		level = color.CyanString("TRACE")
	case slog.LevelDebug:
		level = color.WhiteString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	case LevelFatal:
		level = color.MagentaString("FATAL")
	}

	var str strings.Builder
	str.WriteString(timeStr)
	str.WriteByte(' ')
	str.WriteString(level)
	str.WriteByte(' ')
	str.WriteString(r.Message)

	r.Attrs(func(a slog.Attr) bool {
		if a.Key != slog.LevelKey {
			str.WriteByte(' ')
			str.WriteString(a.Key)
			str.WriteByte('=')
			str.WriteString(fmt.Sprint(a.Value))
		}
		return true
	})

	_, err := fmt.Fprintln(h.w, str.String())
	return err
}

// called with dlog.NewLog() and then you vibe on
func NewLog() *slog.Logger {
	preHandler := NewPrettyHandler(os.Stdout, &HandlerOptions{
		HandlerOptions: slog.HandlerOptions{
			Level: LevelTrace,
		},
		TimeStr: "2006-01-02 15:04:05",
	})
	logger := slog.New(preHandler)
	slog.SetDefault(logger)
	return logger
}

// func main() {
// 	logger := NewLog()
// 	logger.Info("Starting server on :8080", "port", 8080, "status", "initializing")
//
// 	// Add a delay to demonstrate time difference
// 	time.Sleep(2 * time.Second)
// 	logger.Info("Server is now running", "port", 8080, "status", "running")
// 	logger.Info("This is an info message")
// 	logger.Warn("This is a warning message")
// 	logger.Error("This is an error message")
// 	logger.Debug("This is a Debug message")
// 	ctx := context.Background()
// 	logger.Log(ctx, LevelTrace, "This is a trace message") // This should be logged
// 	logger.Log(ctx, LevelFatal, "This is a fatal message") // This should be logged

// }

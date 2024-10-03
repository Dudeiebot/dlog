package dlog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
)

const (
	LevelTrace = slog.LevelDebug - 4
	LevelFatal = slog.LevelError + 4
)

var LevelNames = map[slog.Leveler]string{
	LevelTrace: "TRACE",
	LevelFatal: "FATAL",
}

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
	// var ok bool
	// var a slog.Attr
	//
	// level, _ := a.Value.Any().(slog.Level)
	// levelLabel, ok = LevelNames[level]
	// if !ok {
	// 	levelLabel = level.String()
	// }

	// if levelLabel == "" {
	// 	levelLabel = "UNKNOWN"
	// }

	// Format and print the log entry
	_, err := fmt.Fprintf(h.w, "%s  [%s]  %s", timeStr, level, r.Message)
	if err != nil {
		return err
	}

	// If you want to log attributes, uncomment the following:
	r.Attrs(func(a slog.Attr) bool {
		fmt.Fprintf(h.w, " %s=%v", a.Key, a.Value)
		return true
	})

	fmt.Fprintln(h.w)
	return nil
}

func NewLog() *slog.Logger {
	preHandler := NewPrettyHandler(os.Stdout, &HandlerOptions{
		HandlerOptions: slog.HandlerOptions{
			Level: LevelTrace,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.LevelKey {
					level := a.Value.Any().(slog.Level)
					levelLabel, exists := LevelNames[level]
					if !exists {
						levelLabel = level.String()
					}
					a.Value = slog.StringValue(levelLabel)
				}
				return a
			},
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
//
// 	logger.Info("This is an info message")
// 	logger.Warn("This is a warning message")
// 	logger.Error("This is an error message")
// 	ctx := context.Background()
// 	logger.Log(ctx, LevelTrace, "This is a trace message") // This should be logged
// 	logger.Log(ctx, LevelFatal, "This is a fatal message") // This should be logged
// }

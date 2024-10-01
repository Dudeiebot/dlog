package dlog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"
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
	currTime := time.Unix(0, time.Now().UnixNano())
	timeStr := currTime.Format(h.opts.TimeStr)
	level := r.Level.String()

	for len(level) < 5 {
		level += " "
	}
	fmt.Fprintf(h.w, "%s  [%s]  %s", timeStr, level, r.Message)

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
			Level: slog.LevelInfo,
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
// }

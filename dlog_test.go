package dlog

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/fatih/color"
)

func TestPrettyHandlerHandle(t *testing.T) {
	buf := new(bytes.Buffer)
	opts := &HandlerOptions{
		HandlerOptions: slog.HandlerOptions{
			Level: LevelTrace,
		},
		TimeStr: "2006-01-02 15:04:05",
	}
	handler := NewPrettyHandler(buf, opts)

	testTime := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)
	record := slog.Record{
		Time:    testTime,
		Level:   LevelFatal,
		Message: "Test",
	}
	record.AddAttrs(slog.String("key", "value"))

	err := handler.Handle(context.Background(), record)
	if err != nil {
		t.Fatalf("Handle returned an error: %v", err)
	}

	// Adjust the expected output to include the key=value
	expected := fmt.Sprintf(
		"2023-05-15 10:30:00 %s Test key=value\n",
		color.MagentaString("FATAL"),
	)
	if buf.String() != expected {
		t.Errorf("Output doesn't match expected.\nGot: %s\nWant: %s", buf.String(), expected)
	}
}

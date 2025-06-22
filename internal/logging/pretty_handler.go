package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
)

type prettyHandler struct {
	level slog.Leveler
	out   io.Writer
}

func (h *prettyHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *prettyHandler) Handle(_ context.Context, r slog.Record) error {
	fmt.Fprintf(h.out, "[%s] %s: %s\n", r.Time.Format("2006-01-02 15:04:05"), r.Level.String(), r.Message)
	return nil
}

func (h *prettyHandler) WithAttrs(_ []slog.Attr) slog.Handler { return h }
func (h *prettyHandler) WithGroup(_ string) slog.Handler      { return h }

package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/google/uuid"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

type keyType int

const key = keyType(0)

type logCtx struct {
	SongID    int
	RequestID uuid.UUID
}

type HandlerMiddleware struct {
	next slog.Handler
}

func NewHandlerMiddleware(next slog.Handler) *HandlerMiddleware {
	return &HandlerMiddleware{
		next: next,
	}
}

func (h *HandlerMiddleware) Enabled(ctx context.Context, rec slog.Level) bool {
	return h.next.Enabled(ctx, rec)
}

func (h *HandlerMiddleware) Handle(ctx context.Context, rec slog.Record) error {
	if c, ok := ctx.Value(key).(logCtx); ok {
		if c.SongID != 0 {
			rec.Add("songID", c.SongID)
		}
		if c.RequestID != uuid.Nil {
			rec.Add("requestID", c.RequestID)
		}
	}
	return h.next.Handle(ctx, rec)
}

func (h *HandlerMiddleware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &HandlerMiddleware{next: h.next.WithAttrs(attrs)}
}

func (h *HandlerMiddleware) WithGroup(name string) slog.Handler {
	return &HandlerMiddleware{next: h.next.WithGroup(name)}
}

func New(env string) *slog.Logger {
	var handler slog.Handler

	switch env {
	case envDev:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	case envProd:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	default:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}

	handler = NewHandlerMiddleware(handler)

	return slog.New(handler)
}

func WithSongID(ctx context.Context, id int) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.SongID = id
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{SongID: id})
}

func WithRequestID(ctx context.Context) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.RequestID = uuid.New()
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{RequestID: uuid.New()})
}

package db

import (
	"context"
	"time"
)

type Databaser interface {
	AddSegment(ctx context.Context, slug string) (int64, error)
	GetSegment(ctx context.Context, slug string) (*Segment, error)
	DeleteSegment(ctx context.Context, slug string) (int8, error)
	AddUserLink(ctx context.Context, id int64, slug string) (int64, error)
	DeleteUserLink(ctx context.Context, id int64, slug string) (int8, error)
	GetSegmentsByUserID(ctx context.Context, id int64) ([]string, error)
	GetHistory(ctx context.Context, id int64, dateFrom, dateTo time.Time) ([]HistoryAction, error)
}

type Segment struct {
	ID   string `db:"id"`
	Slug string `db:"slug"`
}

type HistoryAction struct {
	UserID string `db:"user_id"`
	Slug   string `db:"slug"`
	Action string `db:"action"`
	Date   string `db:"date"`
}

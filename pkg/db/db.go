package db

import (
	"context"
)

type Databaser interface {
	AddSegment(ctx context.Context, slug string) (int64, error)
	GetSegment(ctx context.Context, slug string) (*Segment, error)
	DeleteSegment(ctx context.Context, slug string) (int8, error)
	AddUserLink(ctx context.Context, id int64, slug string) (int64, error)
	DeleteUserLink(ctx context.Context, id int64, slug string) (int8, error)
	GetSegmentsByUserID(ctx context.Context, id int64) ([]string, error)
}

type Segment struct {
	ID   string `db:"id"`
	Slug string `db:"slug"`
}

package postgres

import (
	"context"
	"database/sql"

	_ "github.com/cockroachdb/cockroach-go/crdb"

	"github.com/7Maliko7/backend-trainee-assignment-2023/pkg/db"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) (*Repository, error) {
	return &Repository{
		db: db,
	}, nil
}

func (repo *Repository) Close() error {
	return repo.db.Close()
}

func (repo *Repository) AddSegment(ctx context.Context, slug string) (int64, error) {
	var id int64
	err := repo.db.QueryRowContext(ctx, AddSegmentQuery, slug).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (repo *Repository) GetSegment(ctx context.Context, slug string) (*db.Segment, error) {
	rows, err := repo.db.QueryContext(ctx, GetSegmentQuery, slug)
	if err != nil {
		return nil, err
	}

	var s *db.Segment
	for rows.Next() {
		s = &db.Segment{}
		err = rows.Scan(&s.ID, &s.Slug)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (repo *Repository) DeleteSegment(ctx context.Context, slug string) (int8, error) {
	var count int8
	err := repo.db.QueryRowContext(ctx, DeleteSegmentQuery, slug).Scan(&count)
	if err != nil {
		return count, err
	}

	return count, nil
}

func (repo *Repository) AddUserLink(ctx context.Context, id int64, slug string) (int64, error) {
	var linkId int64
	err := repo.db.QueryRowContext(ctx, AddUserLinkQuery, id, slug).Scan(&linkId)
	if err != nil {
		return linkId, err
	}

	return linkId, nil
}

func (repo *Repository) DeleteUserLink(ctx context.Context, id int64, slug string) (int8, error) {
	var count int8
	err := repo.db.QueryRowContext(ctx, DeleteUserLinkQuery, id, slug).Scan(&count)
	if err != nil {
		return count, err
	}

	return count, nil
}

func (repo *Repository) GetSegmentsByUserID(ctx context.Context, id int64) ([]string, error) {
	rows, err := repo.db.QueryContext(ctx, GetSegmentsByUserIDQuery, id)
	if err != nil {
		return nil, err
	}

	segments := make([]string, 0, 2)
	for rows.Next() {
		var s db.Segment
		err = rows.Scan(&s.Slug)
		if err != nil {
			return nil, err
		}
		segments = append(segments, s.Slug)
	}

	return segments, nil
}

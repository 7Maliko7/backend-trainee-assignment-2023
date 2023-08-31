package segments

import (
	"context"
	"fmt"

	svc "github.com/7Maliko7/backend-trainee-assignment-2023/internal/service"
	"github.com/7Maliko7/backend-trainee-assignment-2023/internal/transport/structs"
	"github.com/7Maliko7/backend-trainee-assignment-2023/pkg/db"
	"github.com/7Maliko7/backend-trainee-assignment-2023/pkg/errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
)

const (
	logKeyMethod = "method"
)

type Service struct {
	Repository db.Databaser
	Logger     log.Logger
}

func NewService(rep db.Databaser, logger log.Logger) svc.SegmentsService {
	return &Service{
		Repository: rep,
		Logger:     logger,
	}
}

func (s *Service) AddSegment(ctx context.Context, req structs.AddSegmentRequest) (*structs.AddSegmentResponse, error) {
	logger := log.With(s.Logger, logKeyMethod, "Create")

	segment, err := s.Repository.GetSegment(ctx, req.Slug)
	if err != nil {
		level.Error(logger).Log("repository", err.Error())
		return nil, errors.FailedRequest
	}

	if segment != nil {
		return nil, errors.SegmentExists
	}

	id, err := s.Repository.AddSegment(ctx, req.Slug)
	if err != nil {
		level.Error(logger).Log("repository", err.Error())
		return nil, errors.FailedRequest
	}
	level.Debug(logger).Log("New Segment", id)

	return &structs.AddSegmentResponse{Id: id}, nil
}

func (s *Service) DeleteSegment(ctx context.Context, req structs.DeleteSegmentRequest) (*structs.DeleteSegmentResponse, error) {
	logger := log.With(s.Logger, logKeyMethod, "Delete segment")

	count, err := s.Repository.DeleteSegment(ctx, req.Slug)
	if err != nil {
		level.Error(logger).Log("repository", err.Error())
		return nil, errors.FailedRequest
	}
	level.Debug(logger).Log("Deleted segment count", count)

	if count == 0 {
		return nil, errors.SegmentNotFound
	}

	return &structs.DeleteSegmentResponse{Response: structs.Response{Status: "success"}}, nil
}

func (s *Service) UpdateUserSegment(ctx context.Context, req structs.UpdateUserSegmentRequest) (*structs.UpdateUserSegmentResponse, error) {
	logger := log.With(s.Logger, logKeyMethod, "Update user segment")

	var failed []string
	for _, v := range req.Segments {
		switch v.Action {
		case "add":
			_, err := s.Repository.AddUserLink(ctx, req.UserId, v.Slug)
			if err != nil {
				level.Error(logger).Log("repository", err.Error())
				failed = append(failed, fmt.Sprintf("not added %v ", v.Slug))
			}
		case "delete":
			_, err := s.Repository.DeleteUserLink(ctx, req.UserId, v.Slug)
			if err != nil {
				level.Error(logger).Log("repository", err.Error())
				failed = append(failed, fmt.Sprintf("not deleted %v ", v.Slug))
			}
		default:
			failed = append(failed, fmt.Sprintf("invalid action %v ", v.Slug))
		}
	}

	if len(failed) > 0 {
		return nil, fmt.Errorf("%v", failed)
	}

	return &structs.UpdateUserSegmentResponse{Response: structs.Response{Status: "success"}}, nil
}

func (s *Service) GetSegments(ctx context.Context, req structs.GetSegmentsRequest) (*structs.GetSegmentsResponse, error) {
	logger := log.With(s.Logger, logKeyMethod, "Get Segments")
	segments, err := s.Repository.GetSegmentsByUserID(ctx, req.UserId)
	if err != nil {
		level.Error(logger).Log("repository", err.Error())
		return nil, errors.FailedRequest
	}
	level.Debug(logger).Log("Segments count", len(segments))

	if len(segments) == 0 {
		return nil, errors.UserNotFound
	}

	return &structs.GetSegmentsResponse{Segments: segments}, nil

}

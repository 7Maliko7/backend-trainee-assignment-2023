package middleware

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/7Maliko7/backend-trainee-assignment-2023/internal/service"
	"github.com/7Maliko7/backend-trainee-assignment-2023/internal/transport/structs"
)

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next service.SegmentsService) service.SegmentsService {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   service.SegmentsService
	logger log.Logger
}

func (mw loggingMiddleware) AddSegment(ctx context.Context, req structs.AddSegmentRequest) (*structs.AddSegmentResponse, error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "AddSegment", "id", "slug", "duration", time.Since(begin), "err")
	}(time.Now())
	return mw.next.AddSegment(ctx, req)
}

func (mw loggingMiddleware) DeleteSegment(ctx context.Context, req structs.DeleteSegmentRequest) (*structs.DeleteSegmentResponse, error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "DeleteSegment", "duration", time.Since(begin), "err")
	}(time.Now())
	return mw.next.DeleteSegment(ctx, req)
}

func (mw loggingMiddleware) UpdateUserSegment(ctx context.Context, req structs.UpdateUserSegmentRequest) (*structs.UpdateUserSegmentResponse, error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "UpdateUserSegment", "duration", time.Since(begin), "err")
	}(time.Now())
	return mw.next.UpdateUserSegment(ctx, req)
}

func (mw loggingMiddleware) GetSegments(ctx context.Context, req structs.GetSegmentsRequest) (*structs.GetSegmentsResponse, error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetSegments", "duration", time.Since(begin), "err")
	}(time.Now())
	return mw.next.GetSegments(ctx, req)
}

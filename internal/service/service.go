package service

import (
	"context"

	"github.com/7Maliko7/backend-trainee-assignment-2023/internal/transport/structs"
)

type SegmentsService interface {
	AddSegment(ctx context.Context, req structs.AddSegmentRequest) (*structs.AddSegmentResponse, error)
	DeleteSegment(ctx context.Context, req structs.DeleteSegmentRequest) (*structs.DeleteSegmentResponse, error)
	UpdateUserSegment(ctx context.Context, req structs.UpdateUserSegmentRequest) (*structs.UpdateUserSegmentResponse, error)
	GetSegments(ctx context.Context, req structs.GetSegmentsRequest) (*structs.GetSegmentsResponse, error)
}

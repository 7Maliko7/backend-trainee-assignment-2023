package errors

import "errors"

var (
	InvalidRequest = errors.New("invalid request")
	FailedRequest  = errors.New("failed request")
	DataNotFound   = errors.New("data not found")

	SegmentExists   = errors.New("segment already exists")
	SegmentNotFound = errors.New("segment not found")

	UserNotFound = errors.New("user not found")
)

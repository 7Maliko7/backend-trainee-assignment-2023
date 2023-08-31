package structs

type SegmentAction struct {
	Slug   string `json:"slug,omitempty"`
	Action string `json:"action,omitempty"`
}

type AddSegmentRequest struct {
	Slug string `json:"slug,omitempty"`
}

type AddSegmentResponse struct {
	Id int64 `json:"id,omitempty"`
}

type DeleteSegmentRequest struct {
	Slug string `json:"slug,omitempty"`
}

type DeleteSegmentResponse struct {
	Response
}

type Response struct {
	Status string `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

type UpdateUserSegmentRequest struct {
	UserId   int64           `json:"user_id,omitempty"`
	Segments []SegmentAction `json:"segments,omitempty"`
}

type UpdateUserSegmentResponse struct {
	Response
}

type GetSegmentsRequest struct {
	UserId int64 `json:"user_id,omitempty"`
}

type GetSegmentsResponse struct {
	Segments []string `json:"segments,omitempty"`
}

type GetUserSegmentHistoryRequest struct {
	UserId int64  `json:"user_id,omitempty"`
	Period string `json:"period,omitempty"`
}

type GetUserSegmentHistoryResponse struct {
	Actions [][]string
}

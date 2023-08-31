package http

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"github.com/7Maliko7/backend-trainee-assignment-2023/internal/config"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/7Maliko7/backend-trainee-assignment-2023/internal/transport"
	"github.com/7Maliko7/backend-trainee-assignment-2023/internal/transport/structs"
	statusErr "github.com/7Maliko7/backend-trainee-assignment-2023/pkg/errors"
)

// NewService wires Go kit endpoints to the HTTP transport.
func NewService(
	svcEndpoints transport.Endpoints, options []kithttp.ServerOption, logger log.Logger, appConfig *config.Config,
) http.Handler {
	var (
		r            = mux.NewRouter()
		errorLogger  = kithttp.ServerErrorLogger(logger)
		errorEncoder = kithttp.ServerErrorEncoder(encodeErrorResponse)
	)
	options = append(options, errorLogger, errorEncoder)

	segmentRouter := mux.NewRouter()
	segmentRoute := segmentRouter.PathPrefix("/api/v1/segment").Subrouter()

	segmentRoute.Methods(http.MethodPost).Path("/").Handler(kithttp.NewServer(
		svcEndpoints.AddSegment,
		decodeAddSegmentRequest,
		encodeResponse,
		options...,
	))

	segmentRoute.Methods(http.MethodDelete).Path("/{slug}").Handler(kithttp.NewServer(
		svcEndpoints.DeleteSegment,
		decodeDeleteSegmentRequest,
		encodeResponse,
		options...,
	))

	userRouter := mux.NewRouter()
	userRoute := userRouter.PathPrefix("/api/v1/user").Subrouter()

	userRoute.Methods(http.MethodPost).Path("/{user_id:[0-9]+}/segments").Handler(kithttp.NewServer(
		svcEndpoints.UpdateUserSegment,
		decodehUserSegmentRequest,
		encodeResponse,
		options...,
	))

	userRoute.Methods(http.MethodGet).Path("/{user_id:[0-9]+}/segments").Handler(kithttp.NewServer(
		svcEndpoints.GetSegments,
		decodeGetSegmentsRequest,
		encodeResponse,
		options...,
	))

	userRoute.Methods(http.MethodGet).Path("/{user_id:[0-9]+}/segments/history/{period:[0-9]{4}-[0-9]{2}}").Handler(kithttp.NewServer(
		svcEndpoints.GetUserSegmentHistory,
		decodeGetUserSegmentHistoryRequest,
		encodeGetUserSegmentHistoryResponse,
		options...,
	))

	initDocs(appConfig, "/api/v1")

	r.Handle("/api/v1/segment/{_dummy:.*}", segmentRouter)
	r.Handle("/api/v1/user/{_dummy:.*}", userRouter)
	r.Methods("GET").PathPrefix("/docs").Handler(httpSwagger.WrapHandler)

	return r
}

func decodeAddSegmentRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req structs.AddSegmentRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, statusErr.InvalidRequest
	}

	if req.Slug == "" {
		return nil, statusErr.InvalidRequest
	}

	return req, nil
}

func decodeDeleteSegmentRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	req := structs.DeleteSegmentRequest{
		Slug: vars["slug"],
	}

	return req, nil
}

func decodehUserSegmentRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req structs.UpdateUserSegmentRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, statusErr.InvalidRequest
	}

	if len(req.Segments) == 0 {
		return nil, statusErr.InvalidRequest
	}

	vars := mux.Vars(r)
	req.UserId, err = strconv.ParseInt(vars["user_id"], 10, 64)
	if err != nil {
		return nil, statusErr.InvalidRequest
	}

	return req, nil
}

func decodeGetSegmentsRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req structs.GetSegmentsRequest
	vars := mux.Vars(r)
	req.UserId, err = strconv.ParseInt(vars["user_id"], 10, 64)
	if err != nil {
		return nil, statusErr.InvalidRequest
	}

	return req, nil
}

func decodeGetUserSegmentHistoryRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req structs.GetUserSegmentHistoryRequest
	vars := mux.Vars(r)
	req.UserId, err = strconv.ParseInt(vars["user_id"], 10, 64)
	if err != nil {
		return nil, statusErr.InvalidRequest
	}
	req.Period = vars["period"]

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeErrorResponse(ctx, e.error(), w)

		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	return json.NewEncoder(w).Encode(response)
}

func encodeGetUserSegmentHistoryResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeErrorResponse(ctx, e.error(), w)

		return nil
	}

	resp, _ := response.(*structs.GetUserSegmentHistoryResponse)

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment;filename=report.csv")

	return csv.NewWriter(w).WriteAll(resp.Actions)
}

type errorer interface {
	error() error
}

func encodeErrorResponse(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case statusErr.SegmentNotFound,
		statusErr.UserNotFound,
		statusErr.DataNotFound:
		return http.StatusNotFound
	case statusErr.InvalidRequest,
		statusErr.SegmentExists:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

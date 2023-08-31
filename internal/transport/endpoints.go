package transport

import (
	"context"
	"github.com/go-kit/kit/endpoint"

	"github.com/7Maliko7/backend-trainee-assignment-2023/internal/service"
	"github.com/7Maliko7/backend-trainee-assignment-2023/internal/transport/structs"
)

type Endpoints struct {
	AddSegment            endpoint.Endpoint
	DeleteSegment         endpoint.Endpoint
	UpdateUserSegment     endpoint.Endpoint
	GetSegments           endpoint.Endpoint
	GetUserSegmentHistory endpoint.Endpoint
}

func MakeEndpoints(s service.SegmentsService) Endpoints {
	return Endpoints{
		AddSegment:            makeAddSegmentEndpoint(s),
		DeleteSegment:         makeDeleteSegmentEndpoint(s),
		UpdateUserSegment:     makeUpdateUserSegmentEndpoint(s),
		GetSegments:           makeGetSegmentsEndpoint(s),
		GetUserSegmentHistory: makeGetUserSegmentHistoryEndpoint(s),
	}
}

// /api/v1/segment/
//
//	@Summary		Метод создания сегмента.
//	@Description	Метод создания сегмента. Принимает slug (название) сегмента.
//	@Tags			segment
//	@Accept			json
//	@Produce		json
//	@Param			slug	body	string	true	"Название сегмента"	example({"slug": "AVITO_DISCOUNT_30"})
//	@Success		200		{string}	int
//	@Router			/segment/ [POST]
func makeAddSegmentEndpoint(s service.SegmentsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(structs.AddSegmentRequest)
		response, err := s.AddSegment(ctx, req)
		if err != nil {
			return nil, err
		}
		return response, nil
	}
}

// /api/v1/segment/{slug}
//
//	@Summary		Метод удаления сегмента.
//	@Description	Метод удаления сегмента. Принимает slug (название) сегмента.
//	@Tags			segment
//	@Produce		json
//	@Param			slug	path	string	true	"Название сегмента"	example(AVITO_DISCOUNT_30)
//	@Success		200		{string}	int
//	@Router			/segment/{slug} [DELETE]
func makeDeleteSegmentEndpoint(s service.SegmentsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(structs.DeleteSegmentRequest)
		response, err := s.DeleteSegment(ctx, req)
		if err != nil {
			return nil, err
		}
		return response, nil
	}
}

// /api/v1/user/{user_id}/segments
//
//	@Summary		Метод добавления/удаления пользователя в сегмент.
//	@Description	Метод добавления/удаления пользователя в сегмент. Принимает список slug (названий) сегментов которые нужно добавить пользователю, список slug (названий) сегментов которые нужно удалить у пользователя, id пользователя.
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			user_id	path	int	true	"ID пользователя"	example(1000)
//	@Param			segments	body	[]structs.SegmentAction	true	"Список сегментов"	example({ "segments": [ { "slug": "AVITO_DISCOUNT_30", "action": "add" }, { "slug": "AVITO_PERFORMANCE_VA", "action": "delete" } ] })
//	@Success		200		{string}	string
//	@Router			/user/{user_id}/segments [POST]
func makeUpdateUserSegmentEndpoint(s service.SegmentsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(structs.UpdateUserSegmentRequest)
		response, err := s.UpdateUserSegment(ctx, req)
		if err != nil {
			return nil, err
		}
		return response, nil
	}
}

// /api/v1/user/{user_id}/segments
//
//	@Summary		Метод получения активных сегментов пользователя.
//	@Description	Метод получения активных сегментов пользователя. Принимает на вход id пользователя.
//	@Tags			user
//	@Produce		json
//	@Param			user_id	path	int	true	"ID пользователя"	example(1000)
//	@Success		200		{string}	[]string
//	@Router			/user/{user_id}/segments [GET]
func makeGetSegmentsEndpoint(s service.SegmentsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(structs.GetSegmentsRequest)
		response, err := s.GetSegments(ctx, req)
		if err != nil {
			return nil, err
		}
		return response, nil
	}
}

// /api/v1/user/{user_id}/segments/history/{period}
//
//	@Summary		Метод получения истории изменения сегментов пользователя.
//	@Description	Метод получения истории изменения сегментов пользователя. Принимает на вход id пользователя и период год-месяц.
//	@Tags			user
//	@Produce		text/csv
//	@Param			user_id	path	int	true	"ID пользователя"	example(1000)
//	@Param			period	path	string	true	"период год-месяц"	example(2023-08)
//	@Success		200		{file}	file
//	@Router			/user/{user_id}/segments/history/{period} [GET]
func makeGetUserSegmentHistoryEndpoint(s service.SegmentsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(structs.GetUserSegmentHistoryRequest)
		response, err := s.GetUserSegmentHistory(ctx, req)
		if err != nil {
			return nil, err
		}
		return response, nil
	}
}

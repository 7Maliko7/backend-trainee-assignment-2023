package middleware

import "github.com/7Maliko7/backend-trainee-assignment-2023/internal/service"

// Middleware describes a service middleware.
type Middleware func(service service.SegmentsService) service.SegmentsService

package http

import (
	"github.com/7Maliko7/backend-trainee-assignment-2023/docs"
	"github.com/7Maliko7/backend-trainee-assignment-2023/internal/config"
)

func initDocs(appConfig *config.Config, basePath string) {
	docs.SwaggerInfo.Title = "Segment Service API"
	docs.SwaggerInfo.Description = "This is implementation of backend trainee assignment 2023."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = appConfig.ListenAddress
	docs.SwaggerInfo.BasePath = basePath
	docs.SwaggerInfo.Schemes = []string{"http"}
}

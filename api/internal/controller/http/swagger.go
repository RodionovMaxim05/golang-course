package http

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

// SwaggerHandler serves the interactive Swagger UI for the API, backed
// by generated OpenAPI documentation.
func SwaggerHandler() http.Handler {
	return httpSwagger.WrapHandler
}

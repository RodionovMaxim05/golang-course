package http

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

func SwaggerHandler() http.Handler {
	return httpSwagger.WrapHandler
}

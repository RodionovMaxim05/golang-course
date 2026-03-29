package controllers

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "gateway/docs"
)

func SwaggerHandler() http.Handler {
	return httpSwagger.WrapHandler
}

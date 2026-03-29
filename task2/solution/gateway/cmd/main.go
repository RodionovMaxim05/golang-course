// @title       GitHub Repo API
// @version     1.0
// @description Gateway to information about GitHub repositories
// @host        localhost:8080
// @BasePath    /
package main

import (
	_ "gateway/docs"
	"gateway/internal/adapters"
	"gateway/internal/controllers"
	"gateway/internal/services"
	"gateway/internal/usecases"
	"log"
	"net/http"
	"os"
)

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func main() {
	collectorAddr := getEnv("COLLECTOR_ADDR", "localhost:50051")
	port := getEnv("GATEWAY_PORT", ":8080")

	grpcClient, err := adapters.NewClient(collectorAddr)
	if err != nil {
		log.Fatalf("failed to connect to collector: %v", err)
	}

	repoService := services.NewRepoService(grpcClient)
	getRepoUsecase := usecases.NewGetRepoUsecase(repoService)
	repoHandler := controllers.NewRepoHandler(getRepoUsecase)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/repo", repoHandler.GetRepo)
	mux.Handle("/swagger/", controllers.SwaggerHandler())

	log.Printf("Gateway HTTP server listening on %s", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

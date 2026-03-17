package main

import (
	"collector/internal/adapters/clients"
	"collector/internal/adapters/controllers"
	"collector/internal/services"
	"collector/internal/usecases"
	"log"
	"net"
	"os"

	pb "api/gen"

	"google.golang.org/grpc"
)

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func main() {
	port := getEnv("COLLECTOR_PORT", ":50051")

	githubClient := clients.NewGitHubClient()
	repoService := services.NewRepoService(githubClient)
	getRepoUsecase := usecases.NewRepoUsecase(repoService)
	grpcHandler := controllers.NewRepoHandler(getRepoUsecase)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterRepoServiceServer(server, grpcHandler)

	log.Println("Collector gRPC server listening on :50051")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

package main

import (
	"github.com/Mellanie-Marques/microservices/shipping/internal/adapter/grpc"
	"github.com/Mellanie-Marques/microservices/shipping/internal/application/core/api"
	"github.com/Mellanie-Marques/microservices/shipping/internal/config"
)

func main() {
	application := api.NewApplication()
	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}

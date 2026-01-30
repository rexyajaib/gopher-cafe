package main

import (
	"github.com/rexyajaib/gopher-cafe/internal/worker"
	"os"
)

func main() {
	// Get Default Seeds
	seedRequest := "./seed/seed1.json"
	grpcServerAddress := "localhost:50051"

	requestedSeed, exists := os.LookupEnv("SEED_REQUEST_FILE")
	if exists && requestedSeed != "" {
		seedRequest = requestedSeed
	}

	grpcServerAddressEnv, exists := os.LookupEnv("GRPC_SERVER")
	if exists && grpcServerAddressEnv != "" {
		grpcServerAddress = grpcServerAddressEnv
	}

	data, err := os.ReadFile(seedRequest)
	if err != nil {
		panic(err)
	}

	seed, err := worker.ParseSeedRequest(data)
	if err != nil {
		panic(err)
	}

	taskWorker := worker.NewTaskWorker(seed, grpcServerAddress)

	taskWorker.Start()
}

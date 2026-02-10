package worker

import (
	"context"
	"github.com/rexyajaib/gopher-cafe/internal/domain/model"
	gophercafepb "github.com/rexyajaib/gopher-cafe/pkg/gen/go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"math/rand"
	"time"
)

func NewTaskWorker(seedRequest model.Seed, address string) *TaskWorker {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	gopherCafeClient := gophercafepb.NewGopherCafeServiceClient(conn)

	mapOfRequests := make(map[int]*gophercafepb.ExecuteBrewRequest)
	for i, request := range seedRequest.Requests {
		executeBrewRequest := buildExecuteBrewRequest(request)
		mapOfRequests[i] = &executeBrewRequest
	}

	return &TaskWorker{
		seed:                    seedRequest,
		gopherCafeServiceClient: gopherCafeClient,
		mapOfRequests:           mapOfRequests,
	}
}

func (tw *TaskWorker) Start() {
	totalRequests := 0
	tw.initCtx = context.Background()

	duration := time.Millisecond * time.Duration(tw.seed.RunningDurationMs)
	totalTimer := time.NewTimer(duration)
	defer totalTimer.Stop()

loop:
	for {
		select {
		case <-totalTimer.C:
			break loop
		default:
			go tw.worker()
			totalRequests++
			sleepMs := rand.Intn(tw.seed.PauseDurationMs) // 0..N
			time.Sleep(time.Duration(sleepMs) * time.Millisecond)
		}
	}

	// GetStats RPC call to get stats
	response, err := tw.gopherCafeServiceClient.GetStats(tw.initCtx, &gophercafepb.GetStatsRequest{})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Gopher Cafe Stats: %v", model.GetStatsResponse{
		TotalRequestProcessed:     response.GetTotalRequestProcessed(),
		P90ProcessingMilliseconds: response.GetP90ProcessingMilliseconds(),
	})

	log.Printf("Gopher Cafe Total Requests Processed: %d", totalRequests)

	log.Printf("task worker stopped")
}

func (tw *TaskWorker) worker() {
	requestChoosen := rand.Intn(len(tw.mapOfRequests))
	request := tw.mapOfRequests[requestChoosen]

	// Call RPC to gopher cafe service to execute brew
	executeResponse, err := tw.gopherCafeServiceClient.ExecuteBrew(tw.initCtx, request)
	if err != nil {
		log.Fatal(err)
	}

	if tw.seed.EnableLogging {
		log.Printf("Brew Execution Response: %v", model.ExecuteBrewResponseFromProto(executeResponse))
	}
}

func buildExecuteBrewRequest(request model.SeedRequest) gophercafepb.ExecuteBrewRequest {
	var orders []*gophercafepb.Order
	for _, order := range request.Orders {
		orders = append(orders, &gophercafepb.Order{
			Id:    int64(order.Id),
			Drink: buildDrinkType(order.Drink),
		})
	}

	return gophercafepb.ExecuteBrewRequest{
		Baristas: int32(request.Baristas),
		Orders:   orders,
	}
}

func buildDrinkType(drink string) gophercafepb.DrinkType {
	switch drink {
	case "Espresso":
		return gophercafepb.DrinkType_DRINK_TYPE_ESPRESSO
	case "Latte":
		return gophercafepb.DrinkType_DRINK_TYPE_LATTE
	case "Frappe":
		return gophercafepb.DrinkType_DRINK_TYPE_FRAPPE
	case "Matcha":
		return gophercafepb.DrinkType_DRINK_TYPE_MATCHA
	default:
		return gophercafepb.DrinkType_DRINK_TYPE_UNSPECIFIED
	}
}

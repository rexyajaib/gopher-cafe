package worker

import (
	"context"
	gophercafepb "github.com/rexyajaib/gopher-cafe/pkg/gen/go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func NewTaskWorker(seedRequest SeedRequest, address string) *TaskWorker {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	gopherCafeClient := gophercafepb.NewGopherCafeServiceClient(conn)

	return &TaskWorker{
		SeedRequest:             seedRequest,
		gopherCafeServiceClient: gopherCafeClient,
	}
}

func (tw *TaskWorker) Start() {
	for i := 0; i < tw.SeedRequest.LoopCount; i++ {
		ctx := context.Background()
		request := buildExecuteBrewRequest(tw.SeedRequest)

		//  Call RPC to gopher cafe service to execute brew
		executeResponse, err := tw.gopherCafeServiceClient.ExecuteBrew(ctx, &request)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Brew Execution Response: %v", ExecuteBrewResponseFromProto(executeResponse))

		// GetStats RPC call to get stats
		response, err := tw.gopherCafeServiceClient.GetStats(ctx, &gophercafepb.GetStatsRequest{})
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Gopher Cafe Stats: %v", GetStatsResponse{
			TotalRequestProcessed:       response.GetTotalRequestProcessed(),
			TotalProcessingMilliseconds: response.GetTotalProcessingMilliseconds(),
		})

		time.Sleep(time.Millisecond * time.Duration(tw.SeedRequest.DelayMs))
	}
}

func buildExecuteBrewRequest(request SeedRequest) gophercafepb.ExecuteBrewRequest {
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

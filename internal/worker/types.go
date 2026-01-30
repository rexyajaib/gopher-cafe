package worker

import (
	"context"
	"github.com/rexyajaib/gopher-cafe/internal/domain/model"
	gophercafepb "github.com/rexyajaib/gopher-cafe/pkg/gen/go/v1"
)

type TaskWorker struct {
	initCtx                 context.Context
	gopherCafeServiceClient gophercafepb.GopherCafeServiceClient
	seed                    model.Seed
	mapOfRequests           map[int]*gophercafepb.ExecuteBrewRequest
}

package worker

import (
	"encoding/json"
	gophercafepb "github.com/rexyajaib/gopher-cafe/pkg/gen/go/v1"
)

type TaskWorker struct {
	gopherCafeServiceClient gophercafepb.GopherCafeServiceClient
	SeedRequest             SeedRequest
}

type SeedRequest struct {
	Baristas int     `json:"baristas"`
	Orders   []Order `json:"orders"`
}

func ParseSeedRequest(data []byte) (SeedRequest, error) {
	var seedRequest SeedRequest
	err := json.Unmarshal(data, &seedRequest)
	if err != nil {
		return SeedRequest{}, err
	}
	return seedRequest, nil
}

type Order struct {
	Id    int    `json:"id"`
	Drink string `json:"drink"`
}

type GetStatsResponse struct {
	TotalRequestProcessed       int `json:"total_request_processed"`
	TotalProcessingMilliseconds int `json:"total_processing_milliseconds"`
}

type ExecuteBrewResponse struct {
	Results []ExecuteBrewResult
}

func ExecuteBrewResponseFromProto(protoResp *gophercafepb.ExecuteBrewResponse) ExecuteBrewResponse {
	executeBrewResult := make([]ExecuteBrewResult, len(protoResp.Results))
	for i, res := range protoResp.Results {
		executeBrewResult[i] = ExecuteBrewResult{
			OrderId: res.OrderId,
			Steps:   buildStepsFromProto(res.GetSteps()),
		}
	}

	return ExecuteBrewResponse{
		Results: executeBrewResult,
	}
}

type ExecuteBrewResult struct {
	OrderId int64
	Steps   []Step
}

type Step struct {
	EquipmentType string
	Start         int64
	End           int64
}

func buildStepsFromProto(steps []*gophercafepb.Step) []Step {
	result := make([]Step, len(steps))
	for i, step := range steps {
		result[i] = Step{
			EquipmentType: step.GetEquipment().String(),
			Start:         step.Start,
			End:           step.End,
		}
	}
	return result
}

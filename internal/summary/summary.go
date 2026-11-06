package summary

import (
	"fmt"
	"sort"
	"time"

	"github.com/fredericobenevides/pge-stress-test/internal/stresstest"
)

type Summary struct {
	totalOperationDuration time.Duration
	results                []stresstest.Result
}

type SummaryData struct {
	totalRequests          int
	totalErrors            int
	statusCodes            map[int]int
	totalOperationDuration time.Duration
	attemptThroughput      float64
	responseThroughput     float64
	successThroughput      float64
	minLatency             time.Duration
	maxLatency             time.Duration
	averageLatency         time.Duration
}

func NewSummary(operationDuration time.Duration, results []stresstest.Result) *Summary {
	return &Summary{
		totalOperationDuration: operationDuration,
		results:                results,
	}
}

func (s *Summary) Calculate() SummaryData {
	totalRequests := 0
	totalErrors := 0
	statusCodes := make(map[int]int)
	total2xxRequests := 0

	var totalLatency time.Duration
	var minLatency time.Duration
	var maxLatency time.Duration

	for _, result := range s.results {
		totalRequests++
		if result.Error != nil {
			totalErrors++
		} else {
			statusCodes[result.StatusCode]++
		}

		if result.StatusCode >= 200 && result.StatusCode < 300 {
			total2xxRequests++
		}

		if minLatency == 0 || result.OperationTime < minLatency {
			minLatency = result.OperationTime
		}

		if result.OperationTime > maxLatency {
			maxLatency = result.OperationTime
		}

		totalLatency += result.OperationTime
	}

	return SummaryData{
		totalRequests:          totalRequests,
		totalErrors:            totalErrors,
		statusCodes:            statusCodes,
		totalOperationDuration: s.totalOperationDuration,
		attemptThroughput:      float64(totalRequests) / s.totalOperationDuration.Seconds(),
		responseThroughput:     float64(totalRequests-totalErrors) / s.totalOperationDuration.Seconds(),
		successThroughput:      float64(total2xxRequests) / s.totalOperationDuration.Seconds(),
		minLatency:             minLatency,
		maxLatency:             maxLatency,
		averageLatency:         totalLatency / time.Duration(totalRequests),
	}
}

func (s *Summary) Print() {
	data := s.Calculate()

	responseCount := data.totalRequests - data.totalErrors

	fmt.Println("════════════════════ Stress Test Summary ════════════════════")
	fmt.Printf("Duration:              %s\n", data.totalOperationDuration)
	fmt.Printf("Requests attempted:    %d\n", data.totalRequests)
	fmt.Printf("HTTP responses:        %d\n", responseCount)
	fmt.Printf("Transport errors:      %d\n", data.totalErrors)

	fmt.Println()
	fmt.Println("Throughput")
	fmt.Printf("  Attempted:           %.2f req/s\n", data.attemptThroughput)
	fmt.Printf("  Responses:           %.2f req/s\n", data.responseThroughput)
	fmt.Printf("  Successful (2xx):    %.2f req/s\n", data.successThroughput)

	fmt.Println()
	fmt.Println("Latency")
	fmt.Printf("  Average:             %s\n", data.averageLatency)
	fmt.Printf("  Minimum:             %s\n", data.minLatency)
	fmt.Printf("  Maximum:             %s\n", data.maxLatency)

	fmt.Println()
	fmt.Println("HTTP status codes")
	fmt.Println("  Status   Requests   Rate")
	fmt.Println("  ------   --------   ------")

	statusCodes := make([]int, 0, len(data.statusCodes))
	for statusCode := range data.statusCodes {
		statusCodes = append(statusCodes, statusCode)
	}

	sort.Ints(statusCodes)

	for _, statusCode := range statusCodes {
		count := data.statusCodes[statusCode]

		var rate float64
		if responseCount > 0 {
			rate = 100 * float64(count) / float64(responseCount)
		}

		fmt.Printf("  %-8d %-10d %.2f%%\n", statusCode, count, rate)
	}

	fmt.Println("══════════════════════════════════════════════════════════════")
}

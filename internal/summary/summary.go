package summary

import (
	"fmt"
	"time"

	"github.com/fredericobenevides/pge-stress-test/internal/stresstest"
)

type Summary struct {
	results []stresstest.Result
}

type SummaryData struct {
	totalOperation time.Duration
	totalRequests  int
	totalErrors    int
	statusCodes    map[int]int
}

func NewSummary(results []stresstest.Result) *Summary {
	return &Summary{
		results: results,
	}
}

func (s *Summary) Calculate() SummaryData {
	totalOperation := time.Duration(0)
	totalRequests := 0
	totalErrors := 0
	statusCodes := make(map[int]int)

	for _, result := range s.results {
		totalRequests++
		totalOperation += result.OperationTime
		if result.Error != nil {
			totalErrors++
		} else {
			statusCodes[result.StatusCode]++
		}
	}

	return SummaryData{
		totalRequests:  totalRequests,
		totalOperation: totalOperation,
		totalErrors:    totalErrors,
		statusCodes:    statusCodes,
	}
}

func (s *Summary) Print() {
	summaryData := s.Calculate()

	println("================================")
	println("Total Operation Time:", summaryData.totalOperation.String())
	println("Total Requests:", summaryData.totalRequests)
	println("Total Errors:", summaryData.totalErrors)
	println("Status Codes Received:")
	println("\t Status | Total of Requests ")
	for code, count := range summaryData.statusCodes {
		fmt.Printf("\t    %d | %d  \n", code, count)
	}
	println("================================")
}

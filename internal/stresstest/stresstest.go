package stresstest

import (
	"io"
	"net/http"
	"sync"
	"time"
)

type Result struct {
	OperationTime time.Duration
	StatusCode    int
	Error         error
}

type StressTest struct {
	config *Config
}

func NewStressTest(config *Config) *StressTest {
	return &StressTest{
		config: config,
	}
}

func (s *StressTest) Run() (time.Duration, []Result) {
	startTime := time.Now()

	results := make([]Result, 0)
	resultsCh := make(chan Result, s.config.requests)

	jobs := make(chan struct{})

	var wg sync.WaitGroup

	s.worker(&wg, jobs, resultsCh)

	go func() {
		s.sendRequestsToWorker(jobs)
		wg.Wait()
		close(resultsCh)
	}()

	for result := range resultsCh {
		results = append(results, result)
	}

	operationDuration := time.Since(startTime)

	return operationDuration, results
}

func (s *StressTest) worker(wg *sync.WaitGroup, jobs <-chan struct{}, resultsCh chan<- Result) {
	for i := 0; i < s.config.concurrency; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for range jobs {
				result := s.makeRequest()
				resultsCh <- result
			}
		}()
	}
}

func (s *StressTest) sendRequestsToWorker(jobs chan<- struct{}) {
	for i := 0; i < s.config.requests; i++ {
		jobs <- struct{}{}
	}

	close(jobs)
}

func (s *StressTest) makeRequest() Result {
	latency := time.Now()

	resp, err := http.Get(s.config.url)
	if err != nil {
		return Result{
			OperationTime: time.Since(latency),
			StatusCode:    0,
			Error:         err,
		}
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	return Result{
		OperationTime: time.Since(latency),
		StatusCode:    resp.StatusCode,
		Error:         nil,
	}
}

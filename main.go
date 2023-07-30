package main

import (
	"fmt"
	"harry-pap/beat_assignment/calculator"
	"harry-pap/beat_assignment/concurrency"
	"harry-pap/beat_assignment/model"
	"harry-pap/beat_assignment/parser"
	"os"
	"sync"
	"time"
)

const numberOfWorkers = 10

// Runs the script
// Is responsible for launching the involved goroutines,
// opening the involved files and wiring the needed functions
func main() {
	now := time.Now().UTC()

	var wg sync.WaitGroup

	jobs := make(chan []calculator.RidePart, 100)
	results := make(chan model.RideFareEstimation, numberOfWorkers*20)
	done := make(chan interface{}, numberOfWorkers)

	channelCloserInput := concurrency.ChannelCloserInput{Count: numberOfWorkers, Done: done, Results: results, Wg: &wg}

	launch(func() { concurrency.CloseResultChannelWhenWorkersDone(channelCloserInput) }, &wg)

	for w := 1; w <= 10; w++ {
		workerInput := concurrency.WorkerInput{Jobs: jobs, Results: results, Done: done, Wg: &wg, Fun: calculator.CalculateFareForRide}
		launch(func() { concurrency.RunWorker(workerInput) }, &wg)
	}

	inputFile, inputErr := os.Open(os.Args[1])
	outputFile, outputErr := os.Create(os.Args[2])

	panicIfNotNil(inputErr)
	panicIfNotNil(outputErr)

	defer outputFile.Close()
	defer inputFile.Close()

	launch(func() { concurrency.ResultWriter(outputFile, results, &wg) }, &wg)

	parser.ParseInputCSV(inputFile, jobs)

	close(jobs)

	wg.Wait()

	fmt.Println("Time elapsed: ", time.Since(now))
}

func panicIfNotNil(err error) {
	if err != nil {
		panic(err)
	}
}
func launch(fun func(), wg *sync.WaitGroup) {
	wg.Add(1)
	go fun()
}

package concurrency

import (
	"encoding/csv"
	"fmt"
	"github.com/spf13/afero"
	"harry-pap/beat_assignment/calculator"
	"harry-pap/beat_assignment/model"
	"log"
	"sync"
)

// WorkerInput contains the input of the RunWorker function
type WorkerInput struct {
	Jobs    chan []calculator.RidePart
	Results chan model.RideFareEstimation
	Done    chan interface{}
	Wg      *sync.WaitGroup
	Fun     func([]calculator.RidePart) (model.RideFareEstimation, error)
}

// ChannelCloserInput contains the input of CloseResultChannelWhenWorkersDone
type ChannelCloserInput struct {
	Count   int
	Done    chan interface{}
	Results chan model.RideFareEstimation
	Wg      *sync.WaitGroup
}

// RunWorker reads the ChannelCloserInput.Jobs channel
// by invoking the ChannelCloserInput.Fun function on each message
// Upon completion, a message to WorkerInput.Done channel is sent, and WorkerInput.Wg.Done() is invoked
func RunWorker(workerInput WorkerInput) {
	defer workerInput.Wg.Done()

	for job := range workerInput.Jobs {
		result, err := workerInput.Fun(job)

		if err != nil {
			fmt.Println("Failed to calculate job because of error:", err)
		} else {
			workerInput.Results <- result
		}
	}

	workerInput.Done <- nil
}

// ResultWriter reads the RideFareEstimation channel, and writing each estimate into a *afero.File
// Upon completion sync.WaitGroup.Done() is invoked
func ResultWriter(file *afero.File, inputs chan model.RideFareEstimation, wg *sync.WaitGroup) {
	defer wg.Done()

	writer := csv.NewWriter(*file)

	for input := range inputs {
		err := writer.Write(input.ToStringSlice())

		if err != nil {
			log.Fatal("Cannot Write to file:", err)
		}
	}
	writer.Flush()
}

// CloseResultChannelWhenWorkersDone listens to the ChannelCloserInput.Done channel, and for closing ChannelCloserInput.Done and
// ChannelCloserInput.Results channels, when ChannelCloserInput.Count messages are received from ChannelCloserInput.Done
func CloseResultChannelWhenWorkersDone(input ChannelCloserInput) {
	for range input.Done {
		input.Count--
		if input.Count == 0 {
			close(input.Done)
			close(input.Results)
		}
	}

	input.Wg.Done()
}

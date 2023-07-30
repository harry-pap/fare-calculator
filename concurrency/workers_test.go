package concurrency

import (
	"bytes"
	"harry-pap/beat_assignment/calculator"
	"harry-pap/beat_assignment/model"
	"sync"
	"testing"
	"time"
)

func newGroup() *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)
	return &wg
}

var (
	sampleFareEstimation1 = model.RideFareEstimation{RideID: 1, CostEstimation: 14.51}
	sampleFareEstimation2 = model.RideFareEstimation{RideID: 2, CostEstimation: 45.12}
	sampleFareEstimation3 = model.RideFareEstimation{RideID: 3, CostEstimation: 25.56}
	fares                 = []model.RideFareEstimation{sampleFareEstimation1, sampleFareEstimation2, sampleFareEstimation3}
)

func TestWorker(t *testing.T) {

	type args struct {
		input          WorkerInput
		expectedPushes []model.RideFareEstimation
	}
	tests := []struct {
		name string
		args args
		run  func(args)
	}{
		{
			"for a single ride",
			args{WorkerInput{
				make(chan []calculator.RidePart, 10),
				make(chan model.RideFareEstimation, 10),
				make(chan interface{}, 10),
				newGroup(),
				func(parts []calculator.RidePart) (model.RideFareEstimation, error) {
					return sampleFareEstimation1, nil
				}},
				[]model.RideFareEstimation{sampleFareEstimation1},
			},
			func(args args) {
				args.input.Jobs <- []calculator.RidePart{}
				close(args.input.Jobs)
			},
		},
		{
			"for multiple rides",
			args{WorkerInput{
				make(chan []calculator.RidePart, 10),
				make(chan model.RideFareEstimation, 10),
				make(chan interface{}, 10),
				newGroup(),
				func(parts []calculator.RidePart) (model.RideFareEstimation, error) {
					return fares[parts[0].RideID], nil
				}},
				[]model.RideFareEstimation{sampleFareEstimation1, sampleFareEstimation2, sampleFareEstimation3},
			},
			func(args args) {
				args.input.Jobs <- []calculator.RidePart{{0, calculator.Coordinate{0, 0}, 0}}
				args.input.Jobs <- []calculator.RidePart{{1, calculator.Coordinate{0, 0}, 0}}
				args.input.Jobs <- []calculator.RidePart{{2, calculator.Coordinate{0, 0}, 0}}
				close(args.input.Jobs)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.run(tt.args)

			RunWorker(tt.args.input)

			for _, fareEstimation := range tt.args.expectedPushes {
				if res := <-tt.args.input.Results; res != fareEstimation {
					t.Errorf("RunWorker pushed unexpected result channel = %v, want %v", res, fareEstimation)
				}
			}

			if res := <-tt.args.input.Done; res != nil {
				t.Errorf("RunWorker did not send nil to Done channel, but: %v", res)
			}

			tt.args.input.Wg.Wait()

			select {

			case <-time.After(150 * time.Millisecond):
				break

			case <-tt.args.input.Results:
				t.Errorf("RunWorker pushed more messages to result channel than expected")
			}
		})
	}
}

func TestResultWriter(t *testing.T) {
	type args struct {
		inputs chan model.RideFareEstimation
		wg     *sync.WaitGroup
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"with 2 rides",
			args{
				make(chan model.RideFareEstimation, 10),
				newGroup(),
			},
			"1,14.51\n2,45.12\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.args.inputs <- sampleFareEstimation1
			tt.args.inputs <- sampleFareEstimation2

			close(tt.args.inputs)

			output := bytes.NewBufferString("")

			ResultWriter(output, tt.args.inputs, tt.args.wg)

			if output.String() != tt.want {
				t.Errorf("Bad output, expected `%v`, got `%v`", tt.want, output.String())
			}
		})
	}
}

func TestChannelCloser(t *testing.T) {
	type args struct {
		input ChannelCloserInput
	}
	tests := []struct {
		name string
		args args
		run  func(args)
	}{
		{
			"with 5 workers",
			args{
				ChannelCloserInput{
					5,
					make(chan interface{}, 5),
					make(chan model.RideFareEstimation, 5),
					newGroup(),
				},
			},
			func(args args) {
				for i := 0; i < args.input.Count; i++ {
					args.input.Done <- nil
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.run(tt.args)

			CloseResultChannelWhenWorkersDone(tt.args.input)
		})

		tt.args.input.Wg.Wait()

		for range tt.args.input.Done {
			t.Errorf("Done channel was not closed by CloseResultChannelWhenWorkersDone")
		}
		for range tt.args.input.Results {
			t.Errorf("Results channel was not closed by CloseResultChannelWhenWorkersDone")
		}
	}
}

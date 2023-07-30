package model

import (
	"reflect"
	"testing"
)

func TestRideFareEstimation_toStringSlice(t *testing.T) {
	type args struct {
		rideFareEstimation RideFareEstimation
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"second decimal place is rounded up",
			args{RideFareEstimation{RideID: 100, CostEstimation: 41.1456}},
			[]string{"100", "41.15"},
		},
		{
			"second decimal place is rounded down",
			args{RideFareEstimation{RideID: 100, CostEstimation: 41.1449}},
			[]string{"100", "41.14"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.rideFareEstimation.ToStringSlice(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RideFareEstimation.toStringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

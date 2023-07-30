package calculator

import (
	"fmt"
	"harry-pap/beat_assignment/model"
	"math"
	"reflect"
	"testing"
	"time"
)

var (
	Coord1Part1          = Coordinate{Latitude: 51.365184, Longitude: -2.388245}
	Coord1Part2          = Coordinate{Latitude: 52.052135, Longitude: -1.269958}
	Coord1Part12Distance = float64(108.49)
	Coord2Part1          = Coordinate{Latitude: 51.415597, Longitude: -2.397278}
	Coord2Part2          = Coordinate{Latitude: 51.420850, Longitude: -2.392416}
	Coord2Part12Distance = float64(0.674428)
	Coord3Part1          = Coordinate{Latitude: 49.143352, Longitude: 3.519707}
	Coord3Part2          = Coordinate{Latitude: 49.150684, Longitude: 3.532212}
	Coord3Part3          = Coordinate{Latitude: 49.146490, Longitude: 3.548239}
	Coord3Part4          = Coordinate{Latitude: 49.162513, Longitude: 3.588268}
	Coord3Part12Distance = float64(1.22)
	Coord3Part23Distance = float64(1.26)
	Coord3Part34Distance = float64(3.41)
	Coord3TotalDistance  = Coord3Part12Distance + Coord3Part23Distance + Coord3Part34Distance
)

func TestHarvestineInKilometers(t *testing.T) {
	nameTemplate := "%f distance"
	type args struct {
		from Coordinate
		to   Coordinate
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{fmt.Sprintf(nameTemplate, Coord1Part12Distance),
			args{Coord1Part1, Coord1Part2},
			Coord1Part12Distance,
		}, {fmt.Sprintf(nameTemplate, Coord2Part12Distance),
			args{Coord2Part1, Coord2Part2},
			Coord2Part12Distance,
		}, {fmt.Sprintf(nameTemplate, Coord3Part12Distance),
			args{Coord3Part1, Coord3Part2},
			Coord3Part12Distance,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HarvestineInKilometers(tt.args.from, tt.args.to); !Equal(got, tt.want, 0.01) {
				t.Errorf("HarvestineInKilometers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateKmPerHour(t *testing.T) {
	type args struct {
		start RidePart
		end   RidePart
	}
	type result struct {
		res float64
		err error
	}
	tests := []struct {
		name string
		args args
		want result
	}{
		{"calculates expected speed for coord1",
			args{
				start: RidePart{1, Coord1Part1, int32(parseDatetime("2018-12-12T11:45:00Z").Unix())},
				end:   RidePart{1, Coord1Part2, int32(parseDatetime("2018-12-12T12:45:00Z").Unix())},
			},
			result{Coord1Part12Distance, nil},
		},
		{"calculates expected speed for coord2",
			args{
				start: RidePart{1, Coord2Part1, int32(parseDatetime("2018-12-12T11:00:00Z").Unix())},
				end:   RidePart{1, Coord2Part2, int32(parseDatetime("2018-12-12T11:30:00Z").Unix())},
			},
			result{Coord2Part12Distance * 2, nil},
		},
		{"calculates expected speed for coord2 when start == end",
			args{
				start: RidePart{1, Coord2Part1, int32(parseDatetime("2018-12-12T11:00:00Z").Unix())},
				end:   RidePart{1, Coord2Part2, int32(parseDatetime("2018-12-12T11:00:00Z").Unix())},
			},
			result{Coord2Part12Distance / minimumTimeSlotInHours, nil},
		},
		{"returns an error when start > end",
			args{
				start: RidePart{1, Coord2Part1, int32(parseDatetime("2018-12-12T11:00:00Z").Unix())},
				end:   RidePart{1, Coord2Part2, int32(parseDatetime("2018-12-12T10:00:00Z").Unix())},
			},
			result{0, errInvalidTimestamp},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := CalculateKmPerHour(tt.args.start, tt.args.end); !Equal(got, tt.want.res, 0.01) || tt.want.err != err {
				t.Errorf("CalculateKmPerHour() = %v, %v, want %v, %v", got, err, tt.want.res, tt.want.err)
			}
		})
	}
}

func TestGetValidSegments(t *testing.T) {
	var coord1Part3 = Coordinate{52.052135, -1.269958}
	type args struct {
		entries []RidePart
	}
	tests := []struct {
		name string
		args args
		want []RideSegment
	}{
		{
			"all segments are valid and are retained",
			args{[]RidePart{
				{1, Coord1Part1, int32(parseDatetime("2018-12-12T11:45:00Z").Unix())},
				{1, Coord1Part2, int32(parseDatetime("2018-12-12T13:45:00Z").Unix())},
				{1, coord1Part3, int32(parseDatetime("2018-12-12T15:45:00Z").Unix())},
			},
			},
			[]RideSegment{
				{RidePart{1, Coord1Part1, int32(parseDatetime("2018-12-12T11:45:00Z").Unix())},
					RidePart{1, Coord1Part2, int32(parseDatetime("2018-12-12T13:45:00Z").Unix())}},
				{RidePart{1, Coord1Part2, int32(parseDatetime("2018-12-12T13:45:00Z").Unix())},
					RidePart{1, coord1Part3, int32(parseDatetime("2018-12-12T15:45:00Z").Unix())}},
			}},
		{
			"the invalid segments are discarded",
			args{[]RidePart{
				{1, Coord1Part1, int32(parseDatetime("2018-12-12T11:45:00Z").Unix())},
				{1, Coord1Part2, int32(parseDatetime("2018-12-12T11:45:01Z").Unix())},
				{1, coord1Part3, int32(parseDatetime("2018-12-12T15:45:00Z").Unix())},
			},
			},
			[]RideSegment{
				{RidePart{1, Coord1Part1, int32(parseDatetime("2018-12-12T11:45:00Z").Unix())},
					RidePart{1, coord1Part3, int32(parseDatetime("2018-12-12T15:45:00Z").Unix())}},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetValidSegments(tt.args.entries); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetValidSegments() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRideSegment_GetFare(t *testing.T) {
	type args struct {
		Start RidePart
		End   RidePart
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"day fare",
			args{RidePart{1, Coord1Part1, int32(parseDatetime("2018-12-12T11:45:00Z").Unix())},
				RidePart{1, Coord1Part2, int32(parseDatetime("2018-12-12T13:45:00Z").Unix())}},
			Coord1Part12Distance * dayFarePerKm,
		},
		{"night fare",
			args{RidePart{1, Coord1Part1, int32(parseDatetime("2018-12-12T01:45:00Z").Unix())},
				RidePart{1, Coord1Part2, int32(parseDatetime("2018-12-12T03:45:00Z").Unix())}},
			Coord1Part12Distance * nightFarePerKm,
		},
		{"idle fare",
			args{RidePart{1, Coord1Part1, int32(parseDatetime("2018-12-12T01:45:00Z").Unix())},
				RidePart{1, Coord2Part1, int32(parseDatetime("2018-12-12T03:45:00Z").Unix())}},
			2 * idleFarePerHour,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segment := RideSegment{
				Start: tt.args.Start,
				End:   tt.args.End,
			}
			if got := segment.GetFare(); !Equal(got, tt.want, 0.01) {
				t.Errorf("RideSegment.GetFare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateFareForRide(t *testing.T) {
	type args struct {
		entries []RidePart
		id      int64
	}
	type want struct {
		res model.RideFareEstimation
		err error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"single idle segment",
			args{
				[]RidePart{
					{1, Coord2Part1, int32(parseDatetime("2018-12-12T11:45:00Z").Unix())},
					{1, Coord2Part2, int32(parseDatetime("2018-12-12T12:45:00Z").Unix())},
				},
				1},
			want{model.RideFareEstimation{RideID: 1, CostEstimation: idleFarePerHour + flagValue}, nil},
		},
		{
			"many day segments",
			args{
				[]RidePart{
					{1, Coord3Part1, int32(parseDatetime("2018-12-12T11:00:00Z").Unix())},
					{1, Coord3Part2, int32(parseDatetime("2018-12-12T11:03:00Z").Unix())},
					{1, Coord3Part3, int32(parseDatetime("2018-12-12T11:06:00Z").Unix())},
					{1, Coord3Part4, int32(parseDatetime("2018-12-12T11:14:00Z").Unix())},
				},
				1},
			want{model.RideFareEstimation{RideID: 1, CostEstimation: (Coord3TotalDistance * dayFarePerKm) + flagValue}, nil},
		},
		{
			"many night segments",
			args{
				[]RidePart{
					{1, Coord3Part1, int32(parseDatetime("2018-12-12T03:00:00Z").Unix())},
					{1, Coord3Part2, int32(parseDatetime("2018-12-12T03:03:00Z").Unix())},
					{1, Coord3Part3, int32(parseDatetime("2018-12-12T03:06:00Z").Unix())},
					{1, Coord3Part4, int32(parseDatetime("2018-12-12T03:14:00Z").Unix())},
				},
				1},
			want{model.RideFareEstimation{RideID: 1, CostEstimation: (Coord3TotalDistance * nightFarePerKm) + flagValue}, nil},
		},
		{
			"many idle segments",
			args{
				[]RidePart{
					{1, Coord3Part1, int32(parseDatetime("2018-12-12T03:00:00Z").Unix())},
					{1, Coord3Part2, int32(parseDatetime("2018-12-12T05:00:00Z").Unix())},
					{1, Coord3Part3, int32(parseDatetime("2018-12-12T07:00:00Z").Unix())},
					{1, Coord3Part4, int32(parseDatetime("2018-12-12T09:00:00Z").Unix())},
				},
				1},
			want{model.RideFareEstimation{RideID: 1, CostEstimation: (6 * idleFarePerHour) + flagValue}, nil},
		},
		{
			"combination of all day,night,idle segments",
			args{
				[]RidePart{
					{1, Coord3Part1, int32(parseDatetime("2018-12-12T04:58:00Z").Unix())},
					{1, Coord3Part2, int32(parseDatetime("2018-12-12T05:01:00Z").Unix())},
					{1, Coord3Part3, int32(parseDatetime("2018-12-12T05:03:00Z").Unix())},
					{1, Coord3Part4, int32(parseDatetime("2018-12-12T07:03:00Z").Unix())},
				},
				1},
			want{model.RideFareEstimation{RideID: 1, CostEstimation: (Coord3Part12Distance * nightFarePerKm) +
				(Coord3Part23Distance * dayFarePerKm) +
				(2 * idleFarePerHour) +
				flagValue}, nil},
		},
		{
			"no segments returns an error",
			args{[]RidePart{}, 1},
			want{model.RideFareEstimation{}, errNotEnoughSegments},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateFareForRide(tt.args.entries)
			if err != tt.want.err || got.RideID != tt.want.res.RideID || !Equal(got.CostEstimation, tt.want.res.CostEstimation, 2) {
				t.Errorf("CalculateFareForRide() = %v,%v want %v", got, err, tt.want)
			}
		})
	}
}

func parseDatetime(str string) time.Time {
	t, err := time.Parse(time.RFC3339, str)

	if err != nil {
		panic(err)
	}

	return t
}

func Equal(x, y, tolerance float64) bool {
	return math.Abs(x-y) < tolerance
}

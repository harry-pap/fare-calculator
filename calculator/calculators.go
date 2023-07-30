package calculator

import (
	"errors"
	"fmt"
	"harry-pap/beat_assignment/model"
	"math"
	"time"
)

const (
	dayFarePerKm           = 0.74
	nightFarePerKm         = 1.30
	idleFarePerHour        = 11.90
	earthRadius            = float64(6371)
	minimumRide            = 3.47
	flagValue              = 1.30
	minimumTimeSlotInHours = 0.0001
)

var errNotEnoughSegments = errors.New("not_enough_segments")
var errInvalidTimestamp = errors.New("end_timestamp_not_greater_than_start")

// Coordinate represents a position with Latitude and Longitude
type Coordinate struct {
	Latitude  float64
	Longitude float64
}

// RidePart represents a part of a given ride, with a RideID, a Coordinate and a Unix timestamp
type RidePart struct {
	RideID     int64
	Coordinate Coordinate
	Timestamp  int32
}

// RideSegment represents a segment of two RideParts
type RideSegment struct {
	Start RidePart
	End   RidePart
}

// CalculateFareForRide calculates the fare of the ride. Invalid ride parts(where speed is over 100km/hour) are not included
// If the cost is less than that of the minimum fare(3.47), then the minimum fare is returned.
func CalculateFareForRide(entries []RidePart) (model.RideFareEstimation, error) {
	segments := GetValidSegments(entries)

	if len(segments) == 0 {
		return model.RideFareEstimation{}, errNotEnoughSegments
	}
	sum := flagValue

	for _, segment := range segments {
		sum += segment.GetFare()
	}

	if sum < minimumRide {
		sum = minimumRide
	}

	return model.RideFareEstimation{RideID: entries[0].RideID, CostEstimation: sum}, nil
}

// GetValidSegments filters out the second part of segments, in which the speed is found to be > 100KM/H, as they are considered erroneous
func GetValidSegments(entries []RidePart) []RideSegment {

	result := make([]RideSegment, 0, len(entries))

	for i, j := 0, 1; i < len(entries)-1 && j < len(entries); j++ {
		segment := RideSegment{Start: entries[i], End: entries[j]}

		isValid, err := segment.isValid()
		if err != nil {
			fmt.Println("Ignoring ", segment, " due to error: ", err)
		} else if isValid {
			result = append(result, segment)

			i = j
		}
	}

	return result
}

// CalculateKmPerHour calculates the speed between two given RideParts, in KM/H
// An error is returned if: start timestamp > end timestamp
// If: start timestamp == end timestamp, then time difference is set to 0.4 seconds
func CalculateKmPerHour(start RidePart, end RidePart) (float64, error) {
	kilometers := HarvestineInKilometers(
		Coordinate{Latitude: start.Coordinate.Latitude, Longitude: start.Coordinate.Longitude},
		Coordinate{Latitude: end.Coordinate.Latitude, Longitude: end.Coordinate.Longitude})

	hours := secondsToHours(end.Timestamp - start.Timestamp)

	if hours < 0 {
		return 0, errInvalidTimestamp
	} else if hours == 0 {
		hours = minimumTimeSlotInHours
	}
	return kilometers / hours, nil
}

// GetFare calculates and returns the fare of a ride segment.
// If the ride is NOT idle(speed is > 10 KM/H):
//   - If the START timestamp is between 00:00:00 and 05:00:00(not inclusive), fare = 1.30 * km driven
//   - If the START timestamp is between 05:00:00 and 00:00:00(not inclusive), fare = 0.74 * km driven
//
// If the ride is IDLE: fare = 11.90 * segment hours
func (segment RideSegment) GetFare() float64 {
	isIdle, err := segment.isIdle()

	if err != nil {
		fmt.Println("Ignoring ", segment, " due to error: ", err)
	} else if isIdle {
		return idleFarePerHour * (secondsToHours(segment.End.Timestamp - segment.Start.Timestamp))
	}

	kmDriven := HarvestineInKilometers(segment.Start.Coordinate, segment.End.Coordinate)
	if isNightHours(int64(segment.Start.Timestamp)) {
		return kmDriven * nightFarePerKm
	}
	return kmDriven * dayFarePerKm
}

// HarvestineInKilometers uses the Harvestine formula, to calculate the distance between two Coordinates, in kilometers
func HarvestineInKilometers(from Coordinate, to Coordinate) float64 {
	var deltaLatitude = (to.Latitude - from.Latitude) * (math.Pi / 180)
	var deltaLongitude = (to.Longitude - from.Longitude) * (math.Pi / 180)

	var a = math.Sin(deltaLatitude/2)*math.Sin(deltaLatitude/2) +
		math.Cos(from.Latitude*(math.Pi/180))*
			math.Cos(to.Latitude*(math.Pi/180))*
			math.Sin(deltaLongitude/2)*
			math.Sin(deltaLongitude/2)

	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

func (segment RideSegment) isValid() (bool, error) {
	kmPerHour, err := CalculateKmPerHour(segment.Start, segment.End)
	if err != nil {
		return false, err
	}

	if kmPerHour > 100 {
		return false, nil
	}
	return true, nil
}

func (segment RideSegment) isIdle() (bool, error) {
	kmPerHour, err := CalculateKmPerHour(segment.Start, segment.End)

	if err != nil {
		return false, err
	}

	if kmPerHour > 10 {
		return false, nil
	}
	return true, nil
}

func secondsToHours(seconds int32) float64 {
	return (time.Duration(seconds) * time.Second).Hours()
}

func isNightHours(timestamp int64) bool {
	ts := time.Unix(timestamp, 0).UTC()

	return ts.Hour() >= 0 && ts.Hour() < 5
}

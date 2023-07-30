package model

import "strconv"

// RideFareEstimation contains the fare estimation of a ride, including the ride id, and the cost estimation
type RideFareEstimation struct {
	RideID         int64
	CostEstimation float64
}

// ToStringSlice converts a RideFareEstimation, into a []string, representing its fields
func (rideFareEstimation RideFareEstimation) ToStringSlice() []string {
	return []string{
		strconv.FormatInt(rideFareEstimation.RideID, 10),
		strconv.FormatFloat(rideFareEstimation.CostEstimation, 'f', 2, 64),
	}
}

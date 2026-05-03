package model

import "time"

type BaseEntity struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type DashboardMetrics struct {
	Range             string    `json:"range"`
	RangeLabel        string    `json:"rangeLabel"`
	RangeStart        time.Time `json:"rangeStart"`
	RangeEnd          time.Time `json:"rangeEnd"`
	InRentCarCount    int64     `json:"inRentCarCount"`
	IdleCarCount      int64     `json:"idleCarCount"`
	RentOrderCount    int64     `json:"rentOrderCount"`
	CustomerCount     int64     `json:"customerCount"`
	RefreshedAt       time.Time `json:"refreshedAt"`
	VehicleScopeLabel string    `json:"vehicleScopeLabel"`
}

package search

import "math"

type MinAggregationResult struct {
	Name string
	Value float64	//+inf means missing value
}

func (a *MinAggregationResult) GetName() string {
	return a.Name
}

func (a *MinAggregationResult) GetType() AggregationType {
	return AggregationMinType
}

func (a *MinAggregationResult) HasValue() bool {
	return a != nil && !math.IsInf(a.Value, 1)
}
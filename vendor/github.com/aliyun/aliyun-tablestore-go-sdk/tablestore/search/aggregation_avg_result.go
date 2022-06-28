package search

import "math"

type AvgAggregationResult struct {
	Name string
	Value float64	//+inf means missing value
}

func (a *AvgAggregationResult) GetName() string {
	return a.Name
}

func (a *AvgAggregationResult) GetType() AggregationType {
	return AggregationAvgType
}

func (a *AvgAggregationResult) HasValue() bool {
	return a != nil && !math.IsInf(a.Value, 1)
}
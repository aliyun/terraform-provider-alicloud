package search

import "math"

type MaxAggregationResult struct {
	Name string
	Value float64	//-inf means missing value
}

func (a *MaxAggregationResult) GetName() string {
	return a.Name
}

func (a *MaxAggregationResult) GetType() AggregationType {
	return AggregationMaxType
}

func (a *MaxAggregationResult) HasValue() bool {
	return a != nil && !math.IsInf(a.Value, -1)
}
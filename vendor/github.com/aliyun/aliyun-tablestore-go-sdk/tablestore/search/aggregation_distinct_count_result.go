package search

type DistinctCountAggregationResult struct {
	Name string
	Value int64
}

func (a *DistinctCountAggregationResult) GetName() string {
	return a.Name
}

func (a *DistinctCountAggregationResult) GetType() AggregationType {
	return AggregationDistinctCountType
}
package search

type CountAggregationResult struct {
	Name string
	Value int64
}

func (a *CountAggregationResult) GetName() string {
	return a.Name
}

func (a *CountAggregationResult) GetType() AggregationType {
	return AggregationCountType
}
package search

type SumAggregationResult struct {
	Name string
	Value float64	//0 means either field missing or 'sum is actually 0'
}

func (a *SumAggregationResult) GetName() string {
	return a.Name
}

func (a *SumAggregationResult) GetType() AggregationType {
	return AggregationSumType
}
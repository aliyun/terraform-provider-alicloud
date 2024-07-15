package search

type GroupByDateHistogramItem struct {
	Timestamp int64
	RowCount  int64

	SubAggregations AggregationResults
	SubGroupBys     GroupByResults
}

type GroupByDateHistogramResult struct {
	Name  string
	Items []GroupByDateHistogramItem
}

func (g *GroupByDateHistogramResult) GetName() string {
	return g.Name
}

func (g *GroupByDateHistogramResult) GetType() GroupByType {
	return GroupByDateHistogramType
}

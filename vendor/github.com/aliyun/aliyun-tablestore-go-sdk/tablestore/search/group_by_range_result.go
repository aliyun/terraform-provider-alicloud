package search

type GroupByRangeResultItem struct {
	RowCount        int64
	From			float64
	To				float64

	SubAggregations AggregationResults
	SubGroupBys     GroupByResults
}

type GroupByRangeResult struct {
	Name string
	Items []GroupByRangeResultItem
}

func (g *GroupByRangeResult) GetName() string {
	return g.Name
}

func (g *GroupByRangeResult) GetType() GroupByType {
	return GroupByRangeType
}
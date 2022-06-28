package search

type GroupByFilterResultItem struct {
	RowCount        int64
	SubAggregations AggregationResults
	SubGroupBys     GroupByResults
}

type GroupByFilterResult struct {
	Name string
	Items []GroupByFilterResultItem
}

func (g *GroupByFilterResult) GetName() string {
	return g.Name
}

func (g *GroupByFilterResult) GetType() GroupByType {
	return GroupByFilterType
}
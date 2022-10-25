package search

type GroupByFieldResultItem struct {
	Key             string
	RowCount        int64
	SubAggregations AggregationResults
	SubGroupBys     GroupByResults
}

type GroupByFieldResult struct {
	Name string
	Items []GroupByFieldResultItem
}

func (g *GroupByFieldResult) GetName() string {
	return g.Name
}

func (g *GroupByFieldResult) GetType() GroupByType {
	return GroupByFieldType
}
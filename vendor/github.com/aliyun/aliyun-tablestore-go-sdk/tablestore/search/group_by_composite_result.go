package search

type GroupByCompositeResultItem struct {
	Keys            []*string
	RowCount        int64
	SubAggregations AggregationResults
	SubGroupBys     GroupByResults
}

type GroupByCompositeResult struct {
	Name               string
	Items              []GroupByCompositeResultItem
	SourceGroupByNames []string
	NextToken          *string
}

func (g *GroupByCompositeResult) GetName() string {
	return g.Name
}

func (g *GroupByCompositeResult) GetType() GroupByType {
	return GroupByCompositeType
}

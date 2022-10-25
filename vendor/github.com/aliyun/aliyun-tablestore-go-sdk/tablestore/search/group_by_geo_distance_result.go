package search

type GroupByGeoDistanceResultItem struct {
	RowCount        int64
	From			float64
	To				float64

	SubAggregations AggregationResults
	SubGroupBys     GroupByResults
}

type GroupByGeoDistanceResult struct {
	Name string
	Items []GroupByGeoDistanceResultItem
}

func (g *GroupByGeoDistanceResult) GetName() string {
	return g.Name
}

func (g *GroupByGeoDistanceResult) GetType() GroupByType {
	return GroupByGeoDistanceType
}
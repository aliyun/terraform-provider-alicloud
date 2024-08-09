package search

type GroupByGeoGridResultItem struct {
	Key      string
	GeoGrid  GeoGrid
	RowCount int64

	SubAggregations AggregationResults
	SubGroupBys     GroupByResults
}

type GroupByGeoGridResult struct {
	Name  string
	Items []GroupByGeoGridResultItem
}

func (g *GroupByGeoGridResult) GetName() string {
	return g.Name
}

func (g *GroupByGeoGridResult) GetType() GroupByType {
	return GroupByGeoGridType
}

package search

import "github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/search/model"

type GroupByHistogramItem struct {
	Key   model.ColumnValue
	Value int64

	SubAggregations AggregationResults
	SubGroupBys     GroupByResults
}

type GroupByHistogramResult struct {
	Name  string
	Items []GroupByHistogramItem
}

func (g *GroupByHistogramResult) GetName() string {
	return g.Name
}

func (g *GroupByHistogramResult) GetType() GroupByType {
	return GroupByHistogramType
}

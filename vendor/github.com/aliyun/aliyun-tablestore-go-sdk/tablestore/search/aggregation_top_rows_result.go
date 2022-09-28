package search

import (
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/search/model"
)

type TopRowsAggregationResult struct {
	Name  string
	Value []model.Row
}

func (a *TopRowsAggregationResult) GetName() string {
	return a.Name
}

func (a *TopRowsAggregationResult) GetType() AggregationType {
	return AggregationTopRowsType
}

func (a *TopRowsAggregationResult) HasValue() bool {
	return a != nil && len(a.Value) != 0
}

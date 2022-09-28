package search

import (
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/search/model"
)

type PercentilesAggregationItem struct {
	Key   float64
	Value model.ColumnValue
}

type PercentilesAggregationResult struct {
	Name                    string
	PercentilesAggregationItems []PercentilesAggregationItem
}

func (a *PercentilesAggregationResult) GetName() string {
	return a.Name
}

func (a *PercentilesAggregationResult) GetType() AggregationType {
	return AggregationPercentilesType
}

func (a *PercentilesAggregationResult) HasValue() bool {
	return a != nil && len(a.PercentilesAggregationItems) > 0
}

package search

import (
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/golang/protobuf/proto"
)

type TopRowsAggregation struct {
	AggName string
	Limit   *int32
	Sort    *Sort
}

func (a *TopRowsAggregation) GetName() string {
	return a.AggName
}

func (a *TopRowsAggregation) GetType() AggregationType {
	return AggregationTopRowsType
}

func (a *TopRowsAggregation) Serialize() ([]byte, error) {
	pbAgg := &otsprotocol.TopRowsAggregation{}

	pbAgg.Limit = a.Limit
	if a.Sort != nil {
		var err error
		pbAgg.Sort, err = (*Sort).ProtoBuffer(a.Sort)
		if err != nil {
			return nil, err
		}
	}

	data, err := proto.Marshal(pbAgg)
	return data, err
}

func (a *TopRowsAggregation) ProtoBuffer() (*otsprotocol.Aggregation, error) {
	return BuildPBForAggregation(a)
}

func (a *TopRowsAggregation) SetName(name string) *TopRowsAggregation {
	a.AggName = name
	return a
}

func (a *TopRowsAggregation) SetLimit(limit int32) *TopRowsAggregation {
	a.Limit = &limit
	return a
}

func (a *TopRowsAggregation) SetSort(sort *Sort) *TopRowsAggregation {
	a.Sort = sort
	return a
}

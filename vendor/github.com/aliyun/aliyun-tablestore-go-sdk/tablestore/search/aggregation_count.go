package search

import (
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/golang/protobuf/proto"
)

type CountAggregation struct {
	AggName		string
	Field		string
}

func (a *CountAggregation) GetName() string {
	return a.AggName
}

func (a *CountAggregation) GetType() AggregationType {
	return AggregationCountType
}

func (a *CountAggregation) Serialize() ([]byte, error) {
	pbAgg := &otsprotocol.CountAggregation{}

	//field_name
	pbAgg.FieldName = &a.Field

	data, err := proto.Marshal(pbAgg)
	return data, err
}

func (a *CountAggregation) ProtoBuffer() (*otsprotocol.Aggregation, error) {
	return BuildPBForAggregation(a)
}

//building chain
func (a *CountAggregation) Name(name string) *CountAggregation {
	a.AggName = name
	return a
}

func (a *CountAggregation) FieldName(fieldName string) *CountAggregation {
	a.Field = fieldName
	return a
}

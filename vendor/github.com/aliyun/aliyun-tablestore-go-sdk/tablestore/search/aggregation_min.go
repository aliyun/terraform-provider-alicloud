package search

import (
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/golang/protobuf/proto"
)

type MinAggregation struct {
	AggName		string
	Field	string
	MissingValue 	interface{}
}

func (a *MinAggregation) GetName() string {
	return a.AggName
}

func (a *MinAggregation) GetType() AggregationType {
	return AggregationMinType
}

func (a *MinAggregation) Serialize() ([]byte, error) {
	pbAgg := &otsprotocol.MinAggregation{}

	//field_name
	pbAgg.FieldName = &a.Field

	//missing
	if a.MissingValue != nil {
		vt, err := ToVariantValue(a.MissingValue)
		if err != nil {
			return nil, err
		}
		pbAgg.Missing = []byte(vt)
	}

	data, err := proto.Marshal(pbAgg)
	return data, err
}

func (a *MinAggregation) ProtoBuffer() (*otsprotocol.Aggregation, error) {
	return BuildPBForAggregation(a)
}

//building chain
func (a *MinAggregation) Name(name string) *MinAggregation {
	a.AggName = name
	return a
}

func (a *MinAggregation) FieldName(fieldName string) *MinAggregation {
	a.Field = fieldName
	return a
}

func (a *MinAggregation) Missing(missing interface{}) *MinAggregation {
	a.MissingValue = missing
	return a
}

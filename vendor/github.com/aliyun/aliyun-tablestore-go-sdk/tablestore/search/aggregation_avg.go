package search

import (
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/golang/protobuf/proto"
)

type AvgAggregation struct {
	AggName			string
	Field			string
	MissingValue 	interface{}
}

func (a *AvgAggregation) GetName() string {
	return a.AggName
}

func (a *AvgAggregation) GetType() AggregationType {
	return AggregationAvgType
}

func (a *AvgAggregation) Serialize() ([]byte, error) {
	pbAgg := &otsprotocol.AvgAggregation{}

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

func (a *AvgAggregation) ProtoBuffer() (*otsprotocol.Aggregation, error) {
	return BuildPBForAggregation(a)
}

//building chain
func (a *AvgAggregation) Name(name string) *AvgAggregation {
	a.AggName = name
	return a
}

func (a *AvgAggregation) FieldName(fieldName string) *AvgAggregation {
	a.Field = fieldName
	return a
}

func (a *AvgAggregation) Missing(missing interface{}) *AvgAggregation {
	a.MissingValue = missing
	return a
}

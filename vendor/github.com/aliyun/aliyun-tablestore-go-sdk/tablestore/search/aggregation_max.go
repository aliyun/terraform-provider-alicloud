package search

import (
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/golang/protobuf/proto"
)

type MaxAggregation struct {
	AggName			string
	Field			string
	MissingValue 	interface{}
}

func (a *MaxAggregation) GetName() string {
	return a.AggName
}

func (a *MaxAggregation) GetType() AggregationType {
	return AggregationMaxType
}

func (a *MaxAggregation) Serialize() ([]byte, error) {
	pbAgg := &otsprotocol.MaxAggregation{}

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

func (a *MaxAggregation) ProtoBuffer() (*otsprotocol.Aggregation, error) {
	return BuildPBForAggregation(a)
}

//building chain
func (a *MaxAggregation) Name(name string) *MaxAggregation {
	a.AggName = name
	return a
}

func (a *MaxAggregation) FieldName(fieldName string) *MaxAggregation {
	a.Field = fieldName
	return a
}

func (a *MaxAggregation) Missing(missing interface{}) *MaxAggregation {
	a.MissingValue = missing
	return a
}

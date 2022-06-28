package search

import (
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/golang/protobuf/proto"
)

type SumAggregation struct {
	AggName			string
	Field			string
	MissingValue	interface{}
}

func (a *SumAggregation) GetName() string {
	return a.AggName
}

func (a *SumAggregation) GetType() AggregationType {
	return AggregationSumType
}

func (a *SumAggregation) Serialize() ([]byte, error) {
	pbAgg := &otsprotocol.SumAggregation{}

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

func (a *SumAggregation) ProtoBuffer() (*otsprotocol.Aggregation, error) {
	return BuildPBForAggregation(a)
}

//building chain
func (a *SumAggregation) Name(name string) *SumAggregation {
	a.AggName = name
	return a
}

func (a *SumAggregation) FieldName(fieldName string) *SumAggregation {
	a.Field = fieldName
	return a
}

func (a *SumAggregation) Missing(missing interface{}) *SumAggregation {
	a.MissingValue = missing
	return a
}

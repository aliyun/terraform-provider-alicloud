package search

import (
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/golang/protobuf/proto"
)

type DistinctCountAggregation struct {
	AggName			string
	Field			string
	MissingValue 	interface{}
}

func (a *DistinctCountAggregation) GetName() string {
	return a.AggName
}

func (a *DistinctCountAggregation) GetType() AggregationType {
	return AggregationDistinctCountType
}

func (a *DistinctCountAggregation) Serialize() ([]byte, error) {
	pbAgg := &otsprotocol.DistinctCountAggregation{}

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

func (a *DistinctCountAggregation) ProtoBuffer() (*otsprotocol.Aggregation, error) {
	return BuildPBForAggregation(a)
}

//building chain
func (a *DistinctCountAggregation) Name(name string) *DistinctCountAggregation {
	a.AggName = name
	return a
}

func (a *DistinctCountAggregation) FieldName(fieldName string) *DistinctCountAggregation {
	a.Field = fieldName
	return a
}

func (a *DistinctCountAggregation) Missing(missing interface{}) *DistinctCountAggregation {
	a.MissingValue = missing
	return a
}

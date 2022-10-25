package search

import (
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/golang/protobuf/proto"
)

type PercentilesAggregation struct {
	AggName      string
	Field        string
	Percents     []float64
	MissingValue interface{}
}

func (a *PercentilesAggregation) GetName() string {
	return a.AggName
}

func (a *PercentilesAggregation) GetType() AggregationType {
	return AggregationPercentilesType
}

func (a *PercentilesAggregation) GetPercents() []float64 {
	return a.Percents
}

func (a *PercentilesAggregation) Serialize() ([]byte, error) {
	pbAgg := &otsprotocol.PercentilesAggregation{}

	//field_name
	pbAgg.FieldName = &a.Field

	pbAgg.Percentiles = a.Percents
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

func (a *PercentilesAggregation) ProtoBuffer() (*otsprotocol.Aggregation, error) {
	return BuildPBForAggregation(a)
}

func (a *PercentilesAggregation) SetName(name string) *PercentilesAggregation {
	a.AggName = name
	return a
}

func (a *PercentilesAggregation) SetFieldName(fieldName string) *PercentilesAggregation {
	a.Field = fieldName
	return a
}

func (a *PercentilesAggregation) SetPercents(percents []float64) *PercentilesAggregation {
	a.Percents = percents
	return a
}

func (a *PercentilesAggregation) SetMissing(missing interface{}) *PercentilesAggregation {
	a.MissingValue = missing
	return a
}

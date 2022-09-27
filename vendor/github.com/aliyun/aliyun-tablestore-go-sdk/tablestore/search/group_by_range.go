package search

import (
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/golang/protobuf/proto"
)

type GroupByRange struct {
	AggName string

	Field 		string
	RangeList	[]Range

	SubAggList			[]Aggregation
	SubGroupByList		[]GroupBy
}

func (g *GroupByRange) GetName() string {
	return g.AggName
}

func (g *GroupByRange) GetType() GroupByType {
	return GroupByRangeType
}

func (g *GroupByRange) Serialize() ([]byte, error) {
	pbGroupBy := &otsprotocol.GroupByRange{}

	//field_name
	pbGroupBy.FieldName = &g.Field

	//ranges
	pbRanges, err := BuildPBForRanges(g.RangeList)
	if err != nil {
		return nil, err
	}
	pbGroupBy.Ranges = pbRanges

	//sub aggs
	pbAggregations, err := BuildPBForAggregations(g.SubAggList)
	if err != nil {
		return nil, err
	}
	if pbAggregations != nil {
		pbGroupBy.SubAggs = pbAggregations
	}

	//sub group bys
	pbGroupBys, err := BuildPBForGroupBys(g.SubGroupByList)
	if err != nil {
		return nil, err
	}
	if pbGroupBys != nil {
		pbGroupBy.SubGroupBys = pbGroupBys
	}

	data, err := proto.Marshal(pbGroupBy)
	return data, err
}

func (g *GroupByRange) ProtoBuffer() (*otsprotocol.GroupBy, error) {
	return BuildPBForGroupBy(g)
}

func (g *GroupByRange) SubAggregations(subAggregations ...Aggregation) *GroupByRange {
	g.SubAggList = subAggregations
	return g
}

func (g *GroupByRange) SubAggregation(subAggregation Aggregation) *GroupByRange {
	g.SubAggList = append(g.SubAggList, subAggregation)
	return g
}


func (g *GroupByRange) SubGroupBys(subGroupBys ...GroupBy) *GroupByRange {
	g.SubGroupByList = subGroupBys
	return g
}

func (g *GroupByRange) SubGroupBy(subGroupBy GroupBy) *GroupByRange {
	g.SubGroupByList = append(g.SubGroupByList, subGroupBy)
	return g
}

func (g *GroupByRange) Name(name string) *GroupByRange {
	g.AggName = name
	return g
}

func (g *GroupByRange) FieldName(fieldName string) *GroupByRange {
	g.Field = fieldName
	return g
}

//append a Range[from, to)
func (g *GroupByRange) Range(fromInclusive float64, toExclusive float64) *GroupByRange {
	g.RangeList = append(g.RangeList,
		Range { from: fromInclusive, to: toExclusive})
	return g
}

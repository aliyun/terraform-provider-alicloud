package search

import (
	"errors"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/golang/protobuf/proto"
)

type GroupByField struct {
	AggName string

	Field   string
	Sz    *int32
	Sorters []GroupBySorter

	SubAggList			[]Aggregation
	SubGroupByList		[]GroupBy
}

func (g *GroupByField) GetName() string {
	return g.AggName
}

func (g *GroupByField) GetType() GroupByType {
	return GroupByFieldType
}

func (g *GroupByField) Serialize() ([]byte, error) {
	pbGroupBy := &otsprotocol.GroupByField{}

	//field_name
	pbGroupBy.FieldName = &g.Field

	//size
	if g.Sz != nil {
		pbGroupBy.Size = g.Sz
	}

	//sorters
	var pbGroupBySorters []*otsprotocol.GroupBySorter
	for _, sorter := range g.Sorters {
		pbGroupBySorter, err := sorter.ProtoBuffer()
		if err != nil {
			return nil, errors.New("invalid group by sorter:" + err.Error())
		}
		pbGroupBySorters = append(pbGroupBySorters, pbGroupBySorter)
	}

	pbGroupBySort := otsprotocol.GroupBySort{}
	pbGroupBySort.Sorters = pbGroupBySorters
	pbGroupBy.Sort = &pbGroupBySort

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

func (g *GroupByField) ProtoBuffer() (*otsprotocol.GroupBy, error) {
	return BuildPBForGroupBy(g)
}

func (g *GroupByField) SubAggregations(subAggregations ...Aggregation) *GroupByField {
	g.SubAggList = subAggregations
	return g
}

func (g *GroupByField) SubAggregation(subAggregation Aggregation) *GroupByField {
	g.SubAggList = append(g.SubAggList, subAggregation)
	return g
}


func (g *GroupByField) SubGroupBys(subGroupBys ...GroupBy) *GroupByField {
	g.SubGroupByList = subGroupBys
	return g
}

func (g *GroupByField) SubGroupBy(subGroupBy GroupBy) *GroupByField {
	g.SubGroupByList = append(g.SubGroupByList, subGroupBy)
	return g
}

func (g *GroupByField) Name(name string) *GroupByField {
	g.AggName = name
	return g
}

func (g *GroupByField) FieldName(fieldName string) *GroupByField {
	g.Field = fieldName
	return g
}

func (g *GroupByField) Size(size int32) *GroupByField {
	g.Sz = &size
	return g
}

func (g *GroupByField) GroupBySorters(sorters []GroupBySorter) *GroupByField {
	g.Sorters = sorters
	return g
}

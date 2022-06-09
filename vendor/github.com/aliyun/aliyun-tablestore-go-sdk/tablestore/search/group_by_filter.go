package search

import (
	"errors"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/golang/protobuf/proto"
)

type GroupByFilter struct {
	AggName string

	Queries []Query

	SubAggList			[]Aggregation
	SubGroupByList		[]GroupBy
}

func (g *GroupByFilter) GetName() string {
	return g.AggName
}

func (g *GroupByFilter) GetType() GroupByType {
	return GroupByFilterType
}

func (g *GroupByFilter) Serialize() ([]byte, error) {
	pbGroupBy := &otsprotocol.GroupByFilter{}

	//queries
	var pbQueries []*otsprotocol.Query
	for _, query := range g.Queries {
		pbQuery, err := query.ProtoBuffer()
		if err != nil {
			return nil, errors.New("invalid query: " + err.Error())
		}
		pbQueries = append(pbQueries, pbQuery)
	}
	pbGroupBy.Filters = pbQueries

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

func (g *GroupByFilter) ProtoBuffer() (*otsprotocol.GroupBy, error) {
	return BuildPBForGroupBy(g)
}

func (g *GroupByFilter) SubAggregations(subAggregations ...Aggregation) *GroupByFilter {
	g.SubAggList = subAggregations
	return g
}

func (g *GroupByFilter) SubAggregation(subAggregation Aggregation) *GroupByFilter {
	g.SubAggList = append(g.SubAggList, subAggregation)
	return g
}


func (g *GroupByFilter) SubGroupBys(subGroupBys ...GroupBy) *GroupByFilter {
	g.SubGroupByList = subGroupBys
	return g
}

func (g *GroupByFilter) SubGroupBy(subGroupBy GroupBy) *GroupByFilter {
	g.SubGroupByList = append(g.SubGroupByList, subGroupBy)
	return g
}

func (g *GroupByFilter) Name(name string) *GroupByFilter {
	g.AggName = name
	return g
}

func (g *GroupByFilter) Query(query Query) *GroupByFilter {
	g.Queries = append(g.Queries, query)
	return g
}

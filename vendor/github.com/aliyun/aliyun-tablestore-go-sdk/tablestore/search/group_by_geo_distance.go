package search

import (
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/golang/protobuf/proto"
)

type GroupByGeoDistance struct {
	AggName string

	Field 			string
	Origin			GeoPoint
	RangeList		[]Range

	SubAggList			[]Aggregation
	SubGroupByList		[]GroupBy
}

func (g *GroupByGeoDistance) GetName() string {
	return g.AggName
}

func (g *GroupByGeoDistance) GetType() GroupByType {
	return GroupByGeoDistanceType
}

func (g *GroupByGeoDistance) Serialize() ([]byte, error) {
	pbGroupBy := &otsprotocol.GroupByGeoDistance{}

	//field_name
	pbGroupBy.FieldName = &g.Field

	//origin
	origin := otsprotocol.GeoPoint{}
	origin.Lat = &g.Origin.Lat
	origin.Lon = &g.Origin.Lon
	pbGroupBy.Origin = &origin

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

func (g *GroupByGeoDistance) ProtoBuffer() (*otsprotocol.GroupBy, error) {
	return BuildPBForGroupBy(g)
}

func (g *GroupByGeoDistance) SubAggregations(subAggregations ...Aggregation) *GroupByGeoDistance {
	g.SubAggList = subAggregations
	return g
}

func (g *GroupByGeoDistance) SubAggregation(subAggregation Aggregation) *GroupByGeoDistance {
	g.SubAggList = append(g.SubAggList, subAggregation)
	return g
}


func (g *GroupByGeoDistance) SubGroupBys(subGroupBys ...GroupBy) *GroupByGeoDistance {
	g.SubGroupByList = subGroupBys
	return g
}

func (g *GroupByGeoDistance) SubGroupBy(subGroupBy GroupBy) *GroupByGeoDistance {
	g.SubGroupByList = append(g.SubGroupByList, subGroupBy)
	return g
}

func (g *GroupByGeoDistance) Name(name string) *GroupByGeoDistance {
	g.AggName = name
	return g
}

func (g *GroupByGeoDistance) FieldName(fieldName string) *GroupByGeoDistance {
	g.Field = fieldName
	return g
}

func (g *GroupByGeoDistance) CenterPoint(latitude float64, longitude float64) *GroupByGeoDistance {
	g.Origin = GeoPoint{Lat: latitude, Lon: longitude}
	return g
}

//append a Range[from, to)
func (g *GroupByGeoDistance) Range(fromInclusive float64, toExclusive float64) *GroupByGeoDistance {
	g.RangeList = append(g.RangeList,
		Range { from: fromInclusive, to: toExclusive})
	return g
}

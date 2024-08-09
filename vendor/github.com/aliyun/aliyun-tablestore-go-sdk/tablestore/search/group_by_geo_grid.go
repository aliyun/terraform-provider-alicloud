package search

import (
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/search/model"
	"github.com/golang/protobuf/proto"
)

// GroupByGeoGrid 对GeoPoint类型的字段按照地理区域进行分组统计
type GroupByGeoGrid struct {
	GroupByName string

	Field     string
	Precision model.GeoHashPrecision
	Size      int64

	SubAggList     []Aggregation
	SubGroupByList []GroupBy
}

func (g *GroupByGeoGrid) GetName() string {
	return g.GroupByName
}

func (g *GroupByGeoGrid) GetType() GroupByType {
	return GroupByGeoGridType
}

func (g *GroupByGeoGrid) GetField() string {
	return g.Field
}

func (g *GroupByGeoGrid) GetPrecision() model.GeoHashPrecision {
	return g.Precision
}

func (g *GroupByGeoGrid) GetSize() int64 {
	return g.Size
}

func (g *GroupByGeoGrid) Serialize() ([]byte, error) {
	pbGroupBy := &otsprotocol.GroupByGeoGrid{}

	pbGroupBy.FieldName = &g.Field

	pbGroupBy.Precision = g.Precision.ProtoBuffer()

	pbGroupBy.Size = &g.Size

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

func (g *GroupByGeoGrid) ProtoBuffer() (*otsprotocol.GroupBy, error) {
	return BuildPBForGroupBy(g)
}

func (g *GroupByGeoGrid) SubAggregations(subAggregations ...Aggregation) *GroupByGeoGrid {
	g.SubAggList = subAggregations
	return g
}

func (g *GroupByGeoGrid) SubAggregation(subAggregation Aggregation) *GroupByGeoGrid {
	g.SubAggList = append(g.SubAggList, subAggregation)
	return g
}

func (g *GroupByGeoGrid) SubGroupBys(subGroupBys ...GroupBy) *GroupByGeoGrid {
	g.SubGroupByList = subGroupBys
	return g
}

func (g *GroupByGeoGrid) SubGroupBy(subGroupBy GroupBy) *GroupByGeoGrid {
	g.SubGroupByList = append(g.SubGroupByList, subGroupBy)
	return g
}

func (g *GroupByGeoGrid) SetGroupByName(name string) *GroupByGeoGrid {
	g.GroupByName = name
	return g
}

func (g *GroupByGeoGrid) SetField(fieldName string) *GroupByGeoGrid {
	g.Field = fieldName
	return g
}

func (g *GroupByGeoGrid) SetPrecision(precision model.GeoHashPrecision) *GroupByGeoGrid {
	g.Precision = precision
	return g
}

func (g *GroupByGeoGrid) SetSize(size int64) *GroupByGeoGrid {
	g.Size = size
	return g
}
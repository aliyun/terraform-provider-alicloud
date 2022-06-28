package search

import (
	"errors"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/golang/protobuf/proto"
)

type GroupByType int

const (
	GroupByNoneType        GroupByType = 0
	GroupByFieldType       GroupByType = 1
	GroupByRangeType       GroupByType = 2
	GroupByFilterType      GroupByType = 3
	GroupByGeoDistanceType GroupByType = 4
	GroupByHistogramType   GroupByType = 5
)

func (g GroupByType) Enum() *GroupByType {
	newGroupBy := g
	return &newGroupBy
}

func (a GroupByType) String() string {
	switch a {
	case GroupByFieldType:
		return "group_by_field"
	case GroupByRangeType:
		return "group_by_range"
	case GroupByFilterType:
		return "group_by_filter"
	case GroupByGeoDistanceType:
		return "group_by_geo_distance"
	case GroupByHistogramType:
		return "group_by_histogram"
	default:
		return "unknown"
	}
}

func (g GroupByType) ToPB() *otsprotocol.GroupByType {
	switch g {
	case GroupByFieldType:
		return otsprotocol.GroupByType_GROUP_BY_FIELD.Enum()
	case GroupByRangeType:
		return otsprotocol.GroupByType_GROUP_BY_RANGE.Enum()
	case GroupByFilterType:
		return otsprotocol.GroupByType_GROUP_BY_FILTER.Enum()
	case GroupByGeoDistanceType:
		return otsprotocol.GroupByType_GROUP_BY_GEO_DISTANCE.Enum()
	case GroupByHistogramType:
		return otsprotocol.GroupByType_GROUP_BY_HISTOGRAM.Enum()
	default:
		return nil
	}
}

/*
	message GroupBy {
    	optional string name = 1;
    	optional GroupByType type = 2;
    	optional bytes body = 3;
	}
 */
type GroupBy interface {
	//get group by name
	GetName() string

	//get group by type
	GetType() GroupByType

	//build body, implemented by each concrete agg
	Serialize() ([]byte, error)

	// build the whole aggregation, implemented by each concrete agg,
	// using BuildPBForAggregation() defined in agg interface
	ProtoBuffer() (*otsprotocol.GroupBy, error)
}

func BuildPBForGroupBy(g GroupBy) (*otsprotocol.GroupBy, error) {
	pbGroupBy := &otsprotocol.GroupBy{}

	pbGroupBy.Name = proto.String(g.GetName())
	pbGroupBy.Type = g.GetType().ToPB()
	body, err := g.Serialize()
	if err != nil {
		return nil, err
	}
	pbGroupBy.Body = body
	return pbGroupBy, nil
}

func BuildPBForGroupBys(groupBys []GroupBy) (*otsprotocol.GroupBys, error) {
	if len(groupBys) == 0 {
		return nil, nil
	}

	pbGroupBys := new(otsprotocol.GroupBys)
	for _, subGroupBy := range groupBys {
		pbGroupBy, err := subGroupBy.ProtoBuffer()
		if err != nil {
			return nil, errors.New("invalid group by: " + err.Error())
		}
		pbGroupBys.GroupBys = append(pbGroupBys.GroupBys, pbGroupBy)
	}
	return pbGroupBys, nil
}

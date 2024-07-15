package search

import (
	"errors"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/search/model"
	"github.com/golang/protobuf/proto"
)

type GroupByDateHistogram struct {
	GroupByName string

	Field       string
	Interval    model.DateTimeValue
	FieldRange  model.FiledRange
	Missing     interface{}
	MinDocCount *int64
	TimeZone    *string
	Sorters     []GroupBySorter

	SubAggList     []Aggregation
	SubGroupByList []GroupBy
}

func (g *GroupByDateHistogram) GetName() string {
	return g.GroupByName
}

func (g *GroupByDateHistogram) GetType() GroupByType {
	return GroupByDateHistogramType
}

func (g *GroupByDateHistogram) GetField() string {
	return g.Field
}

func (g *GroupByDateHistogram) GetTimeZone() *string {
	return g.TimeZone
}

func (g *GroupByDateHistogram) GetMinDocCount() *int64 {
	return g.MinDocCount
}

func (g *GroupByDateHistogram) GetInterval() model.DateTimeValue {
	return g.Interval
}

func (g *GroupByDateHistogram) Serialize() ([]byte, error) {

	pbGroupBy := &otsprotocol.GroupByDateHistogram{}

	pbGroupBy.FieldName = &g.Field
	pbGroupBy.MinDocCount = g.MinDocCount
	pbGroupBy.TimeZone = g.TimeZone

	if g.Missing != nil {
		vt, err := ToVariantValue(g.Missing)
		if err != nil {
			return nil, err
		}
		pbGroupBy.Missing = vt
	}
	pbGroupBy.Interval = g.Interval.ProtoBuffer()
	pbGroupBy.FieldRange = &otsprotocol.FieldRange{}
	if g.FieldRange.Max != nil {
		vt, err := ToVariantValue(g.FieldRange.Max)
		if err != nil {
			return nil, err
		}
		pbGroupBy.FieldRange.Max = vt
	}
	if g.FieldRange.Min != nil {
		vt, err := ToVariantValue(g.FieldRange.Min)
		if err != nil {
			return nil, err
		}
		pbGroupBy.FieldRange.Min = vt
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
	subGroupBys, err := BuildPBForGroupBys(g.SubGroupByList)
	if err != nil {
		return nil, err
	}
	if subGroupBys != nil {
		pbGroupBy.SubGroupBys = subGroupBys
	}

	//sub groupbys
	subAggs, err := BuildPBForAggregations(g.SubAggList)
	if err != nil {
		return nil, err
	}
	if subAggs != nil {
		pbGroupBy.SubAggs = subAggs
	}

	data, err := proto.Marshal(pbGroupBy)
	return data, err
}

func (g *GroupByDateHistogram) ProtoBuffer() (*otsprotocol.GroupBy, error) {
	return BuildPBForGroupBy(g)
}

func (g *GroupByDateHistogram) SubAggregations(subAggregations ...Aggregation) *GroupByDateHistogram {
	g.SubAggList = subAggregations
	return g
}

func (g *GroupByDateHistogram) SubAggregation(subAggregation Aggregation) *GroupByDateHistogram {
	g.SubAggList = append(g.SubAggList, subAggregation)
	return g
}

func (g *GroupByDateHistogram) SubGroupBys(subGroupBys ...GroupBy) *GroupByDateHistogram {
	g.SubGroupByList = subGroupBys
	return g
}

func (g *GroupByDateHistogram) SubGroupBy(subGroupBy GroupBy) *GroupByDateHistogram {
	g.SubGroupByList = append(g.SubGroupByList, subGroupBy)
	return g
}

func (g *GroupByDateHistogram) SetName(name string) *GroupByDateHistogram {
	g.GroupByName = name
	return g
}

func (g *GroupByDateHistogram) SetTimeZone(timeZone string) *GroupByDateHistogram {
	g.TimeZone = &timeZone
	return g
}

func (g *GroupByDateHistogram) SetMissing(missing interface{}) *GroupByDateHistogram {
	g.Missing = missing
	return g
}

func (g *GroupByDateHistogram) SetMinDocCount(minDocCount int64) *GroupByDateHistogram {
	g.MinDocCount = &minDocCount
	return g
}

func (g *GroupByDateHistogram) SetInterval(interval model.DateTimeValue) *GroupByDateHistogram {
	g.Interval = interval
	return g
}

func (g *GroupByDateHistogram) SetFieldName(fieldName string) *GroupByDateHistogram {
	g.Field = fieldName
	return g
}

func (g *GroupByDateHistogram) SetFiledRange(min interface{}, max interface{}) *GroupByDateHistogram {
	g.FieldRange.Min = min
	g.FieldRange.Max = max
	return g
}

func (g *GroupByDateHistogram) SetGroupBySorters(sorters []GroupBySorter) *GroupByDateHistogram {
	g.Sorters = sorters
	return g
}

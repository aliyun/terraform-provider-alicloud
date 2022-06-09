package search

import (
	"errors"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/search/model"
	"github.com/golang/protobuf/proto"
)

type GroupByHistogram struct {
	GroupByName string

	Field       string
	Interval   	interface{}
	Missing     interface{}
	MinDocCount *int64
	FieldRange  model.FiledRange
	Sorters     []GroupBySorter

	SubAggList     []Aggregation
	SubGroupByList []GroupBy
}

func (g *GroupByHistogram) GetName() string {
	return g.GroupByName
}

func (g *GroupByHistogram) GetType() GroupByType {
	return GroupByHistogramType
}

func (g *GroupByHistogram) GetField() string {
	return g.Field
}

func (g *GroupByHistogram) Serialize() ([]byte, error) {

	pbGroupBy := &otsprotocol.GroupByHistogram{}

	pbGroupBy.FieldName = &g.Field
	pbGroupBy.MinDocCount = g.MinDocCount

	if g.Missing != nil {
		vt, err := ToVariantValue(g.Missing)
		if err != nil {
			return nil, err
		}
		pbGroupBy.Missing = []byte(vt)
	}
	if g.Interval != nil {
		vt, err := ToVariantValue(g.Interval)
		if err != nil {
			return nil, err
		}
		pbGroupBy.Interval = []byte(vt)
	}
	pbGroupBy.FieldRange = &otsprotocol.FieldRange{}
	if g.FieldRange.Max != nil {
		vt, err := ToVariantValue(g.FieldRange.Max)
		if err != nil {
			return nil, err
		}
		pbGroupBy.FieldRange.Max = []byte(vt)
	}
	if g.FieldRange.Min != nil {
		vt, err := ToVariantValue(g.FieldRange.Min)
		if err != nil {
			return nil, err
		}
		pbGroupBy.FieldRange.Min = []byte(vt)
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

func (g *GroupByHistogram) ProtoBuffer() (*otsprotocol.GroupBy, error) {
	return BuildPBForGroupBy(g)
}

func (g *GroupByHistogram) SubAggregations(subAggregations ...Aggregation) *GroupByHistogram {
	g.SubAggList = subAggregations
	return g
}

func (g *GroupByHistogram) SubAggregation(subAggregation Aggregation) *GroupByHistogram {
	g.SubAggList = append(g.SubAggList, subAggregation)
	return g
}

func (g *GroupByHistogram) SubGroupBys(subGroupBys ...GroupBy) *GroupByHistogram {
	g.SubGroupByList = subGroupBys
	return g
}

func (g *GroupByHistogram) SubGroupBy(subGroupBy GroupBy) *GroupByHistogram {
	g.SubGroupByList = append(g.SubGroupByList, subGroupBy)
	return g
}

func (g *GroupByHistogram) SetName(name string) *GroupByHistogram {
	g.GroupByName = name
	return g
}

func (g *GroupByHistogram) SetMissing(missing interface{}) *GroupByHistogram {
	g.Missing = missing
	return g
}

func (g *GroupByHistogram) SetMinDocCount(minDocCount int64) *GroupByHistogram {
	g.MinDocCount = &minDocCount
	return g
}

func (g *GroupByHistogram) SetInterval(interval interface{}) *GroupByHistogram {
	g.Interval = interval
	return g
}

func (g *GroupByHistogram) SetFieldName(fieldName string) *GroupByHistogram {
	g.Field = fieldName
	return g
}

func (g *GroupByHistogram) SetFiledRange(min interface{}, max interface{}) *GroupByHistogram {
	g.FieldRange.Min = min
	g.FieldRange.Max = max
	return g
}

func (g *GroupByHistogram) SetGroupBySorters(sorters []GroupBySorter) *GroupByHistogram {
	g.Sorters = sorters
	return g
}

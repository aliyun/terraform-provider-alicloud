package search

import (
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/golang/protobuf/proto"
)

// GroupByComposite 支持组合多种GroupBy以扁平模式返回分组结果，并且可通过NextToken翻页。
// 目前source中可组合的GroupBy类型包括：GroupByField/GroupByHistogram/GroupByDateHistogram。
// source中的GroupBy存在参数限制。
type GroupByComposite struct {
	GroupByName       string
	SourceGroupByList []GroupBy
	Size              *int32
	SuggestedSize     *int32
	NextToken         *string
	SubAggList        []Aggregation
	SubGroupByList    []GroupBy
}

func (g *GroupByComposite) GetName() string {
	return g.GroupByName
}

func (g *GroupByComposite) GetType() GroupByType {
	return GroupByCompositeType
}

func (g *GroupByComposite) Serialize() ([]byte, error) {
	pbGroupBy := &otsprotocol.GroupByComposite{}

	//size
	if g.Size != nil {
		pbGroupBy.Size = g.Size
	}

	if g.SuggestedSize != nil {
		pbGroupBy.SuggestedSize = g.SuggestedSize
	}

	// nextToken
	if g.NextToken != nil {
		pbGroupBy.NextToken = g.NextToken
	}

	// sources
	sourceGroupBys, err := BuildPBForGroupBys(g.SourceGroupByList)
	if err != nil {
		return nil, err
	}
	if sourceGroupBys != nil {
		pbGroupBy.Sources = sourceGroupBys
	}

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

func (g *GroupByComposite) ProtoBuffer() (*otsprotocol.GroupBy, error) {
	return BuildPBForGroupBy(g)
}

func (g *GroupByComposite) Name(aggName string) *GroupByComposite {
	g.GroupByName = aggName
	return g
}

func (g *GroupByComposite) SourceGroupBy(groupBy GroupBy) *GroupByComposite {
	g.SourceGroupByList = append(g.SourceGroupByList, groupBy)
	return g
}

func (g *GroupByComposite) SourceGroupBys(groupBys ...GroupBy) *GroupByComposite {
	g.SourceGroupByList = groupBys
	return g
}

func (g *GroupByComposite) SetSize(size int32) *GroupByComposite {
	g.Size = &size
	return g
}

func (g *GroupByComposite) SetSuggestedSize(suggestedSize int32) *GroupByComposite {
	g.SuggestedSize = &suggestedSize
	return g
}

func (g *GroupByComposite) SetNextToken(nextToken *string) *GroupByComposite {
	g.NextToken = nextToken
	return g
}

func (g *GroupByComposite) SubAggregations(subAggregations ...Aggregation) *GroupByComposite {
	g.SubAggList = subAggregations
	return g
}

func (g *GroupByComposite) SubAggregation(subAggregation Aggregation) *GroupByComposite {
	g.SubAggList = append(g.SubAggList, subAggregation)
	return g
}

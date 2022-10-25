package search

import (
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/search/model"
	"github.com/golang/protobuf/proto"
)

type GroupByResult interface {
	GetName() string
	GetType() GroupByType
}

type GroupByResults struct {
	resultMap map[string]GroupByResult
}

func (g *GroupByResults) GetRawResults() map[string]GroupByResult {
	m := make(map[string]GroupByResult, len(g.resultMap))
	for k, v := range g.resultMap {
		m[k] = v
	}
	return m
}

func (g *GroupByResults) Put(name string, result GroupByResult) {
	if g.resultMap == nil {
		g.resultMap = make(map[string]GroupByResult)
	}
	g.resultMap[name] = result
}

func (g GroupByResults) GroupByField(name string) (*GroupByFieldResult, error){
	if result, ok := g.resultMap[name]; ok {
		if result.GetType() != GroupByFieldType {
			return nil, errors.New(fmt.Sprintf("wrong group by type: [%v] needed, [%v] provided", result.GetType().String(), GroupByFieldType.String()))
		}
		return result.(*GroupByFieldResult), nil
	}
	return nil, errors.New(fmt.Sprintf("group by [%v] not found", name))
}

func (g GroupByResults) GroupByRange(name string) (*GroupByRangeResult, error){
	if result, ok := g.resultMap[name]; ok {
		if result.GetType() != GroupByRangeType {
			return nil, errors.New(fmt.Sprintf("wrong group by type: [%v] needed, [%v] provided", result.GetType().String(), GroupByRangeType.String()))
		}
		return result.(*GroupByRangeResult), nil
	}
	return nil, errors.New(fmt.Sprintf("group by [%v] not found", name))
}

func (g GroupByResults) GroupByFilter(name string) (*GroupByFilterResult, error){
	if result, ok := g.resultMap[name]; ok {
		if result.GetType() != GroupByFilterType {
			return nil, errors.New(fmt.Sprintf("wrong group by type: [%v] needed, [%v] provided", result.GetType().String(), GroupByFilterType.String()))
		}
		return result.(*GroupByFilterResult), nil
	}
	return nil, errors.New(fmt.Sprintf("group by [%v] not found", name))
}

func (g GroupByResults) GroupByGeoDistance(name string) (*GroupByGeoDistanceResult, error){
	if result, ok := g.resultMap[name]; ok {
		if result.GetType() != GroupByGeoDistanceType {
			return nil, errors.New(fmt.Sprintf("wrong group by type: [%v] needed, [%v] provided", result.GetType().String(), GroupByGeoDistanceType.String()))
		}
		return result.(*GroupByGeoDistanceResult), nil
	}
	return nil, errors.New(fmt.Sprintf("group by [%v] not found", name))
}

func (g GroupByResults) GroupByHistogram(name string) (*GroupByHistogramResult, error){
	if result, ok := g.resultMap[name]; ok {
		if result.GetType() != GroupByHistogramType {
			return nil, errors.New(fmt.Sprintf("wrong group by type: [%v] needed, [%v] provided", result.GetType().String(), GroupByHistogramType.String()))
		}
		return result.(*GroupByHistogramResult), nil
	}
	return nil, errors.New(fmt.Sprintf("group by [%v] not found", name))
}

func (g GroupByResults) Empty() bool {
	return len(g.resultMap) == 0
}


func ParseGroupByFieldResultFromPB(pbGroupByResult *otsprotocol.GroupByResult) (*GroupByFieldResult, error) {
	groupByResult := new(GroupByFieldResult)
	groupByResult.Name = *pbGroupByResult.Name

	pbGroupByResultBody := new(otsprotocol.GroupByFieldResult)
	err := proto.Unmarshal(pbGroupByResult.GroupByResult, pbGroupByResultBody)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to parse group by body: %v", err.Error()))
	}
	pbItems := pbGroupByResultBody.GroupByFieldResultItems

	if len(pbItems) > 0 {
		var items []GroupByFieldResultItem

		for _, pbItem := range pbItems {
			item := GroupByFieldResultItem{}
			item.Key = *pbItem.Key
			item.RowCount = *pbItem.RowCount
			if pbItem.SubAggsResult != nil && len(pbItem.SubAggsResult.AggResults) > 0 {
				subAggResults, err := ParseAggregationResultsFromPB(pbItem.SubAggsResult.AggResults)
				if err != nil {
					return nil, err
				}
				item.SubAggregations = *subAggResults
			}
			if pbItem.SubGroupBysResult != nil && len(pbItem.SubGroupBysResult.GroupByResults) > 0 {
				subGroupByResults, err :=  ParseGroupByResultsFromPB(pbItem.SubGroupBysResult.GroupByResults)
				if err != nil {
					return nil, err
				}
				item.SubGroupBys = *subGroupByResults
			}
			items = append(items, item)
		}
		groupByResult.Items = items
	}
	return groupByResult, nil
}

func ParseGroupByRangeResultFromPB(pbGroupByResult *otsprotocol.GroupByResult) (*GroupByRangeResult, error) {
	groupByResult := new(GroupByRangeResult)
	groupByResult.Name = *pbGroupByResult.Name

	pbGroupByResultBody := new(otsprotocol.GroupByRangeResult)
	err := proto.Unmarshal(pbGroupByResult.GroupByResult, pbGroupByResultBody)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to parse group by body: %v", err.Error()))
	}
	pbItems := pbGroupByResultBody.GroupByRangeResultItems

	if len(pbItems) > 0 {
		var items []GroupByRangeResultItem

		for _, pbItem := range pbItems {
			item := GroupByRangeResultItem{}
			item.From = *pbItem.From
			item.To = *pbItem.To
			item.RowCount = *pbItem.RowCount
			if pbItem.SubAggsResult != nil && len(pbItem.SubAggsResult.AggResults) > 0 {
				subAggResults, err := ParseAggregationResultsFromPB(pbItem.SubAggsResult.AggResults)
				if err != nil {
					return nil, err
				}
				item.SubAggregations = *subAggResults
			}
			if pbItem.SubGroupBysResult != nil && len(pbItem.SubGroupBysResult.GroupByResults) > 0 {
				subGroupByResults, err :=  ParseGroupByResultsFromPB(pbItem.SubGroupBysResult.GroupByResults)
				if err != nil {
					return nil, err
				}
				item.SubGroupBys = *subGroupByResults
			}
			items = append(items, item)
		}
		groupByResult.Items = items
	}
	return groupByResult, nil
}

func ParseGroupByFilterResultFromPB(pbGroupByResult *otsprotocol.GroupByResult) (*GroupByFilterResult, error) {
	groupByResult := new(GroupByFilterResult)
	groupByResult.Name = *pbGroupByResult.Name

	pbGroupByResultBody := new(otsprotocol.GroupByFilterResult)
	err := proto.Unmarshal(pbGroupByResult.GroupByResult, pbGroupByResultBody)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to parse group by body: %v", err.Error()))
	}
	pbItems := pbGroupByResultBody.GroupByFilterResultItems

	if len(pbItems) > 0 {
		var items []GroupByFilterResultItem

		for _, pbItem := range pbItems {
			item := GroupByFilterResultItem{}
			item.RowCount = *pbItem.RowCount
			if pbItem.SubAggsResult != nil && len(pbItem.SubAggsResult.AggResults) > 0 {
				subAggResults, err := ParseAggregationResultsFromPB(pbItem.SubAggsResult.AggResults)
				if err != nil {
					return nil, err
				}
				item.SubAggregations = *subAggResults
			}
			if pbItem.SubGroupBysResult != nil && len(pbItem.SubGroupBysResult.GroupByResults) > 0 {
				subGroupByResults, err :=  ParseGroupByResultsFromPB(pbItem.SubGroupBysResult.GroupByResults)
				if err != nil {
					return nil, err
				}
				item.SubGroupBys = *subGroupByResults
			}
			items = append(items, item)
		}
		groupByResult.Items = items
	}
	return groupByResult, nil
}

func ParseGroupByGeoDistanceResultFromPB(pbGroupByResult *otsprotocol.GroupByResult) (*GroupByGeoDistanceResult, error) {
	groupByResult := new(GroupByGeoDistanceResult)
	groupByResult.Name = *pbGroupByResult.Name

	pbGroupByResultBody := new(otsprotocol.GroupByGeoDistanceResult)
	err := proto.Unmarshal(pbGroupByResult.GroupByResult, pbGroupByResultBody)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to parse group by body: %v", err.Error()))
	}
	pbItems := pbGroupByResultBody.GroupByGeoDistanceResultItems

	if len(pbItems) > 0 {
		var items []GroupByGeoDistanceResultItem

		for _, pbItem := range pbItems {
			item := GroupByGeoDistanceResultItem{}
			item.From = *pbItem.From
			item.To = *pbItem.To
			item.RowCount = *pbItem.RowCount
			if pbItem.SubAggsResult != nil && len(pbItem.SubAggsResult.AggResults) > 0 {
				subAggResults, err := ParseAggregationResultsFromPB(pbItem.SubAggsResult.AggResults)
				if err != nil {
					return nil, err
				}
				item.SubAggregations = *subAggResults
			}
			if pbItem.SubGroupBysResult != nil && len(pbItem.SubGroupBysResult.GroupByResults) > 0 {
				subGroupByResults, err :=  ParseGroupByResultsFromPB(pbItem.SubGroupBysResult.GroupByResults)
				if err != nil {
					return nil, err
				}
				item.SubGroupBys = *subGroupByResults
			}
			items = append(items, item)
		}
		groupByResult.Items = items
	}
	return groupByResult, nil
}

func ParseGroupByHistogramResultFromPB(pbGroupByResult *otsprotocol.GroupByResult) (*GroupByHistogramResult, error) {
	groupByResult := new(GroupByHistogramResult)
	groupByResult.Name = *pbGroupByResult.Name

	pbGroupByResultBody := new(otsprotocol.GroupByHistogramResult)
	err := proto.Unmarshal(pbGroupByResult.GroupByResult, pbGroupByResultBody)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to parse group by body: %v", err.Error()))
	}
	pbItems := pbGroupByResultBody.GroupByHistograItems

	if len(pbItems) > 0 {
		var items []GroupByHistogramItem

		for _, pbItem := range pbItems {
			item := GroupByHistogramItem{}
			var err error
			var valuePtr *model.ColumnValue
			valuePtr, err = ForceConvertToDestColumnValue(pbItem.Key)
			if err != nil {
				return nil, err
			}
			item.Key = *valuePtr
			item.Value = *pbItem.Value
			if pbItem.SubAggsResult != nil && len(pbItem.SubAggsResult.AggResults) > 0 {
				subAggResults, err := ParseAggregationResultsFromPB(pbItem.SubAggsResult.AggResults)
				if err != nil {
					return nil, err
				}
				item.SubAggregations = *subAggResults
			}
			if pbItem.SubGroupBysResult != nil && len(pbItem.SubGroupBysResult.GroupByResults) > 0 {
				subGroupByResults, err := ParseGroupByResultsFromPB(pbItem.SubGroupBysResult.GroupByResults)
				if err != nil {
					return nil, err
				}
				item.SubGroupBys = *subGroupByResults
			}
			items = append(items, item)
		}
		groupByResult.Items = items
	}
	return groupByResult, nil
}

func ParseGroupByResultsFromPB(pbGroupByResults []*otsprotocol.GroupByResult) (*GroupByResults, error) {
	groupByResults := GroupByResults{}
	for _, pbGroupByResult := range pbGroupByResults {
		switch pbGroupByResult.GetType() {
		case otsprotocol.GroupByType_GROUP_BY_FIELD:
			groupByResult, err := ParseGroupByFieldResultFromPB(pbGroupByResult)
			if err != nil {
				return nil, err
			}
			groupByResults.Put(pbGroupByResult.GetName(), groupByResult)
			break
		case otsprotocol.GroupByType_GROUP_BY_RANGE:
			groupByResult, err := ParseGroupByRangeResultFromPB(pbGroupByResult)
			if err != nil {
				return nil, err
			}
			groupByResults.Put(groupByResult.Name, groupByResult)
			break
		case otsprotocol.GroupByType_GROUP_BY_FILTER:
			groupByResult, err := ParseGroupByFilterResultFromPB(pbGroupByResult)
			if err != nil {
				return nil, err
			}
			groupByResults.Put(groupByResult.Name, groupByResult)
			break
		case otsprotocol.GroupByType_GROUP_BY_GEO_DISTANCE:
			groupByResult, err := ParseGroupByGeoDistanceResultFromPB(pbGroupByResult)
			if err != nil {
				return nil, err
			}
			groupByResults.Put(groupByResult.Name, groupByResult)
			break
		case otsprotocol.GroupByType_GROUP_BY_HISTOGRAM:
			groupByResult, err := ParseGroupByHistogramResultFromPB(pbGroupByResult)
			if err != nil {
				return nil, err
			}
			groupByResults.Put(groupByResult.Name, groupByResult)
			break
		default:
			return nil, errors.New(fmt.Sprintf("unknown group by result type: %v", pbGroupByResult.GetType()))
		}
	}

	return &groupByResults, nil
}

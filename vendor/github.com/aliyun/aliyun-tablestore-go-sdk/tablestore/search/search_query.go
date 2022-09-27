package search

import (
	"encoding/json"
	"errors"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/golang/protobuf/proto"
)

type SearchQuery interface {
	Serialize() ([]byte, error)
}

type queryAlias struct {
	Name string
	Query Query
}

type aggregationAlias struct {
	Name string
	Aggregation Aggregation
}

type searchQuery struct {
	Offset        int32
	Limit         int32
	Query         Query `json:"-"`
	Collapse      *Collapse
	Sort          *Sort
	GetTotalCount bool
	Token         []byte
	Aggregations  []Aggregation `json:"-"`
	GroupBys      []GroupBy `json:"-"`

	// for json marshal and unmarshal
	QueryAlias    queryAlias `json:"Query"`
	AggregationAlias []aggregationAlias `json:"Aggregations"`
}

func (q *searchQuery) MarshalJSON() (data []byte, err error) {
	type searchQueryAlias searchQuery
	query := searchQueryAlias(*q)
	if q.Query != nil {
		query.QueryAlias = queryAlias{
			Name:  q.Query.Type().String(),
			Query: q.Query,
		}
	}

	if q.Aggregations != nil {
		aggs := make([]aggregationAlias, 0)
		for _, agg := range q.Aggregations {
			aggs = append(aggs, aggregationAlias{
				Name:        agg.GetType().String(),
				Aggregation: agg,
			})
		}

		query.AggregationAlias = aggs
	}

	data, err = json.Marshal(query)
	return
}

func (q *searchQuery) UnmarshalJSON(data []byte) (err error) {
	type searchQueryAlias searchQuery
	sqAlias := &searchQueryAlias{}
	err = json.Unmarshal(data, sqAlias)
	if err != nil {
		return
	}

	q.Offset = sqAlias.Offset
	q.Limit = sqAlias.Limit
	q.Query = sqAlias.QueryAlias.Query
	q.Collapse = sqAlias.Collapse
	q.Sort = sqAlias.Sort
	q.GetTotalCount = sqAlias.GetTotalCount
	q.Token = sqAlias.Token

	// aggregations
	if sqAlias.AggregationAlias != nil {
		aggs := make([]Aggregation, 0)
		for _, agg := range sqAlias.AggregationAlias {
			aggs = append(aggs, agg.Aggregation)
		}

		q.Aggregations = aggs
	}

	return
}

func (q *queryAlias) UnmarshalJSON(data []byte) (err error) {
	jm := make(map[string]json.RawMessage)
	err = json.Unmarshal(data, &jm)
	if err != nil {
		return
	}

	nameRM, ok := jm["Name"]
	if !ok {
		err = errors.New("Field 'Name' is missing.")
		return
	}

	var name string
	err = json.Unmarshal(nameRM, &name)
	if err != nil {
		return
	}

	query, ok := jm["Query"]
	if !ok {
		err = errors.New("Field 'Query' is missing.")
		return
	}

	q.Name = name
	q.Query, err = UnmarshalQuery(name, query)
	if err != nil {
		return
	}

	return
}

func (q *aggregationAlias) UnmarshalJSON(data []byte) (err error) {
	jm := make(map[string]json.RawMessage)
	err = json.Unmarshal(data, &jm)
	if err != nil {
		return
	}

	nameRM, ok := jm["Name"]
	if !ok {
		err = errors.New("Field 'Name' is missing.")
		return
	}

	var name string
	err = json.Unmarshal(nameRM, &name)
	if err != nil {
		return
	}

	agg, ok := jm["Aggregation"]
	if !ok {
		err = errors.New("Field 'Aggregation' is missing.")
		return
	}

	q.Name = name
	q.Aggregation, err = UnmarshalAggregation(name, agg)
	if err != nil {
		return
	}

	return
}

func NewSearchQuery() *searchQuery {
	return &searchQuery {
		Offset:        -1,
		Limit:         -1,
		GetTotalCount: false,
	}
}

func (s *searchQuery) SetOffset(offset int32) *searchQuery {
	s.Offset = offset
	return s
}

func (s *searchQuery) SetLimit(limit int32) *searchQuery {
	s.Limit = limit
	return s
}

func (s *searchQuery) SetQuery(query Query) *searchQuery {
	s.Query = query
	return s
}

func NewAvgAggregation(name string, fieldName string) *AvgAggregation {
	return &AvgAggregation {
		AggName: name,
		Field: fieldName,
	}
}

func NewDistinctCountAggregation(name string, fieldName string) *DistinctCountAggregation {
	return &DistinctCountAggregation {
		AggName: name,
		Field: fieldName,
	}
}

func NewMaxAggregation(name string, fieldName string) *MaxAggregation {
	return &MaxAggregation {
		AggName: name,
		Field: fieldName,
	}
}

func NewMinAggregation(name string, fieldName string) *MinAggregation {
	return &MinAggregation {
		AggName: name,
		Field: fieldName,
	}
}

func NewSumAggregation(name string, fieldName string) *SumAggregation {
	return &SumAggregation {
		AggName: name,
		Field: fieldName,
	}
}

func NewCountAggregation(name string, fieldName string) *CountAggregation {
	return &CountAggregation {
		AggName: name,
		Field: fieldName,
	}
}

func NewTopRowsAggregation(name string) *TopRowsAggregation {
	return &TopRowsAggregation {
		AggName: name,
	}
}

func NewPercentilesAggregation(name string, filedName string) *PercentilesAggregation {
	return &PercentilesAggregation {
		AggName: name,
		Field: 	 filedName,
	}
}

//
func NewGroupByField(name string, fieldName string) *GroupByField {
	return &GroupByField {
		AggName: name,
		Field:   fieldName,
	}
}

func NewGroupByRange(name string, fieldName string) *GroupByRange {
	return &GroupByRange {
		AggName: name,
		Field:   fieldName,
	}
}

func NewGroupByFilter(name string) *GroupByFilter {
	return &GroupByFilter {
		AggName: name,
	}
}

func NewGroupByGeoDistance(name string, fieldName string, origin GeoPoint) *GroupByGeoDistance {
	return &GroupByGeoDistance {
		AggName: name,
		Field:   fieldName,
		Origin:  origin,
	}
}

func NewGroupByHistogram(name string, filedName string) *GroupByHistogram {
	return &GroupByHistogram{
		GroupByName: name,
		Field:       filedName,
	}
}

func (s *searchQuery) Aggregation(agg ...Aggregation) *searchQuery {
	for i := 0; i < len(agg); i++ {
		s.Aggregations = append(s.Aggregations, agg[i])
	}
	return s
}

func (s *searchQuery) GroupBy(groupBy ...GroupBy) *searchQuery {
	for i := 0; i < len(groupBy); i++ {
		s.GroupBys = append(s.GroupBys, groupBy[i])
	}
	return s
}

func (s *searchQuery) SetCollapse(collapse *Collapse) *searchQuery {
	s.Collapse = collapse
	return s
}

func (s *searchQuery) SetSort(sort *Sort) *searchQuery {
	s.Sort = sort
	return s
}

func (s *searchQuery) SetGetTotalCount(getTotalCount bool) *searchQuery {
	s.GetTotalCount = getTotalCount
	return s
}

func (s *searchQuery) SetToken(token []byte) *searchQuery {
	s.Token = token
	s.Sort = nil
	return s
}

func (s *searchQuery) Serialize() ([]byte, error) {
	searchQuery := &otsprotocol.SearchQuery{}
	if s.Offset >= 0 {
		searchQuery.Offset = &s.Offset
	}
	if s.Limit >= 0 {
		searchQuery.Limit = &s.Limit
	}
	if s.Query != nil {
		pbQuery, err := s.Query.ProtoBuffer()
		if err != nil {
			return nil, err
		}
		searchQuery.Query = pbQuery
	}
	if s.Collapse != nil {
		pbCollapse, err := s.Collapse.ProtoBuffer()
		if err != nil {
			return nil, err
		}
		searchQuery.Collapse = pbCollapse
	}
	if s.Sort != nil {
		pbSort, err := s.Sort.ProtoBuffer()
		if err != nil {
			return nil, err
		}
		searchQuery.Sort = pbSort
	}
	searchQuery.GetTotalCount = &s.GetTotalCount
	if s.Token != nil && len(s.Token) > 0 {
		searchQuery.Token = s.Token
	}

	if len(s.Aggregations) > 0 {
		pbAggregations := new(otsprotocol.Aggregations)
		for _, aggregation := range s.Aggregations {
			pbAggregation, err := aggregation.ProtoBuffer()
			if err != nil {
				return nil, err
			}
			pbAggregations.Aggs = append(pbAggregations.Aggs, pbAggregation)
		}
		searchQuery.Aggs = pbAggregations
	}

	if len(s.GroupBys) > 0 {
		pbGroupBys := new(otsprotocol.GroupBys)
		for _, groupBy := range s.GroupBys {
			pbGroupBy, err := groupBy.ProtoBuffer()
			if err != nil {
				return nil, err
			}
			pbGroupBys.GroupBys = append(pbGroupBys.GroupBys, pbGroupBy)
		}
		searchQuery.GroupBys = pbGroupBys
	}

	data, err := proto.Marshal(searchQuery)
	return data, err
}

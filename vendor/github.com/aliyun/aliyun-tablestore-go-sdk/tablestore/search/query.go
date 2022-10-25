package search

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
)

type QueryType int

const (
	QueryType_None                QueryType = 0
	QueryType_MatchQuery          QueryType = 1
	QueryType_MatchPhraseQuery    QueryType = 2
	QueryType_TermQuery           QueryType = 3
	QueryType_RangeQuery          QueryType = 4
	QueryType_PrefixQuery         QueryType = 5
	QueryType_BoolQuery           QueryType = 6
	QueryType_ConstScoreQuery     QueryType = 7
	QueryType_FunctionScoreQuery  QueryType = 8
	QueryType_NestedQuery         QueryType = 9
	QueryType_WildcardQuery       QueryType = 10
	QueryType_MatchAllQuery       QueryType = 11
	QueryType_GeoBoundingBoxQuery QueryType = 12
	QueryType_GeoDistanceQuery    QueryType = 13
	QueryType_GeoPolygonQuery     QueryType = 14
	QueryType_TermsQuery          QueryType = 15
	QueryType_ExistsQuery         QueryType = 16
)

func (q QueryType) String() string {
	switch q {
	case QueryType_MatchQuery:
		return "MatchQuery"
	case QueryType_MatchPhraseQuery:
		return "MatchPhraseQuery"
	case QueryType_TermQuery:
		return "TermQuery"
	case QueryType_RangeQuery:
		return "RangeQuery"
	case QueryType_PrefixQuery:
		return "PrefixQuery"
	case QueryType_BoolQuery:
		return "BoolQuery"
	case QueryType_ConstScoreQuery:
		return "ConstScoreQuery"
	case QueryType_FunctionScoreQuery:
		return "FunctionScoreQuery"
	case QueryType_NestedQuery:
		return "NestedQuery"
	case QueryType_WildcardQuery:
		return "WildcardQuery"
	case QueryType_MatchAllQuery:
		return "MatchAllQuery"
	case QueryType_GeoBoundingBoxQuery:
		return "GeoBoundingBoxQuery"
	case QueryType_GeoDistanceQuery:
		return "GeoDistanceQuery"
	case QueryType_GeoPolygonQuery:
		return "GeoPolygonQuery"
	case QueryType_TermsQuery:
		return "TermsQuery"
	case QueryType_ExistsQuery:
		return "ExistsQuery"
	}

	return ""
}

func UnmarshalQuery(name string, data json.RawMessage) (Query, error) {
	var err error
	switch name {
	case "MatchQuery":
		q := &MatchQuery{}
		err = json.Unmarshal(data, q)
		return q, err
	case "MatchPhraseQuery":
		q := &MatchPhraseQuery{}
		err = json.Unmarshal(data, q)
		return q, err
	case "TermQuery":
		q := &TermQuery{}
		err = json.Unmarshal(data, q)
		return q, err
	case "RangeQuery":
		q := &RangeQuery{}
		err = json.Unmarshal(data, q)
		return q, err
	case "PrefixQuery":
		q := &PrefixQuery{}
		err = json.Unmarshal(data, q)
		return q, err
	case "BoolQuery":
		q := &BoolQuery{}
		err = json.Unmarshal(data, q)
		return q, err
	case "ConstScoreQuery":
		q := &ConstScoreQuery{}
		err = json.Unmarshal(data, q)
		return q, err
	case "FunctionScoreQuery":
		q := &FunctionScoreQuery{}
		err = json.Unmarshal(data, q)
		return q, err
	case "NestedQuery":
		q := &NestedQuery{}
		err = json.Unmarshal(data, q)
		return q, err
	case "WildcardQuery":
		q := &WildcardQuery{}
		err = json.Unmarshal(data, q)
		return q, err
	case "MatchAllQuery":
		q := &MatchAllQuery{}
		err = json.Unmarshal(data, q)
		return q, err
	case "GeoBoundingBoxQuery":
		q := &GeoBoundingBoxQuery{}
		err = json.Unmarshal(data, q)
		return q, err
	case "GeoDistanceQuery":
		q := &GeoDistanceQuery{}
		err = json.Unmarshal(data, q)
		return q, err
	case "GeoPolygonQuery":
		q := &GeoPolygonQuery{}
		err = json.Unmarshal(data, q)
		return q, err
	case "TermsQuery":
		q := &TermsQuery{}
		err = json.Unmarshal(data, q)
		return q, err
	case "ExistsQuery":
		q := &ExistsQuery{}
		err = json.Unmarshal(data, q)
		return q, err
	}

	return nil, errors.New(fmt.Sprintf("Unknown query type: %s.", name))
}

func ToQueryType(q string) QueryType {
	switch q {
	case "MatchQuery":
		return QueryType_MatchQuery
	case "MatchPhraseQuery":
		return QueryType_MatchPhraseQuery
	case "TermQuery":
		return QueryType_TermQuery
	case "RangeQuery":
		return QueryType_RangeQuery
	case "PrefixQuery":
		return QueryType_PrefixQuery
	case "BoolQuery":
		return QueryType_BoolQuery
	case "ConstScoreQuery":
		return QueryType_ConstScoreQuery
	case "FunctionScoreQuery":
		return QueryType_FunctionScoreQuery
	case "NestedQuery":
		return QueryType_NestedQuery
	case "WildcardQuery":
		return QueryType_WildcardQuery
	case "MatchAllQuery":
		return QueryType_MatchAllQuery
	case "GeoBoundingBoxQuery":
		return QueryType_GeoBoundingBoxQuery
	case "GeoDistanceQuery":
		return QueryType_GeoDistanceQuery
	case "GeoPolygonQuery":
		return QueryType_GeoPolygonQuery
	case "TermsQuery":
		return QueryType_TermsQuery
	case "ExistsQuery":
		return QueryType_ExistsQuery
	}

	return QueryType_None
}

func (q QueryType) Enum() *QueryType {
	newQuery := q
	return &newQuery
}

func (q QueryType) ToPB() *otsprotocol.QueryType {
	switch q {
	case QueryType_None:
		return nil
	case QueryType_MatchQuery:
		return otsprotocol.QueryType_MATCH_QUERY.Enum()
	case QueryType_MatchPhraseQuery:
		return otsprotocol.QueryType_MATCH_PHRASE_QUERY.Enum()
	case QueryType_TermQuery:
		return otsprotocol.QueryType_TERM_QUERY.Enum()
	case QueryType_RangeQuery:
		return otsprotocol.QueryType_RANGE_QUERY.Enum()
	case QueryType_PrefixQuery:
		return otsprotocol.QueryType_PREFIX_QUERY.Enum()
	case QueryType_BoolQuery:
		return otsprotocol.QueryType_BOOL_QUERY.Enum()
	case QueryType_ConstScoreQuery:
		return otsprotocol.QueryType_CONST_SCORE_QUERY.Enum()
	case QueryType_FunctionScoreQuery:
		return otsprotocol.QueryType_FUNCTION_SCORE_QUERY.Enum()
	case QueryType_NestedQuery:
		return otsprotocol.QueryType_NESTED_QUERY.Enum()
	case QueryType_WildcardQuery:
		return otsprotocol.QueryType_WILDCARD_QUERY.Enum()
	case QueryType_MatchAllQuery:
		return otsprotocol.QueryType_MATCH_ALL_QUERY.Enum()
	case QueryType_GeoBoundingBoxQuery:
		return otsprotocol.QueryType_GEO_BOUNDING_BOX_QUERY.Enum()
	case QueryType_GeoDistanceQuery:
		return otsprotocol.QueryType_GEO_DISTANCE_QUERY.Enum()
	case QueryType_GeoPolygonQuery:
		return otsprotocol.QueryType_GEO_POLYGON_QUERY.Enum()
	case QueryType_TermsQuery:
		return otsprotocol.QueryType_TERMS_QUERY.Enum()
	case QueryType_ExistsQuery:
		return otsprotocol.QueryType_EXISTS_QUERY.Enum()
	default:
		panic("unexpected")
	}
}

type Query interface {
	Type() QueryType
	Serialize() ([]byte, error)
	ProtoBuffer() (*otsprotocol.Query, error)
}

func BuildPBForQuery(q Query) (*otsprotocol.Query, error) {
	query := &otsprotocol.Query{}
	query.Type = q.Type().ToPB()
	data, err := q.Serialize()
	if err != nil {
		return nil, err
	}
	query.Query = data
	return query, nil
}

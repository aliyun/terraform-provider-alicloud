package search

import (
	"encoding/json"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/golang/protobuf/proto"
)

type ScoreModeType int

const (
	ScoreMode_None  ScoreModeType = 1
	ScoreMode_Avg   ScoreModeType = 2
	ScoreMode_Max   ScoreModeType = 3
	ScoreMode_Total ScoreModeType = 4
	ScoreMode_Min   ScoreModeType = 5
)

type NestedQuery struct {
	Path      string
	Query     Query `json:"-"`
	ScoreMode ScoreModeType

	// for json marshal and unmarshal
	QueryAlias queryAlias `json:"Query"`
}

func (q *NestedQuery) MarshalJSON() ([]byte, error) {
	type NestedQueryAlias NestedQuery
	bqAlias := NestedQueryAlias(*q)
	if bqAlias.Query != nil {
		bqAlias.QueryAlias = queryAlias{
			Name:  q.Query.Type().String(),
			Query: q.Query,
		}
	}

	data, err := json.Marshal(bqAlias)
	return data, err
}

func (q *NestedQuery) UnmarshalJSON(data []byte) (err error) {
	type NestedQueryAlias NestedQuery
	bqAlias := &NestedQueryAlias{}
	err = json.Unmarshal(data, bqAlias)
	if err != nil {
		return
	}

	q.Query = bqAlias.QueryAlias.Query
	return
}

func (q *NestedQuery) Type() QueryType {
	return QueryType_NestedQuery
}

func (q *NestedQuery) Serialize() ([]byte, error) {
	query := &otsprotocol.NestedQuery{}
	pbQ, err := q.Query.ProtoBuffer()
	if err != nil {
		return nil, err
	}
	query.Query = pbQ
	query.Path = &q.Path
	switch q.ScoreMode {
	case ScoreMode_None:
		query.ScoreMode = otsprotocol.ScoreMode_SCORE_MODE_NONE.Enum()
	case ScoreMode_Avg:
		query.ScoreMode = otsprotocol.ScoreMode_SCORE_MODE_AVG.Enum()
	case ScoreMode_Max:
		query.ScoreMode = otsprotocol.ScoreMode_SCORE_MODE_MAX.Enum()
	case ScoreMode_Min:
		query.ScoreMode = otsprotocol.ScoreMode_SCORE_MODE_MIN.Enum()
	case ScoreMode_Total:
		query.ScoreMode = otsprotocol.ScoreMode_SCORE_MODE_TOTAL.Enum()
	}
	data, err := proto.Marshal(query)
	return data, err
}

func (q *NestedQuery) ProtoBuffer() (*otsprotocol.Query, error) {
	return BuildPBForQuery(q)
}

package search

import (
	"encoding/json"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/golang/protobuf/proto"
)

type ConstScoreQuery struct {
	Filter Query `json:"-"`

	// for json marshal and unmarshal
	FilterAlias queryAlias `json:"Filter"`
}

func (q *ConstScoreQuery) MarshalJSON() ([]byte, error) {
	type ConstScoreQueryAlias ConstScoreQuery
	bqAlias := ConstScoreQueryAlias(*q)
	if bqAlias.Filter != nil {
		bqAlias.FilterAlias = queryAlias{
			Name:  q.Filter.Type().String(),
			Query: q.Filter,
		}
	}

	data, err := json.Marshal(bqAlias)
	return data, err
}

func (q *ConstScoreQuery) UnmarshalJSON(data []byte) (err error) {
	type ConstScoreQueryAlias ConstScoreQuery
	bqAlias := &ConstScoreQueryAlias{}
	err = json.Unmarshal(data, bqAlias)
	if err != nil {
		return
	}

	q.Filter = bqAlias.FilterAlias.Query
	return
}

func (q *ConstScoreQuery) Type() QueryType {
	return QueryType_ConstScoreQuery
}

func (q *ConstScoreQuery) Serialize() ([]byte, error) {
	query := &otsprotocol.ConstScoreQuery{}
	pbQ, err := q.Filter.ProtoBuffer()
	if err != nil {
		return nil, err
	}
	query.Filter = pbQ
	data, err := proto.Marshal(query)
	return data, err
}

func (q *ConstScoreQuery) ProtoBuffer() (*otsprotocol.Query, error) {
	return BuildPBForQuery(q)
}

package search

import (
	"encoding/json"
	"errors"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/golang/protobuf/proto"
)

type FieldValueFactor struct {
	FieldName string
}

func (f *FieldValueFactor) ProtoBuffer() (*otsprotocol.FieldValueFactor, error) {
	pb := &otsprotocol.FieldValueFactor{}
	pb.FieldName = &f.FieldName
	return pb, nil
}

type FunctionScoreQuery struct {
	Query            Query `json:"-"`
	FieldValueFactor *FieldValueFactor

	// for json marshal and unmarshal
	QueryAlias queryAlias `json:"Query"`
}

func (q *FunctionScoreQuery) MarshalJSON() ([]byte, error) {
	type FunctionScoreQueryAlias FunctionScoreQuery
	bqAlias := FunctionScoreQueryAlias(*q)
	if bqAlias.Query != nil {
		bqAlias.QueryAlias = queryAlias{
			Name:  q.Query.Type().String(),
			Query: q.Query,
		}
	}

	data, err := json.Marshal(bqAlias)
	return data, err
}

func (q *FunctionScoreQuery) UnmarshalJSON(data []byte) (err error) {
	type FunctionScoreQueryAlias FunctionScoreQuery
	bqAlias := &FunctionScoreQueryAlias{}
	err = json.Unmarshal(data, bqAlias)
	if err != nil {
		return
	}

	q.Query = bqAlias.QueryAlias.Query
	return
}

func (q *FunctionScoreQuery) Type() QueryType {
	return QueryType_FunctionScoreQuery
}

func (q *FunctionScoreQuery) Serialize() ([]byte, error) {
	if q.Query == nil || q.FieldValueFactor == nil {
		return nil, errors.New("FunctionScoreQuery: Query or FieldValueFactor is nil")
	}
	query := &otsprotocol.FunctionScoreQuery{}
	pbQ, err := q.Query.ProtoBuffer()
	if err != nil {
		return nil, err
	}
	query.Query = pbQ
	pbF, err := q.FieldValueFactor.ProtoBuffer()
	if err != nil {
		return nil, err
	}
	query.FieldValueFactor = pbF
	data, err := proto.Marshal(query)
	return data, err
}

func (q *FunctionScoreQuery) ProtoBuffer() (*otsprotocol.Query, error) {
	return BuildPBForQuery(q)
}

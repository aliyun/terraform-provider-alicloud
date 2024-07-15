package search

import (
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/golang/protobuf/proto"
)

type KnnVectorQuery struct {
	FieldName          string
	TopK               *int32
	Float32QueryVector []float32
	Filter             Query
}

func (q *KnnVectorQuery) Type() QueryType {
	return QueryType_KnnVectorQuery
}

func (q *KnnVectorQuery) Serialize() ([]byte, error) {
	query := &otsprotocol.KnnVectorQuery{}
	query.FieldName = proto.String(q.FieldName)
	query.TopK = q.TopK
	query.Float32QueryVector = q.Float32QueryVector
	
	if q.Filter != nil {
		if querySerialize, err := q.Filter.ProtoBuffer(); err != nil {
			return nil, err
		} else {
			query.Filter = querySerialize
		}

	}

	return proto.Marshal(query)
}

func (q *KnnVectorQuery) ProtoBuffer() (*otsprotocol.Query, error) {
	return BuildPBForQuery(q)
}

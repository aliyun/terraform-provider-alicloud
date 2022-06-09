package search

import (
	"encoding/json"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/golang/protobuf/proto"
)

type BoolQuery struct {
	MustQueries        []Query `json:"-"`
	MustNotQueries     []Query `json:"-"`
	FilterQueries      []Query `json:"-"`
	ShouldQueries      []Query `json:"-"`
	MinimumShouldMatch *int32

	// for json marshal and unmarshal
	MustQueriesAlias   []queryAlias `json:"MustQueries"`
	MustNotQueriesAlias[]queryAlias `json:"MustNotQueries"`
	FilterQueriesAlias []queryAlias `json:"FilterQueries"`
	ShouldQueriesAlias []queryAlias `json:"ShouldQueries"`
}

func (q *BoolQuery) MarshalJSON() ([]byte, error) {
	type BoolQueryAlias BoolQuery
	bqAlias := BoolQueryAlias(*q)
	if q.MustQueries != nil {
		mqs := make([]queryAlias, 0)
		for _, tq := range q.MustQueries {
			mqs = append(mqs, queryAlias{
				Name: tq.Type().String(),
				Query: tq,
			})
		}

		bqAlias.MustQueriesAlias = mqs
	}

	if q.MustNotQueries != nil {
		mqs := make([]queryAlias, 0)
		for _, tq := range q.MustNotQueries {
			mqs = append(mqs, queryAlias{
				Name: tq.Type().String(),
				Query: tq,
			})
		}

		bqAlias.MustNotQueriesAlias = mqs
	}

	if q.FilterQueries != nil {
		mqs := make([]queryAlias, 0)
		for _, tq := range q.FilterQueries {
			mqs = append(mqs, queryAlias{
				Name: tq.Type().String(),
				Query: tq,
			})
		}

		bqAlias.FilterQueriesAlias = mqs
	}

	if q.ShouldQueries != nil {
		mqs := make([]queryAlias, 0)
		for _, tq := range q.ShouldQueries {
			mqs = append(mqs, queryAlias{
				Name: tq.Type().String(),
				Query: tq,
			})
		}

		bqAlias.ShouldQueriesAlias = mqs
	}

	data, err := json.Marshal(bqAlias)
	return data, err
}

func (q *BoolQuery) UnmarshalJSON(data []byte) (err error) {
	type BoolQueryAlias BoolQuery
	bqAlias := &BoolQueryAlias{}
	err = json.Unmarshal(data, bqAlias)
	if err != nil {
		return
	}

	q.MinimumShouldMatch = bqAlias.MinimumShouldMatch
	if bqAlias.MustQueriesAlias != nil {
		mqs := make([]Query, 0)
		for _, tq := range bqAlias.MustQueriesAlias {
			mqs = append(mqs, tq.Query)
		}
		q.MustQueries = mqs
	}

	if bqAlias.MustNotQueriesAlias != nil {
		mqs := make([]Query, 0)
		for _, tq := range bqAlias.MustNotQueriesAlias {
			mqs = append(mqs, tq.Query)
		}
		q.MustNotQueries = mqs
	}

	if bqAlias.FilterQueriesAlias != nil {
		mqs := make([]Query, 0)
		for _, tq := range bqAlias.FilterQueriesAlias {
			mqs = append(mqs, tq.Query)
		}
		q.FilterQueries = mqs
	}
	if bqAlias.ShouldQueriesAlias != nil {
		mqs := make([]Query, 0)
		for _, tq := range bqAlias.ShouldQueriesAlias {
			mqs = append(mqs, tq.Query)
		}
		q.ShouldQueries = mqs
	}
	return
}

func (q *BoolQuery) Type() QueryType {
	return QueryType_BoolQuery
}

func (q *BoolQuery) Serialize() ([]byte, error) {
	query := &otsprotocol.BoolQuery{}
	if q.MustQueries != nil {
		pbMustQs := make([]*otsprotocol.Query, 0)
		for _, mustQ := range q.MustQueries {
			pbQ, err := mustQ.ProtoBuffer()
			if err != nil {
				return nil, err
			}
			pbMustQs = append(pbMustQs, pbQ)
		}
		query.MustQueries = pbMustQs
	}
	if q.MustNotQueries != nil {
		pbMustNotQs := make([]*otsprotocol.Query, 0)
		for _, mustNotQ := range q.MustNotQueries {
			pbQ, err := mustNotQ.ProtoBuffer()
			if err != nil {
				return nil, err
			}
			pbMustNotQs = append(pbMustNotQs, pbQ)
		}
		query.MustNotQueries = pbMustNotQs
	}
	if q.FilterQueries != nil {
		pbFilterQs := make([]*otsprotocol.Query, 0)
		for _, filterQ := range q.FilterQueries {
			pbQ, err := filterQ.ProtoBuffer()
			if err != nil {
				return nil, err
			}
			pbFilterQs = append(pbFilterQs, pbQ)
		}
		query.FilterQueries = pbFilterQs
	}
	if q.ShouldQueries != nil {
		pbShouldQs := make([]*otsprotocol.Query, 0)
		for _, shouldQ := range q.ShouldQueries {
			pbQ, err := shouldQ.ProtoBuffer()
			if err != nil {
				return nil, err
			}
			pbShouldQs = append(pbShouldQs, pbQ)
		}
		query.ShouldQueries = pbShouldQs
	}
	if (q.MinimumShouldMatch != nil) {
		query.MinimumShouldMatch = q.MinimumShouldMatch
	}
	data, err := proto.Marshal(query)
	return data, err
}

func (q *BoolQuery) ProtoBuffer() (*otsprotocol.Query, error) {
	return BuildPBForQuery(q)
}

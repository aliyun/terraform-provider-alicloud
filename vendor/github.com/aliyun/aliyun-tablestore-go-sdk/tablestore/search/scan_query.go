package search

import (
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/golang/protobuf/proto"
)

type ScanQuery interface {
	Serialize() ([]byte, error)
}

type scanQuery struct {
	Query		Query
	Limit		*int32
	AliveTime	*int32	//in seconds, 60s by default
	Token		[]byte
	CurrentParallelID	*int32
	MaxParallel			*int32
}

func NewScanQuery() *scanQuery {
	return &scanQuery{}
}

func (s *scanQuery) SetQuery(query Query)  *scanQuery {
	s.Query = query
	return s
}

func (s *scanQuery) SetLimit(limit int32) *scanQuery {
	s.Limit = proto.Int32(limit)
	return s
}

func (s *scanQuery) SetAliveTime(aliveTime int32) *scanQuery {
	s.AliveTime = proto.Int32(aliveTime)
	return s
}

func (s *scanQuery) SetToken(token []byte) *scanQuery {
	s.Token = token
	return s
}

func (s *scanQuery) SetCurrentParallelID(currentParallelID int32) *scanQuery {
	s.CurrentParallelID = proto.Int32(currentParallelID)
	return s
}

func (s *scanQuery) SetMaxParallel(maxParallel int32) *scanQuery {
	s.MaxParallel = proto.Int32(maxParallel)
	return s
}

func (s *scanQuery) Serialize() ([]byte, error) {
	scanQuery := &otsprotocol.ScanQuery{}

	if s.Query != nil {
		pbQuery, err := s.Query.ProtoBuffer()
		if err != nil {
			return nil, err
		}
		scanQuery.Query = pbQuery
	}
	if s.Limit != nil && *s.Limit >= 0 {
		scanQuery.Limit = s.Limit
	}
	if s.AliveTime != nil && *s.AliveTime > 0 {
		scanQuery.AliveTime = s.AliveTime
	}
	if s.Token != nil && len(s.Token) > 0 {
		scanQuery.Token = s.Token
	}
	if s.CurrentParallelID != nil && *s.CurrentParallelID >= 0 {
		scanQuery.CurrentParallelId = s.CurrentParallelID
	}
	if s.MaxParallel != nil && *s.MaxParallel > 0 {
		scanQuery.MaxParallel = s.MaxParallel
	}

	data, err := proto.Marshal(scanQuery)
	return data, err
}

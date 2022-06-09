package search

import (
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
)

type SubAggGroupBySort struct {
	Order		*SortOrder
	SubAggName	string
}

func (s *SubAggGroupBySort) ProtoBuffer() (*otsprotocol.GroupBySorter, error) {
	pbSort := &otsprotocol.SubAggSort{}
	if s.Order != nil {
		pbOrder, err := s.Order.ProtoBuffer()
		if err != nil {
			return nil, err
		}
		pbSort.Order = pbOrder
	}
	pbSort.SubAggName = &s.SubAggName
	pbSorter := &otsprotocol.GroupBySorter{
		SubAggSort: pbSort,
	}
	return pbSorter, nil
}

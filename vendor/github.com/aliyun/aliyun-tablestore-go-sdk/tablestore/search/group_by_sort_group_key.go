package search

import "github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"

type GroupKeyGroupBySort struct {
	Order		*SortOrder
}

func (s *GroupKeyGroupBySort) ProtoBuffer() (*otsprotocol.GroupBySorter, error) {
	pbSort := &otsprotocol.GroupKeySort{}
	if s.Order != nil {
		pbOrder, err := s.Order.ProtoBuffer()
		if err != nil {
			return nil, err
		}
		pbSort.Order = pbOrder
	}
	pbSorter := &otsprotocol.GroupBySorter{
		GroupKeySort: pbSort,
	}
	return pbSorter, nil
}

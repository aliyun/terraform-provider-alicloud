package search

import "github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"

type RowCountGroupBySort struct {
	Order		*SortOrder
}

func (s *RowCountGroupBySort) ProtoBuffer() (*otsprotocol.GroupBySorter, error) {
	pbSort := &otsprotocol.RowCountSort{}
	if s.Order != nil {
		pbOrder, err := s.Order.ProtoBuffer()
		if err != nil {
			return nil, err
		}
		pbSort.Order = pbOrder
	}
	pbSorter := &otsprotocol.GroupBySorter{
		RowCountSort: pbSort,
	}
	return pbSorter, nil
}

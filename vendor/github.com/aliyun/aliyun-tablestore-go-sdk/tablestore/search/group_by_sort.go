package search

import (
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
)

//interface of each child 'group by sorter'
type GroupBySorter interface {
	//should be implemented by each 'group by sorter'
	ProtoBuffer() (*otsprotocol.GroupBySorter, error)
}

type GroupBySort struct {
	Sorters []GroupBySorter
}

func (s *GroupBySort) ProtoBuffer() (*otsprotocol.GroupBySort, error) {
	pbGroupBySort := &otsprotocol.GroupBySort{}
	pbGroupBySorters := make([]*otsprotocol.GroupBySorter, 0)
	for _, fs := range s.Sorters {
		pbFs, err := fs.ProtoBuffer()
		if err != nil {
			return nil, err
		}
		pbGroupBySorters = append(pbGroupBySorters, pbFs)
	}
	pbGroupBySort.Sorters = pbGroupBySorters
	return pbGroupBySort, nil
}

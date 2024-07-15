package search

import "github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"

type DocSort struct {
	SortOrder *SortOrder
}

func (docSort *DocSort) ProtoBuffer() (sorter *otsprotocol.Sorter, err error) {
	pbDocSort := &otsprotocol.DocSort{}
	if docSort.SortOrder != nil {
		if pbDocSort.Order, err = docSort.SortOrder.ProtoBuffer(); err != nil {
			return
		}
	} else {
		pbDocSort.Order = otsprotocol.SortOrder_SORT_ORDER_ASC.Enum()
	}
	
	sorter = &otsprotocol.Sorter{
		DocSort: pbDocSort,
	}
	return
}


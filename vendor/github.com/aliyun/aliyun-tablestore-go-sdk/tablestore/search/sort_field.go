package search

import (
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
)

const (
	FirstWhenMissing = "_first"
	LastWhenMissing  = "_last"
)

type NestedFilter struct {
	Path   string
	Filter Query
}

func (f *NestedFilter) ProtoBuffer() (*otsprotocol.NestedFilter, error) {
	pbF := &otsprotocol.NestedFilter{
		Path: &f.Path,
	}
	pbQ, err := f.Filter.ProtoBuffer()
	if err != nil {
		return nil, err
	}
	pbF.Filter = pbQ
	return pbF, err
}

type FieldSort struct {
	FieldName    string
	Order        *SortOrder
	Mode         *SortMode
	NestedFilter *NestedFilter
	MissingValue interface{} // 当排序的字段某些行没有填充值时，排序行为支持三种方式：1、设置为FirstWhenMissing，当排序字段值缺省时候排在最前面；2、设置为LastWhenMissing，当排序字段值缺省时候排在最后面；3、自定义值，当排序字段值缺省时候使用指定的值进行排序。
	MissingField *string
}

func NewFieldSort(fieldName string, order SortOrder) *FieldSort {
	return &FieldSort{
		FieldName: fieldName,
		Order:     order.Enum(),
	}
}

func (s *FieldSort) ProtoBuffer() (*otsprotocol.Sorter, error) {
	pbFieldSort := &otsprotocol.FieldSort{
		FieldName: &s.FieldName,
	}
	if s.Order != nil {
		pbOrder, err := s.Order.ProtoBuffer()
		if err != nil {
			return nil, err
		}
		pbFieldSort.Order = pbOrder
	}
	if s.Mode != nil {
		pbMode, err := s.Mode.ProtoBuffer()
		if err != nil {
			return nil, err
		}
		if pbMode != nil {
			pbFieldSort.Mode = pbMode
		}
	}
	if s.NestedFilter != nil {
		pbFilter, err := s.NestedFilter.ProtoBuffer()
		if err != nil {
			return nil, err
		}
		pbFieldSort.NestedFilter = pbFilter
	}
	if s.MissingField != nil {
		pbFieldSort.MissingField = s.MissingField
	}
	//missingValue
	if s.MissingValue != nil {
		vt, err := ToVariantValue(s.MissingValue)
		if err != nil {
			return nil, err
		}
		pbFieldSort.MissingValue = []byte(vt)
	}
	pbSorter := &otsprotocol.Sorter{
		FieldSort: pbFieldSort,
	}
	return pbSorter, nil
}

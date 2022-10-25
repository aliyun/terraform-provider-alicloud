package search

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
)

type SortOrder int8

const (
	SortOrder_ASC  SortOrder = 0
	SortOrder_DESC SortOrder = 1
)

func (x SortOrder) Enum() *SortOrder {
	p := new(SortOrder)
	*p = x
	return p
}

func (o *SortOrder) ProtoBuffer() (*otsprotocol.SortOrder, error) {
	if o == nil {
		return nil, errors.New("sort order is nil")
	}
	if *o == SortOrder_ASC {
		return otsprotocol.SortOrder_SORT_ORDER_ASC.Enum(), nil
	} else if *o == SortOrder_DESC {
		return otsprotocol.SortOrder_SORT_ORDER_DESC.Enum(), nil
	} else {
		return nil, errors.New("unknown sort order: " + fmt.Sprintf("%#v", *o))
	}
}

func ParseSortOrder(order *otsprotocol.SortOrder) *SortOrder {
	if order == nil {
		return nil
	}
	if *order == otsprotocol.SortOrder_SORT_ORDER_ASC {
		return SortOrder_ASC.Enum()
	} else if *order == otsprotocol.SortOrder_SORT_ORDER_DESC {
		return SortOrder_DESC.Enum()
	} else {
		return nil
	}
}

func (order SortOrder) String() string {
	switch order {
	case SortOrder_ASC:
		return "ASC"
	case SortOrder_DESC:
		return "DESC"
	default:
		return fmt.Sprintf("%d", order)
	}
}

func ToSortOrder(order string) (SortOrder, error) {
	switch strings.ToUpper(order) {
	case "ASC":
		return SortOrder_ASC, nil
	case "DESC":
		return SortOrder_DESC, nil
	default:
		return SortOrder_ASC, errors.New("Invalid sort order: " + order)
	}
}

func (op *SortOrder) UnmarshalJSON(data []byte) (err error) {
	var opStr string
	err = json.Unmarshal(data, &opStr)
	if err != nil {
		return
	}

	*op, err = ToSortOrder(opStr)
	if err != nil {
		return err
	}
	return
}

func (op *SortOrder) MarshalJSON() (data []byte, err error) {
	data, err = json.Marshal(op.String())
	return
}

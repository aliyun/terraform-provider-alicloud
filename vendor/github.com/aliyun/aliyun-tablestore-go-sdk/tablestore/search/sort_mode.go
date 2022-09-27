package search

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"strings"
)

type SortMode int8

const (
	SortMode_Min SortMode = 0
	SortMode_Max SortMode = 1
	SortMode_Avg SortMode = 2
)

func (x SortMode) Enum() *SortMode {
	p := new(SortMode)
	*p = x
	return p
}

func (m *SortMode) ProtoBuffer() (*otsprotocol.SortMode, error) {
	if m == nil {
		return nil, errors.New("sort mode is nil")
	}
	if *m == SortMode_Min {
		return otsprotocol.SortMode_SORT_MODE_MIN.Enum(), nil
	} else if *m == SortMode_Max {
		return otsprotocol.SortMode_SORT_MODE_MAX.Enum(), nil
	} else if *m == SortMode_Avg {
		return otsprotocol.SortMode_SORT_MODE_AVG.Enum(), nil
	} else {
		return nil, errors.New("unknown sort mode: " + fmt.Sprintf("%#v", *m))
	}
}

func (mode SortMode) String() string {
	switch mode {
	case SortMode_Min:
		return "MIN"
	case SortMode_Max:
		return "MAX"
	case SortMode_Avg:
		return "AVG"
	default:
		return fmt.Sprintf("%d", mode)
	}
}

func ToSortMode(mode string) (SortMode, error) {
	switch strings.ToUpper(mode) {
	case "MIN":
		return SortMode_Min, nil
	case "MAX":
		return SortMode_Max, nil
	case "AVG":
		return SortMode_Avg, nil
	default:
		return SortMode_Min, errors.New("Invalid sort mode: " + mode)
	}
}

func (op *SortMode) UnmarshalJSON(data []byte) (err error) {
	var opStr string
	err = json.Unmarshal(data, &opStr)
	if err != nil {
		return
	}

	*op, err = ToSortMode(opStr)
	if err != nil {
		return err
	}
	return
}

func (op *SortMode) MarshalJSON() (data []byte, err error) {
	data, err = json.Marshal(op.String())
	return
}
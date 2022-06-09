package search

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"strings"
)

type GeoDistanceType int8

const (
	GeoDistanceType_ARC   GeoDistanceType = 0
	GeoDistanceType_PLANE GeoDistanceType = 1
)

func (t *GeoDistanceType) ProtoBuffer() (*otsprotocol.GeoDistanceType, error) {
	if t == nil {
		return nil, errors.New("type is nil")
	}
	if *t == GeoDistanceType_ARC {
		return otsprotocol.GeoDistanceType_GEO_DISTANCE_ARC.Enum(), nil
	} else if *t == GeoDistanceType_PLANE {
		return otsprotocol.GeoDistanceType_GEO_DISTANCE_PLANE.Enum(), nil
	} else {
		return nil, errors.New("unknown distance type: " + fmt.Sprintf("%#v", *t))
	}
}

type GeoDistanceSort struct {
	FieldName       string
	Points          []string
	Order           *SortOrder
	Mode            *SortMode
	GeoDistanceType *GeoDistanceType
	NestedFilter    *NestedFilter
}

func (s *GeoDistanceSort) ProtoBuffer() (*otsprotocol.Sorter, error) {
	pbGeoDistanceSort := &otsprotocol.GeoDistanceSort{
		FieldName: &s.FieldName,
		Points:    s.Points,
	}
	if s.Order != nil {
		pbOrder, err := s.Order.ProtoBuffer()
		if err != nil {
			return nil, err
		}
		pbGeoDistanceSort.Order = pbOrder
	}
	if s.Mode != nil {
		pbMode, err := s.Mode.ProtoBuffer()
		if err != nil {
			return nil, err
		}
		if pbMode != nil {
			pbGeoDistanceSort.Mode = pbMode
		}
	}
	if s.GeoDistanceType != nil {
		pbGeoDisType, err := s.GeoDistanceType.ProtoBuffer()
		if err != nil {
			return nil, err
		}
		pbGeoDistanceSort.DistanceType = pbGeoDisType
	}
	if s.NestedFilter != nil {
		pbFilter, err := s.NestedFilter.ProtoBuffer()
		if err != nil {
			return nil, err
		}
		pbGeoDistanceSort.NestedFilter = pbFilter
	}
	pbSorter := &otsprotocol.Sorter{
		GeoDistanceSort: pbGeoDistanceSort,
	}
	return pbSorter, nil
}

func (t GeoDistanceType) String() string {
	switch t {
	case GeoDistanceType_ARC:
		return "ARC"
	case GeoDistanceType_PLANE:
		return "PLANE"
	default:
		return fmt.Sprintf("%d", t)
	}
}

func ToGeoDistanceType(t string) (GeoDistanceType, error) {
	switch strings.ToUpper(t) {
	case "ARC":
		return GeoDistanceType_ARC, nil
	case "PLANE":
		return GeoDistanceType_PLANE, nil
	default:
		return GeoDistanceType_ARC, errors.New("Invalid type: " + t)
	}
}

func (t *GeoDistanceType) UnmarshalJSON(data []byte) (err error) {
	var opStr string
	err = json.Unmarshal(data, &opStr)
	if err != nil {
		return
	}

	*t, err = ToGeoDistanceType(opStr)
	if err != nil {
		return err
	}
	return
}

func (op *GeoDistanceType) MarshalJSON() (data []byte, err error) {
	data, err = json.Marshal(op.String())
	return
}
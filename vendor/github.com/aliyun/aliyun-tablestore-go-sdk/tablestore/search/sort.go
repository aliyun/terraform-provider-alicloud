package search

import (
	"encoding/json"
	"errors"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
)

type Sorter interface {
	ProtoBuffer() (*otsprotocol.Sorter, error)
}

type Sort struct {
	Sorters []Sorter
}

func (s *Sort) ProtoBuffer() (*otsprotocol.Sort, error) {
	pbSort := &otsprotocol.Sort{}
	pbSortors := make([]*otsprotocol.Sorter, 0)
	for _, fs := range s.Sorters {
		pbFs, err := fs.ProtoBuffer()
		if err != nil {
			return nil, err
		}
		pbSortors = append(pbSortors, pbFs)
	}
	pbSort.Sorter = pbSortors
	return pbSort, nil
}

func (s *Sort) MarshalJSON() ([]byte, error) {
	type SorterInJson struct {
		Name string
		Sorter Sorter
	}

	sorters := make(map[string]interface{})
	data := make([]SorterInJson, 0)
	for _, sorter := range s.Sorters {
		sorterName := ""
		switch sorter.(type) {
		case *PrimaryKeySort:
			sorterName = "PrimaryKeySort"
		case *GeoDistanceSort:
			sorterName = "GeoDistanceSort"
		case *ScoreSort:
			sorterName = "ScoreSort"
		case *FieldSort:
			sorterName = "FieldSort"
		default:
			return nil, errors.New("Unknown sort type.")
		}

		data = append(data,
			SorterInJson{
				Name:   sorterName,
				Sorter: sorter,
			})
	}

	sorters["Sorters"] = data
	return json.Marshal(sorters)
}

func (r *Sort) UnmarshalJSON(data []byte) (err error) {
	rawData := make(map[string]json.RawMessage)
	err = json.Unmarshal(data, &rawData)
	if err != nil {
		return
	}

	sortersRawMessage, ok := rawData["Sorters"]
	if !ok {
		return
	}

	sorters := make([]map[string]json.RawMessage, 0)
	err = json.Unmarshal(sortersRawMessage, &sorters)
	if err != nil {
		return
	}

	r.Sorters = make([]Sorter, 0)
	for _, sorter := range sorters {
		sorterNameRawMessage, hasName := sorter["Name"]
		sorterRawMessage, hasData := sorter["Sorter"]
		if !hasName || !hasData {
			err = errors.New("Sorter is invalid.")
			return
		}

		var sorterName = ""
		err = json.Unmarshal(sorterNameRawMessage, &sorterName)
		if err != nil {
			return
		}
		switch string(sorterName) {
		case "PrimaryKeySort":
			s := &PrimaryKeySort{}
			err = json.Unmarshal(sorterRawMessage, s)
			r.Sorters = append(r.Sorters, s)
		case "GeoDistanceSort":
			s := &GeoDistanceSort{}
			err = json.Unmarshal(sorterRawMessage, s)
			r.Sorters = append(r.Sorters, s)
		case "ScoreSort":
			s := &ScoreSort{}
			err = json.Unmarshal(sorterRawMessage, s)
			r.Sorters = append(r.Sorters, s)
		case "FieldSort":
			s := &FieldSort{}
			err = json.Unmarshal(sorterRawMessage, s)
			r.Sorters = append(r.Sorters, s)
		default:
			err = errors.New("Unknown sorter type: " + string(sorterName))
			return
		}
	}

	return
}

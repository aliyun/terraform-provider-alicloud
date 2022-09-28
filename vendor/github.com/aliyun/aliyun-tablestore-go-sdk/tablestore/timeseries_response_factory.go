package tablestore

import (
	"bytes"
	"fmt"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
)

func CreateGetTimeseriesDataResponse(pbResponse *otsprotocol.GetTimeseriesDataResponse) (*GetTimeseriesDataResponse , error) {
	response := new(GetTimeseriesDataResponse)
	if pbResponse.GetRowsData() != nil && len(pbResponse.GetRowsData()) > 0 {
		rows, err := readRowsWithHeader(bytes.NewReader(pbResponse.RowsData))
		if err != nil {
			return nil, fmt.Errorf("parser response failed with err : %v" , err)
		}

		for _ , row := range rows {
			timeseriesKey := NewTimeseriesKey()
			measurement := row.primaryKey[1].cellValue.Value.(string)
			source := row.primaryKey[2].cellValue.Value.(string)
			tagsStr := row.primaryKey[3].cellValue.Value.(string)
			timestamp := row.primaryKey[4].cellValue.Value.(int64)
			tags, err := parseTagsOrAttrs(tagsStr)
			if err != nil {
				return nil, err
			}
			timeseriesKey.SetMeasurementName(measurement)
			timeseriesKey.SetDataSource(source)
			timeseriesKey.AddTags(tags)

			timeseriesRow := NewTimeseriesRow(timeseriesKey)
			timeseriesRow.SetTimeInus(timestamp)
			for _ , field := range row.cells {
				timeseriesRow.AddField(convertColumnName(field.cellName) , field.cellValue)
			}
			response.rows = append(response.rows , timeseriesRow)
		}
	}

	if pbResponse.GetNextToken() != nil && len(pbResponse.GetNextToken()) != 0{
		response.nextToken = pbResponse.NextToken
	}
	return response , nil

}

func parseTimeseriesMeta(pbResponseMeta *otsprotocol.TimeseriesMeta) (*TimeseriesMeta, error) {
	timeseriesKey := NewTimeseriesKey()

	currentMetaTagStr := pbResponseMeta.GetTimeSeriesKey().GetTags()
	tags, err := parseTagsOrAttrs(currentMetaTagStr)
	if err != nil {
		return nil, err
	}
	timeseriesKey.AddTags(tags)

	timeseriesKey.source = pbResponseMeta.GetTimeSeriesKey().GetSource()
	timeseriesKey.measurement = pbResponseMeta.GetTimeSeriesKey().GetMeasurement()

	timeseriesMeta := NewTimeseriesMeta(timeseriesKey)

	if pbResponseMeta.Attributes != nil {
		attrs, err := parseTagsOrAttrs(*pbResponseMeta.Attributes)
		if err != nil {
			return nil, err
		}
		timeseriesMeta.attributes = attrs
	}
	if pbResponseMeta.UpdateTime != nil {
		timeseriesMeta.updateTimeInUs = *pbResponseMeta.UpdateTime
	}
	return timeseriesMeta, nil
}


func ParseTimeseriesTableMeta(pbResponseTableMeta *otsprotocol.TimeseriesTableMeta) (*TimeseriesTableMeta) {
	timeseriesTableMeta := NewTimeseriesTableMeta(pbResponseTableMeta.GetTableName())
	timeseriesTableOptions := NewTimeseriesTableOptions(int64(pbResponseTableMeta.GetTableOptions().GetTimeToLive()))
	timeseriesTableMeta.SetTimeseriesTableOptions(timeseriesTableOptions)
	return timeseriesTableMeta
}

func convertColumnName(serverColName []byte) string {
	for i := len(serverColName) - 1; i >= 0; i-- {
		if serverColName[i] == ':' {
			return string(serverColName[:i])
		}
	}
	return ""
}

func parseTagsOrAttrs(tagsStr string) (map[string]string, error) {
	if tagsStr == "" {
		return nil, fmt.Errorf("tags string is empty")
	}

	if len(tagsStr) < 2 || tagsStr[0] != '[' || tagsStr[len(tagsStr) - 1] != ']' {
		return nil, fmt.Errorf("invalid tags or attributes string: %v" , tagsStr)
	}

	tags := map[string]string{}
	keyStart := -1
	valueStart := -1
	for i := 1; i < len(tagsStr) - 1; i++ {
		if tagsStr[i] != '"' {
			return nil, fmt.Errorf("invalid tags or attributes string: %v" , tagsStr)
		}
		i += 1
		keyStart = i
		for ;(i < len(tagsStr) - 1) && (tagsStr[i] != '=') && (tagsStr[i] != '"'); {
			i++
		}
		if tagsStr[i] != '=' {
			return nil, fmt.Errorf("invalid tags or attributes string: %v" , tagsStr)
		}
		i += 1
		valueStart = i
		for ;(i < len(tagsStr) - 1) && (tagsStr[i] != '"'); {
			i++
		}
		if tagsStr[i] != '"' {
			return nil, fmt.Errorf("invalid tags or attributes string: %v" , tagsStr)
		}
		tags[tagsStr[keyStart:valueStart-1]] = tagsStr[valueStart:i]
		i += 1
		if i < len(tagsStr) - 1 && tagsStr[i] != ',' {
			return nil, fmt.Errorf("invalid tags or attributes string: %v" , tagsStr)
		}
	}
	return tags, nil
}
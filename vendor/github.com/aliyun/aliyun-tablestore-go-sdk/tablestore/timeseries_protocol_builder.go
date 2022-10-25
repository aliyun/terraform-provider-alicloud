package tablestore

import (
	"fmt"
	lruCache "github.com/hashicorp/golang-lru"
	"reflect"

	Fieldvalues "github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/timeseries/flatbuffer"
	flatbuffers "github.com/google/flatbuffers/go"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/golang/protobuf/proto"
)


func BuildFlatbufferRows(rows []*TimeseriesRow , timeseriesTableName string , timeseriesMetaCache *lruCache.Cache) ([]byte , error) {
	rowsNum := len(rows)
	rowGroupOffs := make([]flatbuffers.UOffsetT , rowsNum)
	fbb := flatbuffers.NewBuilder(1024)

	var err error
	for i := 0; i < rowsNum; i++ {
		rowGroupOffs[i] , err = buildTimeseriesRowToRowGroupOffset(rows[i] , fbb , timeseriesTableName , timeseriesMetaCache)
		if err != nil {
			return nil , fmt.Errorf("BuildFlatbufferRows failed! ")
		}
	}

	rowsVectorOffset := createRowGroupsVector(fbb , rowGroupOffs)
	rowsOffset := createFlatBufferRows(fbb , rowsVectorOffset)

	fbb.Finish(rowsOffset)
	return fbb.FinishedBytes() , nil
}


func buildTimeseriesRowToRowGroupOffset(row *TimeseriesRow , fbb *flatbuffers.Builder , timeseriesTableName string , timeseriesMetaCache *lruCache.Cache) (flatbuffers.UOffsetT , error) {
	fieldCount := len(row.fields)
	fieldValueTypes := make([]Fieldvalues.DataType, fieldCount)
	fieldNameOffs := make([]flatbuffers.UOffsetT , fieldCount)

	var idx int = 0
	var doubleValueCount int = 0
	var longValueCount int = 0
	var boolValueCount int = 0
	var binaryValueCount int = 0
	var stringValueCount int = 0

	field_keys , field_values := SortedMapColumnValue(row.fields)
	for i := 0; i < fieldCount; i++ {
		fieldNameOffs[i] = fbb.CreateString(field_keys[i])
		switch field_values[i].Type {
		case ColumnType_INTEGER:
			fieldValueTypes[i] = Fieldvalues.DataTypeLONG		// LONG
			longValueCount++
			break
		case ColumnType_BOOLEAN:
			fieldValueTypes[i] = Fieldvalues.DataTypeBOOLEAN		// BOOLEAN
			boolValueCount++
			break
		case ColumnType_DOUBLE:
			fieldValueTypes[i] = Fieldvalues.DataTypeDOUBLE		// DOUBLE
			doubleValueCount++
			break
		case ColumnType_STRING:
			fieldValueTypes[i] = Fieldvalues.DataTypeSTRING		// STRING
			stringValueCount++
			break
		case ColumnType_BINARY:
			fieldValueTypes[i] = Fieldvalues.DataTypeBINARY		// BINARY
			binaryValueCount++
			break
		default:
			return 0 , fmt.Errorf("Err ColumnType : %v" , field_values[i].Type)
		}
	}

	longValues := make([]int64 , longValueCount)
	boolValues := make([]bool , boolValueCount)
	doubleValues := make([]float64 , doubleValueCount)
	strValueOffs := make([]flatbuffers.UOffsetT , stringValueCount)
	binaryValueOffs := make([]flatbuffers.UOffsetT , binaryValueCount)

	doubleValueCount = 0
	longValueCount  = 0
	boolValueCount  = 0
	stringValueCount  = 0
	binaryValueCount  = 0

	for i := 0; i < fieldCount; i++ {
		switch field_values[i].Type {
		case ColumnType_INTEGER:
			switch value := field_values[i].Value.(type) {
			case int64:
				longValues[longValueCount] = value
			case int:
				longValues[longValueCount] = int64(value)
			default:
				return 0, fmt.Errorf("unsupported field type: %v", reflect.TypeOf(field_values[i].Value))
			}
			longValueCount++
			break
		case ColumnType_BOOLEAN:
			boolValues[boolValueCount] = field_values[i].Value.(bool)
			boolValueCount++
			break
		case ColumnType_DOUBLE:
			doubleValues[doubleValueCount] = field_values[i].Value.(float64)
			doubleValueCount++
			break
		case ColumnType_STRING:
			strValueOffs[stringValueCount] = fbb.CreateString(field_values[i].Value.(string))
			stringValueCount++
			break
		case ColumnType_BINARY:
			binaryValueOffs[binaryValueCount] = CreateBytesValue(fbb , CreateBytesValueVector(fbb , field_values[i].Value.([]byte)))
			binaryValueCount++
			break
		default:
			return 0 , fmt.Errorf("Err ColumnType : %v" , field_values[idx].Type)
		}
	}

	long_valuesOffset := flatbuffers.UOffsetT(0)
	bool_valuesOffset := flatbuffers.UOffsetT(0)
	double_valuesOffset := flatbuffers.UOffsetT(0)
	string_valuesOffset := flatbuffers.UOffsetT(0)
	binary_valuesOffset := flatbuffers.UOffsetT(0)
	fieldValueOff := flatbuffers.UOffsetT(0)

	if longValueCount != 0 {
		fbb.StartVector(8 , longValueCount , 8)
		for i := longValueCount - 1; i >= 0; i-- {
			fbb.PrependInt64(longValues[i])
		}
		long_valuesOffset = fbb.EndVector(longValueCount)
	}

	if boolValueCount != 0 {
		fbb.StartVector(1 , boolValueCount , 1 )
		for i := boolValueCount - 1; i >= 0; i-- {
			fbb.PrependBool(boolValues[i])
		}
		bool_valuesOffset = fbb.EndVector(boolValueCount)
	}

	if doubleValueCount != 0 {
		fbb.StartVector(8 , doubleValueCount , 8)
		for i := doubleValueCount - 1; i >= 0; i-- {
			fbb.PrependFloat64(doubleValues[i])
		}
		double_valuesOffset = fbb.EndVector(doubleValueCount)
	}

	if stringValueCount != 0 {
		fbb.StartVector(4 , stringValueCount , 4)
		for i := stringValueCount - 1; i >= 0; i-- {
			fbb.PrependUOffsetT(strValueOffs[i])
		}
		string_valuesOffset = fbb.EndVector(stringValueCount)
	}

	if binaryValueCount != 0 {
		fbb.StartVector(4 , binaryValueCount , 4)
		for i := binaryValueCount - 1; i >= 0; i-- {
			fbb.PrependUOffsetT(flatbuffers.UOffsetT(binaryValueOffs[i]))
		}
		binary_valuesOffset = fbb.EndVector(binaryValueCount)
	}

	fbb.StartObject(5)
	Fieldvalues.FieldValuesAddBinaryValues(fbb , binary_valuesOffset)
	Fieldvalues.FieldValuesAddStringValues(fbb , string_valuesOffset)
	Fieldvalues.FieldValuesAddDoubleValues(fbb , double_valuesOffset)
	Fieldvalues.FieldValuesAddBoolValues(fbb , bool_valuesOffset)
	Fieldvalues.FieldValuesAddLongValues(fbb , long_valuesOffset)
	fieldValueOff = Fieldvalues.FieldValuesEnd(fbb)

	var source_keyOffset flatbuffers.UOffsetT
	if row.timeseriesKey.source == "" {
		source_keyOffset = fbb.CreateString("")
	} else {
		source_keyOffset = fbb.CreateString(row.timeseriesKey.source)
	}

	var tags_Offset flatbuffers.UOffsetT
	var err error

	if row.timeseriesKey.tagsString == nil {
		row.timeseriesKey.tagsString = new(string)
		if *row.timeseriesKey.tagsString , err = BuildTagString(row.timeseriesKey.tags); err != nil {
			return 0 , fmt.Errorf("Build tags string failed with error: %s" , err)
		}
	}
	tags_Offset = fbb.CreateString(*row.timeseriesKey.tagsString)

	rowInGroupOffs := make([]flatbuffers.UOffsetT , 1)
	if row.timeseriesMetaKey == nil {
		row.timeseriesMetaKey = new(string)
		if *row.timeseriesMetaKey , err = row.timeseriesKey.buildTimeseriesMetaKey(timeseriesTableName); err != nil {
			return 0 , fmt.Errorf("Build meta key failed with error: %s" , err)
		}
	}

	updateTimeInSec , ok := timeseriesMetaCache.Get(*row.timeseriesMetaKey)
	var updateTime uint32
	if ok {
		updateTime = updateTimeInSec.(uint32)
	}

	Fieldvalues.FlatBufferRowInGroupStart(fbb)
	Fieldvalues.FlatBufferRowInGroupAddTime(fbb , row.timeInUs)
	Fieldvalues.FlatBufferRowInGroupAddMetaCacheUpdateTime(fbb , updateTime)
	Fieldvalues.FlatBufferRowInGroupAddFieldValues(fbb , fieldValueOff)
	Fieldvalues.FlatBufferRowInGroupAddTags(fbb , tags_Offset)
	Fieldvalues.FlatBufferRowInGroupAddDataSource(fbb , source_keyOffset)
	rowInGroupOffs[0] = Fieldvalues.FlatBufferRowInGroupEnd(fbb)

	measurement_namesOffset := fbb.CreateString(row.timeseriesKey.measurement)
	field_namesOffset := createFieldNamesVector(fbb , fieldNameOffs)
	field_typesOffset := createFieldTypesVector(fbb , fieldValueTypes)
	rowsOffset := createRowsVector(fbb , rowInGroupOffs)

	Fieldvalues.FlatBufferRowGroupStart(fbb)
	Fieldvalues.FlatBufferRowGroupAddRows(fbb , rowsOffset)
	Fieldvalues.FlatBufferRowGroupAddFieldTypes(fbb , field_typesOffset)
	Fieldvalues.FlatBufferRowGroupAddFieldNames(fbb , field_namesOffset)
	Fieldvalues.FlatBufferRowGroupAddMeasurementName(fbb ,measurement_namesOffset )

	return  Fieldvalues.FlatBufferRowGroupEnd(fbb) , nil
}


func CreateBytesValue(builder *flatbuffers.Builder , offset flatbuffers.UOffsetT) flatbuffers.UOffsetT {
	builder.StartObject(1)
	Fieldvalues.BytesValueAddValue(builder , offset)
	return Fieldvalues.BytesValueEnd(builder)
}

func CreateBytesValueVector(builder *flatbuffers.Builder , data []byte) flatbuffers.UOffsetT {
	builder.StartVector(1 , len(data) , 1)
	for i := len(data) - 1; i >= 0; i-- {
		builder.PlaceByte(data[i])
	}
	return builder.EndVector(len(data))
}

func createTagNamesVector(builder *flatbuffers.Builder , tagNameOffs []flatbuffers.UOffsetT) flatbuffers.UOffsetT {
	builder.StartVector(4 , len(tagNameOffs) , 4)
	for i := len(tagNameOffs) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(tagNameOffs[i])
	}
	return builder.EndVector(len(tagNameOffs))
}


func createTagValuesVector(builder *flatbuffers.Builder , data []flatbuffers.UOffsetT) flatbuffers.UOffsetT {
	builder.StartVector(4 , len(data) , 4)
	for i := len(data) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(data[i])
	}
	return builder.EndVector(len(data))
}

func createRowsVector(builder *flatbuffers.Builder , rowInGroupOffs []flatbuffers.UOffsetT) flatbuffers.UOffsetT {
	builder.StartVector(4 , len(rowInGroupOffs) , 4)
	for i := len(rowInGroupOffs) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(rowInGroupOffs[i])
	}
	return builder.EndVector(len(rowInGroupOffs))
}

func createFieldTypesVector(builder *flatbuffers.Builder , data []Fieldvalues.DataType) flatbuffers.UOffsetT {
	builder.StartVector(1 , len(data) , 1)
	for i := len(data) - 1; i >= 0; i-- {
		builder.PrependInt8(int8(data[i]))
	}
	return builder.EndVector(len(data))
}

func createFieldNamesVector(builder *flatbuffers.Builder , fieldNameOffs []flatbuffers.UOffsetT) flatbuffers.UOffsetT {
	builder.StartVector(4 , len(fieldNameOffs), 4)
	for i := len(fieldNameOffs) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(fieldNameOffs[i])
	}
	return builder.EndVector(len(fieldNameOffs))
}

func createRowGroupsVector(builder *flatbuffers.Builder , rowGroupOffs []flatbuffers.UOffsetT) flatbuffers.UOffsetT {
	builder.StartVector(4 , len(rowGroupOffs) , 4)
	for i := len(rowGroupOffs) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(rowGroupOffs[i])
	}

	return builder.EndVector(len(rowGroupOffs))
}

func createFlatBufferRows(builder *flatbuffers.Builder , row_groupsOffset flatbuffers.UOffsetT) flatbuffers.UOffsetT {
	Fieldvalues.FlatBufferRowsStart(builder)
	Fieldvalues.FlatBufferRowsAddRowGroups(builder , row_groupsOffset)
	return Fieldvalues.FlatBufferRowsEnd(builder)
}

func buildTimeseriesKey(curTimeseriesKey *TimeseriesKey) (*otsprotocol.TimeseriesKey , error) {
	var err error
	timeseriesKey := new(otsprotocol.TimeseriesKey)
	if curTimeseriesKey.tagsString == nil {
		curTimeseriesKey.tagsString = new(string)
		if *curTimeseriesKey.tagsString , err = BuildTagString(curTimeseriesKey.tags); err != nil {
			return nil , err
		}
	}
	timeseriesKey.Tags = curTimeseriesKey.tagsString
	timeseriesKey.Source = proto.String(curTimeseriesKey.source)
	timeseriesKey.Measurement = proto.String(curTimeseriesKey.measurement)

	return timeseriesKey , nil
}

func buildProtocolBufferRows(rows []*TimeseriesRow, timeseriesTableName string, timeseriesMetaCache *lruCache.Cache) ([]byte, error) {
	pbRows := new(otsprotocol.TimeseriesPBRows)
	pbRows.Rows = make([]*otsprotocol.TimeseriesRow, len(rows))
	for i, row := range rows {
		// build tag string
		var err error
		if row.timeseriesKey.tagsString == nil {
			row.timeseriesKey.tagsString = new(string)
			if *row.timeseriesKey.tagsString, err = BuildTagString(row.timeseriesKey.tags); err != nil {
				return nil, fmt.Errorf("Build tags string failed with error: %s", err)
			}
		}
		// build meta key
		if row.timeseriesMetaKey == nil {
			row.timeseriesMetaKey = new(string)
			if *row.timeseriesMetaKey, err = row.timeseriesKey.buildTimeseriesMetaKey(timeseriesTableName); err != nil {
				return nil, fmt.Errorf("Build meta key failed with error: %s", err)
			}
		}
		// fetch meta update time
		updateTimeInSec, ok := timeseriesMetaCache.Get(*row.timeseriesMetaKey)
		var updateTime uint32
		if ok {
			updateTime = updateTimeInSec.(uint32)
		}
		// build fields
		fieldMap := row.GetFieldsMap()
		fields := make([]*otsprotocol.TimeseriesField, 0, len(fieldMap))
		for fieldName, fieldValue := range fieldMap {
			field := &otsprotocol.TimeseriesField{}
			field.FieldName = proto.String(fieldName)
			switch fieldValue.Type {
			case ColumnType_STRING:
				field.ValueString = proto.String(fieldValue.Value.(string))
			case ColumnType_INTEGER:
				switch value := fieldValue.Value.(type) {
				case int:
					field.ValueInt = proto.Int64(int64(value))
				case int64:
					field.ValueInt = proto.Int64(value)
				default:
					return nil, fmt.Errorf("unsupported field type: %v", reflect.TypeOf(fieldValue))
				}
			case ColumnType_BOOLEAN:
				field.ValueBool = proto.Bool(fieldValue.Value.(bool))
			case ColumnType_DOUBLE:
				field.ValueDouble = proto.Float64(fieldValue.Value.(float64))
			case ColumnType_BINARY:
				field.ValueBinary = fieldValue.Value.([]byte)
			}
			fields = append(fields, field)
		}
		// build row
		pbRow := &otsprotocol.TimeseriesRow{
			TimeseriesKey: &otsprotocol.TimeseriesKey{
				Measurement: proto.String(row.GetTimeseriesKey().measurement),
				Source:      proto.String(row.GetTimeseriesKey().GetDataSource()),
				Tags:        row.timeseriesKey.tagsString,
			},
			Time:                proto.Int64(row.timeInUs),
			Fields:              fields,
			MetaCacheUpdateTime: proto.Uint32(updateTime),
		}
		pbRows.Rows[i] = pbRow
	}
	return proto.Marshal(pbRows)
}

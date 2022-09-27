package search

import (
	"encoding/binary"
	"errors"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/search/model"
	"math"
	"reflect"
)

type VariantValue []byte
type VariantType byte

const (
	// variant type
	VT_INTEGER VariantType = 0x0
	VT_DOUBLE  VariantType = 0x1
	VT_BOOLEAN VariantType = 0x2
	VT_STRING  VariantType = 0x3
)

func ToVariantValue(value interface{}) (VariantValue, error) {
	if value == nil {
		return nil, errors.New("interface{} should not be nil")
	}
	t := reflect.TypeOf(value)
	switch t.Kind() {
	case reflect.String:
		return VTString(value.(string)), nil
	case reflect.Int:
		return VTInteger(int64(value.(int))), nil
	case reflect.Int64:
		return VTInteger(value.(int64)), nil
	case reflect.Float64:
		return VTDouble(value.(float64)), nil
	case reflect.Bool:
		return VTBoolean(value.(bool)), nil
	default:
		return nil, errors.New("interface{} type must be string/int64/float64.")
	}
}

func ForceConvertToDestColumnValue(bytes []byte) (*model.ColumnValue, error) {
	if len(bytes) == 0 {
		return nil, errors.New("the length of bytes must greater than 0")
	}
	columnValue := new(model.ColumnValue)
	var err error
	if bytes[0] == byte(VT_INTEGER) {
		columnValue.Value, err = AsInteger(bytes)
		if err != nil {
			return nil, err
		}
		columnValue.Type = model.ColumnType_INTEGER
	} else if bytes[0] == byte(VT_DOUBLE) {
		columnValue.Value, err = AsDouble(bytes)
		if err != nil {
			return nil, err
		}
		columnValue.Type = model.ColumnType_DOUBLE
	} else if bytes[0] == byte(VT_STRING) {
		columnValue.Value, err = AsString(bytes)
		if err != nil {
			return nil, err
		}
		columnValue.Type = model.ColumnType_STRING
	} else if bytes[0] == byte(VT_BOOLEAN) {
		columnValue.Value, err = AsBoolean(bytes)
		if err != nil {
			return nil, err
		}
		columnValue.Type = model.ColumnType_BOOLEAN
	} else {
		return columnValue, errors.New("type must be string/int64/float64/boolean")
	}
	return columnValue, nil
}

func (v *VariantValue) GetType() VariantType {
	return VariantType(([]byte)(*v)[0])
}

func VTInteger(v int64) VariantValue {
	buf := make([]byte, 9)
	buf[0] = byte(VT_INTEGER)
	binary.LittleEndian.PutUint64(buf[1:9], uint64(v))
	return (VariantValue)(buf)
}

func AsInteger(bytes []byte) (int64, error) {
	if len(bytes) < 9 {
		return -1, errors.New("the length of bytes must greater than or equal 9")
	}
	return int64(binary.LittleEndian.Uint64(bytes[1:9])), nil
}

func VTDouble(v float64) VariantValue {
	buf := make([]byte, 9)
	buf[0] = byte(VT_DOUBLE)
	binary.LittleEndian.PutUint64(buf[1:9], math.Float64bits(v))
	return (VariantValue)(buf)
}

func AsDouble(bytes []byte) (float64, error) {
	if len(bytes) < 9 {
		return 0.0, errors.New("the length of bytes must greater than or equal 9")
	}
	bits := binary.LittleEndian.Uint64(bytes[1:9])
	return math.Float64frombits(bits), nil
}

func VTString(v string) VariantValue {
	buf := make([]byte, 5+len(v))
	buf[0] = byte(VT_STRING)
	binary.LittleEndian.PutUint32(buf[1:5], uint32(len(v)))
	copy(buf[5:], v)
	return (VariantValue)(buf)
}

func AsString(bytes []byte) (string, error){
	if len(bytes) < 5 {
		return "", errors.New("the length of bytes must greater than or equal 5")
	}
	length :=binary.LittleEndian.Uint32(bytes[1:5])
	return string(bytes[5:5+length]), nil
}

func VTBoolean(b bool) VariantValue {
	buf := make([]byte, 2)
	buf[0] = byte(VT_BOOLEAN)
	if b {
		buf[1] = 1
	} else {
		buf[1] = 0
	}
	return (VariantValue)(buf)
}

func AsBoolean(bytes []byte) (bool, error)  {
	if len(bytes) < 2 {
		return true, errors.New("the length of bytes must greater than or equal 2")
	}
	if bytes[1] == 1 {
		return true, nil
	} else {
		return false, nil
	}
}
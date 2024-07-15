package tablestore

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	MatrixBlockApiVersion      = 0x304d5253
	TagChecksum           byte = 1 // ignore
	TagRow                byte = 2
	TagRowCount           byte = 3
	TagEntirePrimaryKeys  byte = 10
)

func parseMatrixRows(rows []byte) ([]*Row, error) {
	var (
		dataOffset int32
		pkCount    int32
		fieldCount int32
		rowCount   int32
		fieldsName []string
		err        error
	)
	reader := bytes.NewReader(rows)
	dataOffset, pkCount, fieldCount, rowCount, fieldsName, err = ensureMatrixMeta(reader)
	if err != nil {
		return nil, err
	}
	_, err = reader.Seek(int64(dataOffset), io.SeekStart)
	if err != nil {
		return nil, err
	}
	return readMatrixData(reader, rowCount, pkCount, fieldCount-pkCount, fieldsName)
}

func readMatrixData(reader io.Reader, rowCount, pkCount, attrCount int32, fieldNames []string) ([]*Row, error) {
	rows := make([]*Row, rowCount)
	for r := int32(0); r < rowCount; r++ {
		pk := new(PrimaryKey)
		pk.PrimaryKeys = make([]*PrimaryKeyColumn, 0, pkCount)
		row := &Row{
			PrimaryKey: pk,
			Columns:    make([]*AttributeColumn, 0, attrCount),
		}
		if err := checkRowTag(reader); err != nil {
			return nil, err
		}
		for i := int32(0); i < pkCount; i++ {
			v, err := readPkValue(reader)
			if err != nil {
				return nil, err
			}
			pk.AddPrimaryKeyColumn(fieldNames[i], v)
		}
		for i := pkCount; i < pkCount+attrCount; i++ {
			v, err := readAttrValue(reader)
			if err != nil {
				return nil, err
			}
			if v == nil {
				continue
			}
			row.Columns = append(row.Columns, &AttributeColumn{ColumnName: fieldNames[i], Value: v})
		}
		rows[r] = row
	}
	return rows, nil
}

func readAttrValue(reader io.Reader) (interface{}, error) {
	typ, err := readCellType(reader)
	if err != nil {
		return nil, err
	}
	switch typ {
	case VT_INTEGER:
		var v int64
		err := binary.Read(reader, binary.LittleEndian, &v)
		return v, err
	case VT_DOUBLE:
		var v float64
		err := binary.Read(reader, binary.LittleEndian, &v)
		return v, err
	case VT_BOOLEAN:
		var v bool
		err := binary.Read(reader, binary.LittleEndian, &v)
		return v, err
	case VT_STRING:
		return readDataString(reader)
	case VT_BLOB:
		var l int32
		if err := binary.Read(reader, binary.LittleEndian, &l); err != nil {
			return nil, err
		}
		return readBlob(reader, l)
	case VT_NULL:
		return nil, nil // special case: column not exist in this row
	default:
		return nil, fmt.Errorf("unexpected attribute column type %d", typ)
	}
}

func readPkValue(reader io.Reader) (interface{}, error) {
	typ, err := readCellType(reader)
	if err != nil {
		return nil, err
	}
	switch typ {
	case VT_INTEGER:
		var v int64
		err := binary.Read(reader, binary.LittleEndian, &v)
		return v, err
	case VT_STRING:
		return readDataString(reader)
	case VT_BLOB:
		var l int32
		if err := binary.Read(reader, binary.LittleEndian, &l); err != nil {
			return nil, err
		}
		return readBlob(reader, l)
	default:
		return nil, fmt.Errorf("unexpected primary key column type %d", typ)
	}
}

func readDataString(reader io.Reader) (string, error) {
	var len int32
	if err := binary.Read(reader, binary.LittleEndian, &len); err != nil {
		return "", err
	}
	return readString(reader, len)
}

func readCellType(reader io.Reader) (int, error) {
	var typ byte
	if err := binary.Read(reader, binary.LittleEndian, &typ); err != nil {
		return 0, err
	}
	return int(typ), nil
}

func checkRowTag(reader io.Reader) error {
	var tag byte
	err := binary.Read(reader, binary.LittleEndian, &tag)
	if err != nil {
		return err
	}
	if tag != TagRow {
		return fmt.Errorf("want TagRow, got tag %d", tag)
	}
	return nil
}

// MatrixBlock data layout:
// APIVersion(int32):dataOffset(int32):optionOffset(int32):pkCount(int32):attrCount(int32):[fieldNameArray...]:[options...]:[data...]:CRCTag(byte):CRC(byte)
// fieldNameArray layout:
// Length(int16):StringChars(length size)... from pk columns to attribute columns
// options layout:
// Tag(byte):Value(length is according to tag type), eg: TagRowCount's value length is 4 bytes, TagEntirePrimaryKeys value length is 1 byte.
// data layout:
// Type(byte):~Length when type is string or binary(int32)~:Data(1 boolean/8 integer or double/length size string or binary)
func ensureMatrixMeta(reader io.ReadSeeker) (dataOffset, pkCount, fieldCount, rowCount int32, fieldNames []string, err error) {
	var (
		apiVersion   int32
		optionOffset int32
		attrCount    int32
	)
	err = binary.Read(reader, binary.LittleEndian, &apiVersion)
	if err != nil {
		return
	}
	if apiVersion != MatrixBlockApiVersion {
		err = fmt.Errorf("unknow api version %d", apiVersion)
		return
	}
	if err = binary.Read(reader, binary.LittleEndian, &dataOffset); err != nil {
		return
	}
	if err = binary.Read(reader, binary.LittleEndian, &optionOffset); err != nil {
		return
	}
	if err = binary.Read(reader, binary.LittleEndian, &pkCount); err != nil {
		return
	}
	if err = binary.Read(reader, binary.LittleEndian, &attrCount); err != nil {
		return
	}
	fieldCount = pkCount + attrCount
	fieldNames, err = readMatrixFieldsName(reader, fieldCount)
	if err != nil {
		return
	}
	_, err = reader.Seek(int64(optionOffset), io.SeekStart)
	if err != nil {
		return
	}
	rowCount, err = readRowCountFromOptions(reader)
	return
}

func readRowCountFromOptions(reader io.ReadSeeker) (int32, error) {
	var tag byte
	for {
		err := binary.Read(reader, binary.LittleEndian, &tag)
		switch tag {
		case TagRowCount:
			var rowCount int32
			err = binary.Read(reader, binary.LittleEndian, &rowCount)
			return rowCount, err
		case TagEntirePrimaryKeys:
			//skip TagEntirePrimaryKeys
			err = binary.Read(reader, binary.LittleEndian, new(byte))
			if err != nil {
				return 0, err
			}
		default:
			return 0, fmt.Errorf("unknow option tag %d", tag)
		}
	}
}

func readMatrixFieldsName(reader io.ReadSeeker, fieldCount int32) ([]string, error) {
	filedNames := make([]string, fieldCount)
	for i := int32(0); i < fieldCount; i++ {
		var nameLen int16
		err := binary.Read(reader, binary.LittleEndian, &nameLen)
		if err != nil {
			return nil, err
		}
		name, err := readString(reader, int32(nameLen))
		if err != nil {
			return nil, err
		}
		filedNames[i] = name
	}
	return filedNames, nil
}

func readBlob(reader io.Reader, len int32) ([]byte, error) {
	buf := make([]byte, len)
	_, err := io.ReadFull(reader, buf)
	return buf, err
}

func readString(reader io.Reader, len int32) (string, error) {
	buf, err := readBlob(reader, len)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

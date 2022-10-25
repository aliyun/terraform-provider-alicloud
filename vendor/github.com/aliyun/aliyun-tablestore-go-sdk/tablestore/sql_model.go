package tablestore

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/sql/dataprotocol"
)

type SQLStatementType int32

const (
	SQL_SELECT         SQLStatementType = 1
	SQL_CREATE_TABLE   SQLStatementType = 2
	SQL_SHOW_TABLE     SQLStatementType = 3
	SQL_DESCRIBE_TABLE SQLStatementType = 4
	SQL_DROP_TABLE     SQLStatementType = 5
	SQL_ALTER_TABLE    SQLStatementType = 6
)

type SQLPayloadVersion int32

const (
	SQLPAYLOAD_PLAIN_BUFFER SQLPayloadVersion = 1
	SQLPAYLOAD_FLAT_BUFFERS SQLPayloadVersion = 2
)

// SQLResultSet is the result of a sql query. Its cursor starts before the first row
// of the result set. Use Next to fetch next row.
type SQLResultSet interface {
	// Columns returns the column infos of SQLResultSet.
	Columns() []*SQLColumnInfo
	// Next returns the next row of SQLResultSet.
	Next() SQLRow
	// HasNext returns whether finished.
	HasNext() bool
	// Reset reset the cursor to beginning.
	Reset()
}

type TablestoreSQLResultSet struct {
	meta      *SQLTableMeta
	rows      []SQLRow
	rowCursor int
}

func newSQLResultSetFromPlainBuffer(rowBytes []byte) (SQLResultSet, error) {
	rows, err := readRowsWithHeader(bytes.NewReader(rowBytes))
	if err != nil {
		return nil, err
	}
	rs := new(TablestoreSQLResultSet)
	rs.rows = make([]SQLRow, len(rows))
	if len(rows) > 0 {
		// TODO: plain buffer cells may not full
		rs.meta = formatSQLTableMetaFromPlainBufferRow(rows[0])
		for i, row := range rows {
			rs.rows[i] = newPlainBufferSQLRow(row, rs.meta.colOffset)
		}
	}
	return rs, nil
}

func newSQLResultSetFromFlatBuffers(rowBytes []byte) (SQLResultSet, error) {
	if len(rowBytes) == 0 {
		return nil, nil
	}
	sqlResponse := dataprotocol.GetRootAsSQLResponseColumns(rowBytes, 0)

	rs := new(TablestoreSQLResultSet)
	rs.meta = formatSQLTableMetaFromFlatBufferColumns(sqlResponse)
	rs.rows = make([]SQLRow, sqlResponse.RowCount())
	// fetch and cache the column values to accelerate fetch
	allColValues := make([]*dataprotocol.ColumnValues, sqlResponse.ColumnsLength())
	allRleValues := make([]*dataprotocol.RLEStringValues, sqlResponse.ColumnsLength())
	for j := 0; j < sqlResponse.ColumnsLength(); j++ {
		colInfo := rs.meta.columns[j]
		col := new(dataprotocol.SQLResponseColumn)
		sqlResponse.Columns(col, j)
		colValues := new(dataprotocol.ColumnValues)
		col.ColumnValue(colValues)
		allColValues[j] = colValues
		if colInfo.dataType == dataprotocol.DataTypeSTRING_RLE {
			rleStringValues := new(dataprotocol.RLEStringValues)
			colValues.RleStringValues(rleStringValues)
			allRleValues[j] = rleStringValues
		}
	}
	for i := 0; i < int(sqlResponse.RowCount()); i++ {
		rs.rows[i] = newFlatBuffersSQLRow(allColValues, allRleValues, i, rs.meta)
	}

	return rs, nil
}

func (rs *TablestoreSQLResultSet) Columns() []*SQLColumnInfo {
	if rs.meta == nil {
		return nil
	}
	return rs.meta.columns
}

// Next returns the next row of SQLResultSet.
func (rs *TablestoreSQLResultSet) Next() SQLRow {
	if rs.rowCursor >= len(rs.rows) {
		return nil
	}
	row := rs.rows[rs.rowCursor]
	rs.rowCursor++
	return row
}

// HasNext returns whether finished.
func (rs *TablestoreSQLResultSet) HasNext() bool {
	return rs.rowCursor < len(rs.rows)
}

func (rs *TablestoreSQLResultSet) Reset() {
	rs.rowCursor = 0
}

type SQLTableMeta struct {
	columns   []*SQLColumnInfo
	colOffset map[string]int
}

func formatSQLTableMetaFromPlainBufferRow(row *PlainBufferRow) *SQLTableMeta {
	meta := new(SQLTableMeta)
	meta.colOffset = make(map[string]int)
	for idx, cell := range row.cells {
		columnInfo := &SQLColumnInfo{
			Name: string(cell.cellName),
			Type: cell.cellValue.Type,
		}
		meta.columns = append(meta.columns, columnInfo)
		meta.colOffset[columnInfo.Name] = idx
	}
	return meta
}

func formatColumnTypeFromFlatBuffers(typ dataprotocol.DataType) ColumnType {
	switch typ {
	case dataprotocol.DataTypeLONG:
		return ColumnType_INTEGER
	case dataprotocol.DataTypeBOOLEAN:
		return ColumnType_BOOLEAN
	case dataprotocol.DataTypeDOUBLE:
		return ColumnType_DOUBLE
	case dataprotocol.DataTypeSTRING:
		return ColumnType_STRING
	case dataprotocol.DataTypeBINARY:
		return ColumnType_BINARY
	case dataprotocol.DataTypeSTRING_RLE:
		return ColumnType_STRING
	default:
		return ColumnType(-1)
	}
}
func formatSQLTableMetaFromFlatBufferColumns(columns *dataprotocol.SQLResponseColumns) *SQLTableMeta {
	meta := new(SQLTableMeta)
	meta.colOffset = make(map[string]int)
	for i := 0; i < columns.ColumnsLength(); i++ {
		column := new(dataprotocol.SQLResponseColumn)
		columns.Columns(column, i)
		columnInfo := &SQLColumnInfo{
			Name:     string(column.ColumnName()),
			Type:     formatColumnTypeFromFlatBuffers(column.ColumnType()),
			dataType: column.ColumnType(),
		}
		meta.columns = append(meta.columns, columnInfo)
		meta.colOffset[columnInfo.Name] = i
	}
	return meta
}

func formatSQLStmtTypeFromPB(typ otsprotocol.SQLStatementType) SQLStatementType {
	switch typ {
	case otsprotocol.SQLStatementType_SQL_SELECT:
		return SQL_SELECT
	case otsprotocol.SQLStatementType_SQL_CREATE_TABLE:
		return SQL_CREATE_TABLE
	case otsprotocol.SQLStatementType_SQL_SHOW_TABLE:
		return SQL_SHOW_TABLE
	case otsprotocol.SQLStatementType_SQL_DESCRIBE_TABLE:
		return SQL_DESCRIBE_TABLE
	case otsprotocol.SQLStatementType_SQL_DROP_TABLE:
		return SQL_DROP_TABLE
	case otsprotocol.SQLStatementType_SQL_ALTER_TABLE:
		return SQL_ALTER_TABLE
	default:
		return SQLStatementType(-1)
	}
}

// SQLColumnInfo contains information of a column.
type SQLColumnInfo struct {
	Name string
	Type ColumnType

	dataType dataprotocol.DataType
}

// SQLRow represents a row of data, can be used to access values.
type SQLRow interface {
	// IsNull returns whether the value with the colIdx is nil.
	// 	When return true, means <code>NULL</code> in SQL.
	IsNull(colIdx int) (bool, error)

	// IsNullByName returns whether the value with the column name is nil.
	//	When return true, means <code>NULL</code> in SQL.
	IsNullByName(colName string) (bool, error)

	// GetString returns the string value with the colIdx.
	//	if the value is SQL <code>NULL</code>, the value returned is <code>""</code
	GetString(colIdx int) (string, error)

	// GetStringByName returns the string value with the column Name.
	GetStringByName(colName string) (string, error)

	// GetInt64 returns the int64 value with the colIdx.
	//	if the value is SQL <code>NULL</code>, the value returned is <code>0</code>
	GetInt64(colIdx int) (int64, error)

	// GetInt64ByName returns the int64 value with the column Name.
	GetInt64ByName(colName string) (int64, error)

	// GetBool returns the bool value with the colIdx.
	//	if the value is SQL <code>NULL</code>, the value returned is <code>false</code>
	GetBool(colIdx int) (bool, error)

	// GetBoolByName returns the bool value with the column Name.
	GetBoolByName(colName string) (bool, error)

	// GetBytes returns the bytes value with the colIdx.
	//	if the value is SQL <code>NULL</code>, the value returned is <code>nil</code>
	GetBytes(colIdx int) ([]byte, error)

	// GetBytesByName returns the bytes value with the column Name.
	GetBytesByName(colName string) ([]byte, error)

	// GetFloat64 returns the float64 value with the colIdx.
	//	if the value is SQL <code>NULL</code>, the value returned is <code>0</code>
	GetFloat64(colIdx int) (float64, error)

	// GetFloat64ByName returns the float64 value with the column Name.
	GetFloat64ByName(colName string) (float64, error)

	// DebugString for debug/test/print use
	DebugString() string
}

type flatBuffersSQLRow struct {
	data    []*dataprotocol.ColumnValues
	rleData []*dataprotocol.RLEStringValues
	rowIdx  int

	colOffset   map[string]int
	columnInfos []*SQLColumnInfo
}

func newFlatBuffersSQLRow(data []*dataprotocol.ColumnValues, rleData []*dataprotocol.RLEStringValues,
	rowIdx int, tableMeta *SQLTableMeta) SQLRow {
	row := new(flatBuffersSQLRow)
	row.data = data
	row.rleData = rleData
	row.rowIdx = rowIdx
	row.colOffset = tableMeta.colOffset
	row.columnInfos = tableMeta.columns
	return row
}

func (row *flatBuffersSQLRow) IsNull(colIdx int) (bool, error) {
	if colIdx < len(row.data) {
		colValues := row.data[colIdx]
		return colValues.IsNullvalues(row.rowIdx), nil
	} else {
		return false, errors.New(fmt.Sprintf("colIdx out of bound, max: %d", len(row.data)-1))
	}
}

func (row *flatBuffersSQLRow) IsNullByName(colName string) (bool, error) {
	if colIdx, ok := row.colOffset[colName]; ok {
		return row.IsNull(colIdx)
	} else {
		return false, errors.New("SQLRow doesn't contains Name: " + colName)
	}
}

func (row *flatBuffersSQLRow) GetString(colIdx int) (string, error) {
	if colIdx < len(row.data) {
		colInfo := row.columnInfos[colIdx]
		if colInfo.Type != ColumnType_STRING {
			return "", errors.New("the type of column is not STRING")
		}
		if colInfo.dataType == dataprotocol.DataTypeSTRING_RLE {
			rleStringValues := row.rleData[colIdx]
			if rleStringValues == nil {
				return "", errors.New("invalid RLE string values")
			}
			if row.rowIdx < rleStringValues.IndexMappingLength() {
				dataIdx := rleStringValues.IndexMapping(row.rowIdx)
				return string(rleStringValues.Array(int(dataIdx))), nil
			} else {
				return "", errors.New(fmt.Sprintf("rowIdx out of bound, max: %d", rleStringValues.IndexMappingLength()-1))
			}
		} else {
			colValues := row.data[colIdx]
			if row.rowIdx < colValues.StringValuesLength() {
				return string(colValues.StringValues(row.rowIdx)), nil
			} else {
				return "", errors.New(fmt.Sprintf("rowIdx out of bound, max: %d", colValues.StringValuesLength()-1))
			}
		}
	} else {
		return "", errors.New(fmt.Sprintf("colIdx out of bound, max: %d", len(row.data)-1))
	}
}

func (row *flatBuffersSQLRow) GetStringByName(colName string) (string, error) {
	if colIdx, ok := row.colOffset[colName]; ok {
		return row.GetString(colIdx)
	} else {
		return "", errors.New("SQLRow doesn't contains Name: " + colName)
	}
}

func (row *flatBuffersSQLRow) GetInt64(colIdx int) (int64, error) {
	if colIdx < len(row.data) {
		if row.columnInfos[colIdx].Type != ColumnType_INTEGER {
			return 0, errors.New("the type of column is not INTEGER")
		}
		colValues := row.data[colIdx]
		if row.rowIdx < colValues.LongValuesLength() {
			return colValues.LongValues(row.rowIdx), nil
		} else {
			return 0, errors.New(fmt.Sprintf("rowIdx out of bound, max: %d", colValues.LongValuesLength()-1))
		}
	} else {
		return 0, errors.New(fmt.Sprintf("colIdx out of bound, max: %d", len(row.data)-1))
	}
}

func (row *flatBuffersSQLRow) GetInt64ByName(colName string) (int64, error) {
	if colIdx, ok := row.colOffset[colName]; ok {
		return row.GetInt64(colIdx)
	} else {
		return 0, errors.New("SQLRow doesn't contains Name: " + colName)
	}
}

func (row *flatBuffersSQLRow) GetBool(colIdx int) (bool, error) {
	if colIdx < len(row.data) {
		if row.columnInfos[colIdx].Type != ColumnType_BOOLEAN {
			return false, errors.New("the type of column is not BOOLEAN")
		}
		colValues := row.data[colIdx]
		if row.rowIdx < colValues.BoolValuesLength() {
			return colValues.BoolValues(row.rowIdx), nil
		} else {
			return false, errors.New(fmt.Sprintf("rowIdx out of bound, max: %d", colValues.BoolValuesLength()-1))
		}
	} else {
		return false, errors.New(fmt.Sprintf("colIdx out of bound, max: %d", len(row.data)-1))
	}
}

func (row *flatBuffersSQLRow) GetBoolByName(colName string) (bool, error) {
	if colIdx, ok := row.colOffset[colName]; ok {
		return row.GetBool(colIdx)
	} else {
		return false, errors.New("SQLRow doesn't contains Name: " + colName)
	}
}

func (row *flatBuffersSQLRow) GetBytes(colIdx int) ([]byte, error) {
	if colIdx < len(row.data) {
		if row.columnInfos[colIdx].Type != ColumnType_BINARY {
			return nil, errors.New("the type of column is not BINARY")
		}
		colValues := row.data[colIdx]
		if row.rowIdx < colValues.BinaryValuesLength() {
			bytesVal := new(dataprotocol.BytesValue)
			_ = colValues.BinaryValues(bytesVal, row.rowIdx)
			retBytes := make([]byte, bytesVal.ValueLength())
			for idx := 0; idx < bytesVal.ValueLength(); idx++ {
				retBytes[idx] = byte(bytesVal.Value(idx))
			}
			return retBytes, nil
		} else {
			return nil, errors.New(fmt.Sprintf("rowIdx out of bound, max: %d", colValues.BinaryValuesLength()-1))
		}
	} else {
		return nil, errors.New(fmt.Sprintf("colIdx out of bound, max: %d", len(row.data)-1))
	}
}

func (row *flatBuffersSQLRow) GetBytesByName(colName string) ([]byte, error) {
	if colIdx, ok := row.colOffset[colName]; ok {
		return row.GetBytes(colIdx)
	} else {
		return nil, errors.New("SQLRow doesn't contains Name: " + colName)
	}
}

func (row *flatBuffersSQLRow) GetFloat64(colIdx int) (float64, error) {
	if colIdx < len(row.data) {
		if row.columnInfos[colIdx].Type != ColumnType_DOUBLE {
			return 0, errors.New("the type of column is not DOUBLE")
		}
		colValues := row.data[colIdx]
		if row.rowIdx < colValues.DoubleValuesLength() {
			return colValues.DoubleValues(row.rowIdx), nil
		} else {
			return 0, errors.New(fmt.Sprintf("rowIdx out of bound, max: %d", colValues.DoubleValuesLength()-1))
		}
	} else {
		return 0, errors.New(fmt.Sprintf("colIdx out of bound, max: %d", len(row.data)-1))
	}
}

func (row *flatBuffersSQLRow) GetFloat64ByName(colName string) (float64, error) {
	if colIdx, ok := row.colOffset[colName]; ok {
		return row.GetFloat64(colIdx)
	} else {
		return 0, errors.New("SQLRow doesn't contains Name: " + colName)
	}
}

func (row *flatBuffersSQLRow) DebugString() string {
	rowValues := make([]interface{}, len(row.colOffset))
	for i := 0; i < len(row.colOffset); i++ {
		colValues := row.data[i]
		isnull := colValues.IsNullvalues(row.rowIdx)
		if isnull {
			continue
		}
		rleStringValues := row.rleData[i]
		if colValues.StringValuesLength() > row.rowIdx {
			rowValues[i] = string(colValues.StringValues(row.rowIdx))
		} else if colValues.LongValuesLength() > row.rowIdx {
			rowValues[i] = colValues.LongValues(row.rowIdx)
		} else if colValues.BinaryValuesLength() > row.rowIdx {
			bytesVal := new(dataprotocol.BytesValue)
			_ = colValues.BinaryValues(bytesVal, row.rowIdx)
			retBytes := make([]byte, bytesVal.ValueLength())
			for idx := 0; idx < bytesVal.ValueLength(); idx++ {
				retBytes[idx] = byte(bytesVal.Value(idx))
			}
			rowValues[i] = retBytes
		} else if colValues.DoubleValuesLength() > row.rowIdx {
			rowValues[i] = colValues.DoubleValues(row.rowIdx)
		} else if colValues.BoolValuesLength() > row.rowIdx {
			rowValues[i] = colValues.BoolValues(row.rowIdx)
		} else if rleStringValues.IndexMappingLength() > row.rowIdx {
			dataIdx := rleStringValues.IndexMapping(row.rowIdx)
			rowValues[i] = string(rleStringValues.Array(int(dataIdx)))
		}
	}
	rowBytes, err := json.Marshal(rowValues)
	if err != nil {
		return ""
	}
	return string(rowBytes)
}

// plainBufferSQLRow Legacy SQL payload format
type plainBufferSQLRow struct {
	cellValues []interface{}

	colOffset map[string]int
}

func newPlainBufferSQLRow(pbRow *PlainBufferRow, colOffset map[string]int) SQLRow {
	row := new(plainBufferSQLRow)

	cellValues := make([]interface{}, len(colOffset))
	for _, cell := range pbRow.cells {
		idx := colOffset[string(cell.cellName)]
		cellValues[idx] = cell.cellValue.Value
	}

	row.cellValues = cellValues
	row.colOffset = colOffset
	return row
}

func (row *plainBufferSQLRow) IsNull(colIdx int) (bool, error) {
	if colIdx < len(row.cellValues) {
		return row.cellValues[colIdx] == nil, nil
	} else {
		return false, errors.New(fmt.Sprintf("colIdx out of bound, max: %d", len(row.colOffset)-1))
	}
}

func (row *plainBufferSQLRow) IsNullByName(colName string) (bool, error) {
	if colIdx, ok := row.colOffset[colName]; ok {
		return row.IsNull(colIdx)
	} else {
		return false, errors.New("SQLRow doesn't contains Name: " + colName)
	}
}

func (row *plainBufferSQLRow) GetString(colIdx int) (string, error) {
	if colIdx < len(row.cellValues) {
		if row.cellValues[colIdx] == nil {
			return "", nil
		}
		if val, ok := row.cellValues[colIdx].(string); ok {
			return val, nil
		} else {
			return "", errors.New("the type of column is not STRING")
		}
	} else {
		return "", errors.New(fmt.Sprintf("colIdx out of bound, max: %d", len(row.cellValues)-1))
	}
}

func (row *plainBufferSQLRow) GetStringByName(colName string) (string, error) {
	if colIdx, ok := row.colOffset[colName]; ok {
		return row.GetString(colIdx)
	} else {
		return "", errors.New("SQLRow doesn't contains Name: " + colName)
	}
}

func (row *plainBufferSQLRow) GetInt64(colIdx int) (int64, error) {
	if colIdx < len(row.cellValues) {
		if row.cellValues[colIdx] == nil {
			return 0, nil
		}
		if val, ok := row.cellValues[colIdx].(int64); ok {
			return val, nil
		} else {
			return 0, errors.New("the type of column is not INTEGER")
		}
	} else {
		return 0, errors.New(fmt.Sprintf("colIdx out of bound, max: %d", len(row.cellValues)-1))
	}
}

func (row *plainBufferSQLRow) GetInt64ByName(colName string) (int64, error) {
	if colIdx, ok := row.colOffset[colName]; ok {
		return row.GetInt64(colIdx)
	} else {
		return 0, errors.New("SQLRow doesn't contains Name: " + colName)
	}
}

func (row *plainBufferSQLRow) GetBool(colIdx int) (bool, error) {
	if colIdx < len(row.cellValues) {
		if row.cellValues[colIdx] == nil {
			return false, nil
		}
		if val, ok := row.cellValues[colIdx].(bool); ok {
			return val, nil
		} else {
			return false, errors.New("the type of column is not BOOL")
		}
	} else {
		return false, errors.New(fmt.Sprintf("colIdx out of bound, max: %d", len(row.cellValues)-1))
	}
}

func (row *plainBufferSQLRow) GetBoolByName(colName string) (bool, error) {
	if colIdx, ok := row.colOffset[colName]; ok {
		return row.GetBool(colIdx)
	} else {
		return false, errors.New("SQLRow doesn't contains Name: " + colName)
	}
}

func (row *plainBufferSQLRow) GetBytes(colIdx int) ([]byte, error) {
	if colIdx < len(row.cellValues) {
		if row.cellValues[colIdx] == nil {
			return nil, nil
		}
		if val, ok := row.cellValues[colIdx].([]byte); ok {
			return val, nil
		} else {
			return nil, errors.New("the type of column is not BINARY")
		}
	} else {
		return nil, errors.New(fmt.Sprintf("colIdx out of bound, max: %d", len(row.cellValues)-1))
	}
}

func (row *plainBufferSQLRow) GetBytesByName(colName string) ([]byte, error) {
	if colIdx, ok := row.colOffset[colName]; ok {
		return row.GetBytes(colIdx)
	} else {
		return nil, errors.New("SQLRow doesn't contains Name: " + colName)
	}
}

func (row *plainBufferSQLRow) GetFloat64(colIdx int) (float64, error) {
	if colIdx < len(row.cellValues) {
		if row.cellValues[colIdx] == nil {
			return 0, nil
		}
		if val, ok := row.cellValues[colIdx].(float64); ok {
			return val, nil
		} else {
			return 0, errors.New("the type of column is not DOUBLE")
		}
	} else {
		return 0, errors.New(fmt.Sprintf("colIdx out of bound, max: %d", len(row.cellValues)-1))
	}
}

func (row *plainBufferSQLRow) GetFloat64ByName(colName string) (float64, error) {
	if colIdx, ok := row.colOffset[colName]; ok {
		return row.GetFloat64(colIdx)
	} else {
		return 0, errors.New("SQLRow doesn't contains Name: " + colName)
	}
}

func (row *plainBufferSQLRow) DebugString() string {
	rowBytes, err := json.Marshal(row.cellValues)
	if err != nil {
		return ""
	}
	return string(rowBytes)
}

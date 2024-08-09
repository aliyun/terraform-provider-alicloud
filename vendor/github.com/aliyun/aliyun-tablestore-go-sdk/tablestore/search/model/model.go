package model

import (
	"fmt"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"strconv"
	"strings"
	//"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
)

type PrimaryKey struct {
	PrimaryKeys []*PrimaryKeyColumn
}

type PrimaryKeyType int32

type PrimaryKeyOption int32

const (
	NONE           PrimaryKeyOption = 0
	AUTO_INCREMENT PrimaryKeyOption = 1
	MIN            PrimaryKeyOption = 2
	MAX            PrimaryKeyOption = 3
)

type PrimaryKeyColumn struct {
	ColumnName       string
	Value            interface{}
	PrimaryKeyOption PrimaryKeyOption
}

func (this *PrimaryKeyColumn) String() string {
	xs := make([]string, 0)
	xs = append(xs, fmt.Sprintf("\"Name\": \"%s\"", this.ColumnName))
	switch this.PrimaryKeyOption {
	case NONE:
		xs = append(xs, fmt.Sprintf("\"Value\": \"%s\"", this.Value))
	case MIN:
		xs = append(xs, "\"Value\": -inf")
	case MAX:
		xs = append(xs, "\"Value\": +inf")
	case AUTO_INCREMENT:
		xs = append(xs, "\"Value\": auto-incr")
	}
	return fmt.Sprintf("{%s}", strings.Join(xs, ", "))
}

type AttributeColumn struct {
	ColumnName string
	Value      interface{}
	Timestamp  int64
}

type TimeRange struct {
	Start    int64
	End      int64
	Specific int64
}

type ColumnToUpdate struct {
	ColumnName   string
	Type         byte
	Timestamp    int64
	HasType      bool
	HasTimestamp bool
	IgnoreValue  bool
	Value        interface{}
}

type RowExistenceExpectation int

const (
	RowExistenceExpectation_IGNORE           RowExistenceExpectation = 0
	RowExistenceExpectation_EXPECT_EXIST     RowExistenceExpectation = 1
	RowExistenceExpectation_EXPECT_NOT_EXIST RowExistenceExpectation = 2
)

type ComparatorType int32

const (
	CT_EQUAL         ComparatorType = 1
	CT_NOT_EQUAL     ComparatorType = 2
	CT_GREATER_THAN  ComparatorType = 3
	CT_GREATER_EQUAL ComparatorType = 4
	CT_LESS_THAN     ComparatorType = 5
	CT_LESS_EQUAL    ComparatorType = 6
)

type LogicalOperator int32

const (
	LO_NOT LogicalOperator = 1
	LO_AND LogicalOperator = 2
	LO_OR  LogicalOperator = 3
)

type FilterType int32

type ColumnFilter interface {
	Serialize() []byte
	ToFilter() *otsprotocol.Filter
}

type VariantType int32

type ValueTransferRule struct {
	Regex     string
	Cast_type VariantType
}

type SingleColumnCondition struct {
	Comparator        *ComparatorType
	ColumnName        *string
	ColumnValue       interface{} //[]byte
	FilterIfMissing   bool
	LatestVersionOnly bool
	TransferRule      *ValueTransferRule
}

type ReturnType int32

const (
	ReturnType_RT_NONE         ReturnType = 0
	ReturnType_RT_PK           ReturnType = 1
	ReturnType_RT_AFTER_MODIFY ReturnType = 2
)

type RowCondition struct {
	RowExistenceExpectation RowExistenceExpectation
	ColumnCondition         ColumnFilter
}

type PutRowChange struct {
	TableName     string
	PrimaryKey    *PrimaryKey
	Columns       []AttributeColumn
	Condition     *RowCondition
	ReturnType    ReturnType
	TransactionId *string
}

type PutRowRequest struct {
	PutRowChange *PutRowChange
}

type DeleteRowChange struct {
	TableName     string
	PrimaryKey    *PrimaryKey
	Condition     *RowCondition
	TransactionId *string
}

type DeleteRowRequest struct {
	DeleteRowChange *DeleteRowChange
}

type SingleRowQueryCriteria struct {
	ColumnsToGet  []string
	TableName     string
	PrimaryKey    *PrimaryKey
	MaxVersion    int32
	TimeRange     *TimeRange
	Filter        ColumnFilter
	StartColumn   *string
	EndColumn     *string
	TransactionId *string
}

type UpdateRowChange struct {
	TableName           string
	PrimaryKey          *PrimaryKey
	Columns             []ColumnToUpdate
	Condition           *RowCondition
	TransactionId       *string
	ReturnType          ReturnType
	ColumnNamesToReturn []string
}

type UpdateRowRequest struct {
	UpdateRowChange *UpdateRowChange
}

func (rowQueryCriteria *SingleRowQueryCriteria) AddColumnToGet(columnName string) {
	rowQueryCriteria.ColumnsToGet = append(rowQueryCriteria.ColumnsToGet, columnName)
}

func (rowQueryCriteria *SingleRowQueryCriteria) SetStartColumn(columnName string) {
	rowQueryCriteria.StartColumn = &columnName
}

func (rowQueryCriteria *SingleRowQueryCriteria) SetEndtColumn(columnName string) {
	rowQueryCriteria.EndColumn = &columnName
}

func (rowQueryCriteria *SingleRowQueryCriteria) getColumnsToGet() []string {
	return rowQueryCriteria.ColumnsToGet
}

func (rowQueryCriteria *MultiRowQueryCriteria) AddColumnToGet(columnName string) {
	rowQueryCriteria.ColumnsToGet = append(rowQueryCriteria.ColumnsToGet, columnName)
}

func (rowQueryCriteria *RangeRowQueryCriteria) AddColumnToGet(columnName string) {
	rowQueryCriteria.ColumnsToGet = append(rowQueryCriteria.ColumnsToGet, columnName)
}

func (rowQueryCriteria *MultiRowQueryCriteria) AddRow(pk *PrimaryKey) {
	rowQueryCriteria.PrimaryKey = append(rowQueryCriteria.PrimaryKey, pk)
}

type GetRowRequest struct {
	SingleRowQueryCriteria *SingleRowQueryCriteria
}

type MultiRowQueryCriteria struct {
	PrimaryKey   []*PrimaryKey
	ColumnsToGet []string
	TableName    string
	MaxVersion   int
	TimeRange    *TimeRange
	Filter       ColumnFilter
	StartColumn  *string
	EndColumn    *string
}

type BatchGetRowRequest struct {
	MultiRowQueryCriteria []*MultiRowQueryCriteria
}

type ColumnMap struct {
	Columns    map[string][]*AttributeColumn
	columnsKey []string
}

type Error struct {
	Code    string
	Message string
}

type RowChange interface {
	Serialize() []byte
	getOperationType() otsprotocol.OperationType
	getCondition() *otsprotocol.Condition
	GetTableName() string
}

type Direction int32

const (
	FORWARD  Direction = 0
	BACKWARD Direction = 1
)

type RangeRowQueryCriteria struct {
	TableName       string
	StartPrimaryKey *PrimaryKey
	EndPrimaryKey   *PrimaryKey
	ColumnsToGet    []string
	MaxVersion      int32
	TimeRange       *TimeRange
	Filter          ColumnFilter
	Direction       Direction
	Limit           int32
	StartColumn     *string
	EndColumn       *string
	TransactionId   *string
}

type GetRangeRequest struct {
	RangeRowQueryCriteria *RangeRowQueryCriteria
}

type Row struct {
	PrimaryKey *PrimaryKey
	Columns    []*AttributeColumn
}

type StreamSpecification struct {
	EnableStream   bool
	ExpirationTime int32 // must be positive. in hours
}

type StreamDetails struct {
	EnableStream   bool
	StreamId       *StreamId // nil when stream is disabled.
	ExpirationTime int32     // in hours
	LastEnableTime int64     // the last time stream is enabled, in usec
}

type Split struct {
	LowerBound *PrimaryKey
	UpperBound *PrimaryKey
	Location   string
}

type StreamId string
type ShardId string
type ShardIterator string
type StreamStatus int

/*
 * Shards are possibly splitted into two or merged from two.
 * After splitting, both newly generated shards have the same FatherShard.
 * After merging, the newly generated shard have both FatherShard and MotherShard.
 */
type StreamShard struct {
	SelfShard   *ShardId // required
	FatherShard *ShardId // optional
	MotherShard *ShardId // optional
}

type StreamRecord struct {
	Type       ActionType
	Info       *RecordSequenceInfo // required
	PrimaryKey *PrimaryKey         // required
	Columns    []*RecordColumn
}

func (this *StreamRecord) String() string {
	return fmt.Sprintf(
		"{\"Type\":%s, \"PrimaryKey\":%s, \"Info\":%s, \"Columns\":%s}",
		this.Type,
		*this.PrimaryKey,
		this.Info,
		this.Columns)
}

type ActionType int

const (
	AT_Put ActionType = iota
	AT_Update
	AT_Delete
)

func (this ActionType) String() string {
	switch this {
	case AT_Put:
		return "\"PutRow\""
	case AT_Update:
		return "\"UpdateRow\""
	case AT_Delete:
		return "\"DeleteRow\""
	default:
		panic(fmt.Sprintf("unknown action type: %d", int(this)))
	}
}

type RecordSequenceInfo struct {
	Epoch     int32
	Timestamp int64
	RowIndex  int32
}

func (this *RecordSequenceInfo) String() string {
	return fmt.Sprintf(
		"{\"Epoch\":%d, \"Timestamp\": %d, \"RowIndex\": %d}",
		this.Epoch,
		this.Timestamp,
		this.RowIndex)
}

type RecordColumn struct {
	Type      RecordColumnType
	Name      *string     // required
	Value     interface{} // optional. present when Type is RCT_Put
	Timestamp *int64      // optional, in msec. present when Type is RCT_Put or RCT_DeleteOneVersion
}

func (this *RecordColumn) String() string {
	xs := make([]string, 0)
	xs = append(xs, fmt.Sprintf("\"Name\":%s", strconv.Quote(*this.Name)))
	switch this.Type {
	case RCT_DeleteAllVersions:
		xs = append(xs, "\"Type\":\"DeleteAllVersions\"")
	case RCT_DeleteOneVersion:
		xs = append(xs, "\"Type\":\"DeleteOneVersion\"")
		xs = append(xs, fmt.Sprintf("\"Timestamp\":%d", *this.Timestamp))
	case RCT_Put:
		xs = append(xs, "\"Type\":\"Put\"")
		xs = append(xs, fmt.Sprintf("\"Timestamp\":%d", *this.Timestamp))
		xs = append(xs, fmt.Sprintf("\"Value\":%s", this.Value))
	}
	return fmt.Sprintf("{%s}", strings.Join(xs, ", "))
}

type RecordColumnType int

const (
	RCT_Put RecordColumnType = iota
	RCT_DeleteOneVersion
	RCT_DeleteAllVersions
)

type IndexMeta struct {
	IndexName      string
	Primarykey     []string
	DefinedColumns []string
	IndexType      IndexType
}

type DefinedColumnSchema struct {
	Name       string
	ColumnType DefinedColumnType
}

type IndexType int32

type DefinedColumnType int32

const (
	/**
	 * 64位整数。
	 */
	DefinedColumn_INTEGER DefinedColumnType = 1

	/**
	 * 浮点数。
	 */
	DefinedColumn_DOUBLE DefinedColumnType = 2

	/**
	 * 布尔值。
	 */
	DefinedColumn_BOOLEAN DefinedColumnType = 3

	/**
	 * 字符串。
	 */
	DefinedColumn_STRING DefinedColumnType = 4

	/**
	 * BINARY。
	 */
	DefinedColumn_BINARY DefinedColumnType = 5
)

type StartLocalTransactionRequest struct {
	PrimaryKey *PrimaryKey
	TableName  string
}

// compute splits
type SearchIndexSplitsOptions struct {
	IndexName string
}

type ComputeSplitsRequest struct {
	TableName                string
	searchIndexSplitsOptions *SearchIndexSplitsOptions
}

type DateTimeValue struct {
	Value *int32
	Unit  *DateTimeUnit
}

type DateTimeUnit int32

const (
	DateTimeUnit_YEAR         DateTimeUnit = 1
	DateTimeUnit_QUARTER_YEAR DateTimeUnit = 2
	DateTimeUnit_MONTH        DateTimeUnit = 3
	DateTimeUnit_WEEK         DateTimeUnit = 4
	DateTimeUnit_DAY          DateTimeUnit = 5
	DateTimeUnit_HOUR         DateTimeUnit = 6
	DateTimeUnit_MINUTE       DateTimeUnit = 7
	DateTimeUnit_SECOND       DateTimeUnit = 8
)

func (d DateTimeUnit) Enum() *DateTimeUnit {
	p := new(DateTimeUnit)
	*p = d
	return p
}

func (d *DateTimeUnit) ProtoBuffer() *otsprotocol.DateTimeUnit {
	if d == nil {
		return nil
	}
	switch *d {
	case DateTimeUnit_YEAR:
		return otsprotocol.DateTimeUnit_YEAR.Enum()
	case DateTimeUnit_QUARTER_YEAR:
		return otsprotocol.DateTimeUnit_QUARTER_YEAR.Enum()
	case DateTimeUnit_MONTH:
		return otsprotocol.DateTimeUnit_MONTH.Enum()
	case DateTimeUnit_WEEK:
		return otsprotocol.DateTimeUnit_WEEK.Enum()
	case DateTimeUnit_DAY:
		return otsprotocol.DateTimeUnit_DAY.Enum()
	case DateTimeUnit_HOUR:
		return otsprotocol.DateTimeUnit_HOUR.Enum()
	case DateTimeUnit_MINUTE:
		return otsprotocol.DateTimeUnit_MINUTE.Enum()
	case DateTimeUnit_SECOND:
		return otsprotocol.DateTimeUnit_SECOND.Enum()
	default:
		return nil
	}
}

func (d *DateTimeValue) ProtoBuffer() *otsprotocol.DateTimeValue {
	dateTimeValue := &otsprotocol.DateTimeValue{}
	dateTimeValue.Unit = d.Unit.ProtoBuffer()
	dateTimeValue.Value = d.Value
	return dateTimeValue
}

type GeoHashPrecision int32

const (
	GHP_5009KM_4992KM_1 GeoHashPrecision = 1
	GHP_1252KM_624KM_2  GeoHashPrecision = 2
	GHP_156KM_156KM_3   GeoHashPrecision = 3
	GHP_39KM_19KM_4     GeoHashPrecision = 4
	GHP_4900M_4900M_5   GeoHashPrecision = 5
	GHP_1200M_609M_6    GeoHashPrecision = 6
	GHP_152M_152M_7     GeoHashPrecision = 7
	GHP_38M_19M_8       GeoHashPrecision = 8
	GHP_480CM_480CM_9   GeoHashPrecision = 9
	GHP_120CM_595MM_10  GeoHashPrecision = 10
	GHP_149MM_149MM_11  GeoHashPrecision = 11
	GHP_37MM_19MM_12    GeoHashPrecision = 12
)

func (d *GeoHashPrecision) ProtoBuffer() *otsprotocol.GeoHashPrecision {
	if d == nil {
		return nil
	}
	switch *d {
	case GHP_5009KM_4992KM_1:
		return otsprotocol.GeoHashPrecision_GHP_5009KM_4992KM_1.Enum()
	case GHP_1252KM_624KM_2:
		return otsprotocol.GeoHashPrecision_GHP_1252KM_624KM_2.Enum()
	case GHP_156KM_156KM_3:
		return otsprotocol.GeoHashPrecision_GHP_156KM_156KM_3.Enum()
	case GHP_39KM_19KM_4:
		return otsprotocol.GeoHashPrecision_GHP_39KM_19KM_4.Enum()
	case GHP_4900M_4900M_5:
		return otsprotocol.GeoHashPrecision_GHP_4900M_4900M_5.Enum()
	case GHP_1200M_609M_6:
		return otsprotocol.GeoHashPrecision_GHP_1200M_609M_6.Enum()
	case GHP_152M_152M_7:
		return otsprotocol.GeoHashPrecision_GHP_152M_152M_7.Enum()
	case GHP_38M_19M_8:
		return otsprotocol.GeoHashPrecision_GHP_38M_19M_8.Enum()
	case GHP_480CM_480CM_9:
		return otsprotocol.GeoHashPrecision_GHP_480CM_480CM_9.Enum()
	case GHP_120CM_595MM_10:
		return otsprotocol.GeoHashPrecision_GHP_120CM_595MM_10.Enum()
	case GHP_149MM_149MM_11:
		return otsprotocol.GeoHashPrecision_GHP_149MM_149MM_11.Enum()
	case GHP_37MM_19MM_12:
		return otsprotocol.GeoHashPrecision_GHP_37MM_19MM_12.Enum()
	default:
		return nil
	}
}

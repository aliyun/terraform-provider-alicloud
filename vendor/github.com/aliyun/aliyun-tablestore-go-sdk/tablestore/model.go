package tablestore

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/aliyun-tablestore-go-sdk/common"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/golang/protobuf/proto"
	lruCache "github.com/hashicorp/golang-lru"
	"sync"
)

type internalClient struct {
	endPoint        string
	instanceName    string
	accessKeyId     string
	accessKeySecret string
	securityToken   string

	httpClient IHttpClient
	config     *TableStoreConfig
	random     *rand.Rand
	mu         *sync.Mutex

	externalHeader      map[string]string
	CustomizedRetryFunc CustomizedRetryNotMatterActions

	timeseriesConfiguration *TimeseriesConfiguration
	credentialsProvider     common.CredentialsProvider
}

const initMapLen int = 8

// @class TableStoreClient
// The TableStoreClient, which will connect OTS service for authorization, create/list/
// delete tables/table groups, to get/put/delete a row.
// Note: TableStoreClient is thread-safe.
// TableStoreClient的功能包括连接OTS服务进行验证、创建/列出/删除表或表组、插入/获取/
// 删除/更新行数据
type TableStoreClient struct {
	*internalClient
}

type ClientOption func(*TableStoreClient)

type TimeseriesClientOption func(client *TimeseriesClient)

type TableStoreHttpClient struct {
	httpClient *http.Client
}

// use this to mock http.client for testing
type IHttpClient interface {
	Do(*http.Request) (*http.Response, error)
	New(*http.Client)
}

func (httpClient *TableStoreHttpClient) Do(req *http.Request) (*http.Response, error) {
	return httpClient.httpClient.Do(req)
}

func (httpClient *TableStoreHttpClient) New(client *http.Client) {
	httpClient.httpClient = client
}

type HTTPTimeout struct {
	ConnectionTimeout time.Duration
	RequestTimeout    time.Duration
}

type TableStoreConfig struct {
	RetryTimes         uint
	MaxRetryTime       time.Duration
	HTTPTimeout        HTTPTimeout
	MaxIdleConnections int
	Transport          http.RoundTripper
}

func NewDefaultTableStoreConfig() *TableStoreConfig {
	httpTimeout := &HTTPTimeout{
		ConnectionTimeout: time.Second * 15,
		RequestTimeout:    time.Second * 30}
	config := &TableStoreConfig{
		RetryTimes:         10,
		HTTPTimeout:        *httpTimeout,
		MaxRetryTime:       time.Second * 5,
		MaxIdleConnections: 2000}
	return config
}

type CreateTableRequest struct {
	TableMeta          *TableMeta
	TableOption        *TableOption
	ReservedThroughput *ReservedThroughput
	StreamSpec         *StreamSpecification
	IndexMetas         []*IndexMeta
	SSESpecification   *SSESpecification
}

type CreateIndexRequest struct {
	MainTableName   string
	IndexMeta       *IndexMeta
	IncludeBaseData bool
}

type DeleteIndexRequest struct {
	MainTableName string
	IndexName     string
}

type ResponseInfo struct {
	RequestId string
}

type CreateTableResponse struct {
	ResponseInfo
}

type CreateIndexResponse struct {
	ResponseInfo
}

type DeleteIndexResponse struct {
	ResponseInfo
}

type DeleteTableResponse struct {
	ResponseInfo
}

type TableMeta struct {
	TableName      string
	SchemaEntry    []*PrimaryKeySchema
	DefinedColumns []*DefinedColumnSchema
}

type PrimaryKeySchema struct {
	Name   *string
	Type   *PrimaryKeyType
	Option *PrimaryKeyOption
}

type PrimaryKey struct {
	PrimaryKeys []*PrimaryKeyColumn
}

type TableOption struct {
	TimeToAlive, MaxVersion   int
	DeviationCellVersionInSec int64
	AllowUpdate               *bool
}

type ReservedThroughput struct {
	Readcap, Writecap int
}

type ListTableResponse struct {
	TableNames []string
	ResponseInfo
}

type DeleteTableRequest struct {
	TableName string
}

type DescribeTableRequest struct {
	TableName string
}

type DescribeTableResponse struct {
	TableMeta          *TableMeta
	TableOption        *TableOption
	ReservedThroughput *ReservedThroughput
	StreamDetails      *StreamDetails
	IndexMetas         []*IndexMeta
	SSEDetails         *SSEDetails
	ResponseInfo
}

type UpdateTableRequest struct {
	TableName          string
	TableOption        *TableOption
	ReservedThroughput *ReservedThroughput
	StreamSpec         *StreamSpecification
}

type UpdateTableResponse struct {
	TableOption        *TableOption
	ReservedThroughput *ReservedThroughput
	StreamDetails      *StreamDetails
	ResponseInfo
}

type AddDefinedColumnRequest struct {
	TableName      string
	DefinedColumns []*DefinedColumnSchema
}

type DeleteDefinedColumnRequest struct {
	TableName      string
	DefinedColumns []string
}

type AddDefinedColumnResponse struct {
	ResponseInfo
}

type DeleteDefinedColumnResponse struct {
	ResponseInfo
}

type ConsumedCapacityUnit struct {
	Read  int32
	Write int32
}

type PutRowResponse struct {
	ConsumedCapacityUnit *ConsumedCapacityUnit
	PrimaryKey           PrimaryKey
	ResponseInfo
}

type DeleteRowResponse struct {
	ConsumedCapacityUnit *ConsumedCapacityUnit
	ResponseInfo
}

type UpdateRowResponse struct {
	Columns              []*AttributeColumn
	ConsumedCapacityUnit *ConsumedCapacityUnit
	ResponseInfo
}

type PrimaryKeyType int32

const (
	PrimaryKeyType_INTEGER PrimaryKeyType = 1
	PrimaryKeyType_STRING  PrimaryKeyType = 2
	PrimaryKeyType_BINARY  PrimaryKeyType = 3
)

const (
	DefaultRetryInterval = 10
	MaxRetryInterval     = 320
)

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

const (
	FT_SINGLE_COLUMN_VALUE    FilterType = 1
	FT_COMPOSITE_COLUMN_VALUE FilterType = 2
	FT_COLUMN_PAGINATION      FilterType = 3
)

type ColumnFilter interface {
	Serialize() []byte
	ToFilter() *otsprotocol.Filter
}

type VariantType int32

const (
	Variant_INTEGER VariantType = 0
	Variant_DOUBLE  VariantType = 1
	//VT_BOOLEAN = 2;
	Variant_STRING VariantType = 3
)

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

type PaginationFilter struct {
	Offset int32
	Limit  int32
}

type CompositeColumnValueFilter struct {
	Operator LogicalOperator
	Filters  []ColumnFilter
}

func (ccvfilter *CompositeColumnValueFilter) Serialize() []byte {
	result, _ := proto.Marshal(ccvfilter.ToFilter())
	return result
}

func (ccvfilter *CompositeColumnValueFilter) ToFilter() *otsprotocol.Filter {
	compositefilter := NewCompositeFilter(ccvfilter.Filters, ccvfilter.Operator)
	compositeFilterToBytes, _ := proto.Marshal(compositefilter)
	filter := new(otsprotocol.Filter)
	filter.Type = otsprotocol.FilterType_FT_COMPOSITE_COLUMN_VALUE.Enum()
	filter.Filter = compositeFilterToBytes
	return filter
}

func (ccvfilter *CompositeColumnValueFilter) AddFilter(filter ColumnFilter) {
	ccvfilter.Filters = append(ccvfilter.Filters, filter)
}

func (condition *SingleColumnCondition) ToFilter() *otsprotocol.Filter {
	singlefilter := NewSingleColumnValueFilter(condition)
	singleFilterToBytes, _ := proto.Marshal(singlefilter)
	filter := new(otsprotocol.Filter)
	filter.Type = otsprotocol.FilterType_FT_SINGLE_COLUMN_VALUE.Enum()
	filter.Filter = singleFilterToBytes
	return filter
}

func (condition *SingleColumnCondition) Serialize() []byte {
	result, _ := proto.Marshal(condition.ToFilter())
	return result
}

func (pageFilter *PaginationFilter) ToFilter() *otsprotocol.Filter {
	compositefilter := NewPaginationFilter(pageFilter)
	compositeFilterToBytes, _ := proto.Marshal(compositefilter)
	filter := new(otsprotocol.Filter)
	filter.Type = otsprotocol.FilterType_FT_COLUMN_PAGINATION.Enum()
	filter.Filter = compositeFilterToBytes
	return filter
}

func (pageFilter *PaginationFilter) Serialize() []byte {
	result, _ := proto.Marshal(pageFilter.ToFilter())
	return result
}

func NewTableOptionWithMaxVersion(maxVersion int) *TableOption {
	tableOption := new(TableOption)
	tableOption.TimeToAlive = -1
	tableOption.MaxVersion = maxVersion
	return tableOption
}

func NewTableOption(timeToAlive int, maxVersion int) *TableOption {
	tableOption := new(TableOption)
	tableOption.TimeToAlive = timeToAlive
	tableOption.MaxVersion = maxVersion
	return tableOption
}

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

type GetRowResponse struct {
	PrimaryKey           PrimaryKey
	Columns              []*AttributeColumn
	ConsumedCapacityUnit *ConsumedCapacityUnit
	columnMap            *ColumnMap
	ResponseInfo
}

type Error struct {
	Code    string
	Message string
}

type RowResult struct {
	TableName            string
	IsSucceed            bool
	Error                Error
	PrimaryKey           PrimaryKey
	Columns              []*AttributeColumn
	ConsumedCapacityUnit *ConsumedCapacityUnit
	Index                int32
}

type RowChange interface {
	Serialize() []byte
	getOperationType() otsprotocol.OperationType
	getCondition() *otsprotocol.Condition
	GetTableName() string
}

type BatchGetRowResponse struct {
	TableToRowsResult map[string][]RowResult
	ResponseInfo
}

//IsAtomic设置是否为批量原子写
//如果设置了批量原子写，需要保证写入到同一张表格中的分区键相同，否则会写入失败
type BatchWriteRowRequest struct {
	RowChangesGroupByTable map[string][]RowChange
	IsAtomic               bool
}

type BatchWriteRowResponse struct {
	TableToRowsResult map[string][]RowResult
	ResponseInfo
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

type GetRangeResponse struct {
	Rows                 []*Row
	ConsumedCapacityUnit *ConsumedCapacityUnit
	NextStartPrimaryKey  *PrimaryKey
	ResponseInfo
}

type SQLQueryRequest struct {
	Query string
}

type SearchConsumedCU struct {
	TableName            string
	IndexName            string
	ConsumedCapacityUnit *ConsumedCapacityUnit
}

type TableConsumedCU struct {
	TableName            string
	ConsumedCapacityUnit *ConsumedCapacityUnit
}

type SQLQueryConsumed struct {
	SearchConsumes []*SearchConsumedCU
	TableConsumes  []*TableConsumedCU
}

type SQLQueryResponse struct {
	ResultSet        SQLResultSet
	StmtType         SQLStatementType
	PayloadVersion   SQLPayloadVersion
	SQLQueryConsumed *SQLQueryConsumed
	ResponseInfo
}

type ListStreamRequest struct {
	TableName *string
}

type Stream struct {
	Id           *StreamId
	TableName    *string
	CreationTime int64
}

type ListStreamResponse struct {
	Streams []Stream
	ResponseInfo
}

type StreamSpecification struct {
	EnableStream       bool
	ExpirationTime     int32    // must be positive. in hours
	OriginColumnsToGet []string //origin columns to get for stream data
}

type StreamDetails struct {
	EnableStream       bool
	StreamId           *StreamId // nil when stream is disabled.
	ExpirationTime     int32     // in hours
	LastEnableTime     int64     // the last time stream is enabled, in usec
	OriginColumnsToGet []string  //origin columns to get for stream data
}

type DescribeStreamRequest struct {
	StreamId              *StreamId // required
	InclusiveStartShardId *ShardId  // optional
	ShardLimit            *int32    // optional
}

type DescribeStreamResponse struct {
	StreamId       *StreamId    // required
	ExpirationTime int32        // in hours
	TableName      *string      // required
	CreationTime   int64        // in usec
	Status         StreamStatus // required
	Shards         []*StreamShard
	NextShardId    *ShardId // optional. nil means "no more shards"
	ResponseInfo
}

type GetShardIteratorRequest struct {
	StreamId  *StreamId // required
	ShardId   *ShardId  // required
	Timestamp *int64
	Token     *string
}

type GetShardIteratorResponse struct {
	ShardIterator *ShardIterator // required
	Token         *string
	ResponseInfo
}

type GetStreamRecordRequest struct {
	ShardIterator *ShardIterator // required
	Limit         *int32         // optional. max records which will reside in response
}

type GetStreamRecordResponse struct {
	Records           []*StreamRecord
	NextShardIterator *ShardIterator // optional. an indicator to be used to read more records in this shard
	ResponseInfo
}

type ComputeSplitPointsBySizeRequest struct {
	TableName           string
	SplitSize           int64
	SplitSizeUnitInByte *int64
	SplitPointLimit     *int32
}

type ComputeSplitPointsBySizeResponse struct {
	SchemaEntry []*PrimaryKeySchema
	Splits      []*Split
	ResponseInfo
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

const (
	SS_Enabling StreamStatus = iota
	SS_Active
)

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
	Type          ActionType
	Info          *RecordSequenceInfo // required
	PrimaryKey    *PrimaryKey         // required
	Columns       []*RecordColumn
	OriginColumns []*RecordColumn
}

func (this *StreamRecord) String() string {
	return fmt.Sprintf(
		"{\"Type\":%s, \"PrimaryKey\":%s, \"Info\":%s, \"Columns\":%s, \"OriginColumns\":%s}",
		this.Type,
		*this.PrimaryKey,
		this.Info,
		this.Columns,
		this.OriginColumns)
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

const (
	IT_GLOBAL_INDEX IndexType = 0
	IT_LOCAL_INDEX  IndexType = 1
)

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

type StartLocalTransactionResponse struct {
	TransactionId *string
	ResponseInfo
}

type CommitTransactionRequest struct {
	TransactionId *string
}

type CommitTransactionResponse struct {
	ResponseInfo
}

type AbortTransactionRequest struct {
	TransactionId *string
}

type AbortTransactionResponse struct {
	ResponseInfo
}

// compute splits
type SearchIndexSplitsOptions struct {
	IndexName string
}

type ComputeSplitsRequest struct {
	TableName                string
	searchIndexSplitsOptions *SearchIndexSplitsOptions
}

type ComputeSplitsResponse struct {
	SessionId  []byte
	SplitsSize int32
	ResponseInfo
}

func (r *ComputeSplitsRequest) SetTableName(tableName string) *ComputeSplitsRequest {
	r.TableName = tableName
	return r
}

func (r *ComputeSplitsRequest) SetSearchIndexSplitsOptions(options SearchIndexSplitsOptions) *ComputeSplitsRequest {
	r.searchIndexSplitsOptions = &options
	return r
}

type TimeseriesClient struct {
	*internalClient
	timeseriesMetaCache *lruCache.Cache
}

func (timeseriesClient *TimeseriesClient) SetTimeseriesMetaCache(timeseriesMetaCache *lruCache.Cache) {
	timeseriesClient.timeseriesMetaCache = timeseriesMetaCache
}

func (timeseriesClient *TimeseriesClient) GetTimeseriesMetaCache() *lruCache.Cache {
	return timeseriesClient.timeseriesMetaCache
}

type TimeseriesTableOptions struct {
	timeToLive int64
}

func NewTimeseriesTableOptions(timeToLive int64) *TimeseriesTableOptions {
	return &TimeseriesTableOptions{
		timeToLive: timeToLive,
	}
}

func (timeseriesTableOptions *TimeseriesTableOptions) SetTimetoLive(timetoLive int64) {
	timeseriesTableOptions.timeToLive = timetoLive
}

func (timeseriesTableOptions *TimeseriesTableOptions) GetTimeToLive() int64 {
	return timeseriesTableOptions.timeToLive
}

type TimeseriesTableMeta struct {
	timeseriesTableName    string
	timeseriesTableOptions *TimeseriesTableOptions
}

func NewTimeseriesTableMeta(timeseriesTableName string) *TimeseriesTableMeta {
	return &TimeseriesTableMeta{
		timeseriesTableName: timeseriesTableName,
	}
}

func (timeseriesTableMeta *TimeseriesTableMeta) SetTimeseriesTableName(timeseriesTableName string) {
	timeseriesTableMeta.timeseriesTableName = timeseriesTableName
}

func (timeseriesTableMeta *TimeseriesTableMeta) GetTimeseriesTableName() string {
	return timeseriesTableMeta.timeseriesTableName
}

func (timeseriesTableMeta *TimeseriesTableMeta) SetTimeseriesTableOptions(timeseriesTableOptions *TimeseriesTableOptions) {
	timeseriesTableMeta.timeseriesTableOptions = timeseriesTableOptions
}

func (timeseriesTableMeta *TimeseriesTableMeta) GetTimeseriesTableOPtions() *TimeseriesTableOptions {
	return timeseriesTableMeta.timeseriesTableOptions
}

type CreateTimeseriesTableRequest struct {
	timeseriesTableMeta *TimeseriesTableMeta
}

func NewCreateTimeseriesTableRequest() *CreateTimeseriesTableRequest {
	return &CreateTimeseriesTableRequest{}
}

func (createTimeseriesTableRequest *CreateTimeseriesTableRequest) SetTimeseriesTableMeta(timeseriesTableMeta *TimeseriesTableMeta) {
	createTimeseriesTableRequest.timeseriesTableMeta = timeseriesTableMeta
}

func (createTimeseriesTableRequest *CreateTimeseriesTableRequest) GetTimeseriesTableMeta() *TimeseriesTableMeta {
	return createTimeseriesTableRequest.timeseriesTableMeta
}

type CreateTimeseriesTableResponse struct {
	ResponseInfo
}

type PutTimeseriesDataRequest struct {
	timeseriesTableName string
	rows                []*TimeseriesRow
}

func NewPutTimeseriesDataRequest(timeseriesTableName string) *PutTimeseriesDataRequest {
	return &PutTimeseriesDataRequest{
		timeseriesTableName: timeseriesTableName,
		rows:                make([]*TimeseriesRow, 0, initMapLen),
	}
}

func (putTimeseriesDataRequest *PutTimeseriesDataRequest) SetTimeseriesTableName(timeseriesTableName string) {
	putTimeseriesDataRequest.timeseriesTableName = timeseriesTableName
}

func (putTimeseriesDataRequest *PutTimeseriesDataRequest) GetTimeseriesTableName() string {
	return putTimeseriesDataRequest.timeseriesTableName
}

func (putTimeseriesDataRequest *PutTimeseriesDataRequest) AddTimeseriesRows(timeseriesRows ...*TimeseriesRow) {
	for i := 0; i < len(timeseriesRows); i++ {
		putTimeseriesDataRequest.rows = append(putTimeseriesDataRequest.rows, timeseriesRows[i])
	}
}

func (putTimeseriesDataRequest *PutTimeseriesDataRequest) GetTimeseriesRows() []*TimeseriesRow {
	if len(putTimeseriesDataRequest.rows) == 0 {
		return nil
	}
	return putTimeseriesDataRequest.rows
}

type TimeseriesRow struct {
	timeseriesKey     *TimeseriesKey
	timeInUs          int64
	fields            map[string]*ColumnValue
	timeseriesMetaKey *string
}

func NewTimeseriesRow(timeseriesKey *TimeseriesKey) *TimeseriesRow {
	return &TimeseriesRow{
		timeseriesKey: timeseriesKey,
		fields:        make(map[string]*ColumnValue, initMapLen),
	}
}

func (timeseriesRow *TimeseriesRow) SetTimeseriesKey(timeseriesKey *TimeseriesKey) {
	timeseriesRow.timeseriesKey = timeseriesKey
}

func (timeseriesRow *TimeseriesRow) GetTimeseriesKey() *TimeseriesKey {
	return timeseriesRow.timeseriesKey
}

func (timeseriesRow *TimeseriesRow) GetFieldsMap() map[string]*ColumnValue {
	return timeseriesRow.fields
}

func (timeseriesRow *TimeseriesRow) GetFieldsSlice() []string {
	n := len(timeseriesRow.GetFieldsSlice())
	key := make([]string, 0, n)
	for field_key, _ := range timeseriesRow.GetFieldsMap() {
		key = append(key, field_key)
	}
	sort.Strings(key)
	for i := 0; i < len(key); i++ {
		field_value := timeseriesRow.GetFieldsMap()[key[i]]
		switch field_value.Type {
		case ColumnType_STRING:
			key[i] = key[i] + "=" + field_value.Value.(string)
			break
		case ColumnType_INTEGER:
			key[i] = key[i] + "=" + strconv.Itoa(int(field_value.Value.(int64)))
			break
		case ColumnType_BOOLEAN:
			key[i] = key[i] + "=" + fmt.Sprintf("%v", field_value.Value.(bool))
			break
		case ColumnType_DOUBLE:
			key[i] = key[i] + "=" + fmt.Sprintf("%v", field_value.Value.(float64))
			break
		case ColumnType_BINARY:
			key[i] = key[i] + "=" + fmt.Sprintf("%v", field_value.Value.([]byte))
		default:
			panic("Unknow field Value type")
		}
	}
	return key
}

func (timeseriesRow *TimeseriesRow) AddField(fieldName string, fieldValue *ColumnValue) {
	if fieldValue == nil {
		return
	}
	fieldName = strings.ToLower(fieldName)
	timeseriesRow.fields[fieldName] = fieldValue
}

func (timeseriesRow *TimeseriesRow) AddFields(fieldsMap map[string]*ColumnValue) {
	if fieldsMap == nil {
		return
	}

	for field_key, field_value := range fieldsMap {
		if field_value == nil {
			continue
		}
		field_key = strings.ToLower(field_key)
		timeseriesRow.fields[field_key] = field_value
	}
}

func (timeseriesRow *TimeseriesRow) SetTimeInus(timestamp int64) {
	timeseriesRow.timeInUs = timestamp
}

func (timeseriesRow *TimeseriesRow) GetTimeInus() int64 {
	return timeseriesRow.timeInUs
}

type TimeseriesKey struct {
	measurement string
	source      string
	tags        map[string]string
	tagsString  *string
}

func NewTimeseriesKey() *TimeseriesKey {
	return &TimeseriesKey{
		tags: make(map[string]string, initMapLen),
	}
}

func (timeseriesKey *TimeseriesKey) AddTag(tagName string, tagValue string) {
	timeseriesKey.tags[tagName] = tagValue
	timeseriesKey.tagsString = nil
}

func (timeseriesKey *TimeseriesKey) AddTags(tagsMap map[string]string) {
	if tagsMap == nil {
		return
	}
	for tagName, tagValue := range tagsMap {
		timeseriesKey.tags[tagName] = tagValue
	}
	timeseriesKey.tagsString = nil
}

func (timeseriesKey *TimeseriesKey) GetTags() map[string]string {
	return timeseriesKey.tags
}

func (timeseriesKey *TimeseriesKey) buildTimeseriesMetaKey(timeseriesTableName string) (string, error) {
	var capacity int
	var err error
	capacity += len(timeseriesTableName)
	capacity += len(timeseriesKey.measurement)
	capacity += len(timeseriesKey.source)
	if timeseriesKey.tagsString == nil {
		timeseriesKey.tagsString = new(string)
		if *timeseriesKey.tagsString, err = BuildTagString(timeseriesKey.tags); err != nil {
			return "", err
		}
	}
	capacity += len(*timeseriesKey.tagsString)
	capacity += 3

	sb := strings.Builder{}
	sb.Grow(capacity)

	sb.WriteString(timeseriesTableName)
	sb.WriteString("\t")
	sb.WriteString(timeseriesKey.measurement)
	sb.WriteString("\t")
	sb.WriteString(timeseriesKey.source)
	sb.WriteString("\t")
	sb.WriteString(*timeseriesKey.tagsString)

	return sb.String(), nil
}

func (timeseriesKey *TimeseriesKey) SetMeasurementName(measurementName string) {
	timeseriesKey.measurement = measurementName
}

func (timeseriesKey *TimeseriesKey) GetMeasurementName() string {
	return timeseriesKey.measurement
}

func (timeseriesKey *TimeseriesKey) SetDataSource(source string) {
	timeseriesKey.source = source
}

func (timeseriesKey *TimeseriesKey) GetDataSource() string {
	return timeseriesKey.source
}

type PutTimeseriesDataResponse struct {
	failedRowResults []*FailedRowResult
	ResponseInfo
}

type FailedRowResult struct {
	Index     int32
	Error     error
	ErrorCode string
}

func (putTimeseriesDataResponse *PutTimeseriesDataResponse) GetFailedRowResults() []*FailedRowResult {
	return putTimeseriesDataResponse.failedRowResults
}

type GetTimeseriesDataRequest struct {
	timeseriesTableName string
	timeseriesKey       *TimeseriesKey
	beginTimeInUs       int64
	endTimeInUs         int64
	nextToken           []byte
	limit               int32
}

func NewGetTimeseriesDataRequest(timeseriesTableName string) *GetTimeseriesDataRequest {
	return &GetTimeseriesDataRequest{
		timeseriesTableName: timeseriesTableName,
		limit:               -1,
	}
}

func (getDataRequest *GetTimeseriesDataRequest) SetTimeseriesTableName(timeseriesTableName string) {
	getDataRequest.timeseriesTableName = timeseriesTableName
}

func (getDataRequest *GetTimeseriesDataRequest) GetTimeseriesTableName() string {
	return getDataRequest.timeseriesTableName
}

func (getDataRequest *GetTimeseriesDataRequest) SetTimeseriesKey(timeseriesKey *TimeseriesKey) {
	getDataRequest.timeseriesKey = timeseriesKey
}

func (getDataRequest *GetTimeseriesDataRequest) GetTimeseriesKey() *TimeseriesKey {
	return getDataRequest.timeseriesKey
}

func (getDataRequest *GetTimeseriesDataRequest) SetTimeRange(beginTimeInUs int64, endTimeInUs int64) {
	getDataRequest.beginTimeInUs = beginTimeInUs
	getDataRequest.endTimeInUs = endTimeInUs
}

func (getDataRequest *GetTimeseriesDataRequest) GetTimeRange() (int64, int64) {
	return getDataRequest.beginTimeInUs, getDataRequest.endTimeInUs
}

func (getDataRequest *GetTimeseriesDataRequest) GetBeginTimeInUs() int64 {
	return getDataRequest.beginTimeInUs
}

func (getDataRequest *GetTimeseriesDataRequest) GetEndTimeInUs() int64 {
	return getDataRequest.endTimeInUs
}

func (getDataRequest *GetTimeseriesDataRequest) SetNextToken(nextToken []byte) {
	getDataRequest.nextToken = append([]byte{}, nextToken...)
}

func (getDataRequest *GetTimeseriesDataRequest) GetNextToken() []byte {
	return append([]byte{}, getDataRequest.nextToken...)
}

func (getDataRequest *GetTimeseriesDataRequest) SetLimit(limit int32) {
	getDataRequest.limit = limit
}

func (getDataRequest *GetTimeseriesDataRequest) GetLimit() int32 {
	return getDataRequest.limit
}

type GetTimeseriesDataResponse struct {
	rows      []*TimeseriesRow
	nextToken []byte
	ResponseInfo
}

func (getTimeseriesDataResp *GetTimeseriesDataResponse) GetRows() []*TimeseriesRow {
	return getTimeseriesDataResp.rows
}

func (getTimeseriesDataResp *GetTimeseriesDataResponse) GetNextToken() []byte {
	return getTimeseriesDataResp.nextToken
}

type DescribeTimeseriesTableRequest struct {
	timeseriesTableName string
}

func NewDescribeTimeseriesTableRequset(timeseriesTableName string) *DescribeTimeseriesTableRequest {
	requset := &DescribeTimeseriesTableRequest{}
	requset.SetTimeseriesTableName(timeseriesTableName)
	return requset
}

func (describeTimeseriesReq *DescribeTimeseriesTableRequest) SetTimeseriesTableName(timeseriesTableName string) {
	describeTimeseriesReq.timeseriesTableName = timeseriesTableName
}

func (describeTimeseriesReq *DescribeTimeseriesTableRequest) GetTimeseriesTableName() string {
	return describeTimeseriesReq.timeseriesTableName
}

type DescribeTimeseriesTableResponse struct {
	timeseriesTableMeta *TimeseriesTableMeta
	ResponseInfo
}

func (describeTimeseriesTableResp *DescribeTimeseriesTableResponse) GetTimeseriesTableMeta() *TimeseriesTableMeta {
	return describeTimeseriesTableResp.timeseriesTableMeta
}

type ListTimeseriesTableRequest struct {
}

func NewListTimeseriesTableRequest() *ListTimeseriesTableRequest {
	return &ListTimeseriesTableRequest{}
}

type ListTimeseriesTableResponse struct {
	timeseriesTableMetas []*TimeseriesTableMeta
	ResponseInfo
}

func (listTimeseriesTableResponse *ListTimeseriesTableResponse) GetTimeseriesTableMeta() []*TimeseriesTableMeta {
	return listTimeseriesTableResponse.timeseriesTableMetas
}

func (listTimeseriesTableResponse *ListTimeseriesTableResponse) GetTimeseriesTableNames() []string {
	timeseriesTableNames := []string{}
	for i := 0; i < len(listTimeseriesTableResponse.timeseriesTableMetas); i++ {
		timeseriesTableNames = append(timeseriesTableNames, listTimeseriesTableResponse.timeseriesTableMetas[i].GetTimeseriesTableName())
	}
	return timeseriesTableNames
}

type DeleteTimeseriesTableRequest struct {
	timeseriesTableName string
}

func NewDeleteTimeseriesTableRequest(timeseriesTableName string) *DeleteTimeseriesTableRequest {
	return &DeleteTimeseriesTableRequest{
		timeseriesTableName: timeseriesTableName,
	}
}

func (deleteTimeseriesTableRequest *DeleteTimeseriesTableRequest) SetTimeseriesTableName(timeseriesTableName string) {
	deleteTimeseriesTableRequest.timeseriesTableName = timeseriesTableName
}

func (deleteTimeseriesTableRequest *DeleteTimeseriesTableRequest) GetTimeseriesTableName() string {
	return deleteTimeseriesTableRequest.timeseriesTableName
}

type DeleteTimeseriesTableResponse struct {
	ResponseInfo
}

type UpdateTimeseriesMetaRequest struct {
	timeseriesTableName string
	metas               []*TimeseriesMeta
}

func NewUpdateTimeseriesMetaRequest(timeseriesTableName string) *UpdateTimeseriesMetaRequest {
	return &UpdateTimeseriesMetaRequest{
		timeseriesTableName: timeseriesTableName,
	}
}

func (updateTimeseriesMetaRequest *UpdateTimeseriesMetaRequest) SetTimeseriesTableName(timeseriesTableName string) {
	updateTimeseriesMetaRequest.timeseriesTableName = timeseriesTableName
}

func (updateTimeseriesMetaRequest *UpdateTimeseriesMetaRequest) GetTimeseriesTableName() string {
	return updateTimeseriesMetaRequest.timeseriesTableName
}

func (updateTimeseriesMetaRequest *UpdateTimeseriesMetaRequest) AddTimeseriesMetas(metas ...*TimeseriesMeta) {
	updateTimeseriesMetaRequest.metas = append(updateTimeseriesMetaRequest.metas, metas...)
}

func (updateTimeseriesMetaRequest *UpdateTimeseriesMetaRequest) GetTimeseriesMetas() []*TimeseriesMeta {
	return updateTimeseriesMetaRequest.metas
}

type UpdateTimeseriesMetaResponse struct {
	failedRowResults []*FailedRowResult
	ResponseInfo
}

func (updateTimeseriesMetaResponse *UpdateTimeseriesMetaResponse) GetFailedRowResults() []*FailedRowResult {
	return updateTimeseriesMetaResponse.failedRowResults
}

type DeleteTimeseriesMetaRequest struct {
	timeseriesTableName string
	keys                []*TimeseriesKey
}

func NewDeleteTimeseriesMetaRequest(timeseriesTableName string) *DeleteTimeseriesMetaRequest {
	return &DeleteTimeseriesMetaRequest{
		timeseriesTableName: timeseriesTableName,
	}
}

func (deleteTimeseriesMetaRequest *DeleteTimeseriesMetaRequest) SetTimeseriesTableName(timeseriesTableName string) {
	deleteTimeseriesMetaRequest.timeseriesTableName = timeseriesTableName
}

func (deleteTimeseriesMetaRequest *DeleteTimeseriesMetaRequest) GetTimeseriesTableName() string {
	return deleteTimeseriesMetaRequest.timeseriesTableName
}

func (deleteTimeseriesMetaRequest *DeleteTimeseriesMetaRequest) AddTimeseriesKeys(keys ...*TimeseriesKey) {
	deleteTimeseriesMetaRequest.keys = append(deleteTimeseriesMetaRequest.keys, keys...)
}

func (deleteTimeseriesMetaRequest *DeleteTimeseriesMetaRequest) GetTimeseriesKeys() []*TimeseriesKey {
	return deleteTimeseriesMetaRequest.keys
}

type DeleteTimeseriesMetaResponse struct {
	failedRowResults []*FailedRowResult
	ResponseInfo
}

func (deleteTimeseriesMetaResponse *DeleteTimeseriesMetaResponse) GetFailedRowResults() []*FailedRowResult {
	return deleteTimeseriesMetaResponse.failedRowResults
}

type QueryTimeseriesMetaRequest struct {
	timeseriesTableName string
	condition           MetaQueryCondition
	getTotalHits        bool
	nextToken           []byte
	limit               int32
}

func NewQueryTimeseriesMetaRequest(timeseriesTableName string) *QueryTimeseriesMetaRequest {
	return &QueryTimeseriesMetaRequest{
		timeseriesTableName: timeseriesTableName,
		limit:               -1,
	}
}

func (queryTimeseriesMetaRequest *QueryTimeseriesMetaRequest) SetTimeseriesTableName(timeseriesTableName string) {
	queryTimeseriesMetaRequest.timeseriesTableName = timeseriesTableName
}

func (queryTimeseriesMetaRequest *QueryTimeseriesMetaRequest) GetTimeseriesTableName() string {
	return queryTimeseriesMetaRequest.timeseriesTableName
}

func (queryTimeseriesMetaRequest *QueryTimeseriesMetaRequest) SetCondition(condition MetaQueryCondition) {
	queryTimeseriesMetaRequest.condition = condition
}

func (queryTimeseriesMetaRequest *QueryTimeseriesMetaRequest) GetCondition() MetaQueryCondition {
	return queryTimeseriesMetaRequest.condition
}

func (queryTimeseriesMetaRequest *QueryTimeseriesMetaRequest) SetTotalHits(getTotalHits bool) {
	queryTimeseriesMetaRequest.getTotalHits = getTotalHits
}

func (queryTimeseriesMetaRequest *QueryTimeseriesMetaRequest) GetTotalHits() bool {
	return queryTimeseriesMetaRequest.getTotalHits
}

func (queryTimeseriesMetaRequest *QueryTimeseriesMetaRequest) SetNextToken(nextToken []byte) {
	queryTimeseriesMetaRequest.nextToken = append([]byte{}, nextToken...)
}

func (queryTimeseriesMetaRequest *QueryTimeseriesMetaRequest) GetNextToken() []byte {
	return append([]byte{}, queryTimeseriesMetaRequest.nextToken...)
}

func (queryTimeseriesMetaRequest *QueryTimeseriesMetaRequest) SetLimit(limit int32) {
	queryTimeseriesMetaRequest.limit = limit
}

func (queryTimeseriesMetaRequest *QueryTimeseriesMetaRequest) GetLimit() int32 {
	return queryTimeseriesMetaRequest.limit
}

type MetaQueryConditionType int32

const (
	COMPOSITE_CONDITION   MetaQueryConditionType = 1
	MEASUREMENT_CONDITION MetaQueryConditionType = 2
	SOURCE_CONDITION      MetaQueryConditionType = 3
	TAG_CONDITION         MetaQueryConditionType = 4
	UPDATE_TIME_CONDITION MetaQueryConditionType = 5
	ATTRIBUTE_CONDITION   MetaQueryConditionType = 6
)

func (condType MetaQueryConditionType) String() string {
	switch condType {
	case COMPOSITE_CONDITION:
		return "COMPOSITE"
	case MEASUREMENT_CONDITION:
		return "MEASUREMENT"
	case SOURCE_CONDITION:
		return "SOURCE"
	case TAG_CONDITION:
		return "TAG"
	case UPDATE_TIME_CONDITION:
		return "UPDATE_TIME"
	case ATTRIBUTE_CONDITION:
		return "ATTRIBUTE"
	default:
		return string(condType)
	}
}

func ToMetaQueryConditionType(condType string) (MetaQueryConditionType, error) {
	switch strings.ToUpper(condType) {
	case "COMPOSITE":
		return COMPOSITE_CONDITION, nil
	case "MEASUREMENT":
		return MEASUREMENT_CONDITION, nil
	case "SOURCE":
		return SOURCE_CONDITION, nil
	case "TAG":
		return TAG_CONDITION, nil
	case "UPDATE_TIME":
		return UPDATE_TIME_CONDITION, nil
	case "ATTRIBUTE":
		return ATTRIBUTE_CONDITION, nil
	default:
		return COMPOSITE_CONDITION, errors.New("Invalid condition type: " + condType)
	}
}

func (op *MetaQueryConditionType) UnmarshalJSON(data []byte) (err error) {
	var opStr string
	err = json.Unmarshal(data, &opStr)
	if err != nil {
		return
	}

	*op, err = ToMetaQueryConditionType(opStr)
	if err != nil {
		return err
	}
	return
}

func (op *MetaQueryConditionType) MarshalJSON() (data []byte, err error) {
	data, err = json.Marshal(op.String())
	return
}

type MetaQuerySingleOperator int32

const (
	OP_EQUAL         MetaQuerySingleOperator = 1
	OP_GREATER_THAN  MetaQuerySingleOperator = 2
	OP_GREATER_EQUAL MetaQuerySingleOperator = 3
	OP_LESS_THAN     MetaQuerySingleOperator = 4
	OP_LESS_EQUAL    MetaQuerySingleOperator = 5
	OP_PREFIX        MetaQuerySingleOperator = 6
)

func (op MetaQuerySingleOperator) String() string {
	switch op {
	case OP_EQUAL:
		return "EQUAL"
	case OP_GREATER_THAN:
		return "GREATER_THAN"
	case OP_GREATER_EQUAL:
		return "GREATER_EQUAL"
	case OP_LESS_THAN:
		return "LESS_THAN"
	case OP_LESS_EQUAL:
		return "LESS_EQUAL"
	case OP_PREFIX:
		return "PREFIX"
	default:
		return string(op)
	}
}

func ToMetaQuerySingleOperator(op string) (MetaQuerySingleOperator, error) {
	switch strings.ToUpper(op) {
	case "EQUAL":
		return OP_EQUAL, nil
	case "GREATER_THAN":
		return OP_GREATER_THAN, nil
	case "GREATER_EQUAL":
		return OP_GREATER_EQUAL, nil
	case "LESS_THAN":
		return OP_LESS_THAN, nil
	case "LESS_EQUAL":
		return OP_LESS_EQUAL, nil
	case "PREFIX":
		return OP_PREFIX, nil
	default:
		return OP_EQUAL, errors.New("Invalid operator: " + op)
	}
}

func (op *MetaQuerySingleOperator) UnmarshalJSON(data []byte) (err error) {
	var opStr string
	err = json.Unmarshal(data, &opStr)
	if err != nil {
		return
	}

	*op, err = ToMetaQuerySingleOperator(opStr)
	if err != nil {
		return err
	}
	return
}

func (op *MetaQuerySingleOperator) MarshalJSON() (data []byte, err error) {
	data, err = json.Marshal(op.String())
	return
}

type MetaQueryCompositeOperator int32

const (
	OP_AND MetaQueryCompositeOperator = 1
	OP_OR  MetaQueryCompositeOperator = 2
	OP_NOT MetaQueryCompositeOperator = 3
)

func (op MetaQueryCompositeOperator) String() string {
	switch op {
	case OP_AND:
		return "AND"
	case OP_OR:
		return "OR"
	case OP_NOT:
		return "NOT"
	default:
		return string(op)
	}
}

func ToMetaQueryCompositeOperator(op string) (MetaQueryCompositeOperator, error) {
	switch strings.ToUpper(op) {
	case "AND":
		return OP_AND, nil
	case "OR":
		return OP_OR, nil
	case "NOT":
		return OP_NOT, nil
	default:
		return OP_AND, errors.New("Invalid operator: " + op)
	}
}

func (op *MetaQueryCompositeOperator) UnmarshalJSON(data []byte) (err error) {
	var opStr string
	err = json.Unmarshal(data, &opStr)
	if err != nil {
		return
	}

	*op, err = ToMetaQueryCompositeOperator(opStr)
	if err != nil {
		return err
	}
	return
}

func (op *MetaQueryCompositeOperator) MarshalJSON() (data []byte, err error) {
	data, err = json.Marshal(op.String())
	return
}

type MetaQueryCondition interface {
	GetType() MetaQueryConditionType
	Serialize() []byte
}

type MetaQueryConditionWrapper struct {
	Type           MetaQueryConditionType
	QueryCondition MetaQueryCondition
}

func (op *MetaQueryConditionWrapper) UnmarshalJSON(data []byte) (err error) {
	rawData := make(map[string]json.RawMessage)
	err = json.Unmarshal(data, &rawData)
	if err != nil {
		return
	}

	var condTypeStr string
	condTypeRD, ok := rawData["Type"]
	if !ok {
		err = errors.New("type information is missing")
		return
	}

	err = json.Unmarshal(condTypeRD, &condTypeStr)
	if err != nil {
		return
	}

	condType, err := ToMetaQueryConditionType(condTypeStr)
	if err != nil {
		return
	}

	condRM, ok := rawData["QueryCondition"]
	if !ok {
		err = errors.New("query condition is missing")
		return
	}

	op.Type = condType
	switch condType {
	case COMPOSITE_CONDITION:
		rc := CompositeMetaQueryCondition{}
		err = json.Unmarshal(condRM, &rc)
		op.QueryCondition = &rc
	case MEASUREMENT_CONDITION:
		rc := MeasurementMetaQueryCondition{}
		err = json.Unmarshal(condRM, &rc)
		op.QueryCondition = &rc
	case SOURCE_CONDITION:
		rc := DataSourceMetaQueryCondition{}
		err = json.Unmarshal(condRM, &rc)
		op.QueryCondition = &rc
	case TAG_CONDITION:
		rc := TagMetaQueryCondition{}
		err = json.Unmarshal(condRM, &rc)
		op.QueryCondition = &rc
	case UPDATE_TIME_CONDITION:
		rc := UpdateTimeMetaQueryCondition{}
		err = json.Unmarshal(condRM, &rc)
		op.QueryCondition = &rc
	case ATTRIBUTE_CONDITION:
		rc := AttributeMetaQueryCondition{}
		err = json.Unmarshal(condRM, &rc)
		op.QueryCondition = &rc
	}

	return
}

type MeasurementMetaQueryCondition struct {
	Operator MetaQuerySingleOperator
	Value    string
}

func NewMeasurementQueryCondition(operator MetaQuerySingleOperator, value string) *MeasurementMetaQueryCondition {
	return &MeasurementMetaQueryCondition{
		Operator: operator,
		Value:    value,
	}
}

func (measurementMetaQueryCondition *MeasurementMetaQueryCondition) GetType() MetaQueryConditionType {
	return MEASUREMENT_CONDITION
}

func (measurementMetaQueryCondition *MeasurementMetaQueryCondition) Serialize() []byte {
	metaQueryMeasurementCondition := new(otsprotocol.MetaQueryMeasurementCondition)
	metaQueryMeasurementCondition.Op = otsprotocol.MetaQuerySingleOperator(int32(measurementMetaQueryCondition.Operator)).Enum()
	metaQueryMeasurementCondition.Value = proto.String(measurementMetaQueryCondition.Value)

	result, err := proto.Marshal(metaQueryMeasurementCondition)
	if err != nil {
		panic(fmt.Errorf("MeasurementMetaQueryCondition serialize failed with err : %s", err))
		return nil
	}

	return result
}

type DataSourceMetaQueryCondition struct {
	Operator MetaQuerySingleOperator
	Value    string
}

func NewDataSourceMetaQueryCondition(operator MetaQuerySingleOperator, value string) *DataSourceMetaQueryCondition {
	return &DataSourceMetaQueryCondition{
		Operator: operator,
		Value:    value,
	}
}

func (sourceMetaQueryCondition *DataSourceMetaQueryCondition) GetType() MetaQueryConditionType {
	return SOURCE_CONDITION
}

func (sourceMetaQueryCondition *DataSourceMetaQueryCondition) Serialize() []byte {
	metaQuerySourceCondition := new(otsprotocol.MetaQuerySourceCondition)
	metaQuerySourceCondition.Op = otsprotocol.MetaQuerySingleOperator(int32(sourceMetaQueryCondition.Operator)).Enum()
	metaQuerySourceCondition.Value = proto.String(sourceMetaQueryCondition.Value)

	result, err := proto.Marshal(metaQuerySourceCondition)
	if err != nil {
		panic(fmt.Errorf("SourceMetaQueryCondition serialize failed with err : %s", err))
		return nil
	}
	return result
}

type TagMetaQueryCondition struct {
	Operator MetaQuerySingleOperator
	TagName  string
	Value    string
}

func NewTagMetaQueryCondition(operator MetaQuerySingleOperator, tagName string, value string) *TagMetaQueryCondition {
	return &TagMetaQueryCondition{
		Operator: operator,
		TagName:  tagName,
		Value:    value,
	}
}

func (TagMetaQueryCondition *TagMetaQueryCondition) GetType() MetaQueryConditionType {
	return TAG_CONDITION
}

func (tagMetaQueryCondition *TagMetaQueryCondition) Serialize() []byte {
	metaQueryTagCondition := new(otsprotocol.MetaQueryTagCondition)
	metaQueryTagCondition.Op = otsprotocol.MetaQuerySingleOperator(int32(tagMetaQueryCondition.Operator)).Enum()
	metaQueryTagCondition.TagName = proto.String(tagMetaQueryCondition.TagName)
	metaQueryTagCondition.Value = proto.String(tagMetaQueryCondition.Value)

	result, err := proto.Marshal(metaQueryTagCondition)
	if err != nil {
		panic(fmt.Errorf("TagMetaQueryCondition serialize failed with err : %s", err))
		return nil
	}
	return result
}

type UpdateTimeMetaQueryCondition struct {
	Operator MetaQuerySingleOperator
	TimeInUs int64
}

func NewUpdateTimeMetaQueryCondition(operator MetaQuerySingleOperator, timeInUs int64) *UpdateTimeMetaQueryCondition {
	return &UpdateTimeMetaQueryCondition{
		Operator: operator,
		TimeInUs: timeInUs,
	}
}

func (updateTimeMetaQueryCondition *UpdateTimeMetaQueryCondition) GetType() MetaQueryConditionType {
	return UPDATE_TIME_CONDITION
}

func (updateTimeMetaQueryCondition *UpdateTimeMetaQueryCondition) Serialize() []byte {
	metaQueryUpdateTimeCondition := new(otsprotocol.MetaQueryUpdateTimeCondition)
	metaQueryUpdateTimeCondition.Op = otsprotocol.MetaQuerySingleOperator(int32(updateTimeMetaQueryCondition.Operator)).Enum()
	metaQueryUpdateTimeCondition.Value = proto.Int64(updateTimeMetaQueryCondition.TimeInUs)

	result, err := proto.Marshal(metaQueryUpdateTimeCondition)
	if err != nil {
		panic(fmt.Errorf("UpdateTimeMetaQueryCondition serialize failed with err : %s", err))
		return nil
	}
	return result
}

type AttributeMetaQueryCondition struct {
	Operator      MetaQuerySingleOperator
	AttributeName string
	Value         string
}

func NewAttributeMetaQueryCondition(operator MetaQuerySingleOperator, attributeName string, value string) *AttributeMetaQueryCondition {
	return &AttributeMetaQueryCondition{
		Operator:      operator,
		AttributeName: attributeName,
		Value:         value,
	}
}

func (attributeMetaQueryCondition *AttributeMetaQueryCondition) GetType() MetaQueryConditionType {
	return ATTRIBUTE_CONDITION
}

func (attributeMetaQueryCondition *AttributeMetaQueryCondition) Serialize() []byte {
	metaQueryAttributeCondition := new(otsprotocol.MetaQueryAttributeCondition)
	metaQueryAttributeCondition.Op = otsprotocol.MetaQuerySingleOperator(int32(attributeMetaQueryCondition.Operator)).Enum()
	metaQueryAttributeCondition.AttrName = proto.String(attributeMetaQueryCondition.AttributeName)
	metaQueryAttributeCondition.Value = proto.String(attributeMetaQueryCondition.Value)

	result, err := proto.Marshal(metaQueryAttributeCondition)
	if err != nil {
		panic(fmt.Errorf("AttributeMetaQueryCondition serialize failed with err : %s", err))
		return nil
	}
	return result
}

type CompositeMetaQueryCondition struct {
	Operator      MetaQueryCompositeOperator
	SubConditions []*MetaQueryCondition `json:"-"`

	// for json marshal and unmarshal
	SubConditionsAlias []*MetaQueryConditionWrapper `json:"SubConditions"`
}

func (op *CompositeMetaQueryCondition) UnmarshalJSON(data []byte) (err error) {
	type CompositeMetaQueryConditionAlias CompositeMetaQueryCondition
	condAlias := CompositeMetaQueryConditionAlias{}
	err = json.Unmarshal(data, &condAlias)
	if err != nil {
		return
	}

	op.Operator = condAlias.Operator
	op.SubConditions = make([]*MetaQueryCondition, 0)
	for _, cond := range condAlias.SubConditionsAlias {
		op.SubConditions = append(op.SubConditions, &cond.QueryCondition)
	}

	return
}

func (op *CompositeMetaQueryCondition) MarshalJSON() (data []byte, err error) {
	type CompositeMetaQueryConditionAlias CompositeMetaQueryCondition
	condAlias := CompositeMetaQueryConditionAlias(*op)
	condAlias.Operator = op.Operator
	condAlias.SubConditionsAlias = make([]*MetaQueryConditionWrapper, 0)
	for _, cond := range condAlias.SubConditions {
		condAlias.SubConditionsAlias = append(condAlias.SubConditionsAlias,
			&MetaQueryConditionWrapper{
				Type:           (*cond).GetType(),
				QueryCondition: *cond,
			})
	}

	data, err = json.Marshal(&condAlias)
	return
}

func NewCompositeMetaQueryCondition(operator MetaQueryCompositeOperator, subConditions ...MetaQueryCondition) *CompositeMetaQueryCondition {
	compositeMetaQueryCondition := &CompositeMetaQueryCondition{
		Operator: operator,
	}
	if len(subConditions) > 0 {
		for i := 0; i < len(subConditions); i++ {
			compositeMetaQueryCondition.SubConditions = append(compositeMetaQueryCondition.SubConditions, &subConditions[i])
		}

	}
	return compositeMetaQueryCondition
}

func (compositeMetaQueryCondition *CompositeMetaQueryCondition) GetType() MetaQueryConditionType {
	return COMPOSITE_CONDITION
}

func (compositemetaQueryCondition *CompositeMetaQueryCondition) Serialize() []byte {
	metaQueryCompositeCondition := new(otsprotocol.MetaQueryCompositeCondition)
	metaQueryCompositeCondition.Op = otsprotocol.MetaQueryCompositeOperator(int32(compositemetaQueryCondition.Operator)).Enum()
	for i := 0; i < len(compositemetaQueryCondition.SubConditions); i++ {
		metaQueryCondition := new(otsprotocol.MetaQueryCondition)
		switch value := (*compositemetaQueryCondition.getSubConditions()[i]).(type) {
		case *MeasurementMetaQueryCondition:
			metaQueryCondition.Type = otsprotocol.MetaQueryConditionType_MEASUREMENT_CONDITION.Enum()
			metaQueryCondition.ProtoData = value.Serialize()
			break
		case *DataSourceMetaQueryCondition:
			metaQueryCondition.Type = otsprotocol.MetaQueryConditionType_SOURCE_CONDITION.Enum()
			metaQueryCondition.ProtoData = value.Serialize()
			break
		case *TagMetaQueryCondition:
			metaQueryCondition.Type = otsprotocol.MetaQueryConditionType_TAG_CONDITION.Enum()
			metaQueryCondition.ProtoData = value.Serialize()
			break
		case *UpdateTimeMetaQueryCondition:
			metaQueryCondition.Type = otsprotocol.MetaQueryConditionType_UPDATE_TIME_CONDITION.Enum()
			metaQueryCondition.ProtoData = value.Serialize()
			break
		case *CompositeMetaQueryCondition:
			metaQueryCondition.Type = otsprotocol.MetaQueryConditionType_COMPOSITE_CONDITION.Enum()
			metaQueryCondition.ProtoData = value.Serialize()
			break
		default:
			panic("Unknow singleMetaQueryConditionType in compositeMetaQueryCondition!")
			return nil
		}
		metaQueryCompositeCondition.SubConditions = append(metaQueryCompositeCondition.SubConditions, metaQueryCondition)
	}

	result, err := proto.Marshal(metaQueryCompositeCondition)
	if err != nil {
		panic(fmt.Errorf("otsprotocol.MetaQueryCompositeCondition Serialize Failed with err : %s", err))
		return nil
	}
	return result
}

func (compositeMetaQueryCondition *CompositeMetaQueryCondition) AddSubConditions(subconditions ...MetaQueryCondition) {
	for i := 0; i < len(subconditions); i++ {
		compositeMetaQueryCondition.SubConditions = append(compositeMetaQueryCondition.SubConditions, &subconditions[i])
	}

}

func (compositeMetaQueryCondition *CompositeMetaQueryCondition) getSubConditions() []*MetaQueryCondition {
	return compositeMetaQueryCondition.SubConditions
}

func (compositeMetaQueryCondition *CompositeMetaQueryCondition) SetOperator(operator MetaQueryCompositeOperator) {
	compositeMetaQueryCondition.Operator = operator
}

func (compositeMetaQueryCondition *CompositeMetaQueryCondition) GetOperator() MetaQueryCompositeOperator {
	return compositeMetaQueryCondition.Operator
}

type TimeseriesMeta struct {
	timeseriesKey  *TimeseriesKey
	attributes     map[string]string
	updateTimeInUs int64
}

func NewTimeseriesMeta(timeseriesKey *TimeseriesKey) *TimeseriesMeta {
	return &TimeseriesMeta{
		timeseriesKey: timeseriesKey,
		attributes:    map[string]string{},
	}
}

func (timeseriesMeta *TimeseriesMeta) SetTimeseriesKey(timeseriesKey *TimeseriesKey) {
	timeseriesMeta.timeseriesKey = timeseriesKey
}

func (timeseriesMeta *TimeseriesMeta) GetTimeseriesKey() *TimeseriesKey {
	return timeseriesMeta.timeseriesKey
}

func (timeseriesMeta *TimeseriesMeta) AddAttribute(attr_key string, attr_value string) {
	timeseriesMeta.attributes[attr_key] = attr_value
}

func (timeseriesMeta *TimeseriesMeta) AddAttributes(attributes map[string]string) {
	for key, value := range attributes {
		timeseriesMeta.attributes[key] = value
	}
}

func (timeseriesMeta *TimeseriesMeta) GetAttributes() map[string]string {
	attributes := map[string]string{}
	for key, value := range timeseriesMeta.attributes {
		attributes[key] = value
	}
	return attributes
}

func (timeseriesMeta *TimeseriesMeta) GetAttributeSlice() string {
	if len(timeseriesMeta.attributes) > 0 {
		Atrributes, _ := BuildTagString(timeseriesMeta.attributes)
		return Atrributes
	}
	return ""
}

func (timeseriesMeta *TimeseriesMeta) SetUpdateTimeInUs(updateTimeInUs int64) {
	timeseriesMeta.updateTimeInUs = updateTimeInUs
}

func (timeseriesMeta *TimeseriesMeta) GetUpdateTimeInUs() int64 {
	return timeseriesMeta.updateTimeInUs
}

type QueryTimeseriesMetaResponse struct {
	timeseriesMetas []*TimeseriesMeta
	totalHits       int64
	nextToken       []byte
	ResponseInfo
}

func newQueryTimeseriesMetaResponse() *QueryTimeseriesMetaResponse {
	return &QueryTimeseriesMetaResponse{
		totalHits: -1,
	}
}

func (queryTimeseriesMetaResponse *QueryTimeseriesMetaResponse) GetTimeseriesMetas() []*TimeseriesMeta {
	return queryTimeseriesMetaResponse.timeseriesMetas
}

func (queryTimeseriesMetaResponse *QueryTimeseriesMetaResponse) GetTotalHits() int64 {
	return queryTimeseriesMetaResponse.totalHits
}

func (queryTimeseriesMetaResponse *QueryTimeseriesMetaResponse) GetNextToken() []byte {
	return queryTimeseriesMetaResponse.nextToken
}

// UpdateTimeseriesTableRequest
type UpdateTimeseriesTableRequest struct {
	timeseriesTableName    string
	timeseriesTableOptions *TimeseriesTableOptions
}

func NewUpdateTimeseriesTableRequest(timeseriesTableName string) *UpdateTimeseriesTableRequest {
	return &UpdateTimeseriesTableRequest{
		timeseriesTableName: timeseriesTableName,
	}
}

func (updateTimeseriesTableReq *UpdateTimeseriesTableRequest) SetTimeseriesTalbeName(timeseriesTableName string) {
	updateTimeseriesTableReq.timeseriesTableName = timeseriesTableName
}

func (updateTimeseriesTableReq *UpdateTimeseriesTableRequest) GetTimeseriesTableName() string {
	return updateTimeseriesTableReq.timeseriesTableName
}

func (updateTimeseriesTableReq *UpdateTimeseriesTableRequest) SetTimeseriesTableOptions(timeseriesTableOptions *TimeseriesTableOptions) {
	updateTimeseriesTableReq.timeseriesTableOptions = timeseriesTableOptions
}

func (updateTimeseriesTableReq *UpdateTimeseriesTableRequest) GetTimeseriesTableOptions() *TimeseriesTableOptions {
	return updateTimeseriesTableReq.timeseriesTableOptions
}

type UpdateTimeseriesTableResponse struct {
	ResponseInfo
}

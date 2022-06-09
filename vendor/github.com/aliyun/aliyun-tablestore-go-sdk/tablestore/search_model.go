package tablestore

import (
	"encoding/json"
	"errors"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/search"
	"github.com/golang/protobuf/proto"
	"strings"
)

type ColumnsToGet struct {
	Columns            []string
	ReturnAll          bool
	ReturnAllFromIndex bool
}

type SearchRequest struct {
	TableName     string
	IndexName     string
	SearchQuery   search.SearchQuery
	ColumnsToGet  *ColumnsToGet
	RoutingValues []*PrimaryKey
	TimeoutMs     *int32
}

func (r *SearchRequest) SetTableName(tableName string) *SearchRequest {
	r.TableName = tableName
	return r
}

func (r *SearchRequest) SetIndexName(indexName string) *SearchRequest {
	r.IndexName = indexName
	return r
}

func (r *SearchRequest) SetSearchQuery(searchQuery search.SearchQuery) *SearchRequest {
	r.SearchQuery = searchQuery
	return r
}

func (r *SearchRequest) SetColumnsToGet(columnToGet *ColumnsToGet) *SearchRequest {
	r.ColumnsToGet = columnToGet
	return r
}

func (r *SearchRequest) SetRoutingValues(routingValues []*PrimaryKey) *SearchRequest {
	r.RoutingValues = routingValues
	return r
}

func (r *SearchRequest) AddRoutingValue(routingValue *PrimaryKey) *SearchRequest {
	r.RoutingValues = append(r.RoutingValues, routingValue)
	return r
}

func (r *SearchRequest) SetTimeoutMs(timeoutMs int32) *SearchRequest {
	r.TimeoutMs = proto.Int32(timeoutMs)
	return r
}

func (r *SearchRequest) ProtoBuffer() (*otsprotocol.SearchRequest, error) {
	req := &otsprotocol.SearchRequest{}
	req.TableName = &r.TableName
	req.IndexName = &r.IndexName
	query, err := r.SearchQuery.Serialize()
	if err != nil {
		return nil, err
	}
	req.SearchQuery = query
	pbColumns := &otsprotocol.ColumnsToGet{}
	pbColumns.ReturnType = otsprotocol.ColumnReturnType_RETURN_NONE.Enum()
	if r.ColumnsToGet != nil {
		if r.ColumnsToGet.ReturnAll {
			pbColumns.ReturnType = otsprotocol.ColumnReturnType_RETURN_ALL.Enum()
		} else if r.ColumnsToGet.ReturnAllFromIndex {
			pbColumns.ReturnType = otsprotocol.ColumnReturnType_RETURN_ALL_FROM_INDEX.Enum()
		} else if len(r.ColumnsToGet.Columns) > 0 {
			pbColumns.ReturnType = otsprotocol.ColumnReturnType_RETURN_SPECIFIED.Enum()
			pbColumns.ColumnNames = r.ColumnsToGet.Columns
		}
	}
	req.ColumnsToGet = pbColumns
	if r.RoutingValues != nil {
		for _, routingValue := range r.RoutingValues {
			req.RoutingValues = append(req.RoutingValues, routingValue.Build(false))
		}
	}
	if r.TimeoutMs != nil {
		req.TimeoutMs = r.TimeoutMs
	}
	return req, err
}

type SearchResponse struct {
	TotalCount   int64
	Rows         []*Row
	IsAllSuccess bool
	NextToken    []byte

	AggregationResults search.AggregationResults
	GroupByResults     search.GroupByResults

	ConsumedCapacityUnit *ConsumedCapacityUnit
	ReservedThroughput   *ReservedThroughput
	ResponseInfo
}

func convertFieldSchemaToPBFieldSchema(fieldSchemas []*FieldSchema) []*otsprotocol.FieldSchema {
	var schemas []*otsprotocol.FieldSchema
	for _, value := range fieldSchemas {
		field := new(otsprotocol.FieldSchema)

		field.FieldName = proto.String(*value.FieldName)
		field.FieldType = otsprotocol.FieldType(int32(value.FieldType)).Enum()

		if value.Index != nil {
			field.Index = proto.Bool(*value.Index)
		} else if value.FieldType != FieldType_NESTED {
			field.Index = proto.Bool(true)
		}
		if value.IndexOptions != nil {
			field.IndexOptions = otsprotocol.IndexOptions(int32(*value.IndexOptions)).Enum()
		}
		if value.Analyzer != nil {
			field.Analyzer = proto.String(string(*value.Analyzer))

			if value.AnalyzerParameter != nil {
				if *value.Analyzer == Analyzer_SingleWord {
					param := &otsprotocol.SingleWordAnalyzerParameter{}
					if value.AnalyzerParameter.(SingleWordAnalyzerParameter).CaseSensitive != nil {
						param.CaseSensitive = proto.Bool(*value.AnalyzerParameter.(SingleWordAnalyzerParameter).CaseSensitive)
					}
					if value.AnalyzerParameter.(SingleWordAnalyzerParameter).DelimitWord != nil {
						param.DelimitWord = proto.Bool(*value.AnalyzerParameter.(SingleWordAnalyzerParameter).DelimitWord)
					}
					if paramBytes, err := proto.Marshal(param); err == nil {
						field.AnalyzerParameter = paramBytes
					}
				} else if *value.Analyzer == Analyzer_Split {
					param := &otsprotocol.SplitAnalyzerParameter{}
					if value.AnalyzerParameter.(SplitAnalyzerParameter).Delimiter != nil {
						param.Delimiter = proto.String(*value.AnalyzerParameter.(SplitAnalyzerParameter).Delimiter)
					}
					if paramBytes, err := proto.Marshal(param); err == nil {
						field.AnalyzerParameter = paramBytes
					}
				} else if *value.Analyzer == Analyzer_Fuzzy {
					fuzzyParam := value.AnalyzerParameter.(FuzzyAnalyzerParameter)
					param := &otsprotocol.FuzzyAnalyzerParameter{}
					if fuzzyParam.MaxChars != 0 {
						param.MaxChars = proto.Int32(fuzzyParam.MaxChars)
					}
					if fuzzyParam.MinChars != 0 {
						param.MinChars = proto.Int32(fuzzyParam.MinChars)
					}
					if paramBytes, err := proto.Marshal(param); err == nil {
						field.AnalyzerParameter = paramBytes
					}
				}
			}
		}
		if value.EnableSortAndAgg != nil {
			field.SortAndAgg = proto.Bool(*value.EnableSortAndAgg)
		}
		if value.Store != nil {
			field.Store = proto.Bool(*value.Store)
		} else if value.FieldType != FieldType_NESTED {
			if *field.FieldType == otsprotocol.FieldType_TEXT {
				field.Store = proto.Bool(false)
			} else {
				field.Store = proto.Bool(true)
			}
		}
		if value.IsArray != nil {
			field.IsArray = proto.Bool(*value.IsArray)
		}
		if value.FieldType == FieldType_NESTED {
			field.FieldSchemas = convertFieldSchemaToPBFieldSchema(value.FieldSchemas)
		}
		if value.IsVirtualField != nil {
			field.IsVirtualField = proto.Bool(*value.IsVirtualField)
		}
		if len(value.SourceFieldNames) != 0 {
			sourceFieldNameArray := make([]string, 0)
			for _, element := range value.SourceFieldNames {
				sourceFieldNameArray = append(sourceFieldNameArray, element)
			}
			field.SourceFieldNames = sourceFieldNameArray
		}
		if len(value.DateFormats) != 0 {
			dateFormatsArray := make([]string, 0)
			for _, element := range value.DateFormats {
				dateFormatsArray = append(dateFormatsArray, element)
			}
			field.DateFormats = dateFormatsArray
		}

		schemas = append(schemas, field)
	}

	return schemas
}

func ConvertToPbSchema(schema *IndexSchema) (*otsprotocol.IndexSchema, error) {
	indexSchema := new(otsprotocol.IndexSchema)
	indexSchema.FieldSchemas = convertFieldSchemaToPBFieldSchema(schema.FieldSchemas)
	indexSchema.IndexSetting = new(otsprotocol.IndexSetting)
	var defaultNumberOfShards int32 = 1
	indexSchema.IndexSetting.NumberOfShards = &defaultNumberOfShards
	if schema.IndexSetting != nil {
		indexSchema.IndexSetting.RoutingFields = schema.IndexSetting.RoutingFields
	}
	if schema.IndexSort != nil {
		pbSort, err := schema.IndexSort.ProtoBuffer()
		if err != nil {
			return nil, err
		}
		indexSchema.IndexSort = pbSort
	}
	return indexSchema, nil
}

func convertToPbQueryFlowWeight(queryFlowWeightArray []*QueryFlowWeight) []*otsprotocol.QueryFlowWeight {
	var queryFlowWeights []*otsprotocol.QueryFlowWeight
	for _, value := range queryFlowWeightArray {
		queryFlowWeight := new(otsprotocol.QueryFlowWeight)
		queryFlowWeight.IndexName = proto.String(value.IndexName)
		queryFlowWeight.Weight = proto.Int32(value.Weight)
		queryFlowWeights = append(queryFlowWeights, queryFlowWeight)
	}
	return queryFlowWeights
}
func parseQueryFlowWeightFromPb(queryFlowWeights []*otsprotocol.QueryFlowWeight) []*QueryFlowWeight {
	var flowWeights []*QueryFlowWeight
	for _, value := range queryFlowWeights {
		queryFlowWeight := new(QueryFlowWeight)
		queryFlowWeight.IndexName = *value.IndexName
		queryFlowWeight.Weight = *value.Weight

		flowWeights = append(flowWeights, queryFlowWeight)
	}
	return flowWeights
}
func parseFieldSchemaFromPb(pbFieldSchemas []*otsprotocol.FieldSchema) []*FieldSchema {
	var schemas []*FieldSchema
	for _, value := range pbFieldSchemas {
		field := new(FieldSchema)
		field.FieldName = value.FieldName
		field.FieldType = FieldType(*value.FieldType)
		field.Index = value.Index
		if value.IndexOptions != nil {
			indexOption := IndexOptions(*value.IndexOptions)
			field.IndexOptions = &indexOption
		}
		field.Analyzer = (*Analyzer)(value.Analyzer)
		if field.Analyzer != nil && *field.Analyzer == Analyzer_SingleWord && value.AnalyzerParameter != nil {
			param := new(otsprotocol.SingleWordAnalyzerParameter)
			if err := proto.Unmarshal(value.AnalyzerParameter, param); err == nil && param != nil {
				p := SingleWordAnalyzerParameter{}
				if param.CaseSensitive != nil {
					p.CaseSensitive = proto.Bool(*param.CaseSensitive)
				}
				if param.DelimitWord != nil {
					p.DelimitWord = proto.Bool(*param.DelimitWord)
				}
				field.AnalyzerParameter = p
			}
		} else if field.Analyzer != nil && *field.Analyzer == Analyzer_Split && value.AnalyzerParameter != nil {
			param := new(otsprotocol.SplitAnalyzerParameter)
			if err := proto.Unmarshal(value.AnalyzerParameter, param); err == nil && param != nil {
				p := SplitAnalyzerParameter{}
				if param.Delimiter != nil {
					p.Delimiter = proto.String(*param.Delimiter)
				}
				field.AnalyzerParameter = p
			}
		} else if field.Analyzer != nil && *field.Analyzer == Analyzer_Fuzzy && value.AnalyzerParameter != nil {
			param := new(otsprotocol.FuzzyAnalyzerParameter)
			if err := proto.Unmarshal(value.AnalyzerParameter, param); err == nil && param != nil {
				p := FuzzyAnalyzerParameter{}
				if param.MinChars != nil {
					p.MinChars = *param.MinChars
				}
				if param.MaxChars != nil {
					p.MaxChars = *param.MaxChars
				}
				field.AnalyzerParameter = p
			}
		}
		field.EnableSortAndAgg = value.SortAndAgg
		field.Store = value.Store
		field.IsArray = value.IsArray
		field.IsVirtualField = value.IsVirtualField
		if value.SourceFieldNames != nil {
			field.SourceFieldNames = value.SourceFieldNames
		}
		if value.DateFormats != nil {
			field.DateFormats = value.DateFormats
		}
		if field.FieldType == FieldType_NESTED {
			field.FieldSchemas = parseFieldSchemaFromPb(value.FieldSchemas)
		}
		schemas = append(schemas, field)
	}
	return schemas
}

func parseIndexSortFromPb(pbIndexSort *otsprotocol.Sort) (*search.Sort, error) {
	indexSort := &search.Sort{
		Sorters: make([]search.Sorter, 0),
	}
	for _, sorter := range pbIndexSort.GetSorter() {
		if sorter.GetFieldSort() != nil {
			fieldSort := &search.FieldSort{
				FieldName: *sorter.GetFieldSort().FieldName,
				Order:     search.ParseSortOrder(sorter.GetFieldSort().Order),
			}
			indexSort.Sorters = append(indexSort.Sorters, fieldSort)
		} else if sorter.GetPkSort() != nil {
			pkSort := &search.PrimaryKeySort{
				Order: search.ParseSortOrder(sorter.GetPkSort().Order),
			}
			indexSort.Sorters = append(indexSort.Sorters, pkSort)
		} else {
			return nil, errors.New("unknown index sort type")
		}
	}
	return indexSort, nil
}

func ParseFromPbSchema(pbSchema *otsprotocol.IndexSchema) (*IndexSchema, error) {
	schema := &IndexSchema{
		IndexSetting: &IndexSetting{
			RoutingFields: pbSchema.IndexSetting.RoutingFields,
		},
	}
	schema.FieldSchemas = parseFieldSchemaFromPb(pbSchema.GetFieldSchemas())
	indexSort, err := parseIndexSortFromPb(pbSchema.GetIndexSort())
	if err != nil {
		return nil, err
	}
	schema.IndexSort = indexSort
	return schema, nil
}

type IndexSchema struct {
	IndexSetting *IndexSetting
	FieldSchemas []*FieldSchema
	IndexSort    *search.Sort
}

type FieldType int32

const (
	FieldType_LONG      FieldType = 1
	FieldType_DOUBLE    FieldType = 2
	FieldType_BOOLEAN   FieldType = 3
	FieldType_KEYWORD   FieldType = 4
	FieldType_TEXT      FieldType = 5
	FieldType_NESTED    FieldType = 6
	FieldType_GEO_POINT FieldType = 7
	FieldType_DATE      FieldType = 8
)

func (ft FieldType) String() string {
	switch ft {
	case FieldType_LONG:
		return "LONG"
	case FieldType_DOUBLE:
		return "DOUBLE"
	case FieldType_BOOLEAN:
		return "BOOLEAN"
	case FieldType_KEYWORD:
		return "KEYWORD"
	case FieldType_TEXT:
		return "TEXT"
	case FieldType_NESTED:
		return "NESTED"
	case FieldType_GEO_POINT:
		return "GEO_POINT"
	case FieldType_DATE:
		return "DATE"
	default:
		return string(ft)
	}
}

func ToFieldType(fieldType string) (FieldType, error) {
	switch strings.ToUpper(fieldType) {
	case "LONG":
		return FieldType_LONG, nil
	case "DOUBLE":
		return FieldType_DOUBLE, nil
	case "BOOLEAN":
		return FieldType_BOOLEAN, nil
	case "KEYWORD":
		return FieldType_KEYWORD, nil
	case "TEXT":
		return FieldType_TEXT, nil
	case "NESTED":
		return FieldType_NESTED, nil
	case "GEO_POINT":
		return FieldType_GEO_POINT, nil
	case "DATE":
		return FieldType_DATE, nil
	default:
		return FieldType_LONG, errors.New("Invalid field type: " + fieldType)
	}
}

func (ft *FieldType) UnmarshalJSON(data []byte) (err error) {
	var ftStr string
	err = json.Unmarshal(data, &ftStr)
	if err != nil {
		return
	}

	*ft, err = ToFieldType(ftStr)
	if err != nil {
		return err
	}
	return
}

func (ft *FieldType) MarshalJSON() (data []byte, err error) {
	data, err = json.Marshal(ft.String())
	return
}

type IndexOptions int32

const (
	IndexOptions_DOCS      IndexOptions = 1
	IndexOptions_FREQS     IndexOptions = 2
	IndexOptions_POSITIONS IndexOptions = 3
	IndexOptions_OFFSETS   IndexOptions = 4
)

type Analyzer string

const (
	Analyzer_SingleWord Analyzer = "single_word"
	Analyzer_MaxWord    Analyzer = "max_word"
	Analyzer_MinWord    Analyzer = "min_word"
	Analyzer_Split      Analyzer = "split"
	Analyzer_Fuzzy      Analyzer = "fuzzy"
)

type SingleWordAnalyzerParameter struct {
	CaseSensitive *bool
	DelimitWord   *bool
}

type SplitAnalyzerParameter struct {
	Delimiter *string
}

type FuzzyAnalyzerParameter struct {
	MinChars int32
	MaxChars int32
}

type FieldSchema struct {
	FieldName         *string
	FieldType         FieldType
	Index             *bool
	IndexOptions      *IndexOptions
	Analyzer          *Analyzer
	AnalyzerParameter interface{}
	EnableSortAndAgg  *bool
	Store             *bool
	IsArray           *bool
	FieldSchemas      []*FieldSchema
	IsVirtualField    *bool
	SourceFieldNames  []string
	DateFormats       []string
}

func (r *FieldSchema) UnmarshalJSON(data []byte) (err error) {
	type FieldSchemaAlias FieldSchema
	copyFS := &FieldSchemaAlias{}
	err = json.Unmarshal(data, copyFS)
	if err != nil {
		return
	}

	r.FieldName = copyFS.FieldName
	r.FieldType = copyFS.FieldType
	r.Index = copyFS.Index
	r.IndexOptions = copyFS.IndexOptions
	r.Analyzer = copyFS.Analyzer
	r.AnalyzerParameter = copyFS.AnalyzerParameter
	r.EnableSortAndAgg = copyFS.EnableSortAndAgg
	r.Store = copyFS.Store
	r.IsArray = copyFS.IsArray
	r.FieldSchemas = copyFS.FieldSchemas
	r.IsVirtualField = copyFS.IsVirtualField
	r.SourceFieldNames = copyFS.SourceFieldNames
	r.DateFormats = copyFS.DateFormats

	apJson, err := json.Marshal(r.AnalyzerParameter)
	if err != nil {
		return
	}

	if r.Analyzer != nil {
		switch *r.Analyzer {
		case Analyzer_Fuzzy:
			ap := &FuzzyAnalyzerParameter{}
			err = json.Unmarshal(apJson, ap)
			r.AnalyzerParameter = *ap
		case Analyzer_Split:
			ap := &SplitAnalyzerParameter{}
			err = json.Unmarshal(apJson, ap)
			r.AnalyzerParameter = *ap
		case Analyzer_SingleWord:
			ap := &SingleWordAnalyzerParameter{}
			err = json.Unmarshal(apJson, ap)
			r.AnalyzerParameter = *ap
		}
	}

	return
}

func (fs *FieldSchema) String() string {
	out, err := json.Marshal(fs)
	if err != nil {
		panic(err)
	}
	return string(out)
}

func (queryFlowWeight *QueryFlowWeight) String() string {
	out, err := json.Marshal(queryFlowWeight)
	if err != nil {
		panic(err)
	}
	return string(out)
}

type IndexSetting struct {
	RoutingFields []string
}

type CreateSearchIndexRequest struct {
	TableName       string
	IndexName       string
	IndexSchema     *IndexSchema
	SourceIndexName *string
	TimeToLive      *int32
}

type CreateSearchIndexResponse struct {
	ResponseInfo ResponseInfo
}

type DescribeSearchIndexRequest struct {
	TableName string
	IndexName string
}

type SyncPhase int32

const (
	SyncPhase_FULL SyncPhase = 1
	SyncPhase_INCR SyncPhase = 2
)

type SyncStat struct {
	SyncPhase            SyncPhase
	CurrentSyncTimestamp *int64
}

type MeteringInfo struct {
	StorageSize    int64
	RowCount       int64
	ReservedReadCU int64
	LastUpdateTime int64
}

type DescribeSearchIndexResponse struct {
	Schema           *IndexSchema
	SyncStat         *SyncStat
	MeteringInfo     *MeteringInfo
	QueryFlowWeights []*QueryFlowWeight
	CreateTime       int64
	TimeToLive       int32
	ResponseInfo     ResponseInfo
}

type ListSearchIndexRequest struct {
	TableName string
}

type IndexInfo struct {
	TableName string
	IndexName string
}

type ListSearchIndexResponse struct {
	IndexInfo    []*IndexInfo
	ResponseInfo ResponseInfo
}

type DeleteSearchIndexRequest struct {
	TableName string
	IndexName string
}

type DeleteSearchIndexResponse struct {
	ResponseInfo ResponseInfo
}

type QueryFlowWeight struct {
	IndexName string
	Weight    int32
}

type UpdateSearchIndexRequest struct {
	TableName        string
	IndexName        string
	SwitchIndexName  *string
	QueryFlowWeights []*QueryFlowWeight
	TimeToLive       *int32
}

type UpdateSearchIndexResponse struct {
	ResponseInfo ResponseInfo
}

//ParallelScan

type ParallelScanRequest struct {
	TableName    string
	IndexName    string
	ScanQuery    search.ScanQuery
	ColumnsToGet *ColumnsToGet
	SessionId    []byte
	TimeoutMs    *int32
}

type ParallelScanResponse struct {
	Rows      []*Row
	NextToken []byte

	ResponseInfo
}

func (r *ParallelScanRequest) SetTableName(tableName string) *ParallelScanRequest {
	r.TableName = tableName
	return r
}

func (r *ParallelScanRequest) SetIndexName(indexName string) *ParallelScanRequest {
	r.IndexName = indexName
	return r
}

func (r *ParallelScanRequest) SetScanQuery(scanQuery search.ScanQuery) *ParallelScanRequest {
	r.ScanQuery = scanQuery
	return r
}

func (r *ParallelScanRequest) SetColumnsToGet(columnsToGet *ColumnsToGet) *ParallelScanRequest {
	r.ColumnsToGet = columnsToGet
	return r
}

func (r *ParallelScanRequest) SetSessionId(sessionId []byte) *ParallelScanRequest {
	r.SessionId = sessionId
	return r
}

func (r *ParallelScanRequest) SetTimeoutMs(timeoutMs int32) *ParallelScanRequest {
	r.TimeoutMs = proto.Int32(timeoutMs)
	return r
}

func (r *ParallelScanRequest) ProtoBuffer() (*otsprotocol.ParallelScanRequest, error) {
	req := &otsprotocol.ParallelScanRequest{}
	req.TableName = proto.String(r.TableName)
	req.IndexName = proto.String(r.IndexName)
	req.SessionId = r.SessionId
	if r.TimeoutMs != nil {
		req.TimeoutMs = r.TimeoutMs
	}

	query, err := r.ScanQuery.Serialize()
	if err != nil {
		return nil, err
	}
	req.ScanQuery = query

	pbColumns := &otsprotocol.ColumnsToGet{}
	pbColumns.ReturnType = otsprotocol.ColumnReturnType_RETURN_NONE.Enum()
	if r.ColumnsToGet != nil {
		if r.ColumnsToGet.ReturnAllFromIndex {
			pbColumns.ReturnType = otsprotocol.ColumnReturnType_RETURN_ALL_FROM_INDEX.Enum()
		} else if r.ColumnsToGet.ReturnAll {
			return nil, errors.New("RETURN_ALL is not allowed for parallel scan")
		} else if len(r.ColumnsToGet.Columns) > 0 {
			pbColumns.ReturnType = otsprotocol.ColumnReturnType_RETURN_SPECIFIED.Enum()
			pbColumns.ColumnNames = r.ColumnsToGet.Columns
		}
	}
	req.ColumnsToGet = pbColumns

	return req, err
}

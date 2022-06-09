package tablestore

type TableStoreApi interface {
	CreateTable(request *CreateTableRequest) (*CreateTableResponse, error)
	ListTable() (*ListTableResponse, error)
	DeleteTable(request *DeleteTableRequest) (*DeleteTableResponse, error)
	DescribeTable(request *DescribeTableRequest) (*DescribeTableResponse, error)
	UpdateTable(request *UpdateTableRequest) (*UpdateTableResponse, error)
	PutRow(request *PutRowRequest) (*PutRowResponse, error)
	DeleteRow(request *DeleteRowRequest) (*DeleteRowResponse, error)
	GetRow(request *GetRowRequest) (*GetRowResponse, error)
	UpdateRow(request *UpdateRowRequest) (*UpdateRowResponse, error)
	BatchGetRow(request *BatchGetRowRequest) (*BatchGetRowResponse, error)
	BatchWriteRow(request *BatchWriteRowRequest) (*BatchWriteRowResponse, error)
	GetRange(request *GetRangeRequest) (*GetRangeResponse, error)

	// stream related
	ListStream(request *ListStreamRequest) (*ListStreamResponse, error)
	DescribeStream(request *DescribeStreamRequest) (*DescribeStreamResponse, error)
	GetShardIterator(request *GetShardIteratorRequest) (*GetShardIteratorResponse, error)
	GetStreamRecord(request *GetStreamRecordRequest) (*GetStreamRecordResponse, error)

	// search related
	CreateSearchIndex(request *CreateSearchIndexRequest) (*CreateSearchIndexResponse, error)
	UpdateSearchIndex(request *UpdateSearchIndexRequest) (*UpdateSearchIndexResponse, error)
	DeleteSearchIndex(request *DeleteSearchIndexRequest) (*DeleteSearchIndexResponse, error)
	ListSearchIndex(request *ListSearchIndexRequest) (*ListSearchIndexResponse, error)
	DescribeSearchIndex(request *DescribeSearchIndexRequest) (*DescribeSearchIndexResponse, error)
	Search(request *SearchRequest) (*SearchResponse, error)

	ComputeSplits(request *ComputeSplitsRequest) (*ComputeSplitsResponse, error)
	ParallelScan(request *ParallelScanRequest) (*ParallelScanResponse, error)
	SQLQuery(req *SQLQueryRequest) (*SQLQueryResponse, error)
}

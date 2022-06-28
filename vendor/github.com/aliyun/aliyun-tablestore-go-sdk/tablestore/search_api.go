package tablestore

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/search"
	"github.com/golang/protobuf/proto"
)

func (tableStoreClient *TableStoreClient) CreateSearchIndex(request *CreateSearchIndexRequest) (*CreateSearchIndexResponse, error) {
	req := new(otsprotocol.CreateSearchIndexRequest)
	req.TableName = proto.String(request.TableName)
	req.IndexName = proto.String(request.IndexName)
	if nil != request.SourceIndexName {
		req.SourceIndexName = request.SourceIndexName
	}
	if request.TimeToLive != nil {
		req.TimeToLive = proto.Int32(*request.TimeToLive)
	}
	var err error
	req.Schema, err = ConvertToPbSchema(request.IndexSchema)
	if err != nil {
		return nil, err
	}
	resp := new(otsprotocol.CreateSearchIndexResponse)
	response := &CreateSearchIndexResponse{}
	if err := tableStoreClient.doRequestWithRetry(createSearchIndexUri, req, resp, &response.ResponseInfo); err != nil {
		return nil, err
	}
	return response, nil
}

func (tableStoreClient *TableStoreClient) UpdateSearchIndex(request *UpdateSearchIndexRequest) (*UpdateSearchIndexResponse, error) {
	req := new(otsprotocol.UpdateSearchIndexRequest)
	req.TableName = proto.String(request.TableName)
	req.IndexName = proto.String(request.IndexName)
	req.QueryFlowWeight = convertToPbQueryFlowWeight(request.QueryFlowWeights)
	req.TimeToLive = request.TimeToLive
	req.SwitchIndexName = request.SwitchIndexName
	resp := new(otsprotocol.UpdateSearchIndexResponse)
	response := new(UpdateSearchIndexResponse)
	if err := tableStoreClient.doRequestWithRetry(updateSearchIndexUri, req, resp, &response.ResponseInfo); err != nil {
		return nil, err
	}
	return response, nil
}

func (tableStoreClient *TableStoreClient) DeleteSearchIndex(request *DeleteSearchIndexRequest) (*DeleteSearchIndexResponse, error) {
	req := new(otsprotocol.DeleteSearchIndexRequest)
	req.TableName = proto.String(request.TableName)
	req.IndexName = proto.String(request.IndexName)

	resp := new(otsprotocol.DeleteSearchIndexResponse)
	response := &DeleteSearchIndexResponse{}
	if err := tableStoreClient.doRequestWithRetry(deleteSearchIndexUri, req, resp, &response.ResponseInfo); err != nil {
		return nil, err
	}
	return response, nil
}

func (tableStoreClient *TableStoreClient) ListSearchIndex(request *ListSearchIndexRequest) (*ListSearchIndexResponse, error) {
	req := new(otsprotocol.ListSearchIndexRequest)
	req.TableName = proto.String(request.TableName)

	resp := new(otsprotocol.ListSearchIndexResponse)
	response := &ListSearchIndexResponse{}
	if err := tableStoreClient.doRequestWithRetry(listSearchIndexUri, req, resp, &response.ResponseInfo); err != nil {
		return nil, err
	}
	indexs := make([]*IndexInfo, 0)
	for _, info := range resp.Indices {
		indexs = append(indexs, &IndexInfo{
			TableName: *info.TableName,
			IndexName: *info.IndexName,
		})
	}
	response.IndexInfo = indexs
	return response, nil
}

func (tableStoreClient *TableStoreClient) DescribeSearchIndex(request *DescribeSearchIndexRequest) (*DescribeSearchIndexResponse, error) {
	req := new(otsprotocol.DescribeSearchIndexRequest)
	req.TableName = proto.String(request.TableName)
	req.IndexName = proto.String(request.IndexName)

	resp := new(otsprotocol.DescribeSearchIndexResponse)
	response := &DescribeSearchIndexResponse{}
	if err := tableStoreClient.doRequestWithRetry(describeSearchIndexUri, req, resp, &response.ResponseInfo); err != nil {
		return nil, err
	}
	schema, err := ParseFromPbSchema(resp.Schema)
	if err != nil {
		return nil, err
	}
	response.Schema = schema
	if resp.SyncStat != nil {
		response.SyncStat = &SyncStat{
			CurrentSyncTimestamp: resp.SyncStat.CurrentSyncTimestamp,
		}
		syncPhase := resp.SyncStat.SyncPhase
		if syncPhase == nil {
			return nil, errors.New("missing [SyncPhase] in DescribeSearchIndexResponse")
		} else if *syncPhase == otsprotocol.SyncPhase_FULL {
			response.SyncStat.SyncPhase = SyncPhase_FULL
		} else if *syncPhase == otsprotocol.SyncPhase_INCR {
			response.SyncStat.SyncPhase = SyncPhase_INCR
		} else {
			return nil, errors.New(fmt.Sprintf("unknown SyncPhase: %v", syncPhase))
		}
	}

	if resp.MeteringInfo != nil {
		response.MeteringInfo = &MeteringInfo{}

		if resp.MeteringInfo.StorageSize != nil {
			response.MeteringInfo.StorageSize = *resp.MeteringInfo.StorageSize
		}

		if resp.MeteringInfo.RowCount != nil {
			response.MeteringInfo.RowCount = *resp.MeteringInfo.RowCount
		}

		if resp.MeteringInfo.ReservedReadCu != nil {
			response.MeteringInfo.ReservedReadCU = *resp.MeteringInfo.ReservedReadCu
		}

		if resp.MeteringInfo.Timestamp != nil {
			response.MeteringInfo.LastUpdateTime = *resp.MeteringInfo.Timestamp
		}
	}

	if resp.CreateTime != nil {
		response.CreateTime = *resp.CreateTime
	}

	if resp.TimeToLive != nil {
		response.TimeToLive = *resp.TimeToLive
	}

	if resp.QueryFlowWeight != nil {
		response.QueryFlowWeights = parseQueryFlowWeightFromPb(resp.QueryFlowWeight)
	}

	return response, nil
}

func (tableStoreClient *TableStoreClient) Search(request *SearchRequest) (*SearchResponse, error) {
	req, err := request.ProtoBuffer()
	if err != nil {
		return nil, err
	}
	resp := new(otsprotocol.SearchResponse)
	response := &SearchResponse{}
	if err := tableStoreClient.doRequestWithRetry(searchUri, req, resp, &response.ResponseInfo); err != nil {
		return nil, err
	}
	response.TotalCount = *resp.TotalHits

	rows := make([]*PlainBufferRow, 0)
	for _, buf := range resp.Rows {
		row, err := readRowsWithHeader(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		rows = append(rows, row[0])
	}

	for _, row := range rows {
		currentRow := &Row{}
		currentPk := new(PrimaryKey)
		for _, pk := range row.primaryKey {
			pkColumn := &PrimaryKeyColumn{ColumnName: string(pk.cellName), Value: pk.cellValue.Value}
			currentPk.PrimaryKeys = append(currentPk.PrimaryKeys, pkColumn)
		}
		currentRow.PrimaryKey = currentPk
		for _, cell := range row.cells {
			dataColumn := &AttributeColumn{ColumnName: string(cell.cellName), Value: cell.cellValue.Value, Timestamp: cell.cellTimestamp}
			currentRow.Columns = append(currentRow.Columns, dataColumn)
		}
		response.Rows = append(response.Rows, currentRow)
	}

	response.IsAllSuccess = *resp.IsAllSucceeded
	if resp.NextToken != nil && len(resp.NextToken) > 0 {
		response.NextToken = resp.NextToken
	}

	pbAggResults := new(otsprotocol.AggregationsResult)
	if err := proto.Unmarshal(resp.Aggs, pbAggResults); err == nil && pbAggResults != nil && len(pbAggResults.AggResults) > 0 {
		aggResults, err := search.ParseAggregationResultsFromPB(pbAggResults.AggResults)
		if err != nil {
			return nil, err
		}
		response.AggregationResults = *aggResults
	}

	pbGroupByResults := new(otsprotocol.GroupBysResult)
	if err = proto.Unmarshal(resp.GroupBys, pbGroupByResults); err == nil && pbGroupByResults != nil && len(pbGroupByResults.GroupByResults) > 0 {
		groupByResults, err := search.ParseGroupByResultsFromPB(pbGroupByResults.GroupByResults)
		if err != nil {
			return nil, err
		}
		response.GroupByResults = *groupByResults
	}

	if resp.Consumed != nil && resp.Consumed.CapacityUnit != nil {
		response.ConsumedCapacityUnit = &ConsumedCapacityUnit{
			Read:  resp.Consumed.CapacityUnit.GetRead(),
			Write: resp.Consumed.CapacityUnit.GetWrite(),
		}
	}

	if resp.ReservedConsumed != nil && resp.ReservedConsumed.CapacityUnit != nil {
		response.ReservedThroughput = &ReservedThroughput{
			Readcap:  int(resp.ReservedConsumed.CapacityUnit.GetRead()),
			Writecap: int(resp.ReservedConsumed.CapacityUnit.GetWrite()),
		}
	}

	return response, nil
}

func (TableStoreClient *TableStoreClient) ParallelScan(request *ParallelScanRequest) (*ParallelScanResponse, error) {
	req, err := request.ProtoBuffer()
	if err != nil {
		return nil, err
	}
	resp := new(otsprotocol.ParallelScanResponse)
	response := &ParallelScanResponse{}
	if err := TableStoreClient.doRequestWithRetry(parallelScanUri, req, resp, &response.ResponseInfo); err != nil {
		return nil, err
	}

	rows := make([]*PlainBufferRow, 0)
	for _, buf := range resp.Rows {
		row, err := readRowsWithHeader(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		rows = append(rows, row[0])
	}
	for _, row := range rows {
		currentRow := &Row{}
		currentPk := new(PrimaryKey)
		for _, pk := range row.primaryKey {
			pkColumn := &PrimaryKeyColumn{ColumnName: string(pk.cellName), Value: pk.cellValue.Value}
			currentPk.PrimaryKeys = append(currentPk.PrimaryKeys, pkColumn)
		}
		currentRow.PrimaryKey = currentPk
		for _, cell := range row.cells {
			dataColumn := &AttributeColumn{ColumnName: string(cell.cellName), Value: cell.cellValue.Value, Timestamp: cell.cellTimestamp}
			currentRow.Columns = append(currentRow.Columns, dataColumn)
		}
		response.Rows = append(response.Rows, currentRow)
	}

	if resp.NextToken != nil && len(resp.NextToken) > 0 {
		response.NextToken = resp.NextToken
	}

	return response, nil
}

package alicloud

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/search"

	otsTunnel "github.com/aliyun/aliyun-tablestore-go-sdk/tunnel"

	"time"

	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type OtsService struct {
	client *connectivity.AliyunClient
}

func (s *OtsService) getPrimaryKeyType(primaryKeyType string) tablestore.PrimaryKeyType {
	var keyType tablestore.PrimaryKeyType
	t := PrimaryKeyTypeString(primaryKeyType)
	switch t {
	case IntegerType:
		keyType = tablestore.PrimaryKeyType_INTEGER
	case StringType:
		keyType = tablestore.PrimaryKeyType_STRING
	case BinaryType:
		keyType = tablestore.PrimaryKeyType_BINARY
	}
	return keyType
}

func ParseDefinedColumnType(colType string) (tablestore.DefinedColumnType, error) {
	switch DefinedColumnTypeString(colType) {
	case DefinedColumnInteger:
		return tablestore.DefinedColumn_INTEGER, nil
	case DefinedColumnString:
		return tablestore.DefinedColumn_STRING, nil
	case DefinedColumnBinary:
		return tablestore.DefinedColumn_BINARY, nil
	case DefinedColumnDouble:
		return tablestore.DefinedColumn_DOUBLE, nil
	case DefinedColumnBoolean:
		return tablestore.DefinedColumn_BOOLEAN, nil
	}
	return 0, WrapError(fmt.Errorf("unsupported defined column type: %s", colType))
}

func (s *OtsService) ListOtsTable(instanceName string) (table *tablestore.ListTableResponse, err error) {
	if _, err := s.DescribeOtsInstance(instanceName); err != nil {
		return nil, WrapError(err)
	}
	var raw interface{}
	var requestInfo *tablestore.TableStoreClient
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithTableStoreClient(instanceName, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
			requestInfo = tableStoreClient
			return tableStoreClient.ListTable()
		})
		if err != nil {
			if strings.HasSuffix(err.Error(), "no such host") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("ListTable", raw, requestInfo)
		return nil
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), "OTSObjectNotExist") {
			return table, WrapErrorf(err, NotFoundMsg, AliyunTablestoreGoSdk)
		}
		return nil, WrapErrorf(err, DataDefaultErrorMsg, instanceName, "ListTable", AliyunTablestoreGoSdk)
	}
	table, _ = raw.(*tablestore.ListTableResponse)
	if table == nil {
		return table, WrapErrorf(NotFoundErr("OtsTable", instanceName), NotFoundMsg, ProviderERROR)
	}
	return
}

func (s *OtsService) DescribeOtsTable(id string) (*tablestore.DescribeTableResponse, error) {
	table := &tablestore.DescribeTableResponse{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return table, WrapError(err)
	}
	instanceName, tableName := parts[0], parts[1]
	request := new(tablestore.DescribeTableRequest)
	request.TableName = tableName

	if _, err := s.DescribeOtsInstance(instanceName); err != nil {
		return table, WrapError(err)
	}
	var raw interface{}
	var requestInfo *tablestore.TableStoreClient
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithTableStoreClient(instanceName, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
			requestInfo = tableStoreClient
			return tableStoreClient.DescribeTable(request)
		})
		if err != nil {
			if IsExpectedErrors(err, OtsTableIsTemporarilyUnavailable) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("DescribeTable", raw, requestInfo, request)
		return nil
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), "OTSObjectNotExist") {
			return table, WrapErrorf(err, NotFoundMsg, AliyunTablestoreGoSdk)
		}
		return table, WrapErrorf(err, DefaultErrorMsg, id, "DescribeTable", AliyunTablestoreGoSdk)
	}
	table, _ = raw.(*tablestore.DescribeTableResponse)
	if table == nil || table.TableMeta == nil || table.TableMeta.TableName != tableName {
		return table, WrapErrorf(NotFoundErr("OtsTable", id), NotFoundMsg, ProviderERROR)
	}
	return table, nil
}

func (s *OtsService) WaitForOtsTable(instanceName, tableName string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	id := fmt.Sprintf("%s%s%s", instanceName, COLON_SEPARATED, tableName)

	for {
		object, err := s.DescribeOtsTable(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.TableMeta.TableName == tableName && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.TableMeta.TableName, tableName, ProviderERROR)
		}

	}
}

// Convert tablestore.PrimaryKeyType to PrimaryKeyTypeString
func (s *OtsService) convertPrimaryKeyType(t tablestore.PrimaryKeyType) PrimaryKeyTypeString {
	var typeString PrimaryKeyTypeString
	switch t {
	case tablestore.PrimaryKeyType_INTEGER:
		typeString = IntegerType
	case tablestore.PrimaryKeyType_BINARY:
		typeString = BinaryType
	case tablestore.PrimaryKeyType_STRING:
		typeString = StringType
	}
	return typeString
}

func ConvertDefinedColumnType(t tablestore.DefinedColumnType) (DefinedColumnTypeString, error) {
	switch t {
	case tablestore.DefinedColumn_INTEGER:
		return DefinedColumnInteger, nil
	case tablestore.DefinedColumn_STRING:
		return DefinedColumnString, nil
	case tablestore.DefinedColumn_BINARY:
		return DefinedColumnBinary, nil
	case tablestore.DefinedColumn_DOUBLE:
		return DefinedColumnDouble, nil
	case tablestore.DefinedColumn_BOOLEAN:
		return DefinedColumnBoolean, nil
	}
	return "", WrapError(fmt.Errorf("unsupported defined column type: %v", t))
}

func FindDefinedColumn(columns []*tablestore.DefinedColumnSchema, target *tablestore.DefinedColumnSchema) ColumnFindResult {
	for _, column := range columns {
		if column.Name == target.Name {
			if column.ColumnType == target.ColumnType {
				return ExistEqual
			}
			return ExistNotEqual
		}
	}
	return NotExist
}

type ColumnFindResult int32

const (
	ExistEqual    ColumnFindResult = 1
	ExistNotEqual ColumnFindResult = 2
	NotExist      ColumnFindResult = 3
)

func (s *OtsService) ListOtsInstance(maxResults int) (allInstanceNames []string, err error) {
	actionPath := "/v2/openapi/listinstances"
	request := make(map[string]*string)
	request["RegionId"] = StringPointer(s.client.RegionId)
	request["MaxResults"] = StringPointer(strconv.Itoa(maxResults))

	for {
		resp, err := OtsRestApiGetWithRetry(s.client, "tablestore", "2020-12-09", actionPath, request)
		if err != nil {
			return nil, WrapErrorf(err, DefaultErrorMsg, "alicloud_ots_instances", actionPath, AlibabaCloudSdkGoERROR)
		}
		addDebug(actionPath, resp, request)

		// resp struct: {"_headers": {...}, "body": {...}}
		respBody, ok := resp["body"].(map[string]interface{})
		if !ok {
			return allInstanceNames, WrapErrorf(errors.New("parse resp body to map[string]interface{} failed"), DefaultErrorMsg, "instance:*", actionPath, AlibabaCloudSdkGoERROR)
		}
		// respBody["Instances"] struct: [{}, {}, {}]
		instanceMaps := respBody["Instances"]
		// Convert map to json string
		instancesJSON, err := json.Marshal(instanceMaps)
		if err != nil {
			return allInstanceNames, WrapErrorf(err, DefaultErrorMsg, "instance:*", actionPath, AlibabaCloudSdkGoERROR)
		}
		// Convert json string to obj
		var instances []RestOtsInstanceInfo
		if err := json.Unmarshal(instancesJSON, &instances); err != nil {
			return allInstanceNames, WrapErrorf(err, DefaultErrorMsg, "instance:*", actionPath, AlibabaCloudSdkGoERROR)
		}

		if instances == nil || len(instances) < 1 {
			break
		}

		for _, instance := range instances {
			allInstanceNames = append(allInstanceNames, instance.InstanceName)
		}

		nextToken, _ := resp["NextToken"].(string)
		if len(instances) < maxResults || nextToken == "" {
			break
		} else {
			request["NextToken"] = &nextToken
		}
	}
	return allInstanceNames, nil
}

func (s *OtsService) DescribeOtsInstance(instanceName string) (inst RestOtsInstanceInfo, err error) {
	actionPath := "/v2/openapi/getinstance"
	request := make(map[string]*string)
	request["RegionId"] = StringPointer(s.client.RegionId)
	request["InstanceName"] = StringPointer(instanceName)

	//client := meta.(*connectivity.AliyunClient)
	resp, err := OtsRestApiGetWithRetry(s.client, "tablestore", "2020-12-09", actionPath, request)
	if err != nil {
		if NotFoundError(err) {
			return inst, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return inst, WrapErrorf(err, DefaultErrorMsg, instanceName, actionPath, AlibabaCloudSdkGoERROR)
	}
	addDebug(actionPath, resp, request)

	// resp struct: {"_headers": {...}, "body": {...}}
	instMap := resp["body"]
	// Convert map to json string
	instJSON, err := json.Marshal(instMap)
	if err != nil {
		return inst, WrapErrorf(err, DefaultErrorMsg, instanceName, actionPath, AlibabaCloudSdkGoERROR)
	}
	// Convert json string to obj
	if err := json.Unmarshal(instJSON, &inst); err != nil {
		return inst, WrapErrorf(err, DefaultErrorMsg, instanceName, actionPath, AlibabaCloudSdkGoERROR)
	}

	if inst.InstanceName != instanceName {
		return inst, WrapErrorf(NotFoundErr("OtsInstance", instanceName), NotFoundMsg, ProviderERROR)
	}

	return inst, nil
}

func (s *OtsService) DescribeOtsInstanceAttachment(id string) (inst ots.VpcInfo, err error) {
	request := ots.CreateListVpcInfoByInstanceRequest()
	request.RegionId = s.client.RegionId
	request.Method = "GET"
	request.InstanceName = id
	raw, err := s.client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
		return otsClient.ListVpcInfoByInstance(request)
	})
	if err != nil {
		if NotFoundError(err) {
			return inst, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return inst, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	resp, _ := raw.(*ots.ListVpcInfoByInstanceResponse)
	if resp.TotalCount < 1 {
		return inst, WrapErrorf(NotFoundErr("OtsInstanceAttachment", id), NotFoundMsg, ProviderERROR)
	}
	return resp.VpcInfos.VpcInfo[0], nil
}

func (s *OtsService) WaitForOtsInstanceVpc(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeOtsInstanceAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.InstanceName == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.InstanceName, id, ProviderERROR)
		}

	}
}

func (s *OtsService) ListOtsInstanceVpc(id string) (inst []ots.VpcInfo, err error) {
	request := ots.CreateListVpcInfoByInstanceRequest()
	request.RegionId = s.client.RegionId
	request.Method = "GET"
	request.InstanceName = id
	raw, err := s.client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
		return otsClient.ListVpcInfoByInstance(request)
	})
	if err != nil {
		return inst, WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ots_instance_attachments", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	resp, _ := raw.(*ots.ListVpcInfoByInstanceResponse)
	if resp.TotalCount < 1 {
		return inst, WrapErrorf(NotFoundErr("OtsInstanceAttachment", id), NotFoundMsg, ProviderERROR)
	}

	var retInfos []ots.VpcInfo
	for _, vpcInfo := range resp.VpcInfos.VpcInfo {
		vpcInfo.InstanceName = id
		retInfos = append(retInfos, vpcInfo)
	}
	return retInfos, nil
}

func (s *OtsService) WaitForOtsInstance(id string, instanceInnerStatus string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		instance, err := s.DescribeOtsInstance(id)
		if err != nil {
			if NotFoundError(err) {
				if instanceInnerStatus == string(Deleted) {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if instance.InstanceStatus == instanceInnerStatus {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, fmt.Sprint(instance.InstanceStatus), instanceInnerStatus, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *OtsService) DescribeOtsInstanceTypes() (types []string, err error) {
	request := ots.CreateListClusterTypeRequest()
	request.Method = requests.GET
	raw, err := s.client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
		return otsClient.ListClusterType(request)
	})
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "alicloud_ots_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	resp, _ := raw.(*ots.ListClusterTypeResponse)
	if resp != nil {
		return resp.ClusterTypeInfos.ClusterType, nil
	}
	return
}

func isOtsTunnelNotFound(err error) bool {
	if e, ok := err.(*otsTunnel.TunnelError); ok {
		if e.Code == otsTunnel.ErrCodeParamInvalid && strings.Contains(e.Message, "tunnel not exist") {
			return true
		}
		if e.Code == otsTunnel.ErrCodePermissionDenied && strings.Contains(e.Message, "Instance not found") {
			return true
		}
	}
	return false
}

func (s *OtsService) DescribeOtsTunnel(id string) (resp *otsTunnel.DescribeTunnelResponse, err error) {
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return nil, WrapError(err)
	}

	instanceName, tableName, tunnelName := parts[0], parts[1], parts[2]
	request := new(otsTunnel.DescribeTunnelRequest)
	request.TableName = tableName
	request.TunnelName = tunnelName
	var raw interface{}
	var requestInfo otsTunnel.TunnelClient
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithTableStoreTunnelClient(instanceName, func(tunnelClient otsTunnel.TunnelClient) (interface{}, error) {
			requestInfo = tunnelClient
			return tunnelClient.DescribeTunnel(request)
		})
		if err != nil {
			if IsExpectedErrors(err, OtsTunnelIsTemporarilyUnavailable) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		resp, _ := raw.(*otsTunnel.DescribeTunnelResponse)
		if resp != nil && resp.Tunnel != nil && resp.Tunnel.Stage == "InitBaseDataAndStreamShard" {
			return resource.RetryableError(WrapError(Error("ots tunnel is initial")))
		}
		addDebug("DescribeTunnel", raw, requestInfo, request)
		return nil
	})
	if err != nil {
		if isOtsTunnelNotFound(err) {
			return nil, WrapErrorf(err, NotFoundMsg, AliyunTablestoreGoSdk)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, "DescribeTunnel", AliyunTablestoreGoSdk)
	}
	resp, _ = raw.(*otsTunnel.DescribeTunnelResponse)
	if resp == nil || resp.Tunnel == nil || resp.Tunnel.TableName != tableName || resp.Tunnel.TunnelName != tunnelName {
		return nil, WrapErrorf(NotFoundErr("OtsTunnel", id), NotFoundMsg, ProviderERROR)
	}
	return resp, nil
}

func (s *OtsService) ListOtsTunnels(instanceName string, tableName string) (resp *otsTunnel.ListTunnelResponse, err error) {
	// check table exists
	id := ID(instanceName, tableName)
	if _, err := s.DescribeOtsTable(id); err != nil {
		return nil, WrapError(err)
	}

	var raw interface{}
	var requestInfo otsTunnel.TunnelClient
	request := new(otsTunnel.ListTunnelRequest)
	request.TableName = tableName
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithTableStoreTunnelClient(instanceName, func(tunnelClient otsTunnel.TunnelClient) (interface{}, error) {
			requestInfo = tunnelClient
			return tunnelClient.ListTunnel(request)
		})
		if err != nil {
			if IsExpectedErrors(err, OtsTunnelIsTemporarilyUnavailable) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("ListTunnel", raw, requestInfo, request)
		return nil
	})
	if err != nil {
		if isOtsTunnelNotFound(err) {
			return nil, WrapErrorf(err, NotFoundMsg, AliyunTablestoreGoSdk)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, "DescribeTunnel", AliyunTablestoreGoSdk)
	}
	resp, _ = raw.(*otsTunnel.ListTunnelResponse)
	if resp == nil {
		return nil, WrapErrorf(NotFoundErr("OtsTunnel", id), NotFoundMsg, ProviderERROR)
	}
	return resp, nil
}

func (s *OtsService) WaitForOtsTunnel(id string, status Status, timeout int) error {
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeOtsTunnel(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.Tunnel.TunnelName == parts[2] && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Tunnel.TunnelName, parts[2], ProviderERROR)
		}
	}
}

func ID(segName ...string) string {
	return strings.Join(segName, COLON_SEPARATED)
}

func (s *OtsService) ListOtsSecondaryIndex(instanceName string, tableName string) ([]*tablestore.IndexMeta, error) {
	tableResp, err := s.DescribeOtsTable(ID(instanceName, tableName))
	if err != nil {
		return nil, WrapError(err)
	}
	if tableResp == nil {
		return nil, nil
	}

	return tableResp.IndexMetas, nil
}

// DescribeOtsSecondaryIndex The describe method is depended on by AccTest,
// and the second return value of the describe method needs `error` type
func (s *OtsService) DescribeOtsSecondaryIndex(id string) (index *TableIndex, err error) {
	instanceName, tableName, indexName, _, err := ParseIndexId(id)
	if err != nil {
		return
	}
	tableResp, err := s.DescribeOtsTable(ID(instanceName, tableName))
	if err != nil {
		return
	}
	if tableResp == nil {
		err = WrapError(fmt.Errorf("table not exist: %s", tableName))
		return
	}

	for _, idx := range tableResp.IndexMetas {
		if idx.IndexName == indexName {
			index = &TableIndex{
				InstanceName: instanceName,
				TableName:    tableName,
				Index:        idx,
			}
			return
		}
	}
	err = WrapError(fmt.Errorf("index not exist: %s.%s", tableName, indexName))
	return
}

type TableIndex struct {
	InstanceName string
	TableName    string
	Index        *tablestore.IndexMeta
}

func ConvertSecIndexType(indexType tablestore.IndexType) (SecondaryIndexTypeString, error) {
	switch indexType {
	case tablestore.IT_GLOBAL_INDEX:
		return Global, nil
	case tablestore.IT_LOCAL_INDEX:
		return Local, nil
	default:
		return "", WrapError(fmt.Errorf("unexpected secondary index type: %v", indexType))
	}
}
func ConvertSecIndexTypeString(typeStr SecondaryIndexTypeString) (tablestore.IndexType, error) {
	switch typeStr {
	case Global:
		return tablestore.IT_GLOBAL_INDEX, nil
	case Local:
		return tablestore.IT_LOCAL_INDEX, nil
	default:
		return 0, WrapError(fmt.Errorf("unexpected secondary index type: %v", typeStr))
	}
}

type RegxFilter struct {
	regx           *regexp.Regexp
	getSourceValue func(sourceObj interface{}) interface{}
}

func (f *RegxFilter) filter(sourceObj interface{}) bool {
	return f.regx.MatchString(f.getSourceValue(sourceObj).(string))
}

type ValuesFilter struct {
	allowedValues  []interface{}
	getSourceValue func(sourceObj interface{}) interface{}
}

func (f *ValuesFilter) filter(sourceObj interface{}) bool {
	for _, allowed := range f.allowedValues {
		if allowed != nil && allowed == f.getSourceValue(sourceObj) {
			return true
		}
	}
	// source value not in the enumerated values
	return false
}

type DataFilter interface {
	filter(sourceObj interface{}) bool
}
type InputDataSource struct {
	inputs  []interface{}
	filters []DataFilter
}

func (ds *InputDataSource) doFilters() []interface{} {
	var outputs []interface{}
	for _, input := range ds.inputs {
		pass := true
		for _, filter := range ds.filters {
			if !filter.filter(input) {
				pass = false
				break
			}
		}
		if pass {
			outputs = append(outputs, input)
		}
	}
	return outputs
}

func (s *OtsService) LoopWaitTable(instanceName string, tableName string) (table *tablestore.DescribeTableResponse, err error) {
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		t, e := s.DescribeOtsTable(ID(instanceName, tableName))
		if e != nil {
			if NotFoundError(e) {
				return resource.RetryableError(e)
			}
			return resource.NonRetryableError(e)
		}

		table = t
		return nil
	})
	return
}

func (s *OtsService) WaitForSecondaryIndex(instance string, table string, index string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	id := ID(instance, table)
	for {
		tableResp, err := s.DescribeOtsTable(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
		}
		// table exists and index does not exist
		indexFind := IsSubCollection([]string{index}, simplifySecIndex(tableResp.IndexMetas))
		switch {
		case status == Deleted, !indexFind:
			return nil
		case status != Deleted, indexFind:
			// Non-deleted states cannot be distinguished precisely. If the index exists,
			// it is considered that the index successfully matches the non-deleted state, and end waiting
			return nil
		default:
			break
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, tableResp.TableMeta.TableName, index, ProviderERROR)
		}
	}
}

// ParseIndexId both secondary index IDs and search index IDs will use this method, they consist of the same fields
func ParseIndexId(indexId string) (instanceName, tableName, indexName, indexTypeStr string, err error) {
	splits := strings.Split(indexId, COLON_SEPARATED)
	if len(splits) >= 4 {
		instanceName = splits[0]
		tableName = splits[1]
		indexName = splits[2]
		indexTypeStr = splits[3]
	} else {
		err = WrapError(fmt.Errorf("invalid index id(instanceName:tableName:indexName:indexType): %s", indexId))
	}
	return
}

func (s *OtsService) ListOtsSearchIndex(instanceName string, tableName string) (indexes []*tablestore.IndexInfo, err error) {
	// check table exists
	id := ID(instanceName, tableName)
	if _, err := s.DescribeOtsTable(id); err != nil {
		return nil, WrapError(err)
	}

	req := &tablestore.ListSearchIndexRequest{
		TableName: tableName,
	}
	var raw interface{}
	var reqClient *tablestore.TableStoreClient
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithTableStoreClient(instanceName, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
			reqClient = tableStoreClient
			return tableStoreClient.ListSearchIndex(req)
		})
		if err != nil {
			if IsExpectedErrors(err, OtsSearchIndexIsTemporarilyUnavailable) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("ListSearchIndex", raw, reqClient, req)
		return nil
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), "OTSObjectNotExist") {
			return nil, WrapErrorf(err, NotFoundMsg, AliyunTablestoreGoSdk)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, "ListOtsSearchIndex", AliyunTablestoreGoSdk)
	}
	resp, _ := raw.(*tablestore.ListSearchIndexResponse)
	if resp == nil {
		return nil, WrapErrorf(NotFoundErr("SearchIndex", id), NotFoundMsg, ProviderERROR)
	}
	// IndexInfo slice can be nil when table not has search index
	return resp.IndexInfo, nil
}

func (s *OtsService) DescribeOtsSearchIndex(id string) (indexResp *tablestore.DescribeSearchIndexResponse, err error) {
	instanceName, tableName, indexName, _, err := ParseIndexId(id)
	if err != nil {
		return nil, WrapError(err)
	}
	if _, err = s.DescribeOtsInstance(instanceName); err != nil {
		return nil, WrapError(err)
	}

	if _, err := s.DescribeOtsTable(ID(instanceName, tableName)); err != nil {
		if NotFoundError(err) {
			return nil, WrapError(err)
		}
	}

	req := &tablestore.DescribeSearchIndexRequest{
		TableName: tableName,
		IndexName: indexName,
	}

	var raw interface{}
	var reqClient *tablestore.TableStoreClient
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithTableStoreClient(instanceName, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
			reqClient = tableStoreClient
			return tableStoreClient.DescribeSearchIndex(req)
		})
		defer func() {
			addDebug("DescribeSearchIndex", raw, reqClient, req)
		}()

		if err != nil {
			if IsExpectedErrors(err, OtsSearchIndexIsTemporarilyUnavailable) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		if strings.HasPrefix(err.Error(), "OTSObjectNotExist") {
			return nil, WrapErrorf(err, NotFoundMsg, AliyunTablestoreGoSdk)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, "DescribeSearchIndex", AliyunTablestoreGoSdk)
	}

	indexResp, _ = raw.(*tablestore.DescribeSearchIndexResponse)
	if indexResp == nil || indexResp.SyncStat == nil || indexResp.Schema == nil {
		return nil, WrapErrorf(NotFoundErr("OtsSearchIndex", id), NotFoundMsg, ProviderERROR)
	}
	return indexResp, nil
}

func ConvertSearchIndexSyncPhase(syncPhase tablestore.SyncPhase) (OtsSearchIndexSyncPhaseString, error) {
	switch syncPhase {
	case tablestore.SyncPhase_FULL:
		return Full, nil
	case tablestore.SyncPhase_INCR:
		return Incr, nil
	default:
		return "", WrapError(fmt.Errorf("unexpected search index sync phase: %v", syncPhase))
	}
}

func ConvertSearchIndexFieldTypeString(typeStr SearchIndexFieldTypeString) (tablestore.FieldType, error) {
	switch typeStr {
	case "Long":
		return tablestore.FieldType_LONG, nil
	case "Double":
		return tablestore.FieldType_DOUBLE, nil
	case "Boolean":
		return tablestore.FieldType_BOOLEAN, nil
	case "Keyword":
		return tablestore.FieldType_KEYWORD, nil
	case "Text":
		return tablestore.FieldType_TEXT, nil
	case "Date":
		return tablestore.FieldType_DATE, nil
	case "GeoPoint":
		return tablestore.FieldType_GEO_POINT, nil
	case "Nested":
		return tablestore.FieldType_NESTED, nil
	default:
		return 0, WrapError(fmt.Errorf("unexpected search index field type string: %s", typeStr))
	}
}

func ConvertSearchIndexAnalyzerTypeString(typeStr SearchIndexAnalyzerTypeString) (tablestore.Analyzer, error) {
	switch typeStr {
	case "SingleWord":
		return tablestore.Analyzer_SingleWord, nil
	case "Split":
		return tablestore.Analyzer_Split, nil
	case "MinWord":
		return tablestore.Analyzer_MinWord, nil
	case "MaxWord":
		return tablestore.Analyzer_MaxWord, nil
	case "Fuzzy":
		return tablestore.Analyzer_Fuzzy, nil
	default:
		return "", WrapError(fmt.Errorf("unexpected search index analyzer type string: %s", typeStr))
	}
}

func ConvertSearchIndexSortFieldTypeString(typeStr SearchIndexSortFieldTypeString) (search.Sorter, error) {
	switch typeStr {
	case "PrimaryKeySort":
		return &search.PrimaryKeySort{}, nil
	case "FieldSort":
		return &search.FieldSort{}, nil
	default:
		return nil, WrapError(fmt.Errorf("unexpected search index sort field type string [PrimaryKeySort|FieldSort]: %s", typeStr))
	}
}

func ConvertSearchIndexOrderTypeString(typeStr SearchIndexOrderTypeString) (search.SortOrder, error) {
	switch typeStr {
	case "Asc":
		return search.SortOrder_ASC, nil
	case "Desc":
		return search.SortOrder_DESC, nil
	default:
		return 0, WrapError(fmt.Errorf("unexpected search index sort order string [Asc|Desc]: %s", typeStr))
	}
}

func ConvertSearchIndexSortModeString(typeStr SearchIndexSortModeString) (search.SortMode, error) {
	switch typeStr {
	case "Min":
		return search.SortMode_Min, nil
	case "Max":
		return search.SortMode_Max, nil
	case "Avg":
		return search.SortMode_Avg, nil
	default:
		return 0, WrapError(fmt.Errorf("unexpected search index sort mode string [Min|Max|Avg]: %s", typeStr))
	}
}

func (s *OtsService) WaitForSearchIndex(instance string, table string, indexName string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	id := ID(instance, table, indexName, SearchIndexTypeHolder)
	for {
		index, err := s.DescribeOtsSearchIndex(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if index.SyncStat != nil && status != Deleted {
			return nil
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, table, indexName, ProviderERROR)
		}
	}
}

func (s *OtsService) DeleteSearchIndex(instanceName string, tableName string, indexName string) error {
	request := &tablestore.DeleteSearchIndexRequest{
		TableName: tableName,
		IndexName: indexName,
	}
	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		var requestCli *tablestore.TableStoreClient
		raw, err := s.client.WithTableStoreClient(instanceName, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
			requestCli = tableStoreClient

			return tableStoreClient.DeleteSearchIndex(request)
		})
		defer func() {
			addDebug("DeleteSearchIndex", raw, requestCli, request)
		}()

		if err != nil {
			if IsExpectedErrors(err, OtsTableIsTemporarilyUnavailable) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	return err
}

// OtsRestApiPostWithRetry send POST request by CommonSDK(roa/restful) with retry.
// This method directly passes OpenAPI parameters such as product and version, without relying on SDK version upgrades.
// Retry policy: 3, 3+5, 3+5+5…, retry timeout: d.Timeout(schema.TimeoutCreate)
// product is openapi product code, version is openapi version, actionPath is restful openapi backend api path, requestBody is request body content
func OtsRestApiPostWithRetry(client *connectivity.AliyunClient, product string, version string, actionPath string, requestBody map[string]interface{}) (map[string]interface{}, error) {
	return invokeOtsRestApiWithRetry(client, product, version, actionPath, "POST", nil, nil, requestBody)
}

// OtsRestApiGetWithRetry send GET request by CommonSDK(roa/restful) with retry.
// This method directly passes OpenAPI parameters such as product and version, without relying on SDK version upgrades.
// Retry policy: 3, 3+5, 3+5+5…, retry timeout: d.Timeout(schema.TimeoutCreate)
// product is openapi product code, version is openapi version, actionPath is restful openapi backend api path, urlQuery is url param
func OtsRestApiGetWithRetry(client *connectivity.AliyunClient, product string, version string, actionPath string, urlQuery map[string]*string) (map[string]interface{}, error) {
	return invokeOtsRestApiWithRetry(client, product, version, actionPath, "GET", urlQuery, nil, nil)
}

func invokeOtsRestApiWithRetry(client *connectivity.AliyunClient, product string, version string, actionPath string, httpMethod string, urlQuery map[string]*string, headers map[string]*string, requestBody map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	otsClient, err := client.NewOtsRoaClient(product)
	if err != nil {
		return nil, WrapError(err)
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(20*time.Minute, func() *resource.RetryError {
		response, err = otsClient.DoRequest(StringPointer(version), nil, StringPointer(httpMethod), StringPointer("AK"), StringPointer(actionPath), urlQuery, headers, requestBody, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, OtsTableIsTemporarilyUnavailable) ||
				IsExpectedErrors(err, OtsTunnelIsTemporarilyUnavailable) ||
				IsExpectedErrors(err, OtsSecondaryIndexIsTemporarilyUnavailable) ||
				IsExpectedErrors(err, OtsSearchIndexIsTemporarilyUnavailable) ||
				NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(actionPath, response, requestBody)
		return nil
	})

	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, product, actionPath, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

// ACLString2Slice aclPattern: A,B,C
func ACLString2Slice(aclStr string) (s []string) {
	if aclStr == "" {
		return s
	}
	return strings.Split(aclStr, ",")
}

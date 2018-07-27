package alicloud

import (
	"fmt"
	"strings"

	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
)

func (client *AliyunClient) buildTableClient(instanceName string) *tablestore.TableStoreClient {
	endpoint := LoadEndpoint(client.RegionId, OTSCode)
	if endpoint == "" {
		endpoint = fmt.Sprintf("%s.%s.ots.aliyuncs.com", instanceName, client.RegionId)
	}
	if !strings.HasPrefix(endpoint, string(Https)) && !strings.HasPrefix(endpoint, string(Http)) {
		endpoint = fmt.Sprintf("%s://%s", Https, endpoint)
	}
	return tablestore.NewClient(endpoint, instanceName, client.AccessKey, client.SecretKey)
}

func getPrimaryKeyType(primaryKeyType string) tablestore.PrimaryKeyType {
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

func (client *AliyunClient) DescribeOtsTable(instanceName, tableName string) (table *tablestore.DescribeTableResponse, err error) {
	describeTableReq := new(tablestore.DescribeTableRequest)
	describeTableReq.TableName = tableName

	table, err = client.buildTableClient(instanceName).DescribeTable(describeTableReq)
	if err != nil {
		if strings.HasPrefix(err.Error(), OTSObjectNotExist) {
			err = GetNotFoundErrorFromString(GetNotFoundMessage("OTS Table", tableName))
		}
		return
	}
	if table == nil || table.TableMeta == nil || table.TableMeta.TableName != tableName {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("OTS Table", tableName))
	}
	return
}

func (client *AliyunClient) DeleteOtsTable(instanceName, tableName string) (bool, error) {

	deleteReq := new(tablestore.DeleteTableRequest)
	deleteReq.TableName = tableName
	if _, err := client.buildTableClient(instanceName).DeleteTable(deleteReq); err != nil {
		if NotFoundError(err) {
			return true, nil
		}
		return false, err
	}

	describ, err := client.DescribeOtsTable(instanceName, tableName)

	if err != nil {
		if NotFoundError(err) {
			return true, nil
		}
		return false, err
	}

	if describ.TableMeta != nil {
		return false, err
	}

	return true, err
}

// Convert tablestore.PrimaryKeyType to PrimaryKeyTypeString
func convertPrimaryKeyType(t tablestore.PrimaryKeyType) PrimaryKeyTypeString {
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

func (client *AliyunClient) DescribeOtsInstance(name string) (inst ots.InstanceInfo, err error) {
	req := ots.CreateGetInstanceRequest()
	req.InstanceName = name
	req.Method = "GET"
	resp, err := client.otsconn.GetInstance(req)

	// OTS instance not found error code is "NotFound"
	if err != nil {
		return
	}

	if resp == nil || resp.InstanceInfo.InstanceName != name {
		return inst, GetNotFoundErrorFromString(GetNotFoundMessage("OTS Instance", name))
	}
	return resp.InstanceInfo, nil
}

func (client *AliyunClient) DescribeOtsInstanceVpc(name string) (inst ots.VpcInfo, err error) {
	req := ots.CreateListVpcInfoByInstanceRequest()
	req.Method = "GET"
	req.InstanceName = name
	resp, err := client.otsconn.ListVpcInfoByInstance(req)
	if err != nil {
		return inst, err
	}
	if resp == nil || resp.TotalCount < 1 {
		return inst, GetNotFoundErrorFromString(GetNotFoundMessage("OTS Instance VPC", name))
	}
	return resp.VpcInfos.VpcInfo[0], nil
}

func (client *AliyunClient) WaitForOtsInstance(name string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		inst, err := client.DescribeOtsInstance(name)
		if err != nil {
			return err
		}

		if inst.Status == convertOtsInstanceStatus(status) {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("OTS Instance", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

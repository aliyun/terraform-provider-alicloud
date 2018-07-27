package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
)

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

func describeOtsTable(tableName string, meta interface{}) (*tablestore.DescribeTableResponse, error) {
	client := meta.(*AliyunClient).otsconn

	describeTableReq := new(tablestore.DescribeTableRequest)
	describeTableReq.TableName = tableName

	return client.DescribeTable(describeTableReq)
}

func deleteOtsTable(tableName string, meta interface{}) (bool, error) {
	client := meta.(*AliyunClient).otsconn

	deleteReq := new(tablestore.DeleteTableRequest)
	deleteReq.TableName = tableName
	_, err := client.DeleteTable(deleteReq)

	describ, _ := describeOtsTable(tableName, meta)

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
	resp, err := client.otsconnnew.GetInstance(req)

	// OTS instance not found error code is "NotFound"
	if err != nil {
		return
	}
	//OTS instance status: 3-deleting, 4-deleted
	if resp == nil || resp.InstanceInfo.Status == 4 {
		return inst, GetNotFoundErrorFromString(GetNotFoundMessage("OTS Instance", name))
	}
	return resp.InstanceInfo, nil
}

func (client *AliyunClient) DescribeOtsInstanceVpc(name string) (inst ots.VpcInfo, err error) {
	req := ots.CreateListVpcInfoByInstanceRequest()
	req.Method = "GET"
	req.InstanceName = name
	resp, err := client.otsconnnew.ListVpcInfoByInstance(req)
	if err != nil {
		return inst, err
	}
	if resp == nil || resp.TotalCount < 1 {
		return inst, GetNotFoundErrorFromString(GetNotFoundMessage("OTS Instance VPC", name))
	}
	return resp.VpcInfos.VpcInfo[0], nil
}

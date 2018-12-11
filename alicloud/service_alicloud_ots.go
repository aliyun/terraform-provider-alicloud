package alicloud

import (
	"strings"

	"time"

	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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

func (s *OtsService) DescribeOtsTable(instanceName, tableName string) (table *tablestore.DescribeTableResponse, err error) {
	describeTableReq := new(tablestore.DescribeTableRequest)
	describeTableReq.TableName = tableName

	raw, err := s.client.WithTableStoreClient(instanceName, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
		return tableStoreClient.DescribeTable(describeTableReq)
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), OTSObjectNotExist) {
			err = GetNotFoundErrorFromString(GetNotFoundMessage("OTS Table", tableName))
		}
		return
	}
	table, _ = raw.(*tablestore.DescribeTableResponse)
	if table == nil || table.TableMeta == nil || table.TableMeta.TableName != tableName {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("OTS Table", tableName))
	}
	return
}

func (s *OtsService) DeleteOtsTable(instanceName, tableName string) (bool, error) {

	deleteReq := new(tablestore.DeleteTableRequest)
	deleteReq.TableName = tableName
	_, err := s.client.WithTableStoreClient(instanceName, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
		return tableStoreClient.DeleteTable(deleteReq)
	})
	if err != nil {
		if NotFoundError(err) {
			return true, nil
		}
		return false, err
	}

	describ, err := s.DescribeOtsTable(instanceName, tableName)

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

func (s *OtsService) DescribeOtsInstance(name string) (inst ots.InstanceInfo, err error) {
	req := ots.CreateGetInstanceRequest()
	req.InstanceName = name
	req.Method = "GET"
	raw, err := s.client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
		return otsClient.GetInstance(req)
	})

	// OTS instance not found error code is "NotFound"
	if err != nil {
		return
	}
	resp, _ := raw.(*ots.GetInstanceResponse)
	if resp == nil || resp.InstanceInfo.InstanceName != name {
		return inst, GetNotFoundErrorFromString(GetNotFoundMessage("OTS Instance", name))
	}
	return resp.InstanceInfo, nil
}

func (s *OtsService) DescribeOtsInstanceVpc(name string) (inst ots.VpcInfo, err error) {
	req := ots.CreateListVpcInfoByInstanceRequest()
	req.Method = "GET"
	req.InstanceName = name
	raw, err := s.client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
		return otsClient.ListVpcInfoByInstance(req)
	})
	if err != nil {
		return inst, err
	}
	resp, _ := raw.(*ots.ListVpcInfoByInstanceResponse)
	if resp == nil || resp.TotalCount < 1 {
		return inst, GetNotFoundErrorFromString(GetNotFoundMessage("OTS Instance VPC", name))
	}
	return resp.VpcInfos.VpcInfo[0], nil
}

func (s *OtsService) WaitForOtsInstance(name string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		inst, err := s.DescribeOtsInstance(name)
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

func (s *OtsService) DescribeOtsInstanceTypes() (types []string, err error) {
	req := ots.CreateListClusterTypeRequest()
	req.Method = requests.GET
	raw, err := s.client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
		return otsClient.ListClusterType(req)
	})
	if err != nil {
		err = fmt.Errorf("Failed to list instance type with error: %#v", err)
		return
	}
	resp, _ := raw.(*ots.ListClusterTypeResponse)
	if resp != nil {
		return resp.ClusterTypeInfos.ClusterType, nil
	}
	return
}

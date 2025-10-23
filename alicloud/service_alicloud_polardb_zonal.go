// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type PolarDbServiceV2 struct {
	client *connectivity.AliyunClient
}

// DBNode is a nested struct in polardb response
type DBNode struct {
	MaxIOPS          int    `json:"MaxIOPS" xml:"MaxIOPS"`
	DBNodeClass      string `json:"DBNodeClass" xml:"DBNodeClass"`
	FailoverPriority int    `json:"FailoverPriority" xml:"FailoverPriority"`
	DBNodeRole       string `json:"DBNodeRole" xml:"DBNodeRole"`
	DBNodeStatus     string `json:"DBNodeStatus" xml:"DBNodeStatus"`
	MasterId         string `json:"MasterId" xml:"MasterId"`
	CreationTime     string `json:"CreationTime" xml:"CreationTime"`
	HotReplicaMode   string `json:"HotReplicaMode" xml:"HotReplicaMode"`
	ServerlessType   string `json:"ServerlessType" xml:"ServerlessType"`
	Serverless       string `json:"Serverless" xml:"Serverless"`
	MaxConnections   int    `json:"MaxConnections" xml:"MaxConnections"`
	AddedCpuCores    string `json:"AddedCpuCores" xml:"AddedCpuCores"`
	RegionId         string `json:"RegionId" xml:"RegionId"`
	ZoneId           string `json:"ZoneId" xml:"ZoneId"`
	ServerWeight     string `json:"ServerWeight" xml:"ServerWeight"`
	DBNodeId         string `json:"DBNodeId" xml:"DBNodeId"`
	SccMode          string `json:"SccMode" xml:"SccMode"`
	ImciSwitch       string `json:"ImciSwitch" xml:"ImciSwitch"`
}

type DescribeDBClusterAttributeResponse struct {
	*responses.BaseResponse
	DeletionLock              int      `json:"DeletionLock" xml:"DeletionLock"`
	Category                  string   `json:"Category" xml:"Category"`
	ResourceGroupId           string   `json:"ResourceGroupId" xml:"ResourceGroupId"`
	DataLevel1BackupChainSize int64    `json:"DataLevel1BackupChainSize" xml:"DataLevel1BackupChainSize"`
	DBClusterId               string   `json:"DBClusterId" xml:"DBClusterId"`
	DBType                    string   `json:"DBType" xml:"DBType"`
	DBClusterNetworkType      string   `json:"DBClusterNetworkType" xml:"DBClusterNetworkType"`
	IsLatestVersion           bool     `json:"IsLatestVersion" xml:"IsLatestVersion"`
	HasCompleteStandbyRes     bool     `json:"HasCompleteStandbyRes" xml:"HasCompleteStandbyRes"`
	HotStandbyCluster         string   `json:"HotStandbyCluster" xml:"HotStandbyCluster"`
	DataSyncMode              string   `json:"DataSyncMode" xml:"DataSyncMode"`
	StandbyHAMode             string   `json:"StandbyHAMode" xml:"StandbyHAMode"`
	CompressStorageMode       string   `json:"CompressStorageMode" xml:"CompressStorageMode"`
	StorageMax                int64    `json:"StorageMax" xml:"StorageMax"`
	DBVersion                 string   `json:"DBVersion" xml:"DBVersion"`
	ZoneIds                   string   `json:"ZoneIds" xml:"ZoneIds"`
	MaintainTime              string   `json:"MaintainTime" xml:"MaintainTime"`
	Engine                    string   `json:"Engine" xml:"Engine"`
	RequestId                 string   `json:"RequestId" xml:"RequestId"`
	VPCId                     string   `json:"VPCId" xml:"VPCId"`
	DBClusterStatus           string   `json:"DBClusterStatus" xml:"DBClusterStatus"`
	VSwitchId                 string   `json:"VSwitchId" xml:"VSwitchId"`
	DBClusterDescription      string   `json:"DBClusterDescription" xml:"DBClusterDescription"`
	Expired                   string   `json:"Expired" xml:"Expired"`
	PayType                   string   `json:"PayType" xml:"PayType"`
	StoragePayType            string   `json:"StoragePayType" xml:"StoragePayType"`
	LockMode                  string   `json:"LockMode" xml:"LockMode"`
	StorageUsed               int64    `json:"StorageUsed" xml:"StorageUsed"`
	CompressStorageUsed       int64    `json:"CompressStorageUsed" xml:"CompressStorageUsed"`
	StorageSpace              int64    `json:"StorageSpace" xml:"StorageSpace"`
	DBVersionStatus           string   `json:"DBVersionStatus" xml:"DBVersionStatus"`
	CreationTime              string   `json:"CreationTime" xml:"CreationTime"`
	SQLSize                   int64    `json:"SQLSize" xml:"SQLSize"`
	InodeTotal                int64    `json:"InodeTotal" xml:"InodeTotal"`
	InodeUsed                 int64    `json:"InodeUsed" xml:"InodeUsed"`
	BlktagTotal               int64    `json:"BlktagTotal" xml:"BlktagTotal"`
	BlktagUsed                int64    `json:"BlktagUsed" xml:"BlktagUsed"`
	RegionId                  string   `json:"RegionId" xml:"RegionId"`
	ExpireTime                string   `json:"ExpireTime" xml:"ExpireTime"`
	SubCategory               string   `json:"SubCategory" xml:"SubCategory"`
	DeployUnit                string   `json:"DeployUnit" xml:"DeployUnit"`
	IsProxyLatestVersion      bool     `json:"IsProxyLatestVersion" xml:"IsProxyLatestVersion"`
	StorageType               string   `json:"StorageType" xml:"StorageType"`
	ServerlessType            string   `json:"ServerlessType" xml:"ServerlessType"`
	StrictConsistency         string   `json:"StrictConsistency" xml:"StrictConsistency"`
	ProxyCpuCores             string   `json:"ProxyCpuCores" xml:"ProxyCpuCores"`
	ProxyStandardCpuCores     string   `json:"ProxyStandardCpuCores" xml:"ProxyStandardCpuCores"`
	ProxyType                 string   `json:"ProxyType" xml:"ProxyType"`
	ProxyStatus               string   `json:"ProxyStatus" xml:"ProxyStatus"`
	FeatureHTAPSupported      string   `json:"FeatureHTAPSupported" xml:"FeatureHTAPSupported"`
	ProxyServerlessType       string   `json:"ProxyServerlessType" xml:"ProxyServerlessType"`
	Architecture              string   `json:"Architecture" xml:"Architecture"`
	AiType                    string   `json:"AiType" xml:"AiType"`
	DBNodes                   []DBNode `json:"DBNodes" xml:"DBNodes"`
	Tags                      []Tag    `json:"Tags" xml:"Tags"`
	DBClusterClass            string   `json:"DBClusterClass" xml:"DBClusterClass"`
}

// DescribePolarDbZonalCluster <<< Encapsulated get interface for PolarDb ZonalCluster.

func (s *PolarDbServiceV2) DescribePolarDbZonalCluster(id string) (object *DescribeDBClusterAttributeResponse, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBClusterId"] = id

	action := "DescribeDbClusterAttributeZonal"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("ZonalCluster", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	var responseObj *DescribeDBClusterAttributeResponse
	responseByte, err := json.Marshal(response)
	if err := json.Unmarshal(responseByte, &responseObj); err != nil {
		return object, err
	}

	return responseObj, nil
}
func (s *PolarDbServiceV2) DescribeZonalClusterDescribeDBClusterVersionZonal(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DBClusterId"] = id

	action := "DescribeDBClusterVersionZonal"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcGet("polardb", "2017-08-01", action, query, request)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("ZonalCluster", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}
func (s *PolarDbServiceV2) DescribeZonalClusterDescribeAutoRenewAttribute(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBClusterIds"] = id
	request["RegionId"] = client.RegionId
	request["CloudProvider"] = "ENS"
	action := "DescribeAutoRenewAttribute"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Items.AutoRenewAttribute[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Items.AutoRenewAttribute[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("ZonalCluster", id), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *PolarDbServiceV2) PolarDbZonalClusterStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePolarDbZonalCluster(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.DBClusterStatus == failState {
				return object, object.DBClusterStatus, WrapError(Error(FailedToReachTargetStatus, object.DBClusterStatus))
			}
		}
		return object, object.DBClusterStatus, nil
	}
}

func (s *PolarDbServiceV2) WaitForPolarDBPayType(id string, status string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		clusters, err := s.DescribePolarDbZonalCluster(id)
		if err != nil {
			return WrapError(err)
		}
		v, err := jsonpath.Get("PayType", clusters)
		currentStatus := fmt.Sprint(v)
		if currentStatus == status {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, clusters, status, ProviderERROR)
		}
		time.Sleep(DefaultIntervalMedium * time.Second)
	}
	return nil
}

func (s *PolarDbServiceV2) ModifyDBNodesClass(regionId, DBClusterId, modifyType string, nodes interface{}) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBClusterId"] = DBClusterId
	request["RegionId"] = regionId
	request["CloudProvider"] = "ENS"
	request["ModifyType"] = modifyType
	request["DBNode"] = nodes
	action := "ModifyDBNodesClass"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("ZonalCluster", DBClusterId), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, DBClusterId, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *PolarDbServiceV2) DeleteDBNodes(regionId, DBClusterId string, nodeId []string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBClusterId"] = DBClusterId
	request["RegionId"] = regionId
	request["CloudProvider"] = "ENS"
	request["DBNodeId"] = nodeId
	action := "DeleteDBNodes"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("ZonalCluster", DBClusterId), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, DBClusterId, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *PolarDbServiceV2) CreateDBNodes(regionId, DBClusterId, ImciSwitch string, node interface{}) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBClusterId"] = DBClusterId
	request["RegionId"] = regionId
	request["CloudProvider"] = "ENS"
	request["DBNode"] = node
	request["ImciSwitch"] = ImciSwitch
	action := "CreateDBNodes"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("ZonalCluster", DBClusterId), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, DBClusterId, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *PolarDbServiceV2) ModifyAutoRenewAttribute(regionId, DBClusterId, renewalStatus, duration, periodUnit string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBClusterIds"] = DBClusterId
	request["RegionId"] = regionId
	request["CloudProvider"] = "ENS"
	request["RenewalStatus"] = renewalStatus
	request["Duration"] = duration
	request["PeriodUnit"] = periodUnit
	action := "ModifyAutoRenewAttribute"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, err, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("ZonalCluster", DBClusterId), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, DBClusterId, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *PolarDbServiceV2) WaitForPolarDBDeleted(id string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		_, err := s.DescribePolarDbZonalCluster(id)
		if err != nil {
			if NotFoundError(err) {
				time.Sleep(60 * time.Second)
				return nil
			}
			return WrapError(err)
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

type DBEndpoint struct {
	NodeWithRoles         string    `json:"NodeWithRoles" xml:"NodeWithRoles"`
	Nodes                 string    `json:"Nodes" xml:"Nodes"`
	ReadWriteMode         string    `json:"ReadWriteMode" xml:"ReadWriteMode"`
	DBEndpointId          string    `json:"DBEndpointId" xml:"DBEndpointId"`
	EndpointConfig        string    `json:"EndpointConfig" xml:"EndpointConfig"`
	DBEndpointDescription string    `json:"DBEndpointDescription" xml:"DBEndpointDescription"`
	EndpointType          string    `json:"EndpointType" xml:"EndpointType"`
	AutoAddNewNodes       string    `json:"AutoAddNewNodes" xml:"AutoAddNewNodes"`
	DBClusterId           string    `json:"DBClusterId" xml:"DBClusterId"`
	AddressItems          []Address `json:"AddressItems" xml:"AddressItems"`
}

type Address struct {
	PrivateZoneConnectionString string `json:"PrivateZoneConnectionString" xml:"PrivateZoneConnectionString"`
	VpcInstanceId               string `json:"VpcInstanceId" xml:"VpcInstanceId"`
	VPCId                       string `json:"VPCId" xml:"VPCId"`
	Port                        string `json:"Port" xml:"Port"`
	VSwitchId                   string `json:"VSwitchId" xml:"VSwitchId"`
	SSLEnabled                  string `json:"SSLEnabled" xml:"SSLEnabled"`
	ConnectionString            string `json:"ConnectionString" xml:"ConnectionString"`
	IPAddress                   string `json:"IPAddress" xml:"IPAddress"`
	NetType                     string `json:"NetType" xml:"NetType"`
}

func (s *PolarDbServiceV2) DescribeDBClusterEndpointsZonal(id string) (object *DBEndpoint, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}
	dbClusterId := parts[0]
	dbEndpointId := parts[1]

	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = s.client.RegionId
	request["DBClusterId"] = dbClusterId
	request["DBEndpointId"] = dbEndpointId
	action := "DescribeDBClusterEndpointsZonal"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		addDebug(action, response, request)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("ZonalClusterEndpoint", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	if response["Items"] != nil && len(response["Items"].([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("ZonalClusterEndpoint", id), NotFoundMsg, response)
	}

	var responseObj *DBEndpoint
	responseByte, err := json.Marshal(response["Items"].([]interface{})[0])
	if err := json.Unmarshal(responseByte, &responseObj); err != nil {
		return object, err
	}

	return responseObj, nil
}

func (s *PolarDbServiceV2) DescribeDBClusterEndpointsZonalList(regionId, DBClusterId string) (object *[]DBEndpoint, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = regionId
	request["DBClusterId"] = DBClusterId
	action := "DescribeDBClusterEndpointsZonal"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		addDebug(action, response, request)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	var responseObj *[]DBEndpoint
	responseByte, err := json.Marshal(response["Items"])
	if err := json.Unmarshal(responseByte, &responseObj); err != nil {
		return object, err
	}

	return responseObj, nil
}

type CreateDBClusterEndpointRequest struct {
	*polardb.CreateDBClusterEndpointRequest
	VPCId     string `position:"Query" name:"VPCId"`
	VSwitchId string `position:"Query" name:"VSwitchId"`
}

func CreateCreateDBClusterEndpointRequest() (request *CreateDBClusterEndpointRequest) {
	request = &CreateDBClusterEndpointRequest{
		CreateDBClusterEndpointRequest: &polardb.CreateDBClusterEndpointRequest{RpcRequest: &requests.RpcRequest{}},
	}
	request.InitWithApiInfo("polardb", "2017-08-01", "CreateDBClusterEndpoint", "polardb", "openAPI")
	request.Method = requests.POST
	return
}

func (s *PolarDbServiceV2) CreateDBClusterEndpointZonal(requestObj *CreateDBClusterEndpointRequest) (err error) {
	client := s.client
	var request map[string]interface{}
	var query map[string]interface{}
	query = make(map[string]interface{})
	action := "CreateDBClusterEndpointZonal"

	data, _ := json.Marshal(*requestObj)
	if err = json.Unmarshal(data, &request); err != nil {
		return err
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		addDebug(action, err, request)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if IsExpectedErrorsMessage(err, "Current DB instance state does not support this operation") {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *PolarDbServiceV2) ModifyDBClusterEndpointZonal(requestObj *polardb.ModifyDBClusterEndpointRequest) (err error) {
	client := s.client
	var request map[string]interface{}
	var query map[string]interface{}
	query = make(map[string]interface{})
	action := "ModifyDBClusterEndpointZonal"

	data, _ := json.Marshal(*requestObj)
	if err = json.Unmarshal(data, &request); err != nil {
		return err
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err := client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		log.Printf("ModifyDBClusterEndpointZonal response %s %v %v", requestObj.DBEndpointId, request, response)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if IsExpectedErrorsMessage(err, "Current endpoint status does not support this operation") {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *PolarDbServiceV2) DeleteDBClusterEndpointZonal(regionId, DBClusterId, endpointId string) (err error) {
	client := s.client
	var request map[string]interface{}
	var query map[string]interface{}
	query = make(map[string]interface{})
	request = make(map[string]interface{})
	request["RegionId"] = regionId
	request["DBClusterId"] = DBClusterId
	request["DBEndpointId"] = endpointId
	action := "DeleteDBClusterEndpointZonal"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		addDebug(action, err, request)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *PolarDbServiceV2) WaitForPolarDBEndpoints(d *schema.ResourceData, status Status, endpointIds *schema.Set, endpointType string, timeout int) (string, error) {
	var dbEndpointId string
	if d.Id() != "" {
		parts, err := ParseResourceId(d.Id(), 2)
		if err != nil {
			return "", WrapError(err)
		}
		dbEndpointId = parts[1]
	}
	dbClusterId := d.Get("db_cluster_id").(string)

	newEndpoint := make(map[string]string)
	newEndpoint["endpoint_type"] = endpointType

	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		endpoints, err := s.DescribeDBClusterEndpointsZonalList(s.client.RegionId, dbClusterId)
		if err != nil {
			return "", WrapError(err)
		}

		deleted := true
		completed := false
		for _, value := range *endpoints {
			if status == Deleted && dbEndpointId == value.DBEndpointId {
				deleted = false
				continue
			}

			if endpointIds != nil && !endpointIds.Contains(value.DBEndpointId) &&
				value.EndpointType == endpointType && value.EndpointConfig != "" {
				return value.DBEndpointId, nil
			}

			if status == Active && dbEndpointId == value.DBEndpointId {
				if value.EndpointConfig != "" {
					completed = true
				}
			}
		}
		if status == Deleted && deleted {
			return "", nil
		}

		if status == Active && completed {
			return "", nil
		}

		if time.Now().After(deadline) {
			return "", WrapErrorf(err, WaitTimeoutMsg, dbClusterId, GetFunc(1), timeout, endpoints, newEndpoint, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *PolarDbServiceV2) CreateAccount(requestObj *polardb.CreateAccountRequest) (err error) {
	client := s.client
	action := "CreateAccountZonal"

	var query map[string]interface{}
	query = make(map[string]interface{})

	var request map[string]interface{}
	data, _ := json.Marshal(*requestObj)
	if err = json.Unmarshal(data, &request); err != nil {
		return err
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		addDebug(action, err, request)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if IsExpectedErrorsMessage(err, "Current DB instance state does not support this operation") {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *PolarDbServiceV2) DescribePolarDBAccount(id string) (ds *polardb.DBAccount, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}

	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBClusterId"] = parts[0]
	request["AccountName"] = parts[1]
	action := "DescribeAccountsZonal"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		addDebug(action, response, request)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if IsExpectedErrorsMessage(err, "Current DB instance state does not support this operation") {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if IsExpectedErrorsMessage(err, "The DBClusterId provided does not exist in our records") {
			return nil, WrapErrorf(NotFoundErr("ZonalCluster", id), NotFoundMsg, response)
		}
		return nil, err
	}

	var responseObj []polardb.DBAccount
	responseByte, err := json.Marshal(response["Accounts"])
	if err := json.Unmarshal(responseByte, &responseObj); err != nil {
		return nil, err
	}

	if len(responseObj) == 0 {
		return nil, nil
	}

	return &responseObj[0], nil
}

func (s *PolarDbServiceV2) WaitForPolarDBAccount(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribePolarDBAccount(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if status == Deleted && object == nil {
			return nil
		}

		if object.AccountStatus == string(status) {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.AccountStatus, status, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *PolarDbServiceV2) DeleteAccount(dbClusterId, accountName string) (err error) {
	client := s.client
	action := "DeleteAccountZonal"

	var query map[string]interface{}
	query = make(map[string]interface{})

	var request map[string]interface{}
	request = make(map[string]interface{})
	request["DBClusterId"] = dbClusterId
	request["AccountName"] = accountName

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		addDebug(action, err, request)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if IsExpectedErrorsMessage(err, "Current DB instance state does not support this operation") {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *PolarDbServiceV2) modifyAccountDescription(dbClusterId, accountName, accountDescription string) (err error) {
	client := s.client
	action := "ModifyAccountDescriptionZonal"

	var query map[string]interface{}
	query = make(map[string]interface{})

	var request map[string]interface{}
	request = make(map[string]interface{})
	request["DBClusterId"] = dbClusterId
	request["AccountName"] = accountName
	request["AccountDescription"] = accountDescription

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		addDebug(action, err, request)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *PolarDbServiceV2) modifyAccountPassword(dbClusterId, accountName, newAccountPassword string) (err error) {
	client := s.client
	action := "ModifyAccountPasswordZonal"

	var query map[string]interface{}
	query = make(map[string]interface{})

	var request map[string]interface{}
	request = make(map[string]interface{})
	request["DBClusterId"] = dbClusterId
	request["AccountName"] = accountName
	request["NewAccountPassword"] = newAccountPassword

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		addDebug(action, err, request)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *PolarDbServiceV2) GrantPolarDBAccountPrivilege(id, dbName string) error {
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return WrapError(err)
	}

	var query map[string]interface{}
	query = make(map[string]interface{})

	var request map[string]interface{}
	request = make(map[string]interface{})
	request["DBClusterId"] = parts[0]
	request["AccountName"] = parts[1]
	request["AccountPrivilege"] = parts[2]
	request["DBName"] = dbName

	client := s.client
	action := "GrantAccountPrivilegeZonal"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		addDebug(action, err, request)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err := s.WaitForPolarDBAccountPrivilege(id, dbName, Available, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}

	return nil
}

func (s *PolarDbServiceV2) RevokePolarDBAccountPrivilege(id, dbName string) error {
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return WrapError(err)
	}

	var query map[string]interface{}
	query = make(map[string]interface{})

	var request map[string]interface{}
	request = make(map[string]interface{})
	request["DBClusterId"] = parts[0]
	request["AccountName"] = parts[1]
	request["DBName"] = dbName

	client := s.client
	action := "RevokeAccountPrivilegeZonal"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		addDebug(action, err, request)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	if err := s.WaitForPolarDBAccountPrivilegeRevoked(id, dbName, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}

	return nil
}

func (s *PolarDbServiceV2) WaitForPolarDBAccountPrivilegeRevoked(id, dbName string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribePolarDBAccountPrivilege(id)
		if err != nil {
			return err
		}

		exist := false
		if object != nil {
			for _, dp := range object.DatabasePrivileges {
				if dp.DBName == dbName {
					exist = true
					break
				}
			}
		}

		if !exist {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, "", dbName, ProviderERROR)
		}

	}
	return nil
}

func (s *PolarDbServiceV2) WaitForPolarDBAccountPrivilege(id, dbName string, status Status, timeout int) error {
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribePolarDBAccountPrivilege(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if status == Deleted && object == nil {
			return nil
		}

		ready := false
		if object != nil {
			for _, dp := range object.DatabasePrivileges {
				if dp.DBName == dbName && dp.AccountPrivilege == parts[2] {
					ready = true
					break
				}
			}
		}
		if status == Deleted && !ready {
			break
		}
		if status != Deleted && ready {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, "", dbName, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *PolarDbServiceV2) DescribePolarDBAccountPrivilege(id string) (account *polardb.DBAccount, err error) {
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = WrapError(err)
		return
	}

	ds, err := s.DescribePolarDBAccount(parts[0] + ":" + parts[1])
	if err != nil {
		return nil, err
	}
	return ds, nil
}

type DescribeDatabasesResponse struct {
	*responses.BaseResponse
	PageRecordCount int        `json:"PageRecordCount" xml:"PageRecordCount"`
	RequestId       string     `json:"RequestId" xml:"RequestId"`
	PageNumber      int        `json:"PageNumber" xml:"PageNumber"`
	Databases       []Database `json:"Databases" xml:"Databases"`
}

type Database struct {
	DBDescription    string    `json:"DBDescription" xml:"DBDescription"`
	DBStatus         string    `json:"DBStatus" xml:"DBStatus"`
	DBName           string    `json:"DBName" xml:"DBName"`
	Engine           string    `json:"Engine" xml:"Engine"`
	MasterID         string    `json:"MasterID" xml:"MasterID"`
	CharacterSetName string    `json:"CharacterSetName" xml:"CharacterSetName"`
	Accounts         []Account `json:"Accounts" xml:"Accounts"`
}

type Account struct {
	PrivilegeStatus  string `json:"PrivilegeStatus" xml:"PrivilegeStatus"`
	AccountStatus    string `json:"AccountStatus" xml:"AccountStatus"`
	AccountPrivilege string `json:"AccountPrivilege" xml:"AccountPrivilege"`
	AccountName      string `json:"AccountName" xml:"AccountName"`
}

func (s *PolarDbServiceV2) DescribePolarDBDatabase(id string) (ds *Database, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}

	var query map[string]interface{}
	query = make(map[string]interface{})

	var request map[string]interface{}
	request = make(map[string]interface{})
	request["DBClusterId"] = parts[0]
	request["DBName"] = parts[1]
	client := s.client
	action := "DescribeDatabasesZonal"

	var raw interface{}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		raw, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		addDebug(action, raw, err, request)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if IsExpectedErrorsMessage(err, "The DBClusterId provided does not exist in our records") {
			return nil, WrapErrorf(NotFoundErr("ZonalCluster", id), NotFoundMsg, raw)
		}
		return nil, err
	}

	var responseObj *DescribeDatabasesResponse
	responseByte, err := json.Marshal(raw)
	if err := json.Unmarshal(responseByte, &responseObj); err != nil {
		return nil, err
	}
	if len(responseObj.Databases) < 1 {
		return nil, WrapErrorf(NotFoundErr("DBDatabase", parts[1]), NotFoundMsg, ProviderERROR)
	}
	ds = &responseObj.Databases[0]
	return ds, nil
}

func (s *PolarDbServiceV2) WaitForPolarDBDatabase(id string, status Status, timeout int) error {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribePolarDBDatabase(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
				if status == Running {
					continue
				}
			}
			return WrapError(err)
		}
		if status != Deleted && object != nil && object.DBName == parts[1] {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.DBName, parts[1], ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *PolarDbServiceV2) CreateDatabaseZonal(requestObj *polardb.CreateDatabaseRequest) error {
	var query map[string]interface{}
	query = make(map[string]interface{})

	var request map[string]interface{}
	data, _ := json.Marshal(*requestObj)
	if err := json.Unmarshal(data, &request); err != nil {
		return err
	}
	client := s.client
	action := "CreateDatabaseZonal"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err := resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, err := client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		addDebug(action, err, request)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if IsExpectedErrorsMessage(err, "Current DB instance state does not support this operation") {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *PolarDbServiceV2) ModifyDBDescriptionZonal(dbClusterId, dbName, dbDescription string) error {
	var query map[string]interface{}
	query = make(map[string]interface{})

	var request map[string]interface{}
	request = make(map[string]interface{})
	request["DBClusterId"] = dbClusterId
	request["DBName"] = dbName
	request["DBDescription"] = dbDescription

	client := s.client
	action := "ModifyDBDescriptionZonal"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err := resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, err := client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		addDebug(action, err, request)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *PolarDbServiceV2) DeleteDatabaseZonal(dbClusterId, dbName string) error {
	var query map[string]interface{}
	query = make(map[string]interface{})

	var request map[string]interface{}
	request = make(map[string]interface{})
	request["DBClusterId"] = dbClusterId
	request["DBName"] = dbName

	client := s.client
	action := "DeleteDatabaseZonal"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err := resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, err := client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		addDebug(action, err, request)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if IsExpectedErrorsMessage(err, "Current DB instance state does not support this operation") {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *PolarDbServiceV2) FailoverDBClusterZonal(dbClusterId, dbnodeId string) error {
	var query map[string]interface{}
	query = make(map[string]interface{})

	var request map[string]interface{}
	request = make(map[string]interface{})
	request["DBClusterId"] = dbClusterId
	request["TargetDBNodeId"] = dbnodeId

	client := s.client
	action := "FailoverDBClusterZonal"

	wait := incrementalWait(3*time.Second, 30*time.Second)
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		addDebug(action, err, request)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if IsExpectedErrorsMessage(err, "Current DB instance HA status") {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *PolarDbServiceV2) ModifyDBClusterDescriptionZonal(dbClusterId, descriptions string) error {
	var query map[string]interface{}
	query = make(map[string]interface{})

	var request map[string]interface{}
	request = make(map[string]interface{})
	request["DBClusterId"] = dbClusterId
	request["DBClusterDescription"] = descriptions

	client := s.client
	action := "ModifyDBClusterDescriptionZonal"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err := resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, err := client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		addDebug(action, err, request)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *PolarDbServiceV2) UpgradeDBClusterVersionZonal(dbClusterId string) error {
	var query map[string]interface{}
	query = make(map[string]interface{})

	var request map[string]interface{}
	request = make(map[string]interface{})
	request["DBClusterId"] = dbClusterId
	request["UpgradeType"] = "ALL"
	request["UpgradePolicy"] = "hot"

	client := s.client
	action := "UpgradeDBClusterVersionZonal"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err := resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, err := client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		addDebug(action, err, request)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func IsExpectedErrorsMessage(err error, messageKeyWord string) bool {
	if err == nil {
		return false
	}

	if e, ok := err.(*tea.SDKError); ok {
		if strings.Contains(*e.Message, messageKeyWord) {
			return true
		}
		return false
	}

	return false
}

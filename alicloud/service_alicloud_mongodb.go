package alicloud

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/denverdino/aliyungo/common"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type MongoDBInstance struct {
	ChargeType            string `json:"ChargeType"`
	CreationTime          string `json:"CreationTime"`
	DBInstanceClass       string `json:"DBInstanceClass"`
	DBInstanceDescription string `json:"DBInstanceDescription"`
	DBInstanceID          string `json:"DBInstanceId"`
	DBInstanceStatus      string `json:"DBInstanceStatus"`
	DBInstanceStorage     int    `json:"DBInstanceStorage"`
	DBInstanceType        string `json:"DBInstanceType"`
	Engine                string `json:"Engine"`
	EngineVersion         string `json:"EngineVersion"`
	ExpireTime            string `json:"ExpireTime"`
	LockMode              string `json:"LockMode"`
	NetworkType           string `json:"NetworkType"`
	RegionID              string `json:"RegionId"`
	ReplicationFactor     string `json:"ReplicationFactor"`
	ZoneID                string `json:"ZoneId"`
	BackupID              string `json:"BackupId"`
}

type ItemsInDescribeMongoDBInstances struct {
	DBInstances []MongoDBInstance `json:"DBInstance"`
}

type DescribeMongoDBInstancesResponse struct {
	PageNumber int                             `json:"PageNumber"`
	PageSize   int                             `json:"PageSize"`
	RequestID  string                          `json:"RequestId"`
	TotalCount int                             `json:"TotalCount"`
	Items      ItemsInDescribeMongoDBInstances `json:"DBInstances"`
}

type DescribeDBInstanceAttributeResponse struct {
	Items ItemsInDescribeMongoDBInstances `json:"DBInstances"`
}

type CreateMongoDBInstanceResponse struct {
	DBInstanceId string `json:"DBInstanceId"`
	OrderId      string `json:"OrderId"`
}

type DescribeMongoDBSecurityIpsResponse struct {
	SecurityIps string                                    `json:"SecurityIps"`
	Items       ItemsInDescribeMongoDBSecurityIpsResponse `json"SecurityIpGroups"`
}

type ItemsInDescribeMongoDBSecurityIpsResponse struct {
	SecurityIpGroups []SecurityMongoDBIpGroup `json"SecurityIpGroup"`
}

type SecurityMongoDBIpGroup struct {
	SecurityIpGroupName string `json:"SecurityIpGroupName"`
	SecurityIpList      string `json:"SecurityIpList"`
	SecurityIpAttribute string `json:"SecurityIpAttribute"`
}

type DescribeMongoDBBackupPolicyResponse struct {
	BackupRetentionPeriod string `json:"BackupRetentionPeriod"`
	PreferredBackupTime   string `json:"PreferredBackupTime"`
	PreferredBackupPeriod string `json:"PreferredBackupPeriod"`
}

type ReplicaSets struct {
	Items ReplicaSet `json:"ReplicaSets"`
}

type ReplicaSet struct {
	ReplicaSets []ReplicaSetRole `json:"ReplicaSet"`
}

type ReplicaSetRole struct {
	ReplicaSetRole   string `json:"ReplicaSetRole"`
	ConnectionDomain string `json:"ConnectionDomain`
	ConnectionPort   string `json:"ConnectionPort"`
	NetworkType      string `json:"NetworkType"`
}

type MongoDBService struct {
	client *connectivity.AliyunClient
}

func (client *MongoDBService) DescribeMongoDBInstances(request *requests.CommonRequest, aliyunClient *connectivity.AliyunClient) (response *DescribeMongoDBInstancesResponse, err error) {
	resp, err := aliyunClient.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		request.Version = ApiVersion20151201
		request.ApiName = "DescribeDBInstances"

		resp, err := ecsClient.ProcessCommonRequest(request)
		if err != nil {
			return nil, err
		}
		response = new(DescribeMongoDBInstancesResponse)
		err = json.Unmarshal(resp.BaseResponse.GetHttpContentBytes(), &response)
		return response, err
	})
	if err != nil {
		return nil, err
	}
	return resp.(*DescribeMongoDBInstancesResponse), err
}

func (client *MongoDBService) CreateMongoDBInstance(request *requests.CommonRequest, aliyunClient *connectivity.AliyunClient) (response *CreateMongoDBInstanceResponse, err error) {
	resp, err := aliyunClient.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		request.Version = ApiVersion20151201
		request.ApiName = "CreateDBInstance"
		resp, err := ecsClient.ProcessCommonRequest(request)
		if err != nil {
			return nil, err
		}
		response = new(CreateMongoDBInstanceResponse)
		err = json.Unmarshal(resp.BaseResponse.GetHttpContentBytes(), &response)

		return response, err
	})
	if err != nil {
		return nil, err
	}
	return resp.(*CreateMongoDBInstanceResponse), err
}

// WaitForInstance waits for instance to given status
func (client *MongoDBService) WaitForMongoDBInstance(instanceId string, regionId string, status Status, timeout int, aliyunClient *connectivity.AliyunClient) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		instance, err := client.DescribeMongoDBInstanceById(instanceId, regionId, aliyunClient)
		if err != nil && !NotFoundError(err) && !IsExceptedError(err, InvalidDBInstanceIdNotFound) {
			return err
		}
		if instance != nil && instance.DBInstanceStatus == string(status) {
			break
		}

		if timeout <= 0 {
			return common.GetClientErrorFromString("Timeout")
		}

		timeout = timeout - DefaultIntervalMedium
		time.Sleep(DefaultIntervalMedium * time.Second)
	}
	return nil
}

func (client *MongoDBService) DescribeMongoDBInstanceById(id string, regionId string, aliyunClient *connectivity.AliyunClient) (*MongoDBInstance, error) {
	resp, err := aliyunClient.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		request := CommonRequestInit(regionId, MONGODBCode, MongoDBDomain)
		request.RegionId = regionId
		request.Version = ApiVersion20151201
		request.ApiName = "DescribeDBInstanceAttribute"
		request.QueryParams["DBInstanceId"] = id

		resp, err := ecsClient.ProcessCommonRequest(request)
		if err != nil {
			return nil, err
		}

		response := new(DescribeDBInstanceAttributeResponse)
		err = json.Unmarshal(resp.BaseResponse.GetHttpContentBytes(), &response)

		attr := response.Items.DBInstances

		if len(attr) <= 0 {
			return nil, GetNotFoundErrorFromString(fmt.Sprintf("DB instance %s is not found.", id))
		}

		return &attr[0], nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*MongoDBInstance), err
}

func (client *MongoDBService) DescribeReplicaSetRole(id string, regionId string, aliyunClient *connectivity.AliyunClient) ([]ReplicaSetRole, error) {
	resp, err := aliyunClient.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		request := CommonRequestInit(regionId, MONGODBCode, MongoDBDomain)
		request.RegionId = regionId
		request.Version = ApiVersion20151201
		request.ApiName = "DescribeReplicaSetRole"
		request.QueryParams["DBInstanceId"] = id

		resp, err := ecsClient.ProcessCommonRequest(request)
		if err != nil {
			return nil, err
		}

		response := new(ReplicaSets)
		err = json.Unmarshal(resp.BaseResponse.GetHttpContentBytes(), &response)

		attr := response.Items.ReplicaSets

		if len(attr) <= 0 {
			return nil, GetNotFoundErrorFromString(fmt.Sprintf("MongoDB instance %s is not found.", id))
		}

		return attr, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.([]ReplicaSetRole), err
}

func (client *MongoDBService) DescribeMongoDBSecurityIps(request *requests.CommonRequest, aliyunClient *connectivity.AliyunClient) (*DescribeMongoDBSecurityIpsResponse, error) {
	resp, err := aliyunClient.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		request.Version = ApiVersion20151201
		request.ApiName = "DescribeSecurityIps"
		resp, err := ecsClient.ProcessCommonRequest(request)
		if err != nil {
			return nil, err
		}
		response := new(DescribeMongoDBSecurityIpsResponse)
		err = json.Unmarshal(resp.BaseResponse.GetHttpContentBytes(), &response)
		if err != nil {
			return nil, err
		}

		return response, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*DescribeMongoDBSecurityIpsResponse), err
}

func (client *MongoDBService) DeleteMongoDBInstance(request *requests.CommonRequest, aliyunClient *connectivity.AliyunClient) error {
	_, err := aliyunClient.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		request.Version = ApiVersion20151201
		request.ApiName = "DeleteDBInstance"
		resp, err := ecsClient.ProcessCommonRequest(request)
		return resp, err
	})
	return err
}

func (client *MongoDBService) ModifyMongoDBSecurityIps(request *requests.CommonRequest, aliyunClient *connectivity.AliyunClient) error {
	_, err := aliyunClient.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		request.Version = ApiVersion20151201
		request.ApiName = "ModifySecurityIps"
		if _, ok := request.QueryParams["foo"]; ok {
			request.QueryParams["ModifyMode"] = "Cover"
		}
		resp, err := ecsClient.ProcessCommonRequest(request)
		return resp, err
	})
	return err
}

func (client *MongoDBService) ModifyMongoDBInstanceSpec(request *requests.CommonRequest, aliyunClient *connectivity.AliyunClient) error {
	_, err := aliyunClient.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		request.Version = ApiVersion20151201
		request.ApiName = "ModifyDBInstanceSpec"
		resp, err := ecsClient.ProcessCommonRequest(request)
		return resp, err
	})
	return err
}

func (client *MongoDBService) ModifyMongoDBInstanceDescription(request *requests.CommonRequest, aliyunClient *connectivity.AliyunClient) error {
	_, err := aliyunClient.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		request.ApiName = "ModifyDBInstanceDescription"
		request.Version = ApiVersion20151201
		resp, err := ecsClient.ProcessCommonRequest(request)
		return resp, err
	})
	return err
}

func (client *MongoDBService) ModifyMongoDBBackupPolicy(request *requests.CommonRequest, aliyunClient *connectivity.AliyunClient) error {
	_, err := aliyunClient.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		request.Version = ApiVersion20151201
		request.ApiName = "ModifyBackupPolicy"
		resp, err := ecsClient.ProcessCommonRequest(request)
		return resp, err
	})
	return err
}

func (client *MongoDBService) DescribeMongoDBBackupPolicy(request *requests.CommonRequest, aliyunClient *connectivity.AliyunClient) (*DescribeMongoDBBackupPolicyResponse, error) {
	resp, err := aliyunClient.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		request.Version = ApiVersion20151201
		request.ApiName = "DescribeBackupPolicy"
		resp, err := ecsClient.ProcessCommonRequest(request)
		response := new(DescribeMongoDBBackupPolicyResponse)
		err = json.Unmarshal(resp.BaseResponse.GetHttpContentBytes(), &response)
		if err != nil {
			return nil, err
		}
		return response, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*DescribeMongoDBBackupPolicyResponse), err
}

func (s *MongoDBService) NotFoundDBInstance(err error) bool {
	if NotFoundError(err) || IsExceptedErrors(err, []string{InvalidDBInstanceIdNotFound, InvalidDBInstanceNameNotFound}) {
		return true
	}

	return false
}

package alicloud

import (
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dds"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type MongoDBService struct {
	client *connectivity.AliyunClient
}

func (s *MongoDBService) NotFoundMongoDBInstance(err error) bool {
	if NotFoundError(err) || IsExceptedErrors(err, []string{InvalidMongoDBInstanceIdNotFound, InvalidMongoDBNameNotFound}) {
		return true
	}
	return false
}

func (s *MongoDBService) DescribeMongoDBInstance(id string) (instance *dds.DBInstance, err error) {
	request := dds.CreateDescribeDBInstanceAttributeRequest()
	request.DBInstanceId = id
	raw, err := s.client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
		return client.DescribeDBInstanceAttribute(request)
	})
	response, _ := raw.(*dds.DescribeDBInstanceAttributeResponse)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), response)
	if response == nil || len(response.DBInstances.DBInstance) == 0 {
		return nil, WrapErrorf(Error(GetNotFoundMessage("MongoDB Instance", id)), NotFoundMsg, AlibabaCloudSdkGoERROR)
	}
	return &response.DBInstances.DBInstance[0], nil
}

// WaitForInstance waits for instance to given statusid
func (s *MongoDBService) WaitForMongoDBInstance(instanceId string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		instance, err := s.DescribeMongoDBInstance(instanceId)
		if err != nil {
			if s.NotFoundMongoDBInstance(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if instance.DBInstanceStatus == string(status) {
			return nil
		}

		if status == Updating {
			if instance.DBInstanceStatus == "NodeCreating" ||
				instance.DBInstanceStatus == "NodeDeleting" ||
				instance.DBInstanceStatus == "DBInstanceClassChanging" {
				return nil
			}
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, instanceId, GetFunc(1), timeout, instance.DBInstanceStatus, string(status), ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *MongoDBService) GetSecurityIps(instanceId string) ([]string, error) {
	arr, err := s.DescribeMongoDBSecurityIps(instanceId)

	if err != nil {
		return nil, WrapError(err)
	}

	var ips, separator string
	ipsMap := make(map[string]string)
	for _, ip := range arr {
		ips += separator + ip.SecurityIpList
		separator = COMMA_SEPARATED
	}

	for _, ip := range strings.Split(ips, COMMA_SEPARATED) {
		ipsMap[ip] = ip
	}

	var finalIps []string
	if len(ipsMap) > 0 {
		for key := range ipsMap {
			finalIps = append(finalIps, key)
		}
	}
	return finalIps, nil
}

func (s *MongoDBService) DescribeMongoDBSecurityIps(instanceId string) (ips []dds.SecurityIpGroup, err error) {
	request := dds.CreateDescribeSecurityIpsRequest()
	request.DBInstanceId = instanceId

	raw, err := s.client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
		return client.DescribeSecurityIps(request)
	})
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, instanceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	respone, _ := raw.(*dds.DescribeSecurityIpsResponse)
	addDebug(request.GetActionName(), respone)

	return respone.SecurityIpGroups.SecurityIpGroup, nil
}

func (s *MongoDBService) ModifyMongoDBSecurityIps(instanceId, ips string) error {
	request := dds.CreateModifySecurityIpsRequest()
	request.DBInstanceId = instanceId
	request.SecurityIps = ips

	raw, err := s.client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
		return client.ModifySecurityIps(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, instanceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	respone := raw.(*dds.ModifySecurityIpsResponse)
	addDebug(request.GetActionName(), respone)

	if err := s.WaitForMongoDBInstance(instanceId, Running, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	return nil
}

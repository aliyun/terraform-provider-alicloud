package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
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

func (s *MongoDBService) DescribeMongoDBInstance(id string) (instance dds.DBInstance, err error) {
	request := dds.CreateDescribeDBInstanceAttributeRequest()
	request.DBInstanceId = id
	raw, err := s.client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
		return client.DescribeDBInstanceAttribute(request)
	})
	response, _ := raw.(*dds.DescribeDBInstanceAttributeResponse)
	if err != nil {
		return instance, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), response)
	if response == nil || len(response.DBInstances.DBInstance) == 0 {
		return instance, WrapErrorf(Error(GetNotFoundMessage("MongoDB Instance", id)), NotFoundMsg, AlibabaCloudSdkGoERROR)
	}
	return response.DBInstances.DBInstance[0], nil
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

func (server *MongoDBService) ModifyMongodbShardingInstanceNode(
	instanceID string, nodeType MongoDBShardingNodeType, stateList, diffList []interface{}) error {
	client := server.client

	//create node
	if len(stateList) < len(diffList) {
		createList := diffList[len(stateList):]
		diffList = diffList[:len(stateList)]

		for _, item := range createList {
			node := item.(map[string]interface{})

			request := dds.CreateCreateNodeRequest()
			request.DBInstanceId = instanceID
			request.NodeClass = node["node_class"].(string)
			request.NodeType = string(nodeType)

			if nodeType == MongoDBShardingNodeShard {
				request.NodeStorage = requests.NewInteger(node["node_storage"].(int))
			}

			raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
				return ddsClient.CreateNode(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, instanceID, request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw)

			err = server.WaitForMongoDBInstance(instanceID, Updating, DefaultLongTimeout)
			if err != nil {
				return WrapError(err)
			}

			err = server.WaitForMongoDBInstance(instanceID, Running, DefaultLongTimeout)
			if err != nil {
				return WrapError(err)
			}
		}
	} else if len(stateList) > len(diffList) {
		deleteList := stateList[len(diffList):]
		stateList = stateList[:len(diffList)]

		for _, item := range deleteList {
			node := item.(map[string]interface{})

			request := dds.CreateDeleteNodeRequest()
			request.DBInstanceId = instanceID
			request.NodeId = node["node_id"].(string)

			raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
				return ddsClient.DeleteNode(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, instanceID, request.GetActionName(), AlibabaCloudSdkGoERROR)
			}

			addDebug(request.GetActionName(), raw)

			err = server.WaitForMongoDBInstance(instanceID, Running, DefaultLongTimeout)
			if err != nil {
				return WrapError(err)
			}
		}
	}

	//motify node
	for key := 0; key < len(stateList); key++ {
		state := stateList[key].(map[string]interface{})
		diff := diffList[key].(map[string]interface{})

		if state["node_class"] != diff["node_class"] ||
			state["node_storage"] != diff["node_storage"] {
			request := dds.CreateModifyNodeSpecRequest()
			request.DBInstanceId = instanceID
			request.NodeClass = diff["node_class"].(string)
			if nodeType == MongoDBShardingNodeShard {
				request.NodeStorage = requests.NewInteger(diff["node_storage"].(int))
			}
			request.NodeId = state["node_id"].(string)

			raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
				return ddsClient.ModifyNodeSpec(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, instanceID, request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw)
			err = server.WaitForMongoDBInstance(instanceID, Updating, DefaultLongTimeout)
			if err != nil {
				return WrapError(err)
			}
			err = server.WaitForMongoDBInstance(instanceID, Running, DefaultLongTimeout)
			if err != nil {
				return WrapError(err)
			}
		}
	}
	return nil
}

func (s *MongoDBService) DescribeMongoDBBackupPolicy(id string) (response *dds.DescribeBackupPolicyResponse, err error) {
	request := dds.CreateDescribeBackupPolicyRequest()
	request.DBInstanceId = id
	raw, err := s.client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
		return ddsClient.DescribeBackupPolicy(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response, _ = raw.(*dds.DescribeBackupPolicyResponse)
	return response, nil
}

func (s *MongoDBService) MotifyMongoDBBackupPolicy(d *schema.ResourceData) error {
	if err := s.WaitForMongoDBInstance(d.Id(), Running, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	periodList := expandStringList(d.Get("backup_period").(*schema.Set).List())
	backupPeriod := fmt.Sprintf("%s", strings.Join(periodList[:], COMMA_SEPARATED))
	backupTime := d.Get("backup_time").(string)

	request := dds.CreateModifyBackupPolicyRequest()
	request.DBInstanceId = d.Id()
	request.PreferredBackupPeriod = backupPeriod
	request.PreferredBackupTime = backupTime
	raw, err := s.client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
		return ddsClient.ModifyBackupPolicy(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	if err := s.WaitForMongoDBInstance(d.Id(), Running, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	return nil
}

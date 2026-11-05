package mongodb

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dds"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type MongoDBService struct {
	client *connectivity.AliyunClient
}

func NewMongoDBService(client *connectivity.AliyunClient) *MongoDBService {
	return &MongoDBService{client: client}
}

func (s *MongoDBService) DescribeMongoDBInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeDBInstanceAttribute"
	client := s.client
	request := map[string]interface{}{
		"DBInstanceId": id,
	}

	idExist := false
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("MongoDB:Instance", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.DBInstances.DBInstance", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.DBInstances.DBInstance", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("MongoDB:Instance", id), errs.NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["DBInstanceId"]) == id {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("MongoDB:Instance", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *MongoDBService) RdsMongodbDBInstanceStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeMongoDBInstance(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["DBInstanceStatus"]) == failState {
				return object, fmt.Sprint(object["DBInstanceStatus"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["DBInstanceStatus"])))
			}
		}
		return object, fmt.Sprint(object["DBInstanceStatus"]), nil
	}
}

func (s *MongoDBService) RdsMongodbDBInstanceOrderStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeMongoDBInstance(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["DBInstanceOrderStatus"]) == failState {
				return object, fmt.Sprint(object["DBInstanceOrderStatus"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["DBInstanceOrderStatus"])))
			}
		}
		return object, fmt.Sprint(object["DBInstanceOrderStatus"]), nil
	}
}

func (s *MongoDBService) RdsMongoDBPublicNetworkAddressStateRefreshFunc(id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeReplicaSetRole(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		if replicaSetsMap, ok := object["ReplicaSets"].(map[string]interface{}); ok && replicaSetsMap != nil {
			if replicaSetsList, ok := replicaSetsMap["ReplicaSet"]; ok && replicaSetsList != nil {
				for _, replicaSets := range replicaSetsList.([]interface{}) {
					replicaSetsArg := replicaSets.(map[string]interface{})
					networkType, ok := replicaSetsArg["NetworkType"]
					if ok && networkType == "Public" {
						return object, "Running", nil
					}
				}
			}
		}

		return object, "NotExist", nil
	}
}

func (s *MongoDBService) DescribeMongoDBSecurityIps(instanceId string) (ips []string, err error) {
	request := dds.CreateDescribeSecurityIpsRequest()
	request.DBInstanceId = instanceId

	raw, err := s.client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
		return client.DescribeSecurityIps(request)
	})
	if err != nil {
		return ips, errs.WrapErrorf(err, errs.DefaultErrorMsg, instanceId, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	response, _ := raw.(*dds.DescribeSecurityIpsResponse)
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	var ipstr, separator string
	ipsMap := make(map[string]string)
	for _, ip := range response.SecurityIpGroups.SecurityIpGroup {
		if ip.SecurityIpGroupAttribute == "hidden" {
			continue
		}
		ipstr += separator + ip.SecurityIpList
		separator = helper.CommaSeparated
	}

	for _, ip := range strings.Split(ipstr, helper.CommaSeparated) {
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

func (s *MongoDBService) DescribeMongoDBShardingSecurityIps(instanceId string) (ips []string, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeSecurityIps"
	request := map[string]interface{}{
		"DBInstanceId": instanceId,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return ips, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return ips, errs.WrapErrorf(err, errs.DefaultErrorMsg, instanceId, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.SecurityIpGroups.SecurityIpGroup", response)
	if err != nil {
		return ips, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, instanceId, "$.SecurityIpGroups.SecurityIpGroup", response)
	}

	var ipstr, separator string
	ipsMap := make(map[string]string)
	for _, item := range v.([]interface{}) {
		ip := item.(map[string]interface{})
		if ip["SecurityIpGroupAttribute"] == "hidden" {
			continue
		}
		ipstr += separator + ip["SecurityIpList"].(string)
		separator = helper.CommaSeparated
	}

	for _, ip := range strings.Split(ipstr, helper.CommaSeparated) {
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

func (s *MongoDBService) ModifyMongoDBSecurityIps(d *schema.ResourceData, ips string) error {
	var response map[string]interface{}
	var err error
	client := s.client
	action := "ModifySecurityIps"
	request := make(map[string]interface{})
	request["RegionId"] = s.client.RegionId
	request["DBInstanceId"] = d.Id()
	request["SecurityIps"] = ips

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 0, s.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapError(err)
	}

	return nil
}

func (s *MongoDBService) ModifyMongoDBSecurityIpsGroup(d *schema.ResourceData, groupName, attribute, ips, modifyMode string) error {
	var response map[string]interface{}
	var err error
	client := s.client
	action := "ModifySecurityIps"
	request := make(map[string]interface{})
	request["RegionId"] = s.client.RegionId
	request["DBInstanceId"] = d.Id()
	request["SecurityIps"] = ips
	request["ModifyMode"] = modifyMode
	request["SecurityIpGroupName"] = groupName
	request["SecurityIpGroupAttribute"] = attribute

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 0, s.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapError(err)
	}

	return nil
}

func (s *MongoDBService) DescribeMongoDBSecurityGroupId(id string) (object []interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeSecurityGroupConfiguration"
	request := map[string]interface{}{
		"DBInstanceId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return object, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Items.RdsEcsSecurityGroupRel", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Items.RdsEcsSecurityGroupRel", response)
	}

	object = v.([]interface{})
	return object, nil
}

func (s *MongoDBService) DescribeMongoDBShardingSecurityGroupId(id string) (object []interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeSecurityGroupConfiguration"
	request := map[string]interface{}{
		"DBInstanceId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return object, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Items.RdsEcsSecurityGroupRel", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Items.RdsEcsSecurityGroupRel", response)
	}

	object = v.([]interface{})
	return object, nil
}

func (s *MongoDBService) ModifyMongodbShardingInstanceNode(d *schema.ResourceData, nodeType MongoDBShardingNodeType, stateList, diffList []interface{}) error {
	var response map[string]interface{}
	var err error
	client := s.client

	stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 0, s.RdsMongodbDBShardingInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapError(err)
	}

	//create node
	if len(stateList) < len(diffList) {
		createList := diffList[len(stateList):]
		diffList = diffList[:len(stateList)]

		for _, item := range createList {
			node := item.(map[string]interface{})

			action := "CreateNode"
			request := make(map[string]interface{})

			request["RegionId"] = s.client.RegionId
			request["DBInstanceId"] = d.Id()
			request["NodeClass"] = node["node_class"].(string)
			request["NodeType"] = string(nodeType)
			if node["readonly_replicas"] != nil {
				request["ReadonlyReplicas"] = requests.NewInteger(node["readonly_replicas"].(int))
			}
			request["ClientToken"] = helper.BuildClientToken(action)

			if nodeType == MongoDBShardingNodeShard {
				request["NodeStorage"] = requests.NewInteger(node["node_storage"].(int))
			}

			wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
			err := retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
				if err != nil {
					if errs.IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus"}) || helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			logs.AddDebug(action, response, request)

			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}

			stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, s.RdsMongodbDBShardingInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errs.WrapError(err)
			}
		}
	} else if len(stateList) > len(diffList) {
		deleteList := stateList[len(diffList):]
		stateList = stateList[:len(diffList)]

		for _, item := range deleteList {
			node := item.(map[string]interface{})
			action := "DeleteNode"
			request := make(map[string]interface{})

			request["RegionId"] = s.client.RegionId
			request["DBInstanceId"] = d.Id()
			request["NodeId"] = node["node_id"].(string)
			request["ClientToken"] = helper.BuildClientToken(action)

			wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
			err := retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
				if err != nil {
					if errs.IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus"}) || helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			logs.AddDebug(action, response, request)

			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}

			stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, s.RdsMongodbDBShardingInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errs.WrapError(err)
			}
		}
	}

	//motify node
	for key := 0; key < len(stateList); key++ {
		state := stateList[key].(map[string]interface{})
		diff := diffList[key].(map[string]interface{})

		if state["node_class"] != diff["node_class"] || state["node_storage"] != diff["node_storage"] || state["readonly_replicas"] != diff["readonly_replicas"] {
			var response map[string]interface{}
			action := "ModifyNodeSpec"
			request := make(map[string]interface{})
			if d.Get("instance_charge_type").(string) == "PrePaid" {
				if v, ok := d.GetOk("order_type"); ok {
					request["OrderType"] = v.(string)
				}
			}

			request["RegionId"] = s.client.RegionId
			request["DBInstanceId"] = d.Id()
			request["NodeClass"] = diff["node_class"].(string)
			request["ClientToken"] = helper.BuildClientToken(action)

			if nodeType == MongoDBShardingNodeShard {
				request["NodeStorage"] = diff["node_storage"].(int)
			}

			if diff["readonly_replicas"] != nil {
				request["ReadonlyReplicas"] = requests.NewInteger(diff["readonly_replicas"].(int))
			}

			request["NodeId"] = state["node_id"].(string)
			wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			logs.AddDebug(action, response, request)

			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}

			stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, s.RdsMongodbDBShardingInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errs.WrapError(err)
			}
		}
	}

	return nil
}

func (s *MongoDBService) DescribeMongoDBBackupPolicy(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeBackupPolicy"
	request := map[string]interface{}{
		"DBInstanceId": id,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return object, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})
	return object, nil
}

func (s *MongoDBService) DescribeMongoDBShardingBackupPolicy(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeBackupPolicy"
	request := map[string]interface{}{
		"DBInstanceId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return object, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})
	return object, nil
}

func (s *MongoDBService) DescribeMongoDBTDEInfo(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeDBInstanceTDEInfo"
	request := map[string]interface{}{
		"DBInstanceId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return object, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})
	return object, nil
}

func (s *MongoDBService) DescribeDBInstanceSSL(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeDBInstanceSSL"
	request := map[string]interface{}{
		"DBInstanceId": id,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return object, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})
	return object, nil
}

func (s *MongoDBService) DescribeMongoDBShardingTDEInfo(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeDBInstanceTDEInfo"
	request := map[string]interface{}{
		"DBInstanceId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return object, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})
	return object, nil
}

func (s *MongoDBService) ModifyMongoDBBackupPolicy(d *schema.ResourceData) error {
	stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 0*time.Second, s.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapError(err)
	}

	var response map[string]interface{}
	action := "ModifyBackupPolicy"

	var err error
	client := s.client

	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"DBInstanceId": d.Id(),
	}

	backupTime := d.Get("backup_time").(string)
	periodList := convert.ExpandStringList(d.Get("backup_period").(*schema.Set).List())
	backupPeriod := fmt.Sprintf("%s", strings.Join(periodList[:], helper.CommaSeparated))

	request["PreferredBackupTime"] = backupTime
	request["PreferredBackupPeriod"] = backupPeriod

	if v, ok := d.GetOkExists("backup_retention_period"); ok {
		request["BackupRetentionPeriod"] = v
	}

	if v, ok := d.GetOkExists("backup_retention_policy_on_cluster_deletion"); ok {
		request["BackupRetentionPolicyOnClusterDeletion"] = v
	}

	if v, ok := d.GetOkExists("enable_backup_log"); ok {
		request["EnableBackupLog"] = v
	}

	if v, ok := d.GetOkExists("log_backup_retention_period"); ok {
		request["LogBackupRetentionPeriod"] = v
	}

	if v, ok := d.GetOk("snapshot_backup_type"); ok {
		request["SnapshotBackupType"] = v
	}

	if v, ok := d.GetOk("backup_interval"); ok {
		request["BackupInterval"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapError(err)
	}

	return nil
}

func (s *MongoDBService) ResetAccountPassword(d *schema.ResourceData, password string, instanceType string) error {
	var response map[string]interface{}
	action := "ResetAccountPassword"

	var err error
	client := s.client

	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"DBInstanceId": d.Id(),
	}

	request["AccountName"] = "root"
	request["AccountPassword"] = password

	if instanceType == "shardingInstance" {
		request["CharacterType"] = "cs"
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	if instanceType == "instance" {
		stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, s.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapError(err)
		}
	} else {
		stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, s.RdsMongodbDBShardingInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapError(err)
		}
	}

	return err
}

func (s *MongoDBService) setInstanceTags(d *schema.ResourceData) error {
	oraw, nraw := d.GetChange("tags")
	o := oraw.(map[string]interface{})
	n := nraw.(map[string]interface{})

	create, remove := s.diffTags(s.tagsFromMap(o), s.tagsFromMap(n))

	if len(remove) > 0 {
		var tagKey []string
		for _, v := range remove {
			tagKey = append(tagKey, v.Key)
		}
		request := dds.CreateUntagResourcesRequest()
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = "INSTANCE"
		request.TagKey = &tagKey
		request.RegionId = s.client.RegionId
		raw, err := s.client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.UntagResources(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if len(create) > 0 {
		request := dds.CreateTagResourcesRequest()
		request.ResourceId = &[]string{d.Id()}
		request.Tag = &create
		request.ResourceType = "INSTANCE"
		request.RegionId = s.client.RegionId
		raw, err := s.client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.TagResources(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return nil
}

func (s *MongoDBService) tagsToMap(tags []dds.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.Key] = t.Value
		}
	}
	return result
}

func (s *MongoDBService) ignoreTag(t dds.Tag) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.Key)
		ok, _ := regexp.MatchString(v, t.Key)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific t %s (val: %s), ignoring.\n", t.Key, t.Value)
			return true
		}
	}
	return false
}

func (s *MongoDBService) tagsInAttributeToMap(tags []dds.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTagInAttribute(t) {
			result[t.Key] = t.Value
		}
	}
	return result
}

func (s *MongoDBService) ignoreTagInAttribute(t dds.Tag) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.Key)
		ok, _ := regexp.MatchString(v, t.Key)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific t %s (val: %s), ignoring.\n", t.Key, t.Value)
			return true
		}
	}
	return false
}

func (s *MongoDBService) diffTags(oldTags, newTags []dds.TagResourcesTag) ([]dds.TagResourcesTag, []dds.TagResourcesTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []dds.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return s.tagsFromMap(create), remove
}

func (s *MongoDBService) tagsFromMap(m map[string]interface{}) []dds.TagResourcesTag {
	result := make([]dds.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, dds.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func (s *MongoDBService) DescribeMongodbAuditPolicy(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeAuditPolicy"
	request := map[string]interface{}{
		"DBInstanceId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *MongoDBService) DescribeMongodbAccount(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeAccounts"
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"AccountName":  parts[1],
		"DBInstanceId": parts[0],
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Accounts.Account", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Accounts.Account", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("MongoDB", id), errs.NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["AccountName"]) != parts[1] {
			return object, errs.WrapErrorf(errs.NotFoundErr("MongoDB", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *MongoDBService) MongodbAccountStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeMongodbAccount(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["AccountStatus"]) == failState {
				return object, fmt.Sprint(object["AccountStatus"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["AccountStatus"])))
			}
		}
		return object, fmt.Sprint(object["AccountStatus"]), nil
	}
}

func (s *MongoDBService) DescribeSecurityIps(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeSecurityIps"
	request := map[string]interface{}{
		"DBInstanceId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *MongoDBService) DescribeMongodbServerlessInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeDBInstanceAttribute"
	request := map[string]interface{}{
		"DBInstanceId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("MongoDB:ServerlessInstance", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.DBInstances.DBInstance", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.DBInstances.DBInstance", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("MongoDB", id), errs.NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["DBInstanceId"]) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr("MongoDB", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *MongoDBService) MongodbServerlessInstanceStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeMongodbServerlessInstance(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["DBInstanceStatus"]) == failState {
				return object, fmt.Sprint(object["DBInstanceStatus"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["DBInstanceStatus"])))
			}
		}
		return object, fmt.Sprint(object["DBInstanceStatus"]), nil
	}
}

func (s *MongoDBService) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		added, removed := tags.ParsingTags(d)
		client := s.client
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !tags.IgnoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action := "UntagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			wait := helper.IncrementalWait(2*time.Second, 1*time.Second)
			err := retry.Retry(10*time.Minute, func() *retry.RetryError {
				response, err := client.RpcPost("Dds", "2015-12-01", action, nil, request, false)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)

					}
					return retry.NonRetryableError(err)
				}
				logs.AddDebug(action, response, request)
				return nil
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
		}
		if len(added) > 0 {
			action := "TagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			wait := helper.IncrementalWait(2*time.Second, 1*time.Second)
			err := retry.Retry(10*time.Minute, func() *retry.RetryError {
				response, err := client.RpcPost("Dds", "2015-12-01", action, nil, request, false)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)

					}
					return retry.NonRetryableError(err)
				}
				logs.AddDebug(action, response, request)
				return nil
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
		}
	}
	return nil
}

func (s *MongoDBService) ListTagResources(id string, resourceType string) (object interface{}, err error) {
	client := s.client
	action := "ListTagResources"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"ResourceType": resourceType,
		"ResourceId.1": id,
	}
	tags := make([]interface{}, 0)
	var response map[string]interface{}

	for {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err := client.RpcPost("Dds", "2015-12-01", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			v, err := jsonpath.Get("$.TagResources.TagResource", response)
			if err != nil {
				return retry.NonRetryableError(errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TagResources.TagResource", response))
			}
			if v != nil {
				tags = append(tags, v.([]interface{})...)
			}
			return nil
		})
		if err != nil {
			err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
			return
		}
		if response["NextToken"] == nil {
			break
		}
		request["NextToken"] = response["NextToken"]
	}

	return tags, nil
}

func (s *MongoDBService) DescribeMongodbShardingNetworkPublicAddress(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeShardingNetworkAddress"
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"DBInstanceId": parts[0],
		"NodeId":       parts[1],
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return object, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.NetworkAddresses.NetworkAddress", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.NetworkAddresses.NetworkAddress", response)
	}
	exist := false
	var networkAddress = make([]map[string]interface{}, 0)
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("MongoDB", id), errs.NotFoundWithResponse, response)
	} else {
		for _, item := range v.([]interface{}) {
			if item.(map[string]interface{})["NetworkType"].(string) == "Public" {
				exist = true
				networkAddress = append(networkAddress, item.(map[string]interface{}))
			}
		}
		if !exist {
			return object, errs.WrapErrorf(errs.NotFoundErr("MongoDB", id), errs.NotFoundWithResponse, response)
		}
	}
	object = make(map[string]interface{}, 0)
	object["NetworkAddress"] = networkAddress
	return object, nil
}

func (s *MongoDBService) DescribeShardingNodeType(id string) (string, error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeShardingNetworkAddress"
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return "", errs.WrapError(err)
	}
	request := map[string]interface{}{
		"DBInstanceId": parts[0],
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return "", errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return "", errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.NetworkAddresses.NetworkAddress", response)
	if err != nil {
		return "", errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.NetworkAddresses.NetworkAddress", response)
	}

	var nodeType string
	if len(v.([]interface{})) < 1 {
		return "", errs.WrapErrorf(errs.NotFoundErr("MongoDB", id), errs.NotFoundWithResponse, response)
	} else {
		for _, item := range v.([]interface{}) {
			if item.(map[string]interface{})["NodeId"].(string) == parts[1] {
				nodeType = fmt.Sprint(item.(map[string]interface{})["NodeType"])
				break
			}
		}
	}
	if nodeType == "" {
		return "", errs.WrapErrorf(errs.NotFoundErr("MongoDB", id), errs.NotFoundWithResponse, response)
	}

	return nodeType, nil
}

func (s *MongoDBService) MongodbShardingNetworkPublicAddressStateRefreshFunc(id, nodeType string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		parts, err := helper.ParseResourceId(id, 2)
		if err != nil {
			return nil, "", errs.WrapError(err)
		}
		object, err := s.DescribeMongodbServerlessInstance(parts[0])
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		var status string

		switch nodeType {
		case "mongos":
			status = fmt.Sprint(object["DBInstanceStatus"])
		case "db":
			list := object["ShardList"].(map[string]interface{})["ShardAttribute"]
			for _, item := range list.([]interface{}) {
				mongos := item.(map[string]interface{})
				if mongos["NodeId"] == parts[1] {
					status = fmt.Sprint(mongos["Status"])
					break
				}
			}
		case "cs":
			list := object["ConfigserverList"].(map[string]interface{})["ConfigserverAttribute"]
			for _, item := range list.([]interface{}) {
				mongos := item.(map[string]interface{})
				if mongos["NodeId"] == parts[1] {
					status = fmt.Sprint(mongos["Status"])
					break
				}
			}
		}

		for _, failState := range failStates {
			if status == failState {
				return object, status, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, status))
			}
		}

		return object, status, nil
	}
}

func (s *MongoDBService) DescribeMongodbShardingNetworkPrivateAddress(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeShardingNetworkAddress"
	client := s.client
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"DBInstanceId": parts[0],
		"NodeId":       parts[1],
	}

	idExist := false
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("MongoDB:ShardingNetworkPrivateAddress", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.NetworkAddresses.NetworkAddress", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.NetworkAddresses.NetworkAddress", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("MongoDB:ShardingNetworkPrivateAddress", id), errs.NotFoundWithResponse, response)
	}

	object = make(map[string]interface{}, 0)
	networkAddressMaps := make([]map[string]interface{}, 0)
	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["NodeId"]) == parts[1] && fmt.Sprint(v.(map[string]interface{})["NetworkType"]) != "Public" {
			idExist = true
			object["NodeId"] = fmt.Sprint(v.(map[string]interface{})["NodeId"])
			networkAddressMaps = append(networkAddressMaps, v.(map[string]interface{}))
		}
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("MongoDB:ShardingNetworkPrivateAddress", id), errs.NotFoundWithResponse, response)
	}

	object["NetworkAddress"] = networkAddressMaps

	return object, nil
}

func (s *MongoDBService) DescribeMongoDBShardingInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeDBInstanceAttribute"
	client := s.client
	request := map[string]interface{}{
		"DBInstanceId": id,
	}

	idExist := false
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("MongoDB:ShardingInstance", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.DBInstances.DBInstance", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.DBInstances.DBInstance", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("MongoDB:ShardingInstance", id), errs.NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["DBInstanceId"]) == id {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("MongoDB:ShardingInstance", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *MongoDBService) RdsMongodbDBShardingInstanceStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeMongoDBShardingInstance(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["DBInstanceStatus"]) == failState {
				return object, fmt.Sprint(object["DBInstanceStatus"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["DBInstanceStatus"])))
			}
		}
		return object, fmt.Sprint(object["DBInstanceStatus"]), nil
	}
}

func (s *MongoDBService) ModifyParameters(d *schema.ResourceData, attribute string) error {
	client := s.client
	action := "ModifyParameters"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"DBInstanceId": d.Id(),
	}
	config := make(map[string]string)
	o, n := d.GetChange(attribute)
	os, ns := o.(*schema.Set), n.(*schema.Set)
	add := ns.Difference(os).List()
	if len(add) > 0 {
		for _, i := range add {
			key := i.(map[string]interface{})["name"].(string)
			value := i.(map[string]interface{})["value"].(string)
			config[key] = value
		}
		cfg, _ := json.Marshal(config)
		request["Parameters"] = string(cfg)
		response, err := client.RpcPost("Dds", "2015-12-01", action, nil, request, false)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(action, response, request)

		stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, s.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapError(err)
		}
	}

	return nil
}

// lintignore: R001
func (s *MongoDBService) RefreshParameters(d *schema.ResourceData, attribute string) error {
	var param []map[string]interface{}
	documented, ok := d.GetOk(attribute)
	if !ok {
		return nil
	}

	documentedMap := make(map[string]interface{}, 0)
	for _, v := range documented.(*schema.Set).List() {
		parameter := v.(map[string]interface{})
		documentedMap[parameter["name"].(string)] = struct{}{}
	}

	object, err := s.DescribeParameters(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}
	dBInstanceParameters := object["RunningParameters"].(map[string]interface{})["Parameter"].([]interface{})
	for _, v := range dBInstanceParameters {
		item := v.(map[string]interface{})
		if item["ParameterName"] != "" {
			if _, ok := documentedMap[item["ParameterName"].(string)]; ok || len(documentedMap) == 0 {
				parameter := map[string]interface{}{
					"name":  item["ParameterName"],
					"value": item["ParameterValue"],
				}
				param = append(param, parameter)
			}
		}
	}
	if len(param) > 0 {
		if err := d.Set(attribute, param); err != nil {
			return errs.WrapError(err)
		}
	}
	return nil
}

func (s *MongoDBService) DescribeParameters(id string) (map[string]interface{}, error) {
	client := s.client
	var response map[string]interface{}
	var err error
	action := "DescribeParameters"
	request := map[string]interface{}{
		"DBInstanceId": id,
		"ExtraParam":   "terraform",
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	return response, err
}

func (s *MongoDBService) AllocatePublicNetworkAddress(id string) error {
	client := s.client
	var response map[string]interface{}
	var err error
	action := "AllocatePublicNetworkAddress"
	request := map[string]interface{}{
		"DBInstanceId": id,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	return err
}

func (s *MongoDBService) ReleasePublicNetworkAddress(id string) error {
	client := s.client
	var response map[string]interface{}
	var err error
	action := "ReleasePublicNetworkAddress"
	request := map[string]interface{}{
		"DBInstanceId": id,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	return err
}

func getPrefixOfConnectionDomain(domain string) string {
	parts := strings.Split(domain, ".")
	// in fact, len(parts) should be equal to 5
	return parts[0]
}

func (s *MongoDBService) ModifyDBInstanceConnectionString(d *schema.ResourceData, instanceId string, currentConnectionString, newConnectionStringPrefix string, port string) error {
	client := s.client
	var response map[string]interface{}
	var err error

	action := "ModifyDBInstanceConnectionString"
	request := map[string]interface{}{
		"DBInstanceId":            instanceId,
		"CurrentConnectionString": currentConnectionString,
		"NewConnectionString":     newConnectionStringPrefix,
		"NewPort":                 port,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return err
	}

	stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, s.RdsMongodbDBInstanceStateRefreshFunc(instanceId, []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapError(err)
	}
	return nil
}

func (s *MongoDBService) DescribeReplicaSetRole(id string) (map[string]interface{}, error) {
	client := s.client
	var response map[string]interface{}
	var err error
	action := "DescribeReplicaSetRole"
	request := map[string]interface{}{
		"DBInstanceId": id,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	return response, err
}

func transferToMongoReplicaSets(objects map[string]interface{}, sortByRoleId bool) []map[string]interface{} {
	if replicaSetsMap, ok := objects["ReplicaSets"].(map[string]interface{}); ok && replicaSetsMap != nil {
		if replicaSetsListInterface, ok := replicaSetsMap["ReplicaSet"]; ok && replicaSetsListInterface != nil {
			replicaSetsList := replicaSetsListInterface.([]interface{})
			ret := make([]map[string]interface{}, 0, len(replicaSetsList))
			for _, replicaSets := range replicaSetsList {
				replicaSetsArg := replicaSets.(map[string]interface{})
				replicaSetsItemMap := make(map[string]interface{})

				if connectionPort, ok := replicaSetsArg["ConnectionPort"]; ok {
					replicaSetsItemMap["connection_port"] = connectionPort
				}

				if replicaSetRole, ok := replicaSetsArg["ReplicaSetRole"]; ok {
					replicaSetsItemMap["replica_set_role"] = replicaSetRole
				}

				if connectionDomain, ok := replicaSetsArg["ConnectionDomain"]; ok {
					replicaSetsItemMap["connection_domain"] = connectionDomain
				}

				if networkType, ok := replicaSetsArg["NetworkType"]; ok {
					replicaSetsItemMap["network_type"] = networkType
				}

				if roleID, ok := replicaSetsArg["RoleId"]; ok {
					replicaSetsItemMap["role_id"] = roleID
				}

				ret = append(ret, replicaSetsItemMap)
			}

			if sortByRoleId {
				sort.Slice(ret, func(i, j int) bool {
					r1, ok := ret[i]["role_id"]
					if !ok {
						return false
					}
					r2, ok := ret[j]["role_id"]
					if !ok {
						return true
					}
					return strings.Compare(r1.(string), r2.(string)) >= 0
				})
			}

			return ret
		}
	}
	return nil
}

func (s *MongoDBService) DescribeMongoDBGlobalSecurityGroupIds(id string) (globalSecurityGroupIds []string, err error) {
	var response map[string]interface{}
	action := "DescribeGlobalSecurityIPGroupRelation"

	client := s.client

	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"DBClusterId": id,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("Dds", "2015-12-01", action, request, nil)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return globalSecurityGroupIds, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.GlobalSecurityIPGroupRel", response)
	if err != nil {
		return globalSecurityGroupIds, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.GlobalSecurityIPGroupRel", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return globalSecurityGroupIds, nil
	}

	for _, v := range resp.([]interface{}) {
		globalSecurityGroupIds = append(globalSecurityGroupIds, fmt.Sprint(v.(map[string]interface{})["GlobalSecurityGroupId"]))
	}

	return globalSecurityGroupIds, nil
}

func (s *MongoDBService) DescribeMongoDBRecoverTime(d *schema.ResourceData, id string) (objects []interface{}, err error) {
	var response map[string]interface{}

	action := "DescribeInstanceRecoverTime"

	if fmt.Sprint(d.Get("instance_type")) == "sharding" {
		action = "DescribeClusterRecoverTime"
	}

	client := s.client

	request := map[string]interface{}{
		"DBInstanceId": id,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return objects, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.RestoreRanges", response)
	if err != nil {
		return objects, nil
	}

	objects = resp.([]interface{})

	return objects, nil
}

func (s *MongoDBService) DescribeMongoDBAutoRenewalAttribute(id, instanceType string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeInstanceAutoRenewalAttribute"
	client := s.client
	request := map[string]interface{}{
		"RegionId":       s.client.RegionId,
		"DBInstanceId":   id,
		"DBInstanceType": instanceType,
	}

	idExist := false
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.Items.Item", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Items.Item", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, nil
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["DbInstanceId"]) == id {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, nil
	}

	return object, nil
}

func (s *MongoDBService) DescribeMongoDBRoleZoneInfos(id string) (objects []interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeRoleZoneInfo"

	client := s.client

	request := map[string]interface{}{
		"DBInstanceId": id,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return objects, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.ZoneInfos.ZoneInfo", response)
	if err != nil {
		return objects, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ZoneInfos.ZoneInfo", response)
	}

	objects = resp.([]interface{})

	return objects, nil
}

func (s *MongoDBService) DescribeMongoDBUserEncryptionKeyList(id string) (keyIds []interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeUserEncryptionKeyList"

	client := s.client

	request := map[string]interface{}{
		"DBInstanceId": id,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return keyIds, nil
	}

	resp, err := jsonpath.Get("$.KeyIds.KeyId", response)
	if err != nil {
		return keyIds, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.KeyIds.KeyId", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return keyIds, nil
	}

	keyIds = resp.([]interface{})

	return keyIds, nil
}

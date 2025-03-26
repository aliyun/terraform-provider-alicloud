package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dds"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

type MongoDBService struct {
	client *connectivity.AliyunClient
}

func (s *MongoDBService) DescribeMongoDBInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeDBInstanceAttribute"
	client := s.client
	request := map[string]interface{}{
		"DBInstanceId": id,
	}

	idExist := false
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("MongoDB:Instance", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.DBInstances.DBInstance", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.DBInstances.DBInstance", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, WrapErrorf(NotFoundErr("MongoDB:Instance", id), NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["DBInstanceId"]) == id {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, WrapErrorf(NotFoundErr("MongoDB:Instance", id), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *MongoDBService) RdsMongodbDBInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeMongoDBInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["DBInstanceStatus"]) == failState {
				return object, fmt.Sprint(object["DBInstanceStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["DBInstanceStatus"])))
			}
		}
		return object, fmt.Sprint(object["DBInstanceStatus"]), nil
	}
}

func (s *MongoDBService) RdsMongodbDBInstanceOrderStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeMongoDBInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["DBInstanceOrderStatus"]) == failState {
				return object, fmt.Sprint(object["DBInstanceOrderStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["DBInstanceOrderStatus"])))
			}
		}
		return object, fmt.Sprint(object["DBInstanceOrderStatus"]), nil
	}
}

func (s *MongoDBService) DescribeMongoDBSecurityIps(instanceId string) (ips []string, err error) {
	request := dds.CreateDescribeSecurityIpsRequest()
	request.DBInstanceId = instanceId

	raw, err := s.client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
		return client.DescribeSecurityIps(request)
	})
	if err != nil {
		return ips, WrapErrorf(err, DefaultErrorMsg, instanceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response, _ := raw.(*dds.DescribeSecurityIpsResponse)
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	var ipstr, separator string
	ipsMap := make(map[string]string)
	for _, ip := range response.SecurityIpGroups.SecurityIpGroup {
		if ip.SecurityIpGroupAttribute == "hidden" {
			continue
		}
		ipstr += separator + ip.SecurityIpList
		separator = COMMA_SEPARATED
	}

	for _, ip := range strings.Split(ipstr, COMMA_SEPARATED) {
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return ips, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return ips, WrapErrorf(err, DefaultErrorMsg, instanceId, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.SecurityIpGroups.SecurityIpGroup", response)
	if err != nil {
		return ips, WrapErrorf(err, FailedGetAttributeMsg, instanceId, "$.SecurityIpGroups.SecurityIpGroup", response)
	}

	var ipstr, separator string
	ipsMap := make(map[string]string)
	for _, item := range v.([]interface{}) {
		ip := item.(map[string]interface{})
		if ip["SecurityIpGroupAttribute"] == "hidden" {
			continue
		}
		ipstr += separator + ip["SecurityIpList"].(string)
		separator = COMMA_SEPARATED
	}

	for _, ip := range strings.Split(ipstr, COMMA_SEPARATED) {
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

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 0, s.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return object, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Items.RdsEcsSecurityGroupRel", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Items.RdsEcsSecurityGroupRel", response)
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return object, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Items.RdsEcsSecurityGroupRel", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Items.RdsEcsSecurityGroupRel", response)
	}

	object = v.([]interface{})
	return object, nil
}

func (s *MongoDBService) ModifyMongodbShardingInstanceNode(d *schema.ResourceData, nodeType MongoDBShardingNodeType, stateList, diffList []interface{}) error {
	var response map[string]interface{}
	var err error
	client := s.client

	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 0, s.RdsMongodbDBShardingInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
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
			request["ClientToken"] = buildClientToken(action)

			if nodeType == MongoDBShardingNodeShard {
				request["NodeStorage"] = requests.NewInteger(node["node_storage"].(int))
			}

			wait := incrementalWait(3*time.Second, 3*time.Second)
			err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
				if err != nil {
					if IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus"}) || NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, request)

			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

			stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, s.RdsMongodbDBShardingInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapError(err)
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
			request["ClientToken"] = buildClientToken(action)

			wait := incrementalWait(3*time.Second, 3*time.Second)
			err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
				if err != nil {
					if IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus"}) || NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, request)

			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

			stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, s.RdsMongodbDBShardingInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
			if _, err := stateConf.WaitForState(); err != nil {
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
			request["ClientToken"] = buildClientToken(action)

			if nodeType == MongoDBShardingNodeShard {
				request["NodeStorage"] = diff["node_storage"].(int)
			}
			request["NodeId"] = state["node_id"].(string)
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
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
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

			stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, s.RdsMongodbDBShardingInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapError(err)
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

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return object, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return object, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return object, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
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

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return object, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return object, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})
	return object, nil
}

func (s *MongoDBService) ModifyMongoDBBackupPolicy(d *schema.ResourceData) error {
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 0*time.Second, s.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
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
	periodList := expandStringList(d.Get("backup_period").(*schema.Set).List())
	backupPeriod := fmt.Sprintf("%s", strings.Join(periodList[:], COMMA_SEPARATED))

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

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
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

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	d.SetPartial("tags")
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
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
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *MongoDBService) DescribeMongodbAccount(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeAccounts"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"AccountName":  parts[1],
		"DBInstanceId": parts[0],
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
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
	v, err := jsonpath.Get("$.Accounts.Account", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Accounts.Account", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("MongoDB", id), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["AccountName"]) != parts[1] {
			return object, WrapErrorf(NotFoundErr("MongoDB", id), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *MongoDBService) MongodbAccountStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeMongodbAccount(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["AccountStatus"]) == failState {
				return object, fmt.Sprint(object["AccountStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["AccountStatus"])))
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
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
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("MongoDB:ServerlessInstance", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.DBInstances.DBInstance", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.DBInstances.DBInstance", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("MongoDB", id), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["DBInstanceId"]) != id {
			return object, WrapErrorf(NotFoundErr("MongoDB", id), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *MongoDBService) MongodbServerlessInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeMongodbServerlessInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["DBInstanceStatus"]) == failState {
				return object, fmt.Sprint(object["DBInstanceStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["DBInstanceStatus"])))
			}
		}
		return object, fmt.Sprint(object["DBInstanceStatus"]), nil
	}
}

func (s *MongoDBService) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		added, removed := parsingTags(d)
		client := s.client
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
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
			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := client.RpcPost("Dds", "2015-12-01", action, nil, request, false)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)

					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
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

			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := client.RpcPost("Dds", "2015-12-01", action, nil, request, false)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)

					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}
		d.SetPartial("tags")
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
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err := client.RpcPost("Dds", "2015-12-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			v, err := jsonpath.Get("$.TagResources.TagResource", response)
			if err != nil {
				return resource.NonRetryableError(WrapErrorf(err, FailedGetAttributeMsg, id, "$.TagResources.TagResource", response))
			}
			if v != nil {
				tags = append(tags, v.([]interface{})...)
			}
			return nil
		})
		if err != nil {
			err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
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
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"DBInstanceId": parts[0],
		"NodeId":       parts[1],
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return object, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.NetworkAddresses.NetworkAddress", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.NetworkAddresses.NetworkAddress", response)
	}
	exist := false
	var networkAddress = make([]map[string]interface{}, 0)
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("MongoDB", id), NotFoundWithResponse, response)
	} else {
		for _, item := range v.([]interface{}) {
			if item.(map[string]interface{})["NetworkType"].(string) == "Public" {
				exist = true
				networkAddress = append(networkAddress, item.(map[string]interface{}))
			}
		}
		if !exist {
			return object, WrapErrorf(NotFoundErr("MongoDB", id), NotFoundWithResponse, response)
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
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return "", WrapError(err)
	}
	request := map[string]interface{}{
		"DBInstanceId": parts[0],
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return "", WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return "", WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.NetworkAddresses.NetworkAddress", response)
	if err != nil {
		return "", WrapErrorf(err, FailedGetAttributeMsg, id, "$.NetworkAddresses.NetworkAddress", response)
	}

	var nodeType string
	if len(v.([]interface{})) < 1 {
		return "", WrapErrorf(NotFoundErr("MongoDB", id), NotFoundWithResponse, response)
	} else {
		for _, item := range v.([]interface{}) {
			if item.(map[string]interface{})["NodeId"].(string) == parts[1] {
				nodeType = fmt.Sprint(item.(map[string]interface{})["NodeType"])
				break
			}
		}
	}
	if nodeType == "" {
		return "", WrapErrorf(NotFoundErr("MongoDB", id), NotFoundWithResponse, response)
	}

	return nodeType, nil
}

func (s *MongoDBService) MongodbShardingNetworkPublicAddressStateRefreshFunc(id, nodeType string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		parts, err := ParseResourceId(id, 2)
		if err != nil {
			return nil, "", WrapError(err)
		}
		object, err := s.DescribeMongodbServerlessInstance(parts[0])
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
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
				return object, status, WrapError(Error(FailedToReachTargetStatus, status))
			}
		}

		return object, status, nil
	}
}

func (s *MongoDBService) DescribeMongodbShardingNetworkPrivateAddress(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeShardingNetworkAddress"
	client := s.client
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"DBInstanceId": parts[0],
		"NodeId":       parts[1],
	}

	idExist := false
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("MongoDB:ShardingNetworkPrivateAddress", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.NetworkAddresses.NetworkAddress", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.NetworkAddresses.NetworkAddress", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, WrapErrorf(NotFoundErr("MongoDB:ShardingNetworkPrivateAddress", id), NotFoundWithResponse, response)
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
		return object, WrapErrorf(NotFoundErr("MongoDB:ShardingNetworkPrivateAddress", id), NotFoundWithResponse, response)
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("MongoDB:ShardingInstance", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.DBInstances.DBInstance", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.DBInstances.DBInstance", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, WrapErrorf(NotFoundErr("MongoDB:ShardingInstance", id), NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["DBInstanceId"]) == id {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, WrapErrorf(NotFoundErr("MongoDB:ShardingInstance", id), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *MongoDBService) RdsMongodbDBShardingInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeMongoDBShardingInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["DBInstanceStatus"]) == failState {
				return object, fmt.Sprint(object["DBInstanceStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["DBInstanceStatus"])))
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, s.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}
	}

	d.SetPartial(attribute)

	return nil
}

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
		return WrapError(err)
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
			return WrapError(err)
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

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
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
	return response, err
}

package alicloud

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type R_kvstoreService struct {
	client *connectivity.AliyunClient
}

type RKvstoreService struct {
	client *connectivity.AliyunClient
}

func (s *R_kvstoreService) DescribeInstanceSSL(id string) (object r_kvstore.DescribeInstanceSSLResponse, err error) {
	request := r_kvstore.CreateDescribeInstanceSSLRequest()
	request.RegionId = s.client.RegionId
	request.InstanceId = id

	var raw interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.DescribeInstanceSSL(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"LockTimeout"}) || NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstance.NotFound"}) {
			err = WrapErrorf(NotFoundErr("KvstoreInstance", id), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return object, err
	}

	response, _ := raw.(*r_kvstore.DescribeInstanceSSLResponse)
	return *response, nil
}

func (s *R_kvstoreService) DescribeSecurityIps(id, securityIpGroupName string) (object r_kvstore.SecurityIpGroup, err error) {
	request := r_kvstore.CreateDescribeSecurityIpsRequest()
	request.RegionId = s.client.RegionId
	request.InstanceId = id

	var raw interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.DescribeSecurityIps(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"LockTimeout"}) || NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstance.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("Redis", id), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*r_kvstore.DescribeSecurityIpsResponse)

	if len(response.SecurityIpGroups.SecurityIpGroup) < 1 {
		return object, nil
	}

	if securityIpGroupName == "" {
		securityIpGroupName = "default"
	}

	for _, v := range response.SecurityIpGroups.SecurityIpGroup {
		if v.SecurityIpGroupName == "ali_dms_group" || v.SecurityIpGroupName == "hdm_security_ips" {
			continue
		}

		if v.SecurityIpGroupName == securityIpGroupName {
			return v, nil
		}
	}

	return response.SecurityIpGroups.SecurityIpGroup[0], nil
}

func (s *R_kvstoreService) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	oldItems, newItems := d.GetChange("tags")
	added := make([]r_kvstore.TagResourcesTag, 0)
	for key, value := range newItems.(map[string]interface{}) {
		added = append(added, r_kvstore.TagResourcesTag{
			Key:   key,
			Value: value.(string),
		})
	}
	removed := make([]string, 0)
	for key := range oldItems.(map[string]interface{}) {
		removed = append(removed, key)
	}
	// 对系统 Tag 进行过滤
	removedTagKeys := make([]string, 0)
	for _, v := range removed {
		if !ignoredTags(v, "") {
			removedTagKeys = append(removedTagKeys, v)
		}
	}
	if len(removedTagKeys) > 0 {
		request := r_kvstore.CreateUntagResourcesRequest()
		request.RegionId = s.client.RegionId
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = resourceType
		request.TagKey = &removedTagKeys
		raw, err := s.client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.UntagResources(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidDbInstanceId.NotFound"}) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	if len(added) > 0 {
		request := r_kvstore.CreateTagResourcesRequest()
		request.RegionId = s.client.RegionId
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = resourceType
		request.Tag = &added
		raw, err := s.client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.TagResources(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidDbInstanceId.NotFound"}) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return nil
}

func (s *R_kvstoreService) DescribeInstanceAutoRenewalAttribute(id string) (object r_kvstore.Item, err error) {
	request := r_kvstore.CreateDescribeInstanceAutoRenewalAttributeRequest()
	request.RegionId = s.client.RegionId
	request.DBInstanceId = id

	var raw interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.DescribeInstanceAutoRenewalAttribute(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"LockTimeout"}) || NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstance.NotFound"}) {
			err = WrapErrorf(NotFoundErr("KvstoreInstance", id), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return object, err
	}

	response, _ := raw.(*r_kvstore.DescribeInstanceAutoRenewalAttributeResponse)

	if len(response.Items.Item) < 1 {
		err = WrapErrorf(NotFoundErr("KvstoreInstance", id), NotFoundMsg, ProviderERROR, response.RequestId)
		return object, err
	}
	return response.Items.Item[0], nil
}

func (s *R_kvstoreService) DescribeSecurityGroupConfiguration(id string) (object r_kvstore.ItemsInDescribeSecurityGroupConfiguration, err error) {
	request := r_kvstore.CreateDescribeSecurityGroupConfigurationRequest()
	request.RegionId = s.client.RegionId
	request.InstanceId = id

	var raw interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.DescribeSecurityGroupConfiguration(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"LockTimeout"}) || NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstance.NotFound"}) {
			err = WrapErrorf(NotFoundErr("KvstoreInstance", id), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return object, err
	}

	response, _ := raw.(*r_kvstore.DescribeSecurityGroupConfigurationResponse)

	if len(response.Items.EcsSecurityGroupRelation) < 1 {
		return object, nil
	}
	return response.Items, nil
}

func (s *R_kvstoreService) DescribeKvstoreInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeInstanceAttribute"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = s.client.RpcPost("R-kvstore", "2015-01-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			err = WrapErrorf(NotFoundErr("KvstoreInstance", id), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	v, err := jsonpath.Get("$.Instances.DBInstanceAttribute", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Instances", response)
	}
	if v == nil || len(v.([]interface{})) < 1 || fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["InstanceId"]) != id {
		return object, WrapErrorf(NotFoundErr("Redis", id), NotFoundWithResponse, response)
	}
	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *R_kvstoreService) KvstoreInstanceStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeKvstoreInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["InstanceStatus"] == failState {
				return object, fmt.Sprint(object["InstanceStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["InstanceStatus"])))
			}
		}
		return object, fmt.Sprint(object["InstanceStatus"]), nil
	}
}

func (s *R_kvstoreService) KvstoreInstanceAttributeRefreshFunc(id, attribute string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeKvstoreInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		return object, fmt.Sprint(object[attribute]), nil
	}
}

func (s *R_kvstoreService) DescribeKvstoreInstances(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeInstances"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"InstanceIds": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = s.client.RpcPost("R-kvstore", "2015-01-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			err = WrapErrorf(NotFoundErr("KvstoreInstance", id), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	v, err := jsonpath.Get("$.Instances.KVStoreInstance", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Instances", response)
	}
	if v == nil || len(v.([]interface{})) < 1 || fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["InstanceId"]) != id {
		return object, WrapErrorf(NotFoundErr("Redis", id), NotFoundWithResponse, response)
	}
	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *R_kvstoreService) KvstoreInstancesStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeKvstoreInstances(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["InstanceStatus"] == failState {
				return object, fmt.Sprint(object["InstanceStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["InstanceStatus"])))
			}
		}
		return object, fmt.Sprint(object["InstanceStatus"]), nil
	}
}

func (s *R_kvstoreService) DescribeKvstoreInstanceDeleted(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeInstances"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"InstanceIds": id,
		"Expired":     true,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = s.client.RpcPost("R-kvstore", "2015-01-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			err = WrapErrorf(NotFoundErr("KvstoreInstance", id), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	v, err := jsonpath.Get("$.Instances.KVStoreInstance", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Instances", response)
	}
	if v == nil || len(v.([]interface{})) < 1 || fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["InstanceId"]) != id {
		return object, WrapErrorf(NotFoundErr("Redis", id), NotFoundWithResponse, response)
	}
	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *R_kvstoreService) KvstoreInstanceDeletedStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeKvstoreInstanceDeleted(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["InstanceStatus"] == failState {
				return object, fmt.Sprint(object["InstanceStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["InstanceStatus"])))
			}
		}
		return object, fmt.Sprint(object["InstanceStatus"]), nil
	}
}

func (s *R_kvstoreService) DescribeKvstoreConnection(id string) (object []r_kvstore.InstanceNetInfo, err error) {
	request := r_kvstore.CreateDescribeDBInstanceNetInfoRequest()
	request.RegionId = s.client.RegionId
	request.InstanceId = id

	var raw interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.DescribeDBInstanceNetInfo(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"LockTimeout"}) || NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			err = WrapErrorf(NotFoundErr("KvstoreConnection", id), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}

	response, _ := raw.(*r_kvstore.DescribeDBInstanceNetInfoResponse)

	if len(response.NetInfoItems.InstanceNetInfo) < 1 {
		err = WrapErrorf(NotFoundErr("KvstoreConnection", id), NotFoundMsg, ProviderERROR, response.RequestId)
		return
	}
	return response.NetInfoItems.InstanceNetInfo, nil
}

func (s *R_kvstoreService) DescribeKvstoreAccount(id string) (object r_kvstore.Account, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := r_kvstore.CreateDescribeAccountsRequest()
	request.RegionId = s.client.RegionId
	request.AccountName = parts[1]
	request.InstanceId = parts[0]

	var raw interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.DescribeAccounts(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"LockTimeout"}) || NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			err = WrapErrorf(NotFoundErr("KvstoreAccount", id), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}

	response, _ := raw.(*r_kvstore.DescribeAccountsResponse)

	if len(response.Accounts.Account) < 1 {
		err = WrapErrorf(NotFoundErr("KvstoreAccount", id), NotFoundMsg, ProviderERROR, response.RequestId)
		return
	}
	return response.Accounts.Account[0], nil
}

func (s *R_kvstoreService) KvstoreAccountStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeKvstoreAccount(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.AccountStatus == failState {
				return object, object.AccountStatus, WrapError(Error(FailedToReachTargetStatus, object.AccountStatus))
			}
		}
		return object, object.AccountStatus, nil
	}
}

func (s *R_kvstoreService) DescribeBackupPolicy(id string) (object r_kvstore.DescribeBackupPolicyResponse, err error) {
	request := r_kvstore.CreateDescribeBackupPolicyRequest()
	request.RegionId = s.client.RegionId
	request.InstanceId = id

	var raw interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.DescribeBackupPolicy(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"LockTimeout"}) || NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstance.NotFound"}) {
			err = WrapErrorf(NotFoundErr("KvstoreInstance", id), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return object, err
	}

	response, _ := raw.(*r_kvstore.DescribeBackupPolicyResponse)
	return *response, nil
}

func (s *RKvstoreService) DescribeKvstoreAuditLogConfig(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeAuditLogConfig"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = s.client.RpcPost("R-kvstore", "2015-01-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("Redis:AuditLogConfig", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
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

func (s *RKvstoreService) DescribeInstanceAttribute(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeInstanceAttribute"
	request := map[string]interface{}{
		"InstanceId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = s.client.RpcPost("R-kvstore", "2015-01-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("Redis:AuditLogConfig", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Instances.DBInstanceAttribute", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Instances.DBInstanceAttribute", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("Redis", id), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["InstanceId"]) != id {
			return object, WrapErrorf(NotFoundErr("Redis", id), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *RKvstoreService) KvstoreAuditLogConfigStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeInstanceAttribute(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["InstanceStatus"]) == failState {
				return object, fmt.Sprint(object["InstanceStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["InstanceStatus"])))
			}
		}
		return object, fmt.Sprint(object["InstanceStatus"]), nil
	}
}

func (s *RKvstoreService) DescribeInstanceAutoRenewalAttribute(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeInstanceAutoRenewalAttribute"
	request := map[string]interface{}{
		"DBInstanceId": id,
		"RegionId":     s.client.RegionId,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = s.client.RpcPost("R-kvstore", "2015-01-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidOrderCharge.NotSupport"}) {
			return object, WrapErrorf(NotFoundErr("Redis", id), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Items.Item", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Items.Item", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("Redis", id), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["DBInstanceId"]) != id {
			return object, WrapErrorf(NotFoundErr("Redis", id), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *RKvstoreService) DescribeInstanceSSL(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeInstanceSSL"
	request := map[string]interface{}{
		"InstanceId": id,
		"RegionId":   s.client.RegionId,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = s.client.RpcPost("R-kvstore", "2015-01-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"IncorrectEngineVersion"}) {
			return object, WrapErrorf(NotFoundErr("Redis", id), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (s *RKvstoreService) DescribeSecurityGroupConfiguration(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeSecurityGroupConfiguration"
	request := map[string]interface{}{
		"InstanceId": id,
		"RegionId":   s.client.RegionId,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = s.client.RpcPost("R-kvstore", "2015-01-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Items.EcsSecurityGroupRelation", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Items.EcsSecurityGroupRelation", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("Redis", id), NotFoundWithResponse, response)
	}

	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *RKvstoreService) DescribeSecurityIps(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeSecurityIps"
	request := map[string]interface{}{
		"InstanceId": id,
		"RegionId":   s.client.RegionId,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = s.client.RpcPost("R-kvstore", "2015-01-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.SecurityIpGroups.SecurityIpGroup", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.SecurityIpGroups.SecurityIpGroup", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("Redis", id), NotFoundWithResponse, response)
	} else {
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["SecurityIpGroupName"]) != "default" {
				return v.(map[string]interface{}), nil
			}
		}
	}

	return object, nil
}

func (s *RKvstoreService) DescribeInstanceTDEStatus(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeInstanceTDEStatus"
	request := map[string]interface{}{
		"InstanceId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = s.client.RpcPost("R-kvstore", "2015-01-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		// Instances that do not support TDE reject DescribeInstanceTDEStatus with a 400
		// InstanceType.NotSupport. Treat it as "TDE not supported" by returning an empty
		// object without error, so the read is not aborted.
		if IsExpectedErrors(err, []string{"InstanceType.NotSupport"}) {
			return object, nil
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

func (s *RKvstoreService) DescribeEncryptionKey(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeEncryptionKey"
	request := map[string]interface{}{
		"InstanceId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = s.client.RpcPost("R-kvstore", "2015-01-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
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

func (s *R_kvstoreService) DescribeKvStoreInstanceNetInfo(id string) (objects []interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeDBInstanceNetInfo"
	request := map[string]interface{}{
		"InstanceId": id,
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = s.client.RpcPost("R-kvstore", "2015-01-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return objects, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.NetInfoItems.InstanceNetInfo", response)
	if err != nil {
		return objects, WrapErrorf(err, FailedGetAttributeMsg, id, "$.NetInfoItems.InstanceNetInfo", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return objects, WrapErrorf(NotFoundErr("Redis", id), NotFoundWithResponse, response)
	}

	objects = resp.([]interface{})

	return objects, nil
}

func (s *R_kvstoreService) DescribeKvStoreEngineVersion(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeEngineVersion"
	request := map[string]interface{}{
		"InstanceId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = s.client.RpcPost("R-kvstore", "2015-01-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
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

func (s *R_kvstoreService) DescribeKvStoreGlobalSecurityGroupIds(id string) (globalSecurityGroupIds []string, err error) {
	var response map[string]interface{}
	action := "DescribeGlobalSecurityIPGroupRelation"

	client := s.client

	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"DBClusterId": id,
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("R-kvstore", "2015-01-01", action, request, nil)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return globalSecurityGroupIds, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.GlobalSecurityIPGroupRel", response)
	if err != nil {
		return globalSecurityGroupIds, WrapErrorf(err, FailedGetAttributeMsg, id, "$.GlobalSecurityIPGroupRel", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return globalSecurityGroupIds, nil
	}

	for _, v := range resp.([]interface{}) {
		globalSecurityGroupIds = append(globalSecurityGroupIds, fmt.Sprint(v.(map[string]interface{})["GlobalSecurityGroupId"]))
	}

	return globalSecurityGroupIds, nil
}

// DescribeKvstoreDataNodes lists every data node of an ApsaraDB for Redis/Tair instance with
// its default and current bandwidth. DescribeRoleZoneInfo reports CurrentBandWidth as
// "default bandwidth plus purchased additional bandwidth", which is the only read-only way to
// observe the additional bandwidth bought through EnableAdditionalBandwidth.
func (s *R_kvstoreService) DescribeKvstoreDataNodes(id string) (nodes []map[string]interface{}, err error) {
	action := "DescribeRoleZoneInfo"
	return collectAdditionalBandwidthNodes(func(pageNumber, pageSize int) ([]map[string]interface{}, int, error) {
		request := map[string]interface{}{
			"InstanceId": id,
			"QueryType":  1,
			"PageNumber": pageNumber,
			"PageSize":   pageSize,
		}

		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = s.client.RpcPost("R-kvstore", "2015-01-01", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return nil, 0, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}

		pageNodeList, _ := jsonpath.Get("$.Node.NodeInfo", response)
		pageNodes, _ := pageNodeList.([]interface{})
		nodes := make([]map[string]interface{}, 0, len(pageNodes))
		for _, v := range pageNodes {
			if node, ok := v.(map[string]interface{}); ok {
				nodes = append(nodes, node)
			}
		}

		totalCount := 0
		if v, ok := response["TotalCount"]; ok {
			totalCount, _ = strconv.Atoi(fmt.Sprint(v))
		}
		return nodes, totalCount, nil
	})
}

// collectAdditionalBandwidthNodes walks the DescribeRoleZoneInfo pages until every reported node
// is collected. It is separate from the request so the paging boundary can be unit tested.
// errIncompleteAdditionalBandwidthNodes reports a paginated node listing that ended before
// TotalCount was reached; callers keep the previous state instead of writing a partial vector.
var errIncompleteAdditionalBandwidthNodes = fmt.Errorf("incomplete AdditionalBandwidth node listing")

func collectAdditionalBandwidthNodes(fetch func(pageNumber, pageSize int) ([]map[string]interface{}, int, error)) ([]map[string]interface{}, error) {
	const pageSize = 50
	nodes := make([]map[string]interface{}, 0)
	for pageNumber := 1; ; pageNumber++ {
		pageNodes, totalCount, err := fetch(pageNumber, pageSize)
		if err != nil {
			return nodes, err
		}
		nodes = append(nodes, pageNodes...)
		if totalCount > 0 {
			if len(nodes) >= totalCount {
				return nodes, nil
			}
			// A page shorter than requested is not proof of the end: the server can cap the
			// page size. Keep walking until TotalCount is reached or a page comes back empty;
			// an empty page before TotalCount means the listing is incomplete and returning the
			// partial set would silently shrink the ordered bandwidth vector.
			if len(pageNodes) == 0 || pageNumber >= 1000 {
				return nodes, fmt.Errorf("%w: collected %d of %d", errIncompleteAdditionalBandwidthNodes, len(nodes), totalCount)
			}
			continue
		}
		// Without a total count the end of the listing is signalled by a short page.
		if len(pageNodes) < pageSize || pageNumber >= 1000 {
			return nodes, nil
		}
	}
}

// kvstoreAdditionalBandwidth is one configured additional_bandwidth block: the data shard it
// applies to and the bandwidth to add on that shard.
type kvstoreAdditionalBandwidth struct {
	NodeId    string
	Bandwidth int
}

// additionalBandwidthEntriesFromRawConfig returns the additional_bandwidth blocks the user
// actually configured. The attribute is Optional+Computed, so d.Get() would also return values
// that were only read back from the API; only the raw configuration says what to request.
func additionalBandwidthEntriesFromRawConfig(d *schema.ResourceData) ([]kvstoreAdditionalBandwidth, error) {
	rawConfig := d.GetRawConfig()
	if rawConfig.IsNull() || !rawConfig.IsKnown() {
		return nil, nil
	}
	value := rawConfig.GetAttr("additional_bandwidth")
	if value.IsNull() || !value.IsKnown() {
		return nil, nil
	}

	entries := make([]kvstoreAdditionalBandwidth, 0, value.LengthInt())
	for it := value.ElementIterator(); it.Next(); {
		_, element := it.Element()
		if element.IsNull() || !element.IsKnown() {
			continue
		}
		nodeId := element.GetAttr("node_id")
		bandwidth := element.GetAttr("bandwidth")
		if bandwidth.IsNull() || !bandwidth.IsKnown() {
			// bandwidth is Optional only so that a new block never introduces a required
			// attribute; a block without it has no meaning and must not be sent as zero.
			return nil, fmt.Errorf("additional_bandwidth requires bandwidth in every block")
		}
		bandwidthValue, _ := bandwidth.AsBigFloat().Int64()
		entry := kvstoreAdditionalBandwidth{Bandwidth: int(bandwidthValue)}
		if !nodeId.IsNull() && nodeId.IsKnown() {
			entry.NodeId = strings.TrimSpace(nodeId.AsString())
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

// additionalBandwidthEntriesFromState returns the additional_bandwidth blocks tracked in the
// state. A plain refresh has no raw configuration, so the state is the authoritative source of the
// shards this resource manages.
//
// A standard instance tracks its single block without a node_id, so a block without one is kept
// here: dropping it would leave no managed shard and make the refresh mistake a managed instance
// for an import. Only an empty state means the resource has no managed block yet.
func additionalBandwidthEntriesFromState(d *schema.ResourceData) []kvstoreAdditionalBandwidth {
	value, ok := d.Get("additional_bandwidth").(*schema.Set)
	if !ok || value.Len() == 0 {
		return nil
	}
	entries := make([]kvstoreAdditionalBandwidth, 0, value.Len())
	for _, item := range value.List() {
		entry, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		nodeId := strings.TrimSpace(fmt.Sprint(entry["node_id"]))
		if nodeId == "<nil>" {
			nodeId = ""
		}
		bandwidth, _ := entry["bandwidth"].(int)
		entries = append(entries, kvstoreAdditionalBandwidth{
			NodeId:    nodeId,
			Bandwidth: bandwidth,
		})
	}
	return entries
}

// refreshAdditionalBandwidthState rewrites the additional_bandwidth blocks of a refresh from the
// data nodes returned by DescribeRoleZoneInfo. It reports false when the listing cannot describe a
// configured shard, so the caller keeps the previous state instead of writing a partial value.
func refreshAdditionalBandwidthState(d *schema.ResourceData, dataNodes []map[string]interface{}) bool {
	state, ok := buildAdditionalBandwidthState(dataNodes, additionalBandwidthEntriesFromState(d))
	if !ok {
		return false
	}
	d.Set("additional_bandwidth", state)
	return true
}

// additionalBandwidthNeedsShardNames reports whether any block names a data shard. A standard
// instance configures the bandwidth without a shard name, so its request needs no node listing.
func additionalBandwidthNeedsShardNames(entries []kvstoreAdditionalBandwidth) bool {
	for _, entry := range entries {
		if strings.TrimSpace(entry.NodeId) != "" {
			return true
		}
	}
	return false
}

// additionalBandwidthNodeExtras maps every data shard node of an instance to the additional
// bandwidth purchased for it (CurrentBandWidth minus the specification DefaultBandWidth). A shard
// with a replica reports both a master and a slave row, so the master row wins; read-only nodes
// (InsType 3) are not bandwidth targets.
func additionalBandwidthNodeExtras(nodes []map[string]interface{}) map[string]int64 {
	extras := make(map[string]int64)
	for _, node := range nodes {
		if fmt.Sprint(node["InsType"]) == "3" {
			continue
		}
		insName := strings.TrimSpace(fmt.Sprint(node["InsName"]))
		if insName == "" || insName == "<nil>" {
			continue
		}
		defaultBandWidth, _ := strconv.ParseInt(fmt.Sprint(node["DefaultBandWidth"]), 10, 64)
		currentBandWidth, _ := strconv.ParseInt(fmt.Sprint(node["CurrentBandWidth"]), 10, 64)
		extra := currentBandWidth - defaultBandWidth
		if extra < 0 {
			extra = 0
		}
		isMaster := strings.EqualFold(fmt.Sprint(node["Role"]), "master")
		if _, ok := extras[insName]; ok && !isMaster {
			continue
		}
		extras[insName] = extra
	}
	return extras
}

func additionalBandwidthShardNames(extras map[string]int64) []string {
	names := make([]string, 0, len(extras))
	for name := range extras {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// resolveAdditionalBandwidthShard maps a configured node_id to the shard node name reported by
// DescribeRoleZoneInfo. The node name (r-xxx-db-0) or its full topology ID (r-xxx-db-0#123456) can
// be given directly; a suffix such as db-0 is accepted when exactly one data shard ends with it, so
// a configuration does not have to repeat the instance ID. A token that matches no shard or several
// shards is rejected instead of being sent to an API that would adjust a different shard.
func resolveAdditionalBandwidthShard(shardNames []string, token string) (string, bool) {
	token = strings.TrimSpace(token)
	if idx := strings.Index(token, "#"); idx >= 0 {
		// EnableAdditionalBandwidth only accepts the shard node name; the topology ID carries an
		// internal node suffix after "#" that the API rejects.
		token = strings.TrimSpace(token[:idx])
	}
	if token == "" {
		return "", false
	}

	suffixMatches := make([]string, 0, 1)
	for _, shardName := range shardNames {
		if shardName == token {
			return shardName, true
		}
		if strings.HasSuffix(shardName, "-"+token) {
			suffixMatches = append(suffixMatches, shardName)
		}
	}
	if len(suffixMatches) == 1 {
		return suffixMatches[0], true
	}
	return "", false
}

func isClusterOrReadWriteSplitArchitecture(architectureType string) bool {
	architecture := strings.ToLower(strings.TrimSpace(architectureType))
	if architecture == "" {
		return false
	}
	return strings.Contains(architecture, "cluster") || strings.Contains(architecture, "split")
}

// buildAdditionalBandwidthRequest turns the configured blocks into the two vectors
// EnableAdditionalBandwidth expects. NodeId and Bandwidth are matched positionally by the API, so
// every bandwidth value is emitted next to its own shard and the pairing cannot drift; shards are
// ordered by node name to keep the request stable across plans.
//
// NodeId is only available on cluster and read/write splitting instances, where every block must
// name a data shard. A standard instance has a single data shard and takes Bandwidth alone, so its
// one block must not name a shard: the API does not accept NodeId for that architecture.
func buildAdditionalBandwidthRequest(architectureType string, shardCount int, shardNames []string, entries []kvstoreAdditionalBandwidth) (nodeIds string, bandwidths string, err error) {
	if len(entries) == 0 {
		return "", "", nil
	}

	if !isClusterOrReadWriteSplitArchitecture(architectureType) && shardCount <= 1 {
		if len(entries) > 1 {
			return "", "", fmt.Errorf("additional_bandwidth accepts a single block on a standard instance (architecture type: %s, shard count: %d)", architectureType, shardCount)
		}
		entry := entries[0]
		if strings.TrimSpace(entry.NodeId) != "" {
			return "", "", fmt.Errorf("additional_bandwidth node_id is not supported on a standard instance: EnableAdditionalBandwidth only accepts NodeId for cluster or read/write splitting instances (architecture type: %s, shard count: %d)", architectureType, shardCount)
		}
		if entry.Bandwidth < 0 {
			return "", "", fmt.Errorf("additional_bandwidth bandwidth must not be negative")
		}
		return "", strconv.Itoa(entry.Bandwidth), nil
	}

	byShard := make(map[string]int, len(entries))
	for _, entry := range entries {
		if strings.TrimSpace(entry.NodeId) == "" {
			return "", "", fmt.Errorf("additional_bandwidth node_id must be set on a cluster or read/write splitting instance (architecture type: %s, shard count: %d)", architectureType, shardCount)
		}
		if entry.Bandwidth < 0 {
			return "", "", fmt.Errorf("additional_bandwidth bandwidth must not be negative")
		}
		shard, ok := resolveAdditionalBandwidthShard(shardNames, entry.NodeId)
		if !ok {
			return "", "", fmt.Errorf("additional_bandwidth node_id %q does not match exactly one data shard of the instance (data shards: %s)", entry.NodeId, strings.Join(shardNames, ", "))
		}
		if _, exists := byShard[shard]; exists {
			return "", "", fmt.Errorf("additional_bandwidth configures the data shard %s more than once", shard)
		}
		byShard[shard] = entry.Bandwidth
	}

	shards := make([]string, 0, len(byShard))
	for shard := range byShard {
		shards = append(shards, shard)
	}
	sort.Strings(shards)

	nodeIdList := make([]string, 0, len(shards))
	bandwidthList := make([]string, 0, len(shards))
	for _, shard := range shards {
		nodeIdList = append(nodeIdList, shard)
		bandwidthList = append(bandwidthList, strconv.Itoa(byShard[shard]))
	}
	return strings.Join(nodeIdList, ","), strings.Join(bandwidthList, ","), nil
}

// buildAdditionalBandwidthState derives the additional_bandwidth blocks from the data nodes
// returned by DescribeRoleZoneInfo.
//
// The managed shards drive the result: only shards this resource tracks are reported, and their
// configured node_id is kept in the state, so an instance whose other shards still carry purchased
// bandwidth does not gain unmanaged entries and the configuration stays authoritative. Without
// managed shards (an import) every shard with purchased bandwidth is reported under its node name.
// A managed shard without a row means the listing is incomplete, and the caller keeps the previous
// state rather than reporting a smaller value.
func buildAdditionalBandwidthState(nodes []map[string]interface{}, managed []kvstoreAdditionalBandwidth) ([]map[string]interface{}, bool) {
	extras := additionalBandwidthNodeExtras(nodes)
	if len(extras) == 0 {
		return nil, false
	}

	if len(managed) == 0 {
		shards := make([]string, 0, len(extras))
		for _, shard := range additionalBandwidthShardNames(extras) {
			if extras[shard] > 0 {
				shards = append(shards, shard)
			}
		}
		state := make([]map[string]interface{}, 0, len(shards))
		for _, shard := range shards {
			state = append(state, map[string]interface{}{
				"node_id":   shard,
				"bandwidth": int(extras[shard]),
			})
		}
		return state, true
	}

	shardNames := additionalBandwidthShardNames(extras)
	state := make([]map[string]interface{}, 0, len(managed))
	for _, entry := range managed {
		shard := ""
		if strings.TrimSpace(entry.NodeId) == "" {
			// A standard instance takes the bandwidth without a shard name, so its single data
			// shard is the only possible target.
			if len(shardNames) != 1 {
				return nil, false
			}
			shard = shardNames[0]
		} else if resolved, ok := resolveAdditionalBandwidthShard(shardNames, entry.NodeId); ok {
			shard = resolved
		} else {
			return nil, false
		}
		state = append(state, map[string]interface{}{
			"node_id":   entry.NodeId,
			"bandwidth": int(extras[shard]),
		})
	}
	return state, true
}

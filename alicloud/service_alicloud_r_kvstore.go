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
	pageNumber := 1
	pageSize := 50
	for {
		request := map[string]interface{}{
			"InstanceId": id,
			"QueryType":  1,
			"PageNumber": pageNumber,
			"PageSize":   pageSize,
		}

		var response map[string]interface{}
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
			return nodes, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}

		pageNodeList, _ := jsonpath.Get("$.Node.NodeInfo", response)
		pageNodes, _ := pageNodeList.([]interface{})
		for _, v := range pageNodes {
			if node, ok := v.(map[string]interface{}); ok {
				nodes = append(nodes, node)
			}
		}

		totalCount := 0
		if v, ok := response["TotalCount"]; ok {
			totalCount, _ = strconv.Atoi(fmt.Sprint(v))
		}
		if len(pageNodes) == 0 || len(nodes) >= totalCount {
			break
		}
		pageNumber++
	}

	return nodes, nil
}

// buildAdditionalBandwidthState derives the additional_bandwidth value from the data nodes
// returned by DescribeRoleZoneInfo. The response contains one entry per node, so a shard with a
// replica reports both a master and a slave row: rows are grouped by shard (CustinsId) and the
// master row is used. Read-only nodes (InsType 3) are not bandwidth targets. Shards are ordered
// by the ordinal in their node name so the comma-separated value matches the shard order of a
// multi-shard request, and a value uniform across all shards collapses to a single number, which
// is how an All-scoped request is configured.
func buildAdditionalBandwidthState(nodes []map[string]interface{}) string {
	type shardBandwidth struct {
		order  int
		extra  int64
		master bool
	}

	shards := make(map[string]*shardBandwidth)
	for _, node := range nodes {
		if fmt.Sprint(node["InsType"]) == "3" {
			continue
		}
		shardId := fmt.Sprint(node["CustinsId"])
		if shardId == "" || shardId == "<nil>" {
			shardId = fmt.Sprint(node["InsName"])
		}
		if shardId == "" || shardId == "<nil>" {
			continue
		}

		defaultBandWidth, _ := strconv.ParseInt(fmt.Sprint(node["DefaultBandWidth"]), 10, 64)
		currentBandWidth, _ := strconv.ParseInt(fmt.Sprint(node["CurrentBandWidth"]), 10, 64)
		extra := currentBandWidth - defaultBandWidth
		if extra < 0 {
			extra = 0
		}

		isMaster := strings.EqualFold(fmt.Sprint(node["Role"]), "master")
		existing, ok := shards[shardId]
		if !ok || (isMaster && !existing.master) {
			shards[shardId] = &shardBandwidth{
				order:  shardOrdinal(fmt.Sprint(node["InsName"])),
				extra:  extra,
				master: isMaster,
			}
		}
	}
	if len(shards) == 0 {
		return ""
	}

	ordered := make([]*shardBandwidth, 0, len(shards))
	for _, shard := range shards {
		ordered = append(ordered, shard)
	}
	sort.SliceStable(ordered, func(i, j int) bool {
		return ordered[i].order < ordered[j].order
	})

	uniform := true
	for _, shard := range ordered[1:] {
		if shard.extra != ordered[0].extra {
			uniform = false
			break
		}
	}
	if uniform {
		return strconv.FormatInt(ordered[0].extra, 10)
	}

	values := make([]string, 0, len(ordered))
	for _, shard := range ordered {
		values = append(values, strconv.FormatInt(shard.extra, 10))
	}
	return strings.Join(values, ",")
}

// shardOrdinal extracts the trailing shard number from a node name such as r-xxx-db-1, so shards
// sort numerically instead of lexicographically (db-10 must follow db-9).
func shardOrdinal(nodeName string) int {
	if idx := strings.LastIndex(nodeName, "-"); idx >= 0 {
		if ordinal, err := strconv.Atoi(nodeName[idx+1:]); err == nil {
			return ordinal
		}
	}
	return 0
}

func isClusterOrReadWriteSplitArchitecture(architectureType string) bool {
	architecture := strings.ToLower(strings.TrimSpace(architectureType))
	if architecture == "" {
		return false
	}
	return strings.Contains(architecture, "cluster") || strings.Contains(architecture, "split")
}

// validateAdditionalBandwidthNodeId mirrors the EnableAdditionalBandwidth contract: Bandwidth may
// only be sent without NodeId on non-cluster, non-read/write-splitting instances. A shard count
// above one also identifies a cluster instance when the architecture field is unavailable.
func validateAdditionalBandwidthNodeId(architectureType string, shardCount int, additionalBandwidth, nodeId string) error {
	if strings.TrimSpace(additionalBandwidth) == "" || strings.TrimSpace(nodeId) != "" {
		return nil
	}
	if isClusterOrReadWriteSplitArchitecture(architectureType) || shardCount > 1 {
		return fmt.Errorf("additional_bandwidth_node_id must be set when additional_bandwidth is used on a cluster or read/write splitting instance (architecture type: %s, shard count: %d)", architectureType, shardCount)
	}
	return nil
}

// additionalBandwidthFromRawConfig returns the additional_bandwidth values the user actually
// configured. The fields are Optional+Computed, so d.Get() would also return values that were
// only read back from the API.
func additionalBandwidthFromRawConfig(d *schema.ResourceData) (additionalBandwidth, additionalBandwidthNodeId string) {
	rawConfig := d.GetRawConfig()
	if rawConfig.IsNull() || !rawConfig.IsKnown() {
		return "", ""
	}
	if v := rawConfig.GetAttr("additional_bandwidth"); !v.IsNull() && v.IsKnown() {
		additionalBandwidth = strings.TrimSpace(v.AsString())
	}
	if v := rawConfig.GetAttr("additional_bandwidth_node_id"); !v.IsNull() && v.IsKnown() {
		additionalBandwidthNodeId = strings.TrimSpace(v.AsString())
	}
	return additionalBandwidth, additionalBandwidthNodeId
}

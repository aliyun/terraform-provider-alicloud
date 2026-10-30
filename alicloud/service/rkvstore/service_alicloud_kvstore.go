package rkvstore

import (
	"log"
	"regexp"
	"strings"
	"time"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type KvstoreService struct {
	client *connectivity.AliyunClient
}

func NewKvstoreService(client *connectivity.AliyunClient) *KvstoreService {
	return &KvstoreService{client: client}
}

var KVstoreInstanceStatusCatcher = helper.Catcher{Reason: "OperationDenied.KVstoreInstanceStatus", RetryCount: 60, RetryWaitSeconds: 5}

func (s *KvstoreService) DescribeKVstoreInstance(id string) (*r_kvstore.DBInstanceAttribute, error) {
	instance := &r_kvstore.DBInstanceAttribute{}
	request := r_kvstore.CreateDescribeInstanceAttributeRequest()
	request.RegionId = s.client.RegionId
	request.InstanceId = id
	raw, err := s.client.WithRKvstoreClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.DescribeInstanceAttribute(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return instance, errs.WrapErrorf(errs.NotFoundErr("KVstoreInstance", id), errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return instance, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*r_kvstore.DescribeInstanceAttributeResponse)
	if len(response.Instances.DBInstanceAttribute) <= 0 {
		return instance, errs.WrapErrorf(errs.NotFoundErr("KVstoreInstance", id), errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
	}

	return &response.Instances.DBInstanceAttribute[0], nil
}

func (s *KvstoreService) DescribeKVstoreBackupPolicy(id string) (*r_kvstore.DescribeBackupPolicyResponse, error) {
	response := &r_kvstore.DescribeBackupPolicyResponse{}
	request := r_kvstore.CreateDescribeBackupPolicyRequest()
	request.RegionId = s.client.RegionId
	request.InstanceId = id

	var raw interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithRKvstoreClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.DescribeBackupPolicy(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"LockTimeout"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response, _ = raw.(*r_kvstore.DescribeBackupPolicyResponse)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return response, errs.WrapErrorf(errs.NotFoundErr("KVstore:BackupPolicy", id), errs.NotFoundWithResponse, response)
		}
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *KvstoreService) WaitForKVstoreInstance(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeKVstoreInstance(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object.InstanceStatus == string(status) {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.InstanceStatus, status, errs.ProviderERROR)
		}
	}
	return nil
}

func (s *KvstoreService) RdsKvstoreInstanceStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeKVstoreInstance(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.InstanceStatus == failState {
				return object, object.InstanceStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object.InstanceStatus))
			}
		}
		return object, object.InstanceStatus, nil
	}
}

func (s *KvstoreService) WaitForKVstoreInstanceVpcAuthMode(id string, status string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeKVstoreInstance(id)
		if err != nil && !errs.NotFoundError(err) {
			return err
		}
		if object.VpcAuthMode == string(status) {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.VpcAuthMode, status, errs.ProviderERROR)
		}
	}
	return nil
}

func (s *KvstoreService) DescribeParameters(id string) (*r_kvstore.DescribeParametersResponse, error) {
	response := &r_kvstore.DescribeParametersResponse{}
	request := r_kvstore.CreateDescribeParametersRequest()
	request.RegionId = s.client.RegionId
	request.DBInstanceId = id

	raw, err := s.client.WithRKvstoreClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.DescribeParameters(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return response, errs.WrapErrorf(errs.NotFoundErr("Parameters", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*r_kvstore.DescribeParametersResponse)
	return response, nil
}

func (s *KvstoreService) ModifyInstanceConfig(id string, config string) error {
	request := r_kvstore.CreateModifyInstanceConfigRequest()
	request.RegionId = s.client.RegionId
	request.InstanceId = id
	request.Config = config

	if err := s.WaitForKVstoreInstance(id, Normal, DefaultLongTimeout); err != nil {
		return errs.WrapError(err)
	}
	raw, err := s.client.WithRKvstoreClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.ModifyInstanceConfig(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *KvstoreService) setInstanceTags(d *schema.ResourceData) error {
	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		o := oraw.(map[string]interface{})
		n := nraw.(map[string]interface{})
		create, remove := s.diffTags(s.tagsFromMap(o), s.tagsFromMap(n))

		if len(remove) > 0 {
			var tagKey []string
			for _, v := range remove {
				tagKey = append(tagKey, v.Key)
			}
			request := r_kvstore.CreateUntagResourcesRequest()
			request.ResourceId = &[]string{d.Id()}
			request.ResourceType = strings.ToUpper(string(tags.TagResourceInstance))
			request.TagKey = &tagKey
			request.RegionId = s.client.RegionId
			raw, err := s.client.WithRKvstoreClient(func(client *r_kvstore.Client) (interface{}, error) {
				return client.UntagResources(request)
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
			}
			logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}

		if len(create) > 0 {
			request := r_kvstore.CreateTagResourcesRequest()
			request.ResourceId = &[]string{d.Id()}
			request.Tag = &create
			request.ResourceType = strings.ToUpper(string(tags.TagResourceInstance))
			request.RegionId = s.client.RegionId
			raw, err := s.client.WithRKvstoreClient(func(client *r_kvstore.Client) (interface{}, error) {
				return client.TagResources(request)
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
			}
			logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}

	}

	return nil
}

func (s *KvstoreService) tagsToMap(tags []r_kvstore.TagResource) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.TagKey] = t.TagValue
		}
	}
	return result
}

func (s *KvstoreService) tagsFromMap(m map[string]interface{}) []r_kvstore.TagResourcesTag {
	result := make([]r_kvstore.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, r_kvstore.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func (s *KvstoreService) ignoreTag(t r_kvstore.TagResource) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.TagKey)
		ok, _ := regexp.MatchString(v, t.TagValue)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific t %s (val: %s), ignoring.\n", t.TagKey, t.TagValue)
			return true
		}
	}
	return false
}

func (s *KvstoreService) diffTags(oldTags, newTags []r_kvstore.TagResourcesTag) ([]r_kvstore.TagResourcesTag, []r_kvstore.TagResourcesTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []r_kvstore.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return s.tagsFromMap(create), remove
}

func (s *KvstoreService) DescribeTags(resourceId string, resourceType tags.TagResourceType) (tags []r_kvstore.TagResource, err error) {
	request := r_kvstore.CreateListTagResourcesRequest()
	request.RegionId = s.client.RegionId
	request.ResourceType = strings.ToUpper(string(resourceType))
	request.ResourceId = &[]string{resourceId}
	raw, err := s.client.WithRKvstoreClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.ListTagResources(request)
	})
	if err != nil {
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, resourceId, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		return
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*r_kvstore.ListTagResourcesResponse)

	return response.TagResources.TagResource, nil
}

func (s *KvstoreService) WaitForKVstoreAccount(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeKVstoreAccount(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object != nil && object.AccountStatus == string(status) {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.AccountStatus, status, errs.ProviderERROR)
		}
	}
	return nil
}

func (s *KvstoreService) DescribeKVstoreAccount(id string) (*r_kvstore.Account, error) {
	ds := &r_kvstore.Account{}
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return ds, errs.WrapError(err)
	}
	request := r_kvstore.CreateDescribeAccountsRequest()
	request.RegionId = s.client.RegionId
	request.InstanceId = parts[0]
	request.AccountName = parts[1]
	invoker := helper.NewInvoker()
	invoker.AddCatcher(KVstoreInstanceStatusCatcher)
	var response *r_kvstore.DescribeAccountsResponse
	if err := invoker.Run(func() error {
		raw, err := s.client.WithRKvstoreClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.DescribeAccounts(request)
		})
		if err != nil {
			return err
		}

		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

		response, _ = raw.(*r_kvstore.DescribeAccountsResponse)
		return nil
	}); err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return ds, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return ds, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	if len(response.Accounts.Account) < 1 {
		return ds, errs.WrapErrorf(errs.NotFoundErr("KVstoreAccount", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return &response.Accounts.Account[0], nil
}

func (s *KvstoreService) DescribeKVstoreSecurityGroupId(id string) (*r_kvstore.DescribeSecurityGroupConfigurationResponse, error) {
	response := &r_kvstore.DescribeSecurityGroupConfigurationResponse{}
	request := r_kvstore.CreateDescribeSecurityGroupConfigurationRequest()
	request.RegionId = s.client.RegionId
	request.InstanceId = id
	if err := s.WaitForKVstoreInstance(id, Normal, DefaultLongTimeout); err != nil {
		return response, errs.WrapError(err)
	}
	raw, err := s.client.WithRKvstoreClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.DescribeSecurityGroupConfiguration(request)
	})
	if err != nil {
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*r_kvstore.DescribeSecurityGroupConfigurationResponse)

	return response, nil
}

func (s *KvstoreService) DescribeDBInstanceNetInfo(id string) (*r_kvstore.NetInfoItemsInDescribeDBInstanceNetInfo, error) {
	response := &r_kvstore.DescribeDBInstanceNetInfoResponse{}
	request := r_kvstore.CreateDescribeDBInstanceNetInfoRequest()
	request.RegionId = s.client.RegionId
	request.InstanceId = id
	if err := s.WaitForKVstoreInstance(id, Normal, DefaultLongTimeout); err != nil {
		return &response.NetInfoItems, errs.WrapError(err)
	}
	raw, err := s.client.WithRKvstoreClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.DescribeDBInstanceNetInfo(request)
	})
	if err != nil {
		return &response.NetInfoItems, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*r_kvstore.DescribeDBInstanceNetInfoResponse)

	return &response.NetInfoItems, nil
}

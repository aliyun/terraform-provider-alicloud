package redis

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Service wraps the R-kvstore (Redis/Tair) API calls used by the resources and
// data sources of this package.
type Service struct {
	client *connectivity.AliyunClient
}

// NewService is the exported constructor: the client field stays unexported, so
// callers outside the package (including the flat layout's acceptance tests)
// build values with it.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

// DescribeTairInstance <<< Encapsulated get interface for Redis TairInstance.
func (s *Service) DescribeTairInstance(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeInstanceAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = id
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("TairInstance", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Instances.DBInstanceAttribute[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Instances.DBInstanceAttribute[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("TairInstance", id), errs.NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}
func (s *Service) DescribeTairInstanceDescribeInstanceConfig(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeInstanceConfig"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = id
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("TairInstance", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}
func (s *Service) DescribeDescribeSecurityIps(id, securityIpGroupName string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeSecurityIps"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = id
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("TairInstance", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.SecurityIpGroups.SecurityIpGroup[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.SecurityIpGroups.SecurityIpGroup[*]", response)
	}

	if securityIpGroupName == "" {
		securityIpGroupName = "default"
	}

	for _, item := range v.([]interface{}) {
		securityIpGroup := item.(map[string]interface{})
		if securityIpGroup["SecurityIpGroupName"] == securityIpGroupName {
			return securityIpGroup, nil
		}
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("TairInstance", id), errs.NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}
func (s *Service) DescribeTairInstanceDescribeSecurityGroupConfiguration(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeSecurityGroupConfiguration"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = id
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"IncorrectDBInstanceState", "InvalidInstanceId.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("TairInstance", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Items.EcsSecurityGroupRelation[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Items.EcsSecurityGroupRelation[*]", response)
	}

	relations, ok := v.([]interface{})
	if !ok || len(relations) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("TairInstance", id), errs.NotFoundMsg, response)
	}

	// Join every EcsSecurityGroupRelation[*].SecurityGroupId in sorted order so state is
	// stable across API return-order fluctuations. Build a fresh map (never alias into
	// relations[0]) so we do not mutate the RPC response; the caller only reads RegionId
	// and SecurityGroupId, so those are the only fields we carry over.
	securityGroupIds := make([]string, 0, len(relations))
	for _, relation := range relations {
		rel, ok := relation.(map[string]interface{})
		if !ok {
			continue
		}
		securityGroupIds = append(securityGroupIds, fmt.Sprint(rel["SecurityGroupId"]))
	}
	sort.Strings(securityGroupIds)

	object = map[string]interface{}{
		"SecurityGroupId": strings.Join(securityGroupIds, ","),
	}
	if first, ok := relations[0].(map[string]interface{}); ok {
		if regionId, ok := first["RegionId"]; ok {
			object["RegionId"] = regionId
		}
	}

	return object, nil
}
func (s *Service) DescribeTairInstanceDescribeInstanceTDEStatus(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = id
	request["RegionId"] = client.RegionId
	action := "DescribeInstanceTDEStatus"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("TairInstance", id), errs.NotFoundMsg, response)
		}
		// Defense in depth: instances that do not support TDE reject DescribeInstanceTDEStatus
		// with a 400 InstanceType.NotSupport. Treat it as "TDE not supported" by returning an
		// empty object without error, so the read is not aborted.
		if errs.IsExpectedErrors(err, []string{"InstanceType.NotSupport"}) {
			return object, nil
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}
func (s *Service) DescribeTairInstanceDescribeInstanceSSL(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeInstanceSSL"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = id
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("TairInstance", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}
func (s *Service) DescribeTairInstanceDescribeIntranetAttribute(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeIntranetAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = id
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("TairInstance", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *Service) TairInstanceStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeTairInstance(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeTairInstance >>> Encapsulated.

// SetResourceTags <<< Encapsulated tag function for Redis.
func (s *Service) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var err error
		var action string
		client := s.client
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})

		added, removed := tags.ParsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !tags.IgnoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action = "UntagResources"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ResourceId.1"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ResourceType"] = resourceType
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}

			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, false)
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

		}

		if len(added) > 0 {
			action = "TagResources"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ResourceId.1"] = d.Id()
			request["RegionId"] = client.RegionId
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			request["ResourceType"] = resourceType
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, false)
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

		}
	}

	return nil
}

// SetResourceTags >>> tag function encapsulated.

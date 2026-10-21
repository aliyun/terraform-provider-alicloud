// Package polardb holds the SDKv2 implementation, service and tests for the
// PolarDB resources whose flat-layout code has been migrated into a product
// package. The flat alicloud package cannot be imported from here, so the
// Describe/Delete/state-refresh calls previously living on the shared
// PolarDBService for the application resource were extracted with it.
package polardb

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// Service wraps the PolarDB API calls used by the resources of this package.
type Service struct {
	client *connectivity.AliyunClient
}

// NewService is the exported constructor: the client field stays unexported, so
// callers outside the package (including the acceptance tests) build values
// with it.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

// DescribeApplicationAttribute reads a single PolarDB application by id.
func (s *Service) DescribeApplicationAttribute(id string) (object map[string]interface{}, err error) {
	action := "DescribeApplicationAttribute"
	request := map[string]interface{}{
		"ApplicationId": id,
	}
	var response map[string]interface{}
	client := s.client
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
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
		if errs.IsExpectedErrors(err, []string{"InvalidParameter", "ResourceNotFound", "InvalidApplicationId.NotFound"}) {
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

// DeleteApplication deletes a single PolarDB application by id.
func (s *Service) DeleteApplication(id string) (object map[string]interface{}, err error) {
	action := "DeleteApplication"
	request := map[string]interface{}{
		"ApplicationId": id,
	}
	var response map[string]interface{}
	client := s.client
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
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

// applicationStateRefreshFunc returns a StateRefreshFunc polling the
// application status. After the resource is deleted the API may surface
// several error shapes; those are all treated as "resource gone".
func (s *Service) applicationStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeApplicationAttribute(id)
		if err != nil {
			if errs.NotFoundError(err) || strings.Contains(err.Error(), "InvalidParameter") || strings.Contains(err.Error(), "ResourceNotFound") {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		for _, failState := range failStates {
			if object["Status"] == failState {
				return object, object["Status"].(string), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, object["Status"].(string), nil
	}
}

// describeVSwitch resolves the zone id of a VSwitch. The flat layout used the
// shared VpcService.DescribeVSwitch (vpc SDK + WithVpcClient); a service
// package must not import the flat layout, so this call is reimplemented on the
// generic RpcPost entry, keeping the same request semantics.
func (s *Service) describeVSwitch(id string) (zoneId string, err error) {
	action := "DescribeVSwitchAttributes"
	request := map[string]interface{}{
		"VSwitchId": id,
	}
	var response map[string]interface{}
	client := s.client
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("vpc", "2016-04-28", action, nil, request, false)
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
		if errs.IsExpectedErrors(err, []string{"InvalidVswitchID.NotFound"}) {
			return "", errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return "", errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.ZoneId", response)
	if err != nil {
		return "", errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ZoneId", response)
	}
	zoneId, _ = v.(string)
	return zoneId, nil
}

// ListApplicationTags lists the tags of a PolarDB application through
// ListTagResources (ResourceType="application"). It returns a TagKey->TagValue
// map; an application with no tags yields an empty map. A missing
// TagResource path (no tags) is not an error.
func (s *Service) ListApplicationTags(id string) (map[string]string, error) {
	action := "ListTagResources"
	request := map[string]interface{}{
		"ResourceType": "application",
		"ResourceId.1": id,
	}
	response, err := s.client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)
	result := make(map[string]string)
	v, err := jsonpath.Get("$.TagResources.TagResource", response)
	if err != nil || v == nil {
		return result, nil
	}
	list, ok := v.([]interface{})
	if !ok {
		return result, nil
	}
	for _, r := range list {
		item, ok := r.(map[string]interface{})
		if !ok {
			continue
		}
		key, _ := item["TagKey"].(string)
		value, _ := item["TagValue"].(string)
		if key != "" {
			result[key] = value
		}
	}
	return result, nil
}

// TagApplicationResources tags a PolarDB application through TagResources
// (ResourceType="application"). tagMap is the desired key/value set;
// TagResources overwrites existing keys and adds new ones.
func (s *Service) TagApplicationResources(id string, tagMap map[string]interface{}) error {
	action := "TagResources"
	request := map[string]interface{}{
		"ResourceType": "application",
		"ResourceId.1": id,
	}
	tags.ExpandTagsToMap(request, tags.ConvertTags(tagMap))
	response, err := s.client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)
	return nil
}

// UntagApplicationResources removes the given tag keys from a PolarDB
// application through UntagResources (ResourceType="application").
func (s *Service) UntagApplicationResources(id string, tagKeys []string) error {
	action := "UntagResources"
	request := map[string]interface{}{
		"ResourceType": "application",
		"ResourceId.1": id,
	}
	for i, key := range tagKeys {
		request[fmt.Sprintf("TagKey.%d", i+1)] = key
	}
	response, err := s.client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)
	return nil
}

package oos

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type OosServiceV2 struct {
	client *connectivity.AliyunClient
}

func NewOosServiceV2(client *connectivity.AliyunClient) *OosServiceV2 {
	return &OosServiceV2{client: client}
}

// DescribeOosPatchBaseline <<< Encapsulated get interface for Oos PatchBaseline.

func (s *OosServiceV2) DescribeOosPatchBaseline(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Name"] = id
	request["RegionId"] = client.RegionId
	action := "GetPatchBaseline"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("oos", "2019-06-01", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"EntityNotExists.PatchBaseline"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("PatchBaseline", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.PatchBaseline", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.PatchBaseline", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *OosServiceV2) OosPatchBaselineStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.OosPatchBaselineStateRefreshFuncWithApi(id, field, failStates, s.DescribeOosPatchBaseline)
}

func (s *OosServiceV2) OosPatchBaselineStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
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

// DescribeOosPatchBaseline >>> Encapsulated.

func (s *OosServiceV2) SetOssResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var action string
		var err error
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
			request["ResourceIds"] = "[\"" + d.Id() + "\"]"
			request["RegionId"] = client.RegionId
			request["ResourceType"] = resourceType
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			request["TagKeys"] = convert.ConvertListToJsonString(convert.ConvertListStringToListInterface(removedTagKeys))
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("oos", "2019-06-01", action, query, request, false)

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
			request["ResourceIds"] = "[\"" + d.Id() + "\"]"
			request["RegionId"] = client.RegionId
			tagsString, _ := convert.ConvertArrayObjectToJsonString(added)
			request["Tags"] = tagsString
			request["ResourceType"] = resourceType
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("oos", "2019-06-01", action, query, request, true)
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

// DescribeOosSecretParameter <<< Encapsulated get interface for Oos SecretParameter.

func (s *OosServiceV2) DescribeOosSecretParameter(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "GetSecretParameter"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Name"] = id
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("oos", "2019-06-01", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"EntityNotExists.Parameter"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("SecretParameter", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Parameter", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Parameter", response)
	}

	return v.(map[string]interface{}), nil
}
func (s *OosServiceV2) DescribeSecretParameterListTagResources(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "ListTagResources"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceIds"] = convert.ConvertObjectToJsonString(convert.ExpandSingletonToList(id))
	request["RegionId"] = client.RegionId

	request["ResourceType"] = "secretparameter"
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("oos", "2019-06-01", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"InvalidResourceId.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("SecretParameter", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

// WaitForOosSecretParameterTagsConverged polls ListTagResources until the tag bindings
// read back for the specified secret parameter match the expected tags. Tag bindings
// written through UpdateSecretParameter are propagated asynchronously to the tag
// storage behind ListTagResources, so an immediate read may miss them.
func (s *OosServiceV2) WaitForOosSecretParameterTagsConverged(id string, expectedTags map[string]interface{}, timeout time.Duration) error {
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err := retry.Retry(timeout, func() *retry.RetryError {
		object, err := s.DescribeSecretParameterListTagResources(id)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}

		tagsMaps, _ := jsonpath.Get("$.TagResources.TagResource", object)
		tagResources := make([]interface{}, 0)
		for _, tagResource := range convert.ConvertToInterfaceArray(tagsMaps) {
			if tagResourceMap, ok := tagResource.(map[string]interface{}); ok && fmt.Sprint(tagResourceMap["ResourceId"]) == id {
				tagResources = append(tagResources, tagResource)
			}
		}
		currentTags := tags.TagsToMap(tagResources)
		if len(currentTags) != len(expectedTags) {
			wait()
			return retry.RetryableError(errs.Error("The tags of OOS secret parameter %s have not been propagated yet, current: %v, expected: %v.", id, currentTags, expectedTags))
		}
		for key, value := range expectedTags {
			if fmt.Sprint(currentTags[key]) != fmt.Sprint(value) {
				wait()
				return retry.RetryableError(errs.Error("The tags of OOS secret parameter %s have not been propagated yet, current: %v, expected: %v.", id, currentTags, expectedTags))
			}
		}
		return nil
	})
	if err != nil {
		return errs.WrapError(err)
	}
	return nil
}

func (s *OosServiceV2) OosSecretParameterStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeOosSecretParameter(id)
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

// DescribeOosSecretParameter >>> Encapsulated.

// SetResourceTags <<< Encapsulated tag function for Oos.
func (s *OosServiceV2) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var action string
		var err error
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
			request["ResourceIds"] = convert.ConvertObjectToJsonString(convert.ExpandSingletonToList(d.Id()))
			request["RegionId"] = client.RegionId
			if v, ok := d.GetOk("tags"); ok {
				tagsTagKeyJsonPath, err := jsonpath.Get("$.tag_key", v)
				if err == nil && tagsTagKeyJsonPath != "" {
					request["TagKeys"] = tagsTagKeyJsonPath
				}
			}
			request["ResourceType"] = resourceType
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("oos", "2019-06-01", action, query, request, true)
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
			request["ResourceIds"] = convert.ConvertObjectToJsonString(convert.ExpandSingletonToList(d.Id()))
			request["RegionId"] = client.RegionId
			tagsJsonPath, err := jsonpath.Get("$", d.Get("tags"))
			if err == nil {
				request["Tags"] = tagsJsonPath
			}

			request["ResourceType"] = resourceType
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("oos", "2019-06-01", action, query, request, true)
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

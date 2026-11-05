// Package milvus. This file is generated automatically. Please do not modify it manually, thank you!
package milvus

import (
	"encoding/json"
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
	"github.com/tidwall/sjson"
)

// Service wraps the Milvus API calls used by the resources and data
// sources of this package.
type Service struct {
	client *connectivity.AliyunClient
}

// NewService is the exported constructor: the client field stays
// unexported, so callers outside the package (including the flat
// layout's acceptance tests) build values with it.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

// DescribeInstance <<< Encapsulated get interface for Milvus Instance.

func (s *Service) DescribeInstance(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	request = make(map[string]interface{})
	query = make(map[string]*string)
	query["instanceId"] = helper.StringPointer(id)
	query["RegionId"] = helper.StringPointer(client.RegionId)
	action := fmt.Sprintf("/webapi/instance/get")

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("milvus", "2023-10-12", action, query, nil, nil)

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

	v, err := jsonpath.Get("$.instance", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.instance", response)
	}

	currentStatus := v.(map[string]interface{})["status"]
	if fmt.Sprint(currentStatus) == "deleted" {
		return object, errs.WrapErrorf(errs.NotFoundErr("Instance", id), errs.NotFoundMsg, response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) InstanceStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.InstanceStateRefreshFuncWithApi(id, field, failStates, s.DescribeInstance)
}

func (s *Service) InstanceStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
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

// DescribeInstance >>> Encapsulated.

// SetResourceTags <<< Encapsulated tag function for Milvus.
func (s *Service) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var action string
		var err error
		client := s.client
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]*string)
		body := make(map[string]interface{})

		added, removed := tags.ParsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !tags.IgnoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action = fmt.Sprintf("/webapi/tags")
			request = make(map[string]interface{})
			query = make(map[string]*string)
			body = make(map[string]interface{})
			query["ResourceId"] = helper.StringPointer(convert.ConvertListToJsonString(convert.ExpandSingletonToList(d.Id())))
			query["RegionId"] = helper.StringPointer(client.RegionId)
			query["TagKey"] = helper.StringPointer(convert.ConvertListToJsonString(convert.ConvertListStringToListInterface(removedTagKeys)))
			query["ResourceType"] = helper.StringPointer(resourceType)
			body = request
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RoaDelete("milvus", "2023-10-12", action, query, nil, nil, true)
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
			action = fmt.Sprintf("/webapi/tags")
			request = make(map[string]interface{})
			query = make(map[string]*string)
			body = make(map[string]interface{})

			var tagList []interface{}
			if v, ok := d.GetOk("tags"); ok {
				tagsMap := v.(map[string]interface{})
				for key, value := range tagsMap {
					tagItem := map[string]interface{}{
						"Key":   key,
						"Value": value,
					}
					tagList = append(tagList, tagItem)
				}
			}
			request["Tag"] = tagList
			request["RegionId"] = client.RegionId
			request["ResourceType"] = resourceType
			jsonString := convert.ConvertObjectToJsonString(request)
			jsonString, _ = sjson.Set(jsonString, "ResourceId.0", d.Id())
			_ = json.Unmarshal([]byte(jsonString), &request)

			body = request
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RoaPost("milvus", "2023-10-12", action, query, nil, body, true)
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

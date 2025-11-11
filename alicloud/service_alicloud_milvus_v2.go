// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tidwall/sjson"
)

type MilvusServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeMilvusInstance <<< Encapsulated get interface for Milvus Instance.

func (s *MilvusServiceV2) DescribeMilvusInstance(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	request = make(map[string]interface{})
	query = make(map[string]*string)
	query["instanceId"] = StringPointer(id)
	query["RegionId"] = StringPointer(client.RegionId)
	action := fmt.Sprintf("/webapi/instance/get")

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("milvus", "2023-10-12", action, query, nil, nil)

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

	v, err := jsonpath.Get("$.instance", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.instance", response)
	}

	currentStatus := v.(map[string]interface{})["status"]
	if fmt.Sprint(currentStatus) == "deleted" {
		return object, WrapErrorf(NotFoundErr("Instance", id), NotFoundMsg, response)
	}

	return v.(map[string]interface{}), nil
}

func (s *MilvusServiceV2) MilvusInstanceStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.MilvusInstanceStateRefreshFuncWithApi(id, field, failStates, s.DescribeMilvusInstance)
}

func (s *MilvusServiceV2) MilvusInstanceStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
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
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeMilvusInstance >>> Encapsulated.

// SetResourceTags <<< Encapsulated tag function for Milvus.
func (s *MilvusServiceV2) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var action string
		var err error
		client := s.client
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]*string)
		body := make(map[string]interface{})

		added, removed := parsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action = fmt.Sprintf("/webapi/tags")
			request = make(map[string]interface{})
			query = make(map[string]*string)
			body = make(map[string]interface{})
			query["ResourceId"] = StringPointer(convertListToJsonString(expandSingletonToList(d.Id())))
			query["RegionId"] = StringPointer(client.RegionId)
			query["TagKey"] = StringPointer(convertListToJsonString(convertListStringToListInterface(removedTagKeys)))
			query["ResourceType"] = StringPointer(resourceType)
			body = request
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RoaDelete("milvus", "2023-10-12", action, query, nil, nil, true)
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
			jsonString := convertObjectToJsonString(request)
			jsonString, _ = sjson.Set(jsonString, "ResourceId.0", d.Id())
			_ = json.Unmarshal([]byte(jsonString), &request)

			body = request
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RoaPost("milvus", "2023-10-12", action, query, nil, body, true)
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

		}
	}

	return nil
}

// SetResourceTags >>> tag function encapsulated.

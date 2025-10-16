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

type StarRocksServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeStarRocksInstance <<< Encapsulated get interface for StarRocks Instance.

func (s *StarRocksServiceV2) DescribeStarRocksInstance(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	query["InstanceId"] = StringPointer(id)
	query["RegionId"] = StringPointer(client.RegionId)
	action := fmt.Sprintf("/webapi/starrocks/describeInstances")

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaPost("starrocks", "2022-10-19", action, query, nil, body, true)

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

	v, err := jsonpath.Get("$.Data[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("Instance", id), NotFoundMsg, response)
	}

	currentStatus := v.([]interface{})[0].(map[string]interface{})["InstanceStatus"]
	if fmt.Sprint(currentStatus) == "deleting" {
		return object, WrapErrorf(NotFoundErr("Instance", id), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *StarRocksServiceV2) StarRocksInstanceStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeStarRocksInstance(id)
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

// DescribeStarRocksInstance >>> Encapsulated.
// SetResourceTags <<< Encapsulated tag function for StarRocks.
func (s *StarRocksServiceV2) SetResourceTags(d *schema.ResourceData, resourceType string) error {
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
			query["ResourceType"] = StringPointer(resourceType)

			query["TagKey"] = StringPointer(convertListToJsonString(convertListStringToListInterface(removedTagKeys)))
			body = request
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RoaDelete("starrocks", "2022-10-19", action, query, nil, nil, true)
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

			request["ResourceType"] = resourceType
			request["RegionId"] = client.RegionId
			count := 1
			tagsMaps := make([]map[string]interface{}, 0)
			for key, value := range added {
				tagsMap := make(map[string]interface{})
				tagsMap["Key"] = key
				tagsMap["Value"] = value
				tagsMaps = append(tagsMaps, tagsMap)
				count++
			}
			request["Tag"] = tagsMaps

			jsonString := convertObjectToJsonString(request)
			jsonString, _ = sjson.Set(jsonString, "ResourceId.0", d.Id())
			_ = json.Unmarshal([]byte(jsonString), &request)

			body = request
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RoaPost("starrocks", "2022-10-19", action, query, nil, body, true)
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

// DescribeStarRocksNodeGroup <<< Encapsulated get interface for StarRocks NodeGroup.

func (s *StarRocksServiceV2) DescribeStarRocksNodeGroup(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	query["ClusterId"] = StringPointer(parts[0])
	query["RegionId"] = StringPointer(client.RegionId)
	jsonString := convertObjectToJsonString(query)
	jsonString, _ = sjson.Set(jsonString, "nodeGroupIds.0", parts[1])
	_ = json.Unmarshal([]byte(jsonString), &request)
	body := request

	action := fmt.Sprintf("/webapi/nodegroup/describeNodeGroups")

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaPost("starrocks", "2022-10-19", action, query, nil, body, true)

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

	v, err := jsonpath.Get("$.Data[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("NodeGroup", id), NotFoundMsg, response)
	}

	currentStatus := v.([]interface{})[0].(map[string]interface{})["Status"]
	if fmt.Sprint(currentStatus) == "DELETED" {
		return object, WrapErrorf(NotFoundErr("NodeGroup", id), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *StarRocksServiceV2) StarRocksNodeGroupStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.StarRocksNodeGroupStateRefreshFuncWithApi(id, field, failStates, s.DescribeStarRocksNodeGroup)
}

func (s *StarRocksServiceV2) StarRocksNodeGroupStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
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

// DescribeStarRocksNodeGroup >>> Encapsulated.

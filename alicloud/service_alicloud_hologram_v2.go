package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type HologramServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeHologramInstance <<< Encapsulated get interface for Hologram Instance.

func (s *HologramServiceV2) DescribeHologramInstance(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	instanceId := id
	action := fmt.Sprintf("/api/v1/instances/%s", instanceId)
	request = make(map[string]interface{})
	body = make(map[string]interface{})
	query = make(map[string]*string)
	request["instanceId"] = id

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("Hologram", "2022-06-01", action, query, nil, body)

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
		if IsExpectedErrors(err, []string{"resource not exists failed"}) {
			return object, WrapErrorf(NotFoundErr("Instance", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Instance", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.body.Instance", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *HologramServiceV2) HologramInstanceStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeHologramInstance(id)
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeHologramInstance >>> Encapsulated.

// SetResourceTags <<< Encapsulated tag function for Hologram.
func (s *HologramServiceV2) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var err error
		var action string
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
			action = fmt.Sprintf("/api/v1/tag/unTag")
			request = make(map[string]interface{})
			query = make(map[string]*string)
			body = make(map[string]interface{})
			request["resourceIds"] = expandSingletonToList(d.Id())
			request["tagKeys"] = convertListStringToListInterface(removedTagKeys)

			body = request
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RoaPost("Hologram", "2022-06-01", action, query, nil, body, true)

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
			action = fmt.Sprintf("/api/v1/tag")
			request = make(map[string]interface{})
			query = make(map[string]*string)
			body = make(map[string]interface{})
			request["resourceIds"] = expandSingletonToList(d.Id())
			count := 1
			tagsMaps := make([]map[string]interface{}, 0)
			for key, value := range added {
				tagsMap := make(map[string]interface{})
				tagsMap["key"] = key
				tagsMap["value"] = value
				tagsMaps = append(tagsMaps, tagsMap)
				count++
			}
			request["tags"] = tagsMaps

			body = request
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RoaPost("Hologram", "2022-06-01", action, query, nil, body, true)

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

// SetResourceTags >>> tag function encapsulated.

package alicloud

import (
	"fmt"
	"strings"
	"time"

	roa "github.com/alibabacloud-go/tea-roa/client"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type RocketmqServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeRocketmqInstance <<< Encapsulated get interface for Rocketmq Instance.

func (s *RocketmqServiceV2) DescribeRocketmqInstance(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	instanceId := id
	action := fmt.Sprintf("/instances/%s", instanceId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	request["instanceId"] = id

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("RocketMQ", "2022-08-01", action, query, nil, nil)
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
		if IsExpectedErrors(err, []string{"Instance.NotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Instance", id)), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	response = response["body"].(map[string]interface{})

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.data", response)
	}

	currentStatus := v.(map[string]interface{})["status"]
	if currentStatus == "RELEASED" {
		return object, WrapErrorf(Error(GetNotFoundMessage("Instance", id)), NotFoundMsg, response)
	}

	return v.(map[string]interface{}), nil
}
func (s *RocketmqServiceV2) DescribeGetInstanceAccount(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	instanceId := id
	action := fmt.Sprintf("/instances/%s/account", instanceId)
	conn, err := client.NewRocketmqClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	request["instanceId"] = id

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2022-08-01"), nil, StringPointer("GET"), StringPointer("AK"), StringPointer(action), query, nil, nil, &runtime)

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
		addDebug(action, response, request)
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	response = response["body"].(map[string]interface{})

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *RocketmqServiceV2) RocketmqInstanceStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeRocketmqInstance(id)
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

// DescribeRocketmqInstance >>> Encapsulated.

// SetResourceTags <<< Encapsulated tag function for Rocketmq.
func (s *RocketmqServiceV2) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var err error
		var action string
		var conn *roa.Client
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
			action = fmt.Sprintf("/resourceTag/delete")
			conn, err = client.NewRocketmqClient()
			if err != nil {
				return WrapError(err)
			}
			request = make(map[string]interface{})
			query = make(map[string]*string)
			body = make(map[string]interface{})
			resourceIdString, _ := convertArrayObjectToJsonString(expandSingletonToList(d.Id()))
			query["resourceId"] = StringPointer(resourceIdString)
			query["regionId"] = StringPointer(client.RegionId)
			query["tagKey"] = StringPointer(convertListToJsonString(convertListStringToListInterface(removedTagKeys)))

			query["resourceType"] = StringPointer("instance")
			body = request
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer("2022-08-01"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
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
			action = fmt.Sprintf("/resourceTag/create")
			conn, err = client.NewRocketmqClient()
			if err != nil {
				return WrapError(err)
			}
			request = make(map[string]interface{})
			query = make(map[string]*string)
			body = make(map[string]interface{})
			resourceIdString, _ := convertArrayObjectToJsonString(expandSingletonToList(d.Id()))
			query["resourceId"] = StringPointer(resourceIdString)
			query["regionId"] = StringPointer(client.RegionId)
			query["resourceType"] = StringPointer("instance")
			count := 1
			tagsMaps := make([]map[string]interface{}, 0)
			for key, value := range added {
				tagsMap := make(map[string]interface{})
				tagsMap["key"] = key
				tagsMap["value"] = value
				tagsMaps = append(tagsMaps, tagsMap)
				count++
			}
			tagsString, _ := convertArrayObjectToJsonString(tagsMaps)
			query["tag"] = StringPointer(tagsString)

			body = request
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer("2022-08-01"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
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
	}

	return nil
}

// SetResourceTags >>> tag function encapsulated.

// DescribeRocketmqConsumerGroup <<< Encapsulated get interface for Rocketmq ConsumerGroup.

func (s *RocketmqServiceV2) DescribeRocketmqConsumerGroup(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	consumerGroupId := parts[1]
	instanceId := parts[0]
	action := fmt.Sprintf("/instances/%s/consumerGroups/%s", instanceId, consumerGroupId)
	conn, err := client.NewRocketmqClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	body = make(map[string]interface{})
	query = make(map[string]*string)

	body = request
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2022-08-01"), nil, StringPointer("GET"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)

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
		if IsExpectedErrors(err, []string{"ConsumerGroup.NotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("ConsumerGroup", id)), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.body.data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.body.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *RocketmqServiceV2) RocketmqConsumerGroupStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeRocketmqConsumerGroup(id)
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

// DescribeRocketmqConsumerGroup >>> Encapsulated.

// DescribeRocketmqTopic <<< Encapsulated get interface for Rocketmq Topic.

func (s *RocketmqServiceV2) DescribeRocketmqTopic(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	instanceId := parts[0]
	topicName := parts[1]
	action := fmt.Sprintf("/instances/%s/topics/%s", instanceId, topicName)
	conn, err := client.NewRocketmqClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	body = make(map[string]interface{})
	query = make(map[string]*string)

	body = request
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2022-08-01"), nil, StringPointer("GET"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)

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
		if IsExpectedErrors(err, []string{"Topic.NotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Topic", id)), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.body.data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.body.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *RocketmqServiceV2) RocketmqTopicStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeRocketmqTopic(id)
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

// DescribeRocketmqTopic >>> Encapsulated.

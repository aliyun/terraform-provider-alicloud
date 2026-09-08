package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type EventBridgeServiceV2 struct {
	client *connectivity.AliyunClient
}

func (s *EventBridgeServiceV2) DescribeEventBridgeConnection(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "GetConnection"
	client := s.client

	request := map[string]interface{}{
		"ConnectionName": id,
	}

	idExist := false
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"ConnectionNotExist"}) {
			return object, WrapErrorf(NotFoundErr("EventBridge:Connection", id), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.Data.Connections", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data.Connections", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, WrapErrorf(NotFoundErr("EventBridge:Connection", id), NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["ConnectionName"]) == id {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, WrapErrorf(NotFoundErr("EventBridge:Connection", id), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *EventBridgeServiceV2) DescribeEventBridgeApiDestination(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "GetApiDestination"
	client := s.client
	request := map[string]interface{}{
		"ApiDestinationName": id,
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, false)
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
		if IsExpectedErrors(err, []string{"ApiDestinationNotExist"}) {
			return object, WrapErrorf(NotFoundErr("EventBridge:ApiDestination", id), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

// DescribeEventBridgeEventSource <<< Encapsulated get interface for EventBridge EventSource.

func (s *EventBridgeServiceV2) DescribeEventBridgeEventSource(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Limit"] = PageSizeLarge

	action := "ListUserDefinedEventSources"

	idExist := false
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("eventbridge", "2020-04-01", action, query, request, true)
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

		resp, err := jsonpath.Get("$.Data.EventSourceList", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data.EventSourceList", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, WrapErrorf(NotFoundErr("EventBridge:EventSource", id), NotFoundWithResponse, response)
		}

		for _, v := range resp.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["Name"]) == id {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}

		if nextToken, ok := response["Data"].(map[string]interface{})["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	if !idExist {
		return object, WrapErrorf(NotFoundErr("EventBridge:EventSource", id), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *EventBridgeServiceV2) EventBridgeEventSourceStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.EventBridgeEventSourceStateRefreshFuncWithApi(id, field, failStates, s.DescribeEventBridgeEventSource)
}

func (s *EventBridgeServiceV2) EventBridgeEventSourceStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
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

// DescribeEventBridgeEventSource >>> Encapsulated.

// DescribeEventBridgeEventStreaming <<< Encapsulated get interface for EventBridge EventStreaming.
func (s *EventBridgeServiceV2) DescribeEventBridgeEventStreaming(id string) (object map[string]interface{}, err error) {
	client := s.client
	action := "GetEventStreaming"
	request := map[string]interface{}{
		"EventStreamingName": id,
	}
	var response map[string]interface{}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, true)
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
		if IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"EventStreamingNotExisted"}) {
			return object, WrapErrorf(NotFoundErr("EventBridge:EventStreaming", id), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *EventBridgeServiceV2) DescribeEventBridgeEventStreamingTags(id string) (object map[string]interface{}, err error) {
	client := s.client
	action := "ListTagResources"
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"ResourceType": "eventstreaming",
		"ResourceId.1": id,
	}
	var response map[string]interface{}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, true)
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
		if IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"ResourceNotExist"}) {
			return map[string]interface{}{}, nil
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *EventBridgeServiceV2) EventBridgeEventStreamingStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.EventBridgeEventStreamingStateRefreshFuncWithApi(id, field, failStates, s.DescribeEventBridgeEventStreaming)
}

func (s *EventBridgeServiceV2) EventBridgeEventStreamingStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
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

// SetResourceTags <<< Encapsulated tag function for EventBridge.
func (s *EventBridgeServiceV2) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var action string
		var err error
		client := s.client
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})

		added, removed := parsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action = "UntagResources"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ResourceId.1"] = d.Id()
			request["RegionId"] = client.RegionId
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}

			request["ResourceType"] = resourceType
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("eventbridge", "2020-04-01", action, query, request, true)
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
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("eventbridge", "2020-04-01", action, query, request, true)
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

// updateEventBridgeEventStreamingStatus drives the event stream to the status
// expected in the configuration through StartEventStreaming/PauseEventStreaming.
func (s *EventBridgeServiceV2) updateEventBridgeEventStreamingStatus(d *schema.ResourceData, timeout time.Duration) error {
	client := s.client
	target := d.Get("status").(string)
	if target == "" {
		return nil
	}

	object, err := s.DescribeEventBridgeEventStreaming(d.Id())
	if err != nil {
		return WrapError(err)
	}

	current := fmt.Sprint(object["Status"])
	if current == target {
		return nil
	}

	var action string
	switch target {
	case "RUNNING":
		action = "StartEventStreaming"
	case "PAUSED":
		action = "PauseEventStreaming"
	default:
		return nil
	}

	request := map[string]interface{}{
		"EventStreamingName": d.Id(),
	}

	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(timeout), func() *resource.RetryError {
		response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, true)
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

	if fmt.Sprint(response["Code"]) != "Success" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	pending := []string{"READY", "STARTING", "PAUSED"}
	failStates := []string{"STARTING_FAILED", "RUNNING_FAILED"}
	if target == "PAUSED" {
		pending = []string{"RUNNING", "STARTING", "READY"}
		failStates = []string{}
	}

	stateConf := BuildStateConf(pending, []string{target}, timeout, 10*time.Second, s.EventBridgeEventStreamingStateRefreshFunc(d.Id(), "$.Status", failStates))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

// DescribeEventBridgeEventStreaming >>> Encapsulated.

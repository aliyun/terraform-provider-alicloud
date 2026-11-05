package eventbridge

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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// NewEventbridgeService is the exported constructor for the legacy v1
// EventbridgeService, used by acceptance tests in the external test package.
func NewEventbridgeService(client *connectivity.AliyunClient) *EventbridgeService {
	return &EventbridgeService{client: client}
}

// NewEventBridgeServiceV2 is the exported constructor for EventBridgeServiceV2.
func NewEventBridgeServiceV2(client *connectivity.AliyunClient) *EventBridgeServiceV2 {
	return &EventBridgeServiceV2{client: client}
}

type EventbridgeService struct {
	client *connectivity.AliyunClient
}

func (s *EventbridgeService) DescribeEventBridgeEventBus(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetEventBus"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"EventBusName": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, true)
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
		if errs.IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"EventBusNotExist"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("EventBridge:EventBus", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *EventbridgeService) DescribeEventBridgeRule(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	action := "GetRule"

	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"EventBusName": parts[0],
		"RuleName":     parts[1],
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, true)
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
		if errs.IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"EventBusNotExist", "EventRuleNotExisted"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("EventBridge:Rule", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *EventbridgeService) EventBridgeRuleStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEventBridgeRule(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}

		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *EventbridgeService) DescribeEventBridgeServiceLinkedRole(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "CheckRoleForProduct"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"ProductName": id,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("eventbridge", "2020-04-01", action, request, nil)
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

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *EventbridgeService) CheckRoleForProductRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEventBridgeServiceLinkedRole(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["CheckPass"]) == failState {
				return object, fmt.Sprint(object["CheckPass"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["CheckPass"])))
			}
		}
		return object, fmt.Sprint(object["CheckPass"]), nil
	}
}

func (s *EventbridgeService) DescribeEventBridgeEventSource(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListUserDefinedEventSources"

	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"Limit":    helper.PageSizeLarge,
	}

	idExist := false
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, true)
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

		resp, err := jsonpath.Get("$.Data.EventSourceList", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data.EventSourceList", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, errs.WrapErrorf(errs.NotFoundErr("EventBridge:EventSource", id), errs.NotFoundWithResponse, response)
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
		return object, errs.WrapErrorf(errs.NotFoundErr("EventBridge:EventSource", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

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
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"ConnectionNotExist"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("EventBridge:Connection", id), errs.NotFoundWithResponse, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.Data.Connections", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data.Connections", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("EventBridge:Connection", id), errs.NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["ConnectionName"]) == id {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("EventBridge:Connection", id), errs.NotFoundWithResponse, response)
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

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, false)
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
		if errs.IsExpectedErrors(err, []string{"ApiDestinationNotExist"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("EventBridge:ApiDestination", id), errs.NotFoundWithResponse, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
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
	request["Limit"] = helper.PageSizeLarge

	action := "ListUserDefinedEventSources"

	idExist := false
	for {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("eventbridge", "2020-04-01", action, query, request, true)
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

		resp, err := jsonpath.Get("$.Data.EventSourceList", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data.EventSourceList", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, errs.WrapErrorf(errs.NotFoundErr("EventBridge:EventSource", id), errs.NotFoundWithResponse, response)
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
		return object, errs.WrapErrorf(errs.NotFoundErr("EventBridge:EventSource", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *EventBridgeServiceV2) EventBridgeEventSourceStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.EventBridgeEventSourceStateRefreshFuncWithApi(id, field, failStates, s.DescribeEventBridgeEventSource)
}

func (s *EventBridgeServiceV2) EventBridgeEventSourceStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
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

// DescribeEventBridgeEventSource >>> Encapsulated.

// DescribeEventBridgeEventStreaming <<< Encapsulated get interface for EventBridge EventStreaming.
func (s *EventBridgeServiceV2) DescribeEventBridgeEventStreaming(id string) (object map[string]interface{}, err error) {
	client := s.client
	action := "GetEventStreaming"
	request := map[string]interface{}{
		"EventStreamingName": id,
	}
	var response map[string]interface{}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, true)
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
		if errs.IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"EventStreamingNotExisted"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("EventBridge:EventStreaming", id), errs.NotFoundWithResponse, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
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

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, true)
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
		if errs.IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"ResourceNotExist"}) {
			return map[string]interface{}{}, nil
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *EventBridgeServiceV2) EventBridgeEventStreamingStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.EventBridgeEventStreamingStateRefreshFuncWithApi(id, field, failStates, s.DescribeEventBridgeEventStreaming)
}

func (s *EventBridgeServiceV2) EventBridgeEventStreamingStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
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

// SetResourceTags <<< Encapsulated tag function for EventBridge.
func (s *EventBridgeServiceV2) SetResourceTags(d *schema.ResourceData, resourceType string) error {
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
			request["ResourceId.1"] = d.Id()
			request["RegionId"] = client.RegionId
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}

			request["ResourceType"] = resourceType
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("eventbridge", "2020-04-01", action, query, request, true)
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
				response, err = client.RpcPost("eventbridge", "2020-04-01", action, query, request, true)
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
		return errs.WrapError(err)
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
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(timeout), func() *retry.RetryError {
		response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, true)
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

	if fmt.Sprint(response["Code"]) != "Success" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	pending := []string{"READY", "STARTING", "PAUSED"}
	failStates := []string{"STARTING_FAILED", "RUNNING_FAILED"}
	if target == "PAUSED" {
		pending = []string{"RUNNING", "STARTING", "READY"}
		failStates = []string{}
	}

	stateConf := helper.BuildStateConf(pending, []string{target}, timeout, 10*time.Second, s.EventBridgeEventStreamingStateRefreshFunc(d.Id(), "$.Status", failStates))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}

// DescribeEventBridgeEventStreaming >>> Encapsulated.

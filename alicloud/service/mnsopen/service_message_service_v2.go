package mnsopen

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfconvert "github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	tferrs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	tfhelper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	tflogs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	tftags "github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	jsonata "github.com/blues/jsonata-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/sjson"
)

type MessageServiceServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeMessageServiceQueue <<< Encapsulated get interface for MessageService Queue.

func (s *MessageServiceServiceV2) DescribeMessageServiceQueue(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["QueueName"] = id
	request["RegionId"] = client.RegionId
	action := "GetQueueAttributes"

	wait := tfhelper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Mns-open", "2022-01-19", action, query, request, true)

		if err != nil {
			if tfhelper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	tflogs.AddDebug(action, response, request)
	if err != nil {
		if tferrs.IsExpectedErrors(err, []string{"ResourceNotFound:QueueName", "RESOURCE_NOT_FOUND", "QueueNotExist"}) {
			return object, tferrs.WrapErrorf(tferrs.NotFoundErr("Queue", id), tferrs.NotFoundMsg, response)
		}
		return object, tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, id, action, tferrs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, tferrs.WrapErrorf(err, tferrs.FailedGetAttributeMsg, id, "$.Data", response)
	}

	return v.(map[string]interface{}), nil
}

// DescribeMessageServiceQueue >>> Encapsulated.

// DescribeMessageServiceTopic <<< Encapsulated get interface for MessageService Topic.

func (s *MessageServiceServiceV2) DescribeMessageServiceTopic(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "GetTopicAttributes"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["TopicName"] = id

	wait := tfhelper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Mns-open", "2022-01-19", action, query, request, true)

		if err != nil {
			if tfhelper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	tflogs.AddDebug(action, response, request)
	if err != nil {
		if tferrs.IsExpectedErrors(err, []string{"ResourceNotFound:TopicName", "TopicNotExist"}) {
			return object, tferrs.WrapErrorf(tferrs.NotFoundErr("Topic", id), tferrs.NotFoundMsg, response)
		}
		return object, tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, id, action, tferrs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, tferrs.WrapErrorf(err, tferrs.FailedGetAttributeMsg, id, "$.Data", response)
	}

	return v.(map[string]interface{}), nil
}

// DescribeMessageServiceTopic >>> Encapsulated.

// SetResourceTags <<< Encapsulated tag function for MessageService.
func (s *MessageServiceServiceV2) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var err error
		var action string
		client := s.client
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})

		added, removed := tftags.ParsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !tftags.IgnoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action = "UntagResources"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["RegionId"] = client.RegionId
			request["ResourceType"] = resourceType
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}

			jsonString := tfconvert.ConvertObjectToJsonString(request)
			jsonString, _ = sjson.Set(jsonString, "ResourceId.0", d.Id())
			err = json.Unmarshal([]byte(jsonString), &request)
			if err != nil {
				return tferrs.WrapError(err)
			}

			wait := tfhelper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("Mns-open", "2022-01-19", action, query, request, true)
				if err != nil {
					if tfhelper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			tflogs.AddDebug(action, response, request)
			if err != nil {
				return tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, d.Id(), action, tferrs.AlibabaCloudSdkGoERROR)
			}

		}

		if len(added) > 0 {
			action = "TagResources"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["RegionId"] = client.RegionId
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			request["ResourceType"] = resourceType
			jsonString := tfconvert.ConvertObjectToJsonString(request)
			jsonString, _ = sjson.Set(jsonString, "ResourceId.0", d.Id())
			err = json.Unmarshal([]byte(jsonString), &request)
			if err != nil {
				return tferrs.WrapError(err)
			}

			wait := tfhelper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("Mns-open", "2022-01-19", action, query, request, true)
				if err != nil {
					if tfhelper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			tflogs.AddDebug(action, response, request)
			if err != nil {
				return tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, d.Id(), action, tferrs.AlibabaCloudSdkGoERROR)
			}

		}
	}

	return nil
}

// SetResourceTags >>> tag function encapsulated.

// DescribeMessageServiceEndpoint <<< Encapsulated get interface for MessageService Endpoint.

func (s *MessageServiceServiceV2) DescribeMessageServiceEndpoint(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["EndpointType"] = id
	request["RegionId"] = client.RegionId
	action := "GetEndpointAttribute"

	wait := tfhelper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Mns-open", "2022-01-19", action, query, request, true)

		if err != nil {
			if tfhelper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	tflogs.AddDebug(action, response, request)
	if err != nil {
		return object, tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, id, action, tferrs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, tferrs.WrapErrorf(err, tferrs.FailedGetAttributeMsg, id, "$.Data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *MessageServiceServiceV2) MessageServiceEndpointStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeMessageServiceEndpoint(id)
		if err != nil {
			if tferrs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", tferrs.WrapError(err)
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
				return object, currentStatus, tferrs.WrapError(tferrs.Error(tferrs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeMessageServiceEndpoint >>> Encapsulated.

// DescribeMessageServiceEndpointAcl <<< Encapsulated get interface for MessageService EndpointAcl.

func (s *MessageServiceServiceV2) DescribeMessageServiceEndpointAcl(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = tferrs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}

	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["EndpointType"] = parts[0]
	request["RegionId"] = client.RegionId
	action := "GetEndpointAttribute"

	idExist := false
	wait := tfhelper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Mns-open", "2022-01-19", action, query, request, true)

		if err != nil {
			if tfhelper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	tflogs.AddDebug(action, response, request)

	if err != nil {
		return object, tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, id, action, tferrs.AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.Data.CidrList[*]", response)
	if err != nil {
		return object, tferrs.WrapErrorf(err, tferrs.FailedGetAttributeMsg, id, "$.Data.CidrList[*]", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, tferrs.WrapErrorf(tferrs.NotFoundErr("EndpointAcl", id), tferrs.NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["AclStrategy"]) == parts[1] && fmt.Sprint(v.(map[string]interface{})["Cidr"]) == parts[2] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, tferrs.WrapErrorf(tferrs.NotFoundErr("EndpointAcl", id), tferrs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *MessageServiceServiceV2) MessageServiceEndpointAclStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeMessageServiceEndpointAcl(id)
		if err != nil {
			if tferrs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", tferrs.WrapError(err)
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
				return object, currentStatus, tferrs.WrapError(tferrs.Error(tferrs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeMessageServiceEndpointAcl >>> Encapsulated.

// DescribeMessageServiceSubscription <<< Encapsulated get interface for MessageService Subscription.

func (s *MessageServiceServiceV2) DescribeMessageServiceSubscription(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = tferrs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SubscriptionName"] = parts[1]
	request["TopicName"] = parts[0]
	request["RegionId"] = client.RegionId
	action := "GetSubscriptionAttributes"

	wait := tfhelper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Mns-open", "2022-01-19", action, query, request, true)

		if err != nil {
			if tfhelper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	tflogs.AddDebug(action, response, request)
	if err != nil {
		if tferrs.IsExpectedErrors(err, []string{"SubscriptionNotExist", "ResourceNotFound:TopicName/subscriptionName", "TopicNotExist"}) {
			return object, tferrs.WrapErrorf(tferrs.NotFoundErr("Subscription", id), tferrs.NotFoundMsg, response)
		}
		return object, tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, id, action, tferrs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, tferrs.WrapErrorf(err, tferrs.FailedGetAttributeMsg, id, "$.Data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *MessageServiceServiceV2) MessageServiceSubscriptionStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeMessageServiceSubscription(id)
		if err != nil {
			if tferrs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", tferrs.WrapError(err)
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
				return object, currentStatus, tferrs.WrapError(tferrs.Error(tferrs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeMessageServiceSubscription >>> Encapsulated.

// DescribeMessageServiceService <<< Encapsulated get interface for MessageService Service.

func (s *MessageServiceServiceV2) DescribeMessageServiceService(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	request["SubscriptionType"] = "PayAsYouGo"
	var endpoint string
	request["ProductCode"] = "mns"
	request["ProductType"] = ""
	if client.IsInternationalAccount() {
		request["ProductType"] = ""
	}
	action := "QueryAvailableInstances"

	wait := tfhelper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)

		if err != nil {
			if tfhelper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			if !client.IsInternationalAccount() && tferrs.IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductCode"] = "mns"
				request["ProductType"] = ""
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	tflogs.AddDebug(action, response, request)
	if err != nil {
		if tferrs.IsExpectedErrors(err, []string{"Expired"}) {
			return object, tferrs.WrapErrorf(tferrs.NotFoundErr("Service", id), tferrs.NotFoundMsg, response)
		}
		return object, tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, id, action, tferrs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *MessageServiceServiceV2) MessageServiceServiceStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeMessageServiceService(id)
		if err != nil {
			if tferrs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", tferrs.WrapError(err)
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
				return object, currentStatus, tferrs.WrapError(tferrs.Error(tferrs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeMessageServiceService >>> Encapsulated.

// DescribeMessageServiceEventRule <<< Encapsulated get interface for MessageService EventRule.

func (s *MessageServiceServiceV2) DescribeMessageServiceEventRule(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RuleName"] = id
	request["RegionId"] = client.RegionId
	request["ProductName"] = "oss"
	action := "GetEventRule"

	wait := tfhelper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Mns-open", "2022-01-19", action, query, request, true)

		if err != nil {
			if tfhelper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	tflogs.AddDebug(action, response, request)
	if err != nil {
		if tferrs.IsExpectedErrors(err, []string{"RESOURCE_NOT_FOUND"}) {
			return object, tferrs.WrapErrorf(tferrs.NotFoundErr("EventRule", id), tferrs.NotFoundMsg, response)
		}
		return object, tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, id, action, tferrs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, tferrs.WrapErrorf(err, tferrs.FailedGetAttributeMsg, id, "$.Data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *MessageServiceServiceV2) MessageServiceEventRuleStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeMessageServiceEventRule(id)
		if err != nil {
			if tferrs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", tferrs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)
		if field == "$.Endpoint.EndpointValue" {
			e := jsonata.MustCompile("$split($.Endpoint.EndpointValue, \"/queues/\")[-1]")
			v, _ = e.Eval(object)
			currentStatus = fmt.Sprint(v)
		}

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, tferrs.WrapError(tferrs.Error(tferrs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeMessageServiceEventRule >>> Encapsulated.

// DescribeMessageServiceAccountLogging <<< Encapsulated get interface for MessageService AccountLogging.

func (s *MessageServiceServiceV2) DescribeMessageServiceAccountLogging(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	action := "GetAccountAttributes"

	wait := tfhelper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Mns-open", "2022-01-19", action, query, request, true)

		if err != nil {
			if tfhelper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	tflogs.AddDebug(action, response, request)
	if err != nil {
		return object, tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, id, action, tferrs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, tferrs.WrapErrorf(err, tferrs.FailedGetAttributeMsg, id, "$.Data", response)
	}

	return v.(map[string]interface{}), nil
}

// DescribeMessageServiceAccountLogging >>> Encapsulated.

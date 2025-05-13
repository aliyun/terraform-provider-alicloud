package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type AmqpServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeAmqpInstance <<< Encapsulated get interface for Amqp Instance.

func (s *AmqpServiceV2) DescribeAmqpInstance(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "GetInstance"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["InstanceId"] = id

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcGet("amqp-open", "2019-12-12", action, query, request)

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

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceNotfound"}) {
			return object, WrapErrorf(NotFoundErr("Instance", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}

	currentStatus := v.(map[string]interface{})["Status"]
	if currentStatus == "" {
		return object, WrapErrorf(NotFoundErr("Instance", id), NotFoundMsg, response)
	}

	return v.(map[string]interface{}), nil
}
func (s *AmqpServiceV2) DescribeQueryAvailableInstances(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var endpoint string
	var query map[string]interface{}
	action := "QueryAvailableInstances"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["InstanceIDs"] = id

	request["ProductCode"] = "ons"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductType"] = "ons_onsproxy_public_intl"
				endpoint = connectivity.BssOpenAPIEndpointInternational
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

	v, err := jsonpath.Get("$.Data.InstanceList[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data.InstanceList[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("Instance", id), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *AmqpServiceV2) AmqpInstanceStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAmqpInstance(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
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

// DescribeAmqpInstance >>> Encapsulated.

// DescribeAmqpExchange <<< Encapsulated get interface for Amqp Exchange.

func (s *AmqpServiceV2) DescribeAmqpExchange(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["InstanceId"] = parts[0]
	query["VirtualHost"] = parts[1]
	query["RegionId"] = client.RegionId
	query["MaxResults"] = 100
	action := "ListExchanges"

	for {

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = client.RpcGet("amqp-open", "2019-12-12", action, query, request)

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
			if IsExpectedErrors(err, []string{"107", "501"}) {
				return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}

		v, err := jsonpath.Get("$.Data", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
		}

		item := v.(map[string]interface{})
		exchangeNames, err := jsonpath.Get("$.Exchanges[*]", item)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Exchanges[*]", item)
		}
		for _, vv := range exchangeNames.([]interface{}) {
			if vv.(map[string]interface{})["Name"] == parts[2] {
				return vv.(map[string]interface{}), nil
			}
		}

		if nextToken, ok := v.(map[string]interface{})["NextToken"].(string); ok && nextToken != "" {
			query["NextToken"] = nextToken
		} else {
			break
		}
	}
	return object, WrapErrorf(NotFoundErr("Exchange", id), NotFoundMsg, response)
}

func (s *AmqpServiceV2) AmqpExchangeStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAmqpExchange(id)
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

// DescribeAmqpExchange >>> Encapsulated.

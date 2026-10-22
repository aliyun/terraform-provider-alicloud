package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type AmqpServiceV2 struct {
	client *connectivity.AliyunClient
}

// convertAmqpInstanceInstanceTypeResponse, convertAmqpInstanceRenewalDurationUnitResponse
// and convertAmqpInstanceRenewalDurationUnitRequest were originally defined in the flat
// resource_alicloud_amqp_instance.go. They are shared by the AMQP instances data source
// and the RocketMQ instance resource, so they stay in the root package after the
// alicloud_amqp_instance resource migrated to service/amqp.

func convertAmqpInstanceInstanceTypeResponse(source interface{}) interface{} {
	switch source {
	case "PROFESSIONAL":
		return "professional"
	case "ENTERPRISE":
		return "enterprise"
	case "VIP":
		return "vip"
	}
	return source
}
func convertAmqpInstanceRenewalDurationUnitResponse(source interface{}) interface{} {
	switch source {
	case "M":
		return "Month"
	case "Y":
		return "Year"
	}
	return source
}
func convertAmqpInstanceRenewalDurationUnitRequest(source interface{}) interface{} {
	switch source {
	case "Month":
		return "M"
	case "Year":
		return "Y"
	}
	return source
}

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
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = client.RpcGet("amqp-open", "2019-12-12", action, query, request)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
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

func (s *AmqpServiceV2) AmqpExchangeStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
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

// DescribeAmqpVirtualHost <<< Encapsulated get interface for Amqp Virtual Host.

func (s *AmqpServiceV2) DescribeAmqpVirtualHost(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	request = make(map[string]interface{})
	request["InstanceId"] = parts[0]
	request["MaxResults"] = PageSizeLarge
	request["RegionId"] = client.RegionId
	action := "ListVirtualHosts"

	idExist := false
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = client.RpcGet("amqp-open", "2019-12-12", action, request, nil)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			if IsExpectedErrors(err, []string{"InstanceNotExist", "InstanceIdNotExist"}) {
				return object, WrapErrorf(NotFoundErr("VirtualHost", id), NotFoundMsg, response)
			}
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Data.VirtualHosts", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data.VirtualHosts", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, WrapErrorf(NotFoundErr("VirtualHost", id), NotFoundMsg, response)
		}

		for _, v := range resp.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["Name"]) == parts[1] {
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
		return object, WrapErrorf(NotFoundErr("VirtualHost", id), NotFoundMsg, response)
	}

	return object, nil
}

func (s *AmqpServiceV2) AmqpVirtualHostStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAmqpVirtualHost(id)
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

// DescribeAmqpVirtualHost >>> Encapsulated.

// DescribeAmqpOpenSourcePermission <<< Encapsulated get interface for Amqp OpenSourcePermission.

func (s *AmqpServiceV2) DescribeAmqpOpenSourcePermission(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
		return nil, err
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = parts[2]
	request["UserName"] = parts[0]
	request["RegionId"] = client.RegionId
	action := "ListOpenSourcePermissions"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("amqp-open", "2019-12-12", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InstanceNotExist"}) {
			return object, WrapErrorf(NotFoundErr("OpenSourcePermission", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("OpenSourcePermission", id), NotFoundMsg, response)
	}

	result, _ := v.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if fmt.Sprint(item["Vhost"]) != parts[1] {
			continue
		}
		return item, nil
	}
	return object, WrapErrorf(NotFoundErr("OpenSourcePermission", id), NotFoundMsg, response)
}

func (s *AmqpServiceV2) AmqpOpenSourcePermissionStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.AmqpOpenSourcePermissionStateRefreshFuncWithApi(id, field, failStates, s.DescribeAmqpOpenSourcePermission)
}

func (s *AmqpServiceV2) AmqpOpenSourcePermissionStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
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

// DescribeAmqpOpenSourcePermission >>> Encapsulated.
// DescribeAmqpOpenSourceAccount <<< Encapsulated get interface for Amqp OpenSourceAccount.

func (s *AmqpServiceV2) DescribeAmqpOpenSourceAccount(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
		return nil, err
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = parts[1]
	request["RegionId"] = client.RegionId
	action := "ListOpenSourceAccounts"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("amqp-open", "2019-12-12", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InstanceNotExist"}) {
			return object, WrapErrorf(NotFoundErr("OpenSourceAccount", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("OpenSourceAccount", id), NotFoundMsg, response)
	}

	result, _ := v.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if fmt.Sprint(item["Name"]) != parts[0] {
			continue
		}
		return item, nil
	}
	return object, WrapErrorf(NotFoundErr("OpenSourceAccount", id), NotFoundMsg, response)
}

func (s *AmqpServiceV2) AmqpOpenSourceAccountStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.AmqpOpenSourceAccountStateRefreshFuncWithApi(id, field, failStates, s.DescribeAmqpOpenSourceAccount)
}

func (s *AmqpServiceV2) AmqpOpenSourceAccountStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
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

// DescribeAmqpOpenSourceAccount >>> Encapsulated.

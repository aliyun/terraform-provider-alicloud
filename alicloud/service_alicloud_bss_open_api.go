package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type BssOpenApiService struct {
	client *connectivity.AliyunClient
}

const ModulesSizeLimit = 50

func (s *BssOpenApiService) QueryAvailableInstances(id, instanceRegion, productCode, productType, productCodeIntl, productTypeIntl string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	var endpoint string
	action := "QueryAvailableInstances"
	request := map[string]interface{}{
		"ProductCode": productCode,
		"ProductType": productType,
	}
	if client.IsInternationalAccount() {
		request["ProductCode"] = productCodeIntl
		request["ProductType"] = productTypeIntl
	}
	if id != "" {
		request["InstanceIDs"] = id
	}
	if instanceRegion != "" {
		request["Region"] = instanceRegion
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductCode"] = productCodeIntl
				request["ProductType"] = productTypeIntl
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		resp, _ := jsonpath.Get("$.Data.InstanceList", response)
		if len(resp.([]interface{})) < 1 {
			request["ProductCode"] = productCodeIntl
			if productTypeIntl != "" {
				request["ProductType"] = productTypeIntl
			}
			endpoint = connectivity.BssOpenAPIEndpointInternational
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data.InstanceList", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data.InstanceList", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr(productCode+"Instance", id), NotFoundWithResponse, response)
	} else if id != "" {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["InstanceID"]) != id {
			return object, WrapErrorf(NotFoundErr(productCode+"Instance", id), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *BssOpenApiService) QueryAvailableInstancesWithoutProductType(id, instanceRegion, productCode, productCodeIntl string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	var endpoint string
	action := "QueryAvailableInstances"
	request := map[string]interface{}{
		"InstanceIDs": id,
		"ProductCode": productCode,
	}
	if client.IsInternationalAccount() {
		request["ProductCode"] = productCodeIntl
	}
	if instanceRegion != "" {
		request["Region"] = instanceRegion
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductCode"] = productCodeIntl
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		resp, err := jsonpath.Get("$.Data.InstanceList", response)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if !client.IsInternationalAccount() && len(resp.([]interface{})) < 1 {
			request["ProductCode"] = productCodeIntl
			endpoint = connectivity.BssOpenAPIEndpointInternational
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data.InstanceList", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data.InstanceList", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr(productCode+"Instance", id), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["InstanceID"]) != id {
			return object, WrapErrorf(NotFoundErr(productCode+"Instance", id), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *BssOpenApiService) QueryAvailableInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "QueryAvailableInstances"
	request := map[string]interface{}{
		"InstanceIDs": id,
		"ProductCode": "vipcloudfw",
		"ProductType": "vipcloudfw",
	}
	if client.IsInternationalAccount() {
		request["ProductCode"] = "cfw"
		request["ProductType"] = "cfw_pre_intl"
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("BssOpenApi", "2017-12-14", action, nil, request, true)
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

	v, err := jsonpath.Get("$.Data.InstanceList", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data.InstanceList", response)
	}

	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("CloudFirewall", id), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["InstanceID"]) != id {
			return object, WrapErrorf(NotFoundErr("CloudFirewall", id), NotFoundWithResponse, response)
		}
	}

	object = v.([]interface{})[0].(map[string]interface{})

	return object, nil
}

func (s *BssOpenApiService) DescribeCloudFirewallInstanceOrderDetail(orderId string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	var endpoint string
	action := "GetOrderDetail"
	request := map[string]interface{}{
		"OrderId": orderId,
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, orderId, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data.OrderList.Order", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, orderId, "$.Data.OrderList.Order", response)
	}

	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("CloudFirewallOrderId", orderId), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["OrderId"]) != orderId {
			return object, WrapErrorf(NotFoundErr("CloudFirewallOrderId", orderId), NotFoundWithResponse, response)
		}
	}

	object = v.([]interface{})[0].(map[string]interface{})

	return object, nil
}

func (s *BssOpenApiService) CloudFirewallInstanceOrderDetailStateRefreshFunc(orderId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCloudFirewallInstanceOrderDetail(orderId)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["PaymentStatus"].(string) == failState {
				return object, object["PaymentStatus"].(string), WrapError(Error(FailedToReachTargetStatus, object["PaymentStatus"].(string)))
			}
		}

		return object, object["PaymentStatus"].(string), nil
	}
}

func (s *BssOpenApiService) QueryAvailableInstanceList(instanceRegion, productCode, productType, productCodeIntl, productTypeIntl string) (object []interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	var endpoint string
	action := "QueryAvailableInstances"
	request := map[string]interface{}{
		"ProductCode":      productCode,
		"ProductType":      productType,
		"SubscriptionType": "Subscription",
	}
	if s.client.IsInternationalAccount() {
		request = map[string]interface{}{
			"ProductCode": productCodeIntl,
			"ProductType": productTypeIntl,
		}
	}
	if instanceRegion != "" {
		request["Region"] = instanceRegion
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductType"] = productTypeIntl
				request["ProductCode"] = productCodeIntl
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapError(err)
	}
	v, err := jsonpath.Get("$.Data.InstanceList", response)
	if err != nil {
		return object, WrapError(err)
	}
	return v.([]interface{}), nil
}
func (s *BssOpenApiService) GetInstanceTypePrice(productCode, productType, paymentType string, modules []map[string]interface{}) (object []float64, err error) {
	client := s.client
	var response map[string]interface{}
	var endpoint string
	var priceList []float64
	action := "GetSubscriptionPrice"
	request := map[string]interface{}{
		"OrderType":        "NewOrder",
		"SubscriptionType": "Subscription",
		"Region":           client.RegionId,
		"ProductCode":      productCode,
		"ProductType":      productType,
	}
	if paymentType == "Subscription" {
		request["ServicePeriodQuantity"] = 1
		request["ServicePeriodUnit"] = "Month"
		request["Quantity"] = 1
	}
	if paymentType == "PayAsYouGo" {
		action = "GetPayAsYouGoPrice"
		request["SubscriptionType"] = "PayAsYouGo"
	}
	moduleLength := len(modules)
	for {
		if len(modules) < ModulesSizeLimit {
			request["ModuleList"] = modules
		} else {
			tmp := modules[:ModulesSizeLimit]
			modules = modules[ModulesSizeLimit:]
			moduleLength = len(tmp)
			request["ModuleList"] = &tmp
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
					endpoint = connectivity.BssOpenAPIEndpointInternational
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return object, WrapError(err)
		}
		v, err := jsonpath.Get("$.Data.ModuleDetails.ModuleDetail", response)
		if err != nil {
			return object, WrapError(err)
		}
		for _, o := range v.([]interface{}) {
			priceList = append(priceList, o.(map[string]interface{})["OriginalCost"].(float64))
		}
		if moduleLength < ModulesSizeLimit {
			break
		}
	}
	return priceList, nil
}

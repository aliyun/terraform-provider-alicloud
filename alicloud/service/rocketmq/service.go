package rocketmq

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Service struct {
	client *connectivity.AliyunClient
}

func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

// DescribeInstance <<< Encapsulated get interface for Rocketmq Instance.

func (s *Service) DescribeInstance(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	instanceId := id
	request = make(map[string]interface{})
	query = make(map[string]*string)
	request["instanceId"] = id

	action := fmt.Sprintf("/instances/%s", instanceId)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("RocketMQ", "2022-08-01", action, query, nil, nil)

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
		if errs.IsExpectedErrors(err, []string{"Instance.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Instance", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.data", response)
	}

	currentStatus := v.(map[string]interface{})["status"]
	if currentStatus == "RELEASED" {
		return object, errs.WrapErrorf(errs.NotFoundErr("Instance", id), errs.NotFoundMsg, response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) DescribeInstanceAccount(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	instanceId := id
	action := fmt.Sprintf("/instances/%s/account", instanceId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	request["instanceId"] = id

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("RocketMQ", "2022-08-01", action, query, nil, nil)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		logs.AddDebug(action, response, request)
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) InstanceStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeInstance(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
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

func (s *Service) DescribeInstanceIpWhitelist(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	instanceId := id
	request = make(map[string]interface{})
	query = make(map[string]*string)
	request["instanceId"] = id

	action := fmt.Sprintf("/instances/%s/ip/whitelists", instanceId)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("RocketMQ", "2022-08-01", action, query, nil, nil)

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
		if errs.IsExpectedErrors(err, []string{"InstanceIpOrCidr.NotExisted", "Instance.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Instance", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

// DescribeInstance >>> Encapsulated.

// SetResourceTags <<< Encapsulated tag function for Rocketmq.
func (s *Service) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var err error
		var action string
		client := s.client
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]*string)
		body := make(map[string]interface{})

		added, removed := tags.ParsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !tags.IgnoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action = fmt.Sprintf("/resourceTag/delete")
			request = make(map[string]interface{})
			query = make(map[string]*string)
			body = make(map[string]interface{})
			resourceIdString, _ := convert.ConvertArrayObjectToJsonString(convert.ExpandSingletonToList(d.Id()))
			query["resourceId"] = helper.StringPointer(resourceIdString)
			query["regionId"] = helper.StringPointer(client.RegionId)
			query["tagKey"] = helper.StringPointer(convert.ConvertListToJsonString(convert.ConvertListStringToListInterface(removedTagKeys)))

			query["resourceType"] = helper.StringPointer("instance")
			body = request
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RoaDelete("RocketMQ", "2022-08-01", action, query, nil, body, false)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				logs.AddDebug(action, response, request)
				return nil
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}

		}

		if len(added) > 0 {
			action = fmt.Sprintf("/resourceTag/create")
			request = make(map[string]interface{})
			query = make(map[string]*string)
			body = make(map[string]interface{})
			resourceIdString, _ := convert.ConvertArrayObjectToJsonString(convert.ExpandSingletonToList(d.Id()))
			query["resourceId"] = helper.StringPointer(resourceIdString)
			query["regionId"] = helper.StringPointer(client.RegionId)
			query["resourceType"] = helper.StringPointer("instance")
			count := 1
			tagsMaps := make([]map[string]interface{}, 0)
			for key, value := range added {
				tagsMap := make(map[string]interface{})
				tagsMap["key"] = key
				tagsMap["value"] = value
				tagsMaps = append(tagsMaps, tagsMap)
				count++
			}
			tagsString, _ := convert.ConvertArrayObjectToJsonString(tagsMaps)
			query["tag"] = helper.StringPointer(tagsString)

			body = request
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RoaPost("RocketMQ", "2022-08-01", action, query, nil, body, false)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				logs.AddDebug(action, response, request)
				return nil
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}

		}
	}

	return nil
}

// SetResourceTags >>> tag function encapsulated.

// QueryAvailableInstancesWithoutProductType resolves the billing information of a
// single instance without narrowing the query by product type.
func (s *Service) QueryAvailableInstancesWithoutProductType(id, instanceRegion, productCode, productCodeIntl string) (object map[string]interface{}, err error) {
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
	// Region is intentionally not set: QueryAvailableInstances is a global billing
	// API keyed by InstanceIDs. Passing an international region (e.g. na-south-1) for
	// a domestic-station account makes the domestic BSS return an empty list, which
	// triggers a blind fallback to the international endpoint and fails with
	// "account not exists". Omitting Region lets the domestic BSS resolve the instance.

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, endpoint)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			if !client.IsInternationalAccount() && errs.IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductCode"] = productCodeIntl
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		resp, err := jsonpath.Get("$.Data.InstanceList", response)
		if err != nil {
			return retry.NonRetryableError(err)
		}
		if !client.IsInternationalAccount() && len(resp.([]interface{})) < 1 {
			request["ProductCode"] = productCodeIntl
			endpoint = connectivity.BssOpenAPIEndpointInternational
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, endpoint)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data.InstanceList", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data.InstanceList", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr(productCode+"Instance", id), errs.NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["InstanceID"]) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr(productCode+"Instance", id), errs.NotFoundWithResponse, response)
		}
	}

	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

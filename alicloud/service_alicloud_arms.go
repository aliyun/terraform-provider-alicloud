package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type ArmsService struct {
	client *connectivity.AliyunClient
}

func (s *ArmsService) DescribeArmsAlertContact(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "SearchAlertContact"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"ContactIds": convertListToJsonString([]interface{}{id}),
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
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
	v, err := jsonpath.Get("$.PageBean.Contacts", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.PageBean.Contacts", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("ARMS", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["ContactId"]) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("ARMS", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *ArmsService) DescribeArmsAlertContactGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "SearchAlertContactGroup"
	request := map[string]interface{}{
		"RegionId":        s.client.RegionId,
		"ContactGroupIds": convertListToJsonString([]interface{}{id}),
		"IsDetail":        "true",
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
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
	v, err := jsonpath.Get("$.ContactGroups", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ContactGroups", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("ARMS", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["ContactGroupId"]) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("ARMS", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *ArmsService) DescribeArmsAlertRobot(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeIMRobots"
	request := map[string]interface{}{
		"RobotIds": id,
		"Page":     1,
		"Size":     PageSizeXLarge,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
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
	v, err := jsonpath.Get("$.PageBean.AlertIMRobots", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.PageBean.AlertIMRobots", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("ARMS", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["RobotId"]) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("ARMS", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *ArmsService) DescribeArmsDispatchRule(id string) (object map[string]interface{}, err error) {
	client := s.client

	request := map[string]interface{}{
		"Id":       id,
		"RegionId": s.client.RegionId,
	}

	var response map[string]interface{}
	action := "DescribeDispatchRule"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.RpcPost("ARMS", "2019-08-08", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.DispatchRule", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.DispatchRule", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *ArmsService) DescribeArmsPrometheusAlertRule(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListPrometheusAlertRules"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"ClusterId": parts[0],
	}
	idExist := false
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
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
	v, err := jsonpath.Get("$.PrometheusAlertRules", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.PrometheusAlertRules", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("ARMS", id)), NotFoundWithResponse, response)
	}
	for _, v := range v.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["AlertId"]) == parts[1] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("ARMS", id)), NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *ArmsService) ListArmsNotificationPolicies(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListNotificationPolicies"
	request := map[string]interface{}{
		"Page":     1,
		"Size":     PageSizeXLarge,
		"IsDetail": true,
		"RegionId": s.client.RegionId,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
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
	v, err := jsonpath.Get("$.PageBean.NotificationPolicies", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.PageBean.NotificationPolicies", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("ARMS", id)), NotFoundWithResponse, response)
	}
	for _, v := range v.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["Id"]) == id {
			return v.(map[string]interface{}), nil
		}
	}
	return object, WrapErrorf(Error(GetNotFoundMessage("ARMS", id)), NotFoundWithResponse, response)
}

func (s *ArmsService) ArmsDispatchRuleStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeArmsDispatchRule(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}
		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *ArmsService) DescribeArmsPrometheus(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "GetPrometheusInstance"

	client := s.client

	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"ClusterId": id,
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
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

	if fmt.Sprint(response["Code"]) == "404" {
		return object, WrapErrorf(Error(GetNotFoundMessage("Arms:Prometheus", id)), NotFoundWithResponse, response)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *ArmsService) ListTagResources(id string, resourceType string) (object interface{}, err error) {
	client := s.client
	action := "ListTagResources"

	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"ResourceType": resourceType,
	}

	resourceIdNum := strings.Count(id, ":")

	switch resourceIdNum {
	case 0:
		request["ResourceId.1"] = id
	case 1:
		parts, err := ParseResourceId(id, 2)
		if err != nil {
			return object, WrapError(err)
		}
		request["ResourceId.1"] = parts[resourceIdNum]
	}

	tags := make([]interface{}, 0)
	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			v, err := jsonpath.Get("$.TagResources", response)
			if err != nil {
				return resource.NonRetryableError(WrapErrorf(err, FailedGetAttributeMsg, id, "$.TagResources", response))
			}

			if v != nil {
				tags = append(tags, v.([]interface{})...)
			}

			return nil
		})
		if err != nil {
			err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
			return
		}
		if response["NextToken"] == nil {
			break
		}
		request["NextToken"] = response["NextToken"]
	}

	return tags, nil
}

func (s *ArmsService) SetResourceTags(d *schema.ResourceData, resourceType string) error {

	resourceIdNum := strings.Count(d.Id(), ":")

	if d.HasChange("tags") {
		added, removed := parsingTags(d)
		client := s.client

		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}

		if len(removedTagKeys) > 0 {
			action := "UntagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
			}

			switch resourceIdNum {
			case 0:
				request["ResourceId.1"] = d.Id()
			case 1:
				parts, err := ParseResourceId(d.Id(), 2)
				if err != nil {
					return WrapError(err)
				}
				request["ResourceId.1"] = parts[resourceIdNum]
			}

			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := client.RpcPost("ARMS", "2019-08-08", action, nil, request, false)
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
			action := "TagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
			}

			switch resourceIdNum {
			case 0:
				request["ResourceId.1"] = d.Id()
			case 1:
				parts, err := ParseResourceId(d.Id(), 2)
				if err != nil {
					return WrapError(err)
				}
				request["ResourceId.1"] = parts[resourceIdNum]
			}

			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}
			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := client.RpcPost("ARMS", "2019-08-08", action, nil, request, false)
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

func (s *ArmsService) DescribeArmsIntegrationExporter(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "GetPrometheusIntegration"

	client := s.client

	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":        s.client.RegionId,
		"ClusterId":       parts[0],
		"IntegrationType": parts[1],
		"InstanceId":      parts[2],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
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

	if fmt.Sprint(response["Code"]) == "404" {
		return object, WrapErrorf(Error(GetNotFoundMessage("Arms:IntegrationExporter", id)), NotFoundWithResponse, response)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *ArmsService) DescribeArmsRemoteWrite(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "GetPrometheusRemoteWrite"

	client := s.client

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":        s.client.RegionId,
		"ClusterId":       parts[0],
		"RemoteWriteName": parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
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

	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	if fmt.Sprint(response["Code"]) == "404" {
		return object, WrapErrorf(Error(GetNotFoundMessage("Arms:RemoteWrite", id)), NotFoundWithResponse, response)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

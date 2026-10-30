package arms

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

func (s *Service) DescribeArmsAlertContact(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "SearchAlertContact"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"ContactIds": convert.ConvertListToJsonString([]interface{}{id}),
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
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
	v, err := jsonpath.Get("$.PageBean.Contacts", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.PageBean.Contacts", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("ARMS", id), errs.NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["ContactId"]) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr("ARMS", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *Service) DescribeArmsAlertContactGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "SearchAlertContactGroup"
	request := map[string]interface{}{
		"RegionId":        s.client.RegionId,
		"ContactGroupIds": convert.ConvertListToJsonString([]interface{}{id}),
		"IsDetail":        "true",
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
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
	v, err := jsonpath.Get("$.ContactGroups", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ContactGroups", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("ARMS", id), errs.NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["ContactGroupId"]) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr("ARMS", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *Service) DescribeArmsAlertRobot(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeIMRobots"
	request := map[string]interface{}{
		"RobotIds": id,
		"Page":     1,
		"Size":     helper.PageSizeXLarge,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
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
	v, err := jsonpath.Get("$.PageBean.AlertIMRobots", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.PageBean.AlertIMRobots", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("ARMS", id), errs.NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["RobotId"]) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr("ARMS", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *Service) DescribeArmsDispatchRule(id string) (object map[string]interface{}, err error) {
	client := s.client

	request := map[string]interface{}{
		"Id":       id,
		"RegionId": s.client.RegionId,
	}

	var response map[string]interface{}
	action := "DescribeDispatchRule"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err := client.RpcPost("ARMS", "2019-08-08", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response = resp
		logs.AddDebug(action, resp, request)
		return nil
	})
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.DispatchRule", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.DispatchRule", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *Service) DescribeArmsPrometheusAlertRule(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListPrometheusAlertRules"
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"ClusterId": parts[0],
	}
	idExist := false
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
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
	v, err := jsonpath.Get("$.PrometheusAlertRules", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.PrometheusAlertRules", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("ARMS", id), errs.NotFoundWithResponse, response)
	}
	for _, v := range v.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["AlertId"]) == parts[1] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}
	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("ARMS", id), errs.NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *Service) ListArmsNotificationPolicies(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListNotificationPolicies"
	request := map[string]interface{}{
		"Page":     1,
		"Size":     helper.PageSizeXLarge,
		"IsDetail": true,
		"RegionId": s.client.RegionId,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
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
	v, err := jsonpath.Get("$.PageBean.NotificationPolicies", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.PageBean.NotificationPolicies", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("ARMS", id), errs.NotFoundWithResponse, response)
	}
	for _, v := range v.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["Id"]) == id {
			return v.(map[string]interface{}), nil
		}
	}
	return object, errs.WrapErrorf(errs.NotFoundErr("ARMS", id), errs.NotFoundWithResponse, response)
}

func (s *Service) ArmsDispatchRuleStateRefreshFunc(d *schema.ResourceData, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeArmsDispatchRule(d.Id())
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}
		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *Service) DescribeArmsPrometheus(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "GetPrometheusInstance"

	client := s.client

	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"ClusterId": id,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"404"}) {
			return object, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
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

func (s *Service) ListTagResources(id string, resourceType string) (object interface{}, err error) {
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
		parts, err := helper.ParseResourceId(id, 2)
		if err != nil {
			return object, errs.WrapError(err)
		}
		request["ResourceId.1"] = parts[resourceIdNum]
	}

	tags := make([]interface{}, 0)
	var response map[string]interface{}
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			v, err := jsonpath.Get("$.TagResources", response)
			if err != nil {
				return retry.NonRetryableError(errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TagResources", response))
			}

			if v != nil {
				tags = append(tags, v.([]interface{})...)
			}

			return nil
		})
		if err != nil {
			err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
			return
		}
		if response["NextToken"] == nil {
			break
		}
		request["NextToken"] = response["NextToken"]
	}

	return tags, nil
}

func (s *Service) SetResourceTags(d *schema.ResourceData, resourceType string) error {

	resourceIdNum := strings.Count(d.Id(), ":")

	if d.HasChange("tags") {
		added, removed := tags.ParsingTags(d)
		client := s.client

		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !tags.IgnoredTags(v, "") {
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
				parts, err := helper.ParseResourceId(d.Id(), 2)
				if err != nil {
					return errs.WrapError(err)
				}
				request["ResourceId.1"] = parts[resourceIdNum]
			}

			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			wait := helper.IncrementalWait(2*time.Second, 1*time.Second)
			err := retry.Retry(10*time.Minute, func() *retry.RetryError {
				response, err := client.RpcPost("ARMS", "2019-08-08", action, nil, request, false)
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
			action := "TagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
			}

			switch resourceIdNum {
			case 0:
				request["ResourceId.1"] = d.Id()
			case 1:
				parts, err := helper.ParseResourceId(d.Id(), 2)
				if err != nil {
					return errs.WrapError(err)
				}
				request["ResourceId.1"] = parts[resourceIdNum]
			}

			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}
			wait := helper.IncrementalWait(2*time.Second, 1*time.Second)
			err := retry.Retry(10*time.Minute, func() *retry.RetryError {
				response, err := client.RpcPost("ARMS", "2019-08-08", action, nil, request, false)
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

func (s *Service) DescribeArmsIntegrationExporter(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "GetPrometheusIntegration"

	client := s.client

	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return nil, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":        s.client.RegionId,
		"ClusterId":       parts[0],
		"IntegrationType": parts[1],
		"InstanceId":      parts[2],
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"404"}) {
			return object, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
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

func (s *Service) DescribeArmsRemoteWrite(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "GetPrometheusRemoteWrite"

	client := s.client

	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":        s.client.RegionId,
		"ClusterId":       parts[0],
		"RemoteWriteName": parts[1],
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"404"}) {
			return object, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
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

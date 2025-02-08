package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type NlbService struct {
	client *connectivity.AliyunClient
}

func (s *NlbService) SetResourceTags(d *schema.ResourceData, resourceType string) error {

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
				"ResourceId.1": d.Id(),
			}
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := client.RpcPost("Nlb", "2022-04-30", action, nil, request, false)
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
				"ResourceId.1": d.Id(),
			}
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := client.RpcPost("Nlb", "2022-04-30", action, nil, request, false)
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

func (s *NlbService) DescribeNlbServerGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListServerGroups"
	request := map[string]interface{}{
		"RegionId":       s.client.RegionId,
		"ServerGroupIds": []string{id},
	}
	idExist := false
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Nlb", "2022-04-30", action, nil, request, true)
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
	v, err := jsonpath.Get("$.ServerGroups", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ServerGroups", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("NLB", id)), NotFoundWithResponse, response)
	}
	for _, v := range v.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["ServerGroupId"]) == id {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("NLB", id)), NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *NlbService) ListTagResources(id string, resourceType string) (object interface{}, err error) {
	client := s.client
	action := "ListTagResources"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"ResourceType": resourceType,
		"ResourceId.1": id,
	}
	tags := make([]interface{}, 0)
	var response map[string]interface{}

	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err := client.RpcPost("Nlb", "2022-04-30", action, nil, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{Throttling}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			v, err := jsonpath.Get("$.TagResources", response)
			if formatInt(response["TotalCount"]) != 0 && err != nil {
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

func (s *NlbService) NlbServerGroupStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeNlbServerGroup(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["ServerGroupStatus"]) == failState {
				return object, fmt.Sprint(object["ServerGroupStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["ServerGroupStatus"])))
			}
		}
		return object, fmt.Sprint(object["ServerGroupStatus"]), nil
	}
}

func (s *NlbService) DescribeNlbSecurityPolicy(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListSecurityPolicy"
	request := map[string]interface{}{
		"RegionId":            s.client.RegionId,
		"SecurityPolicyIds.1": id,
	}
	idExist := false
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Nlb", "2022-04-30", action, nil, request, true)
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
	v, err := jsonpath.Get("$.SecurityPolicies", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.SecurityPolicies", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("NLB", id)), NotFoundWithResponse, response)
	}
	for _, v := range v.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["SecurityPolicyId"]) == id {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("NLB", id)), NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *NlbService) NlbSecurityPolicyStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeNlbSecurityPolicy(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["SecurityPolicyStatus"]) == failState {
				return object, fmt.Sprint(object["SecurityPolicyStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["SecurityPolicyStatus"])))
			}
		}
		return object, fmt.Sprint(object["SecurityPolicyStatus"]), nil
	}
}

func (s *NlbService) DescribeNlbLoadBalancer(id string) (object map[string]interface{}, err error) {
	client := s.client
	request := map[string]interface{}{
		"LoadBalancerId": id,
		"RegionId":       s.client.RegionId,
	}

	var response map[string]interface{}
	action := "GetLoadBalancerAttribute"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.RpcPost("Nlb", "2022-04-30", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *NlbService) NlbLoadBalancerStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeNlbLoadBalancer(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["LoadBalancerStatus"]) == failState {
				return object, fmt.Sprint(object["LoadBalancerStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["LoadBalancerStatus"])))
			}
		}
		return object, fmt.Sprint(object["LoadBalancerStatus"]), nil
	}
}

func (s *NlbService) DescribeNlbListener(id string) (object map[string]interface{}, err error) {
	client := s.client
	request := map[string]interface{}{
		"ListenerId": id,
		"RegionId":   s.client.RegionId,
	}

	var response map[string]interface{}
	action := "GetListenerAttribute"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.RpcPost("Nlb", "2022-04-30", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceNotFound.listener"}) {
			return object, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *NlbService) NlbListenerStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeNlbListener(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["ListenerStatus"]) == failState {
				return object, fmt.Sprint(object["ListenerStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["ListenerStatus"])))
			}
		}
		return object, fmt.Sprint(object["ListenerStatus"]), nil
	}
}

func (s *NlbService) DescribeNlbServerGroupServerAttachment(id string) (object map[string]interface{}, err error) {
	client := s.client
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"ServerGroupId": parts[0],
		"ServerIds":     parts[1],
	}

	var response map[string]interface{}
	action := "ListServerGroupServers"
	idExist := false
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.RpcPost("Nlb", "2022-04-30", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Servers", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Servers", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("NLB", id)), NotFoundWithResponse, response)
	}
	for _, v := range v.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["ServerId"]) == parts[1] && fmt.Sprint(v.(map[string]interface{})["ServerType"]) == parts[2] && fmt.Sprint(v.(map[string]interface{})["Port"]) == parts[3] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("NLB", id)), NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *NlbService) NlbServerGroupServerAttachmentStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeNlbServerGroupServerAttachment(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *NlbService) DescribeNlbLoadBalancerSecurityGroupAttachment(id string) (object map[string]interface{}, err error) {
	client := s.client
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":       s.client.RegionId,
		"ClientToken":    buildClientToken("GetLoadBalancerAttribute"),
		"LoadBalancerId": parts[0],
	}

	idExist := false
	var response map[string]interface{}
	action := "GetLoadBalancerAttribute"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Nlb", "2022-04-30", action, nil, request, true)
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

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}

	if securityGroupIdsList, ok := v.(map[string]interface{})["SecurityGroupIds"]; ok {
		for _, securityGroupIds := range securityGroupIdsList.([]interface{}) {
			if fmt.Sprint(securityGroupIds) == parts[1] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
	}

	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("NLB:SecurityGroupAttachment", id)), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *NlbService) GetJobStatus(id string) (object map[string]interface{}, err error) {
	client := s.client
	request := map[string]interface{}{
		"JobId": id,
	}
	var response map[string]interface{}
	action := "GetJobStatus"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.RpcPost("Nlb", "2022-04-30", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *NlbService) NlbLoadBalancerSecurityGroupAttachmentStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.GetJobStatus(id)
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

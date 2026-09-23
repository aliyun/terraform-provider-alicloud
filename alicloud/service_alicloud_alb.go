package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type AlbService struct {
	client *connectivity.AliyunClient
}

func (s *AlbService) ListAclEntries(id string) (objects []map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListAclEntries"
	request := map[string]interface{}{
		"AclId":      id,
		"MaxResults": PageSizeLarge,
	}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
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
			if IsExpectedErrors(err, []string{"ResourceNotFound.Acl"}) {
				return objects, WrapErrorf(NotFoundErr("ALB:Acl", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
			}
			return objects, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.AclEntries", response)
		if formatInt(response["TotalCount"]) != 0 && err != nil {
			return objects, WrapErrorf(err, FailedGetAttributeMsg, id, "$.AclEntries", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			objects = append(objects, item)
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return objects, nil
}

func (s *AlbService) ListTagResources(id string, resourceType string) (object interface{}, err error) {
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
			response, err := client.RpcPost("Alb", "2020-06-16", action, nil, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{Throttling}) {
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

func (s *AlbService) SetResourceTags(d *schema.ResourceData, resourceType string) error {

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
			action := "UnTagResources"
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
				response, err := client.RpcPost("Alb", "2020-06-16", action, nil, request, false)
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
				response, err := client.RpcPost("Alb", "2020-06-16", action, nil, request, false)
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

func (s *AlbService) DescribeAlbAcl(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListAcls"
	request := map[string]interface{}{
		"MaxResults": 100,
	}
	idExist := false
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
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
			if IsExpectedErrors(err, []string{"Forbidden.Acl"}) {
				return object, WrapErrorf(NotFoundErr("ALB:Acl", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
			}
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}
		totalCount, err := jsonpath.Get("$.TotalCount", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.TotalCount", response)
		}
		if fmt.Sprint(totalCount) == "0" {
			return object, WrapErrorf(NotFoundErr("ALB", id), NotFoundWithResponse, response)
		}
		v, err := jsonpath.Get("$.Acls", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Acls", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(NotFoundErr("ALB", id), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["AclId"]) == id {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	if !idExist {
		return object, WrapErrorf(NotFoundErr("ALB", id), NotFoundWithResponse, response)
	}
	return
}

func (s *AlbService) DescribeAlbSecurityPolicy(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListSecurityPolicies"
	request := map[string]interface{}{
		"MaxResults": 100,
	}
	idExist := false
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
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
		v, err := jsonpath.Get("$.SecurityPolicies", response)
		if formatInt(response["TotalCount"]) == 0 {
			return object, WrapErrorf(NotFoundErr("ALB", id), NotFoundWithResponse, response)
		}
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.SecurityPolicies", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(NotFoundErr("ALB", id), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["SecurityPolicyId"]) == id {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	if !idExist {
		return object, WrapErrorf(NotFoundErr("ALB", id), NotFoundWithResponse, response)
	}
	return
}

func (s *AlbService) ListSystemSecurityPolicies(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListSystemSecurityPolicies"
	request := map[string]interface{}{}
	idExist := false
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
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
	v, err := jsonpath.Get("$.SecurityPolicies", response)
	if formatInt(response["TotalCount"]) == 0 {
		return object, WrapErrorf(NotFoundErr("ALB", id), NotFoundWithResponse, response)
	}
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.SecurityPolicies", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("ALB", id), NotFoundWithResponse, response)
	}
	for _, v := range v.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["SecurityPolicyId"]) == id {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}
	if !idExist {
		return object, WrapErrorf(NotFoundErr("ALB", id), NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *AlbService) AlbSecurityPolicyStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAlbSecurityPolicy(id)
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

func (s *AlbService) ListServerGroupServers(id string) (object []interface{}, err error) {
	var response map[string]interface{}
	action := "ListServerGroupServers"

	client := s.client

	request := map[string]interface{}{
		"ServerGroupId": id,
	}

	idExist := false
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
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

		if formatInt(response["TotalCount"]) == 0 {
			return object, nil
		}

		resp, err := jsonpath.Get("$.Servers", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Servers", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, WrapErrorf(NotFoundErr("Alb:ServerGroup", id), NotFoundWithResponse, response)
		}

		for _, v := range resp.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["ServerGroupId"]) == id {
				idExist = true
				return resp.([]interface{}), nil
			}
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	if !idExist {
		return object, WrapErrorf(NotFoundErr("Alb:ServerGroup", id), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *AlbService) DescribeAlbServerGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListServerGroups"

	client := s.client

	request := map[string]interface{}{
		"MaxResults": PageSizeXLarge,
	}

	idExist := false
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
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

		resp, err := jsonpath.Get("$.ServerGroups", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ServerGroups", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, WrapErrorf(NotFoundErr("Alb:ServerGroup", id), NotFoundWithResponse, response)
		}

		for _, v := range resp.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["ServerGroupId"]) == id {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	if !idExist {
		return object, WrapErrorf(NotFoundErr("Alb:ServerGroup", id), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *AlbService) AlbServerGroupStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAlbServerGroup(id)
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

func (s *AlbService) DescribeAlbLoadBalancer(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetLoadBalancerAttribute"
	request := map[string]interface{}{
		"LoadBalancerId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"ResourceNotFound.LoadBalancer"}) {
			return object, WrapErrorf(NotFoundErr("ALB:LoadBalancer", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
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

func (s *AlbService) GetLoadBalancerAttribute(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetLoadBalancerAttribute"
	request := map[string]interface{}{
		"LoadBalancerId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"ResourceNotFound.LoadBalancer"}) {
			return object, WrapErrorf(NotFoundErr("ALB:LoadBalancer", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
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

func (s *AlbService) AlbLoadBalancerStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAlbLoadBalancer(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
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

func (s *AlbService) AlbLoadBalancerEditionRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAlbLoadBalancer(d.Id())
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["LoadBalancerEdition"]) == failState {
				return object, fmt.Sprint(object["LoadBalancerEdition"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["LoadBalancerEdition"])))
			}
		}
		return object, fmt.Sprint(object["LoadBalancerEdition"]), nil
	}
}

func (s *AlbService) AlbAclStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAlbAcl(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {

			if fmt.Sprint(object["AclStatus"]) == failState {
				return object, fmt.Sprint(object["AclStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["AclStatus"])))
			}
		}
		return object, fmt.Sprint(object["AclStatus"]), nil
	}
}

func (s *AlbService) AlbListenerStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAlbListener(id)
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

func (s *AlbService) DescribeAlbListener(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetListenerAttribute"
	request := map[string]interface{}{
		"ListenerId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"ResourceNotFound.Listener"}) {
			return object, WrapErrorf(NotFoundErr("ALB:Listener", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
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

func (s *AlbService) DescribeAlbRule(id, direction string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client

	action := "ListRules"
	request := map[string]interface{}{
		"MaxResults": PageSizeLarge,
	}

	if direction != "" {
		request["Direction"] = direction
	}

	idExist := false
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
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
		v, err := jsonpath.Get("$.Rules", response)
		if formatInt(response["TotalCount"]) != 0 && err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Rules", response)
		}
		if val, ok := v.([]interface{}); !ok || len(val) < 1 {
			return object, WrapErrorf(NotFoundErr("ALB", id), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["RuleId"]) == id {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	if !idExist {
		return object, WrapErrorf(NotFoundErr("ALB", id), NotFoundWithResponse, response)
	}
	return
}

func (s *AlbService) AlbRuleStateRefreshFunc(id, direction string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAlbRule(id, direction)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["RuleStatus"]) == failState {
				return object, fmt.Sprint(object["RuleStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["RuleStatus"])))
			}
		}
		return object, fmt.Sprint(object["RuleStatus"]), nil
	}
}

func (s *AlbService) DescribeAlbHealthCheckTemplate(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "GetHealthCheckTemplateAttribute"

	client := s.client

	request := map[string]interface{}{
		"HealthCheckTemplateId": id,
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"ResourceNotFound.HealthCheckTemplate"}) {
			return object, WrapErrorf(NotFoundErr("ALB:HealthCheckTemplate", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
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

func (s *AlbService) DescribeAlbListenerAdditionalCertificateAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return object, WrapError(err)
	}

	action := "ListListenerCertificates"
	request := map[string]interface{}{
		"ListenerId": parts[0],
		"MaxResults": PageSizeLarge,
	}
	idExist := false
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
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
		v, err := jsonpath.Get("$.Certificates", response)
		if formatInt(response["TotalCount"]) != 0 && err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Certificates", response)
		}
		if val, ok := v.([]interface{}); !ok || len(val) < 1 {
			return object, WrapErrorf(NotFoundErr("ALB", id), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["CertificateId"]) == parts[1] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	if !idExist {
		return object, WrapErrorf(NotFoundErr("ALB", id), NotFoundWithResponse, response)
	}
	return
}

func (s *AlbService) AlbListenerAdditionalCertificateAttachmentStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAlbListenerAdditionalCertificateAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
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

func (s *AlbService) DescribeAlbListenerAclAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return object, WrapError(err)
	}
	action := "GetListenerAttribute"
	request := map[string]interface{}{
		"ListenerId": parts[0],
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"ResourceNotFound.Listener"}) {
			return object, WrapErrorf(NotFoundErr("ALB", id), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	listenerObject := v.(map[string]interface{})

	aclConfig, ok := listenerObject["AclConfig"]
	if !ok {
		return object, WrapErrorf(NotFoundErr("alb_listener_acl_attachment", id), NotFoundWithResponse, response)

	}
	idExist := false
	aclConfigObject := aclConfig.(map[string]interface{})
	object = make(map[string]interface{})
	object["AclType"] = aclConfigObject["AclType"]
	if aclRelationsLis, ok := aclConfigObject["AclRelations"]; ok {
		for _, v := range aclRelationsLis.([]interface{}) {
			aclRelations := v.(map[string]interface{})
			if fmt.Sprint(aclRelations["AclId"]) == parts[1] {
				object["AclId"] = aclRelations["AclId"]
				object["Status"] = aclRelations["Status"]
				idExist = true
				break
			}
		}
	}

	if !idExist {
		return nil, WrapErrorf(NotFoundErr("alb_listener_acl_attachment", id), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *AlbService) DescribeAlbAclEntryAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return object, WrapError(err)
	}

	action := "ListAclEntries"
	request := map[string]interface{}{
		"AclId":      parts[0],
		"MaxResults": PageSizeLarge,
	}

	idExist := false
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
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
			if IsExpectedErrors(err, []string{"ResourceNotFound.Acl"}) {
				return object, WrapErrorf(NotFoundErr("ALB:AclEntryAttachment", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
			}
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.AclEntries", response)
		if formatInt(response["TotalCount"]) != 0 && err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.AclEntries", response)
		}
		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, WrapErrorf(NotFoundErr("ALB:AclEntryAttachment", id), NotFoundWithResponse, response)
		}
		for _, v := range resp.([]interface{}) {
			item := v.(map[string]interface{})
			if fmt.Sprint(item["Entry"]) == parts[1] {
				idExist = true
				return item, nil
			}
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	if !idExist {
		return object, WrapErrorf(NotFoundErr("ALB:AclEntryAttachment", id), NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *AlbService) AlbAclEntryAttachmentStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAlbAclEntryAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
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

func (s *AlbService) DescribeAlbAscript(id string) (object map[string]interface{}, err error) {
	client := s.client

	request := map[string]interface{}{
		"AScriptIds.1": id,
	}

	var response map[string]interface{}
	action := "ListAScripts"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"InternalError"}) {
			return object, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.AScripts", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.AScripts", response)
	}
	resp := v.([]interface{})
	if len(resp) < 1 {
		return object, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
	}

	return resp[0].(map[string]interface{}), nil
}

func (s *AlbService) AlbAscriptStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAlbAscript(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["AScriptStatus"]) == failState {
				return object, fmt.Sprint(object["AScriptStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["AScriptStatus"])))
			}
		}
		return object, fmt.Sprint(object["AScriptStatus"]), nil
	}
}

func (s *AlbService) DescribeAlbLoadBalancerCommonBandwidthPackageAttachment(id string) (object map[string]interface{}, err error) {
	client := s.client
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"LoadBalancerId": parts[0],
	}

	var response map[string]interface{}
	action := "GetLoadBalancerAttribute"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
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

func (s *AlbService) AlbLoadBalancerCommonBandwidthPackageAttachmentStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAlbLoadBalancerCommonBandwidthPackageAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		loadBalancerStatus82 := object["LoadBalancerStatus"]
		for _, failState := range failStates {
			if fmt.Sprint(loadBalancerStatus82) == failState {
				return object, fmt.Sprint(loadBalancerStatus82), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(loadBalancerStatus82)))
			}
		}
		return object, fmt.Sprint(loadBalancerStatus82), nil
	}
}

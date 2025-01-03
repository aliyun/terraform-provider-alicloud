package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	rpc "github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type VpcIpamServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeVpcIpamIpam <<< Encapsulated get interface for VpcIpam Ipam.

func (s *VpcIpamServiceV2) DescribeVpcIpamIpam(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "ListIpams"
	conn, err := client.NewVpcipamClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["IpamIds.1"] = id
	request["RegionId"] = client.RegionId

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-02-28"), StringPointer("AK"), query, request, &runtime)

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

	v, err := jsonpath.Get("$.Ipams[*]", response)
	if err != nil {
		return object, WrapErrorf(Error(GetNotFoundMessage("Ipam", id)), NotFoundMsg, response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(Error(GetNotFoundMessage("Ipam", id)), NotFoundMsg, response)
	}

	currentStatus := v.([]interface{})[0].(map[string]interface{})["IpamId"]
	if currentStatus == nil {
		return object, WrapErrorf(Error(GetNotFoundMessage("Ipam", id)), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *VpcIpamServiceV2) VpcIpamIpamStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcIpamIpam(id)
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

// DescribeVpcIpamIpam >>> Encapsulated.

// SetResourceTags <<< Encapsulated tag function for VpcIpam.
func (s *VpcIpamServiceV2) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var err error
		var action string
		var conn *rpc.Client
		client := s.client
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})

		added, removed := parsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action = "UntagResources"
			conn, err = client.NewVpcipamClient()
			if err != nil {
				return WrapError(err)
			}
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ResourceId.1"] = d.Id()
			request["RegionId"] = client.RegionId
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}

			request["ResourceType"] = resourceType
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-02-28"), StringPointer("AK"), query, request, &runtime)
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
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

		}

		if len(added) > 0 {
			action = "TagResources"
			conn, err = client.NewVpcipamClient()
			if err != nil {
				return WrapError(err)
			}
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
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-02-28"), StringPointer("AK"), query, request, &runtime)
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
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

		}
	}

	return nil
}

// SetResourceTags >>> tag function encapsulated.

// DescribeVpcIpamIpamScope <<< Encapsulated get interface for VpcIpam IpamScope.

func (s *VpcIpamServiceV2) DescribeVpcIpamIpamScope(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "ListIpamScopes"
	conn, err := client.NewVpcipamClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["IpamScopeIds.1"] = id
	request["RegionId"] = client.RegionId

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-02-28"), StringPointer("AK"), query, request, &runtime)

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

	v, err := jsonpath.Get("$.IpamScopes[*]", response)
	if err != nil {
		return object, WrapErrorf(Error(GetNotFoundMessage("IpamScope", id)), NotFoundMsg, response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(Error(GetNotFoundMessage("IpamScope", id)), NotFoundMsg, response)
	}

	currentStatus := v.([]interface{})[0].(map[string]interface{})["IpamScopeId"]
	if currentStatus == nil {
		return object, WrapErrorf(Error(GetNotFoundMessage("IpamScope", id)), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *VpcIpamServiceV2) VpcIpamIpamScopeStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcIpamIpamScope(id)
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

// DescribeVpcIpamIpamScope >>> Encapsulated.

// DescribeVpcIpamIpamPool <<< Encapsulated get interface for VpcIpam IpamPool.

func (s *VpcIpamServiceV2) DescribeVpcIpamIpamPool(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "ListIpamPools"
	conn, err := client.NewVpcipamClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["IpamPoolIds.1"] = id
	request["RegionId"] = client.RegionId

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-02-28"), StringPointer("AK"), query, request, &runtime)

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

	v, err := jsonpath.Get("$.IpamPools[*]", response)
	if err != nil {
		return object, WrapErrorf(Error(GetNotFoundMessage("IpamPool", id)), NotFoundMsg, response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(Error(GetNotFoundMessage("IpamPool", id)), NotFoundMsg, response)
	}

	currentStatus := v.([]interface{})[0].(map[string]interface{})["IpamPoolId"]
	if currentStatus == nil {
		return object, WrapErrorf(Error(GetNotFoundMessage("IpamPool", id)), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *VpcIpamServiceV2) VpcIpamIpamPoolStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcIpamIpamPool(id)
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

// DescribeVpcIpamIpamPool >>> Encapsulated.

// DescribeVpcIpamIpamPoolCidr <<< Encapsulated get interface for VpcIpam IpamPoolCidr.

func (s *VpcIpamServiceV2) DescribeVpcIpamIpamPoolCidr(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	action := "ListIpamPoolCidrs"
	conn, err := client.NewVpcipamClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Cidr"] = parts[1]
	request["IpamPoolId"] = parts[0]
	request["RegionId"] = client.RegionId

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-02-28"), StringPointer("AK"), query, request, &runtime)

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

	v, err := jsonpath.Get("$.IpamPoolCidrs[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.IpamPoolCidrs[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(Error(GetNotFoundMessage("IpamPoolCidr", id)), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *VpcIpamServiceV2) VpcIpamIpamPoolCidrStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcIpamIpamPoolCidr(id)
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

// DescribeVpcIpamIpamPoolCidr >>> Encapsulated.

// DescribeVpcIpamIpamPoolAllocation <<< Encapsulated get interface for VpcIpam IpamPoolAllocation.

func (s *VpcIpamServiceV2) DescribeVpcIpamIpamPoolAllocation(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "GetIpamPoolAllocation"
	conn, err := client.NewVpcipamClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["IpamPoolAllocationId"] = id
	query["RegionId"] = client.RegionId

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2023-02-28"), StringPointer("AK"), query, nil, &runtime)

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
		if IsExpectedErrors(err, []string{"ResourceNotFound.IpamPoolAllocation"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("IpamPoolAllocation", id)), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *VpcIpamServiceV2) VpcIpamIpamPoolAllocationStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcIpamIpamPoolAllocation(id)
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

// DescribeVpcIpamIpamPoolAllocation >>> Encapsulated.

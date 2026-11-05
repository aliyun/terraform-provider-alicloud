package vpcipam

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
)

type Service struct {
	client *connectivity.AliyunClient
}

// NewService constructs the VpcIpam (2023-02-28) Service from a connectivity client.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

// DescribeVpcIpamIpam <<< Encapsulated get interface for VpcIpam Ipam.

func (s *Service) DescribeVpcIpamIpam(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["IpamIds.1"] = id
	request["RegionId"] = client.RegionId
	action := "ListIpams"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("VpcIpam", "2023-02-28", action, query, request, true)

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

	v, err := jsonpath.Get("$.Ipams[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Ipam", id), errs.NotFoundMsg, response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("Ipam", id), errs.NotFoundMsg, response)
	}

	currentStatus := v.([]interface{})[0].(map[string]interface{})["IpamId"]
	if currentStatus == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Ipam", id), errs.NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *Service) VpcIpamIpamStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcIpamIpam(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
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

// DescribeVpcIpamIpam >>> Encapsulated.

// SetResourceTags <<< Encapsulated tag function for VpcIpam.
func (s *Service) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var action string
		var err error
		client := s.client
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})

		added, removed := tags.ParsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !tags.IgnoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action = "UntagResources"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ResourceId.1"] = d.Id()
			request["RegionId"] = client.RegionId
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}

			request["ResourceType"] = resourceType
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("VpcIpam", "2023-02-28", action, query, request, true)
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
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}

		}

		if len(added) > 0 {
			action = "TagResources"
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
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("VpcIpam", "2023-02-28", action, query, request, true)
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
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}

		}
	}

	return nil
}

// SetResourceTags >>> tag function encapsulated.

// DescribeVpcIpamIpamScope <<< Encapsulated get interface for VpcIpam IpamScope.

func (s *Service) DescribeVpcIpamIpamScope(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["IpamScopeIds.1"] = id
	request["RegionId"] = client.RegionId
	action := "ListIpamScopes"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("VpcIpam", "2023-02-28", action, query, request, true)

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

	v, err := jsonpath.Get("$.IpamScopes[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("IpamScope", id), errs.NotFoundMsg, response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("IpamScope", id), errs.NotFoundMsg, response)
	}

	currentStatus := v.([]interface{})[0].(map[string]interface{})["IpamScopeId"]
	if currentStatus == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("IpamScope", id), errs.NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *Service) VpcIpamIpamScopeStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcIpamIpamScope(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
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

// DescribeVpcIpamIpamScope >>> Encapsulated.

// DescribeVpcIpamIpamPool <<< Encapsulated get interface for VpcIpam IpamPool.

func (s *Service) DescribeVpcIpamIpamPool(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["IpamPoolIds.1"] = id
	request["RegionId"] = client.RegionId
	action := "ListIpamPools"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("VpcIpam", "2023-02-28", action, query, request, true)

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

	v, err := jsonpath.Get("$.IpamPools[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("IpamPool", id), errs.NotFoundMsg, response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("IpamPool", id), errs.NotFoundMsg, response)
	}

	currentStatus := v.([]interface{})[0].(map[string]interface{})["IpamPoolId"]
	if currentStatus == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("IpamPool", id), errs.NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *Service) VpcIpamIpamPoolStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.VpcIpamIpamPoolStateRefreshFuncWithApi(id, field, failStates, s.DescribeVpcIpamIpamPool)
}

func (s *Service) VpcIpamIpamPoolStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
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

// DescribeVpcIpamIpamPool >>> Encapsulated.

// DescribeVpcIpamIpamPoolCidr <<< Encapsulated get interface for VpcIpam IpamPoolCidr.

func (s *Service) DescribeVpcIpamIpamPoolCidr(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts, err := splitIpamPoolCidrId(id)
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
		return nil, err
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Cidr"] = parts[1]
	request["IpamPoolId"] = parts[0]
	request["RegionId"] = client.RegionId
	action := "ListIpamPoolCidrs"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("VpcIpam", "2023-02-28", action, query, request, true)

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

	v, err := jsonpath.Get("$.IpamPoolCidrs[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.IpamPoolCidrs[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("IpamPoolCidr", id), errs.NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *Service) VpcIpamIpamPoolCidrStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.VpcIpamIpamPoolCidrStateRefreshFuncWithApi(id, field, failStates, s.DescribeVpcIpamIpamPoolCidr)
}

func (s *Service) VpcIpamIpamPoolCidrStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
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

// DescribeVpcIpamIpamPoolCidr >>> Encapsulated.

// DescribeVpcIpamIpamPoolAllocation <<< Encapsulated get interface for VpcIpam IpamPoolAllocation.

func (s *Service) DescribeVpcIpamIpamPoolAllocation(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "GetIpamPoolAllocation"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["IpamPoolAllocationId"] = id
	query["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("VpcIpam", "2023-02-28", action, query, nil)

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
		if errs.IsExpectedErrors(err, []string{"ResourceNotFound.IpamPoolAllocation"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("IpamPoolAllocation", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *Service) VpcIpamIpamPoolAllocationStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcIpamIpamPoolAllocation(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
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

// DescribeVpcIpamIpamPoolAllocation >>> Encapsulated.

// DescribeVpcIpamService <<< Encapsulated get interface for VpcIpam Service.

func (s *Service) DescribeVpcIpamService(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	action := "GetVpcIpamServiceStatus"
	request["ClientToken"] = helper.BuildClientToken(action)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("VpcIpam", "2023-02-28", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

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

	return response, nil
}

func (s *Service) VpcIpamServiceStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcIpamService(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
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

// DescribeVpcIpamService >>> Encapsulated.

// DescribeVpcIpamIpamResourceDiscovery <<< Encapsulated get interface for VpcIpam IpamResourceDiscovery.

func (s *Service) DescribeVpcIpamIpamResourceDiscovery(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["IpamResourceDiscoveryIds.1"] = id
	request["RegionId"] = client.RegionId
	action := "ListIpamResourceDiscoveries"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("VpcIpam", "2023-02-28", action, query, request, true)

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

	v, err := jsonpath.Get("$.IpamResourceDiscoveries[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("IpamResourceDiscovery", id), errs.NotFoundMsg, response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("IpamResourceDiscovery", id), errs.NotFoundMsg, response)
	}

	currentStatus := v.([]interface{})[0].(map[string]interface{})["IpamResourceDiscoveryId"]
	if currentStatus == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("IpamResourceDiscovery", id), errs.NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *Service) VpcIpamIpamResourceDiscoveryStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcIpamIpamResourceDiscovery(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
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

// DescribeVpcIpamIpamResourceDiscovery >>> Encapsulated.

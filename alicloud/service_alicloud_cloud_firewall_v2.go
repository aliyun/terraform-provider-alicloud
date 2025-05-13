package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type CloudFirewallServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeCloudFirewallNatFirewallControlPolicy <<< Encapsulated get interface for CloudFirewall NatFirewallControlPolicy.

func (s *CloudFirewallServiceV2) DescribeCloudFirewallNatFirewallControlPolicy(id string) (object map[string]interface{}, err error) {
	client := s.client
	var endpoint string
	var response map[string]interface{}
	var request map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	action := "DescribeNatFirewallControlPolicy"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["AclUuid"] = parts[0]
	query["Direction"] = parts[2]
	query["NatGatewayId"] = parts[1]

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, query, request, true, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			} else if IsExpectedErrors(err, []string{"not buy user"}) {
				endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
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

	v, err := jsonpath.Get("$.Policys[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Policys[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("NatFirewallControlPolicy", id), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *CloudFirewallServiceV2) CloudFirewallNatFirewallControlPolicyStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCloudFirewallNatFirewallControlPolicy(id)
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

// DescribeCloudFirewallNatFirewallControlPolicy >>> Encapsulated.

// DescribeCloudFirewallNatFirewall <<< Encapsulated get interface for CloudFirewall NatFirewall.

func (s *CloudFirewallServiceV2) DescribeCloudFirewallNatFirewall(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ProxyId"] = id

	action := "DescribeNatFirewallList"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Cloudfw", "2017-12-07", action, query, request, true)

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

	v, err := jsonpath.Get("$.NatFirewallList[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.NatFirewallList[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("NatFirewall", id), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *CloudFirewallServiceV2) CloudFirewallNatFirewallStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCloudFirewallNatFirewall(id)
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

// DescribeCloudFirewallNatFirewall >>> Encapsulated.

// DescribeCloudFirewallVpcCenTrFirewall <<< Encapsulated get interface for CloudFirewall VpcCenTrFirewall.

func (s *CloudFirewallServiceV2) DescribeCloudFirewallVpcCenTrFirewall(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["FirewallId"] = id

	action := "DescribeTrFirewallsV2Detail"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Cloudfw", "2017-12-07", action, query, request, true)

		if err != nil {
			if IsExpectedErrors(err, []string{"ErrorTrResourceNotReady"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorTrFirewallNotExist"}) {
			return object, WrapErrorf(NotFoundErr("VpcCenTrFirewall", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *CloudFirewallServiceV2) CloudFirewallVpcCenTrFirewallStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCloudFirewallVpcCenTrFirewall(id)
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
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

// DescribeCloudFirewallVpcCenTrFirewall >>> Encapsulated.

// DescribeCloudFirewallControlPolicy <<< Encapsulated get interface for CloudFirewall ControlPolicy.

func (s *CloudFirewallServiceV2) DescribeCloudFirewallControlPolicy(id string) (object map[string]interface{}, err error) {
	client := s.client
	var endpoint string
	var response map[string]interface{}
	action := "DescribeControlPolicy"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"AclUuid":     parts[0],
		"CurrentPage": 1,
		"PageSize":    PageSizeLarge,
	}

	idExist := false
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, nil, request, true, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				} else if IsExpectedErrors(err, []string{"not buy user"}) {
					endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
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

		resp, err := jsonpath.Get("$.Policys", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Policys", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, WrapErrorf(NotFoundErr("CloudFirewall:ControlPolicy", id), NotFoundWithResponse, response)
		}

		for _, v := range resp.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["AclUuid"]) == parts[0] && fmt.Sprint(v.(map[string]interface{})["Direction"]) == parts[1] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}

		if len(resp.([]interface{})) < request["PageSize"].(int) {
			break
		}

		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}

	if !idExist {
		return object, WrapErrorf(NotFoundErr("CloudFirewall:ControlPolicy", id), NotFoundWithResponse, response)
	}

	return object, nil
}

// DescribeCloudFirewallControlPolicy >>> Encapsulated.

// Async Api <<< Encapsulated for CloudFirewall.
// Async Api >>> Encapsulated.

// DescribeCloudFirewallIPSConfig <<< Encapsulated get interface for CloudFirewall IPSConfig.

func (s *CloudFirewallServiceV2) DescribeCloudFirewallIPSConfig(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	action := "DescribeDefaultIPSConfig"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Cloudfw", "2017-12-07", action, query, request, true)

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

	return response, nil
}

func (s *CloudFirewallServiceV2) CloudFirewallIPSConfigStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCloudFirewallIPSConfig(id)
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

// DescribeCloudFirewallIPSConfig >>> Encapsulated.

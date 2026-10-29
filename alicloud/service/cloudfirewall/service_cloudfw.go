package cloudfirewall

import (
	json "encoding/json"
	"fmt"
	"time"

	jsonpath "github.com/PaesslerAG/jsonpath"
	connectivity "github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	errs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	helper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	retry "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (s *Service) DescribeCloudFirewallAddressBook(id string) (object map[string]interface{}, err error) {
	client := s.client
	var endpoint string
	var response map[string]interface{}
	action := "DescribeAddressBook"
	request := map[string]interface{}{
		"GroupType":   "all",
		"CurrentPage": 1,
		"PageSize":    helper.PageSizeLarge,
	}

	idExist := false
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, nil, request, true, endpoint)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				} else if errs.IsExpectedErrors(err, []string{"not buy user"}) {
					endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
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

		resp, err := jsonpath.Get("$.Acls", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Acls", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, errs.WrapErrorf(errs.NotFoundErr("CloudFirewall:AddressBook", id), errs.NotFoundWithResponse, response)
		}

		for _, v := range resp.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["GroupUuid"]) == id {
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
		return object, errs.WrapErrorf(errs.NotFoundErr("CloudFirewall:AddressBook", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *Service) DescribeCloudFirewallInstanceMember(id string) (object map[string]interface{}, err error) {
	client := s.client
	var endpoint string
	var response map[string]interface{}
	action := "DescribeInstanceMembers"
	request := map[string]interface{}{
		"MemberUid": id,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, nil, request, true, endpoint)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			} else if errs.IsExpectedErrors(err, []string{"not buy user"}) {
				endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidResource.NotFound"}) {
			return object, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Members", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Members", response)
	}
	members := v.([]interface{})
	if len(members) < 1 {
		return object, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
	}
	return members[0].(map[string]interface{}), nil
}

func (s *Service) CloudFirewallInstanceMemberStateRefreshFunc(d *schema.ResourceData, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCloudFirewallInstanceMember(d.Id())
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["MemberStatus"]) == failState {
				return object, fmt.Sprint(object["MemberStatus"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["MemberStatus"])))
			}
		}
		return object, fmt.Sprint(object["MemberStatus"]), nil
	}
}

func (s *Service) DescribeCloudFirewallVpcFirewallCen(id string) (object map[string]interface{}, err error) {
	client := s.client
	var endpoint string
	var response map[string]interface{}
	action := "DescribeVpcFirewallCenDetail"
	request := map[string]interface{}{
		"VpcFirewallId": id,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, nil, request, true, endpoint)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			} else if errs.IsExpectedErrors(err, []string{"not buy user"}) {
				endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}

		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ErrorVpcFirewallExist"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("VpcFirewallCen", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	status, err := jsonpath.Get("$.FirewallSwitchStatus", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	if status.(string) == "notconfigured" {
		return object, errs.WrapErrorf(errs.NotFoundErr("VpcFirewallCen", id), errs.NotFoundWithResponse, response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) CreateVpcFirewallTask() (err error) {
	client := s.client
	var endpoint string
	var response map[string]interface{}
	action := "CreateVpcFirewallTask"
	request := map[string]interface{}{
		"TaskAction": "sync",
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, nil, request, false, endpoint)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			} else if errs.IsExpectedErrors(err, []string{"not buy user"}) {
				endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}

		return nil
	})
	logs.AddDebug(action, response, request)
	return err
}

func (s *Service) DescribeVpcFirewallCenList(id string) (object map[string]interface{}, err error) {
	client := s.client
	var endpoint string
	var response map[string]interface{}
	action := "DescribeVpcFirewallCenList"
	request := map[string]interface{}{
		"VpcFirewallId": id,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, nil, request, true, endpoint)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			} else if errs.IsExpectedErrors(err, []string{"not buy user"}) {
				endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
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

	v, err := jsonpath.Get("$.VpcFirewalls[0]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	totalCount, _ := jsonpath.Get("$.TotalCount", response)
	total, _ := totalCount.(json.Number).Int64()
	if err != nil && total == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("VpcFirewallCen", id), errs.NotFoundWithResponse, response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) CloudFirewallVpcFirewallCenStateRefreshFunc(d *schema.ResourceData, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCloudFirewallVpcFirewallCen(d.Id())
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["FirewallSwitchStatus"]) == failState {
				return object, fmt.Sprint(object["FirewallSwitchStatus"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["FirewallSwitchStatus"])))
			}
		}
		return object, fmt.Sprint(object["FirewallSwitchStatus"]), nil
	}
}

func (s *Service) DescribeCloudFirewallVpcFirewall(id string) (object map[string]interface{}, err error) {
	client := s.client
	var endpoint string
	var response map[string]interface{}
	action := "DescribeVpcFirewallDetail"
	request := map[string]interface{}{
		"VpcFirewallId": id,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, nil, request, true, endpoint)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			} else if errs.IsExpectedErrors(err, []string{"not buy user"}) {
				endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}

		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ErrorFirewallNotFound"}) {
			return object, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	switchStatus, err := jsonpath.Get("$.FirewallSwitchStatus", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.FirewallSwitchStatus", response)
	}

	if fmt.Sprint(switchStatus) == "notconfigured" {
		return object, errs.WrapErrorf(errs.NotFoundErr("CloudFirewall:VpcFirewall", id), errs.NotFoundWithResponse, response)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *Service) DescribeVpcFirewallList(id string) (object map[string]interface{}, err error) {
	client := s.client
	var endpoint string
	var response map[string]interface{}
	action := "DescribeVpcFirewallList"
	request := map[string]interface{}{
		"VpcFirewallId": id,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, nil, request, true, endpoint)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			} else if errs.IsExpectedErrors(err, []string{"not buy user"}) {
				endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}

		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ErrorFirewallNotFound"}) {
			return object, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	totalCount, _ := jsonpath.Get("$.TotalCount", response)
	total, _ := totalCount.(json.Number).Int64()
	if err != nil && total == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("CloudFirewall:VpcFirewall", id), errs.NotFoundWithResponse, response)
	}

	v, err := jsonpath.Get("$.VpcFirewalls[0]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) CloudFirewallVpcFirewallStateRefreshFunc(d *schema.ResourceData, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCloudFirewallVpcFirewall(d.Id())
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["FirewallSwitchStatus"]) == failState {
				return object, fmt.Sprint(object["FirewallSwitchStatus"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["FirewallSwitchStatus"])))
			}
		}
		return object, fmt.Sprint(object["FirewallSwitchStatus"]), nil
	}
}

func (s *Service) DescribeCloudFirewallInstanceUserBuyVersion(id string) (object map[string]interface{}, err error) {
	client := s.client
	var endpoint string
	var response map[string]interface{}
	action := "DescribeUserBuyVersion"
	request := map[string]interface{}{
		"InstanceId": id,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, nil, request, true, endpoint)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			} else if errs.IsExpectedErrors(err, []string{"not valid instanceId"}) {
				endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
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

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *Service) CloudFirewallInstanceStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCloudFirewallInstanceUserBuyVersion(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["InstanceStatus"].(string) == failState {
				return object, object["InstanceStatus"].(string), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object["InstanceStatus"].(string)))
			}
		}

		return object, object["InstanceStatus"].(string), nil
	}
}

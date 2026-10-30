package vpngateway

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
)

type VpcService struct {
	client *connectivity.AliyunClient
}

func NewVpcService(client *connectivity.AliyunClient) *VpcService {
	return &VpcService{client: client}
}

func (s *VpcService) DescribeVpnRouteEntry(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeVpnRouteEntries"
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return object, errs.WrapError(err)
	}
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"VpnGatewayId": parts[0],
		"PageNumber":   1,
		"PageSize":     helper.PageSizeMedium,
	}
	idExist := false
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
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
		v, err := jsonpath.Get("$.VpnRouteEntries.VpnRouteEntry", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.VpnRouteEntries.VpnRouteEntry", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, errs.WrapErrorf(errs.NotFoundErr("VPC", id), errs.NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["NextHop"]) == parts[1] && fmt.Sprint(v.(map[string]interface{})["RouteDest"]) == parts[2] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("VPC", id), errs.NotFoundWithResponse, response)
	}
	return
}

func (s *VpcService) VpnRouteEntryStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpnRouteEntry(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["State"].(string) == failState {
				return object, object["State"].(string), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object["State"].(string)))
			}
		}
		return object, object["State"].(string), nil
	}
}

func (s *VpcService) DescribeVpnIpsecServer(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListIpsecServers"
	request := map[string]interface{}{
		"RegionId":      s.client.RegionId,
		"IpsecServerId": []string{id},
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
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
	v, err := jsonpath.Get("$.IpsecServers", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.IpsecServers", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("VPC", id), errs.NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["IpsecServerId"]) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr("VPC", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *VpcService) DescribeVpnPbrRouteEntry(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeVpnPbrRouteEntries"
	parts, err := helper.ParseResourceId(id, 4)
	if err != nil {
		return object, errs.WrapError(err)
	}
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"VpnGatewayId": parts[0],
		"PageNumber":   1,
		"PageSize":     helper.PageSizeMedium,
	}
	idExist := false
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
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
		v, err := jsonpath.Get("$.VpnPbrRouteEntries.VpnPbrRouteEntry", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.VpnPbrRouteEntries.VpnPbrRouteEntry", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, errs.WrapErrorf(errs.NotFoundErr("VPC", id), errs.NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["NextHop"]) == parts[1] && fmt.Sprint(v.(map[string]interface{})["RouteSource"]) == parts[2] && fmt.Sprint(v.(map[string]interface{})["RouteDest"]) == parts[3] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("VPC", id), errs.NotFoundWithResponse, response)
	}
	return
}

func (s *VpcService) VpnPbrRouteEntryStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpnPbrRouteEntry(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["State"].(string) == failState {
				return object, object["State"].(string), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object["State"].(string)))
			}
		}
		return object, object["State"].(string), nil
	}
}

func (s *VpcService) DescribeVpnGatewayVpnAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeVpnConnection"
	request := map[string]interface{}{
		"RegionId":        s.client.RegionId,
		"VpnConnectionId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"InvalidVpnConnectionInstanceId.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("VpnConnection", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *VpcService) VpnGatewayVpnAttachmentStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpnGatewayVpnAttachment(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["State"].(string) == failState {
				return object, object["State"].(string), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object["State"].(string)))
			}
		}
		return object, object["State"].(string), nil
	}
}

func (s *VpcService) DescribeVpnGatewayVcoRoute(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeVcoRouteEntries"
	parts, err := helper.ParseResourceId(id, 4)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":        s.client.RegionId,
		"VpnConnectionId": parts[0],
		"PageNumber":      1,
		"PageSize":        helper.PageSizeMedium,
	}
	idExist := false
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
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
			if errs.IsExpectedErrors(err, []string{"UnknownError"}) {
				return object, errs.WrapErrorf(errs.NotFoundErr("VPC:VPNGateway", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
			}
			return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		}
		v, err := jsonpath.Get("$.VcoRouteEntries", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.VcoRouteEntries", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, errs.WrapErrorf(errs.NotFoundErr("VPNGateway", id), errs.NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			item := v.(map[string]interface{})
			if fmt.Sprint(item["RouteDest"]) == parts[1] && fmt.Sprint(item["NextHop"]) == parts[2] && fmt.Sprint(item["Weight"]) == parts[3] {
				idExist = true
				return item, nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("VPNGateway", id), errs.NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *VpcService) VpnGatewayVcoRouteStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpnGatewayVcoRoute(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		state := fmt.Sprint(object["State"])
		for _, failState := range failStates {
			if state == failState {
				return object, state, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, state))
			}
		}
		return object, state, nil
	}
}

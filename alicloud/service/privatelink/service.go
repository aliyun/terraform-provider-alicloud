package privatelink

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// Service wraps the Privatelink API calls used by the resources and data
// sources of this package.
type Service struct {
	client *connectivity.AliyunClient
}

// NewService is the exported constructor: the client field stays unexported,
// so callers outside the package (including the acceptance tests) build
// values with it.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

func (s *Service) ListVpcEndpointServiceResources(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListVpcEndpointServiceResources"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"ServiceId": id,
	}
	response, err = s.client.RpcPost("Privatelink", "2020-04-15", action, nil, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EndpointServiceNotFound"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("PrivatelinkVpcEndpointService", id), errs.NotFoundMsg, errs.ProviderERROR)
			return object, err
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return object, err
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) DescribePrivatelinkVpcEndpointService(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "GetVpcEndpointServiceAttribute"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"ServiceId": id,
	}
	response, err = s.client.RpcPost("Privatelink", "2020-04-15", action, nil, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EndpointServiceNotFound"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("PrivatelinkVpcEndpointService", id), errs.NotFoundMsg, errs.ProviderERROR)
			return object, err
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return object, err
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) PrivatelinkVpcEndpointServiceStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePrivatelinkVpcEndpointService(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["ServiceStatus"].(string) == failState {
				return object, object["ServiceStatus"].(string), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object["ServiceStatus"].(string)))
			}
		}
		return object, object["ServiceStatus"].(string), nil
	}
}

func (s *Service) DescribePrivatelinkVpcEndpointConnection(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListVpcEndpointConnections"
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"EndpointId": parts[1],
		"ServiceId":  parts[0],
	}
	response, err = s.client.RpcPost("Privatelink", "2020-04-15", action, nil, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EndpointServiceNotFound"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("PrivatelinkVpcEndpointConnection", id), errs.NotFoundMsg, errs.ProviderERROR)
			return object, err
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return object, err
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.Connections", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Connections", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("PrivateLink", id), errs.NotFoundWithResponse, response)
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *Service) PrivatelinkVpcEndpointConnectionStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePrivatelinkVpcEndpointConnection(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["ConnectionStatus"].(string) == failState {
				return object, object["ConnectionStatus"].(string), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object["ConnectionStatus"].(string)))
			}
		}
		return object, object["ConnectionStatus"].(string), nil
	}
}

func (s *Service) ListVpcEndpointSecurityGroups(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListVpcEndpointSecurityGroups"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"EndpointId": id,
	}
	response, err = s.client.RpcPost("Privatelink", "2020-04-15", action, nil, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EndpointNotFound"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("PrivatelinkVpcEndpoint", id), errs.NotFoundMsg, errs.ProviderERROR)
			return object, err
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return object, err
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) ListVpcEndpointZones(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListVpcEndpointZones"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"EndpointId": id,
	}
	response, err = s.client.RpcPost("Privatelink", "2020-04-15", action, nil, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EndpointNotFound"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("PrivatelinkVpcEndpoint", id), errs.NotFoundMsg, errs.ProviderERROR)
			return object, err
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return object, err
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) DescribePrivatelinkVpcEndpoint(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "GetVpcEndpointAttribute"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"EndpointId": id,
	}
	request["ClientToken"] = helper.BuildClientToken("GetVpcEndpointAttribute")
	response, err = s.client.RpcPost("Privatelink", "2020-04-15", action, nil, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EndpointNotFound"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("PrivatelinkVpcEndpoint", id), errs.NotFoundMsg, errs.ProviderERROR)
			return object, err
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return object, err
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) PrivatelinkVpcEndpointStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePrivatelinkVpcEndpoint(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["EndpointStatus"].(string) == failState {
				return object, object["EndpointStatus"].(string), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object["EndpointStatus"].(string)))
			}
		}
		return object, object["EndpointStatus"].(string), nil
	}
}

func (s *Service) DescribePrivatelinkVpcEndpointServiceResource(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListVpcEndpointServiceResources"
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"ServiceId": parts[0],
	}
	for {
		response, err = s.client.RpcPost("Privatelink", "2020-04-15", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"EndpointServiceNotFound"}) {
				err = errs.WrapErrorf(errs.NotFoundErr("PrivatelinkVpcEndpointServiceResource", id), errs.NotFoundMsg, errs.ProviderERROR)
				return object, err
			}
			err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
			return object, err
		}
		logs.AddDebug(action, response, request)
		v, err := jsonpath.Get("$.Resources", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Resources", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, errs.WrapErrorf(errs.NotFoundErr("PrivateLink", id), errs.NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if v.(map[string]interface{})["ResourceId"].(string) == parts[1] {
				return v.(map[string]interface{}), nil
			}
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
		return object, nil
	}
	return
}

func (s *Service) DescribePrivatelinkVpcEndpointServiceUser(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListVpcEndpointServiceUsers"
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"ServiceId": parts[0],
		"UserId":    parts[1],
	}
	response, err = s.client.RpcPost("Privatelink", "2020-04-15", action, nil, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EndpointServiceNotFound"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("PrivatelinkVpcEndpointServiceUser", id), errs.NotFoundMsg, errs.ProviderERROR)
			return object, err
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return object, err
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.Users", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Users", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("PrivateLink", id), errs.NotFoundWithResponse, response)
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *Service) DescribePrivatelinkVpcEndpointZone(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListVpcEndpointZones"
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"EndpointId": parts[0],
	}
	for {
		response, err = s.client.RpcPost("Privatelink", "2020-04-15", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"EndpointNotFound"}) {
				err = errs.WrapErrorf(errs.NotFoundErr("PrivatelinkVpcEndpointZone", id), errs.NotFoundMsg, errs.ProviderERROR)
				return object, err
			}
			err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
			return object, err
		}
		logs.AddDebug(action, response, request)
		v, err := jsonpath.Get("$.Zones", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Zones", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, errs.WrapErrorf(errs.NotFoundErr("PrivateLink", id), errs.NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if v.(map[string]interface{})["ZoneId"].(string) == parts[1] {
				return v.(map[string]interface{}), nil
			}
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
		return object, nil
	}
	return
}

func (s *Service) PrivatelinkVpcEndpointZoneStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePrivatelinkVpcEndpointZone(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["ZoneStatus"].(string) == failState {
				return object, object["ZoneStatus"].(string), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object["ZoneStatus"].(string)))
			}
		}
		return object, object["ZoneStatus"].(string), nil
	}
}

// DescribeVswitch looks up a VSwitch by id via the VPC API so the endpoint
// zone resource can resolve the zone of a given vSwitch. It mirrors the
// stock VpcService.DescribeVswitch call the root package used to make here.
func (s *Service) DescribeVswitch(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeVSwitchAttributes"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"VSwitchId": id,
	}
	response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidVSwitchId.NotFound", "InvalidVswitchID.NotFound"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("Vswitch", id), errs.NotFoundWithError, err)
			return object, err
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return object, err
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	if fmt.Sprint(object["VSwitchId"]) != id {
		return object, errs.WrapErrorf(errs.NotFoundErr("vswitch", id), errs.NotFoundWithResponse, response)
	}
	return object, nil
}

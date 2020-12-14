package alicloud

import (
	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type PrivatelinkService struct {
	client *connectivity.AliyunClient
}

func (s *PrivatelinkService) ListVpcEndpointServiceResources(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewPrivatelinkClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListVpcEndpointServiceResources"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"ServiceId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"EndpointServiceNotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("PrivatelinkVpcEndpointService", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *PrivatelinkService) DescribePrivatelinkVpcEndpointService(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewPrivatelinkClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetVpcEndpointServiceAttribute"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"ServiceId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"EndpointServiceNotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("PrivatelinkVpcEndpointService", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *PrivatelinkService) PrivatelinkVpcEndpointServiceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePrivatelinkVpcEndpointService(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["ServiceStatus"].(string) == failState {
				return object, object["ServiceStatus"].(string), WrapError(Error(FailedToReachTargetStatus, object["ServiceStatus"].(string)))
			}
		}
		return object, object["ServiceStatus"].(string), nil
	}
}

func (s *PrivatelinkService) ListVpcEndpointSecurityGroups(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewPrivatelinkClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListVpcEndpointSecurityGroups"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"EndpointId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"EndpointNotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("PrivatelinkVpcEndpoint", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *PrivatelinkService) ListVpcEndpointZones(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewPrivatelinkClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListVpcEndpointZones"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"EndpointId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"EndpointNotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("PrivatelinkVpcEndpoint", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *PrivatelinkService) DescribePrivatelinkVpcEndpoint(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewPrivatelinkClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetVpcEndpointAttribute"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"EndpointId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	request["ClientToken"] = buildClientToken("GetVpcEndpointAttribute")
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"EndpointNotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("PrivatelinkVpcEndpoint", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *PrivatelinkService) PrivatelinkVpcEndpointStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePrivatelinkVpcEndpoint(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["EndpointStatus"].(string) == failState {
				return object, object["EndpointStatus"].(string), WrapError(Error(FailedToReachTargetStatus, object["EndpointStatus"].(string)))
			}
		}
		return object, object["EndpointStatus"].(string), nil
	}
}

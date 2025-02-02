package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type DgService struct {
	client *connectivity.AliyunClient
}

func (s *DgService) DescribeDatabaseGatewayGateway(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetUserGateways"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	idExist := false
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("dg", "2019-03-27", action, nil, request, true)
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
		m, err := jsonpath.Get("$.Data", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
		}
		v, err := convertJsonStringToList(m.(string))
		if err != nil {
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}
		if len(v) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("DatabaseGateway", id)), NotFoundWithResponse, response)
		}
		for _, v := range v {
			if fmt.Sprint(v.(map[string]interface{})["gatewayId"]) == id {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("DatabaseGateway", id)), NotFoundWithResponse, response)
	}
	return
}

func (s *DgService) GetUserGatewayInstances(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetUserGatewayInstances"
	request := map[string]interface{}{
		"GatewayId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("dg", "2019-03-27", action, nil, request, true)
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
	m, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}
	v, err := convertJsonStringToList(m.(string))
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	if len(v) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("DatabaseGateway", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v[0].(map[string]interface{})["gatewayId"]) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("DatabaseGateway", id)), NotFoundWithResponse, response)
		}
	}
	return object, nil
}

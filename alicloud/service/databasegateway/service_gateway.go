package databasegateway

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// Service wraps the Database Gateway API calls used by the resources and
// data sources of this package.
type Service struct {
	client *connectivity.AliyunClient
}

// NewService is the exported constructor: the client field stays unexported,
// so callers outside the package (including the flat layout's acceptance
// tests) build values with it.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

func (s *Service) DescribeDatabaseGatewayGateway(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetUserGateways"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"PageSize":   helper.PageSizeLarge,
		"PageNumber": 1,
	}
	idExist := false
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("dg", "2019-03-27", action, nil, request, true)
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
		if fmt.Sprint(response["Success"]) == "false" {
			return object, errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		m, err := jsonpath.Get("$.Data", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
		}
		v, err := convert.ConvertJsonStringToList(m.(string))
		if err != nil {
			return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		}
		if len(v) < 1 {
			return object, errs.WrapErrorf(errs.NotFoundErr("DatabaseGateway", id), errs.NotFoundWithResponse, response)
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
		return object, errs.WrapErrorf(errs.NotFoundErr("DatabaseGateway", id), errs.NotFoundWithResponse, response)
	}
	return
}

func (s *Service) GetUserGatewayInstances(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetUserGatewayInstances"
	request := map[string]interface{}{
		"GatewayId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("dg", "2019-03-27", action, nil, request, true)
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
	if fmt.Sprint(response["Success"]) == "false" {
		return object, errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	m, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}
	v, err := convert.ConvertJsonStringToList(m.(string))
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if len(v) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("DatabaseGateway", id), errs.NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v[0].(map[string]interface{})["gatewayId"]) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr("DatabaseGateway", id), errs.NotFoundWithResponse, response)
		}
	}
	return object, nil
}

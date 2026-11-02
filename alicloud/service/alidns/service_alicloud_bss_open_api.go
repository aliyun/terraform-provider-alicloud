package alidns

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// bssOpenApiService is a local 1:1 copy of the root BssOpenApiService defined
// in alicloud/service_alicloud_bss_open_api.go, limited to QueryAvailableInstances
// (the only method alidns_instance consumes). The root definition is unchanged
// for its existing root consumers; alidns cannot import the root package or the
// service/bssopenapi product package (no reverse or cross-product imports).
// Root symbols are translated to their leaf-package equivalents (helper, errs,
// logs) exactly as the rest of this package's migration did.
type bssOpenApiService struct {
	client *connectivity.AliyunClient
}

func (s *bssOpenApiService) QueryAvailableInstances(id, instanceRegion, productCode, productType, productCodeIntl, productTypeIntl string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	var endpoint string
	action := "QueryAvailableInstances"
	request := map[string]interface{}{
		"ProductCode": productCode,
		"ProductType": productType,
	}
	if client.IsInternationalAccount() {
		request["ProductCode"] = productCodeIntl
		request["ProductType"] = productTypeIntl
	}
	if id != "" {
		request["InstanceIDs"] = id
	}
	if instanceRegion != "" {
		request["Region"] = instanceRegion
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, endpoint)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			if !client.IsInternationalAccount() && errs.IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductCode"] = productCodeIntl
				request["ProductType"] = productTypeIntl
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		resp, _ := jsonpath.Get("$.Data.InstanceList", response)
		if len(resp.([]interface{})) < 1 {
			request["ProductCode"] = productCodeIntl
			if productTypeIntl != "" {
				request["ProductType"] = productTypeIntl
			}
			endpoint = connectivity.BssOpenAPIEndpointInternational
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, endpoint)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data.InstanceList", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data.InstanceList", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr(productCode+"Instance", id), errs.NotFoundWithResponse, response)
	} else if id != "" {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["InstanceID"]) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr(productCode+"Instance", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

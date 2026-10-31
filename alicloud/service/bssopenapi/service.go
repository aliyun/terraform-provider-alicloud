package bssopenapi

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// Service wraps the BssOpenApi (billing) calls shared by product packages.
// It is a leaf package: it depends only on provider-level leaf packages
// (connectivity, errs, helper, logs and provider/conns) and never imports the
// root alicloud package, so product packages can call it directly.
type Service struct {
	client *connectivity.AliyunClient
}

func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

// QueryAvailableInstances returns the billing record of a single instance. The
// lookup is keyed by InstanceIDs only, so no Region is sent: passing an
// international region for a domestic-station account makes the domestic BSS
// return an empty list, which would trigger a blind fallback to the
// international endpoint and fail with "account not exists". When the domestic
// endpoint reports the product code is not applicable, or returns no record,
// the call falls back to the international endpoint and product code.
//
// instanceRegion is kept for signature compatibility with the root
// BssOpenApiService helper; it is intentionally not sent to the API.
func (s *Service) QueryAvailableInstances(id, instanceRegion, productCode, productCodeIntl string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	var endpoint string
	action := "QueryAvailableInstances"
	request := map[string]interface{}{
		"InstanceIDs": id,
		"ProductCode": productCode,
	}
	if client.IsInternationalAccount() {
		request["ProductCode"] = productCodeIntl
	}

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
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		resp, err := jsonpath.Get("$.Data.InstanceList", response)
		if err != nil {
			return retry.NonRetryableError(err)
		}
		if !client.IsInternationalAccount() && len(resp.([]interface{})) < 1 {
			request["ProductCode"] = productCodeIntl
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
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["InstanceID"]) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr(productCode+"Instance", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

package polardb

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// DescribeAvailableInstance reads the billing record of a cluster so the
// subscription release path can reuse the commodity ProductType reported by
// the account instead of hard coding one per station.
func (s *Service) DescribeAvailableInstance(id string) (map[string]interface{}, error) {
	action := "QueryAvailableInstances"
	// The API is keyed by InstanceIDs and is resolved by the account's own
	// station, so Region is deliberately not sent.
	request := map[string]interface{}{
		"ProductCode": "polardb",
		"InstanceIDs": id,
	}
	var response map[string]interface{}
	var endpoint string
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		var err error
		response, err = s.client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, endpoint)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			if !s.client.IsInternationalAccount() && errs.IsExpectedErrors(err, []string{"NotApplicable"}) {
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data.InstanceList", response)
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data.InstanceList", response)
	}
	list, _ := v.([]interface{})
	if len(list) < 1 {
		return nil, errs.WrapErrorf(errs.NotFoundErr("PolarDBInstance", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	instance, _ := list[0].(map[string]interface{})
	if instance == nil || fmt.Sprint(instance["InstanceID"]) != id {
		return nil, errs.WrapErrorf(errs.NotFoundErr("PolarDBInstance", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return instance, nil
}

// RefundPrePaidCluster releases a subscription cluster through the billing
// RefundInstance API and returns the refund order id. A cluster that is
// already refunded reports ResourceNotExists, which the caller treats as an
// idempotent success.
func (s *Service) RefundPrePaidCluster(id string, timeout time.Duration) (int64, error) {
	instance, err := s.DescribeAvailableInstance(id)
	if err != nil {
		return 0, err
	}
	productType := fmt.Sprint(instance["ProductType"])
	if productType == "<nil>" {
		productType = ""
	}

	action := "RefundInstance"
	request := map[string]interface{}{
		"ProductCode":        "polardb",
		"ProductType":        productType,
		"InstanceId":         id,
		"ImmediatelyRelease": "1",
		"ClientToken":        helper.BuildClientToken(action),
	}
	var response map[string]interface{}
	var endpoint string
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(timeout, func() *retry.RetryError {
		var err error
		response, err = s.client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, endpoint)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			if !s.client.IsInternationalAccount() && errs.IsExpectedErrors(err, []string{"NotApplicable"}) {
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return 0, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	// A business rejection can also arrive with HTTP 200 and Success=false;
	// keep the response Code so the caller can classify it like an error code.
	// A response that does not confirm success is a rejection, not a refund.
	success, ok := response["Success"].(bool)
	if !ok || !success {
		code := fmt.Sprint(response["Code"])
		if code == "<nil>" {
			code = ""
		}
		message := fmt.Sprint(response["Message"])
		if message == "<nil>" || message == "" {
			message = "the refund response did not confirm success"
		}
		return 0, errs.WrapErrorf(errs.NewProviderError(code, message), errs.DefaultErrorMsg, id, action, errs.ProviderERROR)
	}

	var orderId int64
	if v, err := jsonpath.Get("$.Data.OrderId", response); err == nil {
		orderId, _ = strconv.ParseInt(fmt.Sprint(v), 10, 64)
	}
	log.Printf("[INFO] Refunded PolarDB cluster %s, refund order id: %d", id, orderId)
	return orderId, nil
}

package mscopensubscription

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

// Service wraps the MscOpenSubscription API calls used by the resources
// and data sources of this package.
type Service struct {
	client *connectivity.AliyunClient
}

// NewService is the exported constructor: the client field stays unexported,
// so callers outside the package (including the flat layout's acceptance
// tests) build values with it.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

func (s *Service) GetContact(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetContact"
	request := map[string]interface{}{
		"ContactId": id,
		"Locale":    "en",
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("MscOpenSubscription", "2021-07-13", action, request, nil)
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
	if errs.IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"ResourceNotFound"}) {
		return object, errs.WrapErrorf(errs.NotFoundErr("MscSub:Contact", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	v, err := jsonpath.Get("$.Contact", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Contact", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) DescribeContact(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetContact"
	request := map[string]interface{}{
		"ContactId": id,
		"Locale":    "en",
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("MscOpenSubscription", "2021-07-13", action, request, nil)
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
		if errs.IsExpectedErrors(err, []string{"ResourceNotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("MscSub:Contact", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Contact", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Contact", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) DescribeSubscription(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetSubscriptionItem"
	request := map[string]interface{}{
		"ItemId": id,
		"Locale": "en",
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("MscOpenSubscription", "2021-07-13", action, request, nil)
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
	v, err := jsonpath.Get("$.SubscriptionItem", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.SubscriptionItem", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) DescribeWebhook(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetWebhook"
	request := map[string]interface{}{
		"WebhookId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("MscOpenSubscription", "2021-07-13", action, request, nil)
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
		if errs.IsExpectedErrors(err, []string{"ResourceNotFound"}) {
			return nil, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Webhook", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Webhook", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

package cms

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// CloudMonitor alert webhooks are served by the Cms ROA API, version 2024-03-30.
// The single-resource CRUD actions are:
//
//	CreateAlertWebhook  POST   /webhook
//	UpdateAlertWebhook  PATCH  /webhook/{webhookId}
//	DeleteAlertWebhooks DELETE /webhooks   (batch; the single id is wrapped into
//	                                       the webhookIds query as a JSON array)
//	ListAlertWebhooks   GET    /webhooks
//
// There is no dedicated single-get endpoint, so the resource Read is implemented
// on top of ListAlertWebhooks with the identifier wrapped into the webhookIds
// query parameter and matched client-side — exactly the inherited PR behavior.
const (
	cmsAPIVersion = "2024-03-30"
	cmsProduct    = "Cms"
	// pageSizeLarge mirrors alicloud.PageSizeLarge without dragging the flat
	// layout's root package into this service package.
	pageSizeLarge = 50
)

// Service wraps the CloudMonitor API calls used by the resources and data
// sources of this package.
type Service struct {
	client *connectivity.AliyunClient
}

// NewService is the exported constructor: the client field stays unexported,
// so callers outside the package (including the flat layout's acceptance tests)
// build values with it.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

func ptr(s string) *string { return &s }

// DescribeCmsAlertWebhook reads a single alert webhook by wrapping its
// identifier into the ListAlertWebhooks webhookIds query parameter (JSON array
// string) and matching client-side. Returns nil when the webhook no longer
// exists, preserving the inherited PR not-found contract (404 /
// ResourceNotFound).
func (s *Service) DescribeCmsAlertWebhook(id string) (map[string]interface{}, error) {
	client := s.client
	webhookIdsJson, err := json.Marshal([]string{id})
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}
	action := "/webhooks"
	query := map[string]*string{
		"webhookIds": ptr(string(webhookIdsJson)),
		"pageNumber": ptr("1"),
		"pageSize":   ptr(fmt.Sprintf("%d", pageSizeLarge)),
	}
	wait := helper.IncrementalWait(3*time.Second, 0)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		var e error
		response, e = client.RoaGet(cmsProduct, cmsAPIVersion, action, query, nil, nil)
		if e != nil {
			if isAlertWebhookNotFound(e) {
				return nil
			}
			if helper.NeedRetry(e) {
				wait()
				return retry.RetryableError(e)
			}
			return retry.NonRetryableError(e)
		}
		return nil
	})
	logs.AddDebug(action, response, query)
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, gerr := jsonpath.Get("$.webhooks[*]", response)
	if gerr != nil {
		if response["webhooks"] == nil {
			// The list endpoint returns no webhooks block for a missing id.
			return nil, nil
		}
		return nil, errs.WrapErrorf(gerr, errs.FailedGetAttributeMsg, id, "$.webhooks[*]", response)
	}
	items, _ := v.([]interface{})
	for _, raw := range items {
		item, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if fmt.Sprint(item["webhookId"]) == id {
			return item, nil
		}
	}
	return nil, nil
}

// isAlertWebhookNotFound reports whether err is the alert-webhook not-found
// error. It mirrors the inherited PR contract IsExpectedErrors(["404",
// "ResourceNotFound"]) plus NotFoundError, so a stale Read during destroy or a
// concurrent delete is treated as gone rather than a hard failure.
func isAlertWebhookNotFound(err error) bool {
	if err == nil {
		return false
	}
	if errs.NotFoundError(err) {
		return true
	}
	if e, ok := err.(*tea.SDKError); ok && e.Code != nil {
		switch tea.StringValue(e.Code) {
		case "404", "ResourceNotFound":
			return true
		}
	}
	msg := err.Error()
	return strings.Contains(msg, "ResourceNotFound")
}

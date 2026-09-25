package arms

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const webhookContactReadTimeout = 1 * time.Minute

// maxExactFloatInteger is the largest integer float64 represents without loss;
// a WebhookId above it cannot be recovered from the number form.
const maxExactFloatInteger = 1 << 53

// WebhookContactId converts a WebhookId from CreateOrUpdateWebhookContact or
// DescribeWebhookContacts into the canonical decimal identity used in state,
// filters and delete requests. The published response contract declares
// WebhookId as type number/format float while some accounts return a numeric
// string, so both are accepted. Containers, fractional numbers, non-finite
// values and numbers outside the lossless integer range are contract drift and
// are rejected instead of being stringified into an id.
func WebhookContactId(raw interface{}) (string, error) {
	switch v := raw.(type) {
	case string:
		if v == "" {
			return "", fmt.Errorf("WebhookId is empty")
		}
		return v, nil
	case json.Number:
		return webhookContactDecimalInteger(v.String())
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return "", fmt.Errorf("WebhookId %v is not finite", v)
		}
		if v != math.Trunc(v) {
			return "", fmt.Errorf("WebhookId %v is not an integer", v)
		}
		if v <= 0 {
			return "", fmt.Errorf("WebhookId %v is not positive", v)
		}
		if math.Abs(v) > maxExactFloatInteger {
			return "", fmt.Errorf("WebhookId %v exceeds the lossless integer range", v)
		}
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case int:
		return webhookContactPositiveInt(int64(v))
	case int64:
		return webhookContactPositiveInt(v)
	case uint64:
		if v == 0 || v > uint64(1<<63-1) {
			return "", fmt.Errorf("WebhookId %d is not a positive int64", v)
		}
		return strconv.FormatUint(v, 10), nil
	case nil:
		return "", fmt.Errorf("WebhookId is missing")
	default:
		return "", fmt.Errorf("WebhookId has unexpected type %T", raw)
	}
}

func webhookContactPositiveInt(v int64) (string, error) {
	if v <= 0 {
		return "", fmt.Errorf("WebhookId %d is not positive", v)
	}
	return strconv.FormatInt(v, 10), nil
}

func webhookContactDecimalInteger(s string) (string, error) {
	if s == "" {
		return "", fmt.Errorf("WebhookId is empty")
	}
	if webhookContactHasNonIntegerRune(s) {
		return "", fmt.Errorf("WebhookId %q is not an integer", s)
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return "", fmt.Errorf("WebhookId %q is not an int64: %w", s, err)
	}
	return webhookContactPositiveInt(v)
}

func webhookContactHasNonIntegerRune(s string) bool {
	for _, c := range s {
		if c == '.' || c == 'e' || c == 'E' || c == '+' {
			return true
		}
	}
	return false
}

// DescribeWebhookContact returns the webhook contact whose WebhookId equals id.
// DescribeWebhookContacts is a collection API, so the target is selected by its
// exact ID instead of the first returned item.
func (s *Service) DescribeWebhookContact(id string) (object map[string]interface{}, err error) {
	return s.DescribeWebhookContactContext(context.Background(), id)
}

func (s *Service) DescribeWebhookContactContext(ctx context.Context, id string) (object map[string]interface{}, err error) {
	client := s.client
	query := map[string]interface{}{
		"ContactIds": id,
		"Page":       1,
		"Size":       helper.PageSizeXLarge,
	}
	action := "DescribeWebhookContacts"

	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(ctx, webhookContactReadTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = client.RpcGet("ARMS", "2019-08-08", action, query, map[string]interface{}{})
		if callErr != nil {
			if helper.NeedRetry(callErr) {
				wait()
				return retry.RetryableError(callErr)
			}
			return retry.NonRetryableError(callErr)
		}
		return nil
	})
	logs.AddDebug(action, response, query)
	if err != nil {
		// A transport-level 404 is definitive. A successful RPC carrying Code 404
		// only in its body is classified after the collection container is
		// validated below, so a malformed body cannot be turned into NotFound.
		if errs.IsTypedNotFound(err) && !errs.IsResponseCodeNotFound(err) {
			return object, errs.WrapErrorf(errs.NotFoundErr("WebhookContact", id), errs.NotFoundMsg, response)
		}
		if !errs.IsResponseCodeNotFound(err) {
			return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	code, _ := jsonpath.Get("$.Code", response)

	// The documented contract returns WebhookContacts as an array. Validate the
	// container before any NotFound classification, including the response-level
	// Code: a null, scalar or object shape is a malformed response, not evidence
	// that the contact no longer exists. Only a confirmed array that is empty or
	// lacks the requested WebhookId means the contact is gone.
	raw, err := jsonpath.Get("$.PageBean.WebhookContacts", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.PageBean.WebhookContacts", response)
	}
	list, ok := raw.([]interface{})
	if !ok {
		return object, errs.WrapErrorf(fmt.Errorf("unexpected WebhookContacts response type %T", raw), errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(code) == "404" || len(list) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("WebhookContact", id), errs.NotFoundMsg, response)
	}

	for _, item := range list {
		contact, ok := item.(map[string]interface{})
		if !ok {
			return object, errs.WrapErrorf(fmt.Errorf("unexpected WebhookContacts element type %T", item), errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		}
		itemId, idErr := WebhookContactId(contact["WebhookId"])
		if idErr != nil {
			return object, errs.WrapErrorf(idErr, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		}
		if itemId == id {
			return contact, nil
		}
	}

	return object, errs.WrapErrorf(errs.NotFoundErr("WebhookContact", id), errs.NotFoundMsg, response)
}

// DeleteWebhookContact deletes the contact and requires the API to confirm the
// deletion. A typed 404 means the contact is already gone.
func (s *Service) DeleteWebhookContact(id string) error {
	return s.DeleteWebhookContactContext(context.Background(), id)
}

func (s *Service) DeleteWebhookContactContext(ctx context.Context, id string) error {
	client := s.client
	action := "DeleteWebhookContact"
	request := map[string]interface{}{
		"WebhookId": id,
	}

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(ctx, webhookContactTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
		if callErr != nil {
			if helper.NeedRetry(callErr) {
				wait()
				return retry.RetryableError(callErr)
			}
			return retry.NonRetryableError(callErr)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsTypedNotFound(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	if success, ok := response["IsSuccess"].(bool); !ok || !success {
		return errs.WrapErrorf(fmt.Errorf("delete response did not confirm IsSuccess=true"), errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

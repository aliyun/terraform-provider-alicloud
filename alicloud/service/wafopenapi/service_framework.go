// Package wafopenapi. Framework-native helpers for the alicloud_waf_defense_asset
// resource. The SDK v2 resources in this package call waf-openapi/2021-10-01 through
// client.RpcPost; the framework resource reuses the same RPC path and the shared
// retry/client-token helpers, so only the framework timeout adapters and the
// expected-error matcher need to live here.
package wafopenapi

import (
	"context"
	"time"

	"github.com/alibabacloud-go/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwutil"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	wafProduct    = "waf-openapi"
	wafAPIVersion = "2021-10-01"
)

// fwTimeoutTypes/fwTimeoutBlock/fwTimeoutDuration adapt the shared fwutil timeout
// helpers to this package's create/update/delete operations.
func fwTimeoutTypes() map[string]attr.Type {
	return fwutil.TimeoutTypes("create", "update", "delete")
}

func fwTimeoutBlock() schema.SingleNestedBlock {
	return fwutil.TimeoutBlock("create", "update", "delete")
}

func fwTimeoutDuration(value types.Object, operation string, fallback time.Duration) (time.Duration, diag.Diagnostics) {
	return fwutil.TimeoutDuration(value, operation, fallback)
}

// wafRetryRpcPost wraps client.RpcPost for the waf-openapi product with the shared
// RetryRpcPost retry loop. The retryable classifier is the same helper.NeedRetry the
// SDK v2 resources in this package use, so transient waf-openapi errors
// (throttling, ServiceUnavailable, OperationConflict, ...) are retried identically.
func wafRetryRpcPost(ctx context.Context, client *connectivity.AliyunClient, action string, request map[string]interface{}, timeout time.Duration) (map[string]interface{}, error) {
	return helper.RetryRpcPost(ctx, action, func() (map[string]interface{}, error) {
		return client.RpcPost(wafProduct, wafAPIVersion, action, nil, request, true)
	}, timeout, helper.NeedRetry)
}

// isExpectedErrors matches structured API codes, unwrapping ComplexError/tea.SDKError.
func isExpectedErrors(err error, codes []string) bool {
	for err != nil {
		if sdkError, ok := err.(*tea.SDKError); ok {
			if sdkError.Code == nil {
				return false
			}
			for _, code := range codes {
				if *sdkError.Code == code {
					return true
				}
			}
			return false
		}
		if wrapped, ok := err.(*errs.ComplexError); ok {
			err = wrapped.Cause
		} else if wrapped, ok := err.(interface{ Unwrap() error }); ok {
			err = wrapped.Unwrap()
		} else {
			return false
		}
	}
	return false
}

// defenseAssetNotExist reports whether err means the defense asset is gone.
// The authoritative online (zhenyuan) DeleteDefenseAsset definition returns
// Defense.Control.DefenseResourceNotExist (HTTP 400) when the asset is already
// gone; Delete and Read treat it as success so the resource is removed from
// state. ResourceNotFound is the generic not-found fallback some gateways emit.
func defenseAssetNotExist(err error) bool {
	return isExpectedErrors(err, []string{"Defense.Control.DefenseResourceNotExist", "ResourceNotFound"})
}

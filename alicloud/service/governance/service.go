// Package governance. This file was generated for the framework-native
// alicloud_governance_service_activation resource and then reduced by hand to the declarations that
// resource executes, so regenerating it has to re-apply the same reduction. The timeout block comes
// from alicloud/provider/fwutil, the retrying write wrapper from alicloud/helper; the state waits
// and the response normalisation stay here.
package governance

import (
	"context"
	"fmt"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwutil"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func fwTimeoutTypes() map[string]attr.Type {
	return fwutil.TimeoutTypes("create", "delete")
}

func fwTimeoutBlock() schema.SingleNestedBlock {
	return fwutil.TimeoutBlock("create", "delete")
}

func fwTimeoutDuration(value types.Object, operation string, fallback time.Duration) (time.Duration, diag.Diagnostics) {
	return fwutil.TimeoutDuration(value, operation, fallback)
}

// retryRpcPost wraps RpcPost in the helper retry loop with the governance retry classifier.
// RpcPost returns *tea.SDKError for RPC failures. ctx carries the lifecycle deadline, so a
// cancelled operation stops retrying.
func retryRpcPost(ctx context.Context, client *connectivity.AliyunClient, product, version, action string, request map[string]interface{}, timeout time.Duration) (map[string]interface{}, error) {
	return helper.RetryRpcPost(ctx, action, func() (map[string]interface{}, error) {
		return client.RpcPost(product, version, action, nil, request, true)
	}, timeout, retryableGovernanceError)
}

// retryableGovernanceError retries only what the published Open/Decommission error tables and the
// shared transient classifier cover. The tables checked on 2026-10-24
// (https://api.aliyun.com/api/governance/2021-01-20/OpenGovernanceService and
// .../DecommissionGovernance) list InternalError (500) for both actions; every other documented
// answer is a terminal 4xx, and retrying a terminal error only delays the same failure by the whole
// operation timeout. Answers that mean "a change is already running" or "already in the target
// state" are resolved by the callers' status-convergence loops instead of a blind retry, so only
// InternalError and the shared classifier fire here.
func retryableGovernanceError(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok {
		if e.Code == nil || e.Message == nil {
			return false
		}
		if *e.Code == "InternalError" {
			return true
		}
	}
	return helper.NeedRetry(err)
}

// statusQueryTimeout bounds a single status query when the caller's context carries no deadline.
const statusQueryTimeout = 1 * time.Minute

// statusQueryBudget returns how long a status query may retry: the whole statusQueryTimeout, or the
// remaining lifecycle budget when ctx already carries an earlier deadline.
func statusQueryBudget(ctx context.Context) time.Duration {
	timeout := statusQueryTimeout
	if deadline, ok := ctx.Deadline(); ok {
		if remaining := time.Until(deadline); remaining < timeout {
			return remaining
		}
	}
	return timeout
}

// serviceActivationGet describes the singleton activation by its region id. It returns (nil, nil)
// when the response carries no object; a deactivated service is still described, with Status
// "Disabled", because GetGovernanceServiceStatus keeps answering after decommission.
func serviceActivationGet(ctx context.Context, client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	if id == "" || id != client.RegionId {
		return nil, fmt.Errorf("building GetGovernanceServiceStatus request: the activation region %q does not match the provider region %q", id, client.RegionId)
	}
	budget := statusQueryBudget(ctx)
	if budget <= 0 {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		return nil, context.DeadlineExceeded
	}
	request := map[string]interface{}{"RegionId": id}
	response, err := retryRpcPost(ctx, client, "governance", "2021-01-20", "GetGovernanceServiceStatus", request, budget)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, nil
	}
	if _, ok := response["ServiceStatus"]; !ok {
		return nil, fmt.Errorf("selecting GetGovernanceServiceStatus resource: the response carries no ServiceStatus")
	}
	return serviceActivationCompleteRead(response, id), nil
}

// serviceActivationCompleteRead renames the response field to the attribute name, because
// GetGovernanceServiceStatus answers ServiceStatus while the schema attribute is status.
func serviceActivationCompleteRead(object map[string]interface{}, id string) map[string]interface{} {
	object["Status"] = object["ServiceStatus"]
	object["Region"] = id
	return object
}

// serviceActivationStateRefreshFunc polls the activation status. A nil object reports state "" so
// the caller's wait can distinguish "no object" from every declared status.
func serviceActivationStateRefreshFunc(ctx context.Context, client *connectivity.AliyunClient, id, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := serviceActivationGet(ctx, client, id)
		if err != nil {
			return nil, "", err
		}
		if object == nil {
			return nil, "", nil
		}
		current := ""
		if value, ok := object[field]; ok && value != nil {
			current = fmt.Sprint(value)
		}
		for _, fs := range failStates {
			if current == fs {
				return object, current, fmt.Errorf("status reached fail state %q", current)
			}
		}
		return object, current, nil
	}
}

// serviceActivationWaitForState polls field until it reaches one of target; a value in failStates
// aborts the wait with an error. Returns the last polled object, or nil when no object was seen.
// ctx bounds the wait as well as the timeout, so cancelling the operation stops the polling.
func serviceActivationWaitForState(ctx context.Context, client *connectivity.AliyunClient, id, field string, target, failStates []string, timeout time.Duration) (map[string]interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{},
		Target:     target,
		Refresh:    serviceActivationStateRefreshFunc(ctx, client, id, field, failStates),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	raw, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, err
	}
	if m, ok := raw.(map[string]interface{}); ok {
		return m, nil
	}
	return nil, nil
}

// serviceActivationAbsent reports the statuses that leave no activation behind: Disabled after a
// decommission, and EnableFailed after the automatic rollback that OpenGovernanceService documents
// for a failed activation (a decommission then finds no enabled activation).
func serviceActivationAbsent(status string) bool {
	return status == "Disabled" || status == "EnableFailed"
}

// serviceActivationFailureDetail renders the FailureCode/FailureMessage that
// GetGovernanceServiceStatus returns for EnableFailed/DisableFailed, so a terminal state wait
// reports why it gave up instead of only the last status.
func serviceActivationFailureDetail(object map[string]interface{}) string {
	code, _ := object["FailureCode"].(string)
	message, _ := object["FailureMessage"].(string)
	if code == "" && message == "" {
		return ""
	}
	return fmt.Sprintf(" (failure code: %q, failure message: %q)", code, message)
}

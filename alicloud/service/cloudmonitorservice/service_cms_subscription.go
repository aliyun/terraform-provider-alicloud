package cloudmonitorservice

import (
	"context"
	stderrors "errors"
	"fmt"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// errSubscriptionNotFound marks the contract-level Cms ResourceNotFound. The
// resource only clears state for this signal; a bare 404 without the documented
// code, or a null payload, stays an error.
var errSubscriptionNotFound = stderrors.New("the specified Cms subscription is not found")

func subscriptionNotFound(err error) bool {
	return stderrors.Is(err, errSubscriptionNotFound)
}

// IsSubscriptionNotFound reports whether err is the contract-level absence
// signal returned by DescribeSubscription.
func IsSubscriptionNotFound(err error) bool {
	return subscriptionNotFound(err)
}

func subscriptionIsResourceNotFound(err error) bool {
	switch typed := err.(type) {
	case *tea.SDKError:
		return tea.StringValue(typed.Code) == "ResourceNotFound"
	case *errs.ProviderError:
		return typed.ErrorCode() == "ResourceNotFound"
	case *errs.ComplexError:
		return subscriptionIsResourceNotFound(typed.Cause)
	}
	return false
}

// DescribeSubscription fetches one Cms 2024-03-30 event subscription. The
// workspace query parameter is optional on the API; it is forwarded when the
// caller has it so that workspace-scoped subscriptions resolve consistently.
// The caller's context carries the operation deadline shared by the request,
// its retries and the polling loops.
func (s *CmsService) DescribeSubscription(ctx context.Context, id, workspace string) (object map[string]interface{}, err error) {
	client := s.client
	action := fmt.Sprintf("/subscriptions/%s", id)
	query := make(map[string]*string)
	if workspace != "" {
		query["workspace"] = helper.StringPointer(workspace)
	}

	var response map[string]interface{}
	notFound := false
	wait := subscriptionBackoff(ctx)
	err = retry.RetryContext(ctx, subscriptionTimeout, func() *retry.RetryError {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return retry.NonRetryableError(ctxErr)
		}
		response, err = client.RoaGet("Cms", "2024-03-30", action, query, nil, nil)
		if err != nil {
			if helper.NeedRetry(err) {
				if !wait() {
					return retry.NonRetryableError(ctx.Err())
				}
				return retry.RetryableError(err)
			}
			if subscriptionIsResourceNotFound(err) {
				notFound = true
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, query)
	if notFound {
		return nil, errSubscriptionNotFound
	}
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	rawSubscription, exists := response["subscription"]
	if !exists {
		return nil, errs.WrapErrorf(
			fmt.Errorf("the response does not contain $.subscription"),
			errs.FailedGetAttributeMsg, id, "$.subscription", response,
		)
	}
	subscription, ok := rawSubscription.(map[string]interface{})
	if !ok {
		// The documented success response always carries the subscription object.
		// A null or non-object payload is off-contract and must not clear state.
		return nil, errs.WrapErrorf(
			fmt.Errorf("unexpected $.subscription payload type %T", rawSubscription),
			errs.FailedGetAttributeMsg, id, "$.subscription", response,
		)
	}
	return subscription, nil
}

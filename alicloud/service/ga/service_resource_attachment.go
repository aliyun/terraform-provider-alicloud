package ga

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	resourceAttachmentCreateTimeout = 5 * time.Minute
	resourceAttachmentDeleteTimeout = 5 * time.Minute
	resourceAttachmentPollInterval  = 5 * time.Second
)

// getResourceAttachment returns the association record whose four key fields
// all match, or nil when the association no longer exists. The list API is
// filtered server side, but the identity check is repeated locally so a record
// for another accelerator or resource can never be written into this state.
func getResourceAttachment(ctx context.Context, client *connectivity.AliyunClient, acceleratorId, resourceType, resourceRegionId, resourceId string) (map[string]interface{}, error) {
	request := map[string]interface{}{
		"AcceleratorId":              acceleratorId,
		"RegionId":                   client.RegionId,
		"AssociatedResourceType":     resourceType,
		"AssociatedResourceRegionId": resourceRegionId,
		"AssociatedResourceId":       resourceId,
	}
	var response map[string]interface{}
	err := retry.RetryContext(ctx, resourceAttachmentCreateTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = client.RpcPost("Ga", "2019-11-20", "GetGlobalAcceleratorResources", nil, request, true)
		if callErr != nil {
			if isAcceleratorNotExist(callErr) {
				return retry.NonRetryableError(callErr)
			}
			if helper.NeedRetry(callErr) {
				if sleepErr := sleepWithContext(ctx, resourceAttachmentPollInterval); sleepErr != nil {
					return retry.NonRetryableError(sleepErr)
				}
				return retry.RetryableError(callErr)
			}
			return retry.NonRetryableError(callErr)
		}
		return nil
	})
	if err != nil {
		if isAcceleratorNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	records, _ := response["AssociatedResources"].([]interface{})
	for _, item := range records {
		record, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if fmt.Sprint(record["AcceleratorId"]) == acceleratorId &&
			fmt.Sprint(record["AssociatedResourceType"]) == resourceType &&
			fmt.Sprint(record["AssociatedResourceRegionId"]) == resourceRegionId &&
			fmt.Sprint(record["AssociatedResourceId"]) == resourceId {
			return record, nil
		}
	}
	return nil, nil
}

// waitResourceAttachmentActive polls until the association reports active. The
// record can stay invisible for a short while after AssociateResources returns.
func waitResourceAttachmentActive(ctx context.Context, client *connectivity.AliyunClient, acceleratorId, resourceType, resourceRegionId, resourceId string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := retry.RetryContext(ctx, resourceAttachmentCreateTimeout, func() *retry.RetryError {
		record, err := getResourceAttachment(ctx, client, acceleratorId, resourceType, resourceRegionId, resourceId)
		if err != nil {
			if helper.NeedRetry(err) {
				if sleepErr := sleepWithContext(ctx, resourceAttachmentPollInterval); sleepErr != nil {
					return retry.NonRetryableError(sleepErr)
				}
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		if record == nil {
			if sleepErr := sleepWithContext(ctx, resourceAttachmentPollInterval); sleepErr != nil {
				return retry.NonRetryableError(sleepErr)
			}
			return retry.RetryableError(fmt.Errorf("the association %s is not visible yet", resourceAttachmentId(acceleratorId, resourceType, resourceRegionId, resourceId)))
		}
		state := fmt.Sprint(record["State"])
		if state == "active" {
			result = record
			return nil
		}
		if state == "wafRelease" {
			return retry.NonRetryableError(fmt.Errorf("the WAF instance %s bound to accelerator %s has been released", resourceId, acceleratorId))
		}
		if sleepErr := sleepWithContext(ctx, resourceAttachmentPollInterval); sleepErr != nil {
			return retry.NonRetryableError(sleepErr)
		}
		return retry.RetryableError(fmt.Errorf("the association %s is in state %s, waiting for active", resourceAttachmentId(acceleratorId, resourceType, resourceRegionId, resourceId), state))
	})
	return result, err
}

// waitResourceAttachmentDeleted polls until the association is no longer
// returned by the list API.
func waitResourceAttachmentDeleted(ctx context.Context, client *connectivity.AliyunClient, acceleratorId, resourceType, resourceRegionId, resourceId string) error {
	return retry.RetryContext(ctx, resourceAttachmentDeleteTimeout, func() *retry.RetryError {
		record, err := getResourceAttachment(ctx, client, acceleratorId, resourceType, resourceRegionId, resourceId)
		if err != nil {
			if isAcceleratorNotExist(err) {
				return nil
			}
			if helper.NeedRetry(err) {
				if sleepErr := sleepWithContext(ctx, resourceAttachmentPollInterval); sleepErr != nil {
					return retry.NonRetryableError(sleepErr)
				}
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		if record == nil {
			return nil
		}
		if sleepErr := sleepWithContext(ctx, resourceAttachmentPollInterval); sleepErr != nil {
			return retry.NonRetryableError(sleepErr)
		}
		return retry.RetryableError(fmt.Errorf("the association %s still exists", resourceAttachmentId(acceleratorId, resourceType, resourceRegionId, resourceId)))
	})
}

func resourceAttachmentId(acceleratorId, resourceType, resourceRegionId, resourceId string) string {
	return strings.Join([]string{acceleratorId, resourceType, resourceRegionId, resourceId}, ":")
}

func parseResourceAttachmentId(id string) (acceleratorId, resourceType, resourceRegionId, resourceId string, err error) {
	parts := strings.Split(id, ":")
	if len(parts) != 4 {
		return "", "", "", "", fmt.Errorf("expected four colon-separated segments, got %d: %q", len(parts), id)
	}
	for _, part := range parts {
		if part == "" {
			return "", "", "", "", fmt.Errorf("every segment must be non-empty: %q", id)
		}
	}
	return parts[0], parts[1], parts[2], parts[3], nil
}

func sleepWithContext(ctx context.Context, d time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(d):
		return nil
	}
}

// isAcceleratorNotExist reports the error returned when the accelerator itself
// is gone, which also means its associations are gone.
func isAcceleratorNotExist(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok && e.Code != nil && *e.Code == "NotExist.Accelerator" {
		return true
	}
	return strings.Contains(err.Error(), "NotExist.Accelerator")
}

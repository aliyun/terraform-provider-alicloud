package live

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// livePrivateLineSep separates the four identity parts encoded into the
// composite id: DomainName, AppName, StreamName and AccelerationType. None of
// these values carry a colon in practice, matching the existing live resources.
const livePrivateLineSep = ":"

// buildLivePrivateLineID packs the resource identity. There is no single
// server-side id: CreateLivePrivateLine returns only RequestId, and the
// Describe entry is filtered from a list by the same four fields.
func buildLivePrivateLineID(domainName, appName, streamName, accelerationType string) string {
	return fmt.Sprintf("%s%s%s%s%s%s%s", domainName, livePrivateLineSep, appName, livePrivateLineSep, streamName, livePrivateLineSep, accelerationType)
}

// parseLivePrivateLineID splits the composite id back into its identity parts.
func parseLivePrivateLineID(id string) (domainName, appName, streamName, accelerationType string, err error) {
	parts := strings.Split(id, livePrivateLineSep)
	if len(parts) != 4 {
		return "", "", "", "", errs.WrapErrorf(errs.NotFoundErr("LivePrivateLine", id), errs.NotFoundMsg, id)
	}
	return parts[0], parts[1], parts[2], parts[3], nil
}

// DescribeLivePrivateLine <<< Encapsulated get interface for Live PrivateLine.
//
// CreateLivePrivateLine returns only RequestId, so the resource has no single
// server-side id. DescribeLivePrivateLineAvailGA lists every binding for a
// stream; the unique entry is the one whose AccelerationType matches the
// encoded identity (Delete operates on the same four fields and ignores
// AccelerationArea, confirming the four-tuple is unique).
func (s *Service) DescribeLivePrivateLine(id string) (object map[string]interface{}, err error) {
	client := s.client
	domainName, appName, streamName, accelerationType, err := parseLivePrivateLineID(id)
	if err != nil {
		return object, err
	}

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DomainName"] = domainName
	request["AppName"] = appName
	request["StreamName"] = streamName
	// IsGaInstance=no returns the binding detail (not the bare GA status).
	request["IsGaInstance"] = "no"
	request["RegionId"] = client.RegionId
	action := "DescribeLivePrivateLineAvailGA"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("live", "2016-11-01", action, query, request, true)
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
		if errs.IsExpectedErrors(err, []string{"InvalidDomain.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("LivePrivateLine", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.LivePrivateLineAvailGAs.LivePrivateLineAvailGA[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.LivePrivateLineAvailGAs.LivePrivateLineAvailGA[*]", response)
	}

	items, ok := v.([]interface{})
	if !ok || len(items) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("LivePrivateLine", id), errs.NotFoundMsg, response)
	}

	// The resource models the binding between a GA instance and a live link, not
	// the GA instance itself. DeleteLivePrivateLine only unbinds: the GA may keep
	// surfacing in the list with BindingStatus=no, and a not-yet-bound entry can
	// appear transiently right after Create. Only a bound entry (BindingStatus=yes)
	// represents an existing private line, so anything else is NotFound; this
	// makes Read drop drift, the create waiter retry until bound, the delete
	// waiter converge, and CheckDestroy pass once unbound.
	for _, row := range items {
		item, ok := row.(map[string]interface{})
		if !ok {
			continue
		}
		if fmt.Sprint(item["AccelerationType"]) == accelerationType && fmt.Sprint(item["BindingStatus"]) == "yes" {
			return item, nil
		}
	}

	return object, errs.WrapErrorf(errs.NotFoundErr("LivePrivateLine", id), errs.NotFoundMsg, response)
}

// LivePrivateLineStateRefreshFunc waits for the binding to surface and reach
// the active state after Create. Right after CreateLivePrivateLine the list
// entry may not be visible yet (Reuse=no also provisions a GA instance), so
// NotFound is a pending state rather than terminal.
func (s *Service) LivePrivateLineStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeLivePrivateLine(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Not yet visible after create; keep waiting.
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		currentStatus := fmt.Sprint(object["Status"])
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeLivePrivateLine >>> Encapsulated.

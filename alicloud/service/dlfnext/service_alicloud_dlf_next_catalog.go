package dlfnext

import (
	"context"
	"fmt"
	"math/rand"
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

type Service struct {
	client *connectivity.AliyunClient
}

func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

// DescribeCatalog reads one catalog. The default 1-minute retry budget is for
// callers that do not share a Create/Delete deadline.
func (s *Service) DescribeCatalog(ctx context.Context, id string) (object map[string]interface{}, err error) {
	return s.DescribeCatalogWithRetryTimeout(ctx, id, 1*time.Minute)
}

// DescribeCatalogWithRetryTimeout is the deadline-aware get. retryTimeout
// bounds the retry loop; an in-flight request is not preempted, so one
// trailing call may slightly exceed it.
func (s *Service) DescribeCatalogWithRetryTimeout(ctx context.Context, id string, retryTimeout time.Duration) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string

	action := catalogAction(id)
	request = make(map[string]interface{})
	query = make(map[string]*string)

	backoff := newCatalogBackoff(3*time.Second, 30*time.Second)
	err = retry.RetryContext(ctx, retryTimeout, func() *retry.RetryError {
		response, err = client.RoaGet("DlfNext", "2025-03-10", action, query, nil, nil)
		if err != nil {
			if helper.NeedRetry(err) {
				if waitErr := backoff.Wait(ctx, retryTimeout); waitErr != nil {
					return retry.NonRetryableError(waitErr)
				}
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"CatalogNotFound", "NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("DlfNextCatalog", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

// ListCatalogs lists one page of catalogs.
func (s *Service) ListCatalogs(ctx context.Context, catalogNamePattern string, pageToken string, maxResults int) (objects []interface{}, nextToken string, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	action := "/dlf/v1/catalogs"
	request = make(map[string]interface{})
	query := make(map[string]*string)

	if catalogNamePattern != "" {
		v := catalogNamePattern
		query["catalogNamePattern"] = &v
	}
	if pageToken != "" {
		v := pageToken
		query["pageToken"] = &v
	}
	if maxResults > 0 {
		v := fmt.Sprint(maxResults)
		query["maxResults"] = &v
	}

	backoff := newCatalogBackoff(3*time.Second, 30*time.Second)
	err = retry.RetryContext(ctx, 1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("DlfNext", "2025-03-10", action, query, nil, nil)
		if err != nil {
			if helper.NeedRetry(err) {
				if waitErr := backoff.Wait(ctx, 1*time.Minute); waitErr != nil {
					return retry.NonRetryableError(waitErr)
				}
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return objects, nextToken, errs.WrapErrorf(err, errs.DefaultErrorMsg, "list", action, errs.AlibabaCloudSdkGoERROR)
	}

	nextToken = ""
	if v, ok := response["nextPageToken"]; ok && v != nil {
		nextToken = fmt.Sprint(v)
	}
	if v, ok := response["catalogs"].([]interface{}); ok {
		objects = v
	}
	return objects, nextToken, nil
}

// CatalogStateRefreshFuncWithApi returns a refresh function backed by the
// injectable call, so waiter policy can be unit-tested without a backend.
func (s *Service) CatalogStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Brief NotFound after create is Pending (eventual consistency);
				// the normal Read still clears the id on NotFound.
				return object, "NOT_FOUND", nil
			}
			return nil, "", errs.WrapError(err)
		}
		v, err := jsonpath.Get(field, object)
		if err != nil {
			return nil, "", errs.WrapErrorf(err, "reading %s from DlfNext Catalog %s", field, id)
		}
		currentStatus := fmt.Sprint(v)
		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// WaitForCatalogGone polls GetCatalog until the catalog returns 404.
func (s *Service) WaitForCatalogGone(ctx context.Context, id string, timeout time.Duration) error {
	return s.WaitForCatalogGoneWithApi(ctx, id, timeout, func(id string) (map[string]interface{}, error) {
		return s.DescribeCatalog(ctx, id)
	})
}

// WaitForCatalogGoneWithApi is the injectable form used by tests.
func (s *Service) WaitForCatalogGoneWithApi(ctx context.Context, id string, timeout time.Duration, call func(id string) (map[string]interface{}, error)) error {
	backoff := newCatalogBackoff(3*time.Second, 30*time.Second)
	return retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		_, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return nil
			}
			if helper.NeedRetry(err) {
				if waitErr := backoff.Wait(ctx, timeout); waitErr != nil {
					return retry.NonRetryableError(waitErr)
				}
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return retry.RetryableError(fmt.Errorf("DlfNext Catalog %s is still being deleted", id))
	})
}

// isCatalogBeingCreatedError matches DropCatalog's 403 while the catalog is
// still in NEW/INITIALIZING; NeedRetry does not cover this status.
func isCatalogBeingCreatedError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "is being created")
}

// isCatalogAlreadyExistsError reports an explicit already-exists response.
func isCatalogAlreadyExistsError(err error) bool {
	if err == nil {
		return false
	}
	if errs.IsExpectedErrors(err, []string{"AlreadyExist", "CatalogAlreadyExists", "DuplicateCatalog", "Duplicate"}) {
		return true
	}
	return strings.Contains(strings.ToLower(err.Error()), "already exist")
}

// isCatalogExplicitClientError reports a client-side validation/permission
// error (HTTP 400/401/403) that proves the create could not have succeeded.
// Throttling is matched by NeedRetry regardless of status: the backend can
// deliver throttling with HTTP 400, and those must stay retryable.
func isCatalogExplicitClientError(err error) bool {
	if err == nil {
		return false
	}
	if isCatalogAlreadyExistsError(err) {
		return false
	}
	if isCatalogBeingCreatedError(err) {
		return false
	}
	e, ok := err.(*tea.SDKError)
	if !ok || e.StatusCode == nil {
		return false
	}
	switch tea.IntValue(e.StatusCode) {
	case 400, 401, 403:
		return !helper.NeedRetry(err)
	}
	return false
}

// isCatalogBackendCapacityError matches the specific DLF cleanup-quota /
// environment-capacity phrases. It is intentionally narrow: generic wording
// such as "exceed the maximum" also appears in throttling responses.
func isCatalogBackendCapacityError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	for _, sig := range []string{
		"too many catalogs",
		"catalogs being cleaned up",
		"approximately 48 hours",
		"catalog quota",
	} {
		if strings.Contains(msg, sig) {
			return true
		}
	}
	return false
}

// catalogBackoff is an exponential backoff with +/-25% jitter. Wait honors
// context cancellation and never sleeps past the caller's budget.
type catalogBackoff struct {
	current time.Duration
	max     time.Duration
}

func newCatalogBackoff(initial, max time.Duration) *catalogBackoff {
	return &catalogBackoff{current: initial, max: max}
}

func (b *catalogBackoff) Wait(ctx context.Context, budget time.Duration) error {
	sleep := b.current
	jitterRange := sleep / 4
	if jitterRange > 0 {
		sleep = sleep + time.Duration(rand.Int63n(int64(2*jitterRange+1))) - jitterRange
	}
	if sleep < 0 {
		sleep = 0
	}
	if budget > 0 && sleep > budget {
		sleep = budget
	}
	if sleep > 0 {
		timer := time.NewTimer(sleep)
		defer timer.Stop()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
	}
	b.current *= 2
	if b.current > b.max {
		b.current = b.max
	}
	return nil
}

// Package dlfnext. This file is generated automatically. Please do not modify it manually, thank you!
package dlfnext

import (
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

// DescribeCatalog <<< Encapsulated get interface for DlfNext Catalog.
// Uses a 1-minute retry budget by default, preserved for callers that do not
// share a Create/Delete deadline (Read, the data source, ad-hoc existence
// checks). The Create and Delete paths call DescribeCatalogWithRetryTimeout
// with the time remaining in their shared deadline so this fixed 1-minute retry
// does not nest inside an outer deadline and exceed it.
func (s *Service) DescribeCatalog(id string) (object map[string]interface{}, err error) {
	return s.DescribeCatalogWithRetryTimeout(id, 1*time.Minute)
}

// DescribeCatalogWithRetryTimeout is the deadline-aware get: the retry budget
// is the caller-supplied retryTimeout (typically time.Until(deadline) from a
// Create/Delete that established a single shared deadline) instead of the
// hardcoded 1 minute. This bounds the retry loop so a post-error reconciliation
// Get or a waiter poll nested inside an outer deadline cannot consume a full
// minute per call and blow past the outer budget.
//
// Guarantee boundary: the retry budget bounds the RETRY loop, not a single
// in-flight HTTP call. retry.Retry stops starting new attempts once
// retryTimeout is reached, but a RoaGet already in flight is not interrupted;
// each RoaGet has its own client-level HTTP timeout, so a single trailing call
// may slightly exceed retryTimeout. This is the practical bound, not a hard
// preempt of in-flight I/O.
func (s *Service) DescribeCatalogWithRetryTimeout(id string, retryTimeout time.Duration) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string

	catalog := id
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/dlf/v1/catalogs/%s", catalog)

	wait := catalogBackoffWait(3*time.Second, 30*time.Second)
	err = retry.Retry(retryTimeout, func() *retry.RetryError {
		response, err = client.RoaGet("DlfNext", "2025-03-10", action, query, nil, nil)
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
		if errs.IsExpectedErrors(err, []string{"CatalogNotFound", "NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("DlfNextCatalog", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

// ListCatalogs <<< Encapsulated list interface for DlfNext Catalog.
func (s *Service) ListCatalogs(catalogNamePattern string, pageToken string, maxResults int) (objects []interface{}, nextToken string, err error) {
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

	wait := catalogBackoffWait(3*time.Second, 30*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("DlfNext", "2025-03-10", action, query, nil, nil)
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

func (s *Service) CatalogStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.CatalogStateRefreshFuncWithApi(id, field, failStates, s.DescribeCatalog)
}

func (s *Service) CatalogStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Surface a stable NOT_FOUND state instead of the empty string.
				// The Create waiter lists NOT_FOUND as a Pending state so a brief
				// NotFound immediately after create (the catalog is not yet
				// visible to Get due to eventual consistency) is retried until
				// the catalog appears, rather than treated as an unexpected
				// state. This is distinct from the normal Resource Read, which
				// still clears the TF id on NotFound.
				return object, "NOT_FOUND", nil
			}
			return nil, "", errs.WrapError(err)
		}
		v, err := jsonpath.Get(field, object)
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

// isCatalogBeingCreatedError reports whether the error returned by DropCatalog
// indicates the catalog is still being created and cannot be dropped yet.
// Catalog creation is asynchronous on the DLF backend: right after
// CreateCatalog the catalog is in NEW/INITIALIZING status and DropCatalog
// rejects the deletion with HTTP 403 "... is being created, can not be
// dropped. Please wait for the catalog to be ready.". NeedRetry does not match
// this (it is a 403, not throttling/5xx), so the Delete path must retry it
// explicitly until the catalog reaches a ready state.
func isCatalogBeingCreatedError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "is being created")
}

// WaitForCatalogGone polls GetCatalog until the catalog is gone
// (CatalogNotFound). DropCatalog is asynchronous: it returns 200 while the
// catalog transitions through DELETING/DELETED before finally returning 404.
// A dedicated poll (rather than CatalogStateRefreshFunc) avoids the fragility of
// matching the transient DELETED status, which is neither a stable pending
// state nor a terminal "gone" signal in the generic state-refresh helper: here
// any non-NotFound response simply keeps polling until the catalog returns 404
// or the timeout elapses.
func (s *Service) WaitForCatalogGone(id string, timeout time.Duration) error {
	return s.WaitForCatalogGoneWithApi(id, timeout, s.DescribeCatalog)
}

// WaitForCatalogGoneWithApi polls GetCatalog until the catalog is gone
// (CatalogNotFound). The get function is injectable so the delete-waiter policy
// can be unit-tested without a live backend. DropCatalog is asynchronous: it
// returns 200 while the catalog transitions through DELETING/DELETED before
// finally returning 404. Any non-NotFound response simply keeps polling until
// the catalog returns 404 or the timeout elapses; a transient retryable Get
// error is retried, a non-retryable Get error fails.
func (s *Service) WaitForCatalogGoneWithApi(id string, timeout time.Duration, call func(id string) (map[string]interface{}, error)) error {
	wait := catalogBackoffWait(3*time.Second, 30*time.Second)
	return retry.Retry(timeout, func() *retry.RetryError {
		_, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return nil
			}
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		// Catalog still exists (DELETING/DELETED); keep polling until 404.
		return retry.RetryableError(fmt.Errorf("DlfNext Catalog %s is still being deleted", id))
	})
}

// isCatalogAlreadyExistsError reports whether err is a CreateCatalog "catalog
// already exists" response. CreateCatalog carries no idempotency token, so an
// explicit already-exists signal means a catalog with this name is already
// present and must not be adopted by this run.
func isCatalogAlreadyExistsError(err error) bool {
	if err == nil {
		return false
	}
	if errs.IsExpectedErrors(err, []string{"AlreadyExist", "CatalogAlreadyExists", "DuplicateCatalog", "Duplicate"}) {
		return true
	}
	return strings.Contains(strings.ToLower(err.Error()), "already exist")
}

// isCatalogExplicitClientError reports whether err is a client-side validation
// or permission error (HTTP 400/401/403) that proves the create could not have
// succeeded. Throttling and the Delete path's "being created" 403 are excluded,
// as is an explicit already-exists error (handled separately). Such errors must
// never trigger adoption of a same-name catalog found afterwards.
//
// Throttling is matched by helper.NeedRetry regardless of HTTP status: the
// backend can deliver a throttling error with HTTP 400, and treating that as a
// validation/permission error would fail a create that should have been
// retried. NeedRetry-class errors are therefore never "explicit client" errors.
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

// isCatalogBackendCapacityError reports whether err signals a backend
// cleanup-quota / environment-capacity limitation, NOT a transient API
// throttle. DLF CreateCatalog carries no idempotency token, and the backend
// rejects new catalogs when the account+Region has too many recently-created
// catalogs still being cleaned up (the backend reaps them asynchronously,
// ~48h). The signals are: an error message containing the backend capacity
// phrases, or a catalog that reached INITIALIZE_FAILED after the create was
// accepted. This is NOT retried: re-issuing CreateCatalog would only add
// another catalog to the cleanup queue and worsen the quota pressure. The
// caller surfaces a clear INFRA_BLOCKED error with account/Region/Catalog
// name/RequestId/server message so the test runner can classify it as an
// environment capacity issue (not a regression, not a PASS).
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
		"exceed the maximum",
	} {
		if strings.Contains(msg, sig) {
			return true
		}
	}
	return false
}

// catalogBackoffWait returns a closure that sleeps for an exponentially
// increasing duration with a random jitter, so multiple ACC runners hitting the
// same DLF backend do not retry in lockstep and hammer the API. The base grows
// geometrically from initial up to max, and each sleep is shaken by up to +/-
// 25% jitter. Bounded: once max is reached, subsequent sleeps stay at max
// (+jitter).
func catalogBackoffWait(initial, max time.Duration) func() {
	current := initial
	return func() {
		sleep := current
		// jitter: +/- 25% of the current sleep, clamped at >= 0.
		jitterRange := sleep / 4
		if jitterRange > 0 {
			sleep = sleep + time.Duration(rand.Int63n(int64(2*jitterRange+1))) - jitterRange
			if sleep < 0 {
				sleep = 0
			}
		}
		time.Sleep(sleep)
		// Grow geometrically, capped at max.
		current *= 2
		if current > max {
			current = max
		}
	}
}

// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type DlfNextServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeDlfNextCatalog <<< Encapsulated get interface for DlfNext Catalog.
// Uses a 1-minute retry budget by default, preserved for callers that do not
// share a Create/Delete deadline (Read, the data source, ad-hoc existence
// checks). The Create and Delete paths call DescribeDlfNextCatalogWithRetryTimeout
// with the time remaining in their shared deadline so this fixed 1-minute retry
// does not nest inside an outer deadline and exceed it.
func (s *DlfNextServiceV2) DescribeDlfNextCatalog(id string) (object map[string]interface{}, err error) {
	return s.DescribeDlfNextCatalogWithRetryTimeout(id, 1*time.Minute)
}

// DescribeDlfNextCatalogWithRetryTimeout is the deadline-aware get: the retry
// budget is the caller-supplied retryTimeout (typically time.Until(deadline)
// from a Create/Delete that established a single shared deadline) instead of
// the hardcoded 1 minute. This bounds the retry loop so a post-error
// reconciliation Get or a waiter poll nested inside an outer deadline cannot
// consume a full minute per call and blow past the outer budget.
//
// Guarantee boundary: the retry budget bounds the RETRY loop, not a single
// in-flight HTTP call. resource.Retry stops starting new attempts once
// retryTimeout is reached, but an RoaGet already in flight is not interrupted;
// each RoaGet has its own client-level HTTP timeout, so a single trailing
// call may slightly exceed retryTimeout. This is the practical bound, not a
// hard preempt of in-flight I/O.
func (s *DlfNextServiceV2) DescribeDlfNextCatalogWithRetryTimeout(id string, retryTimeout time.Duration) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string

	catalog := id
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/dlf/v1/catalogs/%s", catalog)

	wait := dlfNextCatalogBackoffWait(3*time.Second, 30*time.Second)
	err = resource.Retry(retryTimeout, func() *resource.RetryError {
		response, err = client.RoaGet("DlfNext", "2025-03-10", action, query, nil, nil)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"CatalogNotFound", "NotFound"}) {
			return object, WrapErrorf(NotFoundErr("DlfCatalog", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

// ListDlfNextCatalogs <<< Encapsulated list interface for DlfNext Catalog.
func (s *DlfNextServiceV2) ListDlfNextCatalogs(catalogNamePattern string, pageToken string, maxResults int) (objects []interface{}, nextToken string, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	action := fmt.Sprintf("/dlf/v1/catalogs")
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

	wait := dlfNextCatalogBackoffWait(3*time.Second, 30*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("DlfNext", "2025-03-10", action, query, nil, nil)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return objects, nextToken, WrapErrorf(err, DefaultErrorMsg, "list", action, AlibabaCloudSdkGoERROR)
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

func (s *DlfNextServiceV2) DlfNextCatalogStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.DlfNextCatalogStateRefreshFuncWithApi(id, field, failStates, s.DescribeDlfNextCatalog)
}

func (s *DlfNextServiceV2) DlfNextCatalogStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if NotFoundError(err) {
				// Surface a stable NOT_FOUND state instead of the empty string.
				// The Create waiter lists NOT_FOUND as a Pending state so a brief
				// NotFound immediately after create (the catalog is not yet
				// visible to Get due to eventual consistency) is retried until
				// the catalog appears, rather than treated as an unexpected
				// state. This is distinct from the normal Resource Read, which
				// still clears the TF id on NotFound.
				return object, "NOT_FOUND", nil
			}
			return nil, "", WrapError(err)
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
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeDlfNextCatalog >>> Encapsulated.

// isDlfNextCatalogBeingCreatedError reports whether the error returned by
// DropCatalog indicates the catalog is still being created and cannot be
// dropped yet. Catalog creation is asynchronous on the DLF backend: right
// after CreateCatalog the catalog is in NEW/INITIALIZING status and DropCatalog
// rejects the deletion with HTTP 403 "... is being created, can not be
// dropped. Please wait for the catalog to be ready.". NeedRetry does not match
// this (it is a 403, not throttling/5xx), so the Delete path must retry it
// explicitly until the catalog reaches a ready state.
func isDlfNextCatalogBeingCreatedError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "is being created")
}

// WaitForDlfNextCatalogGone polls GetCatalog until the catalog is gone
// (CatalogNotFound). DropCatalog is asynchronous: it returns 200 while the
// catalog transitions through DELETING/DELETED before finally returning 404.
// A dedicated poll (rather than DlfNextCatalogStateRefreshFunc) avoids the
// fragility of matching the transient DELETED status, which is neither a
// stable pending state nor a terminal "gone" signal in the generic
// state-refresh helper: here any non-NotFound response simply keeps polling
// until the catalog returns 404 or the timeout elapses.
func (s *DlfNextServiceV2) WaitForDlfNextCatalogGone(id string, timeout time.Duration) error {
	return s.WaitForDlfNextCatalogGoneWithApi(id, timeout, s.DescribeDlfNextCatalog)
}

// WaitForDlfNextCatalogGoneWithApi polls GetCatalog until the catalog is gone
// (CatalogNotFound). The get function is injectable so the delete-waiter
// policy can be unit-tested without a live backend. DropCatalog is
// asynchronous: it returns 200 while the catalog transitions through
// DELETING/DELETED before finally returning 404. Any non-NotFound response
// simply keeps polling until the catalog returns 404 or the timeout elapses;
// a transient retryable Get error is retried, a non-retryable Get error
// fails.
func (s *DlfNextServiceV2) WaitForDlfNextCatalogGoneWithApi(id string, timeout time.Duration, call func(id string) (map[string]interface{}, error)) error {
	wait := dlfNextCatalogBackoffWait(3*time.Second, 30*time.Second)
	return resource.Retry(timeout, func() *resource.RetryError {
		_, err := call(id)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		// Catalog still exists (DELETING/DELETED); keep polling until 404.
		return resource.RetryableError(fmt.Errorf("DlfNext Catalog %s is still being deleted", id))
	})
}

// isDlfNextCatalogAlreadyExistsError reports whether err is a CreateCatalog
// "catalog already exists" response. CreateCatalog carries no idempotency
// token, so an explicit already-exists signal means a catalog with this name
// is already present and must not be adopted by this run.
func isDlfNextCatalogAlreadyExistsError(err error) bool {
	if err == nil {
		return false
	}
	if IsExpectedErrors(err, []string{"AlreadyExist", "CatalogAlreadyExists", "DuplicateCatalog", "Duplicate"}) {
		return true
	}
	return strings.Contains(strings.ToLower(err.Error()), "already exist")
}

// isDlfNextCatalogExplicitClientError reports whether err is a client-side
// validation or permission error (HTTP 400/401/403) that proves the create
// could not have succeeded. Throttling and the Delete path's "being created"
// 403 are excluded, as is an explicit already-exists error (handled
// separately). Such errors must never trigger adoption of a same-name catalog
// found afterwards.
func isDlfNextCatalogExplicitClientError(err error) bool {
	if err == nil {
		return false
	}
	if isDlfNextCatalogAlreadyExistsError(err) {
		return false
	}
	if isDlfNextCatalogBeingCreatedError(err) {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok && e.StatusCode != nil {
		switch tea.IntValue(e.StatusCode) {
		case 400, 401, 403:
			return true
		}
	}
	return false
}

// isDlfNextCatalogBackendCapacityError reports whether err signals a backend
// cleanup-quota / environment-capacity limitation (review §6 category 2), NOT a
// transient API throttle. DLF CreateCatalog carries no idempotency token, and
// the backend rejects new catalogs when the account+Region has too many
// recently-created catalogs still being cleaned up (the backend reaps them
// asynchronously, ~48h). The signals are: an error message containing the
// backend capacity phrases, or a catalog that reached INITIALIZE_FAILED after
// the create was accepted. This is NOT retried: re-issuing CreateCatalog would
// only add another catalog to the cleanup queue and worsen the quota pressure.
// The caller must surface a clear INFRA_BLOCKED error with account/Region/
// Catalog name/RequestId/server message and let the test runner classify it as
// an environment capacity issue (not a regression, not a PASS).
func isDlfNextCatalogBackendCapacityError(err error) bool {
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

// dlfNextCatalogBackoffWait returns a closure that sleeps for an exponentially
// increasing duration with a random jitter (review §6 category 1), so multiple
// ACC runners hitting the same DLF backend do not retry in lockstep and hammer
// the API. The base grows geometrically from initial up to max, and each sleep
// is shaken by up to +/- 25% jitter. Bounded: once max is reached, subsequent
// sleeps stay at max (+jitter). This replaces the linear incrementalWait for
// DLF Next Create/Delete transient-throttle retry loops.
func dlfNextCatalogBackoffWait(initial, max time.Duration) func() {
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

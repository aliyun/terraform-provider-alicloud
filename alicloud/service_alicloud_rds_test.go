package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/stretchr/testify/assert"
)

// newTestServerError builds a synthetic SDK ServerError with the given HTTP status
// and Alibaba Cloud error code, mirroring how the real SDK deserializes an error
// response (the error code is parsed out of the JSON body's "Code" field). It
// keeps message text free of any code substring so IsExpectedErrors' message
// matching cannot produce false positives in the assertions below.
func newTestServerError(t *testing.T, httpStatus int, code, message string) *errors.ServerError {
	t.Helper()
	body := fmt.Sprintf(`{"Code":%q,"Message":%q}`, code, message)
	err := errors.NewServerError(httpStatus, body, "")
	se, ok := err.(*errors.ServerError)
	if !ok {
		t.Fatalf("NewServerError did not return *ServerError: %T", err)
	}
	return se
}

// TestRDSDBInstanceGoneStatusClassification exercises the error-classification
// logic that the recycle-bin / refunded-instance fix relies on.
//
// When a Subscription (pay-by-month/year) RDS instance is unsubscribed out of band via the
// console, it enters refund -> lock -> release -> recycle-bin. During the
// recycle-bin retention window the instance object still exists and the RDS API
// returns 403 OperationDenied.DBInstanceStatus / OperationDenied.ReadDBInstanceStatus
// — NOT a 404. The provider previously treated those 403s as retryable-then-
// hard-fail, so `terraform refresh` died and `terraform destroy` hung. The fix
// elevates these two codes to a first-class "parent instance gone" signal that is
// equal to NotFound: DescribeDBInstance maps them to NotFound, and child-resource
// Read/Delete paths confirm via a follow-up DescribeDBInstance before dropping state.
//
// This test pins the classification primitives on a BARE ServerError (before the
// choke point wraps it):
//   - IsExpectedErrors(err, dbInstanceGoneStatusCodes) matches a terminal 403 (either
//     gone code) but NotFoundError does not — which is why the fix needs the explicit
//     gone-code set plus a follow-up DescribeDBInstance instead of NotFoundError alone.
//   - A genuine 404 InvalidDBInstanceId.NotFound is NotFoundError == true and is NOT a
//     gone-code match, keeping the recycle-bin (403) and not-found (404) paths distinct.
//   - A non-terminal 403 (Throttling) matches neither, so a transient 403 is never
//     mistaken for a gone instance and remains retryable on mutation paths.
//   - isParentGone is false for every bare ServerError: it only matches a ComplexError
//     whose Err was built from NotFoundMsg (the choke point's mapped output), never a
//     raw SDK error. See TestIsParentGone for the wrapped-form behaviour.
func TestRDSDBInstanceGoneStatusClassification(t *testing.T) {
	neutralMessage := "the request was not permitted in the current state"
	tests := []struct {
		name           string
		httpStatus     int
		code           string
		wantGone       bool
		wantNotFound   bool
		wantParentGone bool
	}{
		{
			name:           "403 OperationDenied.DBInstanceStatus — refunded/recycle-bin parent",
			httpStatus:     403,
			code:           "OperationDenied.DBInstanceStatus",
			wantGone:       true,
			wantNotFound:   false,
			wantParentGone: false,
		},
		{
			name:           "403 OperationDenied.ReadDBInstanceStatus — read-side terminal lock",
			httpStatus:     403,
			code:           "OperationDenied.ReadDBInstanceStatus",
			wantGone:       true,
			wantNotFound:   false,
			wantParentGone: false,
		},
		{
			name:           "404 InvalidDBInstanceId.NotFound — real not-found",
			httpStatus:     404,
			code:           "InvalidDBInstanceId.NotFound",
			wantGone:       false,
			wantNotFound:   true,
			wantParentGone: false,
		},
		{
			name:           "403 Throttling — transient, not terminal",
			httpStatus:     403,
			code:           "Throttling",
			wantGone:       false,
			wantNotFound:   false,
			wantParentGone: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := newTestServerError(t, tt.httpStatus, tt.code, neutralMessage)
			assert.Equal(t, tt.wantGone, IsExpectedErrors(err, dbInstanceGoneStatusCodes),
				"IsExpectedErrors(err, dbInstanceGoneStatusCodes) for code %q", tt.code)
			assert.Equal(t, tt.wantNotFound, NotFoundError(err),
				"NotFoundError(err) for code %q", tt.code)
			assert.Equal(t, tt.wantParentGone, isParentGone(err),
				"isParentGone(err) for code %q (bare ServerError, not choke-point-wrapped)", tt.code)
		})
	}
}

// TestRDSDBInstanceGoneStatusCodesContents guards the gone-code set itself: the two
// codes are load-bearing across DescribeDBInstance (choke point) and the
// child-resource Read/Delete follow-up checks. Drift here silently re-opens the
// refresh-hard-fail / destroy-hang regression, so any change must be deliberate.
func TestRDSDBInstanceGoneStatusCodesContents(t *testing.T) {
	assert.Equal(t,
		[]string{"OperationDenied.DBInstanceStatus", "OperationDenied.ReadDBInstanceStatus"},
		dbInstanceGoneStatusCodes)
	assert.NotContains(t, dbInstanceGoneStatusCodes, "OperationDenied.DBStatus",
		"DBStatus is a transient per-database lock, not a parent-gone signal — it must stay in OperationDeniedDBStatus (retryable), not here")
}

// TestIsParentGone pins the isParentGone primitive that replaces the broad
// NotFoundError at every "is the parent RDS instance gone?" decision site.
//
// The DescribeDBInstance choke point (and, after this update, DescribeDBDatabase /
// DescribeRdsDatabase / DescribeRdsAccount) wrap a recycled-parent 403 and a
// genuine 404 with NotFoundMsg (Err starts with "ResourceNotfound") and wrap an
// auth failure (e.g. InvalidAccessKeyId.NotFound — also a 404) with DefaultErrorMsg.
// isParentGone walks the Cause chain and checks the Err prefix at every
// ComplexError level, so it matches the mapped gone/404 signal (even when callers
// stack another WrapErrorf layer on top, e.g. a waiter timeout message) but does
// NOT match an auth 404 reached by recursing to a raw *errors.ServerError — which
// is exactly the false-positive NotFoundError had (F1): an auth failure would
// clear state / finish delete as if the instance were gone.
func TestIsParentGone(t *testing.T) {
	neutralMessage := "the request was not permitted in the current state"
	// Simulate the choke point output for each error class.
	gone403 := WrapErrorf(newTestServerError(t, 403, "OperationDenied.DBInstanceStatus", neutralMessage), NotFoundMsg, AlibabaCloudSdkGoERROR)
	goneRead403 := WrapErrorf(newTestServerError(t, 403, "OperationDenied.ReadDBInstanceStatus", neutralMessage), NotFoundMsg, AlibabaCloudSdkGoERROR)
	real404 := WrapErrorf(newTestServerError(t, 404, "InvalidDBInstanceId.NotFound", neutralMessage), NotFoundMsg, AlibabaCloudSdkGoERROR)
	auth404 := WrapErrorf(newTestServerError(t, 404, "InvalidAccessKeyId.NotFound", neutralMessage), DefaultErrorMsg, "rm-test", "DescribeDBInstanceAttribute", AlibabaCloudSdkGoERROR)
	// A waiter stacks a timeout-message layer on top of the choke point output;
	// the gone signal must survive the extra WrapErrorf layer.
	stackedGone := WrapErrorf(gone403, DefaultTimeoutMsg, "rm-test", "WaitForDBInstance", AlibabaCloudSdkGoERROR)
	// A bare ServerError (not routed through the choke point) is not a gone signal.
	bare403 := newTestServerError(t, 403, "OperationDenied.DBInstanceStatus", neutralMessage)

	assert.True(t, isParentGone(gone403), "choke-point mapped gone 403 must be parent-gone")
	assert.True(t, isParentGone(goneRead403), "choke-point mapped read gone 403 must be parent-gone")
	assert.True(t, isParentGone(real404), "choke-point mapped real 404 must be parent-gone")
	assert.False(t, isParentGone(auth404), "auth 404 wrapped with DefaultErrorMsg must NOT be parent-gone (the F1 fix)")
	assert.True(t, isParentGone(stackedGone), "gone signal must survive an extra WrapErrorf layer (waiter)")
	assert.False(t, isParentGone(bare403), "a bare ServerError not routed through the choke point must not be parent-gone")
	assert.False(t, isParentGone(nil), "nil is not parent-gone")

	// Contrast with NotFoundError: NotFoundError accepts the auth 404 by recursing
	// the Cause chain to a 404 *errors.ServerError — the F1 bug isParentGone fixes.
	assert.True(t, NotFoundError(auth404),
		"NotFoundError accepts auth 404 via Cause recursion (the F1 bug); isParentGone does not")
}

// TestInvokerRunPreservesErrorChain pins the Invoker exhaustion chain-preservation
// fix (PR #10530) using the WRAPPED error form the real RDS callback emits.
//
// DescribeDBAccountPrivilege is the only RDS path that routes its API call through
// Invoker.Run (with a DBInstanceStatusCatcher). When the parent instance has been
// refunded / unsubscribed out of band and sits in the recycle bin, DescribeAccounts
// keeps returning 403 OperationDenied.DBInstanceStatus until the catcher exhausts
// its retry budget. The callback wraps the raw error with WrapErrorf(DefaultErrorMsg)
// (see service_alicloud_rds.go DescribeDBAccountPrivilege). The gone-check then needs
// IsExpectedErrors(exhaustedErr, dbInstanceGoneStatusCodes) to still return true so it
// can fall through to the DescribeDBInstance follow-up and map the privilege to NotFound
// — unblocking refresh/destroy on the child resource.
//
// Before the fix, Run returned fmt.Errorf("Retry timeout and got an error: %#v.",
// err): the %#v verb stringified err into a fresh *errors.errorString. For a BARE
// ServerError the code text ends up in the string and IsExpectedErrors' final
// strings.Contains fallback would still match — so a bare-callback test passes under
// the broken code too and is not a real guard. For the WRAPPED ComplexError the real
// callback emits, %#v renders the Cause as a pointer address (*errors.ServerError)(0x…)
// and the code is lost, so IsExpectedErrors never fired and the privilege refresh
// hard-failed forever on a recycled parent. This test uses the wrapped form.
func TestInvokerRunPreservesErrorChain(t *testing.T) {
	neutralMessage := "the request was not permitted in the current state"
	wrappedCallback := func() error {
		raw := newTestServerError(t, 403, "OperationDenied.DBInstanceStatus", neutralMessage)
		return WrapErrorf(raw, DefaultErrorMsg, "rm-test:acc:priv", "DescribeAccounts", AlibabaCloudSdkGoERROR)
	}
	invoker := Invoker{catchers: []*Catcher{{
		Reason:           "OperationDenied.DBInstanceStatus",
		RetryCount:       1,
		RetryWaitSeconds: 0,
	}}}
	exhausted := invoker.Run(wrappedCallback)

	assert.True(t, IsExpectedErrors(exhausted, dbInstanceGoneStatusCodes),
		"exhausted Invoker.Run error must remain gone-code matchable through the Cause chain; got: %v", exhausted)
	assert.False(t, NotFoundError(exhausted),
		"exhausted Invoker.Run error must NOT be NotFound (terminal 403, not 404); got: %v", exhausted)

	// The raw *ServerError must be reachable by walking the Cause chain so a
	// caller can still inspect ErrorCode / HttpStatus after exhaustion.
	// ComplexError has no Unwrap method, so errors.As cannot be used here.
	var se *errors.ServerError
	for cur := exhausted; cur != nil; {
		if ce, ok := cur.(*ComplexError); ok {
			if s, ok2 := ce.Cause.(*errors.ServerError); ok2 {
				se = s
				break
			}
			cur = ce.Cause
		} else {
			break
		}
	}
	if se == nil {
		t.Fatalf("raw *ServerError Cause not reachable after exhaustion; got: %v", exhausted)
	}
	assert.Equal(t, "OperationDenied.DBInstanceStatus", se.ErrorCode(),
		"the reachable raw ServerError must carry the gone code")
}

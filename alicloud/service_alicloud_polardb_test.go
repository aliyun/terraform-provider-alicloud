package alicloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnitPolarDBEndpointLookupError verifies the error classification that
// DescribePolarDBClusterEndpoint relies on once its retry loop fails to
// surface a visible endpoint. This is the targeted regression for the refresh
// drift bug: a real API error (auth/permission/service) that returns the SDK's
// pre-allocated non-nil empty response must NOT be masked as NotFound, because
// the endpoint Read handlers run d.SetId("") on NotFoundError and would clear
// state for an endpoint that still exists.
//
// The synthetic *errors.ServerError is built via newTestServerError (see
// service_alicloud_rds_test.go), which mirrors how the real SDK deserializes
// an error response (the error code is parsed out of the JSON body's "Code"
// field) and keeps the message free of any expected-code substring so
// IsExpectedErrors' message matching cannot produce false positives.
//
// The retry-recovery path (transient empty list then a visible endpoint on the
// next attempt) and the new-vs-existing call-site difference are exercised by
// the three ConfigUpdate AccTests on a live cluster; the transient
// empty-then-recover state is not deterministically reproducible in a unit
// test, so this test covers the post-retry classification only.
func TestUnitPolarDBEndpointLookupError(t *testing.T) {
	const (
		dbClusterId  = "pc-fakecluster"
		dbEndpointId = "pe-fakeendpoint"
		actionName   = "DescribeDBClusterEndpoints"
	)

	t.Run("real API error is surfaced, not masked as NotFound", func(t *testing.T) {
		// HTTP 403 (not 404) and a non-NotFound code so NotFoundError's
		// 404 / InstanceNotFound fallbacks cannot turn this into a false
		// NotFound. The message avoids any expected-code substring so
		// IsExpectedErrors' message matching stays exact.
		sdkErr := newTestServerError(t, 403, "NoPermission", "request forbidden by RAM policy")
		result := polarDBEndpointLookupError(sdkErr, dbClusterId, dbEndpointId, actionName)

		assert.False(t, NotFoundError(result), "a real API error must not be classified as NotFound")
		assert.Contains(t, result.Error(), "Failed!!!", "the error must use the default failure message, not the NotFound message")
		assert.NotContains(t, result.Error(), "ResourceNotfound!!!", "the error must not carry the NotFound wrapper")
	})

	t.Run("InvalidDBClusterId.NotFound retried to exhaustion is NotFound", func(t *testing.T) {
		// InvalidDBClusterId.NotFound is retried inside the loop; once the
		// retry window is exhausted the cluster is gone, so the endpoint is
		// treated as NotFound.
		sdkErr := newTestServerError(t, 404, "InvalidDBClusterId.NotFound", "the specified dbcluster is gone")
		result := polarDBEndpointLookupError(sdkErr, dbClusterId, dbEndpointId, actionName)

		assert.True(t, NotFoundError(result), "an exhausted InvalidDBClusterId.NotFound (cluster gone) must be NotFound")
	})

	t.Run("persistent empty list on successful calls is NotFound", func(t *testing.T) {
		// lastSDKErr == nil means the last DescribeDBClusterEndpoints call
		// succeeded but still returned an empty Items list after retrying for
		// the whole window; the endpoint has been removed.
		result := polarDBEndpointLookupError(nil, dbClusterId, dbEndpointId, actionName)

		assert.True(t, NotFoundError(result), "a persistent empty list on successful calls must be NotFound")
	})
}

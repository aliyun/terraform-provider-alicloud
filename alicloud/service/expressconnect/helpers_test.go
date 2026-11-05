package expressconnect_test

import (
	"os"
	"strings"
	"sync"
	"testing"
)

// Local replicas of root-package test helpers (alicloud_sweeper_test.go,
// provider_test.go) consumed by the migrated expressconnect unit tests. The
// external test package (expressconnect_test) cannot reach unexported root
// symbols, so the minimal set used here is re-declared with identical bodies.
// The endpoint cache is intentionally local: the tests pass it explicitly to
// connectivity.Config.Endpoints, so a per-region sync.Map behaves identically
// to the root implementation.

var regionalEndpoints sync.Map

// sharedEndpointCacheForRegion returns the product endpoint cache for a region.
func sharedEndpointCacheForRegion(region string) *sync.Map {
	cache, _ := regionalEndpoints.LoadOrStore(region, new(sync.Map))
	return cache.(*sync.Map)
}

// testAccPreCheckWithExpressConnectUidSetting skips the test unless the
// ALICLOUD_EXPRESS_CONNECT_UID environment variable is set.
func testAccPreCheckWithExpressConnectUidSetting(t *testing.T) {
	if v := strings.TrimSpace(os.Getenv("ALICLOUD_EXPRESS_CONNECT_UID")); v == "" {
		t.Skipf("Skipping the test case with no express connect uid  id setting")
		t.Skipped()
	}
}

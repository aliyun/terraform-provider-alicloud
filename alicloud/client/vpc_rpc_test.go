package client

// Wire regressions for the VPC calls migrated off WithVpcClient onto
// VpcSdkRpc: every migrated call site keeps the old HMAC-SHA1 signing unless
// sign_version.default = v4 selects the common V4 gateway, and the typed SDK
// response the callers assert on survives both branches. The requests are
// captured on a local HTTP server with temporary STS credentials, mirroring
// the ECS, CMS and SLS regressions in legacy_rpc_test.go.

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/client/config"
)

// TestUnitVpcSdkRpcWire covers the DescribeVSwitchAttributes migration, the
// most widely shared VPC call in the migration (fourteen cross-product
// lookups): unconfigured sign_version keeps the legacy signature with the
// temporary STS credentials on the query, a default of v4 signs through the
// EO V4 gateway, and a vpc = v1 override keeps the legacy signature under a
// v4 default. Both branches carry the typed serialization (VSwitchId) and
// the region, and the returned value rehydrates into the typed SDK response.
func TestUnitVpcSdkRpcWire(t *testing.T) {
	var endpoint string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		legacy := captureLegacyRpcWire(t, r, endpoint)
		if r.URL.Path != "/" {
			t.Errorf("expected the RPC root path, got %q", r.URL.Path)
		}
		q := r.URL.Query()
		if legacy {
			if q.Get("Action") != "DescribeVSwitchAttributes" || q.Get("Version") != "2016-04-28" {
				t.Error("expected the DescribeVSwitchAttributes action and version on the query")
			}
		} else {
			verifyWireV4(t, r, "vpc", legacyWireRegion)
			if r.Header.Get("x-acs-action") != "DescribeVSwitchAttributes" || r.Header.Get("x-acs-version") != "2016-04-28" {
				t.Error("expected the DescribeVSwitchAttributes action and version headers")
			}
		}
		if q.Get("RegionId") != legacyWireRegion || q.Get("VSwitchId") != "vsw-unit-test" {
			t.Errorf("expected the region and VSwitchId on the query, got %q / %q", q.Get("RegionId"), q.Get("VSwitchId"))
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"VSwitchId":"vsw-unit-test","VpcId":"vpc-unit-test","CidrBlock":"192.168.0.0/16","RequestId":"request-id"}`)
	}))
	defer server.Close()
	endpoint = strings.TrimPrefix(server.URL, "http://")

	for _, tc := range []struct {
		name        string
		defaultSign string
		productSign string
	}{
		{name: "unconfigured-keeps-legacy-signature"},
		{name: "default-v4-signs-eo-v4", defaultSign: "v4"},
		{name: "vpc-v1-overrides-default-v4", defaultSign: "v4", productSign: "v1"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := newLegacyWireClient(t, "vpc", endpoint)
			if tc.defaultSign != "" {
				c.config.SignVersion.Store("default", tc.defaultSign)
			}
			if tc.productSign != "" {
				c.config.SignVersion.Store("vpc", tc.productSign)
			}
			request := vpc.CreateDescribeVSwitchAttributesRequest()
			request.RegionId = legacyWireRegion
			request.VSwitchId = "vsw-unit-test"
			raw, err := c.VpcSdkRpc(request, vpc.CreateDescribeVSwitchAttributesResponse())
			if err != nil {
				t.Fatalf("DescribeVSwitchAttributes: %s", err)
			}
			response, ok := raw.(*vpc.DescribeVSwitchAttributesResponse)
			if !ok || response.VSwitchId != "vsw-unit-test" || response.VpcId != "vpc-unit-test" {
				t.Fatalf("typed response lost after the shared path: %T %+v", raw, response)
			}
		})
	}
}

// TestUnitVpcSdkRpcNestedResponseWire pins the response-side contract of the
// list-shaped VPC calls (route tables, NAT gateways, router interfaces,
// SSL VPN servers): the V4 branch must rehydrate the shared map into the
// nested typed lists the callers iterate, not only the flat envelope.
func TestUnitVpcSdkRpcNestedResponseWire(t *testing.T) {
	var endpoint string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !captureLegacyRpcWire(t, r, endpoint) {
			verifyWireV4(t, r, "vpc", legacyWireRegion)
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"RequestId":"request-id","TotalCount":1,`+
			`"RouteTables":{"RouteTable":[{"RouteTableId":"vtb-unit-test","RouteTableName":"tf-test","VpcId":"vpc-unit-test"}]}}`)
	}))
	defer server.Close()
	endpoint = strings.TrimPrefix(server.URL, "http://")

	for _, tc := range []struct {
		name        string
		defaultSign string
	}{
		{name: "unconfigured-keeps-legacy-signature"},
		{name: "default-v4-signs-eo-v4", defaultSign: "v4"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := newLegacyWireClient(t, "vpc", endpoint)
			if tc.defaultSign != "" {
				c.config.SignVersion.Store("default", tc.defaultSign)
			}
			request := vpc.CreateDescribeRouteTablesRequest()
			request.RegionId = legacyWireRegion
			request.RouteTableId = "vtb-unit-test"
			raw, err := c.VpcSdkRpc(request, vpc.CreateDescribeRouteTablesResponse())
			if err != nil {
				t.Fatalf("DescribeRouteTables: %s", err)
			}
			response, ok := raw.(*vpc.DescribeRouteTablesResponse)
			if !ok {
				t.Fatalf("typed response lost after the shared path: %T", raw)
			}
			if len(response.RouteTables.RouteTable) != 1 || response.RouteTables.RouteTable[0].RouteTableId != "vtb-unit-test" {
				t.Fatalf("nested route table list lost after the shared path: %+v", response)
			}
		})
	}
}

func TestUnitWithVpcClientLegacyFallback(t *testing.T) {
	c := newLegacyWireClient(t, "vpc", "127.0.0.1:1")
	c.config.SignVersion.Store("default", "v4")
	called := false
	_, err := c.WithVpcClient(func(conn *vpc.Client) (interface{}, error) { called = true; return nil, nil })
	if err != nil || !called {
		t.Fatalf("fallback called=%v, err=%v", called, err)
	}
}

// TestUnitVpcSdkRpcCrossRegionWire pins the cross-region contract the legacy
// SDK kept and the V4 branch must not regress: a request carrying its own
// region — the router interface connection describing its opposite side —
// selects that region for the V4 signing key and the query, not the
// provider's client region.
func TestUnitVpcSdkRpcCrossRegionWire(t *testing.T) {
	const oppositeRegion = "cn-opposite-unit-test"
	var endpoint string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if captureLegacyRpcWire(t, r, endpoint) {
			t.Error("expected the request to be V4 signed")
		}
		verifyWireV4(t, r, "vpc", oppositeRegion)
		if q := r.URL.Query(); q.Get("RegionId") != oppositeRegion {
			t.Errorf("expected the opposite region on the query, got %q", q.Get("RegionId"))
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"RequestId":"request-id","TotalCount":1,`+
			`"RouterInterfaceSet":{"RouterInterfaceType":[{"RouterInterfaceId":"ri-unit-test","OppositeRegionId":"`+oppositeRegion+`"}]}}`)
	}))
	defer server.Close()
	endpoint = strings.TrimPrefix(server.URL, "http://")

	c := newLegacyWireClient(t, "vpc", endpoint)
	c.config.SignVersion.Store("default", "v4")
	// An explicitly configured EO endpoint takes priority over a cached location.
	c.config.VpcEndpoint = endpoint
	c.config.Endpoints.Store(config.RegionalEndpointCacheKey("vpc", oppositeRegion), "stale.invalid")
	request := vpc.CreateDescribeRouterInterfacesRequest()
	request.RegionId = oppositeRegion
	values := []string{"ri-unit-test"}
	request.Filter = &[]vpc.DescribeRouterInterfacesFilter{{Key: "RouterInterfaceId", Value: &values}}
	raw, err := c.VpcSdkRpc(request, vpc.CreateDescribeRouterInterfacesResponse())
	if err != nil {
		t.Fatalf("DescribeRouterInterfaces: %s", err)
	}
	response, ok := raw.(*vpc.DescribeRouterInterfacesResponse)
	if !ok || len(response.RouterInterfaceSet.RouterInterfaceType) != 1 ||
		response.RouterInterfaceSet.RouterInterfaceType[0].RouterInterfaceId != "ri-unit-test" {
		t.Fatalf("typed response lost after the shared path: %T %+v", raw, response)
	}
}

// TestUnitVpcSdkRpcCrossRegionAfterSameRegionWarmupWire pins the ordering a
// real apply produces: the same-region VPC calls run first and warm the
// shared per-product endpoint entry, then the cross-region
// DescribeRouterInterfaces must still reach its own region-scoped endpoint
// host and sign with the request's region. Reading the warm shared entry
// instead would send the opposite region's request to the client region's
// host — the host of one region carrying the signature of another.
func TestUnitVpcSdkRpcCrossRegionAfterSameRegionWarmupWire(t *testing.T) {
	const oppositeRegion = "cn-opposite-unit-test"
	var sameCount, crossCount int32

	sameRegionServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isV4WireRequest(r) {
			t.Error("the warm-up call must be V4 signed under the default v4 selection")
		} else {
			verifyWireV4(t, r, "vpc", legacyWireRegion)
		}
		if q := r.URL.Query(); q.Get("RegionId") != legacyWireRegion {
			t.Errorf("expected the client region on the warm-up query, got %q", q.Get("RegionId"))
		}
		atomic.AddInt32(&sameCount, 1)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"VSwitchId":"vsw-unit-test","VpcId":"vpc-unit-test","RequestId":"request-id"}`)
	}))
	defer sameRegionServer.Close()
	sameEndpoint := strings.TrimPrefix(sameRegionServer.URL, "http://")

	var crossRegionServerHost string
	crossRegionServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isV4WireRequest(r) {
			t.Error("expected the request to be V4 signed")
		} else {
			verifyWireV4(t, r, "vpc", oppositeRegion)
		}
		if r.Host != crossRegionServerHost {
			t.Errorf("expected the request on the cross-region endpoint host, got %q", r.Host)
		}
		if q := r.URL.Query(); q.Get("RegionId") != oppositeRegion {
			t.Errorf("expected the opposite region on the query, got %q", q.Get("RegionId"))
		}
		atomic.AddInt32(&crossCount, 1)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"RequestId":"request-id","TotalCount":1,`+
			`"RouterInterfaceSet":{"RouterInterfaceType":[{"RouterInterfaceId":"ri-unit-test","OppositeRegionId":"`+oppositeRegion+`"}]}}`)
	}))
	defer crossRegionServer.Close()
	crossRegionServerHost = strings.TrimPrefix(crossRegionServer.URL, "http://")

	c := newLegacyWireClient(t, "vpc", sameEndpoint)
	c.config.SignVersion.Store("default", "v4")

	// The apply's same-region calls warm the shared per-product entry.
	warmup := vpc.CreateDescribeVSwitchAttributesRequest()
	warmup.RegionId = legacyWireRegion
	warmup.VSwitchId = "vsw-unit-test"
	if _, err := c.VpcSdkRpc(warmup, vpc.CreateDescribeVSwitchAttributesResponse()); err != nil {
		t.Fatalf("DescribeVSwitchAttributes warm-up: %s", err)
	}

	// Location described the opposite region's endpoint for the product.
	c.config.Endpoints.Store(config.RegionalEndpointCacheKey("vpc", oppositeRegion), crossRegionServerHost)

	request := vpc.CreateDescribeRouterInterfacesRequest()
	request.RegionId = oppositeRegion
	values := []string{"ri-unit-test"}
	request.Filter = &[]vpc.DescribeRouterInterfacesFilter{{Key: "RouterInterfaceId", Value: &values}}
	raw, err := c.VpcSdkRpc(request, vpc.CreateDescribeRouterInterfacesResponse())
	if err != nil {
		t.Fatalf("DescribeRouterInterfaces: %s", err)
	}
	response, ok := raw.(*vpc.DescribeRouterInterfacesResponse)
	if !ok || len(response.RouterInterfaceSet.RouterInterfaceType) != 1 ||
		response.RouterInterfaceSet.RouterInterfaceType[0].RouterInterfaceId != "ri-unit-test" {
		t.Fatalf("typed response lost after the shared path: %T %+v", raw, response)
	}
	if got := atomic.LoadInt32(&sameCount); got != 1 {
		t.Errorf("expected exactly one same-region request, got %d", got)
	}
	if got := atomic.LoadInt32(&crossCount); got != 1 {
		t.Errorf("expected the cross-region request on its own endpoint host, got %d", got)
	}
}

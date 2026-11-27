package client

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	sdkerrors "github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

func TestUnitEcsSdkRpcLegacyRetries(t *testing.T) {
	var count atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count.Add(1)
		w.WriteHeader(500)
		io.WriteString(w, `{"Code":"InternalError","Message":"fixture"}`)
	}))
	defer server.Close()
	c := newLegacyWireClient(t, "ecs", server.Listener.Addr().String())
	response := ecs.CreateDescribeRegionsResponse()
	got, err := c.EcsSdkRpc(ecs.CreateDescribeRegionsRequest(), response)
	var remote *sdkerrors.ServerError
	if got != response || !errors.As(err, &remote) || count.Load() != 6 {
		t.Fatalf("attempts=%d, response=%T, err=%v", count.Load(), got, err)
	}
}

func TestUnitEcsSdkRpcWire(t *testing.T) {
	for _, v4 := range []bool{false, true} {
		t.Run(fmt.Sprintf("v4=%t", v4), func(t *testing.T) {
			var endpoint string
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				legacy := captureLegacyRpcWire(t, r, endpoint)
				if legacy == v4 {
					t.Errorf("unexpected signature for v4=%t", v4)
				}
				action := r.URL.Query().Get("Action")
				if v4 {
					verifyWireV4(t, r, "ecs", legacyWireRegion)
					action = r.Header.Get("x-acs-action")
				}
				w.Header().Set("Content-Type", "application/json")
				switch action {
				case "DescribeZones":
					io.WriteString(w, `{"Zones":{"Zone":[{"ZoneId":"zone-one"}]}}`)
				case "ModifySecurityGroupEgressRule":
					if r.URL.Query().Get("SecurityGroupId") != "sg-wire" {
						t.Error("missing group ID")
					}
					calls.Add(1)
					io.WriteString(w, `{"RequestId":"updated"}`)
				default:
					t.Errorf("unexpected action %s", action)
				}
			}))
			defer server.Close()
			endpoint = server.Listener.Addr().String()
			c := newLegacyWireClient(t, "ecs", endpoint)
			if v4 {
				c.config.SignVersion.Store("default", "v4")
			}
			resp := ecs.CreateDescribeZonesResponse()
			_, err := c.EcsSdkRpc(ecs.CreateDescribeZonesRequest(), resp)
			if err != nil || len(resp.Zones.Zone) != 1 || resp.Zones.Zone[0].ZoneId != "zone-one" {
				t.Fatalf("zones=%v, err=%v", resp.Zones, err)
			}
			req := requests.NewCommonRequest()
			req.Product = "Ecs"
			req.Version = "2014-05-26"
			req.ApiName = "ModifySecurityGroupEgressRule"
			req.Method = "POST"
			req.Scheme = "HTTP"
			req.QueryParams["SecurityGroupId"] = "sg-wire"
			if _, err = c.EcsCommonRpc(req); err != nil || calls.Load() != 1 {
				t.Fatalf("calls=%d, err=%v", calls.Load(), err)
			}
		})
	}
}

package esa

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestEsaString(t *testing.T) {
	cases := []struct {
		item map[string]interface{}
		key  string
		want string
	}{
		{map[string]interface{}{"Id": "9987"}, "Id", "9987"},
		{map[string]interface{}{"Id": float64(9987)}, "Id", "9987"},
		{map[string]interface{}{"Id": nil}, "Id", ""},
		{map[string]interface{}{}, "Missing", ""},
		{map[string]interface{}{"Flag": true}, "Flag", "true"},
	}
	for _, c := range cases {
		if got := esaString(c.item, c.key); got != c.want {
			t.Errorf("esaString(%v, %q) = %q, want %q", c.item, c.key, got, c.want)
		}
	}
}

func TestLoadBalancerOriginStatusItems(t *testing.T) {
	t.Run("parses origin status array", func(t *testing.T) {
		resp := map[string]interface{}{
			"OriginStatus": []interface{}{
				map[string]interface{}{
					"LoadBalancerId": float64(99874066052),
					"OriginId":       float64(99750209487),
					"PoolId":         float64(99750209487),
					"PoolType":       "default_pool",
					"Reason":         "",
					"Status":         "healthy",
				},
				map[string]interface{}{
					"LoadBalancerId": "99874066053",
					"OriginId":       "99750209488",
					"PoolId":         "",
					"PoolType":       "",
					"Reason":         "TCP connection error",
					"Status":         "unhealthy",
				},
			},
			"RequestId": "0AEDAF20-4DDF-4165-8750-47FF9C1929C9",
		}
		items := loadBalancerOriginStatusItems(resp)
		if len(items) != 2 {
			t.Fatalf("expected 2 items, got %d", len(items))
		}
		first := items[0]
		assertStringEq(t, first.LoadBalancerId, "99874066052")
		assertStringEq(t, first.OriginId, "99750209487")
		assertStringEq(t, first.PoolType, "default_pool")
		assertStringEq(t, first.Status, "healthy")
		assertStringEq(t, items[1].Reason, "TCP connection error")
		assertStringEq(t, items[1].Status, "unhealthy")
	})

	t.Run("empty response yields empty slice", func(t *testing.T) {
		if got := loadBalancerOriginStatusItems(nil); len(got) != 0 {
			t.Errorf("expected empty slice, got %v", got)
		}
		if got := loadBalancerOriginStatusItems(map[string]interface{}{}); len(got) != 0 {
			t.Errorf("expected empty slice, got %v", got)
		}
		if got := loadBalancerOriginStatusItems(map[string]interface{}{"OriginStatus": "not-an-array"}); len(got) != 0 {
			t.Errorf("expected empty slice for non-array, got %v", got)
		}
	})

	t.Run("ignores non-object entries", func(t *testing.T) {
		resp := map[string]interface{}{
			"OriginStatus": []interface{}{"bad", map[string]interface{}{"Status": "unknown"}},
		}
		items := loadBalancerOriginStatusItems(resp)
		if len(items) != 1 || items[0].Status.ValueString() != "unknown" {
			t.Fatalf("expected 1 valid item, got %v", items)
		}
	})
}

func TestLoadBalancerRegionItems(t *testing.T) {
	t.Run("parses regions and sub-regions", func(t *testing.T) {
		resp := map[string]interface{}{
			"Regions": []interface{}{
				map[string]interface{}{
					"RegionCode":   "SEAS",
					"RegionCnName": "South East Asia",
					"RegionEnName": "South East Asia",
					"SubRegions": []interface{}{
						map[string]interface{}{
							"SubRegionCode":   "ID",
							"SubRegionCnName": "Indonesia",
							"SubRegionEnName": "Indonesia",
						},
					},
				},
			},
			"TotalCount": float64(12),
		}
		items := loadBalancerRegionItems(resp)
		if len(items) != 1 {
			t.Fatalf("expected 1 region, got %d", len(items))
		}
		r := items[0]
		assertStringEq(t, r.RegionCode, "SEAS")
		assertStringEq(t, r.RegionEnName, "South East Asia")
		if len(r.SubRegions) != 1 {
			t.Fatalf("expected 1 sub-region, got %d", len(r.SubRegions))
		}
		assertStringEq(t, r.SubRegions[0].SubRegionCode, "ID")
		assertStringEq(t, r.SubRegions[0].SubRegionEnName, "Indonesia")
	})

	t.Run("region without sub-regions", func(t *testing.T) {
		resp := map[string]interface{}{
			"Regions": []interface{}{
				map[string]interface{}{
					"RegionCode":   "CNM",
					"RegionCnName": "Chinese Mainland",
					"RegionEnName": "Chinese Mainland",
				},
			},
		}
		items := loadBalancerRegionItems(resp)
		if len(items) != 1 || len(items[0].SubRegions) != 0 {
			t.Fatalf("expected empty sub-regions, got %v", items)
		}
	})

	t.Run("empty response yields empty slice", func(t *testing.T) {
		for _, resp := range []map[string]interface{}{nil, {}, {"Regions": nil}} {
			if got := loadBalancerRegionItems(resp); len(got) != 0 {
				t.Errorf("expected empty slice for %v, got %v", resp, got)
			}
		}
	})
}

func assertStringEq(t *testing.T, got types.String, want string) {
	t.Helper()
	if !reflect.DeepEqual(got, types.StringValue(want)) {
		t.Errorf("got %v, want %q", got, want)
	}
}

package eci

import (
	"reflect"
	"testing"
)

func TestUnitBuildEciAcrRegistryInfoRequest(t *testing.T) {
	cases := []struct {
		name string
		in   []interface{}
		want []map[string]interface{}
	}{
		{
			name: "cross-account roles are serialized",
			in: []interface{}{map[string]interface{}{
				"instance_name": "acr",
				"instance_id":   "cri-xxx",
				"domains":       []interface{}{"registry.example.com"},
				"region_id":     "cn-hangzhou",
				"arn_service":   "acs:ram::111:role/eci-side",
				"arn_user":      "acs:ram::222:role/acr-side",
			}},
			want: []map[string]interface{}{{
				"InstanceName": "acr",
				"InstanceId":   "cri-xxx",
				"Domain":       []interface{}{"registry.example.com"},
				"RegionId":     "cn-hangzhou",
				"ArnService":   "acs:ram::111:role/eci-side",
				"ArnUser":      "acs:ram::222:role/acr-side",
			}},
		},
		{
			name: "same-account arn_service only",
			in: []interface{}{map[string]interface{}{
				"instance_name": "",
				"instance_id":   "cri-xxx",
				"domains":       []interface{}{},
				"region_id":     "cn-hangzhou",
				"arn_service":   "acs:ram::111:role/eci-side",
				"arn_user":      "",
			}},
			want: []map[string]interface{}{{
				"InstanceName": "",
				"InstanceId":   "cri-xxx",
				"Domain":       []interface{}{},
				"RegionId":     "cn-hangzhou",
				"ArnService":   "acs:ram::111:role/eci-side",
			}},
		},
		{
			name: "ARNs are omitted when not configured",
			in: []interface{}{map[string]interface{}{
				"instance_name": "acr",
				"instance_id":   "cri-xxx",
				"domains":       []interface{}{},
				"region_id":     "cn-hangzhou",
				"arn_service":   "",
				"arn_user":      "",
			}},
			want: []map[string]interface{}{{
				"InstanceName": "acr",
				"InstanceId":   "cri-xxx",
				"Domain":       []interface{}{},
				"RegionId":     "cn-hangzhou",
			}},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := buildEciAcrRegistryInfoRequest(c.in); !reflect.DeepEqual(got, c.want) {
				t.Fatalf("got %#v, want %#v", got, c.want)
			}
		})
	}
}

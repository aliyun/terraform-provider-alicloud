package privatezone

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestUnitInboundEndpointIpConfigRequest stays in the internal package because the request
// encoder is unexported; the ACC file uses package privatezone_test and cannot reach it.
func TestUnitInboundEndpointIpConfigRequest(t *testing.T) {
	ipConfigs := []inboundEndpointIpConfigModel{
		{
			ZoneId:    types.StringValue("cn-hangzhou-a"),
			CidrBlock: types.StringValue("172.16.0.0/20"),
			VSwitchId: types.StringValue("vsw-explicit"),
			Ip:        types.StringValue("172.16.0.10"),
		},
		{
			ZoneId:    types.StringValue("cn-hangzhou-b"),
			CidrBlock: types.StringValue("172.16.16.0/20"),
			VSwitchId: types.StringValue("vsw-null"),
			Ip:        types.StringNull(),
		},
		{
			ZoneId:    types.StringValue("cn-hangzhou-c"),
			CidrBlock: types.StringValue("172.16.32.0/20"),
			VSwitchId: types.StringValue("vsw-unknown"),
			Ip:        types.StringUnknown(),
		},
		{
			ZoneId:    types.StringValue("cn-hangzhou-d"),
			CidrBlock: types.StringValue("172.16.48.0/20"),
			VSwitchId: types.StringValue("vsw-empty"),
			Ip:        types.StringValue(""),
		},
	}

	request := inboundEndpointIpConfigRequest(ipConfigs)
	if got := request["IpConfig.1.Ip"]; got != "172.16.0.10" {
		t.Fatalf("explicit ip must be sent, got %#v", got)
	}
	if got := request["IpConfig.1.VSwitchId"]; got != "vsw-explicit" {
		t.Fatalf("vswitch id must be sent, got %#v", got)
	}
	if got := request["IpConfig.1.AzId"]; got != "cn-hangzhou-a" {
		t.Fatalf("zone id must be sent as AzId, got %#v", got)
	}
	if got := request["IpConfig.1.CidrBlock"]; got != "172.16.0.0/20" {
		t.Fatalf("cidr block must be sent, got %#v", got)
	}
	for _, key := range []string{"IpConfig.2.Ip", "IpConfig.3.Ip", "IpConfig.4.Ip"} {
		if _, ok := request[key]; ok {
			t.Fatalf("%s must be omitted for null, unknown and empty ips", key)
		}
	}
	if got := request["IpConfig.2.VSwitchId"]; got != "vsw-null" {
		t.Fatalf("vswitch id must still be sent for an auto-assigned ip, got %#v", got)
	}
}

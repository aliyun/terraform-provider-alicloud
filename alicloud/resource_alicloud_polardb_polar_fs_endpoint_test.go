package alicloud

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestUnitFlattenPolarDBPolarFsEndpointAddresses(t *testing.T) {
	raw := []interface{}{map[string]interface{}{
		"ConnectionString": "pfs.example.com", "PrivateZoneConnectionString": "pfs.internal",
		"IPAddress": "10.0.0.1", "Port": "443", "VPCId": "vpc-1", "VSwitchId": "vsw-1", "NetType": "Private",
	}}
	want := []map[string]interface{}{{
		"connection_string": "pfs.example.com", "private_zone_connection_string": "pfs.internal",
		"ip_address": "10.0.0.1", "port": "443", "vpc_id": "vpc-1", "vswitch_id": "vsw-1", "net_type": "Private",
	}}
	if got := flattenPolarDBPolarFsEndpointAddresses(raw); !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected endpoint addresses: %#v", got)
	}
}

func TestAccAliCloudPolarDBPolarFsEndpoint_planOnly(t *testing.T) {
	resourceId := "alicloud_polardb_polar_fs_endpoint.default"
	testAccConfig := resourceTestAccConfigFunc(resourceId, "polarfs-endpoint", func(string) string { return "" })
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{{
			PlanOnly: true,
			Config: testAccConfig(map[string]interface{}{
				"db_cluster_id":           "pc-example",
				"polar_fs_instance_id":    "pfs-example",
				"endpoint_type":           "S3Gateway",
				"vpc_id":                  "vpc-example",
				"vswitch_id":              "vsw-example",
				"db_endpoint_description": "terraform-test",
			}),
			ExpectNonEmptyPlan: true,
		}},
	})
}

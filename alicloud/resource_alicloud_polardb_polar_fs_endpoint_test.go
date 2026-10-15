package alicloud

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestUnitPolarDBPolarFsEndpointTypeValidation(t *testing.T) {
	validate := resourceAlicloudPolarDBPolarFsEndpoint().Schema["endpoint_type"].ValidateFunc
	for _, value := range []string{"Nas", "S3Gateway"} {
		if _, errs := validate(value, "endpoint_type"); len(errs) != 0 {
			t.Fatalf("endpoint_type %q must be accepted: %v", value, errs)
		}
	}
	if _, errs := validate("S3Pvtz", "endpoint_type"); len(errs) == 0 {
		t.Fatal("S3Pvtz is not a documented endpoint type")
	}
}

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

func TestAccAliCloudPolarDBPolarFsEndpoint_complete(t *testing.T) {
	resourceId := "alicloud_polardb_polar_fs_endpoint.default"
	testAccConfig := resourceTestAccConfigFunc(resourceId, "default", testAccPolarDBPolarFsEndpointDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Beijing})
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolarDBPolarFsEndpointDestroy,
		Steps: []resource.TestStep{{
			Config: testAccConfig(map[string]interface{}{
				"db_cluster_id":           "${alicloud_polardb_cluster.polarfs.id}",
				"polar_fs_instance_id":    "${alicloud_polardb_polar_fs.default.id}",
				"endpoint_type":           "S3Gateway",
				"vpc_id":                  "${data.alicloud_vswitches.polarfs.vswitches[0].vpc_id}",
				"vswitch_id":              "${data.alicloud_vswitches.polarfs.vswitches[0].id}",
				"db_endpoint_description": "terraform-acc",
			}),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(resourceId, "endpoint_type", "S3Gateway"),
				resource.TestCheckResourceAttrSet(resourceId, "db_endpoint_id"),
			),
		}},
	})
}

func testAccCheckPolarDBPolarFsEndpointDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_polardb_polar_fs_endpoint" {
			continue
		}
		parts, err := ParseResourceId(rs.Primary.ID, 2)
		if err != nil {
			return err
		}
		endpoints, err := describePolarDBPolarFsEndpoints(client, parts[0], rs.Primary.Attributes["polar_fs_instance_id"], parts[1])
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		if len(endpoints) != 0 {
			return fmt.Errorf("PolarFS endpoint %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccPolarDBPolarFsEndpointDependence(string) string {
	return `
data "alicloud_vswitches" "polarfs" {
  status  = "Available"
  zone_id = "cn-beijing-l"
}

resource "alicloud_polardb_cluster" "polarfs" {
  db_type             = "PostgreSQL"
  db_version          = "14"
  pay_type            = "PostPaid"
  db_node_class       = "polar.pg.g1.tiny.c"
  db_node_num         = 1
  vswitch_id          = data.alicloud_vswitches.polarfs.vswitches[0].id
  zone_id             = data.alicloud_vswitches.polarfs.vswitches[0].zone_id
  creation_category   = "SENormal"
  storage_type        = "ESSDAUTOPL"
  storage_space       = 40
  storage_pay_type    = "PostPaid"
  hot_standby_cluster = "OFF"
  encrypt_new_tables  = "OFF"
  description         = "tf-testAccPolarDBPolarFsEndpoint"
}

resource "alicloud_polardb_polar_fs" "default" {
  db_type                = "PostgreSQL"
  creation_category      = "basic"
  storage_type           = "city_redundancy"
  storage_space          = 50
  pay_type               = "Postpaid"
  accelerate_switch      = "ON"
  accelerate_storage_size = 5000
  accelerate_type        = "juice"
  db_cluster_id          = alicloud_polardb_cluster.polarfs.id
  vpc_id                 = data.alicloud_vswitches.polarfs.vswitches[0].vpc_id
  vswitch_id             = data.alicloud_vswitches.polarfs.vswitches[0].id
  zone_id                = data.alicloud_vswitches.polarfs.vswitches[0].zone_id
}
`
}

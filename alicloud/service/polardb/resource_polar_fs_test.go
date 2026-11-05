package polardb_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	polardb "github.com/aliyun/terraform-provider-alicloud/alicloud/service/polardb"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudPolarDBPolarFs_complete(t *testing.T) {
	resourceId := "alicloud_polardb_polar_fs.default"
	endpointId := "alicloud_polardb_polar_fs_endpoint.default"
	instancesId := "data.alicloud_polardb_polar_fs_instances.default"
	tokenId := "data.alicloud_polardb_lakebase_tenant_token.default"
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, "default", testAccPolarDBPolarFsDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Beijing})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPolarDBPolarFsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":                 "PostgreSQL",
					"creation_category":       "basic",
					"storage_type":            "city_redundancy",
					"storage_space":           50,
					"pay_type":                "Postpaid",
					"accelerate_switch":       "ON",
					"accelerate_storage_size": 5000,
					"accelerate_type":         "juice",
					"db_cluster_id":           "${alicloud_polardb_cluster.polarfs.id}",
					"vpc_id":                  "${data.alicloud_vswitches.polarfs.vswitches[0].vpc_id}",
					"vswitch_id":              "${data.alicloud_vswitches.polarfs.vswitches[0].id}",
					"zone_id":                 "${data.alicloud_vswitches.polarfs.vswitches[0].zone_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "db_type", "PostgreSQL"),
					resource.TestCheckResourceAttr(resourceId, "storage_type", "city_redundancy"),
					resource.TestCheckResourceAttr(resourceId, "storage_space", "50"),
					resource.TestCheckResourceAttr(resourceId, "pay_type", "Postpaid"),
					resource.TestCheckResourceAttr(resourceId, "creation_category", "basic"),
					resource.TestCheckResourceAttr(resourceId, "accelerate_switch", "ON"),
					resource.TestCheckResourceAttr(resourceId, "accelerate_storage_size", "5000"),
					resource.TestCheckResourceAttr(resourceId, "accelerate_type", "juice"),
					resource.TestCheckResourceAttrSet(resourceId, "status"),
					resource.TestCheckResourceAttrSet(resourceId, "meta_instance_name"),
					resource.TestCheckResourceAttr(endpointId, "endpoint_type", "S3Gateway"),
					resource.TestCheckResourceAttrSet(endpointId, "db_endpoint_id"),
					resource.TestCheckResourceAttr(instancesId, "polar_fs_instances.#", "1"),
					resource.TestCheckResourceAttrPair(instancesId, "polar_fs_instances.0.id", resourceId, "id"),
					resource.TestCheckResourceAttrSet(tokenId, "token"),
					resource.TestCheckResourceAttr(tokenId, "subdir", "/terraform-acc"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"polar_fs_path", "storage_used", "accelerate_switch", "accelerate_storage_size"},
			},
		},
	})
}
func testAccPolarDBPolarFsDependence(string) string {
	return `
data "alicloud_vswitches" "polarfs" {
  status  = "Available"
  zone_id = "cn-beijing-l"
}

resource "alicloud_polardb_cluster" "polarfs" {
  db_type           = "PostgreSQL"
  db_version        = "14"
  pay_type          = "PostPaid"
  db_node_class     = "polar.pg.g1.tiny.c"
  db_node_num       = 1
  vswitch_id        = data.alicloud_vswitches.polarfs.vswitches[0].id
  zone_id           = data.alicloud_vswitches.polarfs.vswitches[0].zone_id
  creation_category = "SENormal"
  storage_type      = "ESSDAUTOPL"
  storage_space     = 40
  storage_pay_type  = "PostPaid"
  hot_standby_cluster = "OFF"
  encrypt_new_tables = "OFF"
  description       = "tf-testAccPolarDBPolarFs"
}

resource "alicloud_polardb_polar_fs_endpoint" "default" {
  db_cluster_id           = alicloud_polardb_cluster.polarfs.id
  polar_fs_instance_id    = alicloud_polardb_polar_fs.default.id
  endpoint_type           = "S3Gateway"
  vpc_id                  = data.alicloud_vswitches.polarfs.vswitches[0].vpc_id
  vswitch_id              = data.alicloud_vswitches.polarfs.vswitches[0].id
  db_endpoint_description = "terraform-acc"
}

data "alicloud_polardb_polar_fs_instances" "default" {
  ids = [alicloud_polardb_polar_fs.default.id]
}

data "alicloud_polardb_lakebase_tenant_token" "default" {
  polar_fs_instance_id = alicloud_polardb_polar_fs.default.id
  db_cluster_id        = alicloud_polardb_cluster.polarfs.id
  subdir               = "/terraform-acc"
}
`
}
func testAccCheckPolarDBPolarFsDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_polardb_polar_fs" {
			continue
		}
		_, err := polardb.DescribePolarDBPolarFsAttribute(client, rs.Primary.ID, rs.Primary.Attributes["db_cluster_id"])
		if err == nil {
			return fmt.Errorf("PolarFS instance %s still exists", rs.Primary.ID)
		}
		if !errs.NotFoundError(err) {
			return err
		}
	}
	return nil
}

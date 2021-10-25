package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudGPDBAccount_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_account.default"
	ra := resourceAttrInit(resourceId, AlicloudGPDBAccountMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tftest%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGPDBAccountBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":      "${alicloud_gpdb_elastic_instance.default.id}",
					"account_name":        name,
					"account_password":    "TFTest123",
					"account_description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_name":        name,
						"account_description": name,
						"db_instance_id":      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "TFTest123" + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password"},
			},
		},
	})
}

var AlicloudGPDBAccountMap0 = map[string]string{}

func AlicloudGPDBAccountBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_gpdb_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_gpdb_zones.default.zones.2.id
}

resource "alicloud_vswitch" "default" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id      = data.alicloud_gpdb_zones.default.zones.3.id
  vswitch_name = var.name
}

resource "alicloud_gpdb_elastic_instance" "default" {
  engine                   = "gpdb"
  engine_version           = "6.0"
  seg_storage_type         = "cloud_essd"
  seg_node_num             = 4
  storage_size             = 50
  instance_spec            = "2C16G"
  db_instance_description  = var.name
  instance_network_type    = "VPC"
  payment_type             = "PayAsYouGo"
  vswitch_id               = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.default.*.id, [""])[0]
}
`, name)
}

package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudDasSwitchDasPro_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_das_switch_das_pro.default"
	ra := resourceAttrInit(resourceId, resourceAlicloudDasSwitchDasProMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDasSwitchDasPro")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccDasSwitchDasPro-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlicloudDasSwitchDasProBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":   "${alicloud_polardb_cluster.default.id}",
					"sql_retention": "30",
					"user_id":       "${data.alicloud_account.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":   CHECKSET,
						"sql_retention": "30",
						"user_id":       CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var resourceAlicloudDasSwitchDasProMap = map[string]string{
	"status": CHECKSET,
}

func resourceAlicloudDasSwitchDasProBasicDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_account" "default" {
	}

	resource "alicloud_vpc" "default" {
		vpc_name = var.name
        cidr_block = "172.16.0.0/16"
	}

	data "alicloud_polardb_zones" "default"{}

	resource "alicloud_vswitch" "default" {
        zone_id = data.alicloud_polardb_zones.default.ids[length(data.alicloud_polardb_zones.default.ids)-1]
		vpc_id     = alicloud_vpc.default.id
		vswitch_name = var.name
        cidr_block = "172.16.0.0/24"
	}

	data "alicloud_polardb_node_classes" "default" {
		db_type    = "MySQL"
		db_version = "8.0"
		pay_type   = "PostPaid"
		zone_id    = data.alicloud_polardb_zones.default.ids[length(data.alicloud_polardb_zones.default.ids)-1]
		category   = "Normal"
	}

	resource "alicloud_polardb_cluster" "default" {
		db_type       = "MySQL"
		db_version    = "8.0"
		pay_type      = "PostPaid"
		db_node_class = data.alicloud_polardb_node_classes.default.classes.0.supported_engines.0.available_resources.0.db_node_class
		vswitch_id    = alicloud_vswitch.default.id
		description   = "${var.name}"
	}
`, name)
}

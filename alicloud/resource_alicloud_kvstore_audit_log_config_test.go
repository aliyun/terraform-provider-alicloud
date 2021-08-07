package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudKvstoreAuditLogConfig_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kvstore_audit_log_config.default"
	ra := resourceAttrInit(resourceId, KvstoreAuditLogConfigMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RKvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreAuditLogConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKvstoreAuditLogConfigtftestnormal%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, KvstoreAuditLogConfigBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckPrePaidResources(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_audit":    "true",
					"retention":   "10",
					"instance_id": "${alicloud_kvstore_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_audit":    "true",
						"retention":   "10",
						"instance_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_audit": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_audit": "false",
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

var KvstoreAuditLogConfigMap = map[string]string{}

func KvstoreAuditLogConfigBasicdependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	data "alicloud_kvstore_zones" "default" {
	  instance_charge_type = "PostPaid"
	  engine               = "redis"
	  product_type   = "Local"
	}
	data "alicloud_kvstore_instance_classes" "default" {
	  zone_id        = data.alicloud_kvstore_zones.default.ids.0
	  engine         = "Redis"
	  engine_version = "4.0"
	  product_type   = "Local"
	}
	data "alicloud_vpcs" "default" {
	  is_default = true
	}
	data "alicloud_vswitches" "default" {
	  vpc_id  = data.alicloud_vpcs.default.ids.0
	  zone_id = data.alicloud_kvstore_zones.default.ids.0
	}
	resource "alicloud_vswitch" "vswitch" {
	  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
	  vpc_id       = data.alicloud_vpcs.default.ids.0
	  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
	  zone_id      = data.alicloud_kvstore_zones.default.ids.0
	  vswitch_name = var.name
	}

	locals {
	  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
	}
	
	resource "alicloud_kvstore_instance" "default" {
	  instance_class = data.alicloud_kvstore_instance_classes.default.instance_classes.0
	  db_instance_name  = var.name
	  vswitch_id     = local.vswitch_id
	  security_ips   = ["10.0.0.1"]
	  instance_type  = "Redis"
	  engine_version = "4.0"
	}
	`, name)
}

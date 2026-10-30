package sddp_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/sddp"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudSDDPDataLimit_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sddp_data_limit.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSDDPDataLimitMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sddp.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSddpDataLimit")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%ssddpdatalimit%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSDDPDataLimitBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.CheckoutSupportedRegions(t, true, []connectivity.Region{sddpDataLimitsDataSourceTestRegion})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"lang":              "lang",
					"audit_status":      "0",
					"engine_type":       "MySQL",
					"parent_id":         "${local.parent_id}",
					"resource_type":     "RDS",
					"user_name":         "${var.database_name}",
					"password":          "${var.password}",
					"port":              "3306",
					"service_region_id": "${var.region}",
					"depends_on":        []string{"alicloud_db_account_privilege.default"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"audit_status":  "0",
						"engine_type":   "MySQL",
						"parent_id":     acctest.CHECKSET,
						"resource_type": "RDS",
						"user_name":     acctest.CHECKSET,
						"password":      acctest.CHECKSET,
						"port":          "3306",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"audit_status":  "1",
					"log_store_day": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"audit_status":  "1",
						"log_store_day": "30",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"audit_status": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"audit_status": "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang", "password"},
			},
		},
	})
}

var AlicloudSDDPDataLimitMap0 = map[string]string{}

func AlicloudSDDPDataLimitBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

variable "region" {
  default = "%s"
}

variable "password" {
  default = "Test12345"
}

variable "database_name" {
  default = "tftestdatabase"
}

data "alicloud_db_instance_classes" "default" {
  zone_id                  = "ap-southeast-1c"
  engine                   = "MySQL"
  engine_version           = "8.0"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
  instance_charge_type     = "PostPaid"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 2)
  zone_id      = "ap-southeast-1c"
  vswitch_name = var.name
}

resource "alicloud_db_instance" "default" {
  engine           = "MySQL"
  engine_version   = "8.0"
  instance_type    = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
  vswitch_id       = alicloud_vswitch.default.id
  instance_name    = var.name
}


locals {
  parent_id = join(".", [alicloud_db_instance.default.id, var.database_name])
}

resource "alicloud_rds_account" "default" {
  db_instance_id   = alicloud_db_instance.default.id
  account_name     = var.database_name
  account_password = var.password
}

resource "alicloud_db_database" "default" {
  instance_id = alicloud_db_instance.default.id
  name        = var.database_name
}

resource "alicloud_db_account_privilege" "default" {
  instance_id  = alicloud_db_instance.default.id
  account_name = alicloud_rds_account.default.name
  privilege    = "ReadWrite"
  db_names     = [alicloud_db_database.default.name]
}
`, name, sddpDataLimitsDataSourceTestRegion)
}

package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudDTSConsumerChannel_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dts_consumer_channel.default"
	checkoutSupportedRegions(t, true, connectivity.DTSSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudDTSConsumerChannelMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DtsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDtsConsumerChannel")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tftest%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDTSConsumerChannelBasicDependence0)
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
					"consumer_group_password":  "tftestAcc123",
					"consumer_group_user_name": name,
					"consumer_group_name":      name,
					"dts_instance_id":          "${alicloud_dts_subscription_job.default.dts_instance_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"consumer_group_user_name": name,
						"consumer_group_name":      name,
						"dts_instance_id":          CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"consumer_group_password": "tftestAcc123" + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"consumer_group_password"},
			},
		},
	})
}

var AlicloudDTSConsumerChannelMap0 = map[string]string{}

func AlicloudDTSConsumerChannelBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

variable "creation" {
  default = "Rds"
}

data "alicloud_zones" "default" {
  available_resource_creation = var.creation
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids[0]
  zone_id = data.alicloud_zones.default.zones[0].id
}

resource "alicloud_db_instance" "instance" {
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = "rds.mysql.s1.small"
  instance_storage = "10"
  vswitch_id       = data.alicloud_vswitches.default.ids[0]
  instance_name    = var.name
}

resource "alicloud_db_database" "db" {
  count       = 2
  instance_id = alicloud_db_instance.instance.id
  name        = "tfaccountpri_${count.index}"
  description = "from terraform"
}

resource "alicloud_db_account" "account" {
  db_instance_id      = alicloud_db_instance.instance.id
  account_name        = "tftestprivilege"
  account_password    = "Test12345"
  account_description = "from terraform"
}

resource "alicloud_db_account_privilege" "privilege" {
  instance_id  = alicloud_db_instance.instance.id
  account_name = alicloud_db_account.account.name
  privilege    = "ReadWrite"
  db_names     = alicloud_db_database.db.*.name
}

resource "alicloud_dts_subscription_job" "default" {
  dts_job_name                       = var.name
  payment_type                       = "PayAsYouGo"
  source_endpoint_engine_name        = "MySQL"
  source_endpoint_region             = "cn-hangzhou"
  source_endpoint_instance_type      = "RDS"
  source_endpoint_instance_id        = alicloud_db_instance.instance.id
  source_endpoint_database_name      = "tfaccountpri_0"
  source_endpoint_user_name          = "tftestprivilege"
  source_endpoint_password           = "Test12345"
  subscription_instance_network_type = "vpc"
  db_list                            = <<EOF
        {"dtstestdata": {"name": "tfaccountpri_0", "all": true}}
    EOF
  subscription_instance_vpc_id       = data.alicloud_vpcs.default.ids[0]
  subscription_instance_vswitch_id   = data.alicloud_vswitches.default.ids[0]
  status                             = "Normal"
}
`, name)
}

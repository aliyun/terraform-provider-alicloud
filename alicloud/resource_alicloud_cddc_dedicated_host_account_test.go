package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAliCloudCddcDedicatedHostAccount_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cddc_dedicated_host_account.default"
	ra := resourceAttrInit(resourceId, AlicloudCDDCDedicatedHostAccountMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CddcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCddcDedicatedHostAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCDDCDedicatedHostAccountBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.CDDCSupportRegions)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dedicated_host_id": "${alicloud_cddc_dedicated_host.default.dedicated_host_id}",
					"account_type":      "Normal",
					"account_password":  "Test1234+!",
					"account_name":      name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dedicated_host_id": CHECKSET,
						"account_name":      name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "Test1234+!" + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "Test1234+!",
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

var AlicloudCDDCDedicatedHostAccountMap0 = map[string]string{}

func AlicloudCDDCDedicatedHostAccountBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_cddc_zones" "default" {}

data "alicloud_cddc_host_ecs_level_infos" "default" {
  db_type      = "mysql"
  zone_id      = data.alicloud_cddc_zones.default.ids.0
  storage_type = "cloud_essd"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_cddc_dedicated_host_groups.default.groups.0.vpc_id
  zone_id = data.alicloud_cddc_zones.default.ids.0
}

data "alicloud_cddc_dedicated_host_groups" "default" {
  name_regex = "^NO-DELETING"
  engine     = "MySQL"
}

resource "alicloud_cddc_dedicated_host" "default" {
  host_name               = var.name
  dedicated_host_group_id = data.alicloud_cddc_dedicated_host_groups.default.ids.0
  host_class              = data.alicloud_cddc_host_ecs_level_infos.default.infos.0.res_class_code
  zone_id                 = data.alicloud_cddc_zones.default.ids.0
  vswitch_id              = data.alicloud_vswitches.default.ids.0
  payment_type            = "Subscription"
}
`, name)
}

// lintignore: R001

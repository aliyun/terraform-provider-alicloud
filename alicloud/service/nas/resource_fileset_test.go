package nas_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/nas"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudNASFileset_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_fileset.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.NASCPFSSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudNASFilesetMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewNasService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNasFileset")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snasfileset%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudNASFilesetBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"file_system_path":    "/tf-testAcc-Path/",
					"description":         name,
					"file_system_id":      "${alicloud_nas_file_system.default.id}",
					"deletion_protection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_system_path":    "/tf-testAcc-Path/",
						"description":         name,
						"file_system_id":      acctest.CHECKSET,
						"deletion_protection": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudNASFilesetMap0 = map[string]string{}

func AlicloudNASFilesetBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_nas_zones" "default" {
  file_system_type = "cpfs"
}

locals {
  count_size = length(data.alicloud_nas_zones.default.zones)
  zone_id = data.alicloud_nas_zones.default.zones[local.count_size - 1].zone_id
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = local.zone_id
}

resource "alicloud_nas_file_system" "default" {
  protocol_type    = "cpfs"
  storage_type     = "advance_200"
  file_system_type = "cpfs"
  capacity         = 3600
  zone_id          = local.zone_id
  vpc_id           = data.alicloud_vpcs.default.ids.0
  vswitch_id       = data.alicloud_vswitches.default.ids.0
}
`, name)
}

// lintignore: R001

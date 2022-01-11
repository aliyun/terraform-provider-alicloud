package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudNASSnapshot_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_snapshot.default"
	checkoutSupportedRegions(t, true, connectivity.NASSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudNASSnapshotMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNasSnapshot")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snassnapshot%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNASSnapshotBasicDependence0)
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
					"snapshot_name":  name,
					"file_system_id": "${alicloud_nas_file_system.default.id}",
					"description":    name,
					"retention_days": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snapshot_name":  name,
						"file_system_id": CHECKSET,
						"description":    name,
						"retention_days": "1",
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

var AlicloudNASSnapshotMap0 = map[string]string{
	"status": CHECKSET,
}

func AlicloudNASSnapshotBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_nas_zones" "default" {
  file_system_type = "extreme"
}

locals {
  count_size = length(data.alicloud_nas_zones.default.zones)
}

resource "alicloud_nas_file_system" "default" {
  file_system_type = "extreme"
  protocol_type    = "NFS"
  zone_id          = data.alicloud_nas_zones.default.zones[local.count_size - 1].zone_id
  storage_type     = "standard"
  description      = var.name
  capacity         = 100
}
`, name)
}

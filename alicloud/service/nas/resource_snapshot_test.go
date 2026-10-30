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

func TestAccAliCloudNASSnapshot_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_snapshot.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.NASSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudNASSnapshotMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewNasService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNasSnapshot")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snassnapshot%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudNASSnapshotBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"file_system_id": acctest.CHECKSET,
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
	"status": acctest.CHECKSET,
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
  storage_type     = "advance"
  capacity         = 100
}
`, name)
}

// lintignore: R001

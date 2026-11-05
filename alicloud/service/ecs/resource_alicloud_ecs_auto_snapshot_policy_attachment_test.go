package ecs_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecs"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudEcsAutoSnapshotPolicyAttachment_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_auto_snapshot_policy_attachment.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsAutoSnapshotPolicyAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcsAutoSnapshotPolicyAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccAliCloudEcsAutoSnapshotPolicyAttachment%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsAutoSnapshotPolicyAttachmentBasicDependence0)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_snapshot_policy_id": "${alicloud_ecs_auto_snapshot_policy.default.id}",
					"disk_id":                 "${alicloud_ecs_disk.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_snapshot_policy_id": CHECKSET,
						"disk_id":                 CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudEcsAutoSnapshotPolicyAttachmentMap0 = map[string]string{
	"region_id": CHECKSET,
}

func AliCloudEcsAutoSnapshotPolicyAttachmentBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_zones" "default" {
  		available_resource_creation = "VSwitch"
	}

	resource "alicloud_ecs_auto_snapshot_policy" "default" {
  		name            = "var.name"
  		repeat_weekdays = ["1", "2", "3"]
  		retention_days  = 1
  		time_points     = ["1", "2", "3"]
	}

	resource "alicloud_ecs_disk" "default" {
  		zone_id = data.alicloud_zones.default.zones.0.id
  		size    = "500"
	}
`, name)
}

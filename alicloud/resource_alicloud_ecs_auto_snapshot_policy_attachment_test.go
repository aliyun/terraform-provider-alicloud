package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudEcsAutoSnapshotPolicyAttachmentBasic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_auto_snapshot_policy_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsAutoSnapshotPolicyAttachmentMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccAlicloudEcsAutoSnapshotPolicyAttachment%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsAutoSnapshotPolicyAttachmentBasicDependence)
	testAccCheck := rac.resourceAttrMapUpdateSet()
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
					"disk_id":                 "${alicloud_ecs_disk.default.id}",
					"auto_snapshot_policy_id": "${alicloud_ecs_auto_snapshot_policy.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName: resourceId,
				ImportState:  true,
			},
		},
	})
}

var AlicloudEcsAutoSnapshotPolicyAttachmentMap = map[string]string{
	"auto_snapshot_policy_id": CHECKSET,
	"disk_id":                 CHECKSET,
}

func AlicloudEcsAutoSnapshotPolicyAttachmentBasicDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
			default = "%s"
		}

	resource "alicloud_ecs_auto_snapshot_policy" "default" {
		name              = "var.name"
		repeat_weekdays   = ["1"]
		retention_days    =  -1
		time_points       = ["1"]
		tags 	 = {
			Created = "TF"
			For 	= "acceptance test"
		}
	}
	data "alicloud_zones" default {
	  available_resource_creation = "Instance"
	}
	resource "alicloud_ecs_disk" "default" {
		zone_id = "${data.alicloud_zones.default.zones.0.id}"
		category = "cloud_efficiency"
		delete_auto_snapshot = "true"
		description = "Test For Terraform"
		disk_name = var.name
		enable_auto_snapshot = "true"
		encrypted = "true"
		size = "500"
		tags = {
			Created     = "TF"
			Environment = "Acceptance-test"
		}
	}
		
`, name)
}

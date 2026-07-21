package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ebs EnterpriseSnapshotPolicyAttachment. >>> Resource test cases, automatically generated.
// Case 5526
func TestAccAliCloudEbsEnterpriseSnapshotPolicyAttachment_basic5526(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ebs_enterprise_snapshot_policy_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudEbsEnterpriseSnapshotPolicyAttachmentMap5526)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EbsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEbsEnterpriseSnapshotPolicyAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sebsenterprisesnapshotpolicyattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEbsEnterpriseSnapshotPolicyAttachmentBasicDependence5526)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.EbsEnterpriseSnapshotPolicySupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_id": "${alicloud_ebs_enterprise_snapshot_policy.defaultPE3jjR.id}",
					"disk_id":   "${alicloud_ecs_disk.defaultJkW46o.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_id": CHECKSET,
						"disk_id":   CHECKSET,
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

var AlicloudEbsEnterpriseSnapshotPolicyAttachmentMap5526 = map[string]string{
	"disk_id": CHECKSET,
}

func AlicloudEbsEnterpriseSnapshotPolicyAttachmentBasicDependence5526(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_ecs_disk" "defaultJkW46o" {
  category          = "cloud_essd"
  description       = "esp-attachment-test"
  zone_id           = "cn-hangzhou-i"
  performance_level = "PL1"
  size              = "20"
  disk_name         = var.name

}

resource "alicloud_ebs_enterprise_snapshot_policy" "defaultPE3jjR" {
  status = "DISABLED"
  desc   = "DESC"
  schedule {
    cron_expression = "0 0 0 1 * ?"
  }
  enterprise_snapshot_policy_name = var.name

  target_type = "DISK"
  retain_rule {
    time_interval = "120"
    time_unit     = "DAYS"
  }
}


`, name)
}

// Test Ebs EnterpriseSnapshotPolicyAttachment. <<< Resource test cases, automatically generated.

package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func SkipTestAccAlicloudDiskDiskReplicaPair_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ebs_disk_replica_pair.default"
	ra := resourceAttrInit(resourceId, AlicloudDiskDiskReplicaPairMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EbsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEbsDiskReplicaPair")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sDiskDiskReplicaPair%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDiskDiskReplicaPairBasicDependence)
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
					"destination_disk_id":   "${alicloud_ecs_disk.default.id}",
					"destination_region_id": "cn-hangzhou-onebox-nebula",
					"destination_zone_id":   "cn-hangzhou-onebox-nebula-e",
					"source_zone_id":        "cn-hangzhou-onebox-nebula-b",
					"disk_id":               "${alicloud_ecs_disk.defaultone.id}",
					"description":           name,
					"pair_name":             name,
					"payment_type":          "POSTPAY",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_disk_id":   CHECKSET,
						"destination_region_id": "cn-hangzhou-onebox-nebula",
						"destination_zone_id":   "cn-hangzhou-onebox-nebula-e",
						"source_zone_id":        "cn-hangzhou-onebox-nebula-b",
						"disk_id":               CHECKSET,
						"description":           name,
						"pair_name":             name,
						"payment_type":          "POSTPAY",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pair_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pair_name": name + "_update",
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

var AlicloudDiskDiskReplicaPairMap = map[string]string{}

func AlicloudDiskDiskReplicaPairBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_ecs_disk" "default" {
	zone_id = "cn-hangzhou-onebox-nebula"
	category = "cloud_essd"
	delete_auto_snapshot = "true"
	delete_with_instance = "true"
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

resource "alicloud_ecs_disk" "defaultone" {
	zone_id = "cn-hangzhou-onebox-nebula-b"
	category = "cloud_essd"
	delete_auto_snapshot = "true"
	delete_with_instance = "true"
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

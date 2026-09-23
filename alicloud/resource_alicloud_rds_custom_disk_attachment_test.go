// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Rds CustomDiskAttachment. >>> Resource test cases, automatically generated.
// Case resourceCase_20260407_03_clone_0 12736
func TestAccAliCloudRdsCustomDiskAttachment_basic12736(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_custom_disk_attachment.default"
	ra := resourceAttrInit(resourceId, AliCloudRdsCustomDiskAttachmentMap12736)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsCustomDiskAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRdsCustomDiskAttachmentBasicDependence12736)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${alicloud_rds_custom.default.id}",
					//"delete_with_instance": "true",
					"disk_id": "${alicloud_rds_custom_disk.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
						//"delete_with_instance": "true",
						"disk_id": CHECKSET,
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

func TestAccAliCloudRdsCustomDiskAttachment_basic12736_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_custom_disk_attachment.default"
	ra := resourceAttrInit(resourceId, AliCloudRdsCustomDiskAttachmentMap12736)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsCustomDiskAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRdsCustomDiskAttachmentBasicDependence12736)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":          "${alicloud_rds_custom.default.id}",
					"delete_with_instance": "false",
					"disk_id":              "${alicloud_rds_custom_disk.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":          CHECKSET,
						"delete_with_instance": "false",
						"disk_id":              CHECKSET,
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

var AliCloudRdsCustomDiskAttachmentMap12736 = map[string]string{
	"region_id":            CHECKSET,
	"status":               CHECKSET,
	"delete_with_instance": CHECKSET,
}

func AliCloudRdsCustomDiskAttachmentBasicDependence12736(name string) string {
	return fmt.Sprintf(`
	variable "name" {
    	default = "%s"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}
	
	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = "cn-hangzhou-i"
	}

	data "alicloud_security_groups" "default" {
  		vpc_id     = data.alicloud_vpcs.default.ids.0
  		name_regex = "default-NODELETING"
	}

resource "alicloud_rds_custom" "default" {
  zone_id              = data.alicloud_vswitches.default.zone_id
  instance_charge_type = "PostPaid"
  vswitch_id           = data.alicloud_vswitches.default.ids.0
  amount               = "1"
  security_group_ids   = [data.alicloud_security_groups.default.ids.0]
  system_disk {
    size = "40"
  }
  force         = true
  instance_type = "mysql.x4.xlarge.6cm"
  spot_strategy = "NoSpot"
}

resource "alicloud_rds_custom_disk" "default" {
  zone_id       = data.alicloud_vswitches.default.zone_id
  size          = "40"
  disk_category = "cloud_ssd"
  auto_pay      = true
  disk_name     = "ran_disk_attach"
}
`, name)
}

// Test Rds CustomDiskAttachment. <<< Resource test cases, automatically generated.

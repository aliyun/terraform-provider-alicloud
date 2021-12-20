package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_ecs_disk", &resource.Sweeper{
		Name: "alicloud_ecs_disk",
		F:    testAlicloudEcsDisk,
	})
}

func testAlicloudEcsDisk(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
	}
	request := map[string]interface{}{
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
		"RegionId":   client.RegionId,
	}

	var response map[string]interface{}
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		action := "DescribeDisks"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			log.Printf("[ERROR] %s got an error: %s", action, err)
			return nil
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.Disks.Disk", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Disks.Disk", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["DiskName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Disk: %s", item["DiskName"].(string))
				continue
			}
			action = "DeleteDisk"
			request := map[string]interface{}{
				"DiskId":   item["DiskId"],
				"RegionId": client.RegionId,
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Disk (%s): %s", item["DiskName"].(string), err)
			}
			log.Printf("[INFO] Delete Disk success: %s ", item["DiskName"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudEcsDisk_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsDiskMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsdisk%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsDiskBasicDependence)
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
					"disk_name":    name,
					"zone_id":      "${data.alicloud_zones.default.zones.0.id}",
					"encrypted":    "true",
					"kms_key_id":   "${alicloud_kms_key.key.id}",
					"size":         "500",
					"payment_type": "PayAsYouGo",
					"timeouts": []map[string]interface{}{
						{
							"create": "1h",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_name":  name,
						"zone_id":    CHECKSET,
						"encrypted":  "true",
						"kms_key_id": CHECKSET,
						"size":       "500",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"advanced_features", "type", "dedicated_block_storage_cluster_id", "dry_run", "encrypt_algorithm", "storage_set_id", "storage_set_partition_number"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":          "cloud_essd",
						"performance_level": "PL1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"performance_level": "PL2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"performance_level": "PL2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_auto_snapshot": `true`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_auto_snapshot": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_with_instance": `true`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_with_instance": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "Test For Terraform Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "Test For Terraform Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_name": name + "Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_name": name + "Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_auto_snapshot": `true`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_auto_snapshot": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"size": `700`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"size": "700",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_auto_snapshot": `false`,
					"delete_with_instance": `false`,
					"delete_auto_snapshot": `false`,
					"description":          "Test For Terraform",
					"disk_name":            name,
					"performance_level":    "PL1",
					"size":                 "800",
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_auto_snapshot": `false`,
						"delete_with_instance": `false`,
						"delete_auto_snapshot": `false`,
						"description":          "Test For Terraform",
						"disk_name":            name,
						"performance_level":    "PL1",
						"size":                 "800",
						"tags.%":               "2",
						"tags.Created":         "TF-update",
						"tags.For":             "Test-update",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudEcsDisk_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsDiskMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsdisk%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsDiskBasic1Dependence)
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
					"disk_name":    name,
					"size":         "500",
					"payment_type": "Subscription",
					"instance_id":  "${alicloud_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_name":    name,
						"zone_id":      CHECKSET,
						"size":         "500",
						"instance_id":  CHECKSET,
						"payment_type": "Subscription",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"advanced_features", "type", "dedicated_block_storage_cluster_id", "dry_run", "encrypt_algorithm", "storage_set_id", "storage_set_partition_number"},
			},
		},
	})
}

var AlicloudEcsDiskMap = map[string]string{
	"payment_type":      "PayAsYouGo",
	"performance_level": "",
}

func AlicloudEcsDiskBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
			default = "%s"
		}
data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}
resource "alicloud_kms_key" "key" {
	description             = var.name
	pending_window_in_days  = "7"
	key_state               = "Enabled"
}
`, name)
}

func AlicloudEcsDiskBasic1Dependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
			default = "%s"
		}
data "alicloud_instance_types" "default" {
  cpu_core_count       = 2
  memory_size          = 4
  instance_charge_type = "PrePaid"
}
data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  owners      = "system"
}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
}
resource "alicloud_security_group" "default" {
  name        = var.name
  vpc_id      = data.alicloud_vswitches.default.vswitches.0.vpc_id
}
resource "alicloud_instance" "default" {
  image_id                      = data.alicloud_images.default.images.0.id
  security_groups               = [alicloud_security_group.default.id]
  instance_type                 = data.alicloud_instance_types.default.instance_types.0.id
  system_disk_category          = "cloud_efficiency"
  instance_name                 = var.name
  spot_strategy                 = "NoSpot"
  spot_price_limit              = "0"
  security_enhancement_strategy = "Active"
  user_data                     = "I_am_user_data"
  instance_charge_type          = "PrePaid"
  period                        = 1
  vswitch_id                    = data.alicloud_vswitches.default.ids.0
  force_delete                  = true
}
`, name)
}

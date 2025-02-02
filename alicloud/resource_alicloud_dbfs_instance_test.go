package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_dbfs_instance",
		&resource.Sweeper{
			Name: "alicloud_dbfs_instance",
			F:    testSweepDBFSInstance,
		})
}

func testSweepDBFSInstance(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListDbfs"
	request := map[string]interface{}{}

	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	request["RegionId"] = client.RegionId

	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("DBFS", "2020-04-18", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.DBFSInfo", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.DBFSInfo", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			if _, ok := item["FsName"]; !ok {
				continue
			}
			skip := true
			if !sweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(item["FsName"].(string)), strings.ToLower(prefix)) {
						skip = false
					}
				}
				if skip {
					log.Printf("[INFO] Skipping DBFSDbfs Instance: %s", item["FsName"].(string))
					continue
				}
			}
			action := "DeleteDbfs"
			request := map[string]interface{}{
				"FsId": item["FsId"],
			}
			request["ClientToken"] = buildClientToken("DeleteDbfs")
			_, err = client.RpcPost("DBFS", "2020-04-18", action, nil, request, true)
			if err != nil {
				log.Printf("[ERROR] Failed to delete DBFSDbfs Instance (%s): %s", item["FsId"].(string), err)
			}
			log.Printf("[INFO] Delete DBFSDbfs Instance success: %s ", item["FsId"].(string))
		}

		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAliCloudDBFSInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dbfs_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudDBFSInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DbfsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDbfsInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdbfsinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDBFSInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DBFSSystemSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"category":          "enterprise",
					"zone_id":           "cn-hangzhou-i",
					"performance_level": "PL1",
					"instance_name":     name,
					"size":              "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":          "enterprise",
						"zone_id":           "cn-hangzhou-i",
						"performance_level": "PL1",
						"instance_name":     name,
						"size":              "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ecs_list": []map[string]interface{}{
						{
							"ecs_id": "${alicloud_instance.default.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ecs_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"size": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"size": "200",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF1",
						"For":     "Test1",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF2",
						"For":     "Test2",
						"number":  "2",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ecs_list": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ecs_list.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_snapshot"},
			},
		},
	})
}

var AlicloudDBFSInstanceMap0 = map[string]string{}

func AlicloudDBFSInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

locals {
  zone_id = "cn-hangzhou-i"
}
data "alicloud_instance_types" "example" {
  availability_zone    = local.zone_id
  instance_type_family = "ecs.g7se"
}
data "alicloud_images" "example" {
  instance_type = data.alicloud_instance_types.example.instance_types[length(data.alicloud_instance_types.example.instance_types) - 1].id
  name_regex    = "^aliyun_2"
  owners        = "system"
}

resource "alicloud_vpc" "default" {
    vpc_name = var.name
	cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id  = alicloud_vpc.default.id
  zone_id = local.zone_id
  vswitch_name = var.name
  cidr_block = "172.16.0.0/24"
}

resource "alicloud_security_group" "example" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_instance" "default" {
  availability_zone    = local.zone_id
  instance_name        = var.name
  image_id             = data.alicloud_images.example.images.0.id
  instance_type        = data.alicloud_instance_types.example.instance_types[length(data.alicloud_instance_types.example.instance_types) - 1].id
  security_groups      = [alicloud_security_group.example.id]
  vswitch_id           = alicloud_vswitch.default.id
  system_disk_category = "cloud_essd"
}
`, name)
}

// Test Dbfs DbfsInstance. >>> Resource test cases, automatically generated.
// Case 5069
func TestAccAliCloudDbfsDbfsInstance_basic5069(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dbfs_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudDbfsDbfsInstanceMap5069)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DbfsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDbfsDbfsInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdbfsdbfsinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDbfsDbfsInstanceBasicDependence5069)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DBFSSystemSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"category": "enterprise",
					"zone_id":  "cn-hangzhou-i",
					"size":     "20",
					"fs_name":  "rmc-test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category": "enterprise",
						"zone_id":  "cn-hangzhou-i",
						"size":     "20",
						"fs_name":  "rmc-test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"performance_level": "PL1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"performance_level": "PL1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"used_scene": "MongoDB",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"used_scene": "MongoDB",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "dbfs.small",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": "dbfs.small",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"advanced_features": "{\\\"memorySize\\\":1024,\\\"pageCacheSize\\\":128,\\\"cpuCoreCount\\\":0.5}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"advanced_features": "{\"memorySize\":1024,\"pageCacheSize\":128,\"cpuCoreCount\":0.5}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"size": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"size": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"fs_name": "rmc-test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fs_name": "rmc-test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"fs_name": "rmc-new-name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fs_name": "rmc-new-name",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"used_scene": "PostgreSQL ",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"used_scene": "PostgreSQL ",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "dbfs.medium",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": "dbfs.medium",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"size": "960",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"size": "960",
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
					"advanced_features": "{\\\"memorySize\\\":512,\\\"pageCacheSize\\\":128,\\\"cpuCoreCount\\\":0.5}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"advanced_features": "{\"memorySize\":512,\"pageCacheSize\":128,\"cpuCoreCount\":0.5}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"category":                "enterprise",
					"zone_id":                 "cn-hangzhou-i",
					"size":                    "20",
					"performance_level":       "PL1",
					"fs_name":                 "rmc-test",
					"used_scene":              "MongoDB",
					"instance_type":           "dbfs.small",
					"raid_stripe_unit_number": "2",
					"advanced_features":       "{\\\"memorySize\\\":1024,\\\"pageCacheSize\\\":128,\\\"cpuCoreCount\\\":0.5}",
					"kms_key_id":              "00000000-0000-0000-0000-000000000000",
					"snapshot_id":             "none",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":                "enterprise",
						"zone_id":                 "cn-hangzhou-i",
						"size":                    "20",
						"performance_level":       "PL1",
						"fs_name":                 "rmc-test",
						"used_scene":              "MongoDB",
						"instance_type":           "dbfs.small",
						"raid_stripe_unit_number": "2",
						"advanced_features":       "{\"memorySize\":1024,\"pageCacheSize\":128,\"cpuCoreCount\":0.5}",
						"kms_key_id":              "00000000-0000-0000-0000-000000000000",
						"snapshot_id":             "none",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_snapshot"},
			},
		},
	})
}

var AlicloudDbfsDbfsInstanceMap5069 = map[string]string{
	"status":            CHECKSET,
	"performance_level": "PL1",
	"create_time":       CHECKSET,
	"snapshot_id":       "none",
	"advanced_features": "{}",
}

func AlicloudDbfsDbfsInstanceBasicDependence5069(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 5069  twin
func TestAccAliCloudDbfsDbfsInstance_basic5069_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dbfs_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudDbfsDbfsInstanceMap5069)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DbfsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDbfsDbfsInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdbfsdbfsinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDbfsDbfsInstanceBasicDependence5069)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DBFSSystemSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"category":                "enterprise",
					"zone_id":                 "cn-hangzhou-i",
					"size":                    "960",
					"performance_level":       "PL2",
					"fs_name":                 "rmc-new-name",
					"used_scene":              "PostgreSQL ",
					"instance_type":           "dbfs.medium",
					"raid_stripe_unit_number": "2",
					"advanced_features":       "{\\\"memorySize\\\":512,\\\"pageCacheSize\\\":128,\\\"cpuCoreCount\\\":0.5}",
					"kms_key_id":              "00000000-0000-0000-0000-000000000000",
					"snapshot_id":             "none",
					"encryption":              "true",
					"delete_snapshot":         "true",
					"enable_raid":             "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":                "enterprise",
						"zone_id":                 "cn-hangzhou-i",
						"size":                    "960",
						"performance_level":       "PL2",
						"fs_name":                 "rmc-new-name",
						"used_scene":              "PostgreSQL ",
						"instance_type":           "dbfs.medium",
						"raid_stripe_unit_number": "2",
						"advanced_features":       "{\"memorySize\":512,\"pageCacheSize\":128,\"cpuCoreCount\":0.5}",
						"kms_key_id":              "00000000-0000-0000-0000-000000000000",
						"snapshot_id":             "none",
						"encryption":              "true",
						"delete_snapshot":         "true",
						"enable_raid":             "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_snapshot"},
			},
		},
	})
}

// Test Dbfs DbfsInstance. <<< Resource test cases, automatically generated.

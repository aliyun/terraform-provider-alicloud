package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

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
	conn, err := client.NewDbfsClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
	}

	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-18"), StringPointer("AK"), nil, request, &runtime)
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
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["FsName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping DBFSDbfs Instance: %s", item["FsName"].(string))
				continue
			}
			action := "DeleteDbfs"
			request := map[string]interface{}{
				"FsId": item["FsId"],
			}
			request["ClientToken"] = buildClientToken("DeleteDbfs")
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-18"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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

func TestAccAlicloudDBFSInstance_basic0(t *testing.T) {
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
					"category":          "standard",
					"zone_id":           "cn-hangzhou-i",
					"performance_level": "PL1",
					"instance_name":     name,
					"size":              "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":          "standard",
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
				ImportStateVerifyIgnore: []string{"delete_snapshot", "snapshot_id"},
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
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids[0]
  zone_id = "cn-hangzhou-i"
}
resource "alicloud_security_group" "default" {
  name        = var.name
  description = "tf test"
  vpc_id  = data.alicloud_vpcs.default.ids[0]
}
data "alicloud_images" "default" {
  owners      = "system"
  name_regex  = "^centos_8"
  most_recent = true
}
resource "alicloud_instance" "default" {
  image_id          = data.alicloud_images.default.images[0].id
  instance_name     = var.name
  instance_type     = "ecs.g7se.large"
  availability_zone = "cn-hangzhou-i"
  vswitch_id        = data.alicloud_vswitches.default.ids[0]
  system_disk_category = "cloud_essd"
  security_groups = [
    alicloud_security_group.default.id
  ]
}
`, name)
}

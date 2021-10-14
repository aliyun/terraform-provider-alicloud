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
		"alicloud_dfs_file_system",
		&resource.Sweeper{
			Name: "alicloud_dfs_file_system",
			F:    testSweepDFSFileSystem,
		})
}

func testSweepDFSFileSystem(region string) error {
	rawClient, err := sharedClientForRegionWithBackendRegions(region, true, connectivity.DfsSupportRegions)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	request := map[string]interface{}{
		"InputRegionId": client.RegionId,
	}

	action := "ListFileSystems"
	conn, err := client.NewAlidfsClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
	}
	var response map[string]interface{}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-06-20"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
		return nil
	}

	resp, err := jsonpath.Get("$.FileSystems", response)
	if err != nil {
		log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.FileSystems", action, err)
		return nil
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})

		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(item["FileSystemName"].(string)), strings.ToLower(prefix)) {
				skip = false
			}
		}
		if skip {
			log.Printf("[INFO] Skipping DFS FileSystem: %s", item["FileSystemName"].(string))
			continue
		}

		action := "DeleteFileSystem"
		request := map[string]interface{}{
			"FileSystemId":  item["FileSystemId"].(string),
			"InputRegionId": client.RegionId,
		}

		_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-06-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			log.Printf("[ERROR] Failed to delete DFS FileSystem (%s): %s", item["FileSystemName"].(string), err)
		}
		log.Printf("[INFO] Delete  DFS FileSystem success: %s ", item["FileSystemName"].(string))
	}

	return nil
}

func TestAccAlicloudDFSFileSystem_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dfs_file_system.default"
	checkoutSupportedRegions(t, true, connectivity.DfsSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudDFSFileSystemMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DfsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDfsFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdfsfilesystem%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDFSFileSystemBasicDependence0)
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
					"storage_type":                     "${local.storage_type}",
					"zone_id":                          "${local.zone_id}",
					"protocol_type":                    "HDFS",
					"description":                      name,
					"file_system_name":                 name,
					"space_capacity":                   "1024",
					"throughput_mode":                  "Provisioned",
					"provisioned_throughput_in_mi_bps": "512",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_type":                     CHECKSET,
						"zone_id":                          CHECKSET,
						"protocol_type":                    "HDFS",
						"description":                      name,
						"file_system_name":                 name,
						"space_capacity":                   "1024",
						"throughput_mode":                  "Provisioned",
						"provisioned_throughput_in_mi_bps": "512",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"file_system_name": name + "Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_system_name": name + "Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"space_capacity": "2048",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"space_capacity": "2048",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":      name,
					"file_system_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":      name,
						"file_system_name": name,
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

var AlicloudDFSFileSystemMap0 = map[string]string{}

func AlicloudDFSFileSystemBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_dfs_zones" "default" {}
locals {
  zone_id      = data.alicloud_dfs_zones.default.zones.0.zone_id
  storage_type = data.alicloud_dfs_zones.default.zones.0.options.0.storage_type
}
`, name)
}

package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_nas_file_system",
		&resource.Sweeper{
			Name: "alicloud_nas_file_system",
			F:    testSweepNasFileSystem,
		})
}

func testSweepNasFileSystem(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeFileSystems"
	request := make(map[string]interface{})
	request["RegionId"] = client.Region
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var response map[string]interface{}
	conn, err := client.NewNasClient()
	if err != nil {
		return WrapError(err)
	}
	ids := make([]string, 0)
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-26"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			log.Printf("[ERROR] Error retrieving filesystem: %s", err)
		}
		resp, err := jsonpath.Get("$.FileSystems.FileSystem", response)
		if err != nil {
			log.Println("Get $.FileSystems.FileSystem failed. err:", err)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			description, _ := item["Description"].(string)
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(description), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping FileSystem: %s (%s)", description, item["FileSystemId"])
				continue
			}
			// 删除 fileSystem 时需要先删除其挂载关系
			if v, ok := item["MountTargets"].(map[string]interface{})["MountTarget"].([]interface{}); ok && len(v) > 0 {
				log.Printf("[INFO] Delete mount targets with filesystem: %v", item["FileSystemId"])
				for _, domain := range v {
					domainInfo := domain.(map[string]interface{})
					request := map[string]interface{}{
						"FileSystemId":      item["FileSystemId"],
						"MountTargetDomain": domainInfo["MountTargetDomain"],
					}
					action := "DeleteMountTarget"
					runtime := util.RuntimeOptions{}
					runtime.SetAutoretry(true)
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-26"), StringPointer("AK"), nil, request, &runtime)
					if err != nil {
						log.Printf("[ERROR] Error delete mount target: %v with filesystem: %v err: %v", domainInfo["MountTargetDomain"], item["FileSystemId"], err)
					}
				}
			}
			ids = append(ids, item["FileSystemId"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	for _, filesystemId := range ids {
		request := map[string]interface{}{
			"FileSystemId": filesystemId,
		}
		action := "DeleteFileSystem"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-26"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			log.Printf("[ERROR] Error delete filesystem: %s err: %v", filesystemId, err)
		}
	}
	return nil
}

func TestAccAlicloudNasFileSystem_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_file_system.default"
	ra := resourceAttrInit(resourceId, AlicloudNasFileSystem0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNasFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudNasFileSystem%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNasFileSystemBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.NasNoSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol_type": "${data.alicloud_nas_protocols.example.protocols.0}",
					"storage_type":  "Capacity",
					"zone_id":       "${data.alicloud_nas_zones.default.zones.1.zone_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type": CHECKSET,
						"storage_type":  "Capacity",
						"zone_id":       CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
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
		},
	})
}

func TestAccAlicloudNasFileSystemEncrypt(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_file_system.default"
	ra := resourceAttrInit(resourceId, AlicloudNasFileSystem0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNasFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudNasFileSystem%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNasFileSystemBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.NasNoSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol_type": "NFS",
					"storage_type":  "Capacity",
					"encrypt_type":  "1",
					"zone_id":       "${data.alicloud_nas_zones.default.zones.1.zone_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type": CHECKSET,
						"storage_type":  "Capacity",
						"encrypt_type":  "1",
						"zone_id":       CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
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
		},
	})
}

func TestAccAlicloudNasFileSystemExtreme_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_file_system.default"
	ra := resourceAttrInit(resourceId, AlicloudNasFileSystem0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNasFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudNasFileSystem%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNasFileSystemBasicDependence2)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.NasNoSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol_type":    "NFS",
					"zone_id":          "${local.zone_id}",
					"storage_type":     "standard",
					"file_system_type": "extreme",
					"capacity":         "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type":    CHECKSET,
						"zone_id":          CHECKSET,
						"storage_type":     "standard",
						"file_system_type": "extreme",
						"capacity":         "100",
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

func TestAccAlicloudNasFileSystemExtremeEncrypt(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_file_system.default"
	ra := resourceAttrInit(resourceId, AlicloudNasFileSystem0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNasFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudNasFileSystem%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNasFileSystemBasicDependence3)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.NasClassicSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol_type":    "NFS",
					"zone_id":          "${local.zone_id}",
					"storage_type":     "standard",
					"file_system_type": "extreme",
					"capacity":         "100",
					"encrypt_type":     "2",
					"kms_key_id":       "${alicloud_kms_key.key.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type":    CHECKSET,
						"zone_id":          CHECKSET,
						"storage_type":     "standard",
						"file_system_type": "extreme",
						"capacity":         "100",
						"encrypt_type":     "2",
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

func TestAccAlicloudNasFileSystemCpfs_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_file_system.default"
	checkoutSupportedRegions(t, true, connectivity.NASCPFSSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudNasFileSystem0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNasFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudNasFileSystem%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNasFileSystemBasicDependence4)
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
					"protocol_type":    "cpfs",
					"zone_id":          "${local.zone_id}",
					"storage_type":     "advance_200",
					"file_system_type": "cpfs",
					"capacity":         "3600",
					"vpc_id":           "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":       "${data.alicloud_vswitches.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type":    CHECKSET,
						"zone_id":          CHECKSET,
						"storage_type":     "advance_200",
						"file_system_type": "cpfs",
						"capacity":         "3600",
						"vpc_id":           CHECKSET,
						"vswitch_id":       CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vpc_id", "vswitch_id"},
			},
		},
	})
}

func TestAccAlicloudNasFileSystemTags_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_file_system.default"
	ra := resourceAttrInit(resourceId, AlicloudNasFileSystem0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNasFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudNasFileSystem%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNasFileSystemBasicDependence5)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.NasNoSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol_type": "NFS",
					"storage_type":  "Capacity",
					"zone_id":       "${local.zone_id}",
					"description":   name,
					"tags": map[string]string{
						"Created": "TF2",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type": CHECKSET,
						"storage_type":  "Capacity",
						"zone_id":       CHECKSET,
						"description":   name,
						"tags.%":        "1",
						"tags.Created":  "TF2",
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
						"tags.%":       "2",
						"tags.Created": "TF1",
						"tags.For":     "Test1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF2",
						"For":     "Test2",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF2",
						"tags.For":     "Test2",
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

var AlicloudNasFileSystem0 = map[string]string{}

func AlicloudNasFileSystemBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_nas_protocols" "example" {
        type = "Capacity"
}
data "alicloud_nas_zones" "default" {
}
`, name)
}

func AlicloudNasFileSystemBasicDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_nas_zones" "default" {
}
`, name)
}

func AlicloudNasFileSystemBasicDependence2(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_nas_zones" "default" {
  file_system_type = "extreme"
}

locals {
  count_size = length(data.alicloud_nas_zones.default.zones)
  zone_id = data.alicloud_nas_zones.default.zones[local.count_size - 1].zone_id
}
`, name)
}

func AlicloudNasFileSystemBasicDependence3(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
resource "alicloud_kms_key" "key" {
 description             = var.name
 pending_window_in_days  = "7"
 key_state               = "Enabled"
}

data "alicloud_nas_zones" "default" {
  file_system_type = "extreme"
}

locals {
  count_size = length(data.alicloud_nas_zones.default.zones)
  zone_id = data.alicloud_nas_zones.default.zones[local.count_size - 1].zone_id
}
`, name)
}

func AlicloudNasFileSystemBasicDependence4(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_nas_zones" "default" {
  file_system_type = "cpfs"
}

locals {
  count_size = length(data.alicloud_nas_zones.default.zones)
  zone_id = data.alicloud_nas_zones.default.zones[local.count_size - 1].zone_id
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = local.zone_id
}
`, name)
}

func AlicloudNasFileSystemBasicDependence5(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_nas_zones" "default" {
  file_system_type = "standard"
}

locals {
  count_size = length(data.alicloud_nas_zones.default.zones)
  zone_id = data.alicloud_nas_zones.default.zones[local.count_size - 1].zone_id
}
`, name)
}

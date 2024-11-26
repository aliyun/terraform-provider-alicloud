package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/alibabacloud-go/tea-rpc/client"
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
		return fmt.Errorf("error getting AliCloud client: %s", err)
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
			if !sweepAll() {
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

func TestAccAliCloudNasFileSystem_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.NasClassicSupportedRegions)
	resourceId := "alicloud_nas_file_system.default"
	ra := resourceAttrInit(resourceId, AliCloudNasFileSystem0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNasFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudNasFileSystem%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudNasFileSystemBasicDependence0)
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
					"protocol_type": "NFS",
					"storage_type":  "Capacity",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type": "NFS",
						"storage_type":  "Capacity",
					}),
				),
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
					"recycle_bin": []map[string]interface{}{
						{
							"status":        "Enable",
							"reserved_days": "10",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"recycle_bin.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"recycle_bin": []map[string]interface{}{
						{
							"status": "Disable",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"recycle_bin.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"nfs_acl": []map[string]interface{}{
						{
							"enabled": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nfs_acl.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"nfs_acl": []map[string]interface{}{
						{
							"enabled": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nfs_acl.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "FileSystem",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "FileSystem",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"snapshot_id"},
			},
		},
	})
}

func TestAccAliCloudNasFileSystem_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.NasClassicSupportedRegions)
	resourceId := "alicloud_nas_file_system.default"
	ra := resourceAttrInit(resourceId, AliCloudNasFileSystem0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNasFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudNasFileSystem%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudNasFileSystemBasicDependence0)
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
					"protocol_type":    "NFS",
					"storage_type":     "Capacity",
					"description":      name,
					"encrypt_type":     "2",
					"file_system_type": "standard",
					"kms_key_id":       "${alicloud_kms_key.default.id}",
					"recycle_bin": []map[string]interface{}{
						{
							"status":        "Enable",
							"reserved_days": "10",
						},
					},
					"nfs_acl": []map[string]interface{}{
						{
							"enabled": "true",
						},
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "FileSystem",
					},
					"zone_id": "${data.alicloud_nas_zones.standard.zones.0.zone_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type":     "NFS",
						"storage_type":      "Capacity",
						"description":       name,
						"encrypt_type":      "2",
						"file_system_type":  "standard",
						"kms_key_id":        CHECKSET,
						"recycle_bin.#":     "1",
						"nfs_acl.#":         "1",
						"resource_group_id": CHECKSET,
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "FileSystem",
						"zone_id":           CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"snapshot_id"},
			},
		},
	})
}

func TestAccAliCloudNasFileSystem_basic1(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.NasClassicSupportedRegions)
	resourceId := "alicloud_nas_file_system.default"
	ra := resourceAttrInit(resourceId, AliCloudNasFileSystem0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNasFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudNasFileSystem%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudNasFileSystemBasicDependence0)
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
					"protocol_type":    "NFS",
					"storage_type":     "standard",
					"capacity":         "100",
					"file_system_type": "extreme",
					"zone_id":          "${data.alicloud_nas_zones.extreme.zones.0.zone_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type":    "NFS",
						"storage_type":     "standard",
						"capacity":         "100",
						"file_system_type": "extreme",
						"zone_id":          CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"capacity": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"capacity": "200",
					}),
				),
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
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "FileSystem",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "FileSystem",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"snapshot_id"},
			},
		},
	})
}

func TestAccAliCloudNasFileSystem_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.NasClassicSupportedRegions)
	resourceId := "alicloud_nas_file_system.default"
	ra := resourceAttrInit(resourceId, AliCloudNasFileSystem0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNasFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudNasFileSystem%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudNasFileSystemBasicDependence0)
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
					"protocol_type":     "NFS",
					"storage_type":      "standard",
					"capacity":          "100",
					"description":       name,
					"encrypt_type":      "2",
					"file_system_type":  "extreme",
					"kms_key_id":        "${alicloud_kms_key.default.id}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "FileSystem",
					},
					"zone_id": "${data.alicloud_nas_zones.extreme.zones.0.zone_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type":     "NFS",
						"storage_type":      "standard",
						"capacity":          "100",
						"description":       name,
						"encrypt_type":      "2",
						"file_system_type":  "extreme",
						"kms_key_id":        CHECKSET,
						"resource_group_id": CHECKSET,
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "FileSystem",
						"zone_id":           CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"snapshot_id"},
			},
		},
	})
}

func TestAccAliCloudNasFileSystem_basic2(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.NasClassicSupportedRegions)
	resourceId := "alicloud_nas_file_system.default"
	ra := resourceAttrInit(resourceId, AliCloudNasFileSystem0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNasFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudNasFileSystem%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudNasFileSystemBasicDependence0)
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
					"storage_type":     "advance_100",
					"capacity":         "5000",
					"file_system_type": "cpfs",
					"vswitch_id":       "${data.alicloud_vswitches.default.ids.0}",
					"vpc_id":           "${data.alicloud_vpcs.default.ids.0}",
					"zone_id":          "${data.alicloud_nas_zones.cpfs.zones.1.zone_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type":    "cpfs",
						"storage_type":     "advance_100",
						"capacity":         "5000",
						"file_system_type": "cpfs",
						"vswitch_id":       CHECKSET,
						"vpc_id":           CHECKSET,
						"zone_id":          CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"capacity": "6000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"capacity": "6000",
					}),
				),
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
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "FileSystem",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "FileSystem",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"snapshot_id"},
			},
		},
	})
}

func TestAccAliCloudNasFileSystem_basic2_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.NasClassicSupportedRegions)
	resourceId := "alicloud_nas_file_system.default"
	ra := resourceAttrInit(resourceId, AliCloudNasFileSystem0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNasFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudNasFileSystem%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudNasFileSystemBasicDependence0)
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
					"protocol_type":     "cpfs",
					"storage_type":      "advance_100",
					"capacity":          "5000",
					"description":       name,
					"file_system_type":  "cpfs",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "FileSystem",
					},
					"vswitch_id": "${data.alicloud_vswitches.default.ids.0}",
					"vpc_id":     "${data.alicloud_vpcs.default.ids.0}",
					"zone_id":    "${data.alicloud_nas_zones.cpfs.zones.1.zone_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type":     "cpfs",
						"storage_type":      "advance_100",
						"capacity":          "5000",
						"description":       name,
						"file_system_type":  "cpfs",
						"resource_group_id": CHECKSET,
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "FileSystem",
						"vswitch_id":        CHECKSET,
						"vpc_id":            CHECKSET,
						"zone_id":           CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"snapshot_id"},
			},
		},
	})
}

var AliCloudNasFileSystem0 = map[string]string{
	"capacity":          CHECKSET,
	"create_time":       CHECKSET,
	"file_system_type":  CHECKSET,
	"resource_group_id": CHECKSET,
	"status":            CHECKSET,
	"zone_id":           CHECKSET,
}

func AliCloudNasFileSystemBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	data "alicloud_nas_zones" "standard" {
  		file_system_type = "standard"
	}

	data "alicloud_nas_zones" "extreme" {
  		file_system_type = "extreme"
	}

	data "alicloud_nas_zones" "cpfs" {
  		file_system_type = "cpfs"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_nas_zones.cpfs.zones.0.zone_id
	}

	resource "alicloud_kms_key" "default" {
  		description            = var.name
  		pending_window_in_days = "7"
  		key_state              = "Enabled"
	}
`, name)
}

func TestUnitAliCloudNasFileSystem(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_nas_file_system"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_nas_file_system"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"protocol_type":    "cpfs",
		"zone_id":          "zone_id",
		"storage_type":     "advance_200",
		"file_system_type": "cpfs",
		"capacity":         3600,
		"vpc_id":           "vpc_id",
		"vswitch_id":       "vswitch_id",
		"encrypt_type":     1,
		"kms_key_id":       "kms_key_id",
	} {
		err := dCreate.Set(key, value)
		assert.Nil(t, err)
		err = d.Set(key, value)
		assert.Nil(t, err)
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"FileSystems": map[string]interface{}{
			"FileSystem": []interface{}{
				map[string]interface{}{
					"Description":    "description",
					"ProtocolType":   "protocol_type",
					"StorageType":    "storage_type",
					"EncryptType":    1,
					"FileSystemType": "file_system_type",
					"Capacity":       3600,
					"ZoneId":         "zone_id",
					"KMSKeyId":       "kms_key_id",
					"tags": map[string]string{
						"Created": "TF1",
						"For":     "Test1",
					},
					"Status": "Running",
				},
			},
		},
	}

	responseMock := map[string]func(errorCode string) (map[string]interface{}, error){
		"RetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String(errorCode),
				Data:       String(errorCode),
				Message:    String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"NotFoundError": func(errorCode string) (map[string]interface{}, error) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_nas_file_system", "MockFileSystemId"))
		},
		"NoRetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String(errorCode),
				Data:       String(errorCode),
				Message:    String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"CreateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			result["FileSystemId"] = "MockFileSystemId"
			return result, nil
		},
		"UpdateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"DeleteNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"ReadNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
	}
	// Create
	t.Run("CreateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewNasClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudNasFileSystemCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("InvalidFileSystemStatus.Ordering")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAliCloudNasFileSystemCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAliCloudNasFileSystemCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId("MockFileSystemId")
	// Update
	t.Run("UpdateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewNasClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})

		err := resourceAliCloudNasFileSystemUpdate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateModifyFileSystemAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"description"} {
			switch p["alicloud_nas_file_system"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_nas_file_system"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAliCloudNasFileSystemUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateModifyFileSystemNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"description", "tags"} {
			switch p["alicloud_nas_file_system"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_nas_file_system"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAliCloudNasFileSystemUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("UpdateModifyFileSystemTagAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"tags"} {
			switch p["alicloud_nas_file_system"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_nas_file_system"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAliCloudNasFileSystemUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewNasClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudNasFileSystemDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAliCloudNasFileSystemDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		resourceData, _ := schema.InternalMap(p["alicloud_nas_file_system"].Schema).Data(nil, diff)
		resourceData.SetId(d.Id())
		_ = resourceData.Set("file_system_type", "cpfs")
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("InvalidFileSystemStatus.Ordering")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		patcheDescribeNasFileSystem := gomonkey.ApplyMethod(reflect.TypeOf(&NasService{}), "DescribeNasFileSystem", func(*NasService, string) (map[string]interface{}, error) {
			return responseMock["NotFoundError"]("ResourceNotfound")
		})
		err := resourceAliCloudNasFileSystemDelete(resourceData, rawClient)
		patches.Reset()
		patcheDescribeNasFileSystem.Reset()
		assert.Nil(t, err)
	})

	t.Run("DeleteIsExpectedErrors", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				return responseMock["RetryError"]("InvalidFileSystem.NotFound")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAliCloudNasFileSystemDelete(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})
	//Read
	t.Run("ReadDescribeNasFileSystemNotFound", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			NotFoundFlag := true
			noRetryFlag := false
			if NotFoundFlag {
				return responseMock["NotFoundError"]("ResourceNotfound")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NoRetryError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAliCloudNasFileSystemRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})

	t.Run("ReadDescribeNasFileSystemAbnormal", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			retryFlag := false
			noRetryFlag := true
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAliCloudNasFileSystemRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})
}

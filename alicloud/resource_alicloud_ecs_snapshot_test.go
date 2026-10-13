package alicloud

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers("alicloud_ecs_snapshot", &resource.Sweeper{
		Name: "alicloud_ecs_snapshot",
		F:    testSweepEcsSnapshots,
	})
}

func testSweepEcsSnapshots(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeSnapshots"

	request := map[string]interface{}{
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
		"RegionId":   client.RegionId,
	}

	var response map[string]interface{}

	for {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecs_snapshot", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Snapshots.Snapshot", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Snapshots.Snapshot", response)
		}

		result, _ := resp.([]interface{})

		for _, v := range result {
			item := v.(map[string]interface{})

			name := item["SnapshotName"]
			id := item["SnapshotId"]
			skip := true
			if !sweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(name.(string)), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
				if skip {
					log.Printf("[INFO] Skipping snapshot: %s (%s)", name, id)
					continue
				}
			}
			log.Printf("[INFO] Deleting snapshot: %s (%s)", name, id)
			action = "DeleteSnapshot"
			request := map[string]interface{}{
				"SnapshotId": item["SnapshotId"],
			}

			_, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)

			if err != nil {
				log.Printf("[ERROR] Failed to delete snapshot(%s (%s)): %s", name, id, err)
			}

			log.Printf("[INFO] Delete snapshot success: %s ", item["SnapshotId"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	return nil
}

func TestAccAliCloudECSSnapshot_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_snapshot.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsSnapshotMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsSnapshot")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secssnapshot%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsSnapshotBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_id": "${alicloud_ecs_disk_attachment.default.disk_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_id": CHECKSET,
						"status":  "accomplished",
					}),
					func(s *terraform.State) error {
						if value := s.RootModule().Resources[resourceId].Primary.Attributes["wait_until"]; value != "" {
							return fmt.Errorf("omitted wait_until must remain unset, got %q", value)
						}
						return nil
					},
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
					"snapshot_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snapshot_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Snapshot",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Snapshot",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "wait_until"},
			},
		},
	})
}

func TestAccAliCloudECSSnapshot_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_snapshot.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsSnapshotMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsSnapshot")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secssnapshot%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsSnapshotBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_id":           "${alicloud_ecs_disk_attachment.default.disk_id}",
					"category":          "standard",
					"wait_until":        "accomplished",
					"retention_days":    "50",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"snapshot_name":     name,
					"description":       name,
					"force":             "true",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Snapshot",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_id":           CHECKSET,
						"category":          "standard",
						"wait_until":        "accomplished",
						"status":            "accomplished",
						"retention_days":    "50",
						"resource_group_id": CHECKSET,
						"snapshot_name":     name,
						"description":       name,
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Snapshot",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "wait_until"},
			},
		},
	})
}

func TestAccAliCloudECSSnapshot_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_snapshot.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsSnapshotMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsSnapshot")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secssnapshot%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsSnapshotBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_id":        "${alicloud_ecs_disk_attachment.default.disk_id}",
					"retention_days": "50",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_id":        CHECKSET,
						"retention_days": "50",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention_days": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_days": "60",
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
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Snapshot",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Snapshot",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "wait_until"},
			},
		},
	})
}

func TestAccAliCloudECSSnapshot_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_snapshot.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsSnapshotMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsSnapshot")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secssnapshot%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsSnapshotBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_id":           "${alicloud_ecs_disk_attachment.default.disk_id}",
					"category":          "standard",
					"retention_days":    "50",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"name":              name,
					"description":       name,
					"force":             "true",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Snapshot",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_id":           CHECKSET,
						"category":          "standard",
						"retention_days":    "50",
						"resource_group_id": CHECKSET,
						"name":              name,
						"description":       name,
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Snapshot",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "wait_until"},
			},
		},
	})
}

var AliCloudEcsSnapshotMap0 = map[string]string{
	"available":   "true",
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
	"category":    CHECKSET,
	"status":      CHECKSET,
}

func AliCloudEcsSnapshotBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
  		status = "OK"
	}

	data "alicloud_zones" "default" {
  		available_disk_category     = "cloud_essd"
		available_instance_type     = data.alicloud_instance_types.default.instance_types.0.id
		available_resource_creation = "Instance"
	}
	
	data "alicloud_images" "default" {
  		most_recent = true
  		owners      = "system"
		architecture = "x86_64"
		os_type      = "linux"
	}
	
	data "alicloud_instance_types" "default" {
		instance_type_family = "ecs.g7"
		sorted_by            = "CPU"
  		image_id             = data.alicloud_images.default.images.0.id
  		system_disk_category = "cloud_essd"
	}
	
	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}
	
	resource "alicloud_vswitch" "default" {
		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.192.0/24"
  		zone_id      = data.alicloud_zones.default.zones.0.id
	}
	
	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}
	
	resource "alicloud_instance" "default" {
  		image_id                   = data.alicloud_images.default.images.0.id
  		instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  		security_groups            = alicloud_security_group.default.*.id
  		internet_charge_type       = "PayByTraffic"
  		internet_max_bandwidth_out = "10"
		availability_zone          = data.alicloud_zones.default.zones.0.id
  		instance_charge_type       = "PostPaid"
  		system_disk_category       = "cloud_essd"
		system_disk_encrypted      = true
  		vswitch_id                 = alicloud_vswitch.default.id
  		instance_name              = var.name
		data_disks {
			category = "cloud_essd"
			encrypted = true
			size     = 20
  		}
	}
	
	resource "alicloud_ecs_disk" "default" {
  		disk_name = var.name
		zone_id   = data.alicloud_zones.default.zones.0.id
  		category  = "cloud_essd"
		encrypted = true
		size      = 20
	}
	
	resource "alicloud_ecs_disk_attachment" "default" {
  		disk_id     = alicloud_ecs_disk.default.id
  		instance_id = alicloud_instance.default.id
	}
`, name)
}

func AliCloudEcsSnapshotBasicDependence1(name string) string {
	return AliCloudEcsSnapshotBasicDependence0(name)
}

// lintignore: R001
func TestUnitAliCloudEcsSnapshot(t *testing.T) {
	p := Provider().ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_ecs_snapshot"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_ecs_snapshot"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"disk_id":           "disk_id",
		"category":          `standard`,
		"name":              "name",
		"description":       "description",
		"resource_group_id": "resource_group_id",
		"tags": map[string]string{
			"Created": "TF",
			"For":     "Test",
		},
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
		"Snapshots": map[string]interface{}{
			"Snapshot": []interface{}{
				map[string]interface{}{
					"Category":                   "standard",
					"Description":                "description",
					"SourceDiskId":               "disk_id",
					"InstantAccess":              true,
					"InstantAccessRetentionDays": 20,
					"ResourceGroupId":            "resource_group_id",
					"RetentionDays":              20,
					"SnapshotName":               "snapshot_name",
					"Name":                       "snapshot_name",
					"Status":                     "accomplished",
					"Available":                  true,
					"SnapshotId":                 "MockSnapshotId",
				},
			},
		},
		"TagResources": map[string]interface{}{
			"TagResource": []interface{}{
				map[string]interface{}{
					"TagKey":   "Created",
					"TagValue": "TF",
				},
				map[string]interface{}{
					"TagKey":   "For",
					"TagValue": "Test",
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
			result := ReadMockResponse
			return result, nil
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
			result["SnapshotId"] = "MockSnapshotId"
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
		"ReadDescribeEcsSnapshotNotFound": func(errorCode string) (map[string]interface{}, error) {
			result := map[string]interface{}{
				"Snapshots": map[string]interface{}{
					"Snapshot": []interface{}{},
				},
			}
			return result, nil
		},
	}
	// Create
	t.Run("CreateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEcsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudEcsSnapshotCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateAbnormal", func(t *testing.T) {
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
			return responseMock["CreateNormal"]("")
		})
		err := resourceAliCloudEcsSnapshotCreate(d, rawClient)
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
		err := resourceAliCloudEcsSnapshotCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("CreateNonRetryableError", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := true
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
		err := resourceAliCloudEcsSnapshotCreate(dCreate, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId("MockSnapshotId")
	// Update
	t.Run("UpdateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEcsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})

		err := resourceAliCloudEcsSnapshotUpdate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateHpcClusterAttributeAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"description", "snapshot_name", "name"} {
			switch p["alicloud_ecs_snapshot"].Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &terraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &terraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &terraform.ResourceAttrDiff{Old: "", New: "TF"}
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_ecs_snapshot"].Schema).Data(nil, diff)
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
		err := resourceAliCloudEcsSnapshotUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateHpcClusterAttributeTagsAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"description", "snapshot_name", "name", "tags"} {
			switch p["alicloud_ecs_snapshot"].Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &terraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &terraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &terraform.ResourceAttrDiff{Old: "", New: "TF"}
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_ecs_snapshot"].Schema).Data(nil, diff)
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
		err := resourceAliCloudEcsSnapshotUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateHpcClusterAttributeNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"description", "snapshot_name", "name"} {
			switch p["alicloud_ecs_snapshot"].Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &terraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &terraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &terraform.ResourceAttrDiff{Old: "", New: "TF"}
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_ecs_snapshot"].Schema).Data(nil, diff)
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
		err := resourceAliCloudEcsSnapshotUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEcsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudEcsSnapshotDelete(d, rawClient)
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
		err := resourceAliCloudEcsSnapshotDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockNormal", func(t *testing.T) {
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
			return responseMock["DeleteNormal"]("")
		})
		patcheDescribeEcsSnapshot := gomonkey.ApplyMethod(reflect.TypeOf(&EcsService{}), "DescribeEcsSnapshot", func(*EcsService, string) (map[string]interface{}, error) {
			return responseMock["RetryError"]("InvalidFilterKey.NotFound")
		})
		err := resourceAliCloudEcsSnapshotDelete(d, rawClient)
		patches.Reset()
		patcheDescribeEcsSnapshot.Reset()
		assert.NotNil(t, err)
	})

	//Read
	t.Run("ReadDescribeEcsSnapshotNotFound", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			NotFoundFlag := true
			noRetryFlag := false
			if NotFoundFlag {
				return responseMock["ReadDescribeEcsSnapshotNotFound"]("")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NoRetryError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAliCloudEcsSnapshotRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})

	t.Run("ReadDescribeEcsHpcClusterAbnormal", func(t *testing.T) {
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
		err := resourceAliCloudEcsSnapshotRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})
}

func TestAccAliCloudECSSnapshot_availableDisk(t *testing.T) {
	var snapshot, disk map[string]interface{}
	snapshotID := "alicloud_ecs_snapshot.default"
	diskID := "alicloud_ecs_disk.restored"
	snapshotCheck := resourceCheckInitWithDescribeMethod(snapshotID, &snapshot, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsSnapshot")
	diskCheck := resourceCheckInitWithDescribeMethod(diskID, &disk, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsDisk")
	name := fmt.Sprintf("tf-testacc-snapshot-%d", acctest.RandIntRange(10000, 99999))
	config := AliCloudEcsSnapshotBasicDependence0(name) + `
resource "alicloud_ecs_snapshot" "default" {
  disk_id       = alicloud_ecs_disk_attachment.default.disk_id
  snapshot_name = var.name
  force         = true
  wait_until    = "available"
}

resource "alicloud_ecs_disk" "restored" {
  disk_name   = "${var.name}-restored"
  zone_id     = alicloud_ecs_disk.default.zone_id
  category    = "cloud_essd"
  snapshot_id = alicloud_ecs_snapshot.default.id
  encrypted   = true
  size        = 20
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			snapshotCheck.checkResourceDestroy(), diskCheck.checkResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					snapshotCheck.checkResourceExists(),
					resource.TestCheckResourceAttr(snapshotID, "wait_until", "available"),
					resource.TestCheckResourceAttr(snapshotID, "available", "true"),
					diskCheck.checkResourceExists(),
					resource.TestCheckResourceAttrPair(diskID, "snapshot_id", snapshotID, "id"),
					func(s *terraform.State) error {
						service := EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
						restored, err := service.DescribeEcsDisk(s.RootModule().Resources[diskID].Primary.ID)
						if err != nil {
							return err
						}
						if fmt.Sprint(restored["SourceSnapshotId"]) != s.RootModule().Resources[snapshotID].Primary.ID {
							return fmt.Errorf("restored disk must use the created snapshot")
						}
						return nil
					},
				),
			},
			{
				ResourceName: snapshotID, ImportState: true, ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"force", "wait_until"},
			},
		},
	})
}

func ecsSnapshotTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		for key := range r.Form {
			if strings.EqualFold(strings.ReplaceAll(key, "_", ""), "WaitUntil") {
				t.Errorf("local wait_until must not be sent to ECS: %s", key)
			}
			if strings.EqualFold(key, "Available") {
				t.Errorf("computed available must not be sent to ECS: %s", key)
			}
		}
		handler(w, r)
	}))
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("ecs", endpoint)
	return client
}

func TestUnitEcsSnapshotCreateAvailable(t *testing.T) {
	for _, tc := range []struct {
		name        string
		available   interface{}
		status      string
		progress    string
		pending     bool
		notFound    bool
		apiError    bool
		wantError   string
		wantQueries int32
	}{
		{name: "uploading_available", available: true, status: "progressing", progress: "40%", wantQueries: 2},
		{name: "becomes_available", available: true, status: "progressing", progress: "40%", pending: true, wantQueries: 3},
		{name: "not_yet_available", available: false, status: "accomplished", progress: "100%", wantError: "timeout while waiting"},
		{name: "missing_available", status: "accomplished", progress: "100%", wantError: "timeout while waiting"},
		{name: "eventually_visible", available: true, status: "progressing", progress: "40%", notFound: true, wantQueries: 3},
		{name: "query_error", apiError: true, wantError: "Forbidden", wantQueries: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var describes int32
			client := ecsSnapshotTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Error(err)
				}
				w.Header().Set("Content-Type", "application/json")
				var response interface{}
				switch r.Form.Get("Action") {
				case "CreateSnapshot":
					response = map[string]interface{}{"SnapshotId": "snapshot-test"}
				case "DescribeSnapshots":
					count := atomic.AddInt32(&describes, 1)
					if r.Form.Get("Status") != "" || r.Form.Get("SnapshotIds") != `["snapshot-test"]` {
						t.Error("DescribeSnapshots must query the snapshot ID without filtering status")
					}
					if tc.apiError {
						w.WriteHeader(http.StatusForbidden)
						response = map[string]interface{}{"Code": "Forbidden", "Message": "test query failure"}
						break
					}
					snapshot := map[string]interface{}{"SnapshotId": "snapshot-test", "Status": tc.status, "Progress": tc.progress}
					if tc.available != nil {
						snapshot["Available"] = tc.available
					}
					if tc.pending && count == 1 {
						snapshot["Available"] = false
					}
					snapshots := []interface{}{snapshot}
					if tc.notFound && count == 1 {
						snapshots = []interface{}{}
					}
					response = map[string]interface{}{"Snapshots": map[string]interface{}{"Snapshot": snapshots}}
				default:
					t.Errorf("unexpected action: %s", r.Form.Get("Action"))
					w.WriteHeader(http.StatusBadRequest)
					response = map[string]interface{}{"Code": "UnexpectedAction"}
				}
				if err := json.NewEncoder(w).Encode(response); err != nil {
					t.Error(err)
				}
			})
			r := resourceAliCloudEcsSnapshot()
			timeout := 6 * time.Second
			if tc.pending || tc.notFound {
				timeout = 10 * time.Second
			}
			r.Timeouts = &schema.ResourceTimeout{Create: schema.DefaultTimeout(timeout)}
			d := r.Data(nil)
			d.MarkNewResource()
			if err := d.Set("disk_id", "disk-test"); err != nil {
				t.Fatal(err)
			}
			if err := d.Set("wait_until", "available"); err != nil {
				t.Fatal(err)
			}
			err := r.Create(d, client)
			if tc.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantError) {
					t.Fatalf("expected %q, got %v", tc.wantError, err)
				}
			} else if err != nil {
				t.Fatalf("available snapshot must finish creation while uploading: %v", err)
			}
			if d.Id() != "snapshot-test" {
				t.Fatalf("snapshot ID must be retained: %q", d.Id())
			}
			if tc.wantError == "" && d.Get("status") != tc.status {
				t.Fatalf("Read must preserve the actual status: %v", d.Get("status"))
			}
			if got := atomic.LoadInt32(&describes); tc.wantQueries > 0 && got != tc.wantQueries {
				t.Fatalf("expected %d DescribeSnapshots calls including the final Read, got %d", tc.wantQueries, got)
			}
		})
	}
}

func TestUnitEcsSnapshotCategory(t *testing.T) {
	for _, category := range []string{"standard", "flash"} {
		t.Run(category, func(t *testing.T) {
			var creates int32
			client := ecsSnapshotTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Error(err)
				}
				w.Header().Set("Content-Type", "application/json")
				response := map[string]interface{}{}
				switch r.Form.Get("Action") {
				case "CreateSnapshot":
					atomic.AddInt32(&creates, 1)
					if got := r.Form.Get("Category"); got != category {
						t.Errorf("expected CreateSnapshot category %q, got %q", category, got)
					}
					response["SnapshotId"] = "snapshot-test"
				case "DescribeSnapshots":
					response = ecsSnapshotTestDescribeResponse("accomplished", true, map[string]interface{}{
						"Category": category,
					})
				default:
					t.Errorf("unexpected action: %s", r.Form.Get("Action"))
					w.WriteHeader(http.StatusBadRequest)
					response["Code"] = "UnexpectedAction"
				}
				if err := json.NewEncoder(w).Encode(response); err != nil {
					t.Error(err)
				}
			})

			r := resourceAliCloudEcsSnapshot()
			r.Timeouts = &schema.ResourceTimeout{Create: schema.DefaultTimeout(10 * time.Second)}
			d := r.Data(nil)
			if err := d.Set("disk_id", "disk-test"); err != nil {
				t.Fatal(err)
			}
			if err := d.Set("category", category); err != nil {
				t.Fatal(err)
			}
			d.MarkNewResource()
			if got := d.Timeout(schema.TimeoutCreate); got != 10*time.Second {
				t.Fatalf("expected a 10s create timeout, got %s", got)
			}
			if err := r.Create(d, client); err != nil {
				t.Fatal(err)
			}
			if got := atomic.LoadInt32(&creates); got != 1 {
				t.Fatalf("expected one CreateSnapshot call, got %d", got)
			}
			if err := d.Set("category", ""); err != nil {
				t.Fatal(err)
			}
			if err := r.Read(d, client); err != nil {
				t.Fatal(err)
			}
			if got := d.Get("category"); got != category {
				t.Fatalf("Read must restore category %q, got %v", category, got)
			}
		})
	}
}

func TestUnitEcsSnapshotUpdateAvailable(t *testing.T) {
	t.Run("description while uploading", func(t *testing.T) {
		var describes int32
		var modifies int32
		var firstAction atomic.Value
		client := ecsSnapshotTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if err := r.ParseForm(); err != nil {
				t.Error(err)
			}
			w.Header().Set("Content-Type", "application/json")
			response := map[string]interface{}{}
			if firstAction.Load() == nil {
				firstAction.Store(r.Form.Get("Action"))
			}
			switch r.Form.Get("Action") {
			case "DescribeSnapshots":
				atomic.AddInt32(&describes, 1)
				response = ecsSnapshotTestDescribeResponse("progressing", true, map[string]interface{}{
					"Description": "updated description",
				})
			case "ModifySnapshotAttribute":
				if r.Form.Get("Description") != "updated description" {
					t.Errorf("unexpected description: %q", r.Form.Get("Description"))
				}
				atomic.AddInt32(&modifies, 1)
			default:
				t.Errorf("unexpected action: %s", r.Form.Get("Action"))
				w.WriteHeader(http.StatusBadRequest)
				response = map[string]interface{}{"Code": "UnexpectedAction"}
			}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Error(err)
			}
		})

		d := schema.TestResourceDataRaw(t, resourceAliCloudEcsSnapshot().Schema, map[string]interface{}{
			"disk_id":     "disk-test",
			"description": "updated description",
			"wait_until":  "available",
		})
		d.SetId("snapshot-test")
		if err := resourceAliCloudEcsSnapshotUpdate(d, client); err != nil {
			t.Fatalf("description update must finish while the snapshot is uploading: %v", err)
		}
		if got := atomic.LoadInt32(&modifies); got != 1 {
			t.Fatalf("expected one ModifySnapshotAttribute call, got %d", got)
		}
		if got := atomic.LoadInt32(&describes); got != 1 {
			t.Fatalf("expected only the final Read after the description update, got %d queries", got)
		}
		if got, _ := firstAction.Load().(string); got != "ModifySnapshotAttribute" {
			t.Fatalf("ordinary metadata updates must not wait for upload completion before the API call, first action was %q", got)
		}
	})

	t.Run("retention waits until accomplished", func(t *testing.T) {
		var describes int32
		var lastStatus atomic.Value
		var modifies int32
		client := ecsSnapshotTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if err := r.ParseForm(); err != nil {
				t.Error(err)
			}
			w.Header().Set("Content-Type", "application/json")
			response := map[string]interface{}{}
			switch r.Form.Get("Action") {
			case "DescribeSnapshots":
				count := atomic.AddInt32(&describes, 1)
				status := "accomplished"
				if count == 1 {
					status = "progressing"
				}
				lastStatus.Store(status)
				response = ecsSnapshotTestDescribeResponse(status, true, map[string]interface{}{
					"RetentionDays": 30,
				})
			case "ModifySnapshotAttribute":
				if status, _ := lastStatus.Load().(string); status != "accomplished" {
					w.WriteHeader(http.StatusForbidden)
					response = map[string]interface{}{"Code": "InvalidSnapshotId.NotReady", "Message": "snapshot is not ready"}
					break
				}
				if r.Form.Get("RetentionDays") != "30" {
					t.Errorf("unexpected retention days: %q", r.Form.Get("RetentionDays"))
				}
				atomic.AddInt32(&modifies, 1)
			default:
				t.Errorf("unexpected action: %s", r.Form.Get("Action"))
				w.WriteHeader(http.StatusBadRequest)
				response = map[string]interface{}{"Code": "UnexpectedAction"}
			}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Error(err)
			}
		})

		d := schema.TestResourceDataRaw(t, resourceAliCloudEcsSnapshot().Schema, map[string]interface{}{
			"disk_id":        "disk-test",
			"retention_days": 30,
		})
		d.SetId("snapshot-test")
		if err := resourceAliCloudEcsSnapshotUpdate(d, client); err != nil {
			t.Fatalf("retention update must wait for upload completion: %v", err)
		}
		if got := atomic.LoadInt32(&modifies); got != 1 {
			t.Fatalf("expected one successful ModifySnapshotAttribute call, got %d", got)
		}
		if got := atomic.LoadInt32(&describes); got < 2 {
			t.Fatalf("expected the retention update to observe progressing then accomplished, got %d queries", got)
		}
	})

	t.Run("resource group while uploading", func(t *testing.T) {
		var describes int32
		var joins int32
		client := ecsSnapshotTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if err := r.ParseForm(); err != nil {
				t.Error(err)
			}
			w.Header().Set("Content-Type", "application/json")
			response := map[string]interface{}{}
			switch r.Form.Get("Action") {
			case "DescribeSnapshots":
				atomic.AddInt32(&describes, 1)
				if atomic.LoadInt32(&joins) != 1 {
					t.Error("resource-group updates must call JoinResourceGroup before the final Read")
				}
				response = ecsSnapshotTestDescribeResponse("progressing", true, map[string]interface{}{
					"ResourceGroupId": "rg-test",
				})
			case "JoinResourceGroup":
				if atomic.LoadInt32(&describes) != 0 {
					t.Error("resource-group updates must not wait before JoinResourceGroup")
				}
				if r.Form.Get("ResourceGroupId") != "rg-test" {
					t.Errorf("unexpected resource group: %q", r.Form.Get("ResourceGroupId"))
				}
				atomic.AddInt32(&joins, 1)
			default:
				t.Errorf("unexpected action: %s", r.Form.Get("Action"))
				w.WriteHeader(http.StatusBadRequest)
				response = map[string]interface{}{"Code": "UnexpectedAction"}
			}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Error(err)
			}
		})

		d := schema.TestResourceDataRaw(t, resourceAliCloudEcsSnapshot().Schema, map[string]interface{}{
			"disk_id":           "disk-test",
			"resource_group_id": "rg-test",
		})
		d.SetId("snapshot-test")
		if err := resourceAliCloudEcsSnapshotUpdate(d, client); err != nil {
			t.Fatalf("resource group update must finish while the snapshot is uploading: %v", err)
		}
		if got := atomic.LoadInt32(&joins); got != 1 {
			t.Fatalf("expected one JoinResourceGroup call, got %d", got)
		}
		if got := atomic.LoadInt32(&describes); got != 1 {
			t.Fatalf("expected only the final Read, got %d DescribeSnapshots calls", got)
		}
		if d.Get("resource_group_id") != "rg-test" || d.Get("status") != "progressing" {
			t.Fatal("Read must retain the resource group and actual snapshot status")
		}
	})
}

func ecsSnapshotTestDescribeResponse(status string, available bool, fields map[string]interface{}) map[string]interface{} {
	snapshot := map[string]interface{}{
		"SnapshotId":   "snapshot-test",
		"SourceDiskId": "disk-test",
		"Status":       status,
		"Available":    available,
	}
	for key, value := range fields {
		snapshot[key] = value
	}
	return map[string]interface{}{
		"Snapshots": map[string]interface{}{
			"Snapshot": []interface{}{snapshot},
		},
	}
}

func TestUnitEcsSnapshotWaitUntilSchema(t *testing.T) {
	r := resourceAliCloudEcsSnapshot()
	s := r.Schema["wait_until"]
	if s == nil {
		t.Fatal("wait_until must be an operation-local argument")
	}
	if s.Type != schema.TypeString || !s.Optional || s.Computed || s.ForceNew || s.Default != nil || s.DefaultFunc != nil {
		t.Fatalf("unexpected wait_until schema: %#v", s)
	}
	available := r.Schema["available"]
	if available == nil || available.Type != schema.TypeBool || !available.Computed || available.Optional || available.Required || available.ForceNew {
		t.Fatalf("available must be a computed-only boolean, got %#v", available)
	}
	if errs := r.Validate(terraform.NewResourceConfigRaw(map[string]interface{}{"disk_id": "disk-test"})); len(errs) != 0 {
		t.Fatalf("omitted wait_until must be valid: %v", errs)
	}
	for _, value := range []string{"accomplished", "available", "progressing", "AVAILABLE", ""} {
		errs := r.Validate(terraform.NewResourceConfigRaw(map[string]interface{}{"disk_id": "disk-test", "wait_until": value}))
		valid := value == "accomplished" || value == "available"
		if valid != (len(errs) == 0) {
			t.Errorf("wait_until %q validation errors: %v", value, errs)
		}
	}
}

func TestUnitEcsSnapshotCreateWaitUntil(t *testing.T) {
	for _, tc := range []struct {
		name       string
		policy     string
		inProgress bool
		wantError  bool
	}{
		{name: "default_does_not_return_available", inProgress: true, wantError: true},
		{name: "accomplished_does_not_return_available", policy: "accomplished", inProgress: true, wantError: true},
		{name: "default_reaches_accomplished"},
		{name: "explicit_reaches_accomplished", policy: "accomplished"},
		{name: "available_returns_while_progressing", policy: "available", inProgress: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var describes int32
			client := ecsSnapshotTestClient(t, func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				response := map[string]interface{}{}
				switch req.Form.Get("Action") {
				case "CreateSnapshot":
					response["SnapshotId"] = "snapshot-test"
				case "DescribeSnapshots":
					count := atomic.AddInt32(&describes, 1)
					status := "progressing"
					if !tc.inProgress && count > 1 {
						status = "accomplished"
					}
					response = ecsSnapshotTestDescribeResponse(status, true, nil)
				default:
					t.Errorf("unexpected action: %s", req.Form.Get("Action"))
				}
				if err := json.NewEncoder(w).Encode(response); err != nil {
					t.Error(err)
				}
			})
			r := resourceAliCloudEcsSnapshot()
			timeout := 10 * time.Second
			if tc.wantError {
				timeout = 6 * time.Second
			}
			r.Timeouts = &schema.ResourceTimeout{Create: schema.DefaultTimeout(timeout)}
			config := map[string]interface{}{"disk_id": "disk-test"}
			if tc.policy != "" {
				config["wait_until"] = tc.policy
			}
			diff, err := r.Diff(context.Background(), nil, terraform.NewResourceConfigRaw(config), client)
			if err != nil {
				t.Fatal(err)
			}
			state, diags := r.Apply(context.Background(), nil, diff, client)
			if tc.wantError {
				if !diags.HasError() || !strings.Contains(fmt.Sprintf("%v", diags), "timeout while waiting") {
					t.Fatalf("must keep waiting while Available=true and Status=progressing, got %v", diags)
				}
				return
			}
			if diags.HasError() {
				t.Fatalf("%v", diags)
			}
			wantPolicy, wantStatus, wantQueries := tc.policy, "accomplished", int32(3)
			if tc.policy == "available" {
				wantPolicy, wantStatus, wantQueries = "available", "progressing", 2
			}
			if state.ID != "snapshot-test" || state.Attributes["wait_until"] != wantPolicy || state.Attributes["status"] != wantStatus {
				t.Fatalf("unexpected state: %#v", state)
			}
			if state.Attributes["available"] != "true" {
				t.Fatalf("Read must preserve actual available value, got %q", state.Attributes["available"])
			}
			if got := atomic.LoadInt32(&describes); got != wantQueries {
				t.Fatalf("expected %d queries including Read, got %d", wantQueries, got)
			}
		})
	}
}

func TestUnitEcsSnapshotUpdateWaitUntil(t *testing.T) {
	for _, tc := range []struct {
		name   string
		policy string
		field  string
	}{
		{name: "default_name", field: "snapshot_name"},
		{name: "accomplished_name", policy: "accomplished", field: "snapshot_name"},
		{name: "available_name", policy: "available", field: "snapshot_name"},
		{name: "default_legacy_name", field: "name"},
		{name: "available_legacy_name", policy: "available", field: "name"},
		{name: "default_description", field: "description"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var describes, modifies int32
			client := ecsSnapshotTestClient(t, func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				response := map[string]interface{}{}
				switch req.Form.Get("Action") {
				case "ModifySnapshotAttribute":
					if atomic.LoadInt32(&describes) != 0 {
						t.Error("ordinary metadata must be modified without a readiness precheck")
					}
					key := "SnapshotName"
					if tc.field == "description" {
						key = "Description"
					}
					if req.Form.Get(key) != "after" {
						t.Errorf("expected %s=after, got %q", key, req.Form.Get(key))
					}
					atomic.AddInt32(&modifies, 1)
				case "DescribeSnapshots":
					atomic.AddInt32(&describes, 1)
					if atomic.LoadInt32(&modifies) != 1 {
						t.Error("Read must follow the metadata update")
					}
					response = ecsSnapshotTestDescribeResponse("progressing", true, map[string]interface{}{
						"SnapshotName": "after", "Description": "after",
					})
				default:
					t.Errorf("unexpected action: %s", req.Form.Get("Action"))
				}
				if err := json.NewEncoder(w).Encode(response); err != nil {
					t.Error(err)
				}
			})
			r := resourceAliCloudEcsSnapshot()
			r.Timeouts = &schema.ResourceTimeout{Update: schema.DefaultTimeout(10 * time.Second)}
			state := &terraform.InstanceState{ID: "snapshot-test", Attributes: map[string]string{
				"disk_id": "disk-test", "snapshot_name": "before", "name": "before", "description": "before", "wait_until": "accomplished",
			}}
			config := map[string]interface{}{"disk_id": "disk-test", "description": "before"}
			if tc.field != "name" {
				config["snapshot_name"] = "before"
			}
			config[tc.field] = "after"
			if tc.policy != "" {
				config["wait_until"] = tc.policy
			}
			diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(config), client)
			if err != nil {
				t.Fatal(err)
			}
			updated, diags := r.Apply(context.Background(), state, diff, client)
			if diags.HasError() {
				t.Fatalf("%v", diags)
			}
			if got := atomic.LoadInt32(&describes); got != 1 || atomic.LoadInt32(&modifies) != 1 {
				t.Fatalf("expected one modify and only the final Read, got modify=%d query=%d", modifies, got)
			}
			if updated.Attributes["status"] != "progressing" {
				t.Fatalf("Read must preserve actual progressing status, got %q", updated.Attributes["status"])
			}
			if updated.Attributes[tc.field] != "after" {
				t.Fatalf("Read must refresh %s, got %q", tc.field, updated.Attributes[tc.field])
			}
			if tc.field != "description" && (updated.Attributes["snapshot_name"] != "after" || updated.Attributes["name"] != "after") {
				t.Fatalf("Read must refresh both name aliases, got %#v", updated.Attributes)
			}
		})
	}
}

func TestUnitEcsSnapshotRetentionWaitUntil(t *testing.T) {
	for _, policy := range []string{"accomplished", "available"} {
		for _, mixed := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/mixed_%t", policy, mixed), func(t *testing.T) {
				var describes, modifies int32
				client := ecsSnapshotTestClient(t, func(w http.ResponseWriter, req *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					response := map[string]interface{}{}
					switch req.Form.Get("Action") {
					case "DescribeSnapshots":
						count := atomic.AddInt32(&describes, 1)
						status := "accomplished"
						if count == 1 {
							status = "progressing"
						}
						response = ecsSnapshotTestDescribeResponse(status, true, map[string]interface{}{"RetentionDays": 30})
					case "ModifySnapshotAttribute":
						if atomic.LoadInt32(&describes) < 2 {
							t.Error("retention must observe accomplished before modification in every mode")
						}
						if req.Form.Get("RetentionDays") != "30" || mixed && req.Form.Get("Description") != "after" {
							t.Error("missing retention or mixed metadata update")
						}
						atomic.AddInt32(&modifies, 1)
					default:
						t.Errorf("unexpected action: %s", req.Form.Get("Action"))
					}
					if err := json.NewEncoder(w).Encode(response); err != nil {
						t.Error(err)
					}
				})
				r := resourceAliCloudEcsSnapshot()
				r.Timeouts = &schema.ResourceTimeout{Update: schema.DefaultTimeout(10 * time.Second)}
				state := &terraform.InstanceState{ID: "snapshot-test", Attributes: map[string]string{
					"disk_id": "disk-test", "retention_days": "20", "wait_until": policy,
				}}
				config := map[string]interface{}{"disk_id": "disk-test", "retention_days": 30, "wait_until": policy}
				if mixed {
					config["description"] = "after"
				}
				diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(config), client)
				if err != nil {
					t.Fatal(err)
				}
				if _, diags := r.Apply(context.Background(), state, diff, client); diags.HasError() {
					t.Fatalf("%v", diags)
				}
				if atomic.LoadInt32(&modifies) != 1 || atomic.LoadInt32(&describes) != 3 {
					t.Fatalf("expected progressing/accomplished pre-wait, one modify and final Read, got modify=%d queries=%d", modifies, describes)
				}
			})
		}
	}
}

func TestUnitEcsSnapshotWaitUntilOnly(t *testing.T) {
	for _, tc := range []struct{ name, before, after string }{
		{"to_available", "accomplished", "available"},
		{"to_accomplished", "available", "accomplished"},
		{"remove_override", "available", ""},
		{"remove_accomplished", "accomplished", ""},
		{"upgrade_legacy_state", "", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var describes int32
			client := ecsSnapshotTestClient(t, func(w http.ResponseWriter, req *http.Request) {
				if req.Form.Get("Action") != "DescribeSnapshots" {
					t.Errorf("policy-only change must not mutate ECS: %s", req.Form.Get("Action"))
				}
				atomic.AddInt32(&describes, 1)
				w.Header().Set("Content-Type", "application/json")
				response := ecsSnapshotTestDescribeResponse("progressing", false, map[string]interface{}{
					"SnapshotName": "unchanged", "Description": "unchanged", "WaitUntil": "cloud-value-must-be-ignored",
				})
				if err := json.NewEncoder(w).Encode(response); err != nil {
					t.Error(err)
				}
			})
			r := resourceAliCloudEcsSnapshot()
			r.Timeouts = &schema.ResourceTimeout{Update: schema.DefaultTimeout(1 * time.Second)}
			state := &terraform.InstanceState{ID: "snapshot-test", Attributes: map[string]string{
				"disk_id": "disk-test", "snapshot_name": "unchanged", "name": "unchanged", "description": "unchanged", "status": "progressing",
			}}
			if tc.before != "" {
				state.Attributes["wait_until"] = tc.before
			}
			state, diags := r.RefreshWithoutUpgrade(context.Background(), state, client)
			if diags.HasError() {
				t.Fatalf("%v", diags)
			}
			if state.Attributes["wait_until"] != tc.before {
				t.Fatalf("Refresh must preserve the local policy, got %q", state.Attributes["wait_until"])
			}
			atomic.StoreInt32(&describes, 0)
			config := map[string]interface{}{"disk_id": "disk-test", "snapshot_name": "unchanged", "description": "unchanged"}
			if tc.after != "" {
				config["wait_until"] = tc.after
			}
			diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(config), client)
			if err != nil {
				t.Fatal(err)
			}
			if tc.before == "" && tc.after == "" {
				if diff != nil && !diff.Empty() {
					t.Fatalf("refreshing legacy state with no waiting policy must not create a resource diff: %#v", diff)
				}
				return
			}
			if diff == nil || diff.RequiresNew() || len(diff.Attributes) != 1 || diff.Attributes["wait_until"] == nil {
				t.Fatalf("expected only a non-replacement wait_until diff, got %#v", diff)
			}
			updated, diags := r.Apply(context.Background(), state, diff, client)
			if diags.HasError() {
				t.Fatalf("%v", diags)
			}
			if updated.ID != state.ID || updated.Attributes["wait_until"] != tc.after || updated.Attributes["status"] != "progressing" || updated.Attributes["available"] != "false" {
				t.Fatalf("Read must preserve the local policy and actual status without recreation: %#v", updated)
			}
			if got := atomic.LoadInt32(&describes); got != 1 {
				t.Fatalf("policy-only changes must perform only the final Read, not a waiter, got %d queries", got)
			}
			diff, err = r.Diff(context.Background(), updated, terraform.NewResourceConfigRaw(config), client)
			if err != nil || diff != nil && !diff.Empty() {
				t.Fatalf("applied local policy must have no remaining diff: %#v, %v", diff, err)
			}
		})
	}
}

func TestUnitEcsSnapshotReadAvailable(t *testing.T) {
	values := []interface{}{true, false, true, "missing", true, nil}
	var describes int32
	client := ecsSnapshotTestClient(t, func(w http.ResponseWriter, req *http.Request) {
		if req.Form.Get("Action") != "DescribeSnapshots" {
			t.Errorf("unexpected action: %s", req.Form.Get("Action"))
		}
		count := atomic.AddInt32(&describes, 1)
		response := ecsSnapshotTestDescribeResponse("progressing", false, map[string]interface{}{
			"Available": values[count-1], "WaitUntil": "cloud-value-must-be-ignored",
		})
		if values[count-1] == "missing" {
			snapshot := response["Snapshots"].(map[string]interface{})["Snapshot"].([]interface{})[0].(map[string]interface{})
			delete(snapshot, "Available")
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Error(err)
		}
	})
	r := resourceAliCloudEcsSnapshot()
	d := r.Data(nil)
	d.SetId("snapshot-test")
	imported, err := r.Importer.State(d, client)
	if err != nil || len(imported) != 1 {
		t.Fatalf("unexpected import result: %#v, %v", imported, err)
	}
	state := imported[0].State()
	for _, value := range values {
		newState, diags := r.RefreshWithoutUpgrade(context.Background(), state, client)
		if diags.HasError() {
			t.Fatalf("%v", diags)
		}
		state = newState
		if state.Attributes["available"] != strconv.FormatBool(value == true) || state.Attributes["status"] != "progressing" {
			t.Fatalf("Read must replace availability with the reported boolean or false when absent, response=%v state=%#v", value, state.Attributes)
		}
		if got := r.Data(state).Get("available"); got != (value == true) {
			t.Fatalf("available must remain a boolean in ResourceData, got %#v", got)
		}
		if state.Attributes["wait_until"] != "" {
			t.Fatalf("import/Read must not invent a local policy: %q", state.Attributes["wait_until"])
		}
	}
}

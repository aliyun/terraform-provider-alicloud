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
		"alicloud_ros_template_scratch",
		&resource.Sweeper{
			Name: "alicloud_ros_template_scratch",
			F:    testSweepRosTemplateScratch,
		})
}

func testSweepRosTemplateScratch(region string) error {
	if testSweepPreCheckWithRegions(region, true, connectivity.ROSSupportRegions) {
		log.Printf("[INFO] Skipping Ros Template Scratch unsupported region: %s", region)
		return nil
	}

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListTemplateScratches"
	request := map[string]interface{}{}
	request["RegionId"] = client.RegionId

	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	conn, err := client.NewRosClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
		return nil
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &runtime)
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

		resp, err := jsonpath.Get("$.TemplateScratches", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.TemplateScratches", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			name := fmt.Sprint(item["Description"])
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Ros Template Scratch: %s", name)
				continue
			}
			action := "DeleteTemplateScratch"
			request := map[string]interface{}{
				"TemplateScratchId": item["TemplateScratchId"],
			}
			request["RegionId"] = client.RegionId
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Ros Template Scratch (%s): %s", item["Description"].(string), err)
			}
			log.Printf("[INFO] Delete Ros Template Scratch success: %s ", item["Description"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudROSTemplateScratch_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ros_template_scratch.default"
	checkoutSupportedRegions(t, true, connectivity.ROSSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudROSTemplateScratchMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRosTemplateScratch")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srostemplatescratch%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudROSTemplateScratchBasicDependence0)
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
					"description": "${var.name}",
					"preference_parameters": []map[string]interface{}{
						{
							"parameter_key":   "DeletionPolicy",
							"parameter_value": "Retain",
						},
					},
					"source_resource_group": []map[string]interface{}{
						{
							"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
							"resource_type_filter": []string{"ALIYUN::ECS::VPC", "ALIYUN::ECS::SecurityGroup"},
						},
					},
					"template_scratch_type": "ResourceImport",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":             name,
						"preference_parameters.#": "1",
						"source_resource_group.#": "1",
						"template_scratch_type":   "ResourceImport",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preference_parameters": []map[string]interface{}{
						{
							"parameter_key":   "DeletionPolicy",
							"parameter_value": "Delete",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preference_parameters.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_resource_group": []map[string]interface{}{
						{
							"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
							"resource_type_filter": []string{"ALIYUN::ECS::SecurityGroup"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_resource_group.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"execution_mode"},
			},
		},
	})
}
func TestAccAlicloudROSTemplateScratch_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ros_template_scratch.default"
	checkoutSupportedRegions(t, true, connectivity.ROSSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudROSTemplateScratchMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRosTemplateScratch")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srostemplatescratch%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudROSTemplateScratchBasicDependence1)
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
					"description": "${var.name}",
					"preference_parameters": []map[string]interface{}{
						{
							"parameter_key":   "DeletionPolicy",
							"parameter_value": "Retain",
						},
					},
					"template_scratch_type": "ResourceImport",
					"source_resources": []map[string]interface{}{
						{
							"resource_id":   "${data.alicloud_vpcs.default.ids.0}",
							"resource_type": "ALIYUN::ECS::VPC",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":             name,
						"preference_parameters.#": "1",
						"template_scratch_type":   "ResourceImport",
						"source_resources.#":      "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preference_parameters": []map[string]interface{}{
						{
							"parameter_key":   "DeletionPolicy",
							"parameter_value": "Delete",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preference_parameters.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_resources": []map[string]interface{}{
						{
							"resource_id":   "${local.vswitch_id}",
							"resource_type": "ALIYUN::ECS::VSwitch",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_resources.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"execution_mode"},
			},
		},
	})
}
func TestAccAlicloudROSTemplateScratch_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ros_template_scratch.default"
	checkoutSupportedRegions(t, true, connectivity.ROSSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudROSTemplateScratchMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRosTemplateScratch")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srostemplatescratch%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudROSTemplateScratchBasicDependence2)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			// This test needs to ensure that there is at least one VPC and VSwitch labeled {"Created":"No-Deleting"} in the region.
			testAccPreCheckWithEnvVariable(t, "VPC_WITH_NO-DELETING_TAG")
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":           "${var.name}",
					"template_scratch_type": "ResourceImport",
					"preference_parameters": []map[string]interface{}{
						{
							"parameter_key":   "DeletionPolicy",
							"parameter_value": "Retain",
						},
					},

					"source_tag": []map[string]interface{}{
						{
							"resource_tags": map[string]string{
								"Created": "No-Deleting",
							},
							"resource_type_filter": []string{"ALIYUN::ECS::VPC"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":             name,
						"preference_parameters.#": "1",
						"template_scratch_type":   "ResourceImport",
						"source_tag.#":            "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"logical_id_strategy": "LongTypePrefixAndHashSuffix",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"logical_id_strategy": "LongTypePrefixAndHashSuffix",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"execution_mode": "Sync",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_tag": []map[string]interface{}{
						{
							"resource_tags": map[string]string{
								"Created": "No-Deleting",
							},
							"resource_type_filter": []string{"ALIYUN::ECS::VSwitch"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_tag.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preference_parameters": []map[string]interface{}{
						{
							"parameter_key":   "DeletionPolicy",
							"parameter_value": "Delete",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preference_parameters.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"execution_mode"},
			},
		},
	})
}
func TestAccAlicloudROSTemplateScratch_basic3(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ros_template_scratch.default"
	checkoutSupportedRegions(t, true, connectivity.ROSSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudROSTemplateScratchMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRosTemplateScratch")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srostemplatescratch%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudROSTemplateScratchBasicDependence1)
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
					"description":           "${var.name}",
					"template_scratch_type": "ArchitectureReplication",
					"execution_mode":        "Async",
					"logical_id_strategy":   "LongTypePrefixAndIndexSuffix",
					"preference_parameters": []map[string]interface{}{
						{
							"parameter_key":   "DeletionPolicy",
							"parameter_value": "Retain",
						},
					},
					"source_resources": []map[string]interface{}{
						{
							"resource_id":   "${data.alicloud_vpcs.default.ids.0}",
							"resource_type": "ALIYUN::ECS::VPC",
						},
						{
							"resource_id":   "${local.vswitch_id}",
							"resource_type": "ALIYUN::ECS::VSwitch",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":             name,
						"preference_parameters.#": CHECKSET,
						"template_scratch_type":   "ArchitectureReplication",
						"source_resources.#":      "2",
						"logical_id_strategy":     "LongTypePrefixAndIndexSuffix",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"execution_mode"},
			},
		},
	})
}

var AlicloudROSTemplateScratchMap0 = map[string]string{
	"preference_parameters.#": CHECKSET,
	"status":                  CHECKSET,
	"source_resource_group.#": CHECKSET,
	"source_resources.#":      CHECKSET,
}

func AlicloudROSTemplateScratchBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

`, name)
}

func AlicloudROSTemplateScratchBasicDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
}

locals {
  vswitch_id = data.alicloud_vswitches.default.ids.0
  zone_id     = data.alicloud_vswitches.default.vswitches.0.zone_id
}

`, name)
}

func AlicloudROSTemplateScratchBasicDependence2(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}

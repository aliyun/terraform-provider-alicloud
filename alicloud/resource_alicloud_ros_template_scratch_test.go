package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("ROS", "2019-09-10", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
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
			_, err = client.RpcPost("ROS", "2019-09-10", action, nil, request, false)
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

// collection-order: bypass source_resource_group.resource_type_filter -- pre-existing attribute outside this change; adding a reorder apply reproduces a pre-existing post-apply convergence failure (description and preference parameters revert after the source group update), which needs a separate update-path fix
func TestAccAliCloudROSTemplateScratch_basic0(t *testing.T) {
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
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
func TestAccAliCloudROSTemplateScratch_basic1(t *testing.T) {
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
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
func TestAccAliCloudROSTemplateScratch_basic2(t *testing.T) {
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
	name := fmt.Sprintf("tf-testAcc%srostemplatescratch%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudROSTemplateScratchBasicDependence2)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":           "${var.name}",
					"template_scratch_type": "ResourceImport",
					"depends_on":            []string{"alicloud_vpc.default", "alicloud_vswitch.default"},
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
								"Env":     "Terraform",
							},
							"resource_type_filter": []string{"ALIYUN::ECS::VPC"},
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
					"source_tag": []map[string]interface{}{
						{
							"resource_tags": map[string]string{
								"Created": "No-Deleting",
								"Env":     "Terraform",
							},
							"resource_type_filter": []string{"ALIYUN::ECS::VPC", "ALIYUN::ECS::VSwitch"},
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
					"source_tag": []map[string]interface{}{
						{
							"resource_tags": map[string]string{
								"Created": "No-Deleting",
								"Env":     "Terraform",
							},
							"resource_type_filter": []string{"ALIYUN::ECS::VSwitch", "ALIYUN::ECS::VPC"},
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
func TestAccAliCloudROSTemplateScratch_basic3(t *testing.T) {
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
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
				Config: testAccConfig(map[string]interface{}{
					"preference_parameters": []map[string]interface{}{
						{
							"parameter_key":   "NamePrefix",
							"parameter_value": "tf-replica",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preference_parameters.#": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"execution_mode", "preference_parameters"},
			},
		},
	})
}

func TestAccAliCloudROSTemplateScratch_basic4(t *testing.T) {
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
	name := fmt.Sprintf("tf-testAcc%srostemplatescratch%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudROSTemplateScratchBasicDependence3)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":           "${var.name}",
					"template_scratch_type": "ResourceImport",
					"resource_group_id":     "${alicloud_resource_manager_resource_group.default.id}",
					"depends_on":            []string{"alicloud_vpc.default", "alicloud_vpc.second"},
					"preference_parameters": []map[string]interface{}{
						{
							"parameter_key":   "DeletionPolicy",
							"parameter_value": "Retain",
						},
					},
					"source_resource_group": []map[string]interface{}{
						{
							"resource_group_id":    "${alicloud_resource_manager_resource_group.default.id}",
							"resource_type_filter": []string{"ALIYUN::ECS::VPC"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":             name,
						"resource_group_id":       CHECKSET,
						"source_resource_group.#": "1",
						"stacks.#":                "0",
					}),
					resource.TestCheckResourceAttrPair(resourceId, "resource_group_id", "alicloud_resource_manager_resource_group.default", "id"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${alicloud_resource_manager_resource_group.second.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
					resource.TestCheckResourceAttrPair(resourceId, "resource_group_id", "alicloud_resource_manager_resource_group.second", "id"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_resource_group": []map[string]interface{}{
						{
							"resource_group_id":    "${alicloud_resource_manager_resource_group.second.id}",
							"resource_type_filter": []string{"ALIYUN::ECS::VPC"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_resource_group.#": "1",
					}),
					resource.TestCheckTypeSetElemAttrPair(resourceId, "source_resource_group.*.resource_group_id", "alicloud_resource_manager_resource_group.second", "id"),
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
    name_regex = "^default-NODELETING$"
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

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"

  tags = {
    Created = "No-Deleting"
    Env     = "Terraform"
  }
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name

  tags = {
    Created = "No-Deleting"
    Env     = "Terraform"
  }
}
`, name)
}

func AlicloudROSTemplateScratchBasicDependence3(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_resource_manager_resource_group" "default" {
  resource_group_name = var.name
  display_name        = var.name
}

resource "alicloud_resource_manager_resource_group" "second" {
  resource_group_name = "${var.name}-2"
  display_name        = "${var.name}-2"
}

resource "alicloud_vpc" "default" {
  vpc_name          = var.name
  cidr_block        = "172.16.0.0/12"
  resource_group_id = alicloud_resource_manager_resource_group.default.id
}

resource "alicloud_vpc" "second" {
  vpc_name          = "${var.name}-2"
  cidr_block        = "192.168.0.0/16"
  resource_group_id = alicloud_resource_manager_resource_group.second.id
}

`, name)
}

func TestAccAliCloudROSTemplateScratch_basic5(t *testing.T) {
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
	name := fmt.Sprintf("tf-testAcc%srostemplatescratch%d", defaultRegionToTest, rand)
	stackName := fmt.Sprintf("tf-testAcc%srosstack%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudROSTemplateScratchBasicDependence4)

	var scratchId, stackId, region string
	// The import stack is created out of band, so the scratch cannot be destroyed while the stack
	// stays associated. ErrorCheck runs before the harness destroys the state left by a failed
	// step, which keeps every failure path in stack-then-scratch order.
	deleteImportStack := func() error {
		if err := deleteRosImportStack(stackId, region); err != nil {
			return fmt.Errorf("import stack cleanup failed: %w", err)
		}
		return nil
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		ErrorCheck: func(err error) error {
			if cleanupErr := deleteImportStack(); cleanupErr != nil {
				return fmt.Errorf("%w; %s", err, cleanupErr)
			}
			return err
		},
		CheckDestroy: func(s *terraform.State) error {
			if err := deleteImportStack(); err != nil {
				return err
			}
			return rac.checkResourceDestroy()(s)
		},
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
					"source_resources": []map[string]interface{}{
						{
							"resource_id":   "${alicloud_vpc.default.id}",
							"resource_type": "ALIYUN::ECS::VPC",
						},
						{
							"resource_id":   "${alicloud_vswitch.default.id}",
							"resource_type": "ALIYUN::ECS::VSwitch",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":        name,
						"source_resources.#": "2",
						"stacks.#":           "0",
					}),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceId]
						if !ok {
							return fmt.Errorf("resource %s not found in state", resourceId)
						}
						scratchId = rs.Primary.ID
						region = testAccProvider.Meta().(*connectivity.AliyunClient).RegionId
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					stackId = createRosImportStackFromTemplateScratch(t, stackName, scratchId, region)
				},
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stacks.#": "1",
					}),
					checkRosTemplateScratchStackAssociation(resourceId, &stackId, &region),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"execution_mode"},
			},
			{
				PreConfig: func() {
					if err := deleteRosImportStack(stackId, region); err != nil {
						t.Fatalf("[ERROR] deleting import stack %s before destroying the scenario failed: %s", stackId, err)
					}
				},
				Config:  testAccConfig(map[string]interface{}{}),
				Destroy: true,
			},
		},
	})
}

func createRosImportStackFromTemplateScratch(t *testing.T, stackName, scratchId, region string) string {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	createResponse, err := client.RpcPost("ROS", "2019-09-10", "CreateChangeSet", nil, map[string]interface{}{
		"RegionId":          region,
		"StackName":         stackName,
		"ChangeSetName":     stackName,
		"ChangeSetType":     "IMPORT",
		"TemplateScratchId": scratchId,
		"ClientToken":       buildClientToken("CreateChangeSet"),
	}, true)
	if err != nil {
		t.Fatalf("[ERROR] creating import change set for template scratch %s failed: %s", scratchId, err)
	}
	stackId := fmt.Sprint(createResponse["StackId"])
	changeSetId := fmt.Sprint(createResponse["ChangeSetId"])
	t.Cleanup(func() {
		if err := deleteRosImportStack(stackId, region); err != nil {
			t.Errorf("[ERROR] deleting import stack %s during cleanup failed: %s", stackId, err)
		}
	})
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		if _, err := client.RpcPost("ROS", "2019-09-10", "ExecuteChangeSet", nil, map[string]interface{}{
			"RegionId":    region,
			"ChangeSetId": changeSetId,
			"ClientToken": buildClientToken("ExecuteChangeSet"),
		}, true); err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"ActionInProgress"}) ||
				(IsExpectedErrors(err, []string{"NotSupported"}) && strings.Contains(err.Error(), "CREATE_IN_PROGRESS is not supported")) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("[ERROR] executing import change set %s failed: %s%s", changeSetId, err, deleteRosImportStackOnFailure(stackId, region))
	}
	rosService := RosService{client}
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		object, err := rosService.DescribeRosTemplateScratch(scratchId)
		if err != nil {
			return retry.NonRetryableError(err)
		}
		if stacks, ok := object["Stacks"].([]interface{}); ok {
			for _, item := range stacks {
				if stack, ok := item.(map[string]interface{}); ok && fmt.Sprint(stack["StackId"]) == stackId {
					return nil
				}
			}
		}
		return retry.RetryableError(fmt.Errorf("stack %s is not associated with template scratch %s yet", stackId, scratchId))
	})
	if err != nil {
		t.Fatalf("[ERROR] waiting for stack %s to be associated with template scratch %s failed: %s%s", stackId, scratchId, err, deleteRosImportStackOnFailure(stackId, region))
	}
	return stackId
}

// deleteRosImportStackOnFailure removes a stack the fixture has already created, so a failed
// association step still deletes the stack before the harness destroys the template scratch. It
// returns the cleanup outcome for the fatal message instead of swallowing the error.
func deleteRosImportStackOnFailure(stackId, region string) string {
	if err := deleteRosImportStack(stackId, region); err != nil {
		return fmt.Sprintf("; deleting import stack %s failed: %s", stackId, err)
	}
	return ""
}

func deleteRosImportStack(stackId, region string) error {
	if stackId == "" {
		return nil
	}
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	// Make repeated calls idempotent: cleanup, the destroy step and CheckDestroy may all run.
	if existing, err := client.RpcPost("ROS", "2019-09-10", "GetStack", nil, map[string]interface{}{
		"RegionId": region,
		"StackId":  stackId,
	}, true); err != nil {
		if IsExpectedErrors(err, []string{"StackNotFound"}) {
			return nil
		}
		return err
	} else if fmt.Sprint(existing["Status"]) == "DELETE_COMPLETE" {
		return nil
	}
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		if _, err := client.RpcPost("ROS", "2019-09-10", "DeleteStack", nil, map[string]interface{}{
			"RegionId":           region,
			"StackId":            stackId,
			"RetainAllResources": true,
			"ClientToken":        buildClientToken("DeleteStack"),
		}, true); err != nil {
			if IsExpectedErrors(err, []string{"StackNotFound"}) {
				return nil
			}
			if NeedRetry(err) || IsExpectedErrors(err, []string{"ActionInProgress"}) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(fmt.Errorf("deleting import stack %s: %w", stackId, err))
		}
		return nil
	})
	if err != nil {
		return err
	}
	return retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err := client.RpcPost("ROS", "2019-09-10", "GetStack", nil, map[string]interface{}{
			"RegionId": region,
			"StackId":  stackId,
		}, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"StackNotFound"}) {
				return nil
			}
			return retry.NonRetryableError(err)
		}
		switch status := fmt.Sprint(response["Status"]); status {
		case "DELETE_COMPLETE":
			return nil
		case "DELETE_FAILED":
			return retry.NonRetryableError(fmt.Errorf("import stack %s entered DELETE_FAILED: %s", stackId, fmt.Sprint(response["StatusReason"])))
		default:
			return retry.RetryableError(fmt.Errorf("import stack %s is not deleted yet (status %s)", stackId, status))
		}
	})
}

func AlicloudROSTemplateScratchBasicDependence4(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}
`, name)
}

func TestUnitRosTemplateScratchResourceGroupWire(t *testing.T) {
	var mu sync.Mutex
	createForm := map[string]string{}
	updateForm := map[string]string{}
	getForm := map[string]string{}
	// Keep the value returned by GetTemplateScratch different from the requested one so the
	// read-back assertions below can only pass when Read writes the API response into state.
	currentResourceGroupId := ""
	client := rosTemplateScratchUnitTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		action := r.Form.Get("Action")
		mu.Lock()
		switch action {
		case "CreateTemplateScratch":
			createForm = firstRosMockFormValue(r)
			currentResourceGroupId = "rg-unit-create-readback"
			mu.Unlock()
			writeRosMockResponse(t, w, map[string]interface{}{"RequestId": "mock-1", "TemplateScratchId": "ts-unit-1"})
			return
		case "UpdateTemplateScratch":
			updateForm = firstRosMockFormValue(r)
			if resourceGroupId := r.Form.Get("ResourceGroupId"); resourceGroupId == "rg-unit-updated" {
				currentResourceGroupId = "rg-unit-updated-readback"
			}
			mu.Unlock()
			writeRosMockResponse(t, w, map[string]interface{}{"RequestId": "mock-2"})
			return
		case "GetTemplateScratch":
			getForm = firstRosMockFormValue(r)
			resourceGroupId := currentResourceGroupId
			mu.Unlock()
			writeRosMockResponse(t, w, map[string]interface{}{
				"RequestId": "mock-3",
				"TemplateScratch": map[string]interface{}{
					"TemplateScratchId":   "ts-unit-1",
					"TemplateScratchType": "ResourceImport",
					"Status":              "GENERATE_COMPLETE",
					"Description":         "wire-test",
					"ResourceGroupId":     resourceGroupId,
					"SourceResources": []interface{}{
						map[string]interface{}{"ResourceId": "vpc-unit-1", "ResourceType": "ALIYUN::ECS::VPC"},
					},
				},
			})
			return
		default:
			mu.Unlock()
			t.Errorf("unexpected ROS action: %s", action)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	})

	r := resourceAlicloudRosTemplateScratch()
	r.Timeouts = &schema.ResourceTimeout{
		Create: schema.DefaultTimeout(20 * time.Second),
		Update: schema.DefaultTimeout(20 * time.Second),
	}

	d := r.Data(nil)
	d.MarkNewResource()
	if err := d.Set("description", "wire-test"); err != nil {
		t.Fatalf("setting description failed: %s", err)
	}
	if err := d.Set("template_scratch_type", "ResourceImport"); err != nil {
		t.Fatalf("setting template_scratch_type failed: %s", err)
	}
	if err := d.Set("resource_group_id", "rg-unit-create"); err != nil {
		t.Fatalf("setting resource_group_id failed: %s", err)
	}
	if err := d.Set("source_resources", []interface{}{
		map[string]interface{}{"resource_id": "vpc-unit-1", "resource_type": "ALIYUN::ECS::VPC"},
	}); err != nil {
		t.Fatalf("setting source_resources failed: %s", err)
	}
	if err := r.Create(d, client); err != nil {
		t.Fatalf("Create failed: %s", err)
	}
	if got := createForm["Action"]; got != "CreateTemplateScratch" {
		t.Fatalf("unexpected create action: %q", got)
	}
	if got := createForm["ResourceGroupId"]; got != "rg-unit-create" {
		t.Fatalf("CreateTemplateScratch must carry the configured resource group, got %q", got)
	}
	if got := d.Get("resource_group_id"); got != "rg-unit-create-readback" {
		t.Fatalf("GetTemplateScratch response must be read back into state, got %q", got)
	}
	if got := getForm["TemplateScratchId"]; got != "ts-unit-1" {
		t.Fatalf("GetTemplateScratch must query the created scenario, got %q", got)
	}

	state := &sdktf.InstanceState{ID: "ts-unit-1", Attributes: map[string]string{
		"id":                               "ts-unit-1",
		"description":                      "wire-test",
		"template_scratch_type":            "ResourceImport",
		"resource_group_id":                "rg-unit-create",
		"source_resources.#":               "1",
		"source_resources.0.resource_id":   "vpc-unit-1",
		"source_resources.0.resource_type": "ALIYUN::ECS::VPC",
	}}
	diff := sdktf.NewInstanceDiff()
	diff.Attributes["resource_group_id"] = &sdktf.ResourceAttrDiff{Old: "rg-unit-create", New: "rg-unit-updated"}
	updateData, err := schema.InternalMap(r.Schema).Data(state, diff)
	if err != nil {
		t.Fatalf("building update data failed: %s", err)
	}
	if err := r.Update(updateData, client); err != nil {
		t.Fatalf("Update failed: %s", err)
	}
	if got := updateForm["Action"]; got != "UpdateTemplateScratch" {
		t.Fatalf("unexpected update action: %q", got)
	}
	if got := updateForm["TemplateScratchId"]; got != "ts-unit-1" {
		t.Fatalf("UpdateTemplateScratch must target the scenario, got %q", got)
	}
	if got := updateForm["ResourceGroupId"]; got != "rg-unit-updated" {
		t.Fatalf("UpdateTemplateScratch must carry the changed resource group, got %q", got)
	}
	if got := updateData.Get("resource_group_id"); got != "rg-unit-updated-readback" {
		t.Fatalf("GetTemplateScratch after update must be read back into state, got %q", got)
	}
}

func checkRosTemplateScratchStackAssociation(resourceId string, stackId, region *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceId)
		}
		if *stackId == "" {
			return fmt.Errorf("the controlled import stack was not created")
		}
		attributes := rs.Primary.Attributes
		for key, value := range attributes {
			if !strings.HasPrefix(key, "stacks.") || !strings.HasSuffix(key, ".stack_id") {
				continue
			}
			if value != *stackId {
				continue
			}
			prefix := strings.TrimSuffix(key, "stack_id")
			if got := attributes[prefix+"region_id"]; got != *region {
				return fmt.Errorf("stacks region_id is %q, want the scenario region %q", got, *region)
			}
			return nil
		}
		return fmt.Errorf("state does not report the controlled stack %s in stacks", *stackId)
	}
}

func rosTemplateScratchUnitTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
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
	endpoints.Store("ros", endpoint)
	return client
}

func firstRosMockFormValue(r *http.Request) map[string]string {
	out := make(map[string]string, len(r.Form))
	for key, values := range r.Form {
		if len(values) > 0 {
			out[key] = values[0]
		}
	}
	return out
}

func writeRosMockResponse(t *testing.T, w http.ResponseWriter, body map[string]interface{}) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(body); err != nil {
		t.Error(err)
	}
}

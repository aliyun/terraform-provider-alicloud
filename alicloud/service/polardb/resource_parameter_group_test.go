package polardb_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	polardb "github.com/aliyun/terraform-provider-alicloud/alicloud/service/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_polardb_parameter_group", &resource.Sweeper{
		Name: "alicloud_polardb_parameter_group",
		F:    testSweepPolarDBParameterGroup,
	})
}

func testSweepPolarDBParameterGroup(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeParameterGroups"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	var response map[string]interface{}
	polarDBParameterGroupIds := make([]string, 0)
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_polardb_parameter_group", action, errs.AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.ParameterGroups", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.ParameterGroups", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		skip := true
		item := v.(map[string]interface{})
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(fmt.Sprint(item["ParameterGroupName"])), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping PolarDBParameterGroup Instance: %v", item["ParameterGroupName"])
			continue
		}
		polarDBParameterGroupIds = append(polarDBParameterGroupIds, fmt.Sprint(item["ParameterGroupId"]))
	}

	for _, id := range polarDBParameterGroupIds {
		log.Printf("[INFO] Deleting PolarDBParameterGroup Instance: %s", id)
		deleteAction := "DeleteParameterGroup"
		if err != nil {
			return errs.WrapError(err)
		}
		request = map[string]interface{}{
			"RegionId":         client.RegionId,
			"ParameterGroupId": id,
		}
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(3*time.Minute, func() *retry.RetryError {
			_, err = client.RpcPost("polardb", "2017-08-01", deleteAction, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete PolarDBParameterGroup Instance (%s): %s", polarDBParameterGroupIds, err)
		}
	}
	return nil
}

func TestAccAliCloudPolarDBParameterGroup_basic00(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_polardb_parameter_group.default"
	ra := acctest.ResourceAttrInit(resourceId, resourceAliCloudPolarDbParameterGroupMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return polardb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePolarDBParameterGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sPolarDBParameterGroup-name%d", defaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceAliCloudPolarDbParameterGroupBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":       "tf_testAcc",
					"db_type":    "MySQL",
					"db_version": "8.0",
					"parameters": []map[string]interface{}{
						{
							"param_name":  "wait_timeout",
							"param_value": "86400",
						},
						{
							"param_name":  "innodb_old_blocks_time",
							"param_value": "1000",
						},
					},
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         "tf_testAcc",
						"db_type":      "MySQL",
						"db_version":   "8.0",
						"parameters.#": "2",
						"description":  name,
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

var resourceAliCloudPolarDbParameterGroupMap = map[string]string{}

func resourceAliCloudPolarDbParameterGroupBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

`, name)
}

// Test PolarDb ParameterGroup. >>> Resource test cases, automatically generated.
// Case 属性全覆盖_多参数 11763
func TestAccAliCloudPolarDbParameterGroup_basic11763(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_polardb_parameter_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudPolarDbParameterGroupMap11763)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return polardb.NewPolarDbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePolarDbParameterGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudPolarDbParameterGroupBasicDependence11763)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_version":           "8.0",
					"parameter_group_name": name,
					"db_type":              "MySQL",
					"parameters": []map[string]interface{}{
						{
							"param_name":  "wait_timeout",
							"param_value": "86400",
						},
						{
							"param_name":  "innodb_old_blocks_time",
							"param_value": "1000",
						},
						{
							"param_name":  "default_time_zone",
							"param_value": "SYSTEM",
						},
					},
					"description": "tf_testAcc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_version":           CHECKSET,
						"parameter_group_name": name,
						"db_type":              "MySQL",
						"parameters.#":         "3",
						"description":          "tf_testAcc",
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

var AlicloudPolarDbParameterGroupMap11763 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudPolarDbParameterGroupBasicDependence11763(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test PolarDb ParameterGroup. <<< Resource test cases, automatically generated.

// TestAccAliCloudPolarDbParameterGroup_resourceGroupId covers the create-only resource_group_id mapping.
// Every successful configuration step is followed by an independent import; removing the configured
// value replaces the parameter group while keeping the newly created object readable.
func TestAccAliCloudPolarDbParameterGroup_resourceGroupId(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_polardb_parameter_group.default"
	ra := acctest.ResourceAttrInit(resourceId, map[string]string{
		"create_time": CHECKSET,
	})
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return polardb.NewPolarDbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePolarDbParameterGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudPolarDbParameterGroupBasicDependenceResourceGroupId)
	var idBeforeRemoval string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":              "MySQL",
					"db_version":           "8.0",
					"parameter_group_name": name,
					"parameters": []map[string]interface{}{
						{
							"param_name":  "wait_timeout",
							"param_value": "86400",
						},
						{
							"param_name":  "innodb_old_blocks_time",
							"param_value": "1000",
						},
					},
					"description": "tf_testAcc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_type":              "MySQL",
						"db_version":           "8.0",
						"parameter_group_name": name,
						"parameters.#":         "2",
						"description":          "tf_testAcc",
					}),
					resource.TestCheckNoResourceAttr(resourceId, "resource_group_id"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resource_group_id"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${alicloud_resource_manager_resource_group.group1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters.#": "2",
					}),
					resource.TestCheckResourceAttrPair(resourceId, "resource_group_id", "alicloud_resource_manager_resource_group.group1", "id"),
					testAccCheckPolarDbParameterGroupInResourceGroup(resourceId, "alicloud_resource_manager_resource_group.group1"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resource_group_id"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${alicloud_resource_manager_resource_group.group2.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters.#": "2",
					}),
					resource.TestCheckResourceAttrPair(resourceId, "resource_group_id", "alicloud_resource_manager_resource_group.group2", "id"),
					testAccCheckPolarDbParameterGroupInResourceGroup(resourceId, "alicloud_resource_manager_resource_group.group2"),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceId]
						if !ok || rs.Primary.ID == "" {
							return fmt.Errorf("can't find resource by id: %s", resourceId)
						}
						idBeforeRemoval = rs.Primary.ID
						return nil
					},
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resource_group_id"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.REMOVEKEY,
						"parameters.#":      "2",
					}),
					resource.TestCheckNoResourceAttr(resourceId, "resource_group_id"),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceId]
						if !ok || rs.Primary.ID == "" {
							return fmt.Errorf("can't find resource by id: %s", resourceId)
						}
						if rs.Primary.ID == idBeforeRemoval {
							return fmt.Errorf("expected the parameter group to be replaced after removing resource_group_id, but the id is unchanged: %s", rs.Primary.ID)
						}
						return nil
					},
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resource_group_id"},
			},
		},
	})
}

// testAccCheckPolarDbParameterGroupInResourceGroup verifies that the parameter group is returned
// by DescribeParameterGroups when filtering on the expected resource group, which proves the
// create-time ResourceGroupId is actually applied by the service.
func testAccCheckPolarDbParameterGroupInResourceGroup(parameterGroupResourceId, resourceGroupResourceId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		pgRs, ok := s.RootModule().Resources[parameterGroupResourceId]
		if !ok || pgRs.Primary.ID == "" {
			return fmt.Errorf("can't find resource by id: %s", parameterGroupResourceId)
		}
		rgRs, ok := s.RootModule().Resources[resourceGroupResourceId]
		if !ok || rgRs.Primary.ID == "" {
			return fmt.Errorf("can't find resource by id: %s", resourceGroupResourceId)
		}

		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		request := map[string]interface{}{
			"RegionId":        client.RegionId,
			"ResourceGroupId": rgRs.Primary.ID,
		}
		response, err := client.RpcPost("polardb", "2017-08-01", "DescribeParameterGroups", nil, request, true)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, parameterGroupResourceId, "DescribeParameterGroups", errs.AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.ParameterGroups", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, "DescribeParameterGroups", "$.ParameterGroups", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			if fmt.Sprint(item["ParameterGroupId"]) == pgRs.Primary.ID {
				return nil
			}
		}
		return fmt.Errorf("parameter group %s is not returned by DescribeParameterGroups for resource group %s", pgRs.Primary.ID, rgRs.Primary.ID)
	}
}

func AlicloudPolarDbParameterGroupBasicDependenceResourceGroupId(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

locals {
    resource_group_prefix = replace(var.name, "_", "-")
}

resource "alicloud_resource_manager_resource_group" "group1" {
    resource_group_name = "${local.resource_group_prefix}-1"
    display_name        = "${local.resource_group_prefix}-1"
}

resource "alicloud_resource_manager_resource_group" "group2" {
    resource_group_name = "${local.resource_group_prefix}-2"
    display_name        = "${local.resource_group_prefix}-2"
}
`, name)
}

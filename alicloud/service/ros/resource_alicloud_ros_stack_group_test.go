package ros_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ros"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_ros_stack_group",
		&resource.Sweeper{
			Name: "alicloud_ros_stack_group",
			F:    testSweepRosStackGroup,
		})
}

func testSweepRosStackGroup(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return errs.WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	request := map[string]interface{}{
		"PageSize":   helper.PageSizeLarge,
		"PageNumber": 1,
		"RegionId":   region,
		"Status":     "ACTIVE",
	}
	var response map[string]interface{}
	action := "ListStackGroups"
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = client.RpcPost("ROS", "2019-09-10", action, nil, request, true)
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_ros_stack_group", action, errs.AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.StackGroups", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.StackGroups", response)
		}
		sweeped := false
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["StackGroupName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Ros StackGroup: %s", item["StackGroupName"].(string))
				continue
			}
			sweeped = true
			action := "DeleteStackGroup"
			request := map[string]interface{}{
				"StackGroupName": item["StackGroupName"],
				"RegionId":       region,
			}
			_, err = client.RpcPost("ROS", "2019-09-10", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Ros StackGroup (%s): %s", item["StackGroupName"].(string), err)
			}
			if sweeped {
				// Waiting 5 seconds to ensure Ros StackGroup have been deleted.
				time.Sleep(5 * time.Second)
			}
			log.Printf("[INFO] Delete Ros StackGroup success: %s ", item["StackGroupName"].(string))
		}
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

// collection-order: bypass capabilities -- pre-existing stock TypeList; member order semantics unchanged by this verbatim package migration; no two-distinct-member reorder fixture exists

func TestAccAliCloudROSStackGroup_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ros_stack_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRosStackGroupMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ros.NewRosService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRosStackGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudRosStackGroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRosStackGroupBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"stack_group_name": name,
					"template_body":    `{\"ROSTemplateFormatVersion\": \"2015-09-01\", \"Parameters\": {\"Product\": {\"Default\": \"ack-k8s-audit\", \"Type\": \"String\"}, \"VpcName\": {\"Default\": \"DefaultVpc\", \"Type\": \"String\"}, \"InstanceType\": {\"Default\": \"ecs.g6.large\", \"Type\": \"String\"}}}`,

					"parameters": []map[string]interface{}{
						{
							"parameter_key":   "VpcName",
							"parameter_value": "DefaultVpc",
						},
						{
							"parameter_key":   "Product",
							"parameter_value": "ack-k8s-audit",
						},
						{
							"parameter_key":   "InstanceType",
							"parameter_value": "ecs.g6.large",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stack_group_name": name,
						"template_body":    acctest.CHECKSET,
						"parameters.#":     "3",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"template_body", "template_url", "template_version"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"template_body": `{\"ROSTemplateFormatVersion\": \"2015-09-01\", \"Parameters\": {\"VpcName\": {\"Default\": \"DefaultVpc\", \"Type\": \"String\"}}}`,
					"parameters": []map[string]interface{}{
						{
							"parameter_key":   "VpcName",
							"parameter_value": "DefaultVpc",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"template_body": acctest.CHECKSET,
						"parameters.#":  "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"template_body": `{\"ROSTemplateFormatVersion\": \"2015-09-01\", \"Parameters\": {\"Product\": {\"Default\": \"ack-k8s-audit\", \"Type\": \"String\"}, \"VpcName\": {\"Default\": \"DefaultVpc\", \"Type\": \"String\"}, \"InstanceType\": {\"Default\": \"ecs.g6.large\", \"Type\": \"String\"}}}`,
					"parameters": []map[string]interface{}{
						{
							"parameter_key":   "VpcName",
							"parameter_value": "DefaultVpc",
						},
						{
							"parameter_key":   "Product",
							"parameter_value": "ack-k8s-audit",
						},
						{
							"parameter_key":   "InstanceType",
							"parameter_value": "ecs.g6.large",
						},
					},
					"description": "test for tf",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"template_body": acctest.CHECKSET,
						"parameters.#":  "3",
						"description":   "test for tf",
					}),
				),
			},
		},
	})
}

var AlicloudRosStackGroupMap = map[string]string{
	"stack_group_id": acctest.CHECKSET,
	"status":         acctest.CHECKSET,
}

func AlicloudRosStackGroupBasicDependence(name string) string {
	return ""
}

func TestAccAliCloudRosStackGroup_basic11039(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ros_stack_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRosStackGroupMap11039)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ros.NewRosServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRosStackGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccros%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRosStackGroupBasicDependence11039)
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
					"stack_group_name": name,
					"template_id":      "${alicloud_ros_template.example.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stack_group_name": name,
						"template_id":      acctest.CHECKSET,
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
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"template_id"},
			},
		},
	})
}

var AlicloudRosStackGroupMap11039 = map[string]string{
	"status":         acctest.CHECKSET,
	"stack_group_id": acctest.CHECKSET,
}

func AlicloudRosStackGroupBasicDependence11039(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_ros_template" "example" {
  template_name = "example_value"
  template_body = <<EOF
    {
    	"ROSTemplateFormatVersion": "2015-09-01"
    }
    EOF
}
`, name)
}

// Case 创建自助管理模式的资源栈组（简单）_副本1752666515413_副本1752667635487_副本1752668538043_副本1752669172114 11035
func TestAccAliCloudRosStackGroup_basic11035(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ros_stack_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRosStackGroupMap11035)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ros.NewRosServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRosStackGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccros%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRosStackGroupBasicDependence11035)
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
					"stack_group_name": name,
					"template_body":    `{\"ROSTemplateFormatVersion\":\"2015-09-01\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stack_group_name": name,
						"template_body":    acctest.CHECKSET,
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
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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

var AlicloudRosStackGroupMap11035 = map[string]string{
	"status":         acctest.CHECKSET,
	"stack_group_id": acctest.CHECKSET,
}

func AlicloudRosStackGroupBasicDependence11035(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 创建自助管理模式的资源栈组（简单） 11028
func TestAccAliCloudRosStackGroup_basic11028(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ros_stack_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRosStackGroupMap11028)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ros.NewRosServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRosStackGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccros%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRosStackGroupBasicDependence11028)
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
					"stack_group_name": name,
					"template_url":     "https://ros-template-test.oss-cn-hangzhou.aliyuncs.com/ros-templates/test_output_a.json",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stack_group_name": name,
						"template_url":     acctest.CHECKSET,
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
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"template_url"},
			},
		},
	})
}

var AlicloudRosStackGroupMap11028 = map[string]string{
	"status":         acctest.CHECKSET,
	"stack_group_id": acctest.CHECKSET,
}

func AlicloudRosStackGroupBasicDependence11028(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 资源栈组：自助管理权限模式 11026
func TestAccAliCloudRosStackGroup_basic11026(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ros_stack_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRosStackGroupMap11026)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ros.NewRosServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRosStackGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccros%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRosStackGroupBasicDependence11026)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "self-managed",
					"parameters": []map[string]interface{}{
						{
							"parameter_value": "test",
							"parameter_key":   "Name",
						},
					},
					"stack_group_name":  name,
					"template_body":     `{\"ROSTemplateFormatVersion\":\"2015-09-01\", \"Parameters\": {\"Name\": {\"Type\": \"String\", \"Default\": \"default-name\"}}, \"Outputs\": {\"OutName\": {\"Value\": {\"Ref\": \"Name\"}}}}`,
					"permission_model":  "SELF_MANAGED",
					"capabilities":      []string{"ExpandModules"},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "self-managed",
						"parameters.#":      "1",
						"stack_group_name":  name,
						"template_body":     acctest.CHECKSET,
						"permission_model":  "SELF_MANAGED",
						"capabilities.#":    "1",
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"description":       "self-managed-new",
					"parameters": []map[string]interface{}{
						{
							"parameter_value": "test-new",
							"parameter_key":   "Name",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
						"description":       "self-managed-new",
						"parameters.#":      "1",
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
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"capabilities"},
			},
		},
	})
}

var AlicloudRosStackGroupMap11026 = map[string]string{
	"status":         acctest.CHECKSET,
	"stack_group_id": acctest.CHECKSET,
}

func AlicloudRosStackGroupBasicDependence11026(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}
`, name)
}

// Case StackGroup测试用例 823
func TestAccAliCloudRosStackGroup_basic823(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ros_stack_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRosStackGroupMap823)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ros.NewRosServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRosStackGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccros%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRosStackGroupBasicDependence823)
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
					"description":      "资源栈组描述2",
					"stack_group_name": name,
					"template_body":    `{\"ROSTemplateFormatVersion\":\"2015-09-01\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":      "资源栈组描述2",
						"stack_group_name": name,
						"template_body":    acctest.CHECKSET,
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
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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

var AlicloudRosStackGroupMap823 = map[string]string{
	"status":         acctest.CHECKSET,
	"stack_group_id": acctest.CHECKSET,
}

func AlicloudRosStackGroupBasicDependence823(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case StackGroup测试用例_副本1663566976761 1206
func TestAccAliCloudRosStackGroup_basic1206(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ros_stack_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRosStackGroupMap1206)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ros.NewRosServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRosStackGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccros%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRosStackGroupBasicDependence1206)
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
					"description":      "资源栈组描述2",
					"stack_group_name": name,
					"template_body":    `{\"ROSTemplateFormatVersion\":\"2015-09-01\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":      "资源栈组描述2",
						"stack_group_name": name,
						"template_body":    acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":         "修改资源栈组描述123",
					"stack_group_name":    name + "_update",
					"execution_role_name": "AliyunROSStackGroupExecutionRole",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":         "修改资源栈组描述123",
						"stack_group_name":    name + "_update",
						"execution_role_name": "AliyunROSStackGroupExecutionRole",
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
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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

var AlicloudRosStackGroupMap1206 = map[string]string{
	"status":         acctest.CHECKSET,
	"stack_group_id": acctest.CHECKSET,
}

func AlicloudRosStackGroupBasicDependence1206(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case StackGroup测试用例全生命周期 1558
func TestAccAliCloudRosStackGroup_basic1558(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ros_stack_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRosStackGroupMap1558)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ros.NewRosServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRosStackGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccros%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRosStackGroupBasicDependence1558)
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
					"description":              "资源栈组描述2",
					"stack_group_name":         name,
					"template_body":            `{\"ROSTemplateFormatVersion\":\"2015-09-01\"}`,
					"resource_group_id":        "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"template_version":         "v1",
					"administration_role_name": "AliyunROSStackGroupAdministrationRole",
					"permission_model":         "SELF_MANAGED",
					"execution_role_name":      "AliyunROSStackGroupExecutionRole",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":              "资源栈组描述2",
						"stack_group_name":         name,
						"template_body":            acctest.CHECKSET,
						"resource_group_id":        acctest.CHECKSET,
						"template_version":         "v1",
						"administration_role_name": "AliyunROSStackGroupAdministrationRole",
						"permission_model":         "SELF_MANAGED",
						"execution_role_name":      "AliyunROSStackGroupExecutionRole",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"administration_role_name": acctest.REMOVEKEY,
					"execution_role_name":      acctest.REMOVEKEY,
					"description":              "资源栈组描述11",
					"stack_group_name":         name + "_update",
					"permission_model":         "SERVICE_MANAGED",
					"auto_deployment": []map[string]interface{}{
						{
							"retain_stacks_on_account_removal": "false",
							"enabled":                          "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"administration_role_name": acctest.REMOVEKEY,
						"execution_role_name":      acctest.REMOVEKEY,
						"description":              "资源栈组描述11",
						"stack_group_name":         name + "_update",
						"permission_model":         "SERVICE_MANAGED",
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
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"template_body", "template_version"},
			},
		},
	})
}

var AlicloudRosStackGroupMap1558 = map[string]string{
	"status":         acctest.CHECKSET,
	"stack_group_id": acctest.CHECKSET,
}

func AlicloudRosStackGroupBasicDependence1558(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}
`, name)
}

// Test Ros StackGroup. <<< Resource test cases, automatically generated.

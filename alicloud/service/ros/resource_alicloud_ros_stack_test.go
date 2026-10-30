package ros_test

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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ros"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_ros_stack",
		&resource.Sweeper{
			Name: "alicloud_ros_stack",
			F:    testSweepRosStack,
		})
}

func testSweepRosStack(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return errs.WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
		"tf-example",
		"tf_example",
		"terraform-example",
	}
	request := map[string]interface{}{
		"PageSize":   helper.PageSizeLarge,
		"PageNumber": 1,
		"RegionId":   region,
	}
	var response map[string]interface{}
	action := "ListStacks"
	for {
		response, err = client.RpcPost("ROS", "2019-09-10", action, nil, request, true)
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_ros_stack", action, errs.AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Stacks", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Stacks", response)
		}
		sweeped := false
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["StackName"].(string)), strings.ToLower(prefix)) && item["Status"].(string) != "DELETE_COMPLETE" {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Ros Stack: %s", item["StackName"].(string))
				continue
			}
			sweeped = true
			action = "DeleteStack"
			request := map[string]interface{}{
				"StackId":  item["StackId"],
				"RegionId": region,
			}
			_, err = client.RpcPost("ROS", "2019-09-10", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Ros Stack (%s): %s", item["StackName"].(string), err)
			}
			if sweeped {
				// Waiting 5 seconds to ensure Ros Stack have been deleted.
				time.Sleep(5 * time.Second)
			}
			log.Printf("[INFO] Delete Ros Stack success: %s ", item["StackName"].(string))
		}
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAliCloudROSStack_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ros_stack.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRosStackMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ros.NewRosService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRosStack")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudRosStack%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRosStackBasicDependence)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy: func(state *terraform.State) error {
			rosServiceClient := ros.NewRosService(acctest.Provider.Meta().(*connectivity.AliyunClient))
			for _, rs := range state.RootModule().Resources {
				if "alicloud_ros_stack" != rs.Type {
					continue
				}
				stackId := rs.Primary.ID
				stack, err := rosServiceClient.DescribeRosStack(stackId)
				if err != nil {
					if errs.NotFoundError(err) {
						continue
					}
					return errs.WrapError(err)
				}
				if stack != nil {
					if v, ok := stack["Status"]; ok && v == "DELETE_COMPLETE" { // check delete status
						return nil
					}
				}
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"create_option":           "KeepStackOnCreationComplete",
					"deletion_protection":     "Disabled",
					"disable_rollback":        "false",
					"retain_all_resources":    "true",
					"use_previous_parameters": "true",
					"stack_name":              name,
					"stack_policy_body":       `{\"Statement\": [{\"Action\": \"Update:Delete\", \"Resource\": \"*\", \"Effect\": \"Allow\", \"Principal\": \"*\"}]}`,
					"template_body":           `{\"ROSTemplateFormatVersion\":\"2015-09-01\", \"Parameters\": {\"VpcName\": {\"Type\": \"String\"},\"InstanceType\": {\"Type\": \"String\"}}}`,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "ROS",
					},
					"notification_urls": []string{"http://tf-testacc.example.com/ros-event"},
					"retain_resources":  []string{},
					"parameters": []map[string]interface{}{
						{
							"parameter_key":   "VpcName",
							"parameter_value": "VpcName",
						},
						{
							"parameter_key":   "InstanceType",
							"parameter_value": "InstanceType",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stack_name":        name,
						"stack_policy_body": acctest.CHECKSET,
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "ROS",
						"parameters.#":      "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_option", "notification_urls", "replacement_option", "retain_all_resources", "retain_resources", "stack_policy_during_update_body", "stack_policy_body", "stack_policy_during_update_url", "stack_policy_url", "tags", "template_url", "use_previous_parameters"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disable_rollback":  "true",
					"stack_policy_body": `{\"Statement\": [{\"Action\": \"Update:*\", \"Resource\": \"*\", \"Effect\": \"Allow\", \"Principal\": \"*\"}]}`,
					"template_body":     `{\"ROSTemplateFormatVersion\":\"2015-09-01\", \"Description\" : \"模板描述信息，可用于说明模板的适用场景、架构说明等。\", \"Parameters\": {\"VpcName\": {\"Type\": \"String\"},\"InstanceType\": {\"Type\": \"String\"}}}`,
					"parameters": []map[string]interface{}{
						{
							"parameter_key":   "VpcName",
							"parameter_value": "tf-testacc",
						},
						{
							"parameter_key":   "InstanceType",
							"parameter_value": "ECS",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stack_policy_body": acctest.CHECKSET,
						"parameters.#":      "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF_Update",
						"For":     "ROS Update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF_Update",
						"tags.For":     "ROS Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"stack_policy_body": `{\"Statement\": [{\"Action\": \"Update:Delete\", \"Resource\": \"*\", \"Effect\": \"Allow\", \"Principal\": \"*\"}]}`,
					"template_body":     `{\"ROSTemplateFormatVersion\":\"2015-09-01\", \"Parameters\": {\"VpcName\": {\"Type\": \"String\"},\"InstanceType\": {\"Type\": \"String\"}}}`,
					"parameters": []map[string]interface{}{
						{
							"parameter_key":   "VpcName",
							"parameter_value": "VpcName",
						},
						{
							"parameter_key":   "InstanceType",
							"parameter_value": "InstanceType",
						},
					},
					"timeout_in_minutes": "50",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stack_policy_body":  acctest.CHECKSET,
						"timeout_in_minutes": "50",
						"parameters.#":       "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "ROS",
					},
					"timeout_in_minutes": "60",
					"stack_policy_body":  `{\"Statement\": [{\"Action\": \"Update:*\", \"Resource\": \"*\", \"Effect\": \"Allow\", \"Principal\": \"*\"}]}`,
					"template_body":      `{\"ROSTemplateFormatVersion\":\"2015-09-01\", \"Description\" : \"模板描述信息，可用于说明模板的适用场景、架构说明等。\", \"Parameters\": {\"VpcName\": {\"Type\": \"String\"},\"InstanceType\": {\"Type\": \"String\"}}}`,
					"parameters": []map[string]interface{}{
						{
							"parameter_key":   "VpcName",
							"parameter_value": "tf-testacc",
						},
						{
							"parameter_key":   "InstanceType",
							"parameter_value": "ECS",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":             "2",
						"tags.Created":       "TF",
						"tags.For":           "ROS",
						"timeout_in_minutes": "60",
						"stack_policy_body":  acctest.CHECKSET,
						"parameters.#":       "2",
					}),
				),
			},
		},
	})
}

var AlicloudRosStackMap = map[string]string{
	"deletion_protection": "Disabled",
	"status":              acctest.CHECKSET,
}

func AlicloudRosStackBasicDependence(name string) string {
	return ""
}

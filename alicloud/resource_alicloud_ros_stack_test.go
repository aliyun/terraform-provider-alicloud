package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudRosStack_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ros_stack.default"
	ra := resourceAttrInit(resourceId, AlicloudRosStackMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRosStack")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudRosStack%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRosStackBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"stack_name":        name,
					"stack_policy_body": `{\"Statement\": [{\"Action\": \"Update:Delete\", \"Resource\": \"*\", \"Effect\": \"Allow\", \"Principal\": \"*\"}]}`,
					"template_body":     `{\"ROSTemplateFormatVersion\":\"2015-09-01\", \"Parameters\": {\"VpcName\": {\"Type\": \"String\"},\"InstanceType\": {\"Type\": \"String\"}}}`,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "ROS",
					},
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
						"stack_policy_body": CHECKSET,
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
				ImportStateVerifyIgnore: []string{"create_option", "notification_urls", "replacement_option", "retain_all_resources", "retain_resources", "stack_policy_during_update_body", "stack_policy_body", "stack_policy_during_update_url", "stack_policy_url", "template_body", "tags", "template_url", "use_previous_parameters"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
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
						"stack_policy_body": CHECKSET,
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
						"stack_policy_body":  CHECKSET,
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
						"stack_policy_body":  CHECKSET,
						"parameters.#":       "2",
					}),
				),
			},
		},
	})
}

var AlicloudRosStackMap = map[string]string{
	"deletion_protection": "Disabled",
	"status":              CHECKSET,
}

func AlicloudRosStackBasicDependence(name string) string {
	return ""
}

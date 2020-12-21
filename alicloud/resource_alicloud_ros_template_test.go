package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudRosTemplate_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ros_template.default"
	ra := resourceAttrInit(resourceId, AlicloudRosTemplateMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRosTemplate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudRosTemplate%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRosTemplateBasicDependence)
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
					"template_name": name,
					"template_body": `{\"ROSTemplateFormatVersion\":\"2015-09-01\", \"Parameters\": {\"VpcName\": {\"Type\": \"String\"},\"InstanceType\": {\"Type\": \"String\"}}}`,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "ROS",
					},
					"description": "test for ros template",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"template_name": name,
						"template_body": CHECKSET,
						"tags.%":        "2",
						"tags.Created":  "TF",
						"tags.For":      "ROS",
						"description":   "test for ros template",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"template_url", "template_body"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"template_body": `{\"ROSTemplateFormatVersion\":\"2015-09-01\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"template_body": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"template_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"template_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test for ros template update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test for ros template update",
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
					"template_name": name,
					"template_body": `{\"ROSTemplateFormatVersion\":\"2015-09-01\", \"Parameters\": {\"VpcName\": {\"Type\": \"String\"},\"InstanceType\": {\"Type\": \"String\"}}}`,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "ROS",
					},
					"description": "test for ros template"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"template_name": name,
						"template_body": CHECKSET,
						"tags.%":        "2",
						"tags.Created":  "TF",
						"tags.For":      "ROS",
						"description":   "test for ros template"}),
				),
			},
		},
	})
}

var AlicloudRosTemplateMap = map[string]string{}

func AlicloudRosTemplateBasicDependence(name string) string {
	return ""
}

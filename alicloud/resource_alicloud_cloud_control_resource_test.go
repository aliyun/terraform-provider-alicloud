package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCloudControlResource_basic6159_modify(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_control_resource.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudControlResourceMap6159_modify)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudControlServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudControlResource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudcontrolresource%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudControlResourceBasicDependence6159_modify)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"product":           "Ons",
					"resource_code":     "Instance::Topic",
					"resource_id":       "${alicloud_cloud_control_resource.mq_instance.resource_id}",
					"desire_attributes": "{\\\"InstanceId\\\":\\\"${alicloud_cloud_control_resource.mq_instance.resource_id}\\\", \\\"TopicName\\\":\\\"tf-testacc-ons-topic\\\", \\\"MessageType\\\":\\\"1\\\"}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"product":           "Ons",
						"resource_code":     "Instance::Topic",
						"desire_attributes": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"desire_attributes": "{\\\"Tags\\\":[{\\\"TagKey\\\":\\\"k1\\\",\\\"TagValue\\\":\\\"v1\\\"}]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"desire_attributes": "{\"Tags\":[{\"TagKey\":\"k1\",\"TagValue\":\"v1\"}]}",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"desire_attributes"},
			},
		},
	})
}

var AlicloudCloudControlResourceMap6159_modify = map[string]string{
	"resource_id": CHECKSET,
}

func AlicloudCloudControlResourceBasicDependence6159_modify(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_cloud_control_resource" "mq_instance" {
  desire_attributes = "{\"InstanceName\":\"tf-testacc-ons-instance\"}"
  product = "Ons"
  resource_code = "Instance"
} 

`, name)
}

// Test CloudControl Resource. >>> Resource test cases, automatically generated.
// Case 生命周期测试_全参数 6159
func TestAccAliCloudCloudControlResource_basic6159(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_control_resource.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudControlResourceMap6159)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudControlServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudControlResource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudcontrolresource%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudControlResourceBasicDependence6159)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"product":           "ECS",
					"resource_code":     "KeyPair",
					"desire_attributes": "{\\\"KeyPairName\\\":\\\"testds\\\"}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"product":           "ECS",
						"resource_code":     "KeyPair",
						"desire_attributes": "{\"KeyPairName\":\"testds\"}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"desire_attributes": "{\\\"Tags\\\":[{\\\"TagKey\\\":\\\"k1\\\",\\\"TagValue\\\":\\\"v1\\\"}]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"desire_attributes": "{\"Tags\":[{\"TagKey\":\"k1\",\"TagValue\":\"v1\"}]}",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"desire_attributes"},
			},
		},
	})
}

var AlicloudCloudControlResourceMap6159 = map[string]string{
	"resource_id": CHECKSET,
}

func AlicloudCloudControlResourceBasicDependence6159(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test CloudControl Resource. <<< Resource test cases, automatically generated.

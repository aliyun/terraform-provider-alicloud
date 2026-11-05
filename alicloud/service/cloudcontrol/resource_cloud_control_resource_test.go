package cloudcontrol_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudcontrol"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCloudControlResource_basic6159_modify(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_control_resource.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCloudControlResourceMap6159_modify)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudcontrol.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudControlResource")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudcontrolresource%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCloudControlResourceBasicDependence6159_modify)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
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
						"desire_attributes": acctest.CHECKSET,
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
	"resource_id": acctest.CHECKSET,
}

func AlicloudCloudControlResourceBasicDependence6159_modify(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_cloud_control_resource" "mq_instance" {
  desire_attributes = "{\"InstanceName\":\"${var.name}\"}"
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCloudControlResourceMap6159)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudcontrol.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudControlResource")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudcontrolresource%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCloudControlResourceBasicDependence6159)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
	"resource_id": acctest.CHECKSET,
}

func AlicloudCloudControlResourceBasicDependence6159(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test CloudControl Resource. <<< Resource test cases, automatically generated.

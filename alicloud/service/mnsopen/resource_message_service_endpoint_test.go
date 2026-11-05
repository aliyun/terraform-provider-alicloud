package mnsopen_test

import (
	"fmt"
	"testing"

	tfacctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	mnsopen "github.com/aliyun/terraform-provider-alicloud/alicloud/service/mnsopen"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test MessageService Endpoint. >>> Resource test cases, automatically generated.
// Case Endpoint测试用例_v1.1 10073
func TestAccAliCloudMessageServiceEndpoint_basic10073(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_message_service_endpoint.default"
	ra := tfacctest.ResourceAttrInit(resourceId, AliCloudMessageServiceEndpointMap10073)
	rc := tfacctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return mnsopen.NewMessageServiceServiceV2(tfacctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMessageServiceEndpoint")
	rac := tfacctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%smessageserviceendpoint%d", tfacctest.DefaultRegionToTest, rand)
	testAccConfig := tfacctest.ResourceTestAccConfig(resourceId, name, AliCloudMessageServiceEndpointBasicDependence10073)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tfacctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			tfacctest.PreCheck(t)
		},
		ProviderFactories: tfacctest.SDKProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_type":    "public",
					"endpoint_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_type":    "public",
						"endpoint_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_enabled": "true",
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

var AliCloudMessageServiceEndpointMap10073 = map[string]string{}

func AliCloudMessageServiceEndpointBasicDependence10073(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case Endpoint测试用例 9855
func TestAccAliCloudMessageServiceEndpoint_basic9855(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_message_service_endpoint.default"
	ra := tfacctest.ResourceAttrInit(resourceId, AliCloudMessageServiceEndpointMap9855)
	rc := tfacctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return mnsopen.NewMessageServiceServiceV2(tfacctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMessageServiceEndpoint")
	rac := tfacctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%smessageserviceendpoint%d", tfacctest.DefaultRegionToTest, rand)
	testAccConfig := tfacctest.ResourceTestAccConfig(resourceId, name, AliCloudMessageServiceEndpointBasicDependence9855)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tfacctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			tfacctest.PreCheck(t)
		},
		ProviderFactories: tfacctest.SDKProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_type":    "public",
					"endpoint_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_type":    "public",
						"endpoint_enabled": "true",
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

var AliCloudMessageServiceEndpointMap9855 = map[string]string{}

func AliCloudMessageServiceEndpointBasicDependence9855(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test MessageService Endpoint. <<< Resource test cases, automatically generated.

// Case Endpoint测试用例 9856
func TestAccAliCloudMessageServiceEndpoint_basic9856(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_message_service_endpoint.default"
	ra := tfacctest.ResourceAttrInit(resourceId, AliCloudMessageServiceEndpointMap9855)
	rc := tfacctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return mnsopen.NewMessageServiceServiceV2(tfacctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMessageServiceEndpoint")
	rac := tfacctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%smessageserviceendpoint%d", tfacctest.DefaultRegionToTest, rand)
	testAccConfig := tfacctest.ResourceTestAccConfig(resourceId, name, AliCloudMessageServiceEndpointBasicDependence9855)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tfacctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			tfacctest.PreCheck(t)
		},
		ProviderFactories: tfacctest.SDKProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_type":    "public",
					"endpoint_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_type":    "public",
						"endpoint_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_enabled": "true",
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

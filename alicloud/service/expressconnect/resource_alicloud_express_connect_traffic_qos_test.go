package expressconnect_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/expressconnect"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test ExpressConnect TrafficQos. >>> Resource test cases, automatically generated.
// Case QoS策略用例 6535
func TestAccAliCloudExpressConnectTrafficQos_basic6535(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_traffic_qos.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudExpressConnectTrafficQosMap6535)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnect.NewExpressConnectServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectTrafficQos")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccexpressconnect%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudExpressConnectTrafficQosBasicDependence6535)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"qos_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"qos_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"qos_description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"qos_description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
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

func TestAccAliCloudExpressConnectTrafficQos_basic6535_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_traffic_qos.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudExpressConnectTrafficQosMap6535)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnect.NewExpressConnectServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectTrafficQos")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccexpressconnect%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudExpressConnectTrafficQosBasicDependence6535)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"qos_name":          name,
					"qos_description":   name,
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"qos_name":          name,
						"qos_description":   name,
						"resource_group_id": acctest.CHECKSET,
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Test",
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

var AliCloudExpressConnectTrafficQosMap6535 = map[string]string{
	"status":            acctest.CHECKSET,
	"resource_group_id": acctest.CHECKSET,
}

func AliCloudExpressConnectTrafficQosBasicDependence6535(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Test ExpressConnect TrafficQos. <<< Resource test cases, automatically generated.

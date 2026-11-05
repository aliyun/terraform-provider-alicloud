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

// Test ExpressConnect TrafficQosAssociation. >>> Resource test cases, automatically generated.
// Case QoS关联资源-线上 6830
func TestAccAliCloudExpressConnectTrafficQosAssociation_basic6830(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_traffic_qos_association.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudExpressConnectTrafficQosAssociationMap6830)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnect.NewExpressConnectServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectTrafficQosAssociation")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnecttrafficqosassociation%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudExpressConnectTrafficQosAssociationBasicDependence6830)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":   "${data.alicloud_express_connect_physical_connections.default.ids.1}",
					"qos_id":        "${alicloud_express_connect_traffic_qos.创建QoS策略.id}",
					"instance_type": "PHYSICALCONNECTION",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":   acctest.CHECKSET,
						"qos_id":        acctest.CHECKSET,
						"instance_type": "PHYSICALCONNECTION",
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

var AlicloudExpressConnectTrafficQosAssociationMap6830 = map[string]string{
	"status":        acctest.CHECKSET,
	"instance_id":   acctest.CHECKSET,
	"instance_type": acctest.CHECKSET,
}

func AlicloudExpressConnectTrafficQosAssociationBasicDependence6830(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_express_connect_traffic_qos" "创建QoS策略" {
  qos_name        = "meijian-test"
  qos_description = "meijian-test"
}

data "alicloud_express_connect_physical_connections" "default" {
  name_regex = "preserved-NODELETING"
}

`, name)
}

// Case QoS关联资源-线上 6830  raw
func TestAccAliCloudExpressConnectTrafficQosAssociation_basic6830_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_traffic_qos_association.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudExpressConnectTrafficQosAssociationMap6830)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnect.NewExpressConnectServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectTrafficQosAssociation")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnecttrafficqosassociation%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudExpressConnectTrafficQosAssociationBasicDependence6830)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":   "${data.alicloud_express_connect_physical_connections.default.ids.1}",
					"qos_id":        "${alicloud_express_connect_traffic_qos.创建QoS策略.id}",
					"instance_type": "PHYSICALCONNECTION",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":   acctest.CHECKSET,
						"qos_id":        acctest.CHECKSET,
						"instance_type": "PHYSICALCONNECTION",
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

// Test ExpressConnect TrafficQosAssociation. <<< Resource test cases, automatically generated.

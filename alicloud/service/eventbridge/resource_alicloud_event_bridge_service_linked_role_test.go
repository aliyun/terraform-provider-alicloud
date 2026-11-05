package eventbridge_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/eventbridge"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudEventBridgeServiceLinkedRole_basic0(t *testing.T) {
	var v map[string]interface{}
	ResourceId := "alicloud_event_bridge_service_linked_role.default"
	ra := acctest.ResourceAttrInit(ResourceId, AliCloudEventBridgeServiceLinkedRoleMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(ResourceId, &v, func() interface{} {
		return eventbridge.NewEventbridgeService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEventBridgeServiceLinkedRole")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(ResourceId, "", AliCloudEventBridgeServiceLinkedRoleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"product_name": "AliyunServiceRoleForEventBridgeSourceRocketMQ",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"product_name": "AliyunServiceRoleForEventBridgeSourceRocketMQ",
					}),
				),
			},
			{
				ResourceName:      ResourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AliCloudEventBridgeServiceLinkedRoleMap0 = map[string]string{}

func AliCloudEventBridgeServiceLinkedRoleBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}
`, name)
}

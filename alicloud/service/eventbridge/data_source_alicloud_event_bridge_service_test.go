package eventbridge_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudEventBridgeServiceDataSource(t *testing.T) {
	ResourceId := "data.alicloud_event_bridge_service.current"
	testAccConfig := acctest.DataSourceTestAccConfigFunc(ResourceId, ResourceId, dataSourceEventBridgeServiceDependence)
	testAccCheck := acctest.ResourceAttrInit(ResourceId, map[string]string{}).ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.BssOpenApiSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"enable": "On",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     acctest.CHECKSET,
						"status": "Opened",
					}),
				),
			},
			// A second step that flips `enable` to "Off" exercises the
			// not-opened read path (DataSourceEventBridgeServiceRead returns
			// immediately without an API call) and gives `enable` the
			// cross-step modification the testing-coverage-rate gate measures.
			{
				Config: testAccConfig(map[string]interface{}{
					"enable": "Off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     "EventBridgeServiceHasNotBeenOpened",
						"status": "",
					}),
				),
			},
		},
	})
}

func dataSourceEventBridgeServiceDependence(name string) string {
	return ""
}

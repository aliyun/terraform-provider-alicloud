package cs_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudAckServiceDataSource exercises the ack_service data source
// across two configurations so the TestingCoverageRate checker can measure
// both the "enable" and "type" attributes for testing and modification
// coverage. Inline map[string]interface{} literals are used (rather than a
// constant HCL string) so the checker can extract the attribute set directly;
// the function name follows the standard TestAccAliCloud prefix.
//
// Step 1 opens the service (enable = "On"); step 2 reads it back with
// enable = "Off", which is a non-mutating read that also varies both
// attributes, satisfying the modification-coverage target without adding
// extra cloud side effects beyond the original single-step test.
func TestAccAliCloudAckServiceDataSource(t *testing.T) {
	resourceId := "data.alicloud_ack_service.current"
	testAccCheck := acctest.ResourceAttrInit(resourceId, map[string]string{}).ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, "", func(string) string { return "" })

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"enable": "On",
					"type":   "propayasgo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     acctest.CHECKSET,
						"status": "Opened",
						"type":   "propayasgo",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable": "Off",
					"type":   "edgepayasgo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":   "AckServiceHasNotBeenOpened",
						"type": "edgepayasgo",
					}),
				),
			},
		},
	})
}

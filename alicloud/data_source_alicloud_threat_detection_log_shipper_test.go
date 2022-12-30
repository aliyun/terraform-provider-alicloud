package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudThreatDetectionLogShipperDataSource(t *testing.T) {
	resourceId := "data.alicloud_threat_detection_log_shipper.current"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudThreatDetectionLogShipperDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":                 CHECKSET,
						"status":             "Opened",
						"open_status":        CHECKSET,
						"auth_status":        CHECKSET,
						"buy_status":         CHECKSET,
						"sls_project_status": CHECKSET,
						"sls_service_status": CHECKSET,
					}),
				),
			},
		},
	})
}

const testAccCheckAlicloudThreatDetectionLogShipperDataSource = `
data "alicloud_threat_detection_log_shipper" "current" {
	enable = "On"
}
`

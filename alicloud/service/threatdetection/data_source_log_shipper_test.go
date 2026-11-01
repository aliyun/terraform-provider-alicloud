package threatdetection_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudThreatDetectionLogShipperDataSource(t *testing.T) {
	resourceId := "data.alicloud_threat_detection_log_shipper.current"
	testAccCheck := acctest.ResourceAttrInit(resourceId, map[string]string{}).ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudThreatDetectionLogShipperDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":                 acctest.CHECKSET,
						"status":             "Opened",
						"open_status":        acctest.CHECKSET,
						"auth_status":        acctest.CHECKSET,
						"buy_status":         acctest.CHECKSET,
						"sls_project_status": acctest.CHECKSET,
						"sls_service_status": acctest.CHECKSET,
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

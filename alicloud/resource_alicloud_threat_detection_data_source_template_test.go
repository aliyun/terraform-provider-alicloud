package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudThreatDetectionDataSourceTemplate_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_data_source_template.default"
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionDataSourceTemplate")
	rac := resourceAttrCheckInit(rc, nil)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	templateID := os.Getenv("ALICLOUD_THREAT_DETECTION_DATA_SOURCE_TEMPLATE_ID")
	if templateID == "" {
		t.Skip("ALICLOUD_THREAT_DETECTION_DATA_SOURCE_TEMPLATE_ID is not set; skipping because DataSourceTemplate has no Create API and a pre-existing template ID is required")
	}

	name := templateID
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccThreatDetectionDataSourceTemplateBasicDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"data_source_template_id":       "${var.name}",
					"data_source_recognize_enabled": true,
					"auto_scan_new":                 "enabled",
					"data_source_template_name":     "tf-testacc-dst-name-initial",
					"log_project_pattern":           "tf-testacc-log-project-initial",
					"log_store_pattern":             "tf-testacc-log-store-initial",
					"log_user_ids":                  []string{"tf-testacc-user-1"},
					"log_region_ids":                []string{"cn-hangzhou"},
					"lang":                          "zh",
					"role_for":                      0,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_source_template_id":       templateID,
						"data_source_recognize_enabled": "true",
						"auto_scan_new":                 "enabled",
						"data_source_template_name":     CHECKSET,
						"log_project_pattern":           CHECKSET,
						"log_store_pattern":             CHECKSET,
						"create_time":                   CHECKSET,
						"update_time":                   CHECKSET,
						"data_source_from":              CHECKSET,
						"data_source_recognizer":        CHECKSET,
						"region_id":                     CHECKSET,
						"log_user_ids.#":                "1",
						"log_region_ids.#":              "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_source_template_id":       "${var.name}",
					"data_source_recognize_enabled": false,
					"auto_scan_new":                 "disabled",
					"data_source_template_name":     "tf-testacc-dst-name-updated",
					"log_project_pattern":           "tf-testacc-log-project-updated",
					"log_store_pattern":             "tf-testacc-log-store-updated",
					"log_user_ids":                  []string{"tf-testacc-user-2"},
					"log_region_ids":                []string{"cn-shanghai"},
					"lang":                          "en",
					"role_for":                      1,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_source_recognize_enabled": "false",
						"auto_scan_new":                 "disabled",
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

func testAccThreatDetectionDataSourceTemplateBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

provider "alicloud" {
  region = "cn-hangzhou"
}
`, name)
}

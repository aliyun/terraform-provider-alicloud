package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var AlicloudThreatDetectionDataConnectorMap = map[string]string{
	"data_connector_type":     "oss",
	"data_connector_config":   CHECKSET,
	"src_data_type":           "OSS",
	"dest_data_source_id":     CHECKSET,
	"log_project_name":        CHECKSET,
	"log_store_name":          CHECKSET,
	"log_region_id":           CHECKSET,
	"data_connector_status":   CHECKSET,
	"auth_config_id":          CHECKSET,
	"auth_config_vendor":      CHECKSET,
	"auth_config_product":     CHECKSET,
	"data_connector_id":       CHECKSET,
	"data_connector_name":     CHECKSET,
	"sls_ingestion_job_name":  CHECKSET,
	"sls_ingestion_job_state": CHECKSET,
	"creation_time":           CHECKSET,
	"update_time":             CHECKSET,
}

func TestAccAliCloudThreatDetectionDataConnector_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_data_connector.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionDataConnectorMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudSiemService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionDataConnector")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sThreatDetectionDataConnector%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionDataConnectorBasicDependence)
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
					"data_connector_type":   "oss",
					"data_connector_config": "{\\\"bucket\\\":\\\"tf-testacc-bucket\\\",\\\"prefix\\\":\\\"logs\\\"}",
					"src_data_type":         "OSS",
					"dest_data_source_id":   "${var.dest_data_source_id}",
					"log_project_name":      "${alicloud_log_project.default.name}",
					"log_store_name":        "${alicloud_log_store.default.name}",
					"log_region_id":         "${var.region}",
					"data_connector_status": "enabled",
					"auth_config_id":        "${var.dest_data_source_id}",
					"auth_config_vendor":    "APACHE",
					"auth_config_product":   "oss",
					"lang":                  "en",
					"role_for":              0,
					"region_id":             "${var.region}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_connector_type":   "oss",
						"src_data_type":         "OSS",
						"log_region_id":         "${var.region}",
						"data_connector_status": "enabled",
						"lang":                  "en",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_connector_type":   "oss",
					"data_connector_config": "{\\\"bucket\\\":\\\"tf-testacc-bucket\\\",\\\"prefix\\\":\\\"logs\\\"}",
					"src_data_type":         "OSS",
					"dest_data_source_id":   "${var.dest_data_source_id}",
					"log_project_name":      "${alicloud_log_project.default.name}",
					"log_store_name":        "${alicloud_log_store.default.name}",
					"log_region_id":         "${var.region}",
					"data_connector_status": "disabled",
					"auth_config_id":        "${var.dest_data_source_id}",
					"auth_config_vendor":    "AWS_S3",
					"auth_config_product":   "oss",
					"lang":                  "zh",
					"role_for":              0,
					"region_id":             "${var.region}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_connector_status": "disabled",
						"auth_config_vendor":    "AWS_S3",
						"lang":                  "zh",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"data_connector_config", "lang"},
			},
		},
	})
}

func AlicloudThreatDetectionDataConnectorBasicDependence(name string) string {
	return fmt.Sprintf(`
	variable "region" {
		default = "%s"
	}
	variable "dest_data_source_id" {}
	resource "alicloud_log_project" "default" {
		name = "%s"
	}
	resource "alicloud_log_store" "default" {
		project = "${alicloud_log_project.default.name}"
		name    = "%s"
	}
	`, defaultRegionToTest, name, name)
}

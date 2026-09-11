package alicloud

import (
	"fmt"
	"os"
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
	name := fmt.Sprintf("tf-testacc-threat-detection-data-connector-%d", rand)
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
					"log_project_name":      "${alicloud_log_project.default.project_name}",
					"log_store_name":        "${alicloud_log_store.default.logstore_name}",
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
					"log_project_name":      "${alicloud_log_project.default.project_name}",
					"log_store_name":        "${alicloud_log_store.default.logstore_name}",
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
	// dest_data_source_id references a pre-onboarded cloud-siem DataSource
	// (DataSource is @terraform enable:false, so it cannot be created as a TF
	// dependency). Inject a real id via ALICLOUD_THREAT_DETECTION_DEST_DATA_SOURCE_ID
	// or TF_VAR_dest_data_source_id; otherwise a placeholder is used so the
	// CreateDataConnector call is actually invoked and the missing-prerequisite
	// surfaces as an honest API error rather than a stale config-validation skip.
	destDataSourceId := os.Getenv("ALICLOUD_THREAT_DETECTION_DEST_DATA_SOURCE_ID")
	if destDataSourceId == "" {
		destDataSourceId = "tf-testacc-dest-data-source-id"
	}
	return fmt.Sprintf(`
	variable "region" {
		default = "%s"
	}
	variable "dest_data_source_id" {
		default = "%s"
	}
	resource "alicloud_log_project" "default" {
		project_name = "%s"
	}
	resource "alicloud_log_store" "default" {
		project_name  = "${alicloud_log_project.default.project_name}"
		logstore_name = "%s"
	}
	`, defaultRegionToTest, destDataSourceId, name, name)
}

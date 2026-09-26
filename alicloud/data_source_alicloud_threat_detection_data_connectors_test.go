package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var AlicloudThreatDetectionDataConnectorsMap = map[string]string{
	"ids.#":             CHECKSET,
	"data_connectors.#": CHECKSET,
}

func TestAccAliCloudThreatDetectionDataConnectors_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "data.alicloud_threat_detection_data_connectors.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionDataConnectorsMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudSiemService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionDataConnector")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-threat-detection-data-connectors-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionDataConnectorsBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ids":            []interface{}{"${alicloud_threat_detection_data_connector.default.data_connector_id}"},
					"enable_details": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_details": "true",
					}),
				),
			},
		},
	})
}

func AlicloudThreatDetectionDataConnectorsBasicDependence(name string) string {
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
	resource "alicloud_threat_detection_data_connector" "default" {
		data_connector_type  = "oss"
		data_connector_config = "{\"bucket\":\"tf-testacc-bucket\",\"prefix\":\"logs\"}"
		src_data_type        = "OSS"
		dest_data_source_id  = "${var.dest_data_source_id}"
		log_project_name     = "${alicloud_log_project.default.project_name}"
		log_store_name       = "${alicloud_log_store.default.logstore_name}"
		log_region_id        = "${var.region}"
		data_connector_status = "enabled"
		auth_config_vendor   = "APACHE"
	}
	`, defaultRegionToTest, destDataSourceId, name, name)
}

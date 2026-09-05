package alicloud

import (
	"fmt"
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
	resourceId := "alicloud_threat_detection_data_connectors.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionDataConnectorsMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudSiemService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionDataConnector")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sThreatDetectionDataConnectors%d", defaultRegionToTest, rand)
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
	resource "alicloud_threat_detection_data_connector" "default" {
		data_connector_type  = "oss"
		data_connector_config = "{\"bucket\":\"tf-testacc-bucket\",\"prefix\":\"logs\"}"
		src_data_type        = "OSS"
		dest_data_source_id  = "${var.dest_data_source_id}"
		log_project_name     = "${alicloud_log_project.default.name}"
		log_store_name       = "${alicloud_log_store.default.name}"
		log_region_id        = "${var.region}"
		data_connector_status = "enabled"
		auth_config_vendor   = "APACHE"
	}
	`, defaultRegionToTest, name, name)
}

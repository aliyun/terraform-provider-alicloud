package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudArmsIntegrationExportersDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10, 99)
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudArmsIntegrationExportersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_arms_integration_exporter.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudArmsIntegrationExportersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_arms_integration_exporter.default.id}_fake"]`,
		}),
	}
	var existAlicloudArmsIntegrationExportersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                    "1",
			"integration_exporters.#":                  "1",
			"integration_exporters.0.id":               CHECKSET,
			"integration_exporters.0.cluster_id":       CHECKSET,
			"integration_exporters.0.integration_type": "kafka",
			"integration_exporters.0.instance_id":      CHECKSET,
			"integration_exporters.0.param":            "{\"tls_insecure-skip-tls-verify\":\"none=tls.insecure-skip-tls-verify\",\"tls_enabled\":\"none=tls.enabled\",\"sasl_mechanism\":\"\",\"name\":\"kafka1\",\"sasl_enabled\":\"none=sasl.enabled\",\"ip_ports\":\"abc:888\",\"scrape_interval\":30,\"version\":\"0.10.1.0\"}",
			"integration_exporters.0.instance_name":    CHECKSET,
			"integration_exporters.0.exporter_type":    CHECKSET,
			"integration_exporters.0.target":           CHECKSET,
			"integration_exporters.0.version":          CHECKSET,
		}
	}
	var fakeAlicloudArmsIntegrationExportersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   "0",
			"integration_exporters.#": "0",
		}
	}
	var alicloudArmsIntegrationExportersCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_arms_integration_exporters.default",
		existMapFunc: existAlicloudArmsIntegrationExportersDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudArmsIntegrationExportersDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudArmsIntegrationExportersCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, allConf)
}

func testAccCheckAlicloudArmsIntegrationExportersDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testAcc-ArmsIE-%d"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "default-NODELETING"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	resource "alicloud_security_group" "default" {
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_arms_prometheus" "default" {
  		cluster_type        = "ecs"
  		grafana_instance_id = "free"
  		vpc_id              = data.alicloud_vpcs.default.ids.0
  		vswitch_id          = data.alicloud_vswitches.default.ids.0
  		security_group_id   = alicloud_security_group.default.id
  		cluster_name        = "${var.name}-${data.alicloud_vpcs.default.ids.0}"
  		resource_group_id   = data.alicloud_resource_manager_resource_groups.default.groups.0.id
	}

	resource "alicloud_arms_integration_exporter" "default" {
  		cluster_id       = alicloud_arms_prometheus.default.id
  		integration_type = "kafka"
  		param            = "{\"tls_insecure-skip-tls-verify\":\"none=tls.insecure-skip-tls-verify\",\"tls_enabled\":\"none=tls.enabled\",\"sasl_mechanism\":\"\",\"name\":\"kafka1\",\"sasl_enabled\":\"none=sasl.enabled\",\"ip_ports\":\"abc:888\",\"scrape_interval\":30,\"version\":\"0.10.1.0\"}"
	}

	data "alicloud_arms_integration_exporters" "default" {
  		cluster_id       = alicloud_arms_integration_exporter.default.cluster_id
  		integration_type = alicloud_arms_integration_exporter.default.integration_type
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}

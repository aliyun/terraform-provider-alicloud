// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudThreatDetectionCheckItemConfigDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	ThreatDetectionCheckItemConfigCheckInfo.dataSourceTestCheck(t, rand)
}

var existThreatDetectionCheckItemConfigMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"configs.#":                   "1",
		"configs.0.section_ids.#":     CHECKSET,
		"configs.0.description.#":     CHECKSET,
		"configs.0.check_show_name":   CHECKSET,
		"configs.0.vendor":            CHECKSET,
		"configs.0.instance_sub_type": CHECKSET,
		"configs.0.check_id":          CHECKSET,
		"configs.0.custom_configs.#":  CHECKSET,
		"configs.0.check_type":        CHECKSET,
		"configs.0.estimated_count":   CHECKSET,
		"configs.0.instance_type":     CHECKSET,
		"configs.0.risk_level":        CHECKSET,
	}
}

var fakeThreatDetectionCheckItemConfigMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"configs.#": "0",
	}
}

var ThreatDetectionCheckItemConfigCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_threat_detection_check_item_configs.default",
	existMapFunc: existThreatDetectionCheckItemConfigMapFunc,
	fakeMapFunc:  fakeThreatDetectionCheckItemConfigMapFunc,
}

func testAccCheckAlicloudThreatDetectionCheckItemConfigSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccThreatDetectionCheckItemConfig%d"
}

data "alicloud_threat_detection_check_item_configs" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

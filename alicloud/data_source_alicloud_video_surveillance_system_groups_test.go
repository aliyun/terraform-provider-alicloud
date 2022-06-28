package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudVsGroupDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVsGroupDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_video_surveillance_system_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudVsGroupDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_video_surveillance_system_group.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVsGroupDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_video_surveillance_system_group.default.group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudVsGroupDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_video_surveillance_system_group.default.group_name}_fake"`,
		}),
	}
	inProtocolConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVsGroupDataSourceName(rand, map[string]string{
			"in_protocol": `"${alicloud_video_surveillance_system_group.default.in_protocol}"`,
			"ids":         `["${alicloud_video_surveillance_system_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudVsGroupDataSourceName(rand, map[string]string{
			"in_protocol": `"gb28181"`,
			"ids":         `["${alicloud_video_surveillance_system_group.default.id}"]`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVsGroupDataSourceName(rand, map[string]string{
			"status": `"on"`,
			"ids":    `["${alicloud_video_surveillance_system_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudVsGroupDataSourceName(rand, map[string]string{
			"status": `"off"`,
			"ids":    `["${alicloud_video_surveillance_system_group.default.id}"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVsGroupDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_video_surveillance_system_group.default.id}"]`,
			"name_regex":  `"${alicloud_video_surveillance_system_group.default.group_name}"`,
			"in_protocol": `"${alicloud_video_surveillance_system_group.default.in_protocol}"`,
			"region":      `"cn-beijing"`,
			"status":      `"on"`,
		}),
		fakeConfig: testAccCheckAlicloudVsGroupDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_video_surveillance_system_group.default.id}_fake"]`,
			"name_regex":  `"${alicloud_video_surveillance_system_group.default.group_name}_fake"`,
			"in_protocol": `"gb28181"`,
			"region":      `"cn-shenzhen"`,
			"status":      `"off"`,
		}),
	}
	var existAlicloudVsGroupDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                 "1",
			"names.#":               "1",
			"groups.#":              "1",
			"groups.0.description":  fmt.Sprintf("tf-testAccVsGroup-%d", rand),
			"groups.0.group_name":   fmt.Sprintf("tf-testAccVsGroup-%d", rand),
			"groups.0.in_protocol":  "rtmp",
			"groups.0.out_protocol": "flv",
		}
	}
	var fakeAlicloudVsGroupDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudVsGroupCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_video_surveillance_system_groups.default",
		existMapFunc: existAlicloudVsGroupDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudVsGroupDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.SurveillanceSystemSupportRegions)
	}
	alicloudVsGroupCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, inProtocolConf, statusConf, allConf)
}
func testAccCheckAlicloudVsGroupDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccVsGroup-%d"
}

resource "alicloud_video_surveillance_system_group" "default" {
	group_name =  var.name
	description = var.name
	in_protocol = "rtmp"
	out_protocol ="flv"
	play_domain = "bjtqdh.cn"
	push_domain = "jinliyangbj.com"
}

data "alicloud_video_surveillance_system_groups" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

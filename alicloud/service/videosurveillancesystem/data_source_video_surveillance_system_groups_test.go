package videosurveillancesystem_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudVsGroupDataSource(t *testing.T) {
	rand := tfacctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, connectivity.SurveillanceSystemSupportRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			// ids filter, matching
			{
				Config: testAccCheckAlicloudVsGroupDataSourceName(rand, map[string]string{
					"ids": `["${alicloud_video_surveillance_system_group.default.id}"]`,
				}),
				Check: vsGroupDataSourceCheck(existAlicloudVsGroupDataSourceNameMap(rand)),
			},
			// ids filter, no match
			{
				Config: testAccCheckAlicloudVsGroupDataSourceName(rand, map[string]string{
					"ids": `["${alicloud_video_surveillance_system_group.default.id}_fake"]`,
				}),
				Check: vsGroupDataSourceCheck(fakeAlicloudVsGroupDataSourceNameMap()),
			},
			// name_regex filter, matching
			{
				Config: testAccCheckAlicloudVsGroupDataSourceName(rand, map[string]string{
					"name_regex": `"${alicloud_video_surveillance_system_group.default.group_name}"`,
				}),
				Check: vsGroupDataSourceCheck(existAlicloudVsGroupDataSourceNameMap(rand)),
			},
			// name_regex filter, no match
			{
				Config: testAccCheckAlicloudVsGroupDataSourceName(rand, map[string]string{
					"name_regex": `"${alicloud_video_surveillance_system_group.default.group_name}_fake"`,
				}),
				Check: vsGroupDataSourceCheck(fakeAlicloudVsGroupDataSourceNameMap()),
			},
			// in_protocol filter, matching
			{
				Config: testAccCheckAlicloudVsGroupDataSourceName(rand, map[string]string{
					"in_protocol": `"${alicloud_video_surveillance_system_group.default.in_protocol}"`,
					"ids":         `["${alicloud_video_surveillance_system_group.default.id}"]`,
				}),
				Check: vsGroupDataSourceCheck(existAlicloudVsGroupDataSourceNameMap(rand)),
			},
			// in_protocol filter, no match
			{
				Config: testAccCheckAlicloudVsGroupDataSourceName(rand, map[string]string{
					"in_protocol": `"gb28181"`,
					"ids":         `["${alicloud_video_surveillance_system_group.default.id}"]`,
				}),
				Check: vsGroupDataSourceCheck(fakeAlicloudVsGroupDataSourceNameMap()),
			},
			// status filter, matching
			{
				Config: testAccCheckAlicloudVsGroupDataSourceName(rand, map[string]string{
					"status": `"on"`,
					"ids":    `["${alicloud_video_surveillance_system_group.default.id}"]`,
				}),
				Check: vsGroupDataSourceCheck(existAlicloudVsGroupDataSourceNameMap(rand)),
			},
			// status filter, no match
			{
				Config: testAccCheckAlicloudVsGroupDataSourceName(rand, map[string]string{
					"status": `"off"`,
					"ids":    `["${alicloud_video_surveillance_system_group.default.id}"]`,
				}),
				Check: vsGroupDataSourceCheck(fakeAlicloudVsGroupDataSourceNameMap()),
			},
			// all filters, matching
			{
				Config: testAccCheckAlicloudVsGroupDataSourceName(rand, map[string]string{
					"ids":         `["${alicloud_video_surveillance_system_group.default.id}"]`,
					"name_regex":  `"${alicloud_video_surveillance_system_group.default.group_name}"`,
					"in_protocol": `"${alicloud_video_surveillance_system_group.default.in_protocol}"`,
					"region":      `"cn-beijing"`,
					"status":      `"on"`,
				}),
				Check: vsGroupDataSourceCheck(existAlicloudVsGroupDataSourceNameMap(rand)),
			},
			// all filters, no match
			{
				Config: testAccCheckAlicloudVsGroupDataSourceName(rand, map[string]string{
					"ids":         `["${alicloud_video_surveillance_system_group.default.id}_fake"]`,
					"name_regex":  `"${alicloud_video_surveillance_system_group.default.group_name}_fake"`,
					"in_protocol": `"gb28181"`,
					"region":      `"cn-shenzhen"`,
					"status":      `"off"`,
				}),
				Check: vsGroupDataSourceCheck(fakeAlicloudVsGroupDataSourceNameMap()),
			},
		},
	})
}

func existAlicloudVsGroupDataSourceNameMap(rand int) map[string]string {
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

func fakeAlicloudVsGroupDataSourceNameMap() map[string]string {
	return map[string]string{
		"ids.#":   "0",
		"names.#": "0",
	}
}

func vsGroupDataSourceCheck(m map[string]string) resource.TestCheckFunc {
	var checks []resource.TestCheckFunc
	for k, v := range m {
		checks = append(checks, resource.TestCheckResourceAttr("data.alicloud_video_surveillance_system_groups.default", k, v))
	}
	return resource.ComposeTestCheckFunc(checks...)
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

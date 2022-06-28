package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudDirectMailTagsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.DmSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDirectMailTagsDataSourceName(rand, map[string]string{
			"ids": `[alicloud_direct_mail_tag.default.id]`,
		}),
		fakeConfig: testAccCheckAlicloudDirectMailTagsDataSourceName(rand, map[string]string{
			"ids": `["fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDirectMailTagsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_direct_mail_tag.default.tag_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudDirectMailTagsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_direct_mail_tag.default.tag_name}fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDirectMailTagsDataSourceName(rand, map[string]string{
			"ids":        `[alicloud_direct_mail_tag.default.id]`,
			"name_regex": `"${alicloud_direct_mail_tag.default.tag_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudDirectMailTagsDataSourceName(rand, map[string]string{
			"ids":        `["fake"]`,
			"name_regex": `"${alicloud_direct_mail_tag.default.tag_name}fake"`,
		}),
	}
	var existAlicloudDirectMailTagsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":           "1",
			"ids.0":           CHECKSET,
			"names.#":         "1",
			"names.0":         CHECKSET,
			"tags.#":          "1",
			"tags.0.id":       CHECKSET,
			"tags.0.tag_id":   CHECKSET,
			"tags.0.tag_name": fmt.Sprintf("tftestaccdirectmailtag%d", rand),
		}
	}
	var fakeAlicloudDirectMailTagsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudDirectMailTagsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_direct_mail_tags.default",
		existMapFunc: existAlicloudDirectMailTagsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudDirectMailTagsDataSourceNameMapFunc,
	}

	alicloudDirectMailTagsCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudDirectMailTagsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tftestaccdirectmailtag%d"
}

resource "alicloud_direct_mail_tag" "default" {
	tag_name = var.name
}

data "alicloud_direct_mail_tags" "default" {
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

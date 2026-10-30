package directmail_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudDirectMailTagsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.DmSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDirectMailTagsDataSourceName(rand, map[string]string{
			"ids": `[alicloud_direct_mail_tag.default.id]`,
		}),
		FakeConfig: testAccCheckAlicloudDirectMailTagsDataSourceName(rand, map[string]string{
			"ids": `["fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDirectMailTagsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_direct_mail_tag.default.tag_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudDirectMailTagsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_direct_mail_tag.default.tag_name}fake"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDirectMailTagsDataSourceName(rand, map[string]string{
			"ids":        `[alicloud_direct_mail_tag.default.id]`,
			"name_regex": `"${alicloud_direct_mail_tag.default.tag_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudDirectMailTagsDataSourceName(rand, map[string]string{
			"ids":        `["fake"]`,
			"name_regex": `"${alicloud_direct_mail_tag.default.tag_name}fake"`,
		}),
	}
	var existAlicloudDirectMailTagsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":           "1",
			"ids.0":           acctest.CHECKSET,
			"names.#":         "1",
			"names.0":         acctest.CHECKSET,
			"tags.#":          "1",
			"tags.0.id":       acctest.CHECKSET,
			"tags.0.tag_id":   acctest.CHECKSET,
			"tags.0.tag_name": fmt.Sprintf("tftestaccdirectmailtag%d", rand),
		}
	}
	var fakeAlicloudDirectMailTagsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudDirectMailTagsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_direct_mail_tags.default",
		ExistMapFunc: existAlicloudDirectMailTagsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudDirectMailTagsDataSourceNameMapFunc,
	}

	alicloudDirectMailTagsCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
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

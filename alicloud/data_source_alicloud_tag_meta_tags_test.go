package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudTagMetaTagsDataSource(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.TagSupportRegions)
	rand := acctest.RandInt()

	keyNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudTagMetaTagsDataSourceConfig(rand, map[string]string{}),
		fakeConfig:  "",
	}

	var existTagMetaTagsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"tags.#": CHECKSET,
		}
	}

	var fakeTagMetaTagsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"tags.#": "0",
		}
	}

	var tagMetaTagsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_tag_meta_tags.default",
		existMapFunc: existTagMetaTagsMapFunc,
		fakeMapFunc:  fakeTagMetaTagsMapFunc,
	}

	var perCheck = func() {
		testAccPreCheck(t)
	}

	tagMetaTagsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, perCheck, keyNameConf)

}

func testAccCheckAlicloudTagMetaTagsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testaccTagMetaTag%d"
}

data "alicloud_tag_meta_tags" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}

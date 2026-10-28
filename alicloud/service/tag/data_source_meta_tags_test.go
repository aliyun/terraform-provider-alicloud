// Package tag_test.
package tag_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudTagMetaTagsDataSource(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.TagSupportRegions)
	rand := tfacctest.RandInt()

	keyNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudTagMetaTagsDataSourceConfig(rand, map[string]string{}),
		FakeConfig:  "",
	}

	var existTagMetaTagsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"tags.#": acctest.CHECKSET,
		}
	}

	var fakeTagMetaTagsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"tags.#": "0",
		}
	}

	var tagMetaTagsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_tag_meta_tags.default",
		ExistMapFunc: existTagMetaTagsMapFunc,
		FakeMapFunc:  fakeTagMetaTagsMapFunc,
	}

	var perCheck = func() {
		acctest.PreCheck(t)
	}

	tagMetaTagsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, perCheck, keyNameConf)

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

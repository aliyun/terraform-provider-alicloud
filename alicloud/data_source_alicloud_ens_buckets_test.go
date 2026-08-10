package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEnsBucketsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsBucketsSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ens_bucket.default.bucket_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsBucketsSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ens_bucket.default.bucket_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsBucketsSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ens_bucket.default.bucket_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsBucketsSourceConfig(rand, map[string]string{
			"name_regex": `"TestAccAlicloudEnsBucketsDataSource_fake"`,
		}),
	}

	EnsBucketCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, allConf)
}

func testAccCheckAlicloudEnsBucketsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testacc-ens-buckets-ds-%d"
}

resource "alicloud_ens_bucket" "default" {
  bucket_name         = var.name
  bucket_acl          = "private"
  logical_bucket_type = "sink"
  ens_region_id       = "cn-wuxi-telecom_unicom_cmcc-2"
  dispatch_scope      = "domestic"
}

data "alicloud_ens_buckets" "default" {
	%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

var existEnsBucketMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"buckets.#":             "1",
		"buckets.0.bucket_name": fmt.Sprintf("tf-testacc-ens-buckets-ds-%d", rand),
	}
}

var fakeEnsBucketMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"buckets.#": "0",
	}
}

var EnsBucketCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ens_buckets.default",
	existMapFunc: existEnsBucketMapFunc,
	fakeMapFunc:  fakeEnsBucketMapFunc,
}

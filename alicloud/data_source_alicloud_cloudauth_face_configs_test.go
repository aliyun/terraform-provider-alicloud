package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCloudAuthFaceConfigDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudAuthFaceConfigDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cloudauth_face_config.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCloudAuthFaceConfigDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cloudauth_face_config.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudAuthFaceConfigDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cloudauth_face_config.default.biz_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudAuthFaceConfigDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cloudauth_face_config.default.biz_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudAuthFaceConfigDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cloudauth_face_config.default.id}"]`,
			"name_regex": `"${alicloud_cloudauth_face_config.default.biz_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudAuthFaceConfigDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cloudauth_face_config.default.id}_fake"]`,
			"name_regex": `"${alicloud_cloudauth_face_config.default.biz_name}_fake"`,
		}),
	}
	var existAlicloudCloudAuthFaceConfigDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":              "1",
			"names.#":            "1",
			"configs.#":          "1",
			"configs.0.biz_type": fmt.Sprintf("tf-testaccCAfaceconfig%d", rand),
			"configs.0.biz_name": fmt.Sprintf("tf-testaccCAfaceconfig%d", rand),
		}
	}
	var fakeAlicloudCloudAuthFaceConfigDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudCloudAuthFaceConfigCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cloudauth_face_configs.default",
		existMapFunc: existAlicloudCloudAuthFaceConfigDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCloudAuthFaceConfigDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.CloudAuthSupportRegions)
	}

	alicloudCloudAuthFaceConfigCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudCloudAuthFaceConfigDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testaccCAfaceconfig%d"
}

resource "alicloud_cloudauth_face_config" "default" {
	biz_type = var.name
	biz_name = var.name
}

data "alicloud_cloudauth_face_configs" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

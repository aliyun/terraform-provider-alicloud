package swasopen_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudSimpleApplicationServerImagesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSimpleApplicationServerImageDataSourceName(rand, map[string]string{
			"name_regex": `"CentOS-7.3"`,
		}),
		FakeConfig: "",
	}

	instanceImageTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSimpleApplicationServerImageDataSourceName(rand, map[string]string{
			"image_type": `"system"`,
		}),
		FakeConfig: "",
	}
	platformConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSimpleApplicationServerImageDataSourceName(rand, map[string]string{
			"platform": `"Linux"`,
		}),
		FakeConfig: "",
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSimpleApplicationServerImageDataSourceName(rand, map[string]string{
			"image_type": `"system"`,
			"platform":   `"Linux"`,
		}),
		FakeConfig: "",
	}

	var existDataAlicloudSimpleApplicationServerImagesSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                acctest.CHECKSET,
			"images.#":             acctest.CHECKSET,
			"images.0.description": acctest.CHECKSET,
			"images.0.id":          acctest.CHECKSET,
			"images.0.image_id":    acctest.CHECKSET,
			"images.0.image_name":  acctest.CHECKSET,
			"images.0.image_type":  acctest.CHECKSET,
			"images.0.platform":    acctest.CHECKSET,
		}
	}
	var fakeDataAlicloudSimpleApplicationServerImagesSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"images.#": "0",
		}
	}
	var alicloudSimpleApplicationServerImageCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_simple_application_server_images.default",
		ExistMapFunc: existDataAlicloudSimpleApplicationServerImagesSourceNameMapFunc,
		FakeMapFunc:  fakeDataAlicloudSimpleApplicationServerImagesSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheckWithRegions(t, false, connectivity.SimpleApplicationServerNotSupportRegions)
	}
	alicloudSimpleApplicationServerImageCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, instanceImageTypeConf, platformConf, allConf)
}
func testAccCheckAlicloudSimpleApplicationServerImageDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
data "alicloud_simple_application_server_images" "default" {
	%s
}
`, strings.Join(pairs, " \n "))
	return config
}

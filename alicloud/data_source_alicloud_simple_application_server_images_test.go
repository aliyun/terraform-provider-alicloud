package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSimpleApplicationServerImageDataSource(t *testing.T) {
	rand := acctest.RandInt()

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSimpleApplicationServerImageDataSourceName(rand, map[string]string{
			"name_regex": `"CentOS-7.3"`,
		}),
		fakeConfig: "",
	}

	instanceImageTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSimpleApplicationServerImageDataSourceName(rand, map[string]string{
			"image_type": `"system"`,
		}),
		fakeConfig: "",
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSimpleApplicationServerImageDataSourceName(rand, map[string]string{
			"image_type": `"system"`,
		}),
		fakeConfig: "",
	}

	var existDataAlicloudSimpleApplicationServerImagesSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    CHECKSET,
			"images.#": CHECKSET,
		}
	}
	var fakeDataAlicloudSimpleApplicationServerImagesSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    CHECKSET,
			"images.#": CHECKSET,
		}
	}
	var alicloudSimpleApplicationServerImageCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_simple_application_server_images.default",
		existMapFunc: existDataAlicloudSimpleApplicationServerImagesSourceNameMapFunc,
		fakeMapFunc:  fakeDataAlicloudSimpleApplicationServerImagesSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.AlbSupportRegions)
	}
	alicloudSimpleApplicationServerImageCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, instanceImageTypeConf, allConf)
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

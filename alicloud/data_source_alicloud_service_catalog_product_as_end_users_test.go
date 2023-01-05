package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudServiceCatalogProductAsEndUserDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	checkoutSupportedRegions(t, true, connectivity.ServiceCatalogProvisionedProductSupportRegions)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudServiceCatalogProductAsEndUserSourceConfig(rand, map[string]string{
			"name_regex": `"ram模板创建"`,
		}),
		fakeConfig: testAccCheckAlicloudServiceCatalogProductAsEndUserSourceConfig(rand, map[string]string{
			"name_regex": `"ram模板创建_fake"`,
		}),
	}

	ServiceCatalogProductAsEndUserCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf)
}

var existServiceCatalogProductAsEndUserMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"users.#":                           "1",
		"users.0.id":                        CHECKSET,
		"users.0.create_time":               CHECKSET,
		"users.0.has_default_launch_option": CHECKSET,
		"users.0.product_arn":               CHECKSET,
		"users.0.product_id":                CHECKSET,
		"users.0.product_name":              CHECKSET,
		"users.0.product_type":              CHECKSET,
		"users.0.provider_name":             CHECKSET,
	}
}

var fakeServiceCatalogProductAsEndUserMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"users.#": "0",
	}
}

var ServiceCatalogProductAsEndUserCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_service_catalog_product_as_end_users.default",
	existMapFunc: existServiceCatalogProductAsEndUserMapFunc,
	fakeMapFunc:  fakeServiceCatalogProductAsEndUserMapFunc,
}

func testAccCheckAlicloudServiceCatalogProductAsEndUserSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccServiceCatalogProductAsEndUser%d"
}

data "alicloud_service_catalog_product_as_end_users" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

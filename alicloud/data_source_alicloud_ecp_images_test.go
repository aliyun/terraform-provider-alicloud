package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"strings"
	"testing"
)

func TestAccAlicloudEcpImagesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.ECPSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcpImagesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecp_image.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcpImagesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecp_image.default.id}_fake"]`,
		}),
	}
	imageCategoryConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcpImagesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_ecp_image.default.id}"]`,
			"image_category": `"self"`,
		}),
		fakeConfig: testAccCheckAlicloudEcpImagesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_ecp_image.default.id}"]`,
			"image_category": `"system"`,
		}),
	}

	imageNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcpImagesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecp_image.default.id}"]`,
			"image_name": `"${alicloud_ecp_image.default.image_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcpImagesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecp_image.default.id}"]`,
			"image_name": `"${alicloud_ecp_image.default.image_name}_fake"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcpImagesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecp_image.default.id}"]`,
			"name_regex": `"${alicloud_ecp_image.default.image_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcpImagesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecp_image.default.id}"]`,
			"name_regex": `"${alicloud_ecp_image.default.image_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcpImagesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ecp_image.default.id}"]`,
			"status": `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudEcpImagesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ecp_image.default.id}"]`,
			"status": `"CreateFailed"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcpImagesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_ecp_image.default.id}"]`,
			"image_category": `"self"`,
			"image_name":     `"${alicloud_ecp_image.default.image_name}"`,
			"status":         `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudEcpImagesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_ecp_image.default.id}_fake"]`,
			"image_category": `system`,
			"image_name":     `"${alicloud_ecp_image.default.image_name}_fake"`,
			"status":         `"CreateFailed"`,
		}),
	}
	var existAlicloudEcpImagesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "1",
			"names.#":  "1",
			"images.#": "1",
		}
	}
	var fakeAlicloudEcpImagesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudEcpImagesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ecp_images.default",
		existMapFunc: existAlicloudEcpImagesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEcpImagesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	alicloudEcpImagesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, imageCategoryConf, imageNameConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudEcpImagesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccImage-%d"
}

data "alicloud_ecp_instances" "default" {
}

locals {
  instance_id = data.alicloud_ecp_instances.default.instances[0].instance_id
}

resource "alicloud_ecp_image" "default" {
  image_name  = var.name
  description = var.name
  instance_id = "${local.instance_id}"
  force       = "true"
}

data "alicloud_ecp_images" "default" {	
	enable_details = true
	%s	
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

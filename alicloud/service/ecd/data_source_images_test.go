package ecd_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudECDImagesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	acctest.CheckoutSupportedRegions(t, true, connectivity.EcdUserSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcdImagesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecd_image.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudEcdImagesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecd_image.default.id}_fake"]`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcdImagesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecd_image.default.id}"]`,
			"name_regex": `"${alicloud_ecd_image.default.image_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudEcdImagesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecd_image.default.id}"]`,
			"name_regex": `"${alicloud_ecd_image.default.image_name}_fake"`,
		}),
	}

	imageTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcdImagesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecd_image.default.id}"]`,
			"image_type": `"CUSTOM"`,
		}),
		FakeConfig: testAccCheckAlicloudEcdImagesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecd_image.default.id}"]`,
			"image_type": `"SYSTEM"`,
		}),
	}
	osTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcdImagesDataSourceName(rand, map[string]string{
			"ids":     `["${alicloud_ecd_image.default.id}"]`,
			"os_type": `"${data.alicloud_ecd_bundles.default.bundles.1.os_type}"`,
		}),
		FakeConfig: "",
	}
	desktopInstanceTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcdImagesDataSourceName(rand, map[string]string{
			"ids":                   `["${alicloud_ecd_image.default.id}"]`,
			"desktop_instance_type": `"${data.alicloud_ecd_bundles.default.bundles.1.desktop_type}"`,
		}),
		FakeConfig: "",
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcdImagesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecd_image.default.image_name}"`,
			"status":     `"Available"`,
		}),
		FakeConfig: testAccCheckAlicloudEcdImagesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecd_image.default.image_name}"`,
			"status":     `"Creating"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcdImagesDataSourceName(rand, map[string]string{
			"ids":                   `["${alicloud_ecd_image.default.id}"]`,
			"name_regex":            `"${alicloud_ecd_image.default.image_name}"`,
			"image_type":            `"CUSTOM"`,
			"status":                `"Available"`,
			"os_type":               `"${data.alicloud_ecd_bundles.default.bundles.1.os_type}"`,
			"desktop_instance_type": `"${data.alicloud_ecd_bundles.default.bundles.1.desktop_type}"`,
		}),
		FakeConfig: testAccCheckAlicloudEcdImagesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecd_image.default.id}_fake"]`,
			"name_regex": `"${alicloud_ecd_image.default.image_name}_fake"`,
			"image_type": `"SYSTEM"`,
			"status":     `"Creating"`,
		}),
	}
	var existAlicloudEcdImagesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   "1",
			"names.#":                 "1",
			"images.#":                "1",
			"images.0.id":             acctest.CHECKSET,
			"images.0.create_time":    acctest.CHECKSET,
			"images.0.data_disk_size": acctest.CHECKSET,
			"images.0.description":    fmt.Sprintf("tf-testaccimage%d", rand),
			"images.0.gpu_category":   acctest.CHECKSET,
			"images.0.image_name":     fmt.Sprintf("tf-testaccimage%d", rand),
			"images.0.image_id":       acctest.CHECKSET,
			"images.0.image_type":     acctest.CHECKSET,
			"images.0.os_type":        acctest.CHECKSET,
			"images.0.progress":       acctest.CHECKSET,
			"images.0.size":           acctest.CHECKSET,
			"images.0.status":         "Available",
		}
	}
	var fakeAlicloudEcdImagesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"images.#": "0",
		}
	}
	var alicloudEcdImagesBusesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_ecd_images.default",
		ExistMapFunc: existAlicloudEcdImagesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudEcdImagesDataSourceNameMapFunc,
	}

	alicloudEcdImagesBusesCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, imageTypeConf, statusConf, osTypeConf, desktopInstanceTypeConf, allConf)
}
func testAccCheckAlicloudEcdImagesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
  default = "tf-testaccimage%d"
}

resource "alicloud_ecd_simple_office_site" "default" {
  cidr_block          = "172.16.0.0/12"
  desktop_access_type = "Internet"
  office_site_name    = var.name
}

data "alicloud_ecd_bundles" "default" {
  bundle_type = "SYSTEM"
}

resource "alicloud_ecd_policy_group" "default" {
  policy_group_name = var.name
  clipboard         = "readwrite"
  local_drive       = "read"
  authorize_access_policy_rules {
    description = "example_value"
    cidr_ip     = "1.2.3.4/24"
  }
  authorize_security_policy_rules {
    type        = "inflow"
    policy      = "accept"
    description = "example_value"
    port_range  = "80/80"
    ip_protocol = "TCP"
    priority    = "1"
    cidr_ip     = "0.0.0.0/0"
  }
}

resource "alicloud_ecd_desktop" "default" {
	office_site_id  = alicloud_ecd_simple_office_site.default.id
	policy_group_id = alicloud_ecd_policy_group.default.id
	bundle_id 		= data.alicloud_ecd_bundles.default.bundles.1.id
	desktop_name 	= var.name
}

resource "alicloud_ecd_image" "default" {
	image_name =  var.name
	desktop_id =  alicloud_ecd_desktop.default.id
	description = var.name
}


data "alicloud_ecd_images" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

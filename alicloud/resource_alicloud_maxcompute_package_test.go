package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TestAccAliCloudMaxComputePackage_basic covers create -> update body -> import.
func TestAccAliCloudMaxComputePackage_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_maxcompute_package.default"
	ra := resourceAttrInit(resourceId, AliCloudMaxComputePackageMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MaxComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMaxComputePackage")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	checkoutSupportedRegions(t, true, connectivity.MaxComputeProjectSupportRegions)
	name := fmt.Sprintf("tf_testaccmcpkg%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMaxComputePackageBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// create with is_install=true covers package_name, project_name, is_install.
				Config: testAccConfig(map[string]interface{}{
					"project_name": "${alicloud_maxcompute_project.default.project_name}",
					"package_name": "${var.name}_pkg",
					"is_install":   "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				// update body covers the body attribute (update-only).
				Config: testAccConfig(map[string]interface{}{
					"project_name": "${alicloud_maxcompute_project.default.project_name}",
					"package_name": "${var.name}_pkg",
					"is_install":   "true",
					"body":         "updated package content body",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"is_install", "body"},
				Config: testAccConfig(map[string]interface{}{
					"project_name": "${alicloud_maxcompute_project.default.project_name}",
					"package_name": "${var.name}_pkg",
					"is_install":   "true",
					"body":         "updated package content body",
				}),
			},
		},
	})
}

var AliCloudMaxComputePackageMap = map[string]string{
	"project_name": CHECKSET,
	"package_name": CHECKSET,
}

func AliCloudMaxComputePackageBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_maxcompute_project" "default" {
  project_name = "${var.name}_project"
  product_type = "PayAsYouGo"
  is_logical   = "false"
}
`, name)
}

// TestAccAliCloudMaxComputePackagesDataSource covers the packages data source.
func TestAccAliCloudMaxComputePackagesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	checkoutSupportedRegions(t, true, connectivity.MaxComputeProjectSupportRegions)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMaxComputePackagesSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_maxcompute_package.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudMaxComputePackagesSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_maxcompute_package.default.id}_fake"]`,
		}),
	}

	nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMaxComputePackagesSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudMaxComputePackagesSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMaxComputePackagesSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_maxcompute_package.default.id}"]`,
			"name_regex": `"${var.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudMaxComputePackagesSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_maxcompute_package.default.id}_fake"]`,
			"name_regex": `"${var.name}_fake"`,
		}),
	}

	MaxComputePackagesCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameConf, allConf)
}

var existMaxComputePackagesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"packages.#":              "1",
		"packages.0.package_name": CHECKSET,
	}
}

var fakeMaxComputePackagesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"packages.#": "0",
	}
}

var MaxComputePackagesCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_maxcompute_packages.default",
	existMapFunc: existMaxComputePackagesMapFunc,
	fakeMapFunc:  fakeMaxComputePackagesMapFunc,
}

func testAccCheckAlicloudMaxComputePackagesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
    default = "tf_testaccmcpkgds%d"
}

resource "alicloud_maxcompute_project" "default" {
  project_name = "${var.name}_project"
  product_type = "PayAsYouGo"
  is_logical   = "false"
}

resource "alicloud_maxcompute_package" "default" {
  project_name = "${alicloud_maxcompute_project.default.project_name}"
  package_name = "${var.name}_pkg"
  is_install   = true
}

data "alicloud_maxcompute_packages" "default" {
  project_name = "${alicloud_maxcompute_project.default.project_name}"
%s
}
`, rand, indentedLines(pairs))
	return config
}

func indentedLines(lines []string) string {
	var b string
	for _, l := range lines {
		b += "  " + l + "\n"
	}
	return b
}

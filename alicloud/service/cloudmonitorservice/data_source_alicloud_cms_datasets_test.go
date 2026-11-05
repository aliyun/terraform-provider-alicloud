package cloudmonitorservice_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudCmsDatasetsDataSource_nameRegex(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCmsDatasetsDataSourceName(rand, map[string]string{
			"workspace":          `"${alicloud_cms_workspace.default.workspace_name}"`,
			"dataset_name_regex": `"${alicloud_cms_dataset.default.dataset_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudCmsDatasetsDataSourceName(rand, map[string]string{
			"workspace":          `"${alicloud_cms_workspace.default.workspace_name}"`,
			"dataset_name_regex": `"${alicloud_cms_dataset.default.dataset_name}_fake"`,
		}),
	}
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCmsDatasetsDataSourceName(rand, map[string]string{
			"workspace": `"${alicloud_cms_workspace.default.workspace_name}"`,
			"ids":       `["${alicloud_cms_dataset.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudCmsDatasetsDataSourceName(rand, map[string]string{
			"workspace": `"${alicloud_cms_workspace.default.workspace_name}"`,
			"ids":       `["${alicloud_cms_dataset.default.id}_fake"]`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCmsDatasetsDataSourceName(rand, map[string]string{
			"workspace":          `"${alicloud_cms_workspace.default.workspace_name}"`,
			"dataset_name_regex": `"${alicloud_cms_dataset.default.dataset_name}"`,
			"ids":                `["${alicloud_cms_dataset.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudCmsDatasetsDataSourceName(rand, map[string]string{
			"workspace":          `"${alicloud_cms_workspace.default.workspace_name}"`,
			"dataset_name_regex": `"${alicloud_cms_dataset.default.dataset_name}_fake"`,
			"ids":                `["${alicloud_cms_dataset.default.id}_fake"]`,
		}),
	}
	var existAlicloudCmsDatasetsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   "1",
			"datasets.#":              "1",
			"datasets.0.id":           acctest.CHECKSET,
			"datasets.0.dataset_name": acctest.CHECKSET,
			"datasets.0.workspace":    acctest.CHECKSET,
			"datasets.0.create_time":  acctest.CHECKSET,
			"datasets.0.update_time":  acctest.CHECKSET,
			"datasets.0.region_id":    acctest.CHECKSET,
		}
	}
	var fakeAlicloudCmsDatasetsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"datasets.#": "0",
		}
	}
	var alicloudCmsDatasetsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cms_datasets.default",
		ExistMapFunc: existAlicloudCmsDatasetsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudCmsDatasetsDataSourceNameMapFunc,
	}
	alicloudCmsDatasetsCheckInfo.DataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func testAccCheckAlicloudCmsDatasetsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tftestacccmsdataset%d"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
  sls_project    = alicloud_log_project.default.project_name
  workspace_name = var.name
}

resource "alicloud_cms_dataset" "default" {
  workspace    = alicloud_cms_workspace.default.workspace_name
  dataset_name  = var.name
  description   = var.name
  schema        = "{\"metric\":{\"type\":\"text\"}}"
}

data "alicloud_cms_datasets" "default" {
  %s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

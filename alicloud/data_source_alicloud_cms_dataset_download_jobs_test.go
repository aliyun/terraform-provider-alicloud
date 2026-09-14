package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudCmsDatasetDownloadJobsDataSource_nameRegex(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsDatasetDownloadJobsDataSourceName(rand, map[string]string{
			"workspace":      `"${alicloud_cms_workspace.default.workspace_name}"`,
			"dataset_name":   `"${alicloud_cms_dataset.default.dataset_name}"`,
			"job_name_regex": `"${alicloud_cms_dataset_download_job.default.job_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsDatasetDownloadJobsDataSourceName(rand, map[string]string{
			"workspace":      `"${alicloud_cms_workspace.default.workspace_name}"`,
			"dataset_name":   `"${alicloud_cms_dataset.default.dataset_name}"`,
			"job_name_regex": `"${alicloud_cms_dataset_download_job.default.job_name}_fake"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsDatasetDownloadJobsDataSourceName(rand, map[string]string{
			"workspace":    `"${alicloud_cms_workspace.default.workspace_name}"`,
			"dataset_name": `"${alicloud_cms_dataset.default.dataset_name}"`,
			"ids":          `["${alicloud_cms_dataset_download_job.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCmsDatasetDownloadJobsDataSourceName(rand, map[string]string{
			"workspace":    `"${alicloud_cms_workspace.default.workspace_name}"`,
			"dataset_name": `"${alicloud_cms_dataset.default.dataset_name}"`,
			"ids":          `["${alicloud_cms_dataset_download_job.default.id}_fake"]`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsDatasetDownloadJobsDataSourceName(rand, map[string]string{
			"workspace":      `"${alicloud_cms_workspace.default.workspace_name}"`,
			"dataset_name":   `"${alicloud_cms_dataset.default.dataset_name}"`,
			"job_name_regex": `"${alicloud_cms_dataset_download_job.default.job_name}"`,
			"ids":            `["${alicloud_cms_dataset_download_job.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCmsDatasetDownloadJobsDataSourceName(rand, map[string]string{
			"workspace":      `"${alicloud_cms_workspace.default.workspace_name}"`,
			"dataset_name":   `"${alicloud_cms_dataset.default.dataset_name}"`,
			"job_name_regex": `"${alicloud_cms_dataset_download_job.default.job_name}_fake"`,
			"ids":            `["${alicloud_cms_dataset_download_job.default.id}_fake"]`,
		}),
	}
	var existAlicloudCmsDatasetDownloadJobsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"download_jobs.#":                  "1",
			"download_jobs.0.id":               CHECKSET,
			"download_jobs.0.job_name":         CHECKSET,
			"download_jobs.0.dataset_name":     CHECKSET,
			"download_jobs.0.workspace":        CHECKSET,
			"download_jobs.0.compression_type": CHECKSET,
			"download_jobs.0.content_type":     CHECKSET,
			"download_jobs.0.query":            CHECKSET,
			"download_jobs.0.create_time":      CHECKSET,
			"download_jobs.0.update_time":      CHECKSET,
		}
	}
	var fakeAlicloudCmsDatasetDownloadJobsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":           "0",
			"download_jobs.#": "0",
		}
	}
	var alicloudCmsDatasetDownloadJobsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cms_dataset_download_jobs.default",
		existMapFunc: existAlicloudCmsDatasetDownloadJobsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCmsDatasetDownloadJobsDataSourceNameMapFunc,
	}
	alicloudCmsDatasetDownloadJobsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func testAccCheckAlicloudCmsDatasetDownloadJobsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tftestacccmsddj%d"
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
  dataset_name = var.name
  description  = var.name
  schema        = "{\"metric\":{\"type\":\"text\"}}"
}

resource "alicloud_cms_dataset_download_job" "default" {
  workspace        = alicloud_cms_workspace.default.workspace_name
  dataset_name     = alicloud_cms_dataset.default.dataset_name
  job_name         = var.name
  query            = "* | select count(*) as total"
  compression_type = "gzip"
  content_type     = "csv"
}

data "alicloud_cms_dataset_download_jobs" "default" {
  %s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

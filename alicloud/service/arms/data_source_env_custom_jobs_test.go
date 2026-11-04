package arms_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudArmsEnvCustomJobsDataSource_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_arms_env_custom_jobs.default"
	name := fmt.Sprintf("tf-testacc%sarmsenvcustomjob%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceArmsEnvCustomJobsConfig0)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_env_custom_job.default.environment_id}",
			"ids":            []string{"${alicloud_arms_env_custom_job.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_env_custom_job.default.environment_id}",
			"ids":            []string{"${alicloud_arms_env_custom_job.default.id}_fake"},
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_env_custom_job.default.environment_id}",
			"name_regex":     "${alicloud_arms_env_custom_job.default.env_custom_job_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_env_custom_job.default.environment_id}",
			"name_regex":     "${alicloud_arms_env_custom_job.default.env_custom_job_name}_fake",
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_env_custom_job.default.environment_id}",
			"ids":            []string{"${alicloud_arms_env_custom_job.default.id}"},
			"name_regex":     "${alicloud_arms_env_custom_job.default.env_custom_job_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_env_custom_job.default.environment_id}",
			"ids":            []string{"${alicloud_arms_env_custom_job.default.id}_fake"},
			"name_regex":     "${alicloud_arms_env_custom_job.default.env_custom_job_name}_fake",
		}),
	}
	var existAliCloudArmsEnvCustomJobsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                      "1",
			"names.#":                    "1",
			"jobs.#":                     "1",
			"jobs.0.id":                  acctest.CHECKSET,
			"jobs.0.config_yaml":         acctest.CHECKSET,
			"jobs.0.env_custom_job_name": acctest.CHECKSET,
			"jobs.0.environment_id":      acctest.CHECKSET,
			"jobs.0.region_id":           acctest.CHECKSET,
			"jobs.0.status":              acctest.CHECKSET,
		}
	}
	var fakeAliCloudArmsEnvCustomJobsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"jobs.#":  "0",
		}
	}
	var alicloudArmsEnvCustomJobsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_arms_env_custom_jobs.default",
		ExistMapFunc: existAliCloudArmsEnvCustomJobsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAliCloudArmsEnvCustomJobsDataSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudArmsEnvCustomJobsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}

func dataSourceArmsEnvCustomJobsConfig0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	resource "alicloud_arms_environment" "default" {
  		bind_resource_id     = data.alicloud_vpcs.default.ids.0
  		environment_sub_type = "ECS"
  		environment_type     = "ECS"
  		environment_name     = var.name
  		tags = {
    		Created = "TF"
    		For     = "Environment"
  		}
	}

	resource "alicloud_arms_env_custom_job" "default" {
  		status              = "run"
  		environment_id      = alicloud_arms_environment.default.id
  		env_custom_job_name = var.name
  		config_yaml         = <<EOF
scrape_configs:
- job_name: job-demo1
  honor_timestamps: false
  honor_labels: false
  scrape_interval: 30s
  scheme: http
  metrics_path: /metric
  static_configs:
  - targets:
    - 127.0.0.1:9090
EOF
  		aliyun_lang         = "en"
	}
`, name)
}

package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudArmsEnvCustomJobsDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_arms_env_custom_jobs.default"
	name := fmt.Sprintf("tf-testacc%sarmsenvcustomjob%d", defaultRegionToTest, rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceArmsEnvCustomJobsConfig0)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_env_custom_job.default.environment_id}",
			"ids":            []string{"${alicloud_arms_env_custom_job.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_env_custom_job.default.environment_id}",
			"ids":            []string{"${alicloud_arms_env_custom_job.default.id}_fake"},
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_env_custom_job.default.environment_id}",
			"name_regex":     "${alicloud_arms_env_custom_job.default.env_custom_job_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_env_custom_job.default.environment_id}",
			"name_regex":     "${alicloud_arms_env_custom_job.default.env_custom_job_name}_fake",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_env_custom_job.default.environment_id}",
			"ids":            []string{"${alicloud_arms_env_custom_job.default.id}"},
			"name_regex":     "${alicloud_arms_env_custom_job.default.env_custom_job_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
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
			"jobs.0.id":                  CHECKSET,
			"jobs.0.config_yaml":         CHECKSET,
			"jobs.0.env_custom_job_name": CHECKSET,
			"jobs.0.environment_id":      CHECKSET,
			"jobs.0.region_id":           CHECKSET,
			"jobs.0.status":              CHECKSET,
		}
	}
	var fakeAliCloudArmsEnvCustomJobsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"jobs.#":  "0",
		}
	}
	var alicloudArmsEnvCustomJobsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_arms_env_custom_jobs.default",
		existMapFunc: existAliCloudArmsEnvCustomJobsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAliCloudArmsEnvCustomJobsDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudArmsEnvCustomJobsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
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

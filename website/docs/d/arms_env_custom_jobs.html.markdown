---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_env_custom_jobs"
description: |-
  Provides a list of ARMS Env Custom Jobs to the user.
---

# alicloud_arms_env_custom_jobs

This data source provides the ARMS Env Custom Jobs of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.258.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_arms_environment" "default" {
  bind_resource_id     = data.alicloud_vpcs.default.ids.0
  environment_sub_type = "ECS"
  environment_type     = "ECS"
  environment_name     = "${var.name}-${random_integer.default.result}"
  tags = {
    Created = "TF"
    For     = "Environment"
  }
}

resource "alicloud_arms_env_custom_job" "default" {
  status              = "run"
  environment_id      = alicloud_arms_environment.default.id
  env_custom_job_name = "${var.name}-${random_integer.default.result}"
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

data "alicloud_arms_env_custom_jobs" "ids" {
  environment_id = alicloud_arms_env_custom_job.default.environment_id
  ids            = [alicloud_arms_env_custom_job.default.id]
}

output "arms_env_custom_jobs_id_0" {
  value = data.alicloud_arms_env_custom_jobs.ids.jobs.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List) A list of ARMS Env Custom Job IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by ARMS Env Custom Job name.
* `environment_id` - (Required, ForceNew) The ID of the environment instance.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of ARMS Env Custom Job names.
* `jobs` - A list of ARMS Env Custom Jobs. Each element contains the following attributes:
  * `id` - The ID of the custom job. It formats as `<environment_id>:<env_custom_job_name>`.
  * `config_yaml` - The YAML configuration string.
  * `env_custom_job_name` - The name of the custom job.
  * `environment_id` - The ID of the environment instance.
  * `region_id` - The region ID.
  * `status` - The status of the custom job.

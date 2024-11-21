---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_env_custom_job"
description: |-
  Provides a Alicloud ARMS Env Custom Job resource.
---

# alicloud_arms_env_custom_job

Provides a ARMS Env Custom Job resource. Custom jobs in the arms environment.

For information about ARMS Env Custom Job and how to use it, see [What is Env Custom Job](https://www.alibabacloud.com/help/en/arms/developer-reference/api-arms-2019-08-08-createenvcustomjob).

-> **NOTE:** Available since v1.212.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_arms_env_custom_job&exampleId=3e7e4930-e4fe-cc75-64de-f363080256c218ec9dda&activeTab=example&spm=docs.r.arms_env_custom_job.0.3e7e4930e4&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

resource "random_integer" "default" {
  max = 99999
  min = 10000
}

variable "name" {
  default = "terraform-example"
}

resource "alicloud_vpc" "vpc" {
  description = var.name
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name
}

resource "alicloud_arms_environment" "env-ecs" {
  environment_type     = "ECS"
  environment_name     = "terraform-example-${random_integer.default.result}"
  bind_resource_id     = alicloud_vpc.vpc.id
  environment_sub_type = "ECS"
}

resource "alicloud_arms_env_custom_job" "default" {
  status              = "run"
  environment_id      = alicloud_arms_environment.env-ecs.id
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
```

## Argument Reference

The following arguments are supported:
* `aliyun_lang` - (Optional) The locale. The default is Chinese zh | en.
* `config_yaml` - (Required) Yaml configuration string.
* `env_custom_job_name` - (Required, ForceNew) Custom job name.
* `environment_id` - (Required, ForceNew) Environment id.
* `status` - (Optional, Computed) Status: run, stop.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<environment_id>:<env_custom_job_name>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Env Custom Job.
* `delete` - (Defaults to 5 mins) Used when delete the Env Custom Job.
* `update` - (Defaults to 5 mins) Used when update the Env Custom Job.

## Import

ARMS Env Custom Job can be imported using the id, e.g.

```shell
$ terraform import alicloud_arms_env_custom_job.example <environment_id>:<env_custom_job_name>
```
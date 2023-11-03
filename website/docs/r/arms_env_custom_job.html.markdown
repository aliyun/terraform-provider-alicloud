---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_env_custom_job"
description: |-
  Provides a Alicloud ARMS Env Custom Job resource.
---

# alicloud_arms_env_custom_job

Provides a ARMS Env Custom Job resource. Custom jobs in the arms environment.

For information about ARMS Env Custom Job and how to use it, see [What is Env Custom Job](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.212.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "vpc" {
  description = "api-resource-sub-test1-hz-job"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name

}

resource "alicloud_vswitch" "vsw" {
  description  = "api-resource-test1-hz"
  vpc_id       = alicloud_vpc.vpc.id
  vswitch_name = var.name

  zone_id    = data.alicloud_zones.default.zones.0.id
  cidr_block = "172.16.0.0/24"
}

resource "alicloud_ack_cluster" "ask" {
  kubernetes_version = "1.26.3-aliyun.1"
  cluster_type       = "ManagedKubernetes"
  cluster_spec       = "ack.pro.small"
  vpc_id             = alicloud_vpc.vpc.id
  service_cidr       = "192.168.0.0/24"
  cluster_name       = var.name

  container_cidr = "192.168.1.0/24"
  vswitch_id     = alicloud_vswitch.vsw.id
  profile        = "Serverless"
}

resource "alicloud_arms_environment" "env-cs" {
  environment_type = "CS"
  environment_name = var.name

  bind_resource_id     = alicloud_ack_cluster.ask.id
  environment_sub_type = "ACK"
  tags {
    tag_key   = "api-cs-k1"
    tag_value = "api-cs-v1"
  }
}


resource "alicloud_arms_env_custom_job" "default" {
  status              = "run"
  environment_id      = alicloud_arms_environment.env-cs.id
  env_custom_job_name = var.name

  config_yaml = <<EOF
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
- job_name: job-demo2
  honor_timestamps: false
  honor_labels: false
  scrape_interval: 30s
  scheme: http
  metrics_path: /metric
  static_configs:
  - targets:
    - 127.0.0.1:9090
  http_sd_configs:
  - url: 127.0.0.1:9090
    refresh_interval: 30s
EOF
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
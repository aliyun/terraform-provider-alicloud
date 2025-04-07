---
subcategory: "Serverless App Engine (SAE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sae_application_scaling_rule"
sidebar_current: "docs-alicloud-resource-sae-application-scaling-rule"
description: |-
  Provides a Alicloud Serverless App Engine (SAE) Application Scaling Rule resource.
---

# alicloud_sae_application_scaling_rule

Provides a Serverless App Engine (SAE) Application Scaling Rule resource.

For information about Serverless App Engine (SAE) Application Scaling Rule and how to use it, see [What is Application Scaling Rule](https://next.api.aliyun.com/api/sae/2019-05-06/CreateApplicationScalingRule).

-> **NOTE:** Available since v1.159.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_sae_application_scaling_rule&exampleId=159839e1-57f7-f424-eb00-eabd5f79d14f18bc66e1&activeTab=example&spm=docs.r.sae_application_scaling_rule.0.159839e157&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}
variable "name" {
  default = "tf-example"
}
data "alicloud_regions" "default" {
  current = true
}
resource "random_integer" "default" {
  max = 99999
  min = 10000
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}
resource "alicloud_security_group" "default" {
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_sae_namespace" "default" {
  namespace_id              = "${data.alicloud_regions.default.regions.0.id}:example${random_integer.default.result}"
  namespace_name            = var.name
  namespace_description     = var.name
  enable_micro_registration = false
}

resource "alicloud_sae_application" "default" {
  app_description   = var.name
  app_name          = "${var.name}-${random_integer.default.result}"
  namespace_id      = alicloud_sae_namespace.default.id
  image_url         = "registry-vpc.${data.alicloud_regions.default.regions.0.id}.aliyuncs.com/sae-demo-image/consumer:1.0"
  package_type      = "Image"
  security_group_id = alicloud_security_group.default.id
  vpc_id            = alicloud_vpc.default.id
  vswitch_id        = alicloud_vswitch.default.id
  timezone          = "Asia/Beijing"
  replicas          = "5"
  cpu               = "500"
  memory            = "2048"
}

resource "alicloud_sae_application_scaling_rule" "default" {
  app_id                   = alicloud_sae_application.default.id
  scaling_rule_name        = var.name
  scaling_rule_enable      = true
  scaling_rule_type        = "mix"
  min_ready_instances      = "3"
  min_ready_instance_ratio = "-1"
  scaling_rule_timer {
    period = "* * *"
    schedules {
      at_time      = "08:00"
      max_replicas = 10
      min_replicas = 3
    }
    schedules {
      at_time      = "20:00"
      max_replicas = 50
      min_replicas = 3
    }
  }
  scaling_rule_metric {
    max_replicas = 50
    min_replicas = 3
    metrics {
      metric_type                       = "CPU"
      metric_target_average_utilization = 20
    }
    metrics {
      metric_type                       = "MEMORY"
      metric_target_average_utilization = 30
    }
    metrics {
      metric_type                       = "tcpActiveConn"
      metric_target_average_utilization = 20
    }
    scale_up_rules {
      step                         = 10
      disabled                     = false
      stabilization_window_seconds = 0
    }
    scale_down_rules {
      step                         = 10
      disabled                     = false
      stabilization_window_seconds = 10
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required, ForceNew) Application ID.
* `min_ready_instance_ratio` - (Optional) The min ready instance ratio.
* `min_ready_instances` - (Optional) The min ready instances.
* `scaling_rule_enable` - (Optional) True whether the auto scaling policy is enabled. The value description is as follows: true: enabled state. false: disabled status. Valid values: `false`, `true`.
* `scaling_rule_name` - (Required, ForceNew) The name of a custom elastic scaling policy. In the application, the policy name cannot be repeated. It must start with a lowercase letter, and can only contain lowercase letters, numbers, and dashes (-), and no more than 32 characters. After the scaling policy is successfully created, the policy name cannot be modified.
* `scaling_rule_type` - (Required, ForceNew) Flexible strategy type. Valid values: `mix`, `timing` and `metric`.
* `scaling_rule_timer` - (Optional) Configuration of Timing Resilient Policies. See [`scaling_rule_timer`](#scaling_rule_timer) below.
* `scaling_rule_metric` - (Optional) Monitor the configuration of the indicator elasticity strategy. See [`scaling_rule_metric`](#scaling_rule_metric) below.

### `scaling_rule_timer`

The scaling_rule_timer supports the following:

* `begin_date` - (Optional) The Start date. When the `begin_date` and `end_date` values are empty. it indicates long-term execution and is the default value.
* `end_date` - (Optional) The End Date. When the `begin_date` and `end_date` values are empty. it indicates long-term execution and is the default value.
* `period` - (Optional) The period in which a timed elastic scaling strategy is executed.
* `schedules` - (Optional) Resilient Scaling Strategy Trigger Timing. See [`schedules`](#scaling_rule_timer-schedules) below.

### `scaling_rule_timer-schedules`

The schedules supports the following:

* `at_time` - (Optional) Trigger point in time. When supporting format: minutes, for example: `08:00`.
* `target_replicas` - (Optional) This parameter can specify the number of instances to be applied or the minimum number of surviving instances per deployment. value range [1,50]. -> **NOTE:** The attribute is valid when the attribute `scaling_rule_type` is `timing`.
* `max_replicas` - (Optional) Maximum number of instances applied. -> **NOTE:** The attribute is valid when the attribute `scaling_rule_type` is `mix`.
* `min_replicas` - (Optional) Minimum number of instances applied. -> **NOTE:** The attribute is valid when the attribute `scaling_rule_type` is `mix`.

### `scaling_rule_metric`

The scaling_rule_metric supports the following:

* `max_replicas` - (Optional) Maximum number of instances applied.
* `min_replicas` - (Optional) Minimum number of instances applied.
* `metrics` - (Optional) Indicator rule configuration. See [`metrics`](#scaling_rule_metric-metrics) below.
* `scale_up_rules` - (Optional) Apply expansion rules. See [`scale_up_rules`](#scaling_rule_metric-scale_up_rules) below.
* `scale_down_rules` - (Optional) Apply shrink rules. See [`scale_down_rules`](#scaling_rule_metric-scale_down_rules) below.

### `scaling_rule_metric-scale_up_rules`

The scale_up_rules supports the following:

* `step` - (Optional) Elastic expansion or contraction step size. the maximum number of instances to be scaled in per unit time.
* `stabilization_window_seconds` - (Optional) Cooling time for expansion or contraction. Valid values: `0` to `3600`. Unit: seconds. The default is `0` seconds.
* `disabled` - (Optional) Whether shrinkage is prohibited.

### `scaling_rule_metric-scale_down_rules`

The scale_down_rules supports the following:

* `step` - (Optional) Elastic expansion or contraction step size. the maximum number of instances to be scaled in per unit time.
* `stabilization_window_seconds` - (Optional) Cooling time for expansion or contraction. Valid values: `0` to `3600`. Unit: seconds. The default is `0` seconds.
* `disabled` - (Optional) Whether shrinkage is prohibited.

### `scaling_rule_metric-metrics`

The metrics supports the following:

* `metric_target_average_utilization` - (Optional) According to different `metric_type`, set the target value of the corresponding monitoring index.
* `metric_type` - (Optional) Monitoring indicator trigger condition. Valid values: `CPU`, `MEMORY`, `tcpActiveConn`, `QPS`, `RT`, `SLB_QPS`, `SLB_RT`, `INTRANET_SLB_QPS` and `INTRANET_SLB_RT`. The values are described as follows:
  - CPU: CPU usage.
  - MEMORY: MEMORY usage.
  - tcpActiveConn: The average number of TCP active connections for a single instance in 30 seconds.
  - QPS: The average QPS of a single instance within 1 minute of JAVA application.
  - RT: The average response time of all service interfaces within 1 minute of JAVA application.
  - SLB_QPS: The average public network SLB QPS of a single instance within 15 seconds.
  - SLB_RT: The average response time of public network SLB within 15 seconds.
  - INTRANET_SLB_QPS: The average private network SLB QPS of a single instance within 15 seconds.
  - INTRANET_SLB_RT: The average response time of private network SLB within 15 seconds.
**NOTE:** From version 1.206.0, `metric_type` can be set to `QPS`, `RT`, `INTRANET_SLB_QPS`, `INTRANET_SLB_RT`.
* `slb_id` - (Optional, Available in 1.206.0+) SLB ID.
* `slb_project` - (Optional, Available in 1.206.0+) The project of the Log Service.
* `slb_log_store` - (Optional, Available in 1.206.0+) The log store of the Log Service.
* `vport` - (Optional, Available in 1.206.0+) SLB listening port.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Application Scaling Rule. The value formats as `<app_id>:<scaling_rule_name>`.

## Import

Serverless App Engine (SAE) Application Scaling Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_sae_application_scaling_rule.example <app_id>:<scaling_rule_name>
```

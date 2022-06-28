---
subcategory: "Serverless App Engine (SAE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sae_application_scaling_rule"
sidebar_current: "docs-alicloud-resource-sae-application-scaling-rule"
description: |-
  Provides a Alicloud Serverless App Engine (SAE) Application Scaling Rule resource.
---

# alicloud\_sae\_application\_scaling\_rule

Provides a Serverless App Engine (SAE) Application Scaling Rule resource.

For information about Serverless App Engine (SAE) Application Scaling Rule and how to use it, see [What is Application Scaling Rule](https://help.aliyun.com/document_detail/134120.html).

-> **NOTE:** Available in v1.159.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}
resource "alicloud_sae_namespace" "default" {
  namespace_description = "example_value"
  namespace_id          = "example_value"
  namespace_name        = "example_value"
}
resource "alicloud_sae_application" "default" {
  app_description = "example_value"
  app_name        = "example_value"
  namespace_id    = alicloud_sae_namespace.default.namespace_id
  image_url       = "registry-vpc.cn-hangzhou.aliyuncs.com/lxepoo/apache-php5"
  package_type    = "Image"
  jdk             = "Open JDK 8"
  vswitch_id      = data.alicloud_vswitches.default.ids.0
  vpc_id          = data.alicloud_vpcs.default.ids.0
  timezone        = "Asia/Shanghai"
  replicas        = "5"
  cpu             = "500"
  memory          = "2048"
}
resource "alicloud_sae_application_scaling_rule" "example" {
  app_id              = alicloud_sae_application.default.id
  scaling_rule_name   = "example-value"
  scaling_rule_enable = true
  scaling_rule_type   = "mix"
  scaling_rule_timer {
    begin_date = "2022-02-25"
    end_date   = "2022-03-25"
    period     = "* * *"
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
* `scaling_rule_enable` - (Optional, Computed) True whether the auto scaling policy is enabled. The value description is as follows: true: enabled state. false: disabled status. Valid values: `false`, `true`.
* `scaling_rule_name` - (Required, ForceNew) The name of a custom elastic scaling policy. In the application, the policy name cannot be repeated. It must start with a lowercase letter, and can only contain lowercase letters, numbers, and dashes (-), and no more than 32 characters. After the scaling policy is successfully created, the policy name cannot be modified.
* `scaling_rule_type` - (Required, ForceNew) Flexible strategy type. Valid values: `mix`, `timing` and `metric`.
* `scaling_rule_timer` - (Optional) Configuration of Timing Resilient Policies. See the following `Block scaling_rule_timer`.
* `scaling_rule_metric` - (Optional) Monitor the configuration of the indicator elasticity strategy. See the following `Block scaling_rule_metric`.

#### Block scaling_rule_timer

The scaling_rule_timer supports the following:

* `begin_date` - (Optional) The Start date. When the `begin_date` and `end_date` values are empty. it indicates long-term execution and is the default value.
* `end_date` - (Optional) The End Date. When the `begin_date` and `end_date` values are empty. it indicates long-term execution and is the default value.
* `period` - (Optional) The period in which a timed elastic scaling strategy is executed.
* `schedules` - (Optional) Resilient Scaling Strategy Trigger Timing. See the following `Block schedules`.

#### Block schedules

The schedules supports the following:

* `at_time` - (Optional) Trigger point in time. When supporting format: minutes, for example: `08:00`.
* `target_replicas` - (Optional) This parameter can specify the number of instances to be applied or the minimum number of surviving instances per deployment. value range [1,50]. -> **NOTE:** The attribute is valid when the attribute `scaling_rule_type` is `timing`.
* `max_replicas` - (Optional) Maximum number of instances applied. -> **NOTE:** The attribute is valid when the attribute `scaling_rule_type` is `mix`.
* `min_replicas` - (Optional) Minimum number of instances applied. -> **NOTE:** The attribute is valid when the attribute `scaling_rule_type` is `mix`.

#### Block scaling_rule_metric

The scaling_rule_metric supports the following:

* `max_replicas` - (Optional) Maximum number of instances applied.
* `min_replicas` - (Optional) Minimum number of instances applied.
* `metrics` - (Optional) Indicator rule configuration. See the following `Block metrics`.
* `scale_up_rules` - (Optional) Apply expansion rules. See the following `Block scale_up_rules`.
* `scale_down_rules` - (Optional) Apply shrink rules. See the following `Block scale_down_rules`.

#### Block scale_up_rules

The scale_up_rules supports the following:

* `step` - (Optional) Elastic expansion or contraction step size. the maximum number of instances to be scaled in per unit time.
* `stabilization_window_seconds` - (Optional) Cooling time for expansion or contraction. Valid values: `0` to `3600`. Unit: seconds. The default is `0` seconds.
* `disabled` - (Optional) Whether shrinkage is prohibited.

#### Block scale_down_rules

The scale_down_rules supports the following:

* `step` - (Optional) Elastic expansion or contraction step size. the maximum number of instances to be scaled in per unit time.
* `stabilization_window_seconds` - (Optional) Cooling time for expansion or contraction. Valid values: `0` to `3600`. Unit: seconds. The default is `0` seconds.
* `disabled` - (Optional) Whether shrinkage is prohibited.

#### Block metrics

The metrics supports the following:

* `metric_target_average_utilization` - (Optional) According to different `metric_type`, set the target value of the corresponding monitoring index.
* `metric_type` - (Optional) Monitoring indicator trigger condition. Valid values: `CPU`, `MEMORY`, `tcpActiveConn`, `SLB_QPS` and `SLB_RT`. The values are described as follows:
  - CPU: CPU usage. 
  - MEMORY: MEMORY usage.
  - tcpActiveConn: the average number of TCP active connections for a single instance in 30 seconds.
  - SLB_QPS: the average public network SLB QPS of a single instance within 15 seconds.
  - SLB_RT: the average response time of public network SLB within 15 seconds.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Application Scaling Rule. The value formats as `<app_id>:<scaling_rule_name>`.

## Import

Serverless App Engine (SAE) Application Scaling Rule can be imported using the id, e.g.

```
$ terraform import alicloud_sae_application_scaling_rule.example <app_id>:<scaling_rule_name>
```
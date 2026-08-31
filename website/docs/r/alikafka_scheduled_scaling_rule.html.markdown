---
subcategory: "AliKafka"
layout: "alicloud"
page_title: "Alicloud: alicloud_alikafka_scheduled_scaling_rule"
description: |-
  Provides a Alicloud Alikafka Scheduled Scaling Rule resource.
---

# alicloud_alikafka_scheduled_scaling_rule

Provides a Alikafka Scheduled Scaling Rule resource.

Elastic timing strategy.

For information about Alikafka Scheduled Scaling Rule and how to use it, see [What is Scheduled Scaling Rule](https://next.api.alibabacloud.com/document/alikafka/2019-09-16/CreateScheduledScalingRule).

-> **NOTE:** Available since v1.269.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_alikafka_scheduled_scaling_rule&exampleId=19897ca6-5eab-fbc7-b606-401c772a23206986ac55&activeTab=example&spm=docs.r.alikafka_scheduled_scaling_rule.0.19897ca65e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "10.4.0.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_alikafka_instance" "default" {
  deploy_type     = "4"
  instance_type   = "alikafka_serverless"
  vswitch_id      = alicloud_vswitch.default.id
  spec_type       = "normal"
  service_version = "3.3.1"
  security_group  = alicloud_security_group.default.id
  config          = "{\"enable.acl\":\"true\"}"
  serverless_config {
    reserved_publish_capacity   = 60
    reserved_subscribe_capacity = 60
  }
}

resource "alicloud_alikafka_scheduled_scaling_rule" "default" {
  schedule_type        = "repeat"
  reserved_sub_flow    = "200"
  reserved_pub_flow    = "200"
  time_zone            = "GMT+8"
  duration_minutes     = "100"
  first_scheduled_time = "1769578000000"
  enable               = false
  repeat_type          = "Weekly"
  weekly_types         = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"]
  rule_name            = var.name
  instance_id          = alicloud_alikafka_instance.default.id
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_alikafka_scheduled_scaling_rule&spm=docs.r.alikafka_scheduled_scaling_rule.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `duration_minutes` - (Required, ForceNew, Int) The duration (unit: minutes) of a scheduled elastic task.

-> **NOTE:** The parameter value must be at least 15 minutes.

* `enable` - (Optional, Bool) Enables or disables the scheduled task policy. Valid values:
  -`true`: Enables the policy.
  -`false`: Disables the policy.
* `first_scheduled_time` - (Required, ForceNew, Int) The time when the scheduled policy starts to execute.
* `instance_id` - (Required, ForceNew) The instance ID.
* `repeat_type` - (Optional, ForceNew) When `schedule_type` is `repeat`, the parameter is required. Valid values:
  -`Daily`: Daily scheduled task.
  -`Weekly`: Weekly scheduled task.
* `reserved_pub_flow` - (Required, ForceNew, Int) The scheduled elastic reserved production specification (unit: MB/s).
* `reserved_sub_flow` - (Required, ForceNew, Int) The scheduled elastic reserved consumption specification (unit: MB/s).
* `rule_name` - (Required, ForceNew) The name of the scheduled policy rule.
* `schedule_type` - (Required, ForceNew) The schedule type. Valid values:
  - `at`: Scheduled only once.
  - `repeat`: Scheduled repeatedly.
* `time_zone` - (Required, ForceNew) The time zone (Coordinated Universal Time).
* `weekly_types` - (Optional, ForceNew, List) The weekly types. Supports execution on multiple days. When `repeat_type` is set to `Weekly`, you need to input this parameter. Valid values: `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`, `Sunday`.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<instance_id>:<rule_name>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Scheduled Scaling Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Scheduled Scaling Rule.
* `update` - (Defaults to 5 mins) Used when update the Scheduled Scaling Rule.

## Import

Alikafka Scheduled Scaling Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_alikafka_scheduled_scaling_rule.example <instance_id>:<rule_name>
```

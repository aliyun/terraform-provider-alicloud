---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_metric_rule_black_list"
sidebar_current: "docs-alicloud-resource-cms-metric-rule-black-list"
description: |-
  Provides a Alicloud Cloud Monitor Service Metric Rule Black List resource.
---

# alicloud_cms_metric_rule_black_list

Provides a Cloud Monitor Service Metric Rule Black List resource.

For information about Cloud Monitor Service Metric Rule Black List and how to use it, see [What is Metric Rule Black List](https://www.alibabacloud.com/help/en/cloudmonitor/latest/describemetricruleblacklist).

-> **NOTE:** Available since v1.194.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cms_metric_rule_black_list&exampleId=371257e6-128e-ccd7-fe7d-380dbf5a7d95dd42dd75&activeTab=example&spm=docs.r.cms_metric_rule_black_list.0.371257e612&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Instance"
}
data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones.0.id
  cpu_core_count    = 1
  memory_size       = 2
}
data "alicloud_images" "default" {
  name_regex = "^ubuntu_18.*64"
  owners     = "system"
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
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_instance" "default" {
  availability_zone = data.alicloud_zones.default.zones.0.id
  instance_name     = var.name
  image_id          = data.alicloud_images.default.images.0.id
  instance_type     = data.alicloud_instance_types.default.instance_types.0.id
  security_groups   = [alicloud_security_group.default.id]
  vswitch_id        = alicloud_vswitch.default.id
}

resource "alicloud_cms_metric_rule_black_list" "default" {
  instances = [
    "{\"instancceId\":\"${alicloud_instance.default.id}\"}"
  ]
  metrics {
    metric_name = "disk_utilization"
  }
  category                    = "ecs"
  enable_end_time             = 1799443209000
  namespace                   = "acs_ecs_dashboard"
  enable_start_time           = 1689243209000
  metric_rule_black_list_name = var.name
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cms_metric_rule_black_list&spm=docs.r.cms_metric_rule_black_list.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `category` - (Required) Cloud service classification. For example, Redis includes kvstore_standard, kvstore_sharding, and kvstore_splitrw.
* `effective_time` - (Optional) The effective time range of the alert blacklist policy.
* `enable_end_time` - (Optional) The start timestamp of the alert blacklist policy.Unit: milliseconds.
* `enable_start_time` - (Optional) The end timestamp of the alert blacklist policy.Unit: milliseconds.
* `instances` - (Required) The list of instances of cloud services specified in the alert blacklist policy.
* `is_enable` - (Optional) The status of the alert blacklist policy. Value:-true: enabled.-false: disabled.
* `metric_rule_black_list_name` - (Required, ForceNew) The name of the alert blacklist policy.
* `metrics` - (Optional) Monitoring metrics in the instance. See [`metrics`](#metrics) below. 
* `namespace` - (Required, ForceNew) The data namespace of the cloud service.
* `scope_type` - (Optional) The effective range of the alert blacklist policy. Value:-USER: The alert blacklist policy only takes effect in the current Alibaba cloud account.-GROUP: The alert blacklist policy takes effect in the specified application GROUP.
* `scope_value` - (Optional) Application Group ID list. The format is JSON Array.> This parameter is displayed only when 'ScopeType' is 'GROUP.

### `metrics`

The metrics supports the following:
* `metric_name` - (Required) The name of the monitoring indicator.
* `resource` - (Optional) The extended dimension information of the instance. For example, '{"device":"C:"}' indicates that the blacklist policy is applied to all C disks under the ECS instance.



## Attributes Reference

The following attributes are exported:

* `id` - The ID of the blacklist policy.
* `metric_rule_black_list_id` - The ID of the blacklist policy.
* `create_time` - The timestamp for creating an alert blacklist policy.Unit: milliseconds.
* `update_time` - Modify the timestamp of the alert blacklist policy.Unit: milliseconds.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Metric Rule Black List.
* `delete` - (Defaults to 1 mins) Used when delete the Metric Rule Black List.
* `update` - (Defaults to 5 mins) Used when update the Metric Rule Black List.

## Import

Cloud Monitor Service Metric Rule Black List can be imported using the id, e.g.

```shell
$terraform import alicloud_cms_metric_rule_black_list.example <id>
```
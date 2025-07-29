---
subcategory: "Actiontrail"
layout: "alicloud"
page_title: "Alicloud: alicloud_actiontrail_trail"
description: |-
  Provides a Alicloud Actiontrail Trail resource.
---

# alicloud_actiontrail_trail

Provides a Actiontrail Trail resource.

Trail of ActionTrail. After creating a trail, you need to enable the trail through StartLogging.

For information about Actiontrail Trail and how to use it, see [What is Trail](https://www.alibabacloud.com/help/en/actiontrail/latest/api-actiontrail-2020-07-06-createtrail).

-> **NOTE:** Available since v1.95.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_actiontrail_trail&exampleId=6d6e445a-106e-5b0b-8a89-bfec004f510cdfd5729b&activeTab=example&spm=docs.r.actiontrail_trail.0.6d6e445a10&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_regions" "default" {
  current = true
}

data "alicloud_account" "default" {
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_log_project" "default" {
  project_name = "${var.name}-${random_integer.default.result}"
  description  = "tf actiontrail example"
}

data "alicloud_ram_roles" "default" {
  name_regex = "AliyunServiceRoleForActionTrail"
}

resource "alicloud_actiontrail_trail" "default" {
  trail_name         = var.name
  sls_write_role_arn = data.alicloud_ram_roles.default.roles.0.arn
  sls_project_arn    = "acs:log:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:project/${alicloud_log_project.default.project_name}"
}
```

## Argument Reference

The following arguments are supported:

* `event_rw` - (Optional) The read/write type of the events to be delivered. Default value: `All`. Valid values: `Read`, `Write`, `All`.
* `is_organization_trail` - (Optional, ForceNew, Bool) Specifies whether to create a multi-account trail. Default value: `false`. Valid values:
  - `true`: Creates a multi-account trail.
  - `false`: Creates a single-account trail.
* `max_compute_project_arn` - (Optional, Available since v1.256.0) The ARN of the MaxCompute project to which you want to deliver events.
* `max_compute_write_role_arn` - (Optional, Available since v1.256.0) The ARN of the role that is assumed by ActionTrail to deliver events to the MaxCompute project.
* `oss_bucket_name` - (Optional) The OSS bucket to which the trail delivers logs.
* `oss_key_prefix` - (Optional) The prefix of the file name in the OSS bucket to which the trail delivers logs.
* `oss_write_role_arn` - (Optional) The name of the RAM role that the user allows ActionTrail to access OSS service.
* `sls_project_arn` - (Optional) The ARN of the Simple Log Service project to which the trail delivers logs.
* `sls_write_role_arn` - (Optional) The ARN of the role that ActionTrail assumes to deliver operation events to the Simple Log Service project.
* `status` - (Optional) The status of the trail. Default value: `Enable`. Valid values: `Enable`, `Disable`.
* `trail_name` - (Optional, ForceNew, Available since v1.95.0) The name of the trail to be created.
* `trail_region` - (Optional) The region of the trail.
* `name` - (Optional, ForceNew, Deprecated since v1.95.0) Field `name` has been deprecated from provider version 1.95.0. New field `trail_name` instead.
* `role_name` - (Deprecated since v1.118.0) Field `role_name` has been deprecated from provider version 1.118.0.
* `mns_topic_arn` - (Deprecated since v1.118.0) Field `mns_topic_arn` has been deprecated from provider version 1.118.0.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - (Available since v1.256.0) The time when the trail was created.
* `region_id` - (Available since v1.256.0) The home region of the trail.

## Timeouts

-> **NOTE:** Available since v1.256.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Trail.
* `delete` - (Defaults to 5 mins) Used when delete the Trail.
* `update` - (Defaults to 5 mins) Used when update the Trail.

## Import

Actiontrail Trail can be imported using the id, e.g.

```shell
$ terraform import alicloud_actiontrail_trail.example <id>
```

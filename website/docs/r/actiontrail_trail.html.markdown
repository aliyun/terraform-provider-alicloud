---
subcategory: "Actiontrail"
layout: "alicloud"
page_title: "Alicloud: alicloud_actiontrail_trail"
sidebar_current: "docs-alicloud-resource-actiontrail-trail"
description: |-
  Provides Alibaba Cloud ActionTrail Trail Resource
---

# alicloud_actiontrail_trail

Provides a ActionTrail Trail resource. For information about alicloud actiontrail trail and how to use it, see [What is Resource Alicloud ActionTrail Trail](https://www.alibabacloud.com/help/en/actiontrail/latest/api-actiontrail-2020-07-06-createtrail).

-> **NOTE:** Available since v1.95.0.

-> **NOTE:** You can create a trail to deliver events to Log Service, Object Storage Service (OSS), or both. Before you call this operation to create a trail, make sure that the following requirements are met.
- Deliver events to Log Service: A project is created in Log Service.
- Deliver events to OSS: A bucket is created in OSS.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_actiontrail_trail&exampleId=6d6e445a-106e-5b0b-8a89-bfec004f510cdfd5729b&activeTab=example&spm=docs.r.actiontrail_trail.0.6d6e445a10&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

data "alicloud_regions" "example" {
  current = true
}
data "alicloud_account" "example" {}

resource "alicloud_log_project" "example" {
  project_name = "${var.name}-${random_integer.default.result}"
  description  = "tf actiontrail example"
}

data "alicloud_ram_roles" "example" {
  name_regex = "AliyunServiceRoleForActionTrail"
}

resource "alicloud_actiontrail_trail" "example" {
  trail_name         = var.name
  sls_write_role_arn = data.alicloud_ram_roles.example.roles.0.arn
  sls_project_arn    = "acs:log:${data.alicloud_regions.example.regions.0.id}:${data.alicloud_account.example.id}:project/${alicloud_log_project.example.project_name}"
}
```

## Argument Reference

The following arguments are supported:

* `trail_name` - (Optional, ForceNew) The name of the trail to be created, which must be unique for an account.
* `name` - (Optional, ForceNew, Deprecated from v1.95.0) Field `name` has been deprecated from version 1.95.0. Use `trail_name` instead. 
* `event_rw` - (Optional) Indicates whether the event is a read or a write event. Valid values: `Read`, `Write`, and `All`. Default to `Write`.
* `oss_bucket_name` - (Optional) The OSS bucket to which the trail delivers logs. Ensure that this is an existing OSS bucket.
* `role_name` - (Optional, Deprecated from v1.118.0) Field `name` has been deprecated from version 1.118.0.
* `oss_key_prefix` - (Optional) The prefix of the specified OSS bucket name. This parameter can be left empty.
* `sls_project_arn` - (Optional) The unique ARN of the Log Service project. Ensure that `sls_project_arn` is valid .
* `sls_write_role_arn` - (Optional) The unique ARN of the Log Service role.
* `trail_region` - (Optional) The regions to which the trail is applied. Default to `All`.
* `mns_topic_arn` - (Optional, Deprecated from v1.118.0) Field `mns_topic_arn` has been deprecated from version 1.118.0.
* `status` - (Optional) The status of ActionTrail Trail. After creation, tracking is turned on by default, and you can set the status value to `Disable` to turn off tracking. Valid values: `Enable`, `Disable`. Default to `Enable`.
* `oss_write_role_arn` - (Optional) The unique ARN of the Oss role.
* `is_organization_trail` - (Optional, ForceNew) Specifies whether to create a multi-account trail. Valid values:`true`: Create a multi-account trail.`false`: Create a single-account trail. It is the default value.


## Attributes Reference

The following attributes are exported:

* `id` - The id of ActionTrail Trail. The value is the same as trail_name.

## Import

Action trail can be imported using the id or trail_name, e.g.

```shell
$ terraform import alicloud_actiontrail_trail.default abc12345678
```

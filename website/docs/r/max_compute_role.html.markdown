---
subcategory: "Max Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_max_compute_role"
description: |-
  Provides a Alicloud Max Compute Role resource.
---

# alicloud_max_compute_role

Provides a Max Compute Role resource.



For information about Max Compute Role and how to use it, see [What is Role](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.242.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_max_compute_role&exampleId=4f2786ac-3ae8-3169-a7fc-9e255575c8611ccfa197&activeTab=example&spm=docs.r.max_compute_role.0.4f2786ac3a&intl_lang=EN_US" target="_blank">
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

resource "alicloud_maxcompute_project" "default" {
  default_quota = "默认后付费Quota"
  project_name  = var.name
  comment       = var.name
  product_type  = "PayAsYouGo"
}

resource "alicloud_max_compute_role" "default" {
  type         = "admin"
  project_name = alicloud_maxcompute_project.default.id
  policy       = jsonencode({ "Statement" : [{ "Action" : ["odps:*"], "Effect" : "Allow", "Resource" : ["acs:odps:*:projects/project_name/authorization/roles", "acs:odps:*:projects/project_name/authorization/roles/*/*"] }], "Version" : "1" })
  role_name    = "tf_example112"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_max_compute_role&spm=docs.r.max_compute_role.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `policy` - (Optional, JsonString) Policy Authorization
Refer to [Policy-based access control](https://www.alibabacloud.com/help/en/maxcompute/user-guide/policy-based-access-control-1) and [Authorization practices](https://www.alibabacloud.com/help/en/maxcompute/use-cases/authorization-practices)
* `project_name` - (Required, ForceNew) Project name
* `role_name` - (Required, ForceNew) Role Name

-> **NOTE:** At the beginning of a letter, it can contain letters and numbers and can be no more than 64 characters in length.

* `type` - (Required) Role type Valid values: admin/resource

-> **NOTE:** -- management type (admin) role: You can grant management type permissions through Policy. You cannot grant resource permissions to management type roles. You cannot grant management type permissions to management type roles through ACL. -- resource role: you can authorize resource type permissions through Policy or ACL, but cannot authorize management type permissions. For details, see [role-planning](https://www.alibabacloud.com/help/en/maxcompute/user-guide/role-planning)


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<project_name>:<role_name>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Role.
* `delete` - (Defaults to 5 mins) Used when delete the Role.
* `update` - (Defaults to 5 mins) Used when update the Role.

## Import

Max Compute Role can be imported using the id, e.g.

```shell
$ terraform import alicloud_max_compute_role.example <project_name>:<role_name>
```
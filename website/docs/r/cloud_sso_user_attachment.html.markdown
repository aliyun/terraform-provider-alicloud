---
subcategory: "Cloud SSO"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_sso_user_attachment"
description: |-
  Provides a Alicloud Cloud SSO User Attachment resource.
---

# alicloud_cloud_sso_user_attachment

Provides a Cloud SSO User Attachment resource.

For information about Cloud SSO User Attachment and how to use it, see [What is User Attachment](https://www.alibabacloud.com/help/en/cloudsso/latest/api-cloudsso-2021-05-15-addusertogroup).

-> **NOTE:** Available since v1.141.0.

-> **NOTE:** Cloud SSO Only Support `cn-shanghai` And `us-west-1` Region

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_sso_user_attachment&exampleId=abe154a6-161a-2d40-f0ad-0bd2d5ea1df491ef42e8&activeTab=example&spm=docs.r.cloud_sso_user_attachment.0.abe154a616&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-shanghai"
}

data "alicloud_cloud_sso_directories" "default" {
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_cloud_sso_directory" "default" {
  count          = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? 0 : 1
  directory_name = var.name
}

resource "alicloud_cloud_sso_user" "default" {
  directory_id = local.directory_id
  user_name    = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_cloud_sso_group" "default" {
  directory_id = local.directory_id
  group_name   = var.name
  description  = var.name
}

resource "alicloud_cloud_sso_user_attachment" "default" {
  directory_id = local.directory_id
  user_id      = alicloud_cloud_sso_user.default.user_id
  group_id     = alicloud_cloud_sso_group.default.group_id
}

locals {
  directory_id = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? data.alicloud_cloud_sso_directories.default.ids[0] : concat(alicloud_cloud_sso_directory.default.*.id, [""])[0]
}
```

## Argument Reference

The following arguments are supported:
* `directory_id` - (Required, ForceNew) The ID of the directory.
* `group_id` - (Required, ForceNew) The ID of the group.
* `user_id` - (Required, ForceNew) The ID of the user.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of User Attachment. It formats as `<directory_id>:<group_id>:<user_id>`.

## Import

Cloud SSO User Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_sso_user_attachment.example <directory_id>:<group_id>:<user_id>
```

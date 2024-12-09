---
subcategory: "Cloud SSO"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_sso_user"
sidebar_current: "docs-alicloud-resource-cloud-sso-user"
description: |-
  Provides a Alicloud Cloud SSO User resource.
---

# alicloud_cloud_sso_user

Provides a Cloud SSO User resource.

For information about Cloud SSO User and how to use it, see [What is User](https://www.alibabacloud.com/help/en/cloudsso/latest/api-cloudsso-2021-05-15-createuser).

-> **NOTE:** Available since v1.140.0.

-> **NOTE:** Cloud SSO Only Support `cn-shanghai` And `us-west-1` Region

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_sso_user&exampleId=9ac7ba84-6a49-e81e-4d67-1beed262fa42a7ccbe16&activeTab=example&spm=docs.r.cloud_sso_user.0.9ac7ba846a&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
provider "alicloud" {
  region = "cn-shanghai"
}
data "alicloud_cloud_sso_directories" "default" {}

resource "alicloud_cloud_sso_directory" "default" {
  count          = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? 0 : 1
  directory_name = var.name
}

locals {
  directory_id = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? data.alicloud_cloud_sso_directories.default.ids[0] : concat(alicloud_cloud_sso_directory.default.*.id, [""])[0]
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_cloud_sso_user" "default" {
  directory_id = local.directory_id
  user_name    = "${var.name}-${random_integer.default.result}"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of user. The description can be up to `1024` characters long.
* `directory_id` - (Required, ForceNew) The ID of the Directory.
* `display_name` - (Optional) The display name of user. The display name can be up to `256` characters long.
* `email` - (Optional) The User's Contact Email Address. The email can be up to `128` characters long.
* `first_name` - (Optional) The first name of user. The first_name can be up to `64` characters long.
* `last_name` - (Optional) The last name of user. The last_name can be up to `64` characters long.
* `status` - (Optional) The status of user. Valid values: `Disabled`, `Enabled`.
* `user_name` - (Required, ForceNew) The name of user. The name must be `1` to `64` characters in length and can contain letters, digits, at signs (@), periods (.), underscores (_), and hyphens (-).

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of User. The value formats as `<directory_id>:<user_id>`.
* `user_id` - The User ID of the group.

## Import

Cloud SSO User can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_sso_user.example <directory_id>:<user_id>
```

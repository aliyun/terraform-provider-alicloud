---
subcategory: "Cloud SSO"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_sso_user"
description: |-
  Provides a Alicloud Cloud Sso User resource.
---

# alicloud_cloud_sso_user

Provides a Cloud Sso User resource.



For information about Cloud Sso User and how to use it, see [What is User](https://www.alibabacloud.com/help/en/cloudsso/latest/api-cloudsso-2021-05-15-createuser).

-> **NOTE:** Available since v1.140.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_sso_user&exampleId=e8dcdbdd-325b-c5c5-a0dd-81d8b377f00e91ae9a39&activeTab=example&spm=docs.r.cloud_sso_user.0.e8dcdbdd32&intl_lang=EN_US" target="_blank">
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

locals {
  directory_id = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? data.alicloud_cloud_sso_directories.default.ids[0] : concat(alicloud_cloud_sso_directory.default.*.id, [""])[0]
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) The description of the user. The description can be up to 1,024 characters in length.
* `directory_id` - (Required, ForceNew) The ID of the directory.
* `display_name` - (Optional) The display name of the user. The display name can be up to 256 characters in length.
* `email` - (Optional) The email address of the user. The email address must be unique within the directory. The email address can be up to 128 characters in length.
* `first_name` - (Optional) The first name of the user. The first name can be up to 64 characters in length.
* `last_name` - (Optional) The last name of the user. The last name can be up to 64 characters in length.
* `mfa_authentication_settings` - (Optional, Available since v1.262.1) Specifies whether to enable MFA for the user. Default value: `Enabled`. Valid values: `Enabled`, `Disabled`.
* `password` - (Optional, Available since v1.262.1) The new password. The password must contain the following types of characters: uppercase letters, lowercase letters, digits, and special characters. The password must be 8 to 32 characters in length.
* `status` - (Optional, Computed) The status of the user. Default value: `Enabled`. Valid values: `Enabled`, `Disabled`.
* `tags` - (Optional, Map, Available since v1.262.1) The tag of the resource.
* `user_name` - (Required, ForceNew) The username of the user. The username can contain digits, letters, and the following special characters: @_-. The username can be up to 64 characters in length.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<directory_id>:<user_id>`.
* `create_time` - (Available since v1.262.1) The time when the user was created.
* `user_id` - The ID of the user.

## Timeouts

-> **NOTE:** Available since v1.262.1.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the User.
* `delete` - (Defaults to 5 mins) Used when delete the User.
* `update` - (Defaults to 5 mins) Used when update the User.

## Import

Cloud Sso User can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_sso_user.example <directory_id>:<user_id>
```

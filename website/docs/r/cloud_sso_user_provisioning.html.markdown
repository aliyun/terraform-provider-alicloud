---
subcategory: "Cloud SSO"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_sso_user_provisioning"
description: |-
  Provides a Alicloud Cloud SSO User Provisioning resource.
---

# alicloud_cloud_sso_user_provisioning

Provides a Cloud SSO User Provisioning resource.

RAM user synchronization.

For information about Cloud SSO User Provisioning and how to use it, see [What is User Provisioning](https://next.api.alibabacloud.com/document/cloudsso/2021-05-15/CreateUserProvisioning).

-> **NOTE:** Available since v1.260.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_sso_user_provisioning&exampleId=609ae8da-3e36-dc0d-322f-a5d23c5f4e5a0ef86ae9&activeTab=example&spm=docs.r.cloud_sso_user_provisioning.0.609ae8da3e&intl_lang=EN_US" target="_blank">
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

data "alicloud_account" "default" {
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

resource "alicloud_cloud_sso_user_provisioning" "default" {
  description          = "description"
  principal_id         = alicloud_cloud_sso_user.default.user_id
  target_type          = "RD-Account"
  deletion_strategy    = "Keep"
  duplication_strategy = "KeepBoth"
  principal_type       = "User"
  target_id            = data.alicloud_account.default.id
  directory_id         = alicloud_cloud_sso_user.default.directory_id
}

locals {
  directory_id = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? data.alicloud_cloud_sso_directories.default.ids[0] : concat(alicloud_cloud_sso_directory.default.*.id, [""])[0]
}
```

## Argument Reference

The following arguments are supported:
* `deletion_strategy` - (Required) The processing policy for users who have been synchronized when deleting synchronization
* `description` - (Optional) Description of User Synchronization
* `directory_id` - (Required, ForceNew) The ID of the directory to which the synchronization belongs
* `duplication_strategy` - (Required) Processing Policy for Synchronization Conflicts
* `principal_id` - (Required, ForceNew) The ID of the CloudSSO user/group associated with the synchronization.
* `principal_type` - (Required, ForceNew) The ID of the CloudSSO user/group associated with the synchronization.
* `target_id` - (Required, ForceNew) The ID of the destination associated with the synchronization.
* `target_type` - (Required, ForceNew) The target type associated with the synchronization

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<directory_id>:<user_provisioning_id>`.
* `create_time` - The creation time of the synchronization
* `status` - The status of the resource
* `user_provisioning_id` - The first ID of the resource
* `user_provisioning_statistics` - User Provisioning statistics
  * `failed_event_count` - Number of failed events
  * `gmt_latest_sync` - Last Provisioning time

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the User Provisioning.
* `delete` - (Defaults to 5 mins) Used when delete the User Provisioning.
* `update` - (Defaults to 5 mins) Used when update the User Provisioning.

## Import

Cloud SSO User Provisioning can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_sso_user_provisioning.example <directory_id>:<user_provisioning_id>
```

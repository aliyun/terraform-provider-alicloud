---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_handshake_acceptance"
description: |-
  Provides a Alicloud Resource Manager Handshake Acceptance resource.
---

# alicloud_resource_manager_handshake_acceptance

Provides a Resource Manager Handshake Acceptance resource. It is used by the **invited account** to accept an invitation (handshake) to join a resource directory.

For information about Resource Manager Handshake Acceptance and how to use it, see [What is Handshake](https://www.alibabacloud.com/help/en/doc-detail/135287.htm) and the [AcceptHandshake](https://next.api.alibabacloud.com/document/ResourceDirectoryMaster/2022-04-19/AcceptHandshake) API.

-> **NOTE:** This resource must be applied with the credentials of the **invited account**, not the management (master) account that sent the invitation. The invitation itself is created by the management account via `alicloud_resource_manager_handshake`. Use a separate [provider alias](https://developer.hashicorp.com/terraform/language/providers/configuration#alias-multiple-provider-configurations) configured with the invited account's credentials.

-> **NOTE:** Destroying this resource only removes the acceptance record from Terraform state. If the invitation is also managed by `alicloud_resource_manager_handshake` with `target_type = "Account"`, destroying the management-side handshake resource removes the invited cloud account from the resource directory.

-> **NOTE:** Available since v1.284.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_resource_manager_handshake_acceptance&exampleId=0ef84e5d-067a-1c4e-7d0b-9bf85d757c146e5a9370&activeTab=example&spm=docs.r.resource_manager_handshake_acceptance.0.0ef84e5d06&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

variable "region" {
  default = "cn-hangzhou"
}

variable "management_access_key" {
  type      = string
  sensitive = true
}

variable "management_secret_key" {
  type      = string
  sensitive = true
}

variable "invited_access_key" {
  type      = string
  sensitive = true
}

variable "invited_secret_key" {
  type      = string
  sensitive = true
}

variable "invited_account_id" {
  type = string
}

provider "alicloud" {
  alias      = "management"
  region     = var.region
  access_key = var.management_access_key
  secret_key = var.management_secret_key
}

provider "alicloud" {
  alias      = "invited"
  region     = var.region
  access_key = var.invited_access_key
  secret_key = var.invited_secret_key
}

# The management account sends the invitation.
resource "alicloud_resource_manager_handshake" "example" {
  provider = alicloud.management

  target_entity = var.invited_account_id
  target_type   = "Account"
  note          = "test resource manager handshake"
}

# The invited account accepts it.
resource "alicloud_resource_manager_handshake_acceptance" "example" {
  provider     = alicloud.invited
  handshake_id = alicloud_resource_manager_handshake.example.id
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_resource_manager_handshake_acceptance&spm=docs.r.resource_manager_handshake_acceptance.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:
* `handshake_id` - (Required, ForceNew) The ID of the invitation.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource. The value is the same as `handshake_id`.
* `create_time` - The time when the invitation was created. The time is displayed in UTC.
* `expire_time` - The time when the invitation expires. The time is displayed in UTC.
* `invited_account_real_name` - The real-name verification information of the invited account.
* `master_account_id` - The ID of the management account of the resource directory.
* `master_account_name` - The name of the management account of the resource directory.
* `master_account_real_name` - The real-name verification information of the management account.
* `modify_time` - The time when the invitation was modified. The time is displayed in UTC.
* `note` - The note of the invitation.
* `resource_directory_id` - The ID of the resource directory.
* `status` - The status of the invitation. The acceptance resource exists only when the status is `Accepted`.
* `target_entity` - The invited account. The value is an account ID when `target_type` is `Account`, or a logon email address when `target_type` is `Email`.
* `target_type` - The type of the invited account. Valid values: `Account` and `Email`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when accept the Handshake.
* `delete` - (Defaults to 5 mins) Used when delete the Handshake Acceptance.

## Import

Resource Manager Handshake Acceptance can be imported using the id (handshake id), e.g.

```shell
$ terraform import alicloud_resource_manager_handshake_acceptance.example <id>
```

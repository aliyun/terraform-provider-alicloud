---
subcategory: "Max Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_max_compute_role_user_attachment"
description: |-
  Provides a Alicloud Max Compute Role User Attachment resource.
---

# alicloud_max_compute_role_user_attachment

Provides a Max Compute Role User Attachment resource.

Resources associated with a user and a project-level role.

For information about Max Compute Role User Attachment and how to use it, see [What is Role User Attachment](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.242.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_max_compute_role_user_attachment&exampleId=6b4c228c-a906-0297-68af-d4adf3c269ef36980d2f&activeTab=example&spm=docs.r.max_compute_role_user_attachment.0.6b4c228ca9&intl_lang=EN_US" target="_blank">
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

variable "aliyun_user" {
  default = "ALIYUN$openapiautomation@test.aliyunid.com"
}

variable "ram_user" {
  default = "RAM$openapiautomation@test.aliyunid.com:tf-example"
}

variable "ram_role" {
  default = "RAM$openapiautomation@test.aliyunid.com:role/terraform-no-ak-assumerole-no-deleting"
}

variable "role_name" {
  default = "role_project_admin"
}

variable "project_name" {
  default = "default_project_669886c"
}

resource "alicloud_max_compute_role_user_attachment" "default" {
  role_name    = var.role_name
  user         = var.ram_role
  project_name = var.project_name
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_max_compute_role_user_attachment&spm=docs.r.max_compute_role_user_attachment.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `project_name` - (Required, ForceNew) Project Name
* `role_name` - (Required, ForceNew) Role Name, Valid Values: super_administrator, admin, Custom Role

-> **NOTE:** -- super_administrator: the built-in management role of MaxCompute. The Super Administrator of the project has the permission to operate all resources in the project and the management permission. Project owners or users with the Super_Administrator role can assign the Super_Administrator role to other users. -- admin: the built-in management role of MaxCompute, which has the permission to operate all resources in the project and some basic management permissions. Project owners can assign the Admin role to other users. -- Custom role: a role that is not built-in to MaxCompute and needs to be customized. You can refer to the role (starting with role_) definition in DataWorks.

* `user` - (Optional, ForceNew, Computed) Supported input: Alibaba Cloud account, RAM user, and RAM role

-> **NOTE:** -- Alibaba Cloud account: the account registered on the Alibaba Cloud official website. - RAM User: a user created by an Alibaba Cloud account to assist the Alibaba Cloud account to complete data processing. -- RAM role: a RAM role, like a RAM user, is a type of RAM identity. A RAM role is a virtual user that does not have a specific identity authentication key and needs to be played by a trusted entity user for normal use.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<project_name>-<role_name>-<user>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Role User Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Role User Attachment.

## Import

Max Compute Role User Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_max_compute_role_user_attachment.example <project_name>-<role_name>-<user>
```
---
subcategory: "IMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ims_application_assignment"
sidebar_current: "docs-alicloud-resource-ims-application-assignment"
description: |-
  Provides a IMS Application Assignment resource.
---

# alicloud_ims_application_assignment

Provides a IMS Application Assignment resource.

An assignment binds a RAM principal (a RAM user or a RAM role) to an IMS application so that the principal can access the application. When the resource is destroyed, the binding is removed.

For information about IMS Application Assignment and how to use it, see [What is Application Assignment](https://www.alibabacloud.com/help/en/ims/developer-reference/api-ims-2019-08-15-createapplicationassignment).

-> **NOTE:** Available since v2.2.0.

## Example Usage

Basic Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

# The IMS application itself is not managed by Terraform; supply its ID.
variable "ims_application_id" {
  description = "The ID of the IMS application to assign the principal to."
  type        = string
}

data "alicloud_ims_default_domain" "example" {}

resource "alicloud_ims_user" "example" {
  user_principal_name = "terraform-example@${data.alicloud_ims_default_domain.example.default_domain}"
  display_name        = "Terraform Example"
}

resource "alicloud_ims_application_assignment" "example" {
  app_id         = var.ims_application_id
  principal_type = "user"
  principal_name = alicloud_ims_user.example.user_principal_name
}
```

Assign a RAM role to an application by setting `principal_type` to `role` and `principal_name` to the RAM role name.

## Argument Reference

The following arguments are supported:

* `app_id` - (Required, ForceNew) The ID of the application to which the principal is assigned.
* `principal_type` - (Required, ForceNew) The type of the principal. Valid values: `user` and `role`.
* `principal_name` - (Required, ForceNew) The name of the principal. When `principal_type` is `user`, it is the logon name of the RAM user; when `principal_type` is `role`, it is the name of the RAM role.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the IMS Application Assignment. It is formatted as `<app_id>:<principal_type>:<principal_name>`.
* `create_time` - The time when the assignment was created, in RFC 3339 format (UTC).

## Import

IMS Application Assignment can be imported using the composite id formatted as `<app_id>:<principal_type>:<principal_name>`, e.g.

```shell
$ terraform import alicloud_ims_application_assignment.example app-xxxxx:user:terraform-example@onaliyun.com
```

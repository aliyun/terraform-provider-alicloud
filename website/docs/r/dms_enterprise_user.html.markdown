---
subcategory: "DMS Enterprise"
layout: "alicloud"
page_title: "Alicloud: alicloud_dms_enterprise_user"
sidebar_current: "docs-alicloud-resource-dms-enterprise-user"
description: |-
  Provides a DMS Enterprise User resource.
---

# alicloud_dms_enterprise_user

Provides a DMS Enterprise User resource. For information about Alidms Enterprise User and how to use it, see [What is Resource Alidms Enterprise User](https://www.alibabacloud.com/help/en/dms/developer-reference/api-dms-enterprise-2018-11-01-registeruser).

-> **NOTE:** Available since v1.90.0.

## Example Usage

```terraform
variable "name" {
  default = "tfexamplename"
}
resource "alicloud_ram_user" "default" {
  name         = var.name
  display_name = var.name
  mobile       = "86-18688888888"
  email        = "hello.uuu@aaa.com"
  comments     = "example"
}

resource "alicloud_dms_enterprise_user" "default" {
  uid        = alicloud_ram_user.default.id
  user_name  = var.name
  role_names = ["DBA"]
  mobile     = "86-18688888888"
}
```

## Argument Reference

The following arguments are supported:

* `tid` - (Optional) The tenant ID. 
* `uid` - (Required, ForceNew) The Alibaba Cloud unique ID (UID) of the user to add.
* `status` - (Optional) The state of DMS Enterprise User. Valid values: `NORMAL`, `DISABLE`.
* `role_names` - (Optional) The roles that the user plays.
* `nick_name` - (Optional, Deprecated) It has been deprecated from 1.100.0 and use `user_name` instead.
* `user_name` - (Optional, Available in 1.100.0+) The nickname of the user.
* `mobile` - (Optional) The DingTalk number or mobile number of the user.
* `max_result_count` - (Optional) Query the maximum number of rows on the day.
* `max_execute_count` - (Optional) Maximum number of inquiries on the day.
                         
## Attributes Reference

The following attributes are exported:

* `id` - The Alibaba Cloud unique ID of the user. The value is same as the UID.

## Import

DMS Enterprise User can be imported using the id, e.g.

```shell
$ terraform import alicloud_dms_enterprise_user.example 24356xxx
```

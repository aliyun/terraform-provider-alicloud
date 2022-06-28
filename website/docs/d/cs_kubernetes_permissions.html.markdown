---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_kubernetes_permissions"
sidebar_current: "docs-alicloud-datasource-cs-kubernetes-permissions"
description: |-
  Provides a list of Ram user permissions.
---

# alicloud\_cs\_kubernetes\_permissions

This data source provides a list of Ram user permissions.

-> **NOTE:** Available in v1.122.0+.

## Example Usage

```terraform
# Declare the data source
data "alicloud_ram_users" "users_ds" {
  name_regex = "your_user_name"
}

# permissions
data "alicloud_cs_kubernetes_permissions" "default" {
  uid = data.alicloud_ram_users.users_ds.users.0.id
}

output "permissions" {
  value = data.alicloud_cs_kubernetes_permissions.default.permissions
}
```

## Argument Reference

The following arguments are supported.
* `uid` - (Required) The ID of the RAM user. If you want to query the permissions of a RAM role, specify the ID of the RAM role.

## Attributes Reference

* `id` - Resource ID.
* `uid` - The ID of the RAM user. If you want to query the permissions of a RAM role, specify the ID of the RAM role.
* `permissions` - A list of user permission.
  * `resource_id` - The permission settings to manage ACK clusters. 
  * `resource_type` - The authorization type. Valid values `cluster`, `namespace` and `console`.
  * `role_name` - The name of the predefined role. If a custom role is assigned, the value is the name of the assigined custom role.
  * `role_type` - The predefined role. Valid values `admin`,`ops`,`dev`,`restricted` and `custom`.
  * `is_owner` - ndicates whether the permissions are granted to the cluster owner. Valid values `0`, `1`.
  * `is_ram_role` -Indicates whether the permissions are granted to the RAM role. Valid values `0`,`1`.

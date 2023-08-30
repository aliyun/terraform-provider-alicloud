---
subcategory: "MaxCompute"
layout: "alicloud"
page_title: "Alicloud: alicloud_maxcompute_project"
sidebar_current: "docs-alicloud-resource-maxcompute-project"
description: |-
  Provides a Alicloud Max Compute Project resource.
---

# alicloud_maxcompute_project

Provides a Max Compute Project resource.

For information about Max Compute Project and how to use it, see [What is Project](https://www.alibabacloud.com/help/en/maxcompute).

-> **NOTE:** Available since v1.77.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example"
}
resource "alicloud_maxcompute_project" "default" {
  default_quota = "默认后付费Quota"
  project_name  = var.name
  comment       = var.name
  product_type  = "PayAsYouGo"
}
```

## Argument Reference

The following arguments are supported:
* `project_name` - (Required, ForceNew) The name of the project
* `comment` - (Optional) Comments of project
* `default_quota` - (Optional) Default Computing Resource Group
* `ip_white_list` - (Optional) IP whitelist. See [`ip_white_list`](#ip_white_list) below.
* `properties` - (Optional) Project base attributes. See [`properties`](#properties) below.
* `security_properties` - (Optional) Security-related attributes. See [`security_properties`](#security_properties) below.
* `product_type` - (Optional) Quota payment type, support `PayAsYouGo`, `Subscription`, `Dev`.
* `name` - (Removed from v1.196.0) It has been deprecated from provider version 1.110.0 and `project_name` instead.
* `specification_type` - (Removed from v1.196.0)  The type of resource Specification, only `OdpsStandard` supported currently.
* `order_type` - (Removed from v1.196.0) The type of payment, only `PayAsYouGo` supported currently.


### `ip_white_list`

The ip_white_list supports the following:
* `ip_list` - (Optional) Classic network IP white list.
* `vpc_ip_list` - (Optional) VPC network whitelist.

### `properties`

The properties support the following:
* `allow_full_scan` - (Optional) Whether to allow full table scan.
* `enable_decimal2` - (Optional) Whether to turn on Decimal2.0.
* `encryption` - (Optional, ForceNew) Whether encryption is turned on. See [`encryption`](#properties-encryption) below.
* `retention_days` - (Optional) Job default retention time.
* `sql_metering_max` - (Optional) SQL charge limit.
* `table_lifecycle` - (Optional) Life cycle of tables. See [`table_lifecycle`](#properties-table_lifecycle) below.
* `timezone` - (Optional) Project time zone.
* `type_system` - (Optional) Type system.

### `properties-encryption`

The encryption supports the following:
* `algorithm` - (Optional, ForceNew) Algorithm.
* `enable` - (Optional, ForceNew) Whether to open.
* `key` - (Optional, ForceNew) Encryption algorithm key.

### `properties-table_lifecycle`

The table_lifecycle supports the following:
* `type` - (Optional) Life cycle type.
* `value` - (Optional) The value of the life cycle.

### `security_properties`

The security_properties supports the following:
* `enable_download_privilege` - (Optional) Whether to enable download permission check.
* `label_security` - (Optional) Label authorization.
* `object_creator_has_access_permission` - (Optional) Project creator permissions.
* `object_creator_has_grant_permission` - (Optional) Does the project creator have authorization rights.
* `project_protection` - (Optional) Project protection. See [`project_protection`](#security_properties-project_protection) below.
* `using_acl` - (Optional) Whether to turn on ACL.
* `using_policy` - (Optional) Whether to enable Policy.

### `security_properties-project_protection`

The project_protection supports the following:
* `exception_policy` - (Optional) Exclusion policy.
* `protected` - (Optional) Is it turned on.

## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above. The value is the same as `project_name`.
* `owner` - Project owner
* `security_properties` - Security-related attributes
  * `enable_download_privilege` - Whether to enable download permission check.
  * `label_security` - Label authorization.
  * `object_creator_has_access_permission` - Project creator permissions.
  * `object_creator_has_grant_permission` - Does the project creator have authorization rights.
  * `project_protection` - Project protection.
    * `exception_policy` - Exclusion policy.
    * `protected` - Is it turned on.
  * `using_acl` - Whether to turn on ACL.
  * `using_policy` - Whether to enable Policy.
* `status` - The status of the resource
* `type` - Project type

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Project.
* `delete` - (Defaults to 5 mins) Used when delete the Project.
* `update` - (Defaults to 5 mins) Used when update the Project.
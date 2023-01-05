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

For information about Max Compute Project and how to use it, see [What is Project](https://help.aliyun.com/document_detail/473237.html).

-> **NOTE:** Available in v1.77.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_maxcompute_project" "default" {
  default_quota = "默认后付费Quota"
  project_name  = "test_create_spec_one"
  comment       = "test_for_terraform"
  product_type  = "PAYASYOUGO"
}
```

## Argument Reference

The following arguments are supported:
* `project_name` - (Required, ForceNew) The name of the project
* `comment` - (Optional) Comments of project
* `default_quota` - (Optional) Default Computing Resource Group
* `ip_white_list` - (Computed,Optional) IP whitelistSee the following `Block IpWhiteList`.
* `properties` - (Computed,Optional) Project base attributesSee the following `Block Properties`.
* `security_properties` - (Computed,Optional) Security-related attributesSee the following `Block SecurityProperties`.
* `product_type` - (Optional) Quota payment type, support `PayAsYouGo`, `Subscription`, `Dev`.
* `name` - (Optional, Remove from v1.196.0+) It has been deprecated from provider version 1.110.0 and `project_name` instead.
* `specification_type` - (Optional, Remove from v1.196.0+)  The type of resource Specification, only `OdpsStandard` supported currently.
* `order_type` - (Optional, Remove from v1.196.0+) The type of payment, only `PayAsYouGo` supported currently.


#### Block IpWhiteList

The IpWhiteList supports the following:
* `ip_list` - (Computed,Optional) Classic network IP white list.
* `vpc_ip_list` - (Computed,Optional) VPC network whitelist.

#### Block Encryption

The Encryption supports the following:
* `algorithm` - (Optional, ForceNew, Computed) Algorithm.
* `enable` - (Optional, ForceNew, Computed) Whether to open.
* `key` - (Optional, ForceNew, Computed) Encryption algorithm key.

#### Block Properties

The Properties support the following:
* `allow_full_scan` - (Computed,Optional) Whether to allow full table scan.
* `enable_decimal2` - (Computed,Optional) Whether to turn on Decimal2.0.
* `encryption` - (Optional, ForceNew, Computed) Whether encryption is turned on.See the following `Block Encryption`.
* `retention_days` - (Computed,Optional) Job default retention time.
* `sql_metering_max` - (Optional, Computed) SQL charge limit.
* `table_lifecycle` - (Computed,Optional) Life cycle of tables.See the following `Block TableLifecycle`.
* `timezone` - (Computed,Optional) Project time zone.
* `type_system` - (Optional, Computed) Type system.

#### Block TableLifecycle

The TableLifecycle supports the following:
* `type` - (Computed,Optional) Life cycle type.
* `value` - (Computed,Optional) The value of the life cycle.

#### Block SecurityProperties

The SecurityProperties supports the following:
* `enable_download_privilege` - (Computed,Optional) Whether to enable download permission check.
* `label_security` - (Computed,Optional) Label authorization.
* `object_creator_has_access_permission` - (Computed,Optional) Project creator permissions.
* `object_creator_has_grant_permission` - (Computed,Optional) Does the project creator have authorization rights.
* `project_protection` - (Computed,Optional) Project protection.See the following `Block ProjectProtection`.
* `using_acl` - (Computed,Optional) Whether to turn on ACL.
* `using_policy` - (Computed,Optional) Whether to enable Policy.

#### Block ProjectProtection

The ProjectProtection supports the following:
* `exception_policy` - (Computed,Optional) Exclusion policy.
* `protected` - (Computed,Optional) Is it turned on.

## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above. The value is the same as `project_name`.
* `ip_white_list` - IP whitelist
  * `ip_list` - Classic network IP white list.
  * `vpc_ip_list` - VPC network whitelist.
* `owner` - Project owner
* `properties` - Project base attributes
  * `allow_full_scan` - Whether to allow full table scan.
  * `enable_decimal2` - Whether to turn on Decimal2.0.
  * `retention_days` - Job default retention time.
  * `table_lifecycle` - Life cycle of tables.
    * `type` - Life cycle type.
    * `value` - The value of the life cycle.
  * `timezone` - Project time zone.
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

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Project.
* `delete` - (Defaults to 5 mins) Used when delete the Project.
* `update` - (Defaults to 5 mins) Used when update the Project.
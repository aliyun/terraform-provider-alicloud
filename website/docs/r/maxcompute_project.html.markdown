---
subcategory: "Max Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_maxcompute_project"
description: |-
  Provides a Alicloud Max Compute Project resource.
---

# alicloud_maxcompute_project

Provides a Max Compute Project resource.

MaxCompute project .

For information about Max Compute Project and how to use it, see [What is Project](https://www.alibabacloud.com/help/en/maxcompute/).

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
* `comment` - (Optional, ForceNew) Project description information. The length is 1 to 256 English or Chinese characters. The default value is blank.
* `default_quota` - (Optional) Used to implement computing resource allocation. Valid values: subQuota Nickname
If the calculation Quota is not specified, the default Quota resource will be consumed by jobs initiated by the project. For more information about computing resource usage, see [Computing Resource Usage](https://www.alibabacloud.com/help/en/maxcompute/user-guide/use-of-computing-resources).
* `ip_white_list` - (Optional, List) IP whitelist See [`ip_white_list`](#ip_white_list) below.
* `is_logical` - (Optional) Whether to logically delete. Default value: true. Value: (ture/false),

-> **NOTE:** -- ture: In this case, the project status will be changed to' deleting' and completely deleted after 14 days. -- false: delete immediately, that is, completely deleted and permanently irrecoverable.

* `project_name` - (Optional, ForceNew, Computed) The name begins with a letter, containing letters, digits, and underscores (_). It can be 3 to 28 characters in length and is globally unique.
* `properties` - (Optional, List) Project base attributes See [`properties`](#properties) below.
* `security_properties` - (Optional, List) Security-related attributes See [`security_properties`](#security_properties) below.
* `status` - (Optional, Computed) The project status. Default value: AVAILABLE. Value: (AVAILABLE/READONLY/FROZEN/DELETING)
* `tags` - (Optional, Map) The tag of the resource

### `ip_white_list`

The ip_white_list supports the following:
* `ip_list` - (Optional) Set the IP address whitelist in the classic network. Only devices in the whitelist are allowed to access the project.

-> **NOTE:** If you only configure a classic network IP address whitelist, access to the classic network is restricted and all access to the VPC is prohibited.

* `vpc_ip_list` - (Optional) Set the IP address whitelist in the VPC network to allow only devices in the whitelist to access the project space.

-> **NOTE:** If you only configure a VPC network IP address whitelist, access to the VPC network is restricted and access to the classic network is prohibited.


### `properties`

The properties supports the following:
* `allow_full_scan` - (Optional) Whether to allow full table scan. Default: false
* `enable_decimal2` - (Optional) Whether to turn on Decimal2.0
* `encryption` - (Optional, List) Storage encryption. For details, see [Storage Encryption](https://www.alibabacloud.com/help/en/maxcompute/security-and-compliance/storage-encryption)
  -> **NOTE :**:
To enable storage encryption, you need to modify the parameters of the basic attributes of the MaxCompute project. This operation permission is authenticated by RAM, and you need to have the Super_Administrator role permission of the corresponding project.

To configure the permissions and IP whitelist parameters of the MaxCompute project, you must have the management permissions (Admin) of the corresponding project, including Super_Administrator, Admin, or custom management permissions. For more information, see the project management permissions list.

You can turn on storage encryption only for projects that have not turned on storage encryption. For projects that have turned on storage encryption, you cannot turn off storage encryption or change the encryption algorithm. See [`encryption`](#properties-encryption) below.
* `retention_days` - (Optional, Int) Set the number of days to retain backup data. During this time, you can restore the current version to any backup version. The value range of days is [0,30], and the default value is 1. 0 means backup is turned off.
The effective policy after adjusting the backup cycle is:
Extend the backup cycle: The new backup cycle takes effect on the same day.
Shorten the backup cycle: The system will automatically delete backup data that has exceeded the retention cycle.
* `sql_metering_max` - (Optional) Set the maximum threshold for single SQL Consumption, that is, set the ODPS. SQL. metering.value.max attribute. For more information, see [Consumption control](https://www.alibabacloud.com/help/en/maxcompute/product-overview/consumption-controll).
Unit: scan volume (GB)* complexity.
* `table_lifecycle` - (Optional, List) Set whether the lifecycle of the table in the project needs to be configured, that is, set the ODPS. table.lifecycle property, See [`table_lifecycle`](#properties-table_lifecycle) below.
* `timezone` - (Optional) Project time zone, example value: Asia/Shanghai
* `type_system` - (Optional) Data type version. Value:(1/2/hive)
1: The original MaxCompute type system.
2: New type system introduced by MaxCompute 2.0.
hive: the type system of the Hive compatibility mode introduced by MaxCompute 2.0.

### `properties-encryption`

The properties-encryption supports the following:
* `algorithm` - (Optional) The encryption algorithm supported by the key, including AES256, AESCTR, and RC4.
* `enable` - (Optional) Only enable function is supported. Value: (true)

-> **NOTE:** cannot be turned off after the function is turned on

* `key` - (Optional) The encryption algorithm Key, the Key type used by the project, including the Default Key (MaxCompute Default Key) and the self-contained Key (BYOK). The MaxCompute Default Key is the Default Key created inside MaxCompute.

### `properties-table_lifecycle`

The properties-table_lifecycle supports the following:
* `type` - (Optional) Optional: When creating a table, the Lifecycle clause is optional. If the Lifecycle of the table is not set, the table is permanently valid.
mandatory: the Lifecycle clause is required. You must set the Lifecycle of the table.
inherit: If you do not set the lifecycle of the table when creating a table, the lifecycle of the table is the value of ODPS. table.lifecycle.value, and the ODPS. table.lifecycle.value attribute sets the lifecycle of the table.
* `value` - (Optional) The value of the life cycle, in days. The value range is 1~37231, and the default value is 37231.

### `security_properties`

The security_properties supports the following:
* `enable_download_privilege` - (Optional) Set whether to enable the [Download permission control function](https://www.alibabacloud.com/help/en/maxcompute/user-guide/download-control), that is, set the ODPS. security.enabledownloadprivilege property.
* `label_security` - (Optional) Set whether to use the [Label permission control function](https://www.alibabacloud.com/help/en/maxcompute/user-guide/label-based-access-control), that is, set the LabelSecurity attribute, which is not used by default.
* `object_creator_has_access_permission` - (Optional) Sets whether to allow the creator of the object to have access to the object, I .e. sets the attribute. The default is the allowed state.
* `object_creator_has_grant_permission` - (Optional) The ObjectCreatorHasGrantPermission attribute is set to allow the object creator to have the authorization permission on the object. The default is the allowed state.
* `project_protection` - (Optional, List) Project protection See [`project_protection`](#security_properties-project_protection) below.
* `using_acl` - (Optional) Set whether to use the [ACL permission control function](https://www.alibabacloud.com/help/en/maxcompute/user-guide/maxcompute-permissions), that is, set the CheckPermissionUsingACL attribute, which is in use by default.
* `using_policy` - (Optional) Set whether to use the Policy permission control function (https://www.alibabacloud.com/help/en/maxcompute/user-guide/policy-based-access-control-1), that is, set the CheckPermissionUsingACL attribute, which is in use by default.

### `security_properties-project_protection`

The security_properties-project_protection supports the following:
* `exception_policy` - (Optional, JsonString) Set [Exceptions or Trusted Items](https://www.alibabacloud.com/help/en/maxcompute/security-and-compliance/project-data-protection)
* `protected` - (Optional) Whether enabled, value:(true/false)

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Represents the creation time of the project
* `owner` - Project owner
* `region_id` - The region ID of the resource
* `type` - Project type

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Project.
* `delete` - (Defaults to 5 mins) Used when delete the Project.
* `update` - (Defaults to 5 mins) Used when update the Project.

## Import

Max Compute Project can be imported using the id, e.g.

```shell
$ terraform import alicloud_maxcompute_project.example <id>
```
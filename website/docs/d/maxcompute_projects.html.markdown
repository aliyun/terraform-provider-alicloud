---
subcategory: "Max Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_maxcompute_projects"
sidebar_current: "docs-alicloud-datasource-maxcompute-projects"
description: |-
  Provides a datasource of Max Compute Project owned by an Alibaba Cloud account.
---

# alicloud_maxcompute_projects

This data source provides Max Compute Project available to the user.[What is Project](https://www.alibabacloud.com/help/en/maxcompute/)

-> **NOTE:** Available since v1.196.0.

## Example Usage

```terraform
variable "name" {
  default = "tf_example_acc"
}

resource "alicloud_maxcompute_project" "default" {
  default_quota = "默认后付费Quota"
  project_name  = var.name
  comment       = var.name
  product_type  = "PayAsYouGo"
}

data "alicloud_maxcompute_projects" "default" {
  name_regex = alicloud_maxcompute_project.default.project_name
}

output "alicloud_maxcompute_project_example_id" {
  value = data.alicloud_maxcompute_projects.default.projects.0.project_name
}
```

## Argument Reference

The following arguments are supported:
* `enable_details` - (Optional) Whether to query the detailed attributes of each Project, such as `properties` and `security_properties`. Default value: `false`.
* `ids` - (Optional, ForceNew, Computed) A list of Project IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group Metric Rule name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `product_type` - (Optional) Quota payment type. Valid values: `PayAsYouGo`, `Subscription`, `Dev`.


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Project IDs.
* `names` - A list of name of Projects.
* `projects` - A list of Project Entries. Each element contains the following attributes:

### `projects`

The projects supports the following attributes:
* `comment` - Project description information. The length is 1 to 256 English or Chinese characters. The default value is blank.
* `cost_storage` - View the current storage size of the Project. The storage size is the same as the measurement size, that is, the compressed logical storage size collected by the Project.
* `create_time` - Represents the creation time of the project
* `default_quota` - Used to implement computing resource allocation. If the calculation Quota is not specified, the default Quota resource will be consumed by jobs initiated by the project. For more information about computing resource usage, see [Computing Resource Usage](https://www.alibabacloud.com/help/en/maxcompute/user-guide/use-of-computing-resources).
* `id` - The ID of the Project.
* `ip_white_list` - IP whitelist.
* `owner` - Project owner
* `project_name` - The name begins with a letter, containing letters, digits, and underscores (_). It can be 3 to 28 characters in length and is globally unique.
* `properties` - Project base attributes.
* `region_id` - The region ID of the Project.
* `security_properties` - Security-related attributes.
* `status` - The project status. Default value: AVAILABLE. Value: (AVAILABLE/READONLY/FROZEN/DELETING)
* `tags` - The tags of the Project.
* `three_tier_model` - Whether data storage by schema is supported.
* `trusted_projects` - The list of trusted projects.
* `type` - Project type

### `projects.ip_white_list`

The ip_white_list supports the following attributes:
* `ip_list` - Set the IP address whitelist in the classic network. Only devices in the whitelist are allowed to access the project.-> **NOTE:** If you only configure a classic network IP address whitelist, access to the classic network is restricted and all access to the VPC is prohibited.
* `vpc_ip_list` - Set the IP address whitelist in the VPC network to allow only devices in the whitelist to access the project space.-> **NOTE:** If you only configure a VPC network IP address whitelist, access to the VPC network is restricted and access to the classic network is prohibited.

### `projects.properties`

The properties supports the following attributes:
* `allow_full_scan` - Whether to allow full table scan. Default: false.
* `enable_data_masking` - Whether to enable data masking. Available since v1.287.0.
* `enable_decimal2` - Whether to turn on Decimal2.0.
* `enable_dr` - Whether multi-AZ storage disaster tolerance is enabled. Valid values: true/false.
* `enable_tunnel_quota_route` - Whether to enable data tunnel service resource group routing. Available since v1.287.0.
* `encryption` - Storage encryption. For details, see [Storage Encryption](https://www.alibabacloud.com/help/en/maxcompute/security-and-compliance/storage-encryption).
* `retention_days` - Set the number of days to retain backup data. During this time, you can restore the current version to any backup version. The value range of days is [0,30], and the default value is 1. 0 means backup is turned off. The effective policy after adjusting the backup cycle is: Extend the backup cycle: The new backup cycle takes effect on the same day. Shorten the backup cycle: The system will automatically delete backup data that has exceeded the retention cycle.
* `sql_metering_max` - Set the maximum threshold of single SQL consumption, that is, set the ODPS. SQL. metering.value.max attribute. For details, see [Consumption Monitoring Alarm](https://www.alibabacloud.com/help/en/maxcompute/product-overview/consumption-control). Unit: scan volume (GB)* complexity.
* `table_lifecycle` - Set whether the lifecycle of the table in the project needs to be configured, that is, set the ODPS. table.lifecycle property.
* `timezone` - Project time zone, example value: Asia/Shanghai.
* `tunnel_quota` - The data tunnel service resource group bound to the project. Available since v1.287.0.
* `type_system` - Data type version. Value:(1/2/hive)1: The original MaxCompute type system.2: New type system introduced by MaxCompute 2.0.hive: the type system of the Hive compatibility mode introduced by MaxCompute 2.0.

### `projects.properties.encryption`

The encryption supports the following attributes:
* `algorithm` - The encryption algorithm supported by the key, including AES256, AESCTR, and RC4.
* `enable` - Only enable function is supported. Value: (true).
* `key` - The encryption algorithm Key, the Key type used by the project, including the Default Key (MaxCompute Default Key) and the self-contained Key (BYOK). The MaxCompute Default Key is the Default Key created inside MaxCompute.

### `projects.properties.table_lifecycle`

The table_lifecycle supports the following attributes:
* `type` - Optional: When creating a table, the Lifecycle clause is optional. If the Lifecycle of the table is not set, the table is permanently valid. mandatory: the Lifecycle clause is required. You must set the Lifecycle of the table. inherit: If you do not set the lifecycle of the table when creating a table, the lifecycle of the table is the value of ODPS. table.lifecycle.value, and the ODPS. table.lifecycle.value attribute sets the lifecycle of the table.
* `value` - The value of the life cycle, in days. The value range is 1~37231, and the default value is 37231.

### `projects.security_properties`

The security_properties supports the following attributes:
* `enable_download_privilege` - Set whether to enable the [Download permission control function](https://www.alibabacloud.com/help/en/maxcompute/user-guide/download-control), that is, set the ODPS. security.enabledownloadprivilege property.
* `label_security` - Set whether to use the [Label permission control function](https://www.alibabacloud.com/help/en/maxcompute/user-guide/label-based-access-control), that is, set the LabelSecurity attribute, which is not used by default.
* `object_creator_has_access_permission` - Sets whether to allow the creator of the object to have access to the object, I .e. sets the attribute. The default is the allowed state.
* `object_creator_has_grant_permission` - The ObjectCreatorHasGrantPermission attribute is set to allow the object creator to have the authorization permission on the object. The default is the allowed state.
* `project_protection` - Project protection.
* `using_acl` - Set whether to use the [ACL permission control function](https://www.alibabacloud.com/help/en/maxcompute/user-guide/maxcompute-permissions), that is, set the CheckPermissionUsingACL attribute, which is in use by default.
* `using_policy` - Set whether to use the [Policy permission control function](https://www.alibabacloud.com/help/en/maxcompute/user-guide/policy-based-access-control-1), that is, set the CheckPermissionUsingACL attribute, which is in use by default.

### `projects.security_properties.project_protection`

The project_protection supports the following attributes:
* `exception_policy` - Set [Exceptions or Trusted Items](https://www.alibabacloud.com/help/en/maxcompute/security-and-compliance/project-data-protection).
* `protected` - Whether enabled, value:(true/false).

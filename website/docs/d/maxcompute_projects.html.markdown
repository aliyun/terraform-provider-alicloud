---
subcategory: "Max Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_maxcompute_projects"
sidebar_current: "docs-alicloud-datasource-maxcompute-projects"
description: |-
  Provides a list of Max Compute Project owned by an Alibaba Cloud account.
---

# alicloud_maxcompute_projects

This data source provides Max Compute Project available to the user.[What is Project](https://help.aliyun.com/document_detail/473479.html)

-> **NOTE:** Available since v.1.196.0.

## Example Usage

```
variable "name" {
  default = "tf_testaccmp"
}

resource "alicloud_maxcompute_project" "default" {
  default_quota = "默认后付费Quota"
  project_name  = var.name
  comment       = var.name
  product_type  = "PAYASYOUGO"
}

data "alicloud_maxcompute_projects" "default" {
  ids        = ["${alicloud_maxcompute_project.default.id}"]
  name_regex = alicloud_maxcompute_project.default.name
}

output "alicloud_maxcompute_project_example_id" {
  value = data.alicloud_maxcompute_projects.default.projects.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, ForceNew, Computed) A list of Project IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group Metric Rule name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Project IDs.
* `names` - A list of name of Projects.
* `projects` - A list of Project Entries. Each element contains the following attributes:
    * `id` - Project ID. The value is the same as `project_name`.
    * `default_quota` - Default Computing Resource Group
    * `ip_white_list` - IP whitelist
        * `ip_list` - Classic network IP white list.
        * `vpc_ip_list` - VPC network whitelist.
    * `owner` - Project owner
    * `project_name` - The name of the resource
    * `properties` - Project base attributes
        * `allow_full_scan` - Whether to allow full table scan.
        * `enable_decimal2` - Whether to turn on Decimal2.0.
        * `encryption` - Whether encryption is turned on.
            * `algorithm` - Algorithm.
            * `enable` - Whether to open.
            * `key` - Encryption algorithm key.
        * `retention_days` - Job default retention time.
        * `sql_metering_max` - SQL charge limit.
        * `table_lifecycle` - Life cycle of tables.
            * `type` - Life cycle type.
            * `value` - The value of the life cycle.
        * `timezone` - Project time zone.
        * `type_system` - Type system.
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

---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_data_service_groups"
sidebar_current: "docs-alicloud-datasource-data-works-data-service-groups"
description: |-
  Provides a list of Data Works Data Service Group owned by an Alibaba Cloud account.
---

# alicloud_data_works_data_service_groups

This data source provides Data Works Data Service Group available to the user.[What is Data Service Group](https://next.api.alibabacloud.com/document/dataworks-public/2020-05-18/CreateDataServiceGroup)

-> **NOTE:** Available since v1.289.0.

## Example Usage

```terraform
variable "project_id" {
  description = "The ID of the DataWorks workspace"
  type        = number
}

data "alicloud_data_works_data_service_groups" "default" {
  project_id = var.project_id
}

output "first_group_id" {
  value = data.alicloud_data_works_data_service_groups.default.groups.0.id
}
```

## Argument Reference

The following arguments are supported:
* `project_id` - (Required, ForceNew) The ID of the DataWorks workspace.
* `tenant_id` - (ForceNew, Optional, Deprecated) The tenant ID. This field is deprecated.
* `ids` - (Optional, Computed) A list of Data Service Group IDs. Each element is formulated as `<project_id>:<data_service_group_id>`.
* `name_regex` - (Optional) A regex string to filter results by Data Service Group name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Data Service Group IDs.
* `names` - A list of name of Data Service Groups.
* `groups` - A list of Data Service Group Entries. Each element contains the following attributes:
  * `api_gateway_group_id` - The ID of the API gateway group to which the business process belongs.
  * `create_time` - The time when the business process was created.
  * `creator_id` - The ID of the user who created the business process.
  * `data_service_group_id` - The ID of the business process.
  * `data_service_group_name` - The name of the business process.
  * `description` - The description of the business process.
  * `modified_time` - The time when the business process was last modified.
  * `project_id` - The ID of the DataWorks workspace to which the business process belongs.
  * `tenant_id` - The tenant ID.
  * `id` - The ID of the Data Service Group. The value is formulated as `<project_id>:<data_service_group_id>`.

---
subcategory: "EDAS"
layout: "alicloud"
page_title: "Alicloud: alicloud_edas_deploy_groups"
sidebar_current: "docs-alicloud-datasource-edas-deploy-groups"
description: |-
    Provides a list of EDAS deploy groups available to the user.
---

# alicloud\_edas\_deploy\_groups

This data source provides a list of EDAS deploy groups in an Alibaba Cloud account according to the specified filters.

-> **NOTE:** Available in 1.82.0+

## Example Usage

```
data "alicloud_edas_deploy_groups" "groups" {
  app_id = "xxx"
  ids = ["xxx"]
  output_file = "groups.txt"
}

output "first_group_name" {
  value = data.alicloud_edas_deploy_groups.groups.groups.0.group_name
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required) ID of the EDAS application.
* `ids` - (Optional) An ids string to filter results by the deploy group id. 
* `name_regex` - (Optional) A regex string to filter results by the deploy group name. 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of deploy group IDs.
* `names` - A list of deploy group names.
* `groups` - A list of consumer group ids.
  * `group_id` - The ID of the instance group.
  * `group_name` - The name of the instance group. The length cannot exceed 64 characters.
  * `group_type` - The type of the instance group. Valid values: 0: Default group. 1: Phased release is disabled for traffic management. 2: Phased release is enabled for traffic management.
  * `create_time` - The time when the instance group was created.
  * `update_time` - The time when the instance group was updated.
  * `app_id` - The ID of the application that you want to deploy.
  * `cluster_id` - The ID of the cluster that you want to create the application.
  * `package_version_id` - The version of the deployment package for the instance group that was created.
  * `app_version_id` - The version of the deployment package for the application.
  

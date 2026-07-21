---
subcategory: "Cloud SSO"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_sso_access_assignments"
sidebar_current: "docs-alicloud-datasource-cloud-sso-access-assignments"
description: |-
  Provides a list of Cloud Sso Access Assignments to the user.
---

# alicloud\_cloud\_sso\_access\_assignments

This data source provides the Cloud Sso Access Assignments of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.193.0+.

-> **NOTE:** Cloud SSO Only Support `cn-shanghai` And `us-west-1` Region

## Example Usage

Basic Usage

```terraform
data "alicloud_cloud_sso_access_assignments" "ids" {
  directory_id = "example_value"
  ids          = ["example_value-1", "example_value-2"]
}
output "cloud_sso_access_assignment_id_1" {
  value = data.alicloud_cloud_sso_access_assignments.ids.assignments.0.id
}
```

## Argument Reference

The following arguments are supported:

* `access_configuration_id` - (Optional, ForceNew) Access configuration ID.
* `directory_id` - (Required, ForceNew) Directory ID.
* `ids` - (Optional, ForceNew, Computed)  A list of Access Assignment IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `principal_type` - (Optional, ForceNew) Create the identity type of the access assignment, which can be a user or a user group. Valid values: `Group`, `User`.
* `target_id` - (Optional, ForceNew) The ID of the target to create the resource range.
* `target_type` - (Optional, ForceNew) The type of the resource range target to be accessed. Only a single RD primary account or member account can be specified in the first phase. Valid values: `RD-Account`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `assignments` - A list of Cloud Sso Access Assignments. Each element contains the following attributes:
	* `access_configuration_id` - Access configuration ID.
	* `access_configuration_name` - The name of the access configuration.
	* `create_time` - The creation time of the resource.
	* `directory_id` - Directory ID.
	* `id` - The ID of the Access Assignment.
	* `principal_id` - The ID of the access assignment.
	* `principal_name` - Cloud SSO identity name.
	* `principal_type` - Create the identity type of the access assignment, which can be a user or a user group.
	* `target_id` - The ID of the target to create the resource range.
	* `target_name` - Task target name.
	* `target_path_name` - The path name of the task target in the resource directory.
	* `target_type` - The type of the resource range target to be accessed. Only a single RD primary account or member account can be specified in the first phase.
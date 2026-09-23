---
subcategory: "Cloud SSO"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_sso_groups"
sidebar_current: "docs-alicloud-datasource-cloud-sso-groups"
description: |-
  Provides a list of Cloud Sso Groups to the user.
---

# alicloud\_cloud\_sso\_groups

This data source provides the Cloud Sso Groups of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.138.0+.

-> **NOTE:** Cloud SSO Only Support `cn-shanghai` And `us-west-1` Region

## Example Usage

Basic Usage

```terraform
data "alicloud_cloud_sso_groups" "ids" {
  directory_id = "example_value"
  ids          = ["example_value-1", "example_value-2"]
}
output "cloud_sso_group_id_1" {
  value = data.alicloud_cloud_sso_groups.ids.groups.0.id
}

data "alicloud_cloud_sso_groups" "nameRegex" {
  directory_id = "example_value"
  name_regex   = "^my-Group"
}
output "cloud_sso_group_id_2" {
  value = data.alicloud_cloud_sso_groups.nameRegex.groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `directory_id` - (Required, ForceNew)  The ID of the Directory.
* `ids` - (Optional, ForceNew, Computed)  A list of Group IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `provision_type` - (Optional, ForceNew) The ProvisionType of the Group. Valid values: `Manual`, `Synchronized`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Group names.
* `groups` - A list of Cloud Sso Groups. Each element contains the following attributes:
	* `create_time` - The Created Time of the Directory.
	* `description` - The Description of the Directory.
	* `directory_id` -  The ID of the Directory.
	* `group_id` - The Group ID of the group.
	* `group_name` - The Name of the group.
	* `id` - The ID of the Group.
	* `provision_type` - The Provision Type of the Group. Valid values: `Manual`, `Synchronized`.

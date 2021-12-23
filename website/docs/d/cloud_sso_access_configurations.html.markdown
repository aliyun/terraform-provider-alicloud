---
subcategory: "Cloud SSO"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_sso_access_configurations"
sidebar_current: "docs-alicloud-datasource-cloud-sso-access-configurations"
description: |-
  Provides a list of Cloud Sso Access Configurations to the user.
---

# alicloud\_cloud\_sso\_access\_configurations

This data source provides the Cloud Sso Access Configurations of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.140.0+.

-> **NOTE:** Cloud SSO Only Support `cn-shanghai` And `us-west-1` Region

## Example Usage

Basic Usage

```terraform
data "alicloud_cloud_sso_access_configurations" "ids" {
  directory_id = "example_value"
  ids          = ["example_value-1", "example_value-2"]
}
output "cloud_sso_access_configuration_id_1" {
  value = data.alicloud_cloud_sso_access_configurations.ids.configurations.0.id
}

data "alicloud_cloud_sso_access_configurations" "nameRegex" {
  directory_id = "example_value"
  name_regex   = "^my-AccessConfiguration"
}
output "cloud_sso_access_configuration_id_2" {
  value = data.alicloud_cloud_sso_access_configurations.nameRegex.configurations.0.id
}

```

## Argument Reference

The following arguments are supported:

* `directory_id` - (Required, ForceNew) The ID of the Directory.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Access Configuration IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Access Configuration name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Access Configuration names.
* `configurations` - A list of Cloud Sso Access Configurations. Each element contains the following attributes:
	* `access_configuration_id` - The AccessConfigurationId of the Access Configuration.
	* `access_configuration_name` - The AccessConfigurationName of the Access Configuration.
	* `create_time` - The Created Time of the Directory.
	* `description` - The Description of the Directory.
	* `directory_id` - The ID of the Directory.
	* `id` - The ID of the Access Configuration.
	* `permission_policies` - The Policy List.
		* `add_time` - The Creation time of policy.
		* `permission_policy_document` - The Content of Policy.
		* `permission_policy_name` - The Policy Name of policy.
		* `permission_policy_type` - The Policy Type of policy. Valid values: `System`, `Inline`.
	* `relay_state` - The RelayState of the Access Configuration.
	* `session_duration` - The SessionDuration of the Access Configuration.
	* `status_notifications` - The StatusNotifications of the Access Configuration.

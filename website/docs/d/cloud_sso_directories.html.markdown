---
subcategory: "Cloud SSO"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_sso_directories"
sidebar_current: "docs-alicloud-datasource-cloud-sso-directories"
description: |-
  Provides a list of Cloud Sso Directories to the user.
---

# alicloud\_cloud\_sso\_directories

This data source provides the Cloud Sso Directories of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.135.0+.

-> **NOTE:** Cloud SSO Only Support `cn-shanghai` And `us-west-1` Region

## Example Usage

Basic Usage

```terraform
data "alicloud_cloud_sso_directories" "ids" {
  ids = ["example_id"]
}
output "cloud_sso_directory_id_1" {
  value = data.alicloud_cloud_sso_directories.ids.directories.0.id
}

data "alicloud_cloud_sso_directories" "nameRegex" {
  name_regex = "^my-Directory"
}
output "cloud_sso_directory_id_2" {
  value = data.alicloud_cloud_sso_directories.nameRegex.directories.0.id
}

```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Directory IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Directory name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Directory names.
* `directories` - A list of Cloud Sso Directories. Each element contains the following attributes:
    * `create_time` - The CreateTime of the CloudSSO directory.
    * `directory_id` - The DirectoryId of the CloudSSO directory.
    * `directory_name` - The name of the CloudSSO directory.
    * `id` - The ID of the Directory.
    * `mfa_authentication_status` - The mfa authentication status. Valid values: `Enabled` or `Disabled`. Default to `Disabled`.
    * `region` - The Region of the CloudSSO directory.
    * `saml_identity_provider_configuration` - The saml identity provider configuration.
        * `create_time` - Saml identifies the creation time of the provider configuration.
        * `encoded_metadata_document` - Base64 encoded IdP metadata document.
        * `entity_id` - SAML IdPEntityID.
        * `login_url` - SAML IdP http-post Binding address.
        * `sso_status` - SAML SSO login enabled status. Valid values: `Enabled` or `Disabled`. Default to `Disabled`.
    * `scim_synchronization_status` - The scim synchronization status. Valid values: `Enabled` or `Disabled`. Default to `Disabled`.
    * `tasks` - Asynchronous Task Information Array.
        * `end_time` - The End Time of Task.
        * `failure_reason` - the Reason for the Failure of  the task.
        * `principal_name` - The Name of Cloud SSO Identity.
        * `start_time` - The Start Time of Task.
        * `principal_id` - The ID of Cloud SSO Identity.
        * `target_path` - The Path in RD of Deploy Target.
        * `task_type` - The Type of the Task.
        * `principal_type` - The Type of Cloud SSO Identity.
        * `target_id` - The Id of deploy target.
        * `target_name` - The Name of Deploy Target.
        * `target_type` - The Type of Deploy Target.
        * `access_configuration_id` - The ID of Access Configuration.
        * `access_configuration_name` - The Name of Access Configuration.
        * `status` - The Task Status.
        * `task_id` - The ID of the Task.

---
subcategory: "Cloud SSO"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_sso_users"
sidebar_current: "docs-alicloud-datasource-cloud-sso-users"
description: |-
  Provides a list of Cloud Sso Users to the user.
---

# alicloud\_cloud\_sso\_users

This data source provides the Cloud Sso Users of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.140.0+.

-> **NOTE:** Cloud SSO Only Support `cn-shanghai` And `us-west-1` Region

## Example Usage

Basic Usage

```terraform
data "alicloud_cloud_sso_users" "ids" {
  directory_id = "example_value"
  ids          = ["example_value-1", "example_value-2"]
}
output "cloud_sso_user_id_1" {
  value = data.alicloud_cloud_sso_users.ids.users.0.id
}

data "alicloud_cloud_sso_users" "nameRegex" {
  directory_id = "example_value"
  name_regex   = "^my-User"
}
output "cloud_sso_user_id_2" {
  value = data.alicloud_cloud_sso_users.nameRegex.users.0.id
}

data "alicloud_cloud_sso_users" "provisionType" {
  directory_id   = "example_value"
  ids            = ["example_value-1"]
  provision_type = "Manual"
}
output "cloud_sso_user_id_3" {
  value = data.alicloud_cloud_sso_users.provisionType.users.0.id
}

data "alicloud_cloud_sso_users" "status" {
  directory_id = "example_value"
  ids          = ["example_value-1"]
  status       = "Enabled"
}
output "cloud_sso_user_id_4" {
  value = data.alicloud_cloud_sso_users.status.users.0.id
}

```

## Argument Reference

The following arguments are supported:

* `directory_id` - (Required, ForceNew) The ID of the Directory.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of User IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by User name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `provision_type` - (Optional, ForceNew) ProvisionType. Valid values: `Manual`, `Synchronized`.
* `status` - (Optional, ForceNew) The status of user. Valid values: `Disabled`, `Enabled`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of User names.
* `users` - A list of Cloud Sso Users. Each element contains the following attributes:
    * `create_time` - The create time of the user.
    * `description` - The description of user.
    * `directory_id` - The ID of the Directory.
    * `display_name` - The display name of user.
    * `email` - The User's Contact Email Address.
    * `first_name` - The first name of user.
    * `id` - The ID of the User.
    * `last_name` - The last name of user.
    * `mfa_devices` - The List of MFA Device for User.
      * `device_id` - The MFA Device ID.
      * `device_name` - The MFA Device Name.
      * `device_type` - The MFA Device Type.
      * `effective_time` - The Effective Time of MFA Device.
    * `provision_type` - ProvisionType.
    * `status` - User status. Valid values: `Enabled` and `Disabled`.
    * `user_id` - The User ID of the group.
    * `user_name` - The name of user.

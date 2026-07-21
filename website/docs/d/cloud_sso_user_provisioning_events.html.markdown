---
subcategory: "Cloud SSO"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_sso_user_provisioning_events"
sidebar_current: "docs-alicloud-datasource-cloud-sso-user-provisioning-events"
description: |-
  Provides a list of Cloud Sso User Provisioning Event owned by an Alibaba Cloud account.
---

# alicloud_cloud_sso_user_provisioning_events

This data source provides Cloud Sso User Provisioning Event available to the user.[What is User Provisioning Event](https://next.api.alibabacloud.com/document/cloudsso/2021-05-15/GetUserProvisioningEvent)

-> **NOTE:** Available since v1.261.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-shanghai"
}

resource "alicloud_cloud_sso_directory" "defaultQSrGmc" {
  directory_global_access_status = "Disabled"
  password_policy {
    min_password_length          = "8"
    min_password_different_chars = "8"
    max_password_age             = "90"
    password_reuse_prevention    = "1"
    max_login_attempts           = "5"
  }
  mfa_authentication_setting_info {
    mfa_authentication_advance_settings = "OnlyRiskyLogin"
    operation_for_risk_login            = "EnforceVerify"
  }
  directory_name = "tfexample"
}


data "alicloud_cloud_sso_user_provisioning_events" "default" {
  directory_id = alicloud_cloud_sso_directory.defaultQSrGmc.id
}

output "alicloud_cloud_sso_user_provisioning_event_example_id" {
  value = data.alicloud_cloud_sso_user_provisioning_events.default.events.0.id
}
```

## Argument Reference

The following arguments are supported:
* `directory_id` - (Required, ForceNew) Directory ID
* `user_provisioning_id` - (ForceNew, Optional) The ID of the User Provisioning.
* `ids` - (Optional, ForceNew, Computed) A list of User Provisioning Event IDs. The value is formulated as `<directory_id>:<event_id>`.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of User Provisioning Event IDs.
* `events` - A list of User Provisioning Event Entries. Each element contains the following attributes:
  * `content` - Event content
  * `create_time` - The creation time of the resource
  * `deletion_strategy` - Processing policy when you delete a RAM user
  * `directory_id` - Directory ID
  * `duplication_strategy` - Conflict strategy
  * `error_count` - Number of manual retry failures
  * `error_info` - Error message for last failure
  * `event_id` - Dead letter event ID
  * `last_sync_time` - Last synchronization time
  * `principal_id` - User Provisioning body ID
  * `principal_name` - User Provisioning body name
  * `principal_type` - User Provisioning body type
  * `source_type` - The type of the source action that triggered the event.
  * `target_id` - User Provisioning target ID
  * `target_name` - User Provisioning target name
  * `target_path` - RD path of User Provisioning target
  * `target_type` - User Provisioning target type
  * `update_time` - Event update time
  * `user_provisioning_id` - The ID of the User Provisioning.
  * `id` - The ID of the resource supplied above.

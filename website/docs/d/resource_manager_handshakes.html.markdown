---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_handshakes"
sidebar_current: "docs-alicloud-datasource-resource-manager-handshakes"
description: |-
    Provides a list of Resource Manager Handshakes to the user.
---

# alicloud\_resource\_manager\_handshakes

This data source provides the Resource Manager Handshakes of the current Alibaba Cloud user.

-> **NOTE:**  Available in 1.86.0+.

## Example Usage

```terraform
data "alicloud_resource_manager_handshakes" "example" {}

output "first_handshake_id" {
  value = "${data.alicloud_resource_manager_handshakes.example.handshakes.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Resource Manager Handshake IDs.
* `status` - (Optional, ForceNew, Available in v1.114.0+) The status of handshake, valid values: `Accepted`, `Cancelled`, `Declined`, `Deleted`, `Expired` and `Pending`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `enable_details` -(Optional, Available in v1.114.0+) Default to `false`. Set it to true can output more details.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Resource Manager Handshake IDs.
* `handshakes` - A list of Resource Manager Handshakes. Each element contains the following attributes:
    * `id` - The ID of the resource.
    * `handshake_id`- The ID of the invitation.
    * `expire_time` - The time when the invitation expires.
    * `master_account_id` - The ID of the master account of the resource directory.
    * `master_account_name` - The name of the master account of the resource directory.
    * `modify_time` - The time when the invitation was modified.
    * `note` - The invitation note.
    * `resource_directory_id` - The ID of the resource directory.
    * `status` - The status of the invitation.
    * `target_entity` - The ID or logon email address of the invited account.
    * `target_type` - The type of the invited account. 
    * `invited_account_real_name` - (Available in v1.114.0+) The real name of the invited account. 
    * `master_account_real_name` - (Available in v1.114.0+) The real name of the master account of the resource directory. 
    

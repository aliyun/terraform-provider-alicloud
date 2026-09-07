---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_instance_members"
sidebar_current: "docs-alicloud-datasource-cloud_firewall-instance-members"
description: |-
  Provides a list of Cloud Firewall Instance Member owned by an Alibaba Cloud account.
---

# alicloud_cloud_firewall_instance_members

This data source provides Cloud Firewall Instance Member available to the user.[What is Instance Member](https://help.aliyun.com/document_detail/261237.html)

-> **NOTE:** Available since v1.194.0.

## Example Usage

```terraform
data "alicloud_cloud_firewall_instance_members" "default" {
  ids = ["${alicloud_cloud_firewall_instance_member.default.id}"]
}

output "alicloud_cloud_firewall_instance_member_example_id" {
  value = data.alicloud_cloud_firewall_instance_members.default.members.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, ForceNew, Computed) A list of Instance Member IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Instance Member IDs.
* `members` - A list of Instance Member Entries. Each element contains the following attributes:
  * `id` - The instance id.
  * `create_time` - When the cloud firewall member account was added.> use second-level timestamp format.
  * `member_desc` - Remarks of cloud firewall member accounts.
  * `member_display_name` - The name of the cloud firewall member account.
  * `member_uid` - The UID of the cloud firewall member account.
  * `modify_time` - The last modification time of the cloud firewall member account.> use second-level timestamp format.
  * `status` - The resource attribute field that represents the resource status.

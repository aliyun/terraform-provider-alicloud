---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_instance_member"
sidebar_current: "docs-alicloud-resource-cloud_firewall-instance-member"
description: |-
  Provides a Alicloud Cloud Firewall Instance Member resource.
---

# alicloud_cloud_firewall_instance_member

Provides a Cloud Firewall Instance Member resource.

For information about Cloud Firewall Instance Member and how to use it, see [What is Instance Member](https://www.alibabacloud.com/help/en/server-load-balancer/latest/createloadbalancer).

-> **NOTE:** Available in v1.194.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "AliyunTerraform"
}

resource "alicloud_resource_manager_account" "default" {
  display_name = var.name
}

resource "alicloud_cloud_firewall_instance_member" "default" {
  member_desc = var.name
  member_uid  = alicloud_resource_manager_account.default.id
}
```

## Argument Reference

The following arguments are supported:
* `member_desc` - (Optional) Remarks of cloud firewall member accounts.
* `member_uid` - (Required,ForceNew) The UID of the cloud firewall member account.



## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above. Its value same as `member_uid`.
* `create_time` - When the cloud firewall member account was added.> use second-level timestamp format.
* `member_display_name` - The name of the cloud firewall member account.
* `modify_time` - The last modification time of the cloud firewall member account.> use second-level timestamp format.
* `status` - The resource attribute field that represents the resource status.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Instance Member.
* `delete` - (Defaults to 1 mins) Used when delete the Instance Member.
* `update` - (Defaults to 5 mins) Used when update the Instance Member.

## Import

Cloud Firewall Instance Member can be imported using the id, e.g.

```shell
$terraform import alicloud_cloud_firewall_instance_member.example <id>
```
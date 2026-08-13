---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_security_groups"
sidebar_current: "docs-alicloud-datasource-ens-security-groups"
description: |-
  Provides a list of ENS Security Groups to the user.
---

# alicloud\_ens\_security_groups

This data source provides the ENS Security Groups of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.287.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_ens_security_groups" "default" {
  security_group_name = "tf-example"
  ids                 = ["sg-xxx"]
}
output "first_group_id" {
  value = data.alicloud_ens_security_groups.default.groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` to get more details of the security groups.
* `ids` - (Optional, ForceNew) A list of Security Group IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Security Group name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `security_group_id` - (Optional, ForceNew) The ID of the Security Group.
* `security_group_name` - (Optional, ForceNew) The name of the Security Group.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of Security Group IDs.
* `names` - A list of Security Group names.
* `groups` - A list of Security Groups.

### `groups`

The groups supports the following attributes:
* `create_time` - Creation time of the security group, UTC time.
* `description` - Security group description information.
* `id` - The ID of the Security Group.
* `instance_count` - Number of instances associated with a security group.
* `permissions` - A collection of rules for a security group instance.
* `security_group_id` - The ID of the Security Group.
* `security_group_name` - Security group name.

### `groups.permissions`

The permissions supports the following attributes:
* `creation_time` - Creation time, UTC time.
* `description` - Rule description information.
* `dest_cidr_ip` - Destination IP address segment for outbound authorization.
* `direction` - Authorized direction.
* `ip_protocol` - IP protocol.
* `ipv6_dest_cidr_ip` - The target IPv6 address segment.
* `ipv6_source_cidr_ip` - The source IPv6 address segment.
* `policy` - Authorization Policy.
* `port_range` - Source end port range.
* `priority` - Rule Priority.
* `source_cidr_ip` - Source IP address segment, used for inbound authorization.
* `source_port_range` - The port range of the source security group.

---
layout: "alicloud"
page_title: "Alicloud: alicloud_security_group_rules"
sidebar_current: "docs-alicloud-datasource-security-group-rules"
description: |-
    Provides a collection of Security Group Rules available to the user.
---

# alicloud\_security\_group\_rules

The `alicloud_security_group_rules` data source provides a collection of security permissions of a specific security group.
Each collection item represents a single `ingress` or `egress` permission rule.
The ID of the security group can be provided via a variable or the result from the other data source `alicloud_security_groups`.

## Example Usage

The following example shows how to obtain details about a security group rule and how to pass its data to an instance at launch time.

```
# Get the security group id from a variable
variable "security_group_id" {}

# Or get it from the alicloud_security_groups data source.
# Please note that the data source arguments must be enough to filter results to one security group.
data "alicloud_security_groups" "groups_ds" {
  name_regex = "api"
}

# Filter the security group rule by group
data "alicloud_security_group_rules" "ingress_rules_ds" {
  group_id = "${data.alicloud_security_groups.groups_ds.groups.0.id}" # or ${var.security_group_id}
  nic_type = "internet"
  direction = "ingress"
  ip_protocol = "TCP"
}

# Pass port_range to the backend service
resource "alicloud_instance" "backend" {
  # ...
  user_data = "config_service.sh --portrange=${data.alicloud_security_group_rules.ingress_rules_ds.rules.0.port_range}"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required) The ID of the security group that owns the rules.
* `nic_type` - (Optional) Refers to the network type. Can be either `internet` or `intranet`. The default value is `internet`.
* `direction` - (Optional) Authorization direction. Valid values are: `ingress` or `egress`.
* `ip_protocol` - (Optional) The IP protocol. Valid values are: `tcp`, `udp`, `icmp`, `gre` and `all`.
* `policy` - (Optional) Authorization policy. Can be either `accept` or `drop`. The default value is `accept`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `rules` - A list of rules. Each element contains the following attributes:
  * `group_name` - The name of the security group that owns the rules.
  * `group_desc` - The description of the security group that owns the rules.
  * `rules` - A list of security group rules. Each element contains the following attributes:
    * `ip_protocol` - The protocol. Can be `tcp`, `udp`, `icmp`, `gre` or `all`.
    * `port_range` - The range of port numbers.
    * `source_cidr_ip` - Source IP address segment for ingress authorization.
    * `source_security_group_id` - Source security group ID for ingress authorization.
    * `source_group_owner_account` - Alibaba Cloud account of the source security group.
    * `dest_cidr_ip` - Target IP address segment for egress authorization.
    * `dest_security_group_id` - Target security group id for ingress authorization.
    * `dest_group_owner_account` - Alibaba Cloud account of the target security group.
    * `policy` - Authorization policy. Can be either `accept` or `drop`.
    * `nic_type` - Network type, `internet` or `intranet`.
    * `priority` - Rule priority.
    * `direction` - Authorization direction, `ingress` or `egress`.
    * `description` - The description of the rule.

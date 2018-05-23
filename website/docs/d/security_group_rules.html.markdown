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
The id of the security group can be provided via variable or filtered by another data source `alicloud_security_groups`.

## Example Usage

The following example shows how to obtain details of the security group rule and passing the data to the instance at launch.

```
# accept a security group id as a variable

variable "security_group_id" {}

# or filter using data source
# note the filter must select only one specific group

data "alicloud_security_groups" "api" {
  name_regex = "api"
}

# filter the rule

data "alicloud_security_group_rules" "ingress" {
  id          = "${alicloud_security_groups.api.0.id}"
                # or ${var.security_group_id}
  nic_type    = "internet"
  direction   = "ingress"
  ip_protocol = "TCP"
}

# pass port_range to the backend service

resource "alicloud_instance" "backend" {
  ...
  user_data = "config_service.sh --portrange=${data.alicloud_security_group_rules.ingress.0.port_range}"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required) The id of security group wich owns the rules.
* `nic_type` - (Optional) Refers to the network type. Can be either `internet` or `intranet`. The default value is `internet`.
* `direction` - (Optional) Authorization direction, `ingress` or `egress`.
* `ip_protocol` - (Optional) The protocol. Can be `tcp`, `udp`, `icmp`, `gre` or `all`.
* `policy` - (Optional) Authorization policy. Can be either `accept` or `drop`. The default value is `accept`.
* `output_file` - (Optional) The name of file that can save security group rules after running `terraform plan`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `group_name` - The name of the security group which owns the rules.
* `group_desc` - The description of the security group which owns the rules.
* `rules` - A list of security group rules. Its every element contains the following attributes:
  * `ip_protocol` - The protocol. Can be `tcp`, `udp`, `icmp`, `gre` or `all`.
  * `port_range` - The range of port numbers.
  * `source_cidr_ip` - Source ip address segment for ingress authorization.
  * `source_security_group_id` - Source security group id for ingress authorization.
  * `source_group_owner_account` - Alibaba Cloud account of the source security group.
  * `dest_cidr_ip` - Target ip address segment for egress authorization.
  * `dest_security_group_id` - Target security group id for ingress authorization.
  * `dest_group_owner_account` - Alibaba Cloud account of the target security group.
  * `policy` - Authorization policy. Can be either `accept` or `drop`.
  * `nic_type` - Network type, `internet` or `intranet`.
  * `priority` - Rule priority.
  * `direction` - Authorization direction, `ingress` or `egress`.
  * `description` - The description of the rule.

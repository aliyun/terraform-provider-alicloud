---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_security_group_rule"
sidebar_current: "docs-alicloud-resource-security-group-rule"
description: |-
  Provides a Alicloud Security Group Rule resource.
---

# alicloud\_security\_group\_rule

Provides a security group rule resource.
Represents a single `ingress` or `egress` group rule, which can be added to external Security Groups.

-> **NOTE:**  `nic_type` should set to `intranet` when security group type is `vpc` or specifying the `source_security_group_id`. In this situation it does not distinguish between intranet and internet, the rule is effective on them both.


## Example Usage

Basic Usage

```
resource "alicloud_security_group" "default" {
  name = "default"
}

resource "alicloud_security_group_rule" "allow_all_tcp" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "internet"
  policy            = "accept"
  port_range        = "1/65535"
  priority          = 1
  security_group_id = alicloud_security_group.default.id
  cidr_ip           = "0.0.0.0/0"
}
```

## Module Support

You can use the existing [security-group module](https://registry.terraform.io/modules/alibaba/security-group/alicloud) 
to create a security group and add several rules one-click.

## Argument Reference

The following arguments are supported:

* `type` - (Required, ForceNew) The type of rule being created. Valid options are `ingress` (inbound) or `egress` (outbound).
* `ip_protocol` - (Required, ForceNew) The protocol. Can be `tcp`, `udp`, `icmp`, `gre` or `all`.
* `port_range` - (ForceNew) The range of port numbers relevant to the IP protocol. Default to "-1/-1". When the protocol is tcp or udp, each side port number range from 1 to 65535 and '-1/-1' will be invalid.
  For example, `1/200` means that the range of the port numbers is 1-200. Other protocols' 'port_range' can only be "-1/-1", and other values will be invalid.
* `security_group_id` - (Required, ForceNew) The security group to apply this rule to.
* `nic_type` - (Optional, ForceNew) Network type, can be either `internet` or `intranet`, the default value is `internet`.
* `policy` - (Optional, ForceNew) Authorization policy, can be either `accept` or `drop`, the default value is `accept`.
* `priority` - (Optional, ForceNew) Authorization policy priority, with parameter values: `1-100`, default value: 1.
* `cidr_ip` - (Optional, ForceNew) The target IP address range. The default value is 0.0.0.0/0 (which means no restriction will be applied). Other supported formats include 10.159.6.18/12. Only IPv4 is supported.
* `source_security_group_id` - (Optional, ForceNew) The target security group ID within the same region. If this field is specified, the `nic_type` can only select `intranet`.
* `source_group_owner_account` - (Optional, ForceNew) The Alibaba Cloud user account Id of the target security group when security groups are authorized across accounts.  This parameter is invalid if `cidr_ip` has already been set.
* `description` - (Optional) The description of the security group rule. The description can be up to 1 to 512 characters in length. Defaults to null.
* `prefix_list_id`- (Optional) The ID of the source/destination prefix list to which you want to control access. **NOTE:** If you specify `cidr_ip`,`source_security_group_id`,`ipv6_cidr_ip` parameter, this parameter is ignored.
* `ipv6_cidr_ip`- (Optional, ForceNew, Available in 1.174.0) Source IPv6 CIDR address block that requires access. Supports IP address ranges in CIDR format and IPv6 format. **NOTE:** This parameter cannot be set at the same time as the `cidr_ip` parameter.

-> **NOTE:**  You must specify one of the following field: `cidr_ip`,`source_security_group_id`,`prefix_list_id`,`ipv6_cidr_ip`. 

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the security group rule
* `type` - The type of rule, `ingress` or `egress`
* `port_range` - The range of port numbers
* `ip_protocol` - The protocol of the security group rule

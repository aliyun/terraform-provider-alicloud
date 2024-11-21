---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_security_group_rule"
sidebar_current: "docs-alicloud-resource-security-group-rule"
description: |-
  Provides a Alicloud Security Group Rule resource.
---

# alicloud_security_group_rule

Provides a Security Group Rule resource.

For information about Security Group Rule and how to use it, see [What is Rule](https://www.alibabacloud.com/help/en/ecs/user-guide/security-group-rules).

-> **NOTE:** Available since v0.1.0.

Represents a single `ingress` or `egress` group rule, which can be added to external Security Groups.

-> **NOTE:**  `nic_type` should set to `intranet` when security group type is `vpc` or specifying the `source_security_group_id`. In this situation it does not distinguish between intranet and internet, the rule is effective on them both.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_security_group_rule&exampleId=a52eccbb-c323-784b-a99d-acf4a0d03bcd98f6d5fc&activeTab=example&spm=docs.r.security_group_rule.0.a52eccbbc3&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
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

* `security_group_id` - (Required, ForceNew) The ID of the Security Group.
* `type` - (Required, ForceNew) The type of the Security Group Rule. Valid values:
  - `ingress`: inbound.
  - `egress`: outbound.
* `ip_protocol` - (Required, ForceNew) The transport layer protocol of the Security Group Rule. Valid values: `tcp`, `udp`, `icmp`, `gre`, `all`.
* `policy` - (Optional, ForceNew) The action of the Security Group Rule that determines whether to allow inbound access. Default value: `accept`. Valid values: `accept`, `drop`.
* `priority` - (Optional, ForceNew, Int) The priority of the Security Group Rule. Default value: `1`. Valid values: `1` to `100`.
* `cidr_ip` - (Optional, ForceNew) The target IP address range. The default value is 0.0.0.0/0 (which means no restriction will be applied). Other supported formats include 10.159.6.18/12. Only IPv4 is supported.
* `ipv6_cidr_ip`- (Optional, ForceNew, Available since v1.174.0) Source IPv6 CIDR address block that requires access. Supports IP address ranges in CIDR format and IPv6 format. **NOTE:** This parameter cannot be set at the same time as the `cidr_ip` parameter.
* `source_security_group_id` - (Optional, ForceNew) The target security group ID within the same region. If this field is specified, the `nic_type` can only select `intranet`.
* `source_group_owner_account` - (Optional, ForceNew) The Alibaba Cloud user account Id of the target security group when security groups are authorized across accounts.  This parameter is invalid if `cidr_ip` has already been set.
* `prefix_list_id`- (Optional, ForceNew) The ID of the source/destination prefix list to which you want to control access. **NOTE:** If you specify `cidr_ip`,`source_security_group_id`,`ipv6_cidr_ip` parameter, this parameter is ignored.
* `port_range` - (Optional, ForceNew) The range of port numbers relevant to the IP protocol. Default to "-1/-1". When the protocol is tcp or udp, each side port number range from 1 to 65535 and '-1/-1' will be invalid.
  For example, `1/200` means that the range of the port numbers is 1-200. Other protocols' 'port_range' can only be "-1/-1", and other values will be invalid.
* `nic_type` - (Optional, ForceNew) Network type, can be either `internet` or `intranet`, the default value is `internet`.
* `description` - (Optional) The description of the security group rule. The description can be up to 1 to 512 characters in length. Defaults to null.

-> **NOTE:**  You must specify one of the following field: `cidr_ip`,`source_security_group_id`,`prefix_list_id`,`ipv6_cidr_ip`. 

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Security Group Rule.

**NOTE:**  if `cidr_ip` is set, the `id` formats as `<security_group_id>:<type>:<ip_protocol>:<port_range>:<nic_type>:<cidr_ip>:<policy>:<priority>`.

**NOTE:**  if `ipv6_cidr_ip` is set, the `id` formats as `<security_group_id>:<type>:<ip_protocol>:<port_range>:<nic_type>:<ipv6_cidr_ip>:<policy>:<priority>`.

**NOTE:**  if `source_security_group_id` is set, the `id` formats as `<security_group_id>:<type>:<ip_protocol>:<port_range>:<nic_type>:<source_security_group_id>:<policy>:<priority>`.

**NOTE:**  if `prefix_list_id` is set, the `id` formats as `<security_group_id>:<type>:<ip_protocol>:<port_range>:<nic_type>:<prefix_list_id>:<policy>:<priority>`.

## Import

-> **NOTE:** Available since v1.224.0.

Security Group Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_security_group_rule.example <id>
```

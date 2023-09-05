---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_security_group"
description: |-
  Provides a Alicloud ECS Security Group resource.
---

# alicloud_security_group

Provides a ECS Security Group resource. 

For information about ECS Security Group and how to use it, see [What is Security Group](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.210.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_resource_manager_resource_group" "ResourceGroup" {
  display_name        = "test"
  resource_group_name = var.name
}

resource "alicloud_vpc" "VPC" {
  resource_group_id = alicloud_resource_manager_resource_group.ResourceGroup.id
  cidr_block        = "172.16.0.0/12"
  vpc_name          = "${var.name}1"
}


resource "alicloud_security_group" "default" {
  vpc_id              = alicloud_vpc.VPC.id
  security_group_name = var.name
  resource_group_id   = alicloud_resource_manager_resource_group.ResourceGroup.id
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) The description.
* `inner_access_policy` - (Optional) Network connectivity policy within the security group. Possible values:
  - Accept: intranet interworking
  - Drop: intranet isolation.
* `permissions` - (Optional, ForceNew) List of security group rules. The value range of N is 1~100. See [`permissions`](#permissions) below.
* `resource_group_id` - (Optional, ForceNew, Computed) The enterprise resource group ID where the security group resides.
* `security_group_name` - (Optional) The security group name.
* `security_group_type` - (Optional, ForceNew) Security group type.
* `service_managed` - (Optional, ForceNew) Whether the owner of the security group is a cloud product or vendor.
* `tags` - (Optional, Map) The tags.
* `vpc_id` - (Optional, ForceNew) Secure the group's proprietary network.

### `permissions`

The permissions supports the following:
* `description` - (Optional, ForceNew) The description of the security group rule. The length is 1~512 characters.
* `dest_cidr_ip` - (Optional, ForceNew) The IPv4 CIDR address segment of the destination. Supports IP address ranges in CIDR and IPv4 formats.For more information about how to support five-element rules, see [Security group five-element rules](~~ 97439 ~~).
* `ip_protocol` - (Optional, ForceNew) Transport layer protocol. The value is not case sensitive. Value range:
  - TCP.
  - UDP.
  - ICMP.
  - ICMPv6.
  - GRE.
  - ALL: ALL agreements are supported.

The value range of N is 1~100.

.
* `ipv6_dest_cidr_ip` - (Optional, ForceNew) The IPv6 CIDR address segment of the destination. Supports IP address ranges in CIDR and IPv6 formats.

For more information about how to support five-element rules, see [Security group five-element rules](~~ 97439 ~~).
-> **NOTE:**  valid only on VPC ECS instances that support IPv6, and this parameter cannot be set at the same time as the 'DestCidrIp' parameter.
* `ipv6_source_cidr_ip` - (Optional, ForceNew) The IPv6 CIDR address segment of the source terminal that needs to set access permissions. Supports IP address ranges in CIDR and IPv6 formats.
-> **NOTE:**  valid only on VPC ECS instances that support IPv6, and this parameter cannot be set at the same time as the 'SourceCidrIp' parameter.
* `nic_type` - (Optional, ForceNew, Computed) The NIC type of the classic network type security group rule. Value range:
  - internet: public network card.
  - intranet: intranet network card.

VPC type security group rules do not need to set the network interface card type. The default value is intranet. Only intranet.

When you set mutual access between security groups, you must specify only the DestGroupId parameter.

Default value: internet.

The value range of N is 1~100.
* `policy` - (Optional, ForceNew, Computed) Set access rights. Value range:
  - accept: accept access.
  - drop: Deny access and does not return a denial message, which indicates that the initiator request timed out or the connection cannot be established.

Default value: accept.

The value range of N is 1~100.
* `port_range` - (Optional, ForceNew) The destination port range related to the transport layer protocol open by the security group. Value range:
  - TCP/UDP: the value range is 1 to 65535. Use a forward slash (/) to separate the start and end ports. For example: 1/200.
  - ICMP:-1/-1.
  - GRE:-1/-1.
  - IpProtocol takes the value of ALL:-1/-1.

For more information about port scenarios, see [Common ports for typical applications](~~ 40724 ~~).
* `priority` - (Optional, ForceNew, Computed) Security Group rule priority. The smaller the number, the higher the priority. Value range: 1~100.Default value: 1.The value range of N is 1~100.
* `source_cidr_ip` - (Optional, ForceNew) The IPv4 CIDR address segment of the source terminal that needs to set access permissions. Supports IP address ranges in CIDR and IPv4 formats.
* `source_group_id` - (Optional, ForceNew) The ID of the source security group that needs to be granted access permissions.
  - Set at least one of the 'SourceGroupId', 'SourceCidrIp', 'Ipv6SourceCidrIp', or 'sourcecrefixlistid' parameters.
  - If the parameter 'SourceGroupId' is specified and the parameter 'SourceCidrIp' or 'Ipv6SourceCidrIp' is not specified, the value of the parameter 'NicType' can only be 'intranet '.
  - If both 'SourceGroupId' and 'SourceCidrIp' are specified, the default is 'SourceCidrIp.

The value range of N is 1~100.

You need to pay attention:
  - The Enterprise Security Group does not support authorized security group access.
  - The maximum number of authorized security groups supported by a normal security group is 20.
* `source_group_owner_account` - (Optional, ForceNew) When setting security group rules across accounts, the Alibaba Cloud account to which the source Security Group belongs.
  - If neither 'sourcegroupereraccount' nor 'sourcegroupernerid' is set, it is considered to set the access rights of your other security groups.
  - If the parameter' sourcecidrip' has been set, the parameter' sourcegroupereraccount' is invalid.
* `source_port_range` - (Optional, ForceNew) The source port range related to the transport layer protocol open by the security group. Value range:
  - TCP/UDP protocol: the value range is 1 to 65535. Use a forward slash (/) to separate the start and end ports. For example: 1/200.
  - ICMP protocol:-1/-1.
  - GRE agreement:-1/-1.
  - IpProtocol takes the value of ALL:-1/-1.

For more information about how to support five-element rules, see [Security group five-element rules](~~ 97439 ~~).
* `source_prefix_list_id` - (Optional, ForceNew) The ID of the source prefix list for which you need to set access permissions. You can call [DescribePrefixLists](~~ 205046 ~~) to query the Prefix List ID that can be used.

Precautions:
  - When the network type of the security group is classic network, Prefix List is not supported. For more information about security groups and Prefix List usage restrictions, see [Security Group Usage Restrictions](~~ 25412#SecurityGroupQuota1 ~~).
  - When you specify one of the'sourcecidrip', 'Ipv6SourceCidrIp', or'sourcegrouupid' parameters, this parameter is ignored.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The create time.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Security Group.
* `delete` - (Defaults to 5 mins) Used when delete the Security Group.
* `update` - (Defaults to 5 mins) Used when update the Security Group.

## Import

ECS Security Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_security_group.example <id>
```
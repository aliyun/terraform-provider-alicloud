---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_nat_firewall_control_policy"
description: |-
  Provides a Alicloud Cloud Firewall Nat Firewall Control Policy resource.
---

# alicloud_cloud_firewall_nat_firewall_control_policy

Provides a Cloud Firewall Nat Firewall Control Policy resource. Nat firewall access control policy.

For information about Cloud Firewall Nat Firewall Control Policy and how to use it, see [What is Nat Firewall Control Policy](https://www.alibabacloud.com/help/en/cloud-firewall/developer-reference/api-cloudfw-2017-12-07-createnatfirewallcontrolpolicy).

-> **NOTE:** Available since v1.224.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_firewall_nat_firewall_control_policy&exampleId=88066f9b-5886-fca1-809b-d729c073c467ad0b3162&activeTab=example&spm=docs.r.cloud_firewall_nat_firewall_control_policy.0.88066f9b58&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "eu-central-1"
}

variable "direction" {
  default = "out"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultDEiWfM" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "defaultFHDM3F" {
  vpc_id     = alicloud_vpc.defaultDEiWfM.id
  zone_id    = data.alicloud_zones.default.zones.0.id
  cidr_block = "172.16.2.0/24"
}

resource "alicloud_nat_gateway" "defaultMbS2Ts" {
  vpc_id           = alicloud_vpc.defaultDEiWfM.id
  nat_gateway_name = var.name
  payment_type     = "PayAsYouGo"
  vswitch_id       = alicloud_vswitch.defaultFHDM3F.id
  nat_type         = "Enhanced"
}

resource "alicloud_cloud_firewall_address_book" "port" {
  description  = format("%s%s", var.name, "port")
  group_name   = format("%s%s", var.name, "port")
  group_type   = "port"
  address_list = ["22/22", "23/23", "24/24"]
}

resource "alicloud_cloud_firewall_address_book" "port-update" {
  description  = format("%s%s", var.name, "port-update")
  group_name   = format("%s%s", var.name, "port-update")
  group_type   = "port"
  address_list = ["22/22", "23/23", "24/24"]
}

resource "alicloud_cloud_firewall_address_book" "domain" {
  description  = format("%s%s", var.name, "domain")
  group_name   = format("%s%s", var.name, "domain")
  group_type   = "domain"
  address_list = ["alibaba.com", "aliyun.com", "alicloud.com"]
}

resource "alicloud_cloud_firewall_address_book" "ip" {
  description  = var.name
  group_name   = var.name
  group_type   = "ip"
  address_list = ["1.1.1.1/32", "2.2.2.2/32"]
}

resource "alicloud_cloud_firewall_nat_firewall_control_policy" "default" {
  application_name_list = [
    "ANY"
  ]
  description = var.name
  release     = "false"
  ip_version  = "4"
  repeat_days = [
    "1",
    "2",
    "3"
  ]
  repeat_start_time   = "21:00"
  acl_action          = "log"
  dest_port_group     = alicloud_cloud_firewall_address_book.port.group_name
  repeat_type         = "Weekly"
  nat_gateway_id      = alicloud_nat_gateway.defaultMbS2Ts.id
  source              = "1.1.1.1/32"
  direction           = "out"
  repeat_end_time     = "21:30"
  start_time          = "1699156800"
  destination         = "1.1.1.1/32"
  end_time            = "1888545600"
  source_type         = "net"
  proto               = "TCP"
  new_order           = "1"
  destination_type    = "net"
  dest_port_type      = "group"
  domain_resolve_type = "0"
}
```

## Argument Reference

The following arguments are supported:
* `acl_action` - (Required) The method (action) of access traffic passing through Cloud Firewall in the security access control policy. Valid values:
  - **accept**: Release
  - **drop**: Refused
  - **log**: Observation.
* `application_name_list` - (Required) The list of application types supported by the access control policy.
* `description` - (Required) The description of the access control policy.
* `dest_port` - (Optional) The destination port of traffic access in the access control policy. Value:
  - When the protocol type is set to ICMP, the value of DestPort is null.
-> **NOTE:**  When the protocol type is ICMP, access control on the destination port is not supported.
  - When the protocol type is TCP, UDP, or ANY, and the destination port type (DestPortType) IS group, the value of DestPort is null.
-> **NOTE:**  When you select group (destination port address book) for the destination port type of the access control policy, you do not need to set a specific destination port number. All ports that need to be controlled by this access control policy are included in the destination port address book.
  - When the protocol type is TCP, UDP, or ANY, and the destination port type (DestPortType) is port, the value of DestPort is the destination port number.
* `dest_port_group` - (Optional) The address book name of the destination port of the access traffic in the access control policy.
-> **NOTE:**  When DestPortType is set to group, you need to set the destination port address book name.
* `dest_port_type` - (Optional) The destination port type of the access traffic in the security access control policy.
  - **port**: port
  - **group**: Port Address Book.
* `destination` - (Required) The destination address segment in the access control policy. Valid values:
  - When DestinationType is net, Destination is the Destination CIDR. For example: 1.2.XX.XX/24
  - When DestinationType IS group, Destination is the name of the Destination address book. For example: db_group
  - When DestinationType is domain, Destination is the Destination domain name. For example: * .aliyuncs.com
  - When DestinationType is location, Destination is the Destination region. For example: \["BJ11", "ZB"\].
* `destination_type` - (Required) The destination address type in the access control policy. Valid values:
  - **net**: Destination Network segment (CIDR address)
  - **group**: Destination Address Book
  - **domain**: the destination domain name.
* `direction` - (Required, ForceNew) The traffic direction of the access control policy. Valid values:
  - **out**: Internal and external traffic access control.
* `domain_resolve_type` - (Optional) The domain name resolution method of the access control policy. The policy is enabled by default after it is created. Valid values:
  - **0**: Based on FQDN
  - **1**: DNS-based dynamic resolution
  - **2**: dynamic resolution based on FQDN and DNS.
* `end_time` - (Optional) The end time of the policy validity period of the access control policy. Expresses using the second-level timestamp format. Must be full or half time and at least half an hour greater than the start time.
-> **NOTE:**  When RepeatType is set to permit, EndTime is null. When the RepeatType is None, Daily, Weekly, or Monthly, EndTime must have a value and you need to set the end time.
* `ip_version` - (Optional) Supported IP address version. Value:
  - **4** (default): indicates the IPv4 address.
* `nat_gateway_id` - (Required, ForceNew) The ID of the NAT gateway instance.
* `new_order` - (Required) The priority for the access control policy to take effect. The priority number increases sequentially from 1, and the smaller the priority number, the higher the priority.
* `proto` - (Required) The security protocol type for traffic access in the access control policy. Valid values:
  - ANY (indicates that all protocol types are queried)
  - TCP
  - UDP
  - ICMP.
* `release` - (Optional) The enabled status of the access control policy. The policy is enabled by default after it is created. Value:
  - **true**: Enable access control policy
  - **false**: Do not enable access control policies.
* `repeat_days` - (Optional) Collection of recurring dates for the policy validity period of the access control policy.
  - When RepeatType is 'Permanent', 'None', 'Daily', RepeatDays is an empty collection. For example:[]
  - When RepeatType is Weekly, RepeatDays cannot be empty. For example:["0", "6"]. When the RepeatType is set to Weekly, RepeatDays cannot be repeated.
  - RepeatDays cannot be empty when RepeatType is 'Monthly. For example:[1, 31]. When RepeatType is set to Monthly, RepeatDays cannot be repeated.
* `repeat_end_time` - (Optional) The recurring end time of the policy validity period of the access control policy. For example: 23:30, it must be the whole point or half point time, and at least half an hour greater than the repeat start time.
-> **NOTE:**  When RepeatType is set to normal or None, RepeatEndTime is null. When the RepeatType is Daily, Weekly, or Monthly, the RepeatEndTime must have a value, and you need to set the repeat end time.
* `repeat_start_time` - (Optional) The recurring start time of the policy validity period of the access control policy. For example: 08:00, it must be the whole point or half point time, and at least half an hour less than the repeat end time.
-> **NOTE:**  When RepeatType is set to permit or None, RepeatStartTime is empty. When the RepeatType is Daily, Weekly, or Monthly, the RepeatStartTime must have a value and you need to set the repeat start time.
* `repeat_type` - (Optional) The type of repetition for the policy validity period of the access control policy. Value:
  - **Permit** (default): Always
  - **None**: Specify a single time
  - **Daily**: Daily
  - **Weekly**: Weekly
  - **Monthly**: Monthly.
* `source` - (Required) The source address in the access control policy. Valid values:
  - When **SourceType** is set to 'net', Source is the Source CIDR address. For example: 10.2.4.0/24
  - When **SourceType** is set to 'group', Source is the name of the Source address book. For example: db_group.
* `source_type` - (Required) The source address type in the access control policy. Valid values:
  - **net**: the source network segment (CIDR address)
  - **group**: source address book
* `start_time` - (Optional) The start time of the policy validity period of the access control policy. Expresses using the second-level timestamp format. It must be a full or half hour and at least half an hour less than the end time.
-> **NOTE:**  When RepeatType is set to normal, StartTime is null. When the RepeatType is None, Daily, Weekly, or Monthly, StartTime must have a value and you need to set the start time.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<acl_uuid>:<nat_gateway_id>:<direction>`.
* `acl_uuid` - The unique ID of the security access control policy.
-> **NOTE:**  To modify a security access control policy, you need to provide the unique ID of the policy. You can call the DescribeNatFirewallControlPolicy interface to obtain the ID.
* `create_time` - The time when the policy was created.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Nat Firewall Control Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Nat Firewall Control Policy.
* `update` - (Defaults to 5 mins) Used when update the Nat Firewall Control Policy.

## Import

Cloud Firewall Nat Firewall Control Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_firewall_nat_firewall_control_policy.example <acl_uuid>:<nat_gateway_id>:<direction>
```
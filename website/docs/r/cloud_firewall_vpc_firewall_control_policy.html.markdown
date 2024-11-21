---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_vpc_firewall_control_policy"
sidebar_current: "docs-alicloud-resource-cloud-firewall-vpc-firewall-control-policy"
description: |-
  Provides a Alicloud Cloud Firewall Vpc Firewall Control Policy resource.
---

# alicloud_cloud_firewall_vpc_firewall_control_policy

Provides a Cloud Firewall Vpc Firewall Control Policy resource.

For information about Cloud Firewall Vpc Firewall Control Policy and how to use it, see [What is Vpc Firewall Control Policy](https://www.alibabacloud.com/help/en/cloud-firewall/latest/createvpcfirewallcontrolpolicy).

-> **NOTE:** Available since v1.194.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_firewall_vpc_firewall_control_policy&exampleId=6c065235-0ee8-66b0-1ddf-c9ee5c40e5e0f9ad03f2&activeTab=example&spm=docs.r.cloud_firewall_vpc_firewall_control_policy.0.6c0652350e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_account" "default" {
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  description       = "example_value"
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}

resource "alicloud_cloud_firewall_vpc_firewall_control_policy" "default" {
  order            = "1"
  destination      = "127.0.0.2/32"
  application_name = "ANY"
  description      = "example_value"
  source_type      = "net"
  dest_port        = "80/88"
  acl_action       = "accept"
  lang             = "zh"
  destination_type = "net"
  source           = "127.0.0.1/32"
  dest_port_type   = "port"
  proto            = "TCP"
  release          = true
  member_uid       = data.alicloud_account.default.id
  vpc_firewall_id  = alicloud_cen_instance.default.id
}
```

## Argument Reference

The following arguments are supported:
* `vpc_firewall_id` - (Required, ForceNew) The ID of the VPC firewall instance. Valid values:
  - When the VPC firewall protects traffic between two VPCs connected through the cloud enterprise network, the policy group ID uses the cloud enterprise network instance ID.
  - When the VPC firewall protects traffic between two VPCs connected through the express connection, the policy group ID uses the ID of the VPC firewall instance.
* `application_name` - (Required) The type of the applications that the access control policy supports. Valid values: `FTP`, `HTTP`, `HTTPS`, `MySQL`, `SMTP`, `SMTPS`, `RDP`, `VNC`, `SSH`, `Redis`, `MQTT`, `MongoDB`, `Memcache`, `SSL`, `ANY`.
* `description` - (Required) Access control over VPC firewalls description of the strategy information.
* `acl_action` - (Required) The action that Cloud Firewall performs on the traffic. Valid values: `accept`, `drop`, `log`.
* `source` - (Required) Access control over VPC firewalls strategy in the source address.
* `source_type` - (Required) The type of the source address in the access control policy. Valid values: `net`, `group`.
* `destination` - (Required) The destination address in the access control policy. Valid values:
  - If `destination_type` is set to `net`, the value of `destination` must be a CIDR block.
  - If `destination_type` is set to `group`, the value of `destination` must be an address book.
  - If `destination_type` is set to `domain`, the value of `destination` must be a domain name.
* `destination_type` - (Required) The type of the destination address in the access control policy. Valid values: `net`, `group`, `domain`.
* `proto` - (Required) The type of the protocol in the access control policy. Valid values: `ANY`, `TCP`, `UDP`, `ICMP`.
* `order` - (Required, ForceNew, Int) The priority of the access control policy. The priority value starts from 1. A smaller priority value indicates a higher priority.
* `dest_port` - (Optional) The destination port in the access control policy. **Note:** If `dest_port_type` is set to `port`, you must specify this parameter.
* `dest_port_group` - (Optional) Access control policy in the access traffic of the destination port address book name. **Note:** If `dest_port_type` is set to `group`, you must specify this parameter.
* `dest_port_type` - (Optional) The type of the destination port in the access control policy. Valid values: `port`, `group`.
* `release` - (Optional, Bool) The enabled status of the access control policy. The policy is enabled by default after it is created.. Valid values:
  - `true`: Enable access control policies
  - `false`: does not enable access control policies.
* `member_uid` - (Optional, ForceNew) The UID of the member account of the current Alibaba cloud account.
* `lang` - (Optional) The language of the content within the request and response. Valid values: `zh`, `en`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Vpc Firewall Control Policy. The value formats as `<vpc_firewall_id>:<acl_uuid>`.
* `acl_uuid` - Access control over VPC firewalls strategy unique identifier.
* `application_id` - Policy specifies the application ID.
* `source_group_cidrs` - SOURCE address of the address list.
* `source_group_type` - The source address type in the access control policy.
* `destination_group_cidrs` - Destination address book defined in the address list.
* `destination_group_type` - The destination address book type in the access control policy.
* `dest_port_group_ports` - Port Address Book port list.
* `hit_times` - Control strategy of hits per second.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Vpc Firewall Control Policy.
* `update` - (Defaults to 5 mins) Used when update the Vpc Firewall Control Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Vpc Firewall Control Policy.

## Import

Cloud Firewall Vpc Firewall Control Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_firewall_vpc_firewall_control_policy.example <vpc_firewall_id>:<acl_uuid>
```

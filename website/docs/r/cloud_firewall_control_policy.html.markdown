---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_control_policy"
sidebar_current: "docs-alicloud-resource-cloud-firewall-control-policy"
description: |-
  Provides a Alicloud Cloud Firewall Control Policy resource.
---

# alicloud_cloud_firewall_control_policy

Provides a Cloud Firewall Control Policy resource.

For information about Cloud Firewall Control Policy and how to use it, see [What is Control Policy](https://www.alibabacloud.com/help/doc-detail/138867.htm).

-> **NOTE:** Available since v1.129.0.

## Example Usage
<div class="oics-button" style="float: right;margin: 0 0 -40px 0;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_cloud_firewall_control_policy&exampleId=9466770b-c784-5cd7-73b1-db97948c811470394f57&activeTab=example&spm=docs.r.cloud_firewall_control_policy.0.9466770bc7" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_cloud_firewall_control_policy" "default" {
  direction        = "in"
  application_name = "ANY"
  description      = var.name
  acl_action       = "accept"
  source           = "127.0.0.1/32"
  source_type      = "net"
  destination      = "127.0.0.2/32"
  destination_type = "net"
  proto            = "ANY"
}
```

## Argument Reference

The following arguments are supported:

* `direction` - (Required, ForceNew) The direction of the traffic to which the access control policy applies. Valid values: `in`, `out`.
* `application_name` - (Required) The application type supported by the access control policy. Valid values: `ANY`, `HTTP`, `HTTPS`, `MQTT`, `Memcache`, `MongoDB`, `MySQL`, `RDP`, `Redis`, `SMTP`, `SMTPS`, `SSH`, `SSL`, `VNC`.
-> **NOTE:** If `proto` is set to `TCP`, you can set `application_name` to any valid value. If `proto` is set to `UDP`, `ICMP`, or `ANY`, you can only set `application_name` to `ANY`.
* `description` - (Required) The description of the access control policy.
* `acl_action` - (Required) The action that Cloud Firewall performs on the traffic. Valid values: `accept`, `drop`, `log`.
* `source` - (Required) The source address in the access control policy.
* `source_type` - (Required) The type of the source address in the access control policy. Valid values: `net`, `group`, `location`.
* `destination` - (Required) The destination address in the access control policy.
* `destination_type` - (Required) The type of the destination address in the access control policy. Valid values: `net`, `group`, `domain`, `location`.
* `proto` - (Required) The protocol type supported by the access control policy. Valid values: `ANY`, ` TCP`, `UDP`, `ICMP`.
* `dest_port` - (Optional) The destination port in the access control policy. **Note:** If `dest_port_type` is set to `port`, you must specify `dest_port`.
* `dest_port_group` - (Optional) The name of the destination port address book in the access control policy. **Note:** If `dest_port_type` is set to `group`, you must specify `dest_port_group`.
* `dest_port_type` - (Optional) The type of the destination port in the access control policy. Valid values: `port`, `group`.
* `ip_version` - (Optional, ForceNew) The IP version supported by the access control policy. Default value: `4`. Valid values:
  - `4`: IPv4.
  - `6`: IPv6.
* `release` - (Optional) The status of the access control policy. Valid values: `true`, `false`.
* `source_ip` - (Optional) The source IP address of the request.
* `lang` - (Optional) The language of the content within the request and response. Valid values: `zh`, `en`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Control Policy. It formats as `<acl_uuid>:<direction>`.
* `acl_uuid` - (Available since v1.148.0) The unique ID of the access control policy.

## Import

Cloud Firewall Control Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_firewall_control_policy.example <acl_uuid>:<direction>
```

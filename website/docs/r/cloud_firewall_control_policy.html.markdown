---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_control_policy"
sidebar_current: "docs-alicloud-resource-cloud-firewall-control-policy"
description: |-
  Provides a Alicloud Cloud Firewall Control Policy resource.
---

# alicloud\_cloud\_firewall\_control\_policy

Provides a Cloud Firewall Control Policy resource.

For information about Cloud Firewall Control Policy and how to use it, see [What is Control Policy](https://www.alibabacloud.com/help/doc-detail/138867.htm).

-> **NOTE:** Available in v1.129.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cloud_firewall_control_policy" "example" {
  application_name = "ANY"
  acl_action       = "accept"
  description      = "example"
  destination_type = "net"
  destination      = "100.1.1.0/24"
  direction        = "out"
  proto            = "ANY"
  source           = "1.2.3.0/24"
  source_type      = "net"
}

```

## Argument Reference

The following arguments are supported:

* `acl_action` - (Required) The action that Cloud Firewall performs on the traffic. Valid values: `accept`, `drop`, `log`.
* `application_name` - (Required) The application type that the access control policy supports.If `direction` is `in`, the valid value is `ANY`. If `direction` is `out`, the valid values are `ANY`, `HTTP`, `HTTPS`, `MQTT`, `Memcache`, `MongoDB`, `MySQL`, `RDP`, `Redis`, `SMTP`, `SMTPS`, `SSH`, `SSL`, `VNC`.
* `description` - (Required) The description of the access control policy.
* `dest_port` - (Optional) The destination port defined in the access control policy. 
* `dest_port_group` - (Optional) The destination port address book defined in the access control policy.
* `dest_port_type` - (Optional) The destination port type defined in the access control policy. Valid values: `group`, `port`.
* `destination` - (Required) The destination address defined in the access control policy.
* `destination_type` - (Required) DestinationType. Valid values: If Direction is `in`, the valid values are `net`, `group`. If `direction` is `out`, the valid values are `net`, `group`, `domain`, `location`.
* `direction` - (Required, ForceNew) Direction. Valid values: `in`, `out`.
* `ip_version` - (Optional) The ip version.
* `lang` - (Optional) DestPortGroupPorts. Valid values: `en`, `zh`.
* `proto` - (Required) Proto. Valid values: ` TCP`, ` UDP`, `ANY`, `ICMP`.
* `release` - (Optional) Specifies whether the access control policy is enabled. By default, an access control policy is enabled after it is created. Valid values: `true`, `false`.
* `source` - (Required) Source.
* `source_ip` - (Optional) The source ip.
* `source_type` - (Required) SourceType. Valid values: If `direction` is `in`, the valid values are `net`, `group`, `location`. If `direction` is `out`, the valid values are `net`, `group`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Control Policy. The value formats as `<acl_uuid>:<direction>`.
* `acl_uuid` - (Available in v1.148.0+) The unique ID of the access control policy.

## Import

Cloud Firewall Control Policy can be imported using the id, e.g.

```
$ terraform import alicloud_cloud_firewall_control_policy.example <acl_uuid>:<direction>
```

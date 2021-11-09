---
subcategory: "Smart Access Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_sag_dnat_entry"
sidebar_current: "docs-alicloud-resource-sag-dnat-entry"
description: |-
  Provides a Sag DnatEntry resource.
---

# alicloud\_sag\_dnat_entry

Provides a Sag DnatEntry resource. This topic describes how to add a DNAT entry to a Smart Access Gateway (SAG) instance to enable the DNAT function. By using the DNAT function, you can forward requests received by public IP addresses to Alibaba Cloud instances according to custom mapping rules.

For information about Sag DnatEntry and how to use it, see [What is Sag DnatEntry](https://www.alibabacloud.com/help/doc-detail/124312.htm).

-> **NOTE:** Available in 1.63.0+

-> **NOTE:** Only the following regions suppor. [`cn-shanghai`, `cn-shanghai-finance-1`, `cn-hongkong`, `ap-southeast-1`, `ap-southeast-2`, `ap-southeast-3`, `ap-southeast-5`, `ap-northeast-1`, `eu-central-1`]

## Example Usage

Basic Usage

```
resource "alicloud_sag_dnat_entry" "default" {
  sag_id        = "sag-3rb1t3iagy3w0zgwy9"
  type          = "Intranet"
  ip_protocol   = "tcp"
  external_ip   = "1.0.0.2"
  external_port = "1"
  internal_ip   = "10.0.0.2"
  internal_port = "20"
}
```
## Argument Reference

The following arguments are supported:

* `sag_id` - (Required) The ID of the SAG instance.
* `type` - (Required) The DNAT type. Valid values: Intranet: DNAT of private IP addresses. Internet: DNAT of public IP addresses
* `ip_protocol` - (Required) The protocol type. Valid values: TCP: Forwards packets of the TCP protocol. UDP: Forwards packets of the UDP protocol. Any: Forwards packets of all protocols.
* `external_ip` - (Optional) The external public IP address.when "type" is "Internet",automatically identify the external ip.
* `external_port` - (Required) The public port.Value range: 1 to 65535 or "any".
* `internal_ip` - (Required) The destination private IP address.
* `internal_port` - (Required) The destination private port.Value range: 1 to 65535 or "any".


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the DNAT entry Id and formates as `<sag_id>:<dnat_id>`.

## Import

The Sag DnatEntry can be imported using the id, e.g.

```
$ terraform import alicloud_sag_dnat_entry.example sag-abc123456:dnat-abc123456
```

---
subcategory: "Smart Access Gateway (Smartag)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sag_dnat_entry"
sidebar_current: "docs-alicloud-resource-sag-dnat-entry"
description: |-
  Provides a Sag DnatEntry resource.
---

# alicloud_sag_dnat_entry

Provides a Sag DnatEntry resource. This topic describes how to add a DNAT entry to a Smart Access Gateway (SAG) instance to enable the DNAT function. By using the DNAT function, you can forward requests received by public IP addresses to Alibaba Cloud instances according to custom mapping rules.

For information about Sag DnatEntry and how to use it, see [What is Sag DnatEntry](https://www.alibabacloud.com/help/en/smart-access-gateway/latest/adddnatentry).

-> **NOTE:** Available since v1.63.0.

-> **NOTE:** Only the following regions suppor. [`cn-shanghai`, `cn-shanghai-finance-1`, `cn-hongkong`, `ap-southeast-1`, `ap-southeast-2`, `ap-southeast-3`, `ap-southeast-5`, `ap-northeast-1`, `eu-central-1`]

## Example Usage

Basic Usage

```terraform
variable "sag_id" {
  default = "sag-9bifkfaz***"
}
provider "alicloud" {
  region = "cn-shanghai"
}

resource "alicloud_sag_dnat_entry" "default" {
  sag_id        = var.sag_id
  type          = "Intranet"
  ip_protocol   = "any"
  external_ip   = "172.32.0.2"
  external_port = "any"
  internal_ip   = "172.16.0.4"
  internal_port = "any"
}
```
## Argument Reference

The following arguments are supported:

* `sag_id` - (Required, ForceNew) The ID of the SAG instance.
* `type` - (Required, ForceNew) The DNAT type. Valid values: Intranet: DNAT of private IP addresses. Internet: DNAT of public IP addresses
* `ip_protocol` - (Required, ForceNew) The protocol type. Valid values: TCP: Forwards packets of the TCP protocol. UDP: Forwards packets of the UDP protocol. Any: Forwards packets of all protocols.
* `external_ip` - (Optional, ForceNew) The external public IP address.when "type" is "Internet",automatically identify the external ip.
* `external_port` - (Required, ForceNew) The public port.Value range: 1 to 65535 or "any".
* `internal_ip` - (Required, ForceNew) The destination private IP address.
* `internal_port` - (Required, ForceNew) The destination private port.Value range: 1 to 65535 or "any".


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the DNAT entry Id and formates as `<sag_id>:<dnat_id>`.

## Import

The Sag DnatEntry can be imported using the id, e.g.

```shell
$ terraform import alicloud_sag_dnat_entry.example sag-abc123456:dnat-abc123456
```

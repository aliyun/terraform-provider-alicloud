---
subcategory: "Smart Access Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_sag_snat_entry"
sidebar_current: "docs-alicloud-resource-sag-snat-entry"
description: |-
  Provides a Sag SnatEntry resource.
---

# alicloud\_sag\_snat_entry

Provides a Sag SnatEntry resource. This topic describes how to add a SNAT entry to enable the SNAT function. The SNAT function can hide internal IP addresses and resolve private IP address conflicts. With this function, on-premises sites can access internal IP addresses, but cannot be accessed by internal IP addresses. If you do not add a SNAT entry, on-premises sites can access each other only when all related IP addresses do not conflict.

For information about Sag SnatEntry and how to use it, see [What is Sag SnatEntry](https://www.alibabacloud.com/help/doc-detail/124231.htm).

-> **NOTE:** Available in 1.61.0+

-> **NOTE:** Only the following regions support. [`cn-shanghai`, `cn-shanghai-finance-1`, `cn-hongkong`, `ap-southeast-1`, `ap-southeast-2`, `ap-southeast-3`, `ap-southeast-5`, `ap-northeast-1`, `eu-central-1`]

## Example Usage

Basic Usage

```
resource "alicloud_sag_snat_entry" "default" {
  sag_id     = "sag-3rb1t3iagy3w0zgwy9"
  cidr_block = "192.168.7.0/24"
  snat_ip    = "192.0.0.2"
}
```
## Argument Reference

The following arguments are supported:

* `sag_id` - (Required) The ID of the SAG instance.
* `cidr_block` - (Required) The destination CIDR block.
* `snat_ip` - (Required) The public IP address.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the SNAT entry Id and formates as `<sag_id>:<snat_id>`.

## Import

The Sag SnatEntry can be imported using the id, e.g.

```
$ terraform import alicloud_sag_snat_entry.example sag-abc123456:snat-abc123456
```


---
subcategory: "Smart Access Gateway (Smartag)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sag_snat_entry"
sidebar_current: "docs-alicloud-resource-sag-snat-entry"
description: |-
  Provides a Sag SnatEntry resource.
---

# alicloud_sag_snat_entry

Provides a Sag SnatEntry resource. This topic describes how to add a SNAT entry to enable the SNAT function. The SNAT function can hide internal IP addresses and resolve private IP address conflicts. With this function, on-premises sites can access internal IP addresses, but cannot be accessed by internal IP addresses. If you do not add a SNAT entry, on-premises sites can access each other only when all related IP addresses do not conflict.

For information about Sag SnatEntry and how to use it, see [What is Sag SnatEntry](https://www.alibabacloud.com/help/en/smart-access-gateway/latest/addsnatentry).

-> **NOTE:** Available since v1.61.0.

-> **NOTE:** Only the following regions support. [`cn-shanghai`, `cn-shanghai-finance-1`, `cn-hongkong`, `ap-southeast-1`, `ap-southeast-3`, `ap-southeast-5`, `ap-northeast-1`, `eu-central-1`]

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_sag_snat_entry&exampleId=4452b9ba-56f4-3487-f5e3-d51ba657dfc5889edb12&activeTab=example&spm=docs.r.sag_snat_entry.0.4452b9ba56&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "sag_id" {
  default = "sag-9bifk***"
}
provider "alicloud" {
  region = "cn-shanghai"
}

resource "alicloud_sag_snat_entry" "default" {
  sag_id     = var.sag_id
  cidr_block = "192.168.7.0/24"
  snat_ip    = "192.0.0.2"
}
```
## Argument Reference

The following arguments are supported:

* `sag_id` - (Required, ForceNew) The ID of the SAG instance.
* `cidr_block` - (Required, ForceNew) The destination CIDR block.
* `snat_ip` - (Required, ForceNew) The public IP address.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the SNAT entry Id and formates as `<sag_id>:<snat_id>`.

## Import

The Sag SnatEntry can be imported using the id, e.g.

```shell
$ terraform import alicloud_sag_snat_entry.example sag-abc123456:snat-abc123456
```


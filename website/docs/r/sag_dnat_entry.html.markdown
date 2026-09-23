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

-> **NOTE:** Only the following regions suppor. [`cn-shanghai`, `cn-shanghai-finance-1`, `cn-hongkong`, `ap-southeast-1`, `ap-southeast-3`, `ap-southeast-5`, `ap-northeast-1`, `eu-central-1`]

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_sag_dnat_entry&exampleId=0a75c941-e9e0-38bf-abfb-aff3c517ade8a8afe990&activeTab=example&spm=docs.r.sag_dnat_entry.0.0a75c941e9&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_sag_dnat_entry&spm=docs.r.sag_dnat_entry.example&intl_lang=EN_US)

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

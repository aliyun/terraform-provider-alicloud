---
subcategory: "Smart Access Gateway (Smartag)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sag_qos"
sidebar_current: "docs-alicloud-resource-sag-qos"
description: |-
  Provides a Sag Qos resource.
---

# alicloud_sag_qos

Provides a Sag Qos resource. Smart Access Gateway (SAG) supports quintuple-based QoS functions to differentiate traffic of different services and ensure high-priority traffic bandwidth.

For information about Sag Qos and how to use it, see [What is Qos](https://www.alibabacloud.com/help/en/smart-access-gateway/latest/createqos).

-> **NOTE:** Available since v1.60.0.

-> **NOTE:** Only the following regions support. [`cn-shanghai`, `cn-shanghai-finance-1`, `cn-hongkong`, `ap-southeast-1`, `ap-southeast-3`, `ap-southeast-5`, `ap-northeast-1`, `eu-central-1`]

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_sag_qos&exampleId=21072698-3757-8ecd-3386-113c53ca1a8151e11df8&activeTab=example&spm=docs.r.sag_qos.0.2107269837&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-shanghai"
}
resource "alicloud_sag_qos" "default" {
  name = "terraform-example"
}

```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_sag_qos&spm=docs.r.sag_qos.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the QoS policy to be created. The name can contain 2 to 128 characters including a-z, A-Z, 0-9, periods, underlines, and hyphens. The name must start with an English letter, but cannot start with http:// or https://.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Qos. For example "qos-xxx".

## Import

The Sag Qos can be imported using the id, e.g.

```shell
$ terraform import alicloud_sag_qos.example qos-abc123456
```


---
subcategory: "Smart Access Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_sag_qos"
sidebar_current: "docs-alicloud-resource-sag-qos"
description: |-
  Provides a Sag Qos resource.
---

# alicloud\_sag\_qos

Provides a Sag Qos resource. Smart Access Gateway (SAG) supports quintuple-based QoS functions to differentiate traffic of different services and ensure high-priority traffic bandwidth.

For information about Sag Qos and how to use it, see [What is Qos](https://www.alibabacloud.com/help/doc-detail/131306.htm).

-> **NOTE:** Available in 1.60.0+

-> **NOTE:** Only the following regions support. [`cn-shanghai`, `cn-shanghai-finance-1`, `cn-hongkong`, `ap-southeast-1`, `ap-southeast-2`, `ap-southeast-3`, `ap-southeast-5`, `ap-northeast-1`, `eu-central-1`]

## Example Usage

Basic Usage

```
resource "alicloud_sag_qos" "default" {
  name        = "tf-testAccSagQosName"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the QoS policy to be created. The name can contain 2 to 128 characters including a-z, A-Z, 0-9, periods, underlines, and hyphens. The name must start with an English letter, but cannot start with http:// or https://.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Qos. For example "qos-xxx".

## Import

The Sag Qos can be imported using the id, e.g.

```
$ terraform import alicloud_sag_qos.example qos-abc123456
```


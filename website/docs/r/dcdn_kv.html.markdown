---
subcategory: "DCDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_dcdn_kv"
sidebar_current: "docs-alicloud-resource-dcdn-kv"
description: |-
  Provides a Alicloud Dcdn Kv resource.
---

# alicloud_dcdn_kv

Provides a Dcdn Kv resource.

For information about Dcdn Kv and how to use it, see [What is Kv](https://www.alibabacloud.com/help/en/dcdn/developer-reference/api-dcdn-2018-01-15-putdcdnkv).

-> **NOTE:** Available since v1.198.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_dcdn_kv&exampleId=8bd8218a-a8d0-b4ba-c5b7-56e9974141a15612ff32&activeTab=example&spm=docs.r.dcdn_kv.0.8bd8218aa8&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_dcdn_kv_namespace" "default" {
  description = var.name
  namespace   = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_dcdn_kv" "default" {
  value     = "example-value"
  key       = "${var.name}-${random_integer.default.result}"
  namespace = alicloud_dcdn_kv_namespace.default.namespace
}
```

## Argument Reference

The following arguments are supported:
* `key` - (Required, ForceNew) The name of the key to Put, the longest 512, cannot contain spaces.
* `namespace` - (Required, ForceNew) The name specified when the customer calls PutDcdnKvNamespace.
* `value` - (Required) The content of key, up to 2M(2*1000*1000).



## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.The value is formulated as `<namespace>:<key>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Kv.
* `delete` - (Defaults to 5 mins) Used when delete the Kv.
* `update` - (Defaults to 5 mins) Used when update the Kv.

## Import

Dcdn Kv can be imported using the id, e.g.

```shell
$ terraform import alicloud_dcdn_kv.example <namespace>:<key>
```
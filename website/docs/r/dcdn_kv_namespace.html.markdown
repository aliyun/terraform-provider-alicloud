---
subcategory: "DCDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_dcdn_kv_namespace"
sidebar_current: "docs-alicloud-resource-dcdn-kv-namespace"
description: |-
  Provides a Alicloud Dcdn Kv Namespace resource.
---

# alicloud_dcdn_kv_namespace

Provides a Dcdn Kv Namespace resource.

For information about Dcdn Kv Namespace and how to use it, see [What is Kv Namespace](https://www.alibabacloud.com/help/en/dcdn/developer-reference/api-dcdn-2018-01-15-putdcdnkvnamespace).

-> **NOTE:** Available since v1.198.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_dcdn_kv_namespace&exampleId=47c848a2-795c-9450-24f5-ffe0ec6b734e3430d0f0&activeTab=example&spm=docs.r.dcdn_kv_namespace.0.47c848a279&intl_lang=EN_US" target="_blank">
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
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_dcdn_kv_namespace&spm=docs.r.dcdn_kv_namespace.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `description` - (Required, ForceNew) Namespace description information
* `namespace` - (Required, ForceNew) Namespace name. The name can contain letters, digits, hyphens (-), and underscores (_).


## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Kv Namespace.
* `delete` - (Defaults to 5 mins) Used when delete the Kv Namespace.

## Import

Dcdn Kv Namespace can be imported using the id, e.g.

```shell
$ terraform import alicloud_dcdn_kv_namespace.example 
```
---
subcategory: "Corpus Store (CS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_corpusdoc"
description: |-
  Fixture resource without RpcPost source — 产品目录回落 subcategory。不对应真实资源。
---

# alicloud_corpusdoc

Fixture. corpus-gen 用:无 RpcPost 源码(SDK 风格),`_source_pv` 抽不到产品 →
产品目录回落 docs subcategory("Corpus Store (CS)" → 净化为 `corpus`)。

## Example Usage

```terraform
resource "alicloud_corpusdoc" "default" {
  corpusdoc_name = "tf-example"
}
```

## Argument Reference

The following arguments are supported:
* `corpusdoc_name` - (Required) 名称。

## Attributes Reference

* `id` - 资源 ID。

---
subcategory: "CorpusNet"
layout: "alicloud"
page_title: "Alicloud: alicloud_corpusfree"
description: |-
  Fixture free resource for corpus-gen hermetic tests. 不对应真实 provider 资源。
---

# alicloud_corpusfree

Fixture. corpus-gen 用:免费资源,含可命名 `name` 字段 + `tags` 块,供注入/产品推导断言。

## Example Usage

Basic Usage

```terraform
resource "alicloud_corpusfree" "default" {
  name        = "tf-example"
  description = "corpus fixture"
  tags = {
    Created = "TF"
  }
}
```

## Argument Reference

The following arguments are supported:
* `name` - (Required) 名称。
* `description` - (Optional) 描述。
* `tags` - (Optional) 标签。

## Attributes Reference

* `id` - 资源 ID。

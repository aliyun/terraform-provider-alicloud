---
subcategory: "CorpusPay"
layout: "alicloud"
page_title: "Alicloud: alicloud_corpuschargeds"
description: |-
  Fixture data source for alicloud_corpuscharged (ds- 只读变体来源). 不对应真实资源。
---

# alicloud_corpuschargeds

Fixture data source. corpus-gen 用:订阅类资源 `alicloud_corpuscharged` 的 ds- 只读变体从此抽 data 块。
生成器按 `<short>` / `<short>s` / `<short>es` 三试命中此文档(corpuscharged → corpuschargeds)。

## Example Usage

```terraform
data "alicloud_corpuschargeds" "default" {
  name_regex = "^tf"
  status     = "Running"
}

output "first_id" {
  value = data.alicloud_corpuschargeds.default.ids.0
}
```

## Argument Reference

The following arguments are supported:
* `name_regex` - (Optional) 名称正则。
* `status` - (Optional) 状态过滤。

## Attributes Reference

* `ids` - 匹配到的 ID 列表。

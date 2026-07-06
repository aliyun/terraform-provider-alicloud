---
subcategory: "CorpusSub"
layout: "alicloud"
page_title: "Alicloud: alicloud_corpusperiod"
description: |-
  Fixture resource with a standalone subscription period. 不对应真实资源。
---

# alicloud_corpusperiod

Fixture. corpus-gen 值敏感门用:独立 `period = 1`(订阅时长,整数 ≤ period_subscription_max)
→ 命中订阅门 → 生成器写 `apply: false`。无对应 data source 文档 → 只留规范注记、不生成 ds- 变体。

## Example Usage

```terraform
resource "alicloud_corpusperiod" "default" {
  sub_name = "tf-example"
  period   = 1
}
```

## Argument Reference

The following arguments are supported:
* `sub_name` - (Optional) 名称。
* `period` - (Optional) 订阅时长(月)。

## Attributes Reference

* `id` - 资源 ID。

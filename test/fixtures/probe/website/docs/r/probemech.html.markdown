---
subcategory: "ProbeMech"
layout: "alicloud"
page_title: "Alicloud: alicloud_probemech"
description: |-
  Fixture resource for tier-0 mechanical-diff hermetic tests.
---

# alicloud_probemech

Fixture resource. 文档故意与源码 flag 一致(零 doc↔source 噪声),只让机械 TF↔API diff 的陷阱浮出,供 test/probe_test.sh 断言。不对应真实 provider 资源。

## Example Usage

```terraform
resource "alicloud_probemech" "x" {
  name = "example"
}
```

## Argument Reference

The following arguments are supported:
* `name` - (Required) 名称。
* `storage_class` - (Optional) 存储类型。
* `required_field` - (Optional) 必填字段(源码 Optional,API required)。
* `mask` - (Optional) 掩码。
* `mode_value` - (Optional) 模式值。
* `conflict_type` - (Optional) 类型冲突字段(源码 list,API string)。
* `client_token` - (Optional) 幂等 token(应被抑制表命中)。
* `safe_enum` - (Optional) TF 更严枚举。
* `renamed_field` - (Optional) convert 改名字段,机械层映射不上。
* `opaque_enum` - (Optional) 枚举取自变量,机械层不可解析。

## Attributes Reference

The following attributes are exported:
* `id` - 资源 ID。

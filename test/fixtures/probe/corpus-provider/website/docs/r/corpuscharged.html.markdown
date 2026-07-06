---
subcategory: "CorpusPay"
layout: "alicloud"
page_title: "Alicloud: alicloud_corpuscharged"
description: |-
  Fixture resource whose Example HCL carries a subscription billing value. 不对应真实资源。
---

# alicloud_corpuscharged

Fixture. corpus-gen 值敏感门用:Example HCL 含 `instance_charge_type = "PrePaid"`(订阅语义值)
→ 命中成本安全门 → 生成器写 `apply: false`;并因存在对应 data source 文档而生成 `ds-corpuscharged` 只读变体。

## Example Usage

```terraform
resource "alicloud_corpuscharged" "default" {
  db_name              = "tf-example"
  instance_charge_type = "PrePaid"
}
```

## Argument Reference

The following arguments are supported:
* `db_name` - (Optional) 名称。
* `instance_charge_type` - (Optional) 计费方式。Valid values: `PrePaid`, `PostPaid`.

## Attributes Reference

* `id` - 资源 ID。

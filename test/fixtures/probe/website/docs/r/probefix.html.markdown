---
subcategory: "ProbeFix"
layout: "alicloud"
page_title: "Alicloud: alicloud_probefix"
description: |-
  Fixture resource for tier-0 parser hermetic tests.
---

# alicloud_probefix

Fixture resource. 刻意布置五类 doc↔source gap 供 test/probe_test.sh 断言,不对应真实 provider 资源。

## Example Usage

```terraform
resource "alicloud_probefix" "x" {
  name = "example"
}
```

## Argument Reference

The following arguments are supported:
* `name` - (Required) 名称。文档=源码 Required,应无 finding。
* `flag_field` - (Required) 文档标 Required,但源码是 Optional(flag_mismatch)。
* `forcenew_field` - (Optional) 文档未标 ForceNew,但源码是 ForceNew(forcenew,可能意外重建)。
* `deprecated_field` - (Optional) 源码已 Deprecated,文档仍作正常参数列出(deprecated)。
* `nested_block` - (Optional) 一个嵌套块;其内层 `inner_field` 不应被源码解析器当顶层键。
* `phantom_field` - (Optional) 只存在于文档,源码 schema 无(phantom)。

## Attributes Reference

The following attributes are exported:
* `id` - 资源 ID(纯 Computed,不应被当未文档化参数)。
* `create_time` - 创建时间。

---
subcategory: "CorpusCompute"
layout: "alicloud"
page_title: "Alicloud: alicloud_corpus_instance"
description: |-
  Fixture big-ticket 按量付费资源 — 值敏感门下应直接 apply(放行)。不对应真实资源。
---

# alicloud_corpus_instance

Fixture. corpus-gen 值敏感门用:资源名含 `instance`(付费大件),但计费值为按量(`PostPaid`)、
`period` 取秒级 metric 值(900)、另有 `backup_retention_period` —— 三者均**不应命中**订阅门,生成器写默认 `apply: true`。

## Example Usage

```terraform
resource "alicloud_corpus_instance" "default" {
  instance_name           = "tf-example"
  instance_type           = "small"
  instance_charge_type    = "PostPaid"
  period                  = 900
  backup_retention_period = 7
}
```

## Argument Reference

The following arguments are supported:
* `instance_name` - (Optional) 名称。
* `instance_type` - (Required) 规格。
* `instance_charge_type` - (Optional) 计费方式。Valid values: `PrePaid`, `PostPaid`.
* `period` - (Optional) 采集周期(秒)。
* `backup_retention_period` - (Optional) 备份保留天数。

## Attributes Reference

* `id` - 资源 ID。

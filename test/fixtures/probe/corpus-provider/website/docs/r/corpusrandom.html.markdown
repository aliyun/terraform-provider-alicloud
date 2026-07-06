---
subcategory: "CorpusMisc"
layout: "alicloud"
page_title: "Alicloud: alicloud_corpusrandom"
description: |-
  Fixture exercising provider/terraform block strip + extra provider + heredoc. 不对应真实资源。
---

# alicloud_corpusrandom

Fixture. corpus-gen 用:Example 含 `terraform{}`/`provider{}` 块(应剥离)、`random_integer`
(额外 provider 应声明)、名称带插值(应保留不改写)、heredoc JSON(brace 不应干扰剥离)。

## Example Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_corpusrandom" "default" {
  policy_name     = "tf-example-${random_integer.default.result}"
  policy_document = <<EOF
{
  "Statement": [
    { "Effect": "Allow", "Action": "*" }
  ],
  "Version": "1"
}
EOF
}
```

## Argument Reference

The following arguments are supported:
* `policy_name` - (Required) 名称。
* `policy_document` - (Required) 策略文档。

## Attributes Reference

* `id` - 资源 ID。

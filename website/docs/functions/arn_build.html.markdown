---
subcategory: ""
layout: "alicloud"
page_title: "Alibaba Cloud: arn_build"
description: |-
  Builds an ARN from its constituent parts.
---

# Function: arn_build

Builds an ARN from its constituent parts.

See the [Alibaba Cloud documentation](https://www.alibabacloud.com/help/en/ram/policy-elements) for additional information on Alibaba Cloud Resource Names (ARNs).

-> **NOTE:** Available since v2.0.0-beta3.

## Example Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

# result: acs:ram::123456789012****:role/example
output "role_arn" {
  value = provider::alicloud::arn_build("ram", "", "123456789012****", "role/example")
}

# result: acs:ecs:cn-hangzhou:123456789012****:instance/i-bp1234567890abcdef
output "instance_arn" {
  value = provider::alicloud::arn_build("ecs", "cn-hangzhou", "123456789012****", "instance/i-bp1234567890abcdef")
}

# result: acs:oss:*:*:my-bucket/*
output "bucket_arn" {
  value = provider::alicloud::arn_build("oss", "*", "*", "my-bucket/*")
}
```

## Signature

```text
arn_build(ram_code string, region string, account_id string, relative_id string) string
```

## Arguments

1. `ram_code` (String) RAM code of the Alibaba Cloud service, such as `ecs`, `oss`, `ram`, `fc`, `mns`, `fnf`, or `kms`.
1. `region` (String) Region of the resource, such as `cn-hangzhou`. Use `*` to match every region, or an empty string for services that carry no region component, such as `ram`.
1. `account_id` (String) ID of the Alibaba Cloud account.
1. `relative_id` (String) Service-specific resource path, typically composed of a resource type and identifier, such as `instance/i-bp1234567890abcdef`.

---
subcategory: ""
layout: "alicloud"
page_title: "Alibaba Cloud: arn_parse"
description: |-
  Parses an ARN into its constituent parts.
---

# Function: arn_parse

Parses an ARN into its constituent parts.

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

locals {
  instance = provider::alicloud::arn_parse("acs:ecs:cn-hangzhou:123456789012****:instance/i-bp1234567890abcdef")
}

# result: ecs
output "service" {
  value = local.instance.service
}

# result: cn-hangzhou
output "region" {
  value = local.instance.region
}

# result: 123456789012****
output "account_id" {
  value = local.instance.account_id
}

# result: instance/i-bp1234567890abcdef
output "resource" {
  value = local.instance.resource
}
```

The returned attributes are the arguments of [`arn_build`](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/functions/arn_build), so an ARN can be taken apart, adjusted and put back together:

```terraform
locals {
  bucket = provider::alicloud::arn_parse("acs:oss:cn-hangzhou:123456789012****:my-bucket")
}

# result: acs:oss:cn-hangzhou:123456789012****:my-bucket/logs/*
output "prefix_arn" {
  value = provider::alicloud::arn_build(local.bucket.service, local.bucket.region, local.bucket.account_id, "${local.bucket.resource}/logs/*")
}
```

## Signature

```text
arn_parse(arn string) object
```

## Arguments

1. `arn` (String) ARN to parse, such as `acs:ecs:cn-hangzhou:123456789012****:instance/i-bp1234567890abcdef`.

The ARN must have five colon-separated sections and must start with `acs`. The region and the account ID may be empty, because a global service such as `ram` carries no region and a policy may match every account, but the resource may not: `acs:ram:*::*` parses, `acs:ecs:cn-hangzhou:123456789012****:` does not. Any colons beyond the fourth belong to the resource and are kept as part of it, so `acs:log:cn-hangzhou:123456789012****:project/my-project:logstore/my-logstore` yields a resource of `project/my-project:logstore/my-logstore`.

Unlike `arn_build`, which formats whatever it is given, this function reports an error for an ARN that is not of that form.

-> **NOTE:** The sections are not otherwise validated: an ARN naming a service or region that does not exist parses successfully.

## Return Type

An object with the following attributes:

* `service` (String) RAM code of the Alibaba Cloud service, such as `ecs`, `oss`, `ram`, `fc`, `mns`, `fnf`, or `kms`.
* `region` (String) Region of the resource, such as `cn-hangzhou`. Empty for services that carry no region component, such as `ram`.
* `account_id` (String) ID of the Alibaba Cloud account.
* `resource` (String) Service-specific resource path, typically composed of a resource type and identifier, such as `instance/i-bp1234567890abcdef`.

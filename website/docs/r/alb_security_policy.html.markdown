---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_security_policy"
sidebar_current: "docs-alicloud-resource-alb-security-policy"
description: |-
  Provides a Alicloud ALB Security Policy resource.
---

# alicloud\_alb\_security\_policy

Provides a ALB Security Policy resource.

For information about ALB Security Policy and how to use it, see [What is Security Policy](https://www.alibabacloud.com/help/doc-detail/213607.htm).

-> **NOTE:** Available in v1.130.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "testAccSecurityPolicy"
}

resource "alicloud_alb_security_policy" "default" {
  security_policy_name = var.name
  tls_versions         = ["TLSv1.0"]
  ciphers              = ["ECDHE-ECDSA-AES128-SHA", "AES256-SHA"]
}

```

## Argument Reference

The following arguments are supported:

* `resource_group_id` - (Optional) The ID of the resource group.
* `security_policy_name` - (Required) The name of the resource. The name must be 2 to 128 characters in length and must start with a letter. It can contain digits, periods (.), underscores (_), and hyphens (-).
* `dry_run` - (Optional) The dry run.
* `tls_versions` - (Required) The TLS protocol versions that are supported. Valid values: TLSv1.0, TLSv1.1, TLSv1.2 and TLSv1.3.
* `ciphers` - (Required) The supported cipher suites, which are determined by the TLS protocol version.The specified cipher suites must be supported by at least one TLS protocol version that you select. 
* `tags` - (Optional) A mapping of tags to assign to the resource.
## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Security Policy.
* `status` - The status of the resource.

## Import

ALB Security Policy can be imported using the id, e.g.

```
$ terraform import alicloud_alb_security_policy.example <id>
```

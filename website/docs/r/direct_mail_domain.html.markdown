---
subcategory: "Direct Mail"
layout: "alicloud"
page_title: "Alicloud: alicloud_direct_mail_domain"
sidebar_current: "docs-alicloud-resource-direct-mail-domain"
description: |-
  Provides a Alicloud Direct Mail Domain resource.
---

# alicloud\_direct\_mail\_domain

Provides a Direct Mail Domain resource.

For information about Direct Mail Domain and how to use it, see [What is Domain](https://www.alibabacloud.com/help/en/doc-detail/29414.htm).

-> **NOTE:** Available in v1.134.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_direct_mail_domain" "example" {
  domain_name = "example_value"
}

```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, ForceNew) Domain, length `1` to `50`, including numbers or capitals or lowercase letters or `.` or `-`

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Domain.
* `status` - The status of the domain name. Valid values:`0` to `4`. `0`:Available, Passed. `1`: Unavailable, No passed. `2`: Available, cname no passed, icp no passed. `3`: Available, icp no passed. `4`: Available, cname no passed.

## Import

Direct Mail Domain can be imported using the id, e.g.

```
$ terraform import alicloud_direct_mail_domain.example <id>
```

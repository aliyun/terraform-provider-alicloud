---
subcategory: "Web Application Firewall(WAF)"
layout: "alicloud"
page_title: "Alicloud: alicloud_waf_certificate"
sidebar_current: "docs-alicloud-resource-waf-certificate"
description: |-
  Provides a Web Application Firewall Certificate resource.
---

# alicloud\_waf\_certificate

Provides a WAF Certificate resource.

For information about WAF Certificate and how to use it, see [What is Certificate](https://www.alibabacloud.com/help/doc-detail/28517.htm).

-> **NOTE:** Available in v1.135.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_waf_certificate" "default" {
  certificate_name = "your_certificate_name"
  instance_id      = "your_instance_id"
  domain           = "your_domain_name"
  private_key      = "your_private_key"
  certificate      = "your_certificate"
}
```

## Argument Reference

The following arguments are supported:

* `certificate` - (Optional) Certificate file content.
* `certificate_name` - (Required, ForceNew) Certificate file name.
* `instance_id` - (Required, ForceNew) The ID of the WAF instance.
* `private_key` - (Required) The private key.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Certificate. The value formats as `<instance_id>:<domain>:<certificate_id>`.
* `certificate_id` - Certificate recording ID.
* `domain` - The domain that you want to add to WAF.

## Import

WAF Certificate can be imported using the id, e.g.

```
$ terraform import alicloud_waf_certificate.example <instance_id>:<domain>:<certificate_id>
```
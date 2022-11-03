---
subcategory: "DCDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_dcdn_waf_domain"
sidebar_current: "docs-alicloud-resource-dcdn-waf-domain"
description: |-
  Provides a Alicloud DCDN Waf Domain resource.
---

# alicloud\_dcdn\_waf\_domain

Provides a DCDN Waf Domain resource.

For information about DCDN Waf Domain and how to use it, see [What is Waf Domain](https://www.alibabacloud.com/help/en/dynamic-route-for-cdn/latest/batch-configure-domain-name-protection).

-> **NOTE:** Available in v1.185.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_dcdn_domain" "default" {
  domain_name = var.domain_name
  sources {
    content  = "1.1.1.1"
    port     = "80"
    priority = "20"
    type     = "ipaddr"
  }
}
resource "alicloud_dcdn_waf_domain" "default" {
  domain_name   = alicloud_dcdn_domain.default.domain_name
  client_ip_tag = "X-Forwarded-For"
}
```

## Argument Reference

The following arguments are supported:

* `client_ip_tag` - (Optional) The client ip tag.
* `domain_name` - (Required, ForceNew) The accelerated domain name.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Waf Domain. Its value is same as `domain_name`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the Waf Domain.
* `delete` - (Defaults to 1 mins) Used when deleting the Waf Domain.
* `update` - (Defaults to 1 mins) Used when updating the Waf Domain.

## Import

DCDN Waf Domain can be imported using the id, e.g.

```shell
$ terraform import alicloud_dcdn_waf_domain.example <domain_name>
```
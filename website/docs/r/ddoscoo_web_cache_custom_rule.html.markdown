---
subcategory: "Anti-DDoS Pro (DdosCoo)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ddoscoo_web_cache_custom_rule"
description: |-
  Provides a Alicloud DdosCoo Web Cache Custom Rule resource.
---

# alicloud_ddoscoo_web_cache_custom_rule

Provides a DdosCoo Web Cache Custom Rule resource.

Static page cache custom rules for a website domain protected by Anti-DDoS Pro (DdosCoo).

For information about DdosCoo Web Cache Custom Rule and how to use it, see [ModifyWebCacheCustomRule](https://next.api.alibabacloud.com/document/ddoscoo/2020-01-01/ModifyWebCacheCustomRule).

-> **NOTE:** Available since v1.271.0.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "terraform"
}

variable "domain" {
  default = "terraform-example.alibaba.com"
}

data "alicloud_ddoscoo_instances" "default" {
}

resource "alicloud_ddoscoo_domain_resource" "default" {
  domain       = var.domain
  instance_ids = [data.alicloud_ddoscoo_instances.default.ids.0]
  proxy_types {
    proxy_ports = [443]
    proxy_type  = "https"
  }
  real_servers = ["177.167.32.11"]
  rs_type      = 0
}

resource "alicloud_ddoscoo_web_cache_custom_rule" "default" {
  domain    = alicloud_ddoscoo_domain_resource.default.id
  name      = var.name
  uri       = "/blog/"
  mode      = "standard"
  cache_ttl = 3600
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, ForceNew) The website domain that the static page cache custom rule belongs to. The domain must already be configured with a website forwarding rule and associated with an enhanced package DDoS Pro instance.
* `name` - (Required, ForceNew) The name of the custom cache rule. Together with `domain` it forms the resource ID.
* `uri` - (Required) The path of the cached page, e.g. `/blog/`.
* `mode` - (Required) The cache mode. Valid values: `standard` (standard mode), `aggressive` (aggressive mode) and `bypass` (no caching).
* `cache_ttl` - (Required) The expiration time of the cached page, in seconds.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, in the format `<domain>:<name>`.
* `request_id` - The request ID returned by the most recent `ModifyWebCacheCustomRule` call, for traceability. It is refreshed on every apply and is not stable across reads.

## Import

DdosCoo Web Cache Custom Rule can be imported using the id, e.g. `<domain>:<name>`, as the import ID.

```shell
$ terraform import alicloud_ddoscoo_web_cache_custom_rule.example www.example.com:my-rule
```

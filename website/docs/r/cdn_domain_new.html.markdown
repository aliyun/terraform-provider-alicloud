---
layout: "alicloud"
page_title: "Alicloud: alicloud_cdn_domain_new"
sidebar_current: "docs-alicloud-resource-cdn-domain-new"
description: |-
  Provides a Alicloud Cdn Domain New Resource.
---

# alicloud_cdn_domain_new

Provides a CDN Accelerated Domain resource. This resource is based on CDN's new version OpenAPI.

For information about Cdn Domain New and how to use it, see [Add a domain](https://www.alibabacloud.com/help/doc-detail/91176.html).

-> **NOTE:** Available in v1.34.0+.

## Example Usage

Basic Usage

```
# Create a new Domain.
resource "alicloud_cdn_domain_new" "domain" {
  domain_name = "terraform.test.com"
  cdn_type    = "web"
  scope       = "overseas"
  sources {
    content  = "1.1.1.1"
    type     = "ipaddr"
    priority = 20
    port     = 80
    weight   = 10
  }
}

```
## Argument Reference

The following arguments are supported:

* `domain_name` - (Required) Name of the accelerated domain. This name without suffix can have a string of 1 to 63 characters, must contain only alphanumeric characters or "-", and must not begin or end with "-", and "-" must not in the 3th and 4th character positions at the same time. Suffix `.sh` and `.tel` are not supported.
* `cdn_type` - (Required, ForceNew) Cdn type of the accelerated domain. Valid values are `web`, `download`, `video`.
* `scope` - (Optional) Scope of the accelerated domain. Valid values are `domestic`, `overseas`, `global`. Default value is `domestic`. This parameter's setting is valid Only for the international users and domestic L3 and above users .
* `sources` - (Optional, Type: list) The source address list of the accelerated domain. Defaults to null. See Block Sources.

### Block sources

The `sources` block supports the following:

* `content` - (Required) The adress of source. Valid values can be ip or doaminName. Each item's `content` can not be repeated.
* `type` - (Required) The type of the source. Valid values are `ipaddr`, `domain` and `oss`.
* `port` - (Optional, Type: int) The port of source. Valid values are `443` and `80`. Default value is `80`.
* `priority` - (Optional, Type: int) Priority of the source. Valid values are `0` and `100`. Default value is `20`.
* `weight` - (Optional, Type: int) Weight of the source. Valid values are from `0` to `100`. Default value is `10`, but if type is `ipaddr`, the value can only be `10`. 

## Attributes Reference

The following attributes are exported:

* `id` - The cdn domain id. The value is same as the domain name.

## Import

CDN domain can be imported using the id, e.g.

```
terraform import alicloud_cdn_domain_new.example cdn-abc123456
```

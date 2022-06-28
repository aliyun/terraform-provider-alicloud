---
subcategory: "DNS"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_domain_attachment"
sidebar_current: "docs-alicloud-resource-alidns-domain-attachment"
description: |-
  Provides bind the domain name to the Alidns instance resource.
---

# alicloud\_alidns\_domain\_attachment

Provides bind the domain name to the Alidns instance resource.

-> **NOTE:** Available in v1.99.0+.

## Example Usage

```terraform
resource "alicloud_alidns_domain_attachment" "dns" {
  instance_id  = "dns-cn-mp91lyq9xxxx"
  domain_names = ["test111.abc", "test222.abc"]
}
```
## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The id of the DNS instance.
* `domain_names` - (Required) The domain names bound to the DNS instance.

## Attributes Reference

The following attributes are exported:

* `id` - This ID of this resource. The value is same as `instance_id`. 
* `domain_names` - Domain names bound to DNS instance.

## Import

DNS domain attachment can be imported using the id, e.g.

```
$ terraform import alicloud_alidns_domain_attachment.example dns-cn-v0h1ldjhxxx
```

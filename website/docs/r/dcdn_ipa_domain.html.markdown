---
subcategory: "DCDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_dcdn_ipa_domain"
sidebar_current: "docs-alicloud-resource-dcdn-ipa-domain"
description: |-
  Provides a Alicloud DCDN Ipa Domain resource.
---

# alicloud\_dcdn\_ipa\_domain

Provides a DCDN Ipa Domain resource.

For information about DCDN Ipa Domain and how to use it, see [What is Ipa Domain](https://www.alibabacloud.com/help/en/doc-detail/130634.html).

-> **NOTE:** Available in v1.158.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "default"
}
resource "alicloud_dcdn_ipa_domain" "example" {
  domain_name       = "example.com"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  sources {
    content  = "1.1.1.1"
    port     = 80
    priority = "20"
    type     = "ipaddr"
    weight   = 10
  }
  scope  = "overseas"
  status = "online"
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, ForceNew) The domain name to be added to IPA. Wildcard domain names are supported. A wildcard domain name must start with a period (.).
* `resource_group_id` - (Optional, Computed) The ID of the resource group. If you do not set this parameter, the system automatically assigns the ID of the default resource group.
* `scope` - (Optional, Computed, ForceNew) The accelerated region. Valid values: `domestic`, `global`, `overseas`.
* `sources` - (Required) Sources. See the following `Block sources`.
* `status` - (Optional, Computed) The status of DCDN Ipa Domain. Valid values: `online`, `offline`. Default to `online`.

#### Block sources

The sources supports the following: 

* `content` - (Required) The address of the origin server. You can specify an IP address or a domain name.
* `port` - (Required) The custom port number. Valid values: `0` to `65535`.
* `priority` - (Required) The priority of the origin server. Valid values: `20` and `30`. Default value: `20`. A value of 20 specifies that the origin is a primary origin. A value of 30 specifies that the origin is a secondary origin.
* `type` - (Required) The type of the origin server. Valid values: `ipaddr`, `domain`, `oss`.
* `weight` - (Required) The weight of the origin server. You must specify a value that is less than `100`. Default value: `10`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Ipa Domain. Its value is same as `domain_name`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when Creating DCDN Ipa domain instance.
* `update` - (Defaults to 5 mins) Used when Creating DCDN Ipa domain instance.
* `delete` - (Defaults to 10 mins) Used when terminating the DCDN Ipa domain instance.


## Import

DCDN Ipa Domain can be imported using the id, e.g.

```
$ terraform import alicloud_dcdn_ipa_domain.example <domain_name>
```
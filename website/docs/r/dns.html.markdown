---
subcategory: "Alidns"
layout: "alicloud"
page_title: "Alicloud: alicloud_dns"
sidebar_current: "docs-alicloud-resource-dns"
description: |-
  Provides a DNS resource.
---

# alicloud\_dns

-> **DEPRECATED:** This resource has been renamed to [alicloud_alidns_domain](https://www.terraform.io/docs/providers/alicloud/r/alidns_domain) from version 1.95.0.

Provides a DNS resource.

-> **NOTE:** The domain name which you want to add must be already registered and had not added by another account. Every domain name can only exist in a unique group.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_dns&exampleId=063bc164-a125-fcd7-f4c6-b79efb0c321933b85481&activeTab=example&spm=docs.r.dns.0.063bc164a1&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
# Add a new Domain.
resource "alicloud_dns" "dns" {
  name     = "starmove.com"
  group_id = "85ab8713-4a30-4de4-9d20-155ff830f651"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) Name of the domain. This name without suffix can have a string of 1 to 63 characters, must contain only alphanumeric characters or "-", and must not begin or end with "-", and "-" must not in the 3th and 4th character positions at the same time. Suffix `.sh` and `.tel` are not supported.
* `group_id` - (Optional) Id of the group in which the domain will add. If not supplied, then use default group.
* `resource_group_id` - (Optional, ForceNew, Available in 1.59.0+) The Id of resource group which the dns belongs.

## Attributes Reference

The following attributes are exported:

* `id` - This ID of this resource. The value is set to `domain_name`.
* `domain_id` - The domain ID.
* `name` - The domain name.
* `group_id` - The group id of domain.
* `dns_server` - A list of the dns server name.

## Import

DNS can be imported using the id or domain name, e.g.

```shell
$ terraform import alicloud_dns.example "aliyun.com"
```

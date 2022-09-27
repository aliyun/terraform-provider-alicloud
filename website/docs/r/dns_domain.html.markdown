---
subcategory: "DNS"
layout: "alicloud"
page_title: "Alicloud: alicloud_dns_domain"
sidebar_current: "docs-alicloud-resource-dns-domain"
description: |-
  Provides a DNS domain resource.
---

# alicloud\_dns\_domain

Provides a DNS domain resource.

-> DEPRECATED: This resource has been renamed to [alicloud_alidns_domain](https://www.terraform.io/docs/providers/alicloud/r/alidns_domain) from version 1.95.0.

-> **NOTE:** The domain name which you want to add must be already registered and had not added by another account. Every domain name can only exist in a unique group.

-> **NOTE:** Available in v1.81.0+.

## Example Usage

```terraform 
# Add a new Domain.
resource "alicloud_dns_domain" "dns" {
  domain_name = "starmove.com"
  group_id    = "85ab8713-4a30-4de4-9d20-155ff830****"
  tags = {
    Created     = "Terraform"
    Environment = "test"
  }
}
```
## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, ForceNew) Name of the domain. This name without suffix can have a string of 1 to 63 characters(domain name subject, excluding suffix), must contain only alphanumeric characters or "-", and must not begin or end with "-", and "-" must not in the 3th and 4th character positions at the same time. Suffix `.sh` and `.tel` are not supported.
* `group_id` - (Optional) Id of the group in which the domain will add. If not supplied, then use default group.
* `resource_group_id` - (Optional, ForceNew) The Id of resource group which the dns domain belongs.
* `lang` - (Optional) User language.
* `remark` - (Optional) Remarks information for your domain name.
* `tags` - (Optional) A mapping of tags to assign to the resource.
    - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
    - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.

## Attributes Reference

The following attributes are exported:

* `id` - This ID of this resource. The value is set to `domain_name`.
* `domain_id` - The domain ID.
* `dns_server` - A list of the dns server name.

## Import

DNS domain can be imported using the id or domain name, e.g.

```
$ terraform import alicloud_dns_domain.example aliyun.com
```

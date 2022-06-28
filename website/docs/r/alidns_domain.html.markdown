---
subcategory: "DNS"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_domain"
sidebar_current: "docs-alicloud-resource-alidns-domain"
description: |-
  Provides a Alidns domain resource.
---

# alicloud\_alidns\_domain

Provides a Alidns domain resource.

-> **NOTE:** The domain name which you want to add must be already registered and had not added by another account. Every domain name can only exist in a unique group.

-> **NOTE:** Available in v1.95.0+.

## Example Usage

```terraform
# Add a new Domain.
resource "alicloud_alidns_domain" "dns" {
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
    - Key: It can be [1, 20] characters in length. It can contain A-Z, a-z, numbers, underscores (_), and hyphens (-). It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
    - Value: It can be [1, 20] characters in length. It can contain A-Z, a-z, numbers, underscores (_), and hyphens (-). It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.

## Attributes Reference

The following attributes are exported:

* `id` - This ID of this resource. The value is set to `domain_name`.
* `domain_id` - The domain ID.
* `dns_servers` - A list of the dns server name.
* `group_name` - Domain name group name.
* `puny_code` - Only return punycode codes for Chinese domain names.

### Timeouts

-> **NOTE:** Available in 1.97.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `delete` - (Defaults to 6 mins) Used when terminating the Alidns domain instance.  

## Import

Alidns domain can be imported using the id or domain name, e.g.

```
$ terraform import alicloud_alidns_domain.example aliyun.com
```

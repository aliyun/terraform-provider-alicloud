---
subcategory: "Alidns"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_domain"
sidebar_current: "docs-alicloud-resource-alidns-domain"
description: |-
  Provides a Alidns domain resource.
---

# alicloud_alidns_domain

Provides a Alidns domain resource.

-> **NOTE:** The domain name which you want to add must be already registered and had not added by another account. Every domain name can only exist in a unique group.

-> **NOTE:** Available since v1.95.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_alidns_domain&exampleId=b0eea6b5-3a36-f056-c28d-d77f58768eccacc1eca1&activeTab=example&spm=docs.r.alidns_domain.0.b0eea6b53a&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_alidns_domain_group" "default" {
  domain_group_name = "tf-example"
}
resource "alicloud_alidns_domain" "default" {
  domain_name = "starmove.com"
  group_id    = alicloud_alidns_domain_group.default.id
  tags = {
    Created = "TF",
    For     = "example",
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

## Timeouts

-> **NOTE:** Available since v1.97.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `delete` - (Defaults to 6 mins) Used when terminating the Alidns domain instance.  

## Import

Alidns domain can be imported using the id or domain name, e.g.

```shell
$ terraform import alicloud_alidns_domain.example aliyun.com
```

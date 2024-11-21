---
subcategory: "Alidns"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_domain_attachment"
sidebar_current: "docs-alicloud-resource-alidns-domain-attachment"
description: |-
  Provides bind the domain name to the Alidns instance resource.
---

# alicloud_alidns_domain_attachment

Provides bind the domain name to the Alidns instance resource.

-> **NOTE:** Available since v1.99.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_alidns_domain_attachment&exampleId=e14db959-f1fa-2e75-dfc2-433b7b48d98b3f2729cb&activeTab=example&spm=docs.r.alidns_domain_attachment.0.e14db959f1&intl_lang=EN_US" target="_blank">
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

resource "alicloud_alidns_instance" "default" {
  dns_security   = "basic"
  domain_numbers = 3
  version_code   = "version_personal"
  period         = 1
  renewal_status = "ManualRenewal"
}

resource "alicloud_alidns_domain_attachment" "default" {
  instance_id  = alicloud_alidns_instance.default.id
  domain_names = [alicloud_alidns_domain.default.domain_name]
}
```
## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The id of the DNS instance.
* `domain_names` - (Required) The domain names bound to the DNS instance.

## Attributes Reference

The following attributes are exported:

* `id` - This ID of this resource. The value is same as `instance_id`. 

## Import

DNS domain attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_alidns_domain_attachment.example dns-cn-v0h1ldjhxxx
```

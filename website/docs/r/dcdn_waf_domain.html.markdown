---
subcategory: "DCDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_dcdn_waf_domain"
sidebar_current: "docs-alicloud-resource-dcdn-waf-domain"
description: |-
  Provides a Alicloud DCDN Waf Domain resource.
---

# alicloud_dcdn_waf_domain

Provides a DCDN Waf Domain resource.

For information about DCDN Waf Domain and how to use it, see [What is Waf Domain](https://www.alibabacloud.com/help/en/dcdn/developer-reference/api-dcdn-2018-01-15-batchsetdcdnwafdomainconfigs).

-> **NOTE:** Available since v1.185.0.

## Example Usage
<div class="oics-button" style="float: right;margin: 0 0 -40px 0;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_dcdn_waf_domain&exampleId=59d20c11-eddc-d7c7-931e-2a917c85f34c2841b8dc&activeTab=example&spm=docs.r.dcdn_waf_domain.0.59d20c11ed" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>

Basic Usage

```terraform
variable "domain_name" {
  default = "tf-example.com"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_dcdn_domain" "example" {
  domain_name = "${var.domain_name}-${random_integer.default.result}"
  scope       = "overseas"
  sources {
    content  = "1.1.1.1"
    port     = "80"
    priority = "20"
    type     = "ipaddr"
    weight   = "10"
  }
}

resource "alicloud_dcdn_waf_domain" "example" {
  domain_name   = alicloud_dcdn_domain.example.domain_name
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

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the Waf Domain.
* `delete` - (Defaults to 1 mins) Used when deleting the Waf Domain.
* `update` - (Defaults to 1 mins) Used when updating the Waf Domain.

## Import

DCDN Waf Domain can be imported using the id, e.g.

```shell
$ terraform import alicloud_dcdn_waf_domain.example <domain_name>
```
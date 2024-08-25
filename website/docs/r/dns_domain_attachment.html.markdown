---
subcategory: "Alidns"
layout: "alicloud"
page_title: "Alicloud: alicloud_dns_domain_attachment"
sidebar_current: "docs-alicloud-resource-dns-domain-attachment"
description: |-
  Provides bind the domain name to the DNS instance resource.
---

# alicloud\_dns\_domain\_attachment

Provides bind the domain name to the DNS instance resource.

-> **NOTE:** Available in v1.80.0+.

-> **DEPRECATED:**  This resource has been deprecated from version `1.99.0`. Please use new resource [alicloud_alidns_domain_attachment](https://www.terraform.io/docs/providers/alicloud/r/alidns_domain_attachment).

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_dns_domain_attachment&exampleId=b0483a98-4ec0-cc26-1037-df40741f6a442df66981&activeTab=example&spm=docs.r.dns_domain_attachment.0.b0483a984e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_dns_domain_attachment" "dns" {
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

```shell
$ terraform import alicloud_dns_domain_attachment.example dns-cn-v0h1ldjhxxx
```

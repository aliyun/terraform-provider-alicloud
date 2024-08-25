---
subcategory: "Alidns"
layout: "alicloud"
page_title: "Alicloud: alicloud_dns_group"
sidebar_current: "docs-alicloud-resource-dns-group"
description: |-
  Provides a DNS Group resource.
---

# alicloud\_dns\_group

-> **DEPRECATED:**  This resource  has been deprecated from version `1.84.0`. Please use new resource [alicloud_alidns_domain_group](https://www.terraform.io/docs/providers/alicloud/r/alidns_domain_group).

Provides a DNS Group resource.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_dns_group&exampleId=2cdc16b4-9616-eae6-aff0-f762c2a1f9ec4d6ff2b3&activeTab=example&spm=docs.r.dns_group.0.2cdc16b496&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
# Add a new Domain group.
resource "alicloud_dns_group" "group" {
  name = "testgroup"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the domain group.    

## Attributes Reference

The following attributes are exported:

* `id` - The group id.
* `name` - The group name.

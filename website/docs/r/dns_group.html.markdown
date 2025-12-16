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

```terraform
# Add a new Domain group.
resource "alicloud_dns_group" "group" {
  name = "testgroup"
}

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_dns_group&spm=docs.r.dns_group.example&intl_lang=EN_US)
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the domain group.    

## Attributes Reference

The following attributes are exported:

* `id` - The group id.
* `name` - The group name.

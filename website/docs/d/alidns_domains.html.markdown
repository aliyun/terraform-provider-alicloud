---
subcategory: "Alidns"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_domains"
sidebar_current: "docs-alicloud-datasource-alidns-domains"
description: |-
    Provides a list of Alidns Domains available to the user.
---

# alicloud\_alidns\_domains

This data source provides a list of Alidns Domains in an Alibaba Cloud account according to the specified filters.

-> **NOTE:**  Available in 1.85.0+.

## Example Usage

```
data "alicloud_alidns_domains" "domains_ds" {
  name_regex   = "^hegu"
  output_file  = "domains.txt"
}

output "first_domain_id" {
  value = "${data.alicloud_alidns_domains.domains_ds.domains.0.domain_id}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter results by the domain name. 
* `ids` (Optional) - A list of domain IDs.
* `lang` (Optional) - User language.
* `group_id` (Optional) - Id of the domain group.
* `tags` - (Optional) A mapping of tags to assign to the resource.
      - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
      - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional) The Id of resource group which the Alidns Domain belongs.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of domain IDs.
* `names` - A list of domain names.
* `domains` - A list of domains. Each element contains the following attributes:
  * `id` - ID of the resource.
  * `domain_id` - ID of the domain.
  * `domain_name` - Name of the domain.
  * `group_id` - Id of group that contains the domain.
  * `tags` - A mapping of tags to assign to the resource.
  * `remark` - The remark of the domain.

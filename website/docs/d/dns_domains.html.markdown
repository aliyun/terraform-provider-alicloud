---
subcategory: "DNS"
layout: "alicloud"
page_title: "Alicloud: alicloud_dns_domains"
sidebar_current: "docs-alicloud-datasource-dns-domains"
description: |-
    Provides a list of domains available to the user.
---

# alicloud\_dns\_domains

-> **DEPRECATED:**  This resource has been renamed to [alicloud_alidns_domains](https://www.terraform.io/docs/providers/alicloud/d/alidns_domains) from version 1.95.0. 

This data source provides a list of DNS Domains in an Alibaba Cloud account according to the specified filters.

## Example Usage

```terraform
data "alicloud_dns_domains" "domains_ds" {
  domain_name_regex = "^hegu"
  output_file       = "domains.txt"
}

output "first_domain_id" {
  value = "${data.alicloud_dns_domains.domains_ds.domains.0.domain_id}"
}
```

## Argument Reference

The following arguments are supported:

* `domain_name_regex` - (Optional) A regex string to filter results by the domain name. 
* `group_name_regex` - (Optional)  A regex string to filter results by the group name.
* `ali_domain` - (Optional, type: bool) Specifies whether the domain is from Alibaba Cloud or not.
* `instance_id` - (Optional) Cloud analysis product ID.
* `version_code` - (Optional) Cloud analysis version code.
* `ids` (Optional, Available in 1.53.0+) - A list of domain IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional, ForceNew, Available in 1.59.0+) The Id of resource group which the dns belongs.
* `group_id` - (Optional, ForceNew, Available in 1.87.0+) Domain group ID, if not filled, the default is all groups.
* `key_word` - (Optional, ForceNew, Available in 1.87.0+) The keywords are searched according to the `%KeyWord%` mode, which is not case sensitive.
* `lang` - (Optional, ForceNew, Available in 1.87.0+) User language.
* `search_mode` - (Optional, ForceNew, Available in 1.87.0+) Search mode, `LIKE` fuzzy search, `EXACT` exact search.
* `starmark` - (Optional, ForceNew, Available in 1.87.0+) Whether to query the domain name star.
* `tags` - (Optional, ForceNew, Available in 1.87.0+) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of domain IDs.
* `names` - A list of domain names.
* `domains` - A list of domains. Each element contains the following attributes:
  * `domain_id` - ID of the domain.
  * `domain_name` - Name of the domain.
  * `ali_domain` - Indicates whether the domain is an Alibaba Cloud domain.
  * `group_id` - Id of group that contains the domain.
  * `group_name` - Name of group that contains the domain.
  * `instance_id` - Cloud analysis product ID of the domain.
  * `version_code` - Cloud analysis version code of the domain.
  * `puny_code` - Punycode of the Chinese domain.
  * `dns_servers` - DNS list of the domain in the analysis system.
  * `resource_group_id` - The Id of resource group which the dns belongs.
  * `dns_servers` - DNS list of domain names in the resolution system.
  * `id` - The Id of resource.
  * `in_black_hole` - Whether it is in black hole.
  * `in_clean` - Whether it is cleaning.
  * `min_ttl` - Minimum TTL.
  * `version_code` - Cloud resolution version ID.
  * `record_line_tree_json` - Tree-like analytical line list.
  * `region_lines` - Whether it is a regional route.
  * `remark` - The Id of resource group which the dns belongs.
  * `slave_dns` - Whether to allow auxiliary dns.

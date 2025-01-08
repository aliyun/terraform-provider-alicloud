---
subcategory: "Alidns"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_domains"
sidebar_current: "docs-alicloud-datasource-alidns-domains"
description: |-
    Provides a list of domains available to the user.
---

# alicloud_alidns_domains

This data source provides a list of Alidns Domains in an Alibaba Cloud account according to the specified filters.

-> **NOTE:** Available since v1.95.0.

## Example Usage

```terraform
data "alicloud_alidns_domains" "domains_ds" {
  domain_name_regex = "^hegu"
  output_file       = "domains.txt"
}

output "first_domain_id" {
  value = "${data.alicloud_alidns_domains.domains_ds.domains.0.domain_id}"
}
```

## Argument Reference

The following arguments are supported:

* `domain_name_regex` - (Optional, ForceNew) A regex string to filter results by the domain name. 
* `group_name_regex` - (Optional, ForceNew)  A regex string to filter results by the group name.
* `ali_domain` - (Optional, Bool, ForceNew) Specifies whether the domain is from Alibaba Cloud or not.
* `instance_id` - (Optional, ForceNew) Cloud analysis product ID.
* `version_code` - (Optional, ForceNew) Cloud analysis version code.
* `ids` (Optional, ForceNew) - A list of domain IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional, ForceNew, Available in 1.59.0+) The Id of resource group which the dns belongs.
* `group_id` - (Optional, ForceNew) Domain group ID, if not filled, the default is all groups.
* `key_word` - (Optional, ForceNew) The keywords are searched according to the `%KeyWord%` mode, which is not case sensitive.
* `lang` - (Optional, ForceNew) User language.
* `search_mode` - (Optional, ForceNew) Search mode, `LIKE` fuzzy search, `EXACT` exact search.
* `starmark` - (Optional, ForceNew) Whether to query the domain name star.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.

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
  * `version_name` - Cloud analysis version name of the domain. 
  * `puny_code` - Punycode of the Chinese domain.
  * `dns_servers` - DNS list of the domain in the analysis system.
  * `resource_group_id` - The Id of resource group which the dns belongs.
  * `dns_servers` - DNS list of domain names in the resolution system.
  * `id` - The Id of resource.
  * `in_black_hole` - Whether it is in black hole.
  * `line_type` - The type of the DNS request line.
  * `in_clean` - Whether it is cleaning.
  * `min_ttl` - Minimum TTL.
  * `version_code` - Cloud resolution version ID.
  * `record_line_tree_json` - Tree-like analytical line list.
  * `region_lines` - Whether it is a regional route.
  * `remark` - The Id of resource group which the dns belongs.
  * `slave_dns` - Whether to allow auxiliary dns.
  * `available_ttls` - List of available TTLs.
  * `tags` - Tags of the domain.
  * `record_lines` - Parse the line data list.
    * `father_code` - The code of the parent line, or empty if there is none.
    * `line_display_name` - Parent line display name.
    * `line_code` - Sub-line Code.
    * `line_name` - Sub-line display name.


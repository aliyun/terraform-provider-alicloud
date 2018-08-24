---
layout: "alicloud"
page_title: "Alicloud: alicloud_dns_domains"
sidebar_current: "docs-alicloud-datasource-dns-domains"
description: |-
    Provides a list of domains available to the user.
---

# alicloud\_dns\_domains

This data source provides a list of DNS Domains in an Alibaba Cloud account according to the specified filters.

## Example Usage

```
data "alicloud_dns_domains" "domains_ds" {
  domain_name_regex = "^hegu"
  output_file = "domains.txt"
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
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

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
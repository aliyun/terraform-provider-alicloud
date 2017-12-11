---
layout: "alicloud"
page_title: "Alicloud: alicloud_dns_domains"
sidebar_current: "docs-alicloud-datasource-dns-domains"
description: |-
    Provides a list of domains available to the user.
---

# alicloud\_dns\_domains

The Dns Domains data source provides a list of Alicloud Dns Domains in an Alicloud account according to the specified filters.

## Example Usage

```
data "alicloud_dns_domains" "domain" {
  domain_name_regex = "^hegu"
  output_file = "domains.txt"
}

```

## Argument Reference

The following arguments are supported:

* `domain_name_regex` - (Optional) A regex string to apply to the domain list returned by Alicloud. 
* `group_name_regex` - (Optional)  Limit search to provide group name regex.
* `ali_domain` - (Optional, type: bool) Limit search to specific whether it is Alicloud domain.
* `instance_id` - (Optional) Limit search to specific cloud analysis product ID.
* `version_code` - (Optional) Limit search to specific cloud analysis version code.
* `output_file` - (Optional) The name of file that can save domains data source after running `terraform plan`.

## Attributes Reference

A list of domains will be exported and its every element contains the following attributes:

* `domain_id` - ID of the domain.
* `domain_name` - Name of the domain.
* `ali_domain` - Indicates whether the domain is Alicloud domain.
* `group_id` - Id of group which the domain in.
* `group_name` - Name of group which the domain in.
* `instance_id` - Cloud analysis product id of the domain.
* `version_code` - Cloud analysis version code of the domain.
* `puny_code` - Punycode of the Chinese domain.
* `dns_servers` - DNS list of the domain in the analysis system.
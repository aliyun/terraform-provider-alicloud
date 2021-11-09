---
subcategory: "DCDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_dcdn_domains"
sidebar_current: "docs-alicloud-dcdn-domains"
description: |-
  Provides a collection of DCDN Domains to the specified filters.
---

# alicloud\_dcdn\_domains

Provides a collection of DCDN Domains to the specified filters.

~> **NOTE:** Available in 1.94.0+.

## Example Usage

 ```
data "alicloud_dcdn_domains" "example" {
  ids = ["example.com"]
}

output "domain_id" {
  value = data.alicloud_dcdn_domains.example.domains.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list ids of DCDN Domain.
* `name_regex` - (Optional) A regex string to filter results by the DCDN Domain.
* `change_end_time` - (Optional) The end time of the update. Specify the time in the ISO 8601 standard in the `yyyy-MM-ddTHH:mm:ssZ` format. The time must be in UTC.
* `change_start_time` - (Optional) The start time of the update. Specify the time in the ISO 8601 standard in the `yyyy-MM-ddTHH:mm:ssZ` format. The time must be in UTC.
* `check_domain_show` - (Optional) Specifies whether to display the domains in the checking, check_failed, or configure_failed status. Valid values: `true` or `false`.
* `domain_search_type` - (Optional) The search method. Default value: `fuzzy_match`. Valid values: `fuzzy_match`, `pre_match`, `suf_match`, `full_match`.
* `resource_group_id` - (Optional) The ID of the resource group.
* `status` - (Optional) The status of DCDN Domain.
* `enable_details` - (Optional) Default to `false`. Set it to true can output more details.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

 * `ids` - A list ids of DCDN Domain.
 * `names` - A list of DCDN Domain names.
 * `domains` - A list of domains. Each element contains the following attributes:
   * `id` - The ID of the DCDN Domain.
   * `cert_name` - Indicates the name of the certificate.
   * `resource_group_id` - The ID of the resource group.
   * `domain_name` - The name of the DCDN Domain.
   * `ssl_protocol` - Indicates whether the SSL certificate is enabled.
   * `ssl_pub` -  Indicates the public key of the certificate.
   * `scope` - The acceleration region.
   * `status` - The status of DCDN Domain. Valid values: `online`, `offline`, `check_failed`, `checking`, `configure_failed`, `configuring`.
   * `cname` - The canonical name (CNAME) of the accelerated domain.
   * `description` - The reason that causes the review failure.
   * `gmt_modified` - The time when the accelerated domain was last modified.
   * `sources` - The origin information.
     * `content` - The origin address.
     * `type` - The type of the origin. Valid values:
     * `port` - The port number.
     * `priority` - The priority of the origin if multiple origins are specified.
     * `weight` - The weight of the origin if multiple origins are specified.
     * `enabled` - The status of the origin.

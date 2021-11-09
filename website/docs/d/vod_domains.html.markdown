---
subcategory: "ApsaraVideo VoD"
layout: "alicloud"
page_title: "Alicloud: alicloud_vod_domains"
sidebar_current: "docs-alicloud-datasource-vod-domains"
description: |-
  Provides a list of Vod Domains to the user.
---

# alicloud\_vod\_domains

This data source provides the Vod Domains of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.136.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_vod_domain" "default" {
  domain_name = "your_domain_name"
  scope       = "domestic"
  sources {
    source_type    = "domain"
    source_content = "your_source_content"
    source_port    = "80"
  }
  tags = {
    key1 = "value1"
    key2 = "value2"
  }
}

data "alicloud_vod_domains" "default" {
  ids = [alicloud_vod_domain.default.id]
  tags = {
    key1 = "value1"
    key2 = "value2"
  }
}
output "vod_domain" {
  value = data.alicloud_vod_domains.default.domains.0
}

```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional, ForceNew) A regex string to filter results by Domain name.
* `ids` - (Optional, ForceNew, Computed)  A list of Domain IDs. Its element value is same as Domain Name.
* `domain_search_type` - (Optional, ForceNew) The search method. Valid values:
  * `fuzzy_match`: fuzzy match. This is the default value.
  * `pre_match`: prefix match.
  * `suf_match`: suffix match.
  * `full_match`: exact match
* `status` - (Optional, ForceNew) The status of the domain name. The value of this parameter is used as a condition to filter domain names. Value values:
* `tags` - (Optional) A mapping of tags to assign to the resource.
  * `Key`: It can be up to 64 characters in length. It cannot be a null string. 
  * `Value`: It can be up to 128 characters in length. It can be a null string.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
  * `onine`: indicates that the domain name is enabled.
  * `offline`: indicates that the domain name is disabled.
  * `configuring`: indicates that the domain name is being configured.
  * `configure_failed`: indicates that the domain name failed to be configured.
  * `checking`: indicates that the domain name is under review.
  * `check_failed`: indicates that the domain name failed the review.


## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Domain names.
* `domains` - A list of Vod Domains. Each element contains the following attributes:
  * `gmt_created` - The time when the domain name for CDN was added. The time follows the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time is displayed in UTC.
  * `gmt_modified` - The last time when the domain name for CDN was modified. The time follows the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time is displayed in UTC.
  * `cname` - The CNAME that is assigned to the domain name for CDN. You must add a CNAME record in the system of your Domain Name System (DNS) service provider to map the domain name for CDN to the CNAME.
  * `description` - The description of the domain name for CDN.
  * `ssl_protocol` - Indicates whether the Secure Sockets Layer (SSL) certificate is enabled. Valid values: `on`,`off`.
  * `domain_name` - The domain name for CDN.
  * `sources` - The information about the address of the origin server. For more information about the Sources parameter, See the following `Block sources`.
  * `status` - The status of the resource.
  * `id` - The ID of the Domain. Its value is same as Queue Name.
  * `sand_box` - Indicates whether the domain name for CDN is in a sandbox environment.

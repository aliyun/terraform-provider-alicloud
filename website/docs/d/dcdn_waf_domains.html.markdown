---
subcategory: "DCDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_dcdn_waf_domains"
sidebar_current: "docs-alicloud-datasource-dcdn-waf-domains"
description: |-
  Provides a list of Dcdn Waf Domains to the user.
---

# alicloud_dcdn_waf_domains

This data source provides the Dcdn Waf Domains of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.185.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_dcdn_waf_domains" "ids" {}
output "dcdn_waf_domain_id_1" {
  value = data.alicloud_dcdn_waf_domains.ids.domains.0.id
}
```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Waf Domain IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `query_args` - (Optional, ForceNew) The query conditions. You can filter domain names by name. Fuzzy match is supported `QueryArgs={"DomainName":"Accelerated domain name"}`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `domains` - A list of Dcdn Waf Domains. Each element contains the following attributes:
  * `defense_scenes` - Protection policy type.
    * `defense_scene` - The type of protection policy.
    * `policy_id` - The protection policy ID.
  * `domain_name` - The accelerated domain name.
  * `client_ip_tag` - The client ip tag.
  * `id` - The ID of the Waf Domain.
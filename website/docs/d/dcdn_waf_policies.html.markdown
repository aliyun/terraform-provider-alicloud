---
subcategory: "DCDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_dcdn_waf_policies"
sidebar_current: "docs-alicloud-datasource-dcdn-waf-policies"
description: |-
  Provides a list of Dcdn Waf Policies to the user.
---

# alicloud_dcdn_waf_policies

This data source provides the Dcdn Waf Policies of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.184.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_dcdn_waf_policies" "ids" {}
output "dcdn_waf_policy_id_1" {
  value = data.alicloud_dcdn_waf_policies.ids.policies.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Waf Policy IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `query_args` - (Optional, ForceNew) The query conditions. The value is a string in the JSON format. Format: `{"PolicyIds":"The ID of the proteuleIds":"Thection policy","R range of protection rule IDs","PolicyNameLike":"The name of the protection policy","DomainNames":"The protected domain names","PolicyType":"default","DefenseScenes":"waf_group","PolicyStatus":"on","OrderBy":"GmtModified","Desc":"false"}`.
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `on`, `off`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `policies` - A list of Dcdn Waf Policies. Each element contains the following attributes:
  * `dcdn_waf_policy_id` - The first ID of the resource.
  * `defense_scene` - The type of protection policy.
  * `domain_count` - The number of domain names that use this protection policy.
  * `gmt_modified` - The time when the protection policy was modified.
  * `id` - The ID of the Waf Policy.
  * `policy_name` - The name of the protection policy.
  * `policy_type` - The type of the protection policy.
  * `rule_count` - The number of protection rules in this protection policy.
  * `status` - The status of the resource.
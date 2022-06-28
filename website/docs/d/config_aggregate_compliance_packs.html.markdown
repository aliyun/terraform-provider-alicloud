---
subcategory: "Cloud Config"
layout: "alicloud"
page_title: "Alicloud: alicloud_config_aggregate_compliance_packs"
sidebar_current: "docs-alicloud-datasource-config-aggregate-compliance-packs"
description: |-
  Provides a list of Config Aggregate Compliance Packs to the user.
---

# alicloud\_config\_aggregate\_compliance\_packs

This data source provides the Config Aggregate Compliance Packs of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.124.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_config_aggregate_compliance_packs" "example" {
  aggregator_id = "ca-3a9b626622af001d****"
  ids           = ["cp-152a626622af00bc****"]
  name_regex    = "the_resource_name"
}

output "first_config_aggregate_compliance_pack_id" {
  value = data.alicloud_config_aggregate_compliance_packs.example.packs.0.id
}
```

## Argument Reference

The following arguments are supported:

* `aggregator_id` - (Required, ForceNew) The ID of aggregator.
* `name_regex` - (Optional)  A regex string to filter results by Aggregate Compliance Pack name.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Aggregate Compliance Pack IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the resource. Valid values `ACTIVE`, `CREATING`, `INACTIVE`. 


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Config Aggregate Compliance Pack names.
* `packs` - A list of Config Aggregate Compliance Packs. Each element contains the following attributes:
	* `account_id` - The Aliyun User Id.
	* `aggregator_compliance_pack_id` - The Aggregate Compliance Package Id.
	* `aggregate_compliance_pack_name` -The Aggregate Compliance Package Name.
	* `compliance_pack_template_id` - The template ID of the Compliance Package.
	* `config_rules` - A list of The Aggregate Compliance Package Rules.
		* `config_rule_id` - The ID of the rule.
		* `config_rule_parameters` - A list of parameter rules.
			* `parameter_name` - The Parameter Name.
			* `parameter_value` - The Parameter Value.
			* `required` - Required.
		* `managed_rule_identifier` - Managed Rule Identifier.
	* `description` - The description of aggregate compliance pack.
	* `id` - The ID of the Aggregate Compliance Pack.
	* `risk_level` - The Risk Level.
	* `status` - The status of the resource.

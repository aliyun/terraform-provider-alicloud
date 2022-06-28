---
subcategory: "Cloud Config"
layout: "alicloud"
page_title: "Alicloud: alicloud_config_compliance_packs"
sidebar_current: "docs-alicloud-datasource-config-compliance-packs"
description: |-
  Provides a list of Config Compliance Packs to the user.
---

# alicloud\_config\_compliance\_packs

This data source provides the Config Compliance Packs of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.124.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_config_compliance_packs" "example" {
  ids        = ["cp-152a626622af00bc****"]
  name_regex = "the_resource_name"
}

output "first_config_compliance_pack_id" {
  value = data.alicloud_config_compliance_packs.example.packs.0.id
}
```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Compliance Pack IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Compliance Pack name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the resource. Valid values `ACTIVE`, `CREATING`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Compliance Pack names.
* `packs` - A list of Config Compliance Packs. Each element contains the following attributes:
	* `account_id` - The Aliyun User Id.
	* `compliance_pack_id` - The Compliance Package ID.
	* `compliance_pack_name` - The Compliance Package Name.
	* `compliance_pack_template_id` - The template ID of the Compliance Package.
	* `config_rules` - A list of The Compliance Package Rules.
		* `config_rule_id` - The ID of the rule.
		* `config_rule_parameters` - A list of parameter rules.
			* `parameter_name` - The Parameter Name.
			* `parameter_value` - The Parameter Value.
			* `required` - Required.
		* `managed_rule_identifier` - Managed Rule Identifier.
	* `description` - The description of compliance pack.
	* `id` - The ID of the Compliance Pack.
	* `risk_level` - The Ris Level.
	* `status` - The status of the resource.

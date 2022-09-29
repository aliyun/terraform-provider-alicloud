---
subcategory: "DCDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_dcdn_waf_rules"
sidebar_current: "docs-alicloud-datasource-dcdn-waf-rules"
description: |-
  Provides a list of Dcdn Waf Rule owned by an Alibaba Cloud account.
---

# alicloud_dcdn_waf_rules

This data source provides Dcdn Waf Rule available to the user.[What is Waf Rule](https://www.alibabacloud.com/help/en/dynamic-route-for-cdn/latest/configure-protection-rules)

-> **NOTE:** Available in 1.201.0+

## Example Usage

```terraform
data "alicloud_dcdn_waf_rules" "default" {
  ids = ["${alicloud_dcdn_waf_rule.default.id}"]
}

output "alicloud_dcdn_waf_rule_example_id" {
  value = data.alicloud_dcdn_waf_rules.default.waf_rules.0.id
}
```

## Argument Reference

The following arguments are supported:
* `query_args` - (Optional, ForceNew) The query conditions. The value is a string in the JSON format.
* `ids` - (Optional, ForceNew, Computed) A list of Waf Rule IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `waf_rules` - A list of Waf Rule Entries. Each element contains the following attributes:
  * `defense_scene` - The type of protection policy. The following scenarios are supported:-waf_group:Web basic protection-custom_acl: Custom protection policy-whitelist: whitelist
  * `gmt_modified` - Revised the time. The date format is based on ISO8601 notation and uses UTC +0 time in the format of yyyy-MM-ddTHH:mm:ssZ.
  * `policy_id` - The protection policy ID.
  * `rule_name` - The name of the protection rule.
  * `status` - The status of the resource
  * `waf_rule_id` - The first ID of the resource
  * `id` - The ID of the Waf Rule.
  * `action` - Specifies the action of the rule.
  * `cn_region_list` - The blocked regions in the Chinese mainland, separated by commas (,).
  * `conditions` - The trigger condition of the rule.
    * `key` - The match field.
    * `op_value` - The logical symbol.
    * `sub_key` - The match subfield.
    * `values` - The match content. Separate multiple values with commas (,).
  * `rate_limit` - The rules of rate limiting.
    * `target` - The statistical field for frequency control.
    * `interval` - The statistical interval. Valid values: 5 to 1800. Unit: seconds.
    * `sub_key` - The subfield of the target field. 
    * `threshold` - The trigger threshold of rate limiting. Valid values: 2 to 500000. Unit: requests.
    * `ttl` - The validity period of the blacklist. Valid values: 60 to 86400. Unit: seconds.
    * `status` - The information about the HTTP status code.
      * `codet` - The HTTP status code returned.
      * `ratio` - The percentage of HTTP status codes.
      * `count` - The number of times that the HTTP status code that was returned.
  * `regular_types` - The type of the regular expression. If the value of the tags field contains waf_group, you can specify this field.
  * `remote_addr` - Filter by IP address.
  * `waf_group_ids` - The id of the waf rule group.
  * `effect` - The effective range of the frequency control blacklist.
  * `other_region_list` - The effective range of the frequency control blacklist.
  * `cc_status` - Whether to turn on Frequency Control, on/off
  * `scenes` - List of protection scenarios
  * `regular_rules` - The regular expression.

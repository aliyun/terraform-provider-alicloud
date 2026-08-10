---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_bucket_lifecycles"
description: |-
  Provides a list of ENS Bucket Lifecycle rules to the user.
---

# alicloud\_ens\_bucket\_lifecycles

This data source provides the ENS Bucket Lifecycle rules of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.244.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ens_bucket_lifecycles" "default" {
  bucket_name = "your-bucket-name"
}

output "lifecycle_rule_id" {
  value = data.alicloud_ens_bucket_lifecycles.default.rules.0.id
}
```

## Argument Reference

The following arguments are supported:

* `bucket_name` - (Required) The name of the ENS bucket whose lifecycle rules are to be listed.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `rule_id` - (Optional) A specific rule ID used to filter results to a single rule.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Bucket Lifecycle IDs, formatted as `<bucket_name>:<rule_id>`.
* `rules` - A list of ENS Bucket Lifecycle rules. Each element contains the following attributes:
  * `id` - The ID of the lifecycle rule, formatted as `<bucket_name>:<rule_id>`.
  * `bucket_name` - The name of the ENS bucket.
  * `rule_id` - The unique ID of the rule.
  * `prefix` - The prefix that the rule applies to.
  * `status` - The rule status. Valid values: `Enabled`, `Disabled`.
  * `expiration_days` - The number of days after the last update of the object before the lifecycle rule takes effect.
  * `created_before_date` - The expiration date in ISO8601 format.

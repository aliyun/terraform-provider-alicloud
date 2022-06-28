---
subcategory: "Quotas"
layout: "alicloud"
page_title: "Alicloud: alicloud_quotas_quotas"
sidebar_current: "docs-alicloud-datasource-quotas-quotas"
description: |-
  Provides a list of Quotas Quotas to the user.
---

# alicloud\_quotas\_quotas

This data source provides the Quotas Quotas of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.115.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_quotas_quotas" "example" {
  product_code = "ecs"
  name_regex   = "专有宿主机总数量上限"
}

output "first_quotas_quota_id" {
  value = data.alicloud_quotas_quotas.example.quotas.0.id
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional, ForceNew) A regex string to filter results by Quota name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `product_code` - (Required, ForceNew) The product code.
* `quota_action_code` - (Optional, ForceNew) The quota action code.
* `group_code` - (Optional, ForceNew) The group code.
* `key_word` - (Optional, ForceNew) The key word.
* `quota_category` - (Optional, ForceNew) The category of quota. Valid Values: `FlowControl` and `CommonQuota`.
* `sort_field` - (Optional, ForceNew) Cloud service ECS specification quota supports setting sorting fields. Valid Values: `TIME`, `TOTAL` and `RESERVED`.
* `sort_order` - (Optional, ForceNew) Ranking of cloud service ECS specification quota support. Valid Values: `Ascending` and `Descending`.
* `dimensions` - (Optional, ForceNew) The dimensions.

#### Block dimensions

The dimensions supports the following:

* `key` - The key of dimensions.
* `value` - The value of dimensions.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Quota IDs.
* `names` - A list of Quota names.
* `quotas` - A list of Quotas Quotas. Each element contains the following attributes:
	* `adjustable` - Is the quota adjustable.
	* `applicable_range` - The range of quota adjustment.
	* `applicable_type` - The type of quota.
	* `consumable` - Show used quota.
	* `id` - The ID of the Quota.
	* `quota_action_code` - The quota action code.
	* `quota_description` - The quota description.
	* `quota_name` - The quota name.
	* `quota_type` - The quota type.
	* `quota_unit` - The quota unit.
	* `total_quota` - TotalQuota.
	* `total_usage` - The total of usage.
	* `unadjustable_detail` - The unadjustable detail.

---
subcategory: "Security Center"
layout: "alicloud"
page_title: "Alicloud: alicloud_security_center_grouped_instances"
sidebar_current: "docs-alicloud-datasource-security-center-grouped-instances"
description: |-
  Provides the statistics of the assets in Security Center, grouped by the specified dimension.
---

# alicloud\_security\_center\_grouped\_instances

This data source provides the statistics of the assets in Security Center, grouped by the specified dimension.

-> **NOTE:** Available since v2.2.0.

## Example Usage

```terraform
data "alicloud_security_center_grouped_instances" "default" {
  group_field   = "regionId"
  machine_types = "ecs"
  page_size     = 10
  current_page  = 1
}

output "security_center_grouped_instances" {
  value = data.alicloud_security_center_grouped_instances.default.instances
}
```

## Argument Reference

The following arguments are supported:

* `group_field` - (Required, String) The dimension by which the assets are grouped. Valid values: `groupId`, `regionId` and `vpcInstanceId`.
* `machine_types` - (Optional, String) The type of the assets to query. The value is fixed to `ecs`.
* `field_value` - (Optional, String) The name of the asset group to query. Fuzzy match is supported.
* `vendor` - (Optional, Int) The vendor of the server. Valid values: `0` (Alibaba Cloud), `1` (outside Alibaba Cloud), `2` (IDC), `3`, `4`, `5`, `7` (other clouds) and `8` (lightweight server).
* `vendors` - (Optional, String) The vendors of the servers. Separate multiple values with commas (,). Valid values are the same as those of `vendor`.
* `sale_version_check_code` - (Optional, String) The sale condition used to filter the assets. Valid values: `sas_gte_advanced`, `sas_gte_enterprise`, `sas_gt_basic`, `sas_eq_advanced` and `sas_gt_anti_virus`.
* `lang` - (Optional, String) The language of the request and the response. Defaults to `zh`. Valid values: `zh` and `en`.
* `resource_directory_account_id` - (Optional, Int) The ID of the Alibaba Cloud account that is a member of the resource directory.
* `no_page` - (Optional, Bool) Specifies whether to disable the paged query. When it is set to `true`, all the results are returned in one response and `page_size` and `current_page` do not take effect. The query is paged when it is not set or set to `false`.
* `page_size` - (Optional, Computed, Int) The number of entries per page. Defaults to `20`.
* `current_page` - (Optional, Computed, Int) The number of the page to return. Defaults to `1`.

## Attributes Reference

The following attributes are exported:

* `id` - (String) The ID of the data source, which is derived from the query condition.
* `total_count` - (Int) The total number of the assets that match the query condition. This attribute is left unset when the API does not return it, for example when no asset matches the query condition.
* `page_item_count` - (Int) The number of the assets on the current page.
* `instances` - (List) A list of the asset statistics. Each element contains the following attributes. An attribute the API does not return for a group is left unset:
  * `asap_vul_instance_count` - (Int) The number of the assets in which high-risk vulnerabilities exist.
  * `unprotected_instance_count` - (String) The number of the assets that are not protected by Security Center.
  * `field_alias_name` - (String) The name of the server group.
  * `instance_count` - (String) The number of all the assets of the specified asset type.
  * `risk_instance_count` - (String) The number of the assets in which security risks exist.
  * `group_flag` - (Int) The type of the server group. `0` indicates the default group. `1` indicates another group. This attribute is returned when `group_field` is set to `groupId`.
  * `auth_version_check_count` - (Int) The number of the assets that meet the sale conditions.
  * `os` - (String) The operating system of the assets. Valid values: `windows` and `linux`. This attribute is returned only when `lang` is `zh`, and is left unset when the group has no operating system information.
  * `instance_core_count` - (Int) The number of the compute cores of the assets of the specified asset type.

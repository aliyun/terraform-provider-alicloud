---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_traffic_mirror_filters"
sidebar_current: "docs-alicloud-datasource-vpc-traffic-mirror-filters"
description: |-
  Provides a list of Vpc Traffic Mirror Filters to the user.
---

# alicloud\_vpc\_traffic\_mirror\_filters

This data source provides the Vpc Traffic Mirror Filters of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.140.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_vpc_traffic_mirror_filters" "ids" {
  ids = ["example_id"]
}
output "vpc_traffic_mirror_filter_id_1" {
  value = data.alicloud_vpc_traffic_mirror_filters.ids.filters.0.id
}

data "alicloud_vpc_traffic_mirror_filters" "nameRegex" {
  name_regex = "^my-TrafficMirrorFilter"
}
output "vpc_traffic_mirror_filter_id_2" {
  value = data.alicloud_vpc_traffic_mirror_filters.nameRegex.filters.0.id
}

data "alicloud_vpc_traffic_mirror_filters" "filterName" {
  traffic_mirror_filter_name = "example_traffic_mirror_filter_name"
}
output "vpc_traffic_mirror_filter_id_3" {
  value = data.alicloud_vpc_traffic_mirror_filters.filterName.filters.0.id
}

data "alicloud_vpc_traffic_mirror_filters" "status" {
  status = "^my-TrafficMirrorFilter"
}
output "vpc_traffic_mirror_filter_id_4" {
  value = data.alicloud_vpc_traffic_mirror_filters.status.filters.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Traffic Mirror Filter IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Traffic Mirror Filter name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The state of the filter. Valid values:`Creating`, `Created`, `Modifying` and `Deleting`. `Creating`: The filter is being created. `Created`: The filter is created. `Modifying`: The filter is being modified. `Deleting`: The filter is being deleted.
* `traffic_mirror_filter_ids` - (Optional, ForceNew) The traffic mirror filter ids.
* `traffic_mirror_filter_name` - (Optional, ForceNew) The name of the filter. The name must be `2` to `128` characters in length, and can contain digits, periods (.), underscores (_), and hyphens (-). It must start with a letter and cannot start with `http://` or `https://`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Traffic Mirror Filter names.
* `filters` - A list of Vpc Traffic Mirror Filters. Each element contains the following attributes:
	* `id` - The ID of the Traffic Mirror Filter.
	* `status` - The state of the filter. Valid values:`Creating`, `Created`, `Modifying` and `Deleting`. `Creating`: The filter is being created. `Created`: The filter is created. `Modifying`: The filter is being modified. `Deleting`: The filter is being deleted.
	* `traffic_mirror_filter_description` - The description of the filter.
	* `traffic_mirror_filter_id` - The ID of the filter.
	* `traffic_mirror_filter_name` - The name of the filter.

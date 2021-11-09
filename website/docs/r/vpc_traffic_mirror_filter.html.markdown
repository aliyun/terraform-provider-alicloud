---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_traffic_mirror_filter"
sidebar_current: "docs-alicloud-resource-vpc-traffic-mirror-filter"
description: |-
  Provides a Alicloud VPC Traffic Mirror Filter resource.
---

# alicloud\_vpc\_traffic\_mirror\_filter

Provides a VPC Traffic Mirror Filter resource.

For information about VPC Traffic Mirror Filter and how to use it, see [What is Traffic Mirror Filter](https://www.alibabacloud.com/help/doc-detail/207513.htm).

-> **NOTE:** Available in v1.140.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_vpc_traffic_mirror_filter" "example" {
  traffic_mirror_filter_name = "example_value"
}

```

## Argument Reference

The following arguments are supported:

* `dry_run` - (Optional) The dry run.
* `traffic_mirror_filter_description` - (Optional) The description of the filter. The description must be `2` to `256` characters in length. It must start with a letter and cannot start with `http://` or `https://`.
* `traffic_mirror_filter_name` - (Optional) The name of the filter. The name must be `2` to `128` characters in length, and can contain digits, periods (.), underscores (_), and hyphens (-). It must start with a letter and cannot start with `http://` or `https://`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Traffic Mirror Filter.
* `status` - The state of the filter. Valid values:`Creating`, `Created`, `Modifying` and `Deleting`. `Creating`: The filter is being created. `Created`: The filter is created. `Modifying`: The filter is being modified. `Deleting`: The filter is being deleted.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Traffic Mirror Filter.
* `delete` - (Defaults to 5 mins) Used when delete the Traffic Mirror Filter.

## Import

VPC Traffic Mirror Filter can be imported using the id, e.g.

```
$ terraform import alicloud_vpc_traffic_mirror_filter.example <id>
```

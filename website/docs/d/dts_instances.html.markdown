---
subcategory: "Data Transmission Service (DTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dts_instances"
sidebar_current: "docs-alicloud-datasource-dts-instances"
description: |-
  Provides a list of Dts Instance owned by an Alibaba Cloud account.
---

# alicloud_dts_instances

This data source provides Dts Instance available to the user.[What is Instance](https://www.alibabacloud.com/help/en/data-transmission-service/latest/createdtsinstance)

-> **NOTE:** Available in 1.198.0+

## Example Usage

```terraform
data "alicloud_dts_instances" "default" {
  ids               = ["${alicloud_dts_instance.default.id}"]
  resource_group_id = "example_value"
}

output "alicloud_dts_instance_example_id" {
  value = data.alicloud_dts_instances.default.instances.0.id
}
```

## Argument Reference

The following arguments are supported:
* `resource_group_id` - (ForceNew,Optional) Resource Group ID
* `ids` - (Optional, ForceNew, Computed) A list of Instance IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by trail name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Instance IDs.
* `names` - A list of Instance names.
* `instances` - A list of Instance Entries. Each element contains the following attributes:
  * `id` - The ID of the instance.
  * `create_time` - Instance creation time
  * `destination_endpoint_engine_name` - The target database engine type.
  * `dts_instance_id` - The ID of the subscription instance.
  * `instance_class` - The type of the migration or synchronization instance.- The specifications of the migration instance: **xxlarge**, **xlarge**, **large**, **medium**, **small**.- The types of synchronization instances: **large**, **medium**, **small**, **micro**.
  * `payment_type` - The payment type of the resource.
  * `destination_region` - The target instance region. 
  * `resource_group_id` - Resource Group ID.
  * `source_endpoint_engine_name` - Source instance database engine type.
  * `source_region` - The source instance region.
  * `destination_region` - The destination instance region.
  * `status` - Instance status.
  * `tags` - The tag value corresponding to the tag key.
  * `type` - The instance type. Valid values: -**MIGRATION**: MIGRATION.-**SYNC**: synchronization.-**SUBSCRIBE**: SUBSCRIBE.

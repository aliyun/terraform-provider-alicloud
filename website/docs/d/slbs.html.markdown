---
subcategory: "Classic Load Balancer (CLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slbs"
sidebar_current: "docs-alicloud-datasource-slbs"
description: |-
    Provides a list of server load balancers to the user.
---

# alicloud\_slbs

-> **DEPRECATED:** This datasource has been renamed to [alicloud_slb_load_balancers](https://www.terraform.io/docs/providers/alicloud/d/slb_load_balancers) from version 1.123.1.

This data source provides the server load balancers of the current Alibaba Cloud user.

## Example Usage

```
resource "alicloud_slb" "default" {
  name = "sample_slb"
}

data "alicloud_slbs" "slbs_ds" {
  name_regex = "sample_slb"
}

output "first_slb_id" {
  value = data.alicloud_slbs.slbs_ds.slbs[0].id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of SLBs IDs.
* `name_regex` - (Optional) A regex string to filter results by SLB name.
* `master_availability_zone` - (Optional) Master availability zone of the SLBs.
* `slave_availability_zone` - (Optional) Slave availability zone of the SLBs.
* `network_type` - (Optional) Network type of the SLBs. Valid values: `vpc` and `classic`.
* `vpc_id` - (Optional) ID of the VPC linked to the SLBs.
* `vswitch_id` - (Optional) ID of the VSwitch linked to the SLBs.
* `address` - (Optional) Service address of the SLBs.
* `tags` - (Optional) A map of tags assigned to the SLB instances. The `tags` can have a maximum of 5 tag. It must be in the format:
  ```
  data "alicloud_slbs" "taggedInstances" {
    tags = {
      tagKey1 = "tagValue1",
      tagKey2 = "tagValue2"
    }
  }
  ```
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional, ForceNew, Available in 1.60.0+) The Id of resource group which SLB belongs.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of slb IDs.
* `names` - A list of slb names.
* `slbs` - A list of SLBs. Each element contains the following attributes:
  * `id` - ID of the SLB.
  * `region_id` - Region ID the SLB belongs to.
  * `master_availability_zone` - Master availability zone of the SLBs.
  * `slave_availability_zone` - Slave availability zone of the SLBs.
  * `status` - SLB current status. Possible values: `inactive`, `active` and `locked`.
  * `name` - SLB name.
  * `network_type` - Network type of the SLB. Possible values: `vpc` and `classic`.
  * `vpc_id` - ID of the VPC the SLB belongs to.
  * `vswitch_id` - ID of the VSwitch the SLB belongs to.
  * `address` - Service address of the SLB.
  * `internet` - SLB addressType: internet if `true`, intranet if `false`. Must be `false` when `network_type` is `vpc`.
  * `creation_time` - SLB creation time.
  * `tags` - A map of tags assigned to the SLB instance.

---
subcategory: "Elastic Block Storage(EBS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ebs_regions"
sidebar_current: "docs-alicloud-datasource-ebs-regions"
description: |-
  Provides a list of Ebs Regions to the user.
---

# alicloud\_ebs\_regions

This data source provides the Ebs Regions of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.187.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ebs_regions" "default" {
  region_id = "cn-hangzhou"
}

output "regions" {
  value = data.alicloud_ebs_regions.default.regions
}
```

## Argument Reference

The following arguments are supported:

* `region_id` - (Optional)  A list of Disk Replica Group IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `regions` - A list of Ebs Regions. Each element contains the following attributes:
    * `region_id` - The ID of the region.
    * `zones` - A list of Ebs Zones.
      * `zone_id` - The ID of the zone.

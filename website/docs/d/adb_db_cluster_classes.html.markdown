---
subcategory: "AnalyticDB for MySQL (ADB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_adb_db_cluster_classes"
sidebar_current: "docs-alicloud-datasource-adb-db-cluster-classes"
description: |-
    Provides a list of available cluster classes for ADB that can be used by an Alibaba Cloud account.
---

# alicloud\_adb\_db\_cluster\_classes

This data source provides available cluster classes for ADB that can be accessed by an Alibaba Cloud account within the region configured in the provider.

-> **NOTE:** Available in v1.122.0+.

## Example Usage

```
# Declare the data source
data "alicloud_adb_db_cluster_classes" "availabile_classes" {}

output "classes" {
  value = alicloud_adb_db_cluster_classes.availabile_classes.available_zone_list[0].classes[0]
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Optional) The zone ID of the resource.
* `payment_type` - (Optional) The payment type of the resource. Valid values are `PayAsYouGo` and `Subscription`. Default to `PayAsYouGo`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `available_zone_list` - A list of availability zones. Each element contains the following attributes:
  * `zone_id` - ID of the zone.
  * `classes` - The available classes of zone.
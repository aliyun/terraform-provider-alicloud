---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_node_classes"
sidebar_current: "docs-alicloud-datasource-polardb-node-classes"
description: |-
    Provides a list of PolarDB node classes info.
---

# alicloud\_polardb\_node\_classes

This data source provides the PolarDB node classes resource available info of Alibaba Cloud.

-> **NOTE:** Available in v1.81.0+

## Example Usage

```terraform
data "alicloud_zones" "resources" {
  available_resource_creation = "PolarDB"
}

data "alicloud_polardb_node_classes" "resources" {
  zone_id    = data.alicloud_zones.resources.zones.0.id
  pay_type   = "PostPaid"
  db_type    = "MySQL"
  db_version = "5.6"
}

output "polardb_node_classes" {
  value = data.alicloud_polardb_node_classes.resources.classes
}
```

## Argument Reference

The following arguments are supported:

* `pay_type` - (Required) Filter the results by charge type. Valid values: `PrePaid` and `PostPaid`.
* `db_type` - (Optional) Database type. Options are `MySQL`, `PostgreSQL`, `Oracle`. If db_type is set, db_version also needs to be set.
* `db_version` - (Optional) Database version required by the user. Value options can refer to the latest docs [detail info](https://www.alibabacloud.com/help/doc-detail/98169.htm) `DBVersion`. If db_version is set, db_type also needs to be set.
* `db_node_class` - (Optional) The PolarDB node class type by the user.
* `region_id` - (Optional) The Region to launch the PolarDB cluster.
* `zone_id` - (Optional) The Zone to launch the PolarDB cluster.
* `output_file` - (Optional) File name where to save data source results (after running `terraform apply`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `classes` - A list of PolarDB node classes. Each element contains the following attributes:
  * `zone_id` - The Zone to launch the PolarDB cluster.
  * `supported_engines` - A list of PolarDB node classes in the zone.
    * `engine` - In the zone, the database type supports classes in the following available_resources.
    * `available_resources` - A list of PolarDB node available classes.
      * `db_node_class` - PolarDB node available class.

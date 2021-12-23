---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_db_instance_engines"
sidebar_current: "docs-alicloud-datasource-db-instance-engines"
description: |-
    Provides a list of RDS instacne engines resource info.
---

# alicloud\_db\_instance\_engines

This data source provides the RDS instance engines resource available info of Alibaba Cloud.

-> **NOTE:** Available in v1.46.0+

## Example Usage

```tf
data "alicloud_db_instance_engines" "resources" {
  instance_charge_type = "PostPaid"
  engine               = "MySQL"
  engine_version       = "5.6"
  output_file          = "./engines.txt"
}

output "first_db_category" {
  value = "${data.alicloud_db_instance_engines.resources.instance_engines.0.category}"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Optional) The Zone to launch the DB instance.
* `instance_charge_type` - (Optional) Filter the results by charge type. Valid values: `PrePaid` and `PostPaid`. Default to `PostPaid`.
* `engine` - (Optional) Database type. Valid values: "MySQL", "SQLServer", "PostgreSQL", "PPAS", "MariaDB". If not set, it will match all of engines.
* `engine_version` - (Optional) Database version required by the user. Value options can refer to the latest docs [detail info](https://www.alibabacloud.com/help/doc-detail/26228.htm) `EngineVersion`.
* `category` - (Optional, Available in 1.134.0+) DB Instance category. the value like [`Basic`, `HighAvailability`, `Finance`, `AlwaysOn`], [detail info](https://www.alibabacloud.com/help/doc-detail/69795.htm).
* `db_instance_storage_type` - (Optional, Available in 1.134.0+) The DB instance storage space required by the user. Valid values: "cloud_ssd", "local_ssd", "cloud_essd", "cloud_essd2", "cloud_essd3".
* `multi_zone` - (Optional, Available in v1.48.0+) Whether to show multi available zone. Default false to not show multi availability zone.
* `output_file` - (Optional) File name where to save data source results (after running `terraform apply`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of engines.
* `instance_engines` - A list of Rds available resource. Each element contains the following attributes:
  * `zone_ids` - A list of Zone to launch the DB instance.
    * `id` - The Zone to launch the DB instance
    * `sub_zone_ids` - A list of sub zone ids which in the id - e.g If `id` is `cn-beijing-MAZ5(a,b)`, `sub_zone_ids` will be `["cn-beijing-a", "cn-beijing-b"]`.
  * `engine` - Database type.
  * `engine_version` - DB Instance version.
  * `category` - DB Instance category.
  

---
layout: "alicloud"
page_title: "Alicloud: alicloud_db_instance_classes"
sidebar_current: "docs-alicloud-datasource-db-instance-classes"
description: |-
    Provides a list of RDS instacne classes info.
---

# alicloud\_db\_instances\_classes

This data source provides the RDS instance classes resource available info of Alibaba Cloud.

-> **NOTE:** Available in v1.46.0+

## Example Usage

```tf
data "alicloud_db_instance_classes" "resources" {
  instance_charge_type = "PostPaid"
  engine               = "MySQL"
  engine_version       = "5.6"
  output_file          = "./classes.txt"
}

output "first_db_instance_class" {
  value = "${data.alicloud_db_instance_classes.resources.instance_classes.0.instance_class}"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Optional) The Zone to launch the DB instance.
* `instance_charge_type` - (Optional) Filter the results by charge type. Valid values: `PrePaid` and `PostPaid`. Default to `PostPaid`.
* `engine` - (Optional) Database type. Options are `MySQL`, `SQLServer`, `PostgreSQL` and `PPAS`. If no value is specified, all types are returned.
* `category` - (Optional) DB Instance category. the value like [`Basic`, `HighAvailability`, `Finance`], [detail info](https://www.alibabacloud.com/help/doc-detail/69795.htm).
* `engine_version` - (Optional) Database version required by the user. Value options can refer to the latest docs [detail info](https://www.alibabacloud.com/help/doc-detail/26228.htm) `EngineVersion`.
* `storage_type` - (Optional) The DB instance storage space required by the user. Valid values: `cloud_ssd` and `local_ssd`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform apply`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `instance_classes` - A list of Rds available resource. Each element contains the following attributes:
  * `zone_id` - The Zone to launch the DB instance.
  * `instance_class` - DB Instance available class.
  * `storage_range` - DB Instance available storage range.
    * `min` - DB Instance available storage min value.
    * `max` - DB Instance available storage max value.
    * `step` - DB Instance available storage increase step.
    
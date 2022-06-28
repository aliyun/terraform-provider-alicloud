---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_db_instance_classes"
sidebar_current: "docs-alicloud-datasource-db-instance-classes"
description: |-
    Provides a list of RDS instacne classes info.
---

# alicloud\_db\_instance\_classes

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
* `engine` - (Optional) Database type. Valid values:"MySQL", "SQLServer", "PostgreSQL", "PPAS", "MariaDB". If not set, it will match all of engines.
* `category` - (Optional) DB Instance category. the value like [`Basic`, `HighAvailability`, `Finance`, `AlwaysOn`], [detail info](https://www.alibabacloud.com/help/doc-detail/69795.htm).
* `engine_version` - (Optional) Database version required by the user. Value options can refer to the latest docs [detail info](https://www.alibabacloud.com/help/doc-detail/26228.htm) `EngineVersion`.
* `db_instance_class` - (Optional, Available in 1.51.0+) The DB instance class type by the user.
* `storage_type` - (Optional) It has been deprecated from version 1.134.0+ and using `db_instance_storage_type` instead.
* `db_instance_storage_type` - (Optional, Available in 1.134.0+) The DB instance storage space required by the user. Valid values: "cloud_ssd", "local_ssd", "cloud_essd", "cloud_essd2", "cloud_essd3".
* `commodity_code` - (Optional, Available in 1.166.0+) The commodity code of the instance. Valid values:
  * **bards**: The instance is a pay-as-you-go primary instance. This value is available on the China site (aliyun.com).
  * **rds**: The instance is a subscription primary instance. This value is available on the China site (aliyun.com).
  * **rords**: The instance is a pay-as-you-go read-only instance. This value is available on the China site (aliyun.com).
  * **rds_rordspre_public_cn**: The instance is a subscription read-only instance. This value is available on the China site (aliyun.com).
  * **bards_intl**: The instance is a pay-as-you-go primary instance. This value is available on the International site (alibabacloud.com).
  * **rds_intl**: The instance is a subscription primary instance. This value is available on the International site (alibabacloud.com).
  * **rords_intl**: The instance is a pay-as-you-go read-only instance. This value is available on the International site (alibabacloud.com).
  * **rds_rordspre_public_intl**: The instance is a subscription read-only instance. This value is available on the International site (alibabacloud.com).

-> **NOTE**: The field `db_instance_id` will be ignored when `commodity_code` is not a read-only type.
* `db_instance_id` - (Optional, Available in v1.166.0+) The ID of the instance.

-> **NOTE**: The field `db_instance_id` will be ignored when `commodity_code` is not a read-only type.
* `multi_zone` - (Optional, Available in v1.48.0+) Whether to show multi available zone. Default false to not show multi availability zone.
* `output_file` - (Optional) File name where to save data source results (after running `terraform apply`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - (Available in 1.60.0+) A list of Rds instance class codes.
* `instance_classes` - A list of Rds available resource. Each element contains the following attributes:
  * `zone_ids` - A list of Zone to launch the DB instance.
    * `id` - The Zone to launch the DB instance
    * `sub_zone_ids` - A list of sub zone ids which in the id - e.g If `id` is `cn-beijing-MAZ5(a,b)`, `sub_zone_ids` will be `["cn-beijing-a", "cn-beijing-b"]`.
  * `instance_class` - DB Instance available class.
  * `storage_range` - DB Instance available storage range.
    * `min` - DB Instance available storage min value.
    * `max` - DB Instance available storage max value.
    * `step` - DB Instance available storage increase step.
    

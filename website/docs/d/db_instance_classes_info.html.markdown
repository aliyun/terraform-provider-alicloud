---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_db_instance_classes_info"
sidebar_current: "docs-alicloud-datasource-db-instance-classes-info"
description: |-
    Provides a list of RDS instacne classes detailed info.
---

# alicloud\_db\_instance\_classes\_info

This data source operation to query the instance types that are available to specific instances of Alibaba Cloud.

-> **NOTE:** Available in v1.195.0+

## Example Usage

```tf
data "alicloud_db_instance_classes_info" "resources" {
  commodity_code = "bards"
  order_type     = "BUY"
  output_file    = "./classes.txt"
}

output "first_db_instance_class" {
  value = "${data.alicloud_db_instance_classes_info.resources.instance_classes_infos.0}"
}
```

## Argument Reference

The following arguments are supported:

* `commodity_code` - (Required) The commodity code of the instance. Valid values:
  * **bards**: The instance is a pay-as-you-go primary instance. This value is available on the China site (aliyun.com).
  * **rds**: The instance is a subscription primary instance. This value is available on the China site (aliyun.com).
  * **rords**: The instance is a pay-as-you-go read-only instance. This value is available on the China site (aliyun.com).
  * **rds_rordspre_public_cn**: The instance is a subscription read-only instance. This value is available on the China site (aliyun.com).
  * **bards_intl**: The instance is a pay-as-you-go primary instance. This value is available on the International site (alibabacloud.com).
  * **rds_intl**: The instance is a subscription primary instance. This value is available on the International site (alibabacloud.com).
  * **rords_intl**: The instance is a pay-as-you-go read-only instance. This value is available on the International site (alibabacloud.com).
  * **rds_rordspre_public_intl**: The instance is a subscription read-only instance. This value is available on the International site (alibabacloud.com).
* `orderType` - (Required) FThe type of order that you want to query. Valid values:
  * **BUY**: specifies the query orders that are used to purchase instances.
  * **UPGRADE**: specifies the query orders that are used to change the specifications of instances.
  * **RENEW**: specifies the query orders that are used to renew instances.
  * **CONVERT**: specifies the query orders that are used to change the billing methods of instances.
* `output_file` - (Optional) File name where to save data source results (after running `terraform apply`).

-> **NOTE**: The field `db_instance_id` will be ignored when `commodity_code` is not a read-only type.
* `db_instance_id` - (Optional) The ID of the primary instance.

-> **NOTE**: The field `region_id` if you are using an Alibaba Cloud account on the International site (alibabacloud.com), you must specify the RegionId parameter.
* `region_id` - (Optional) The region ID of the instance.


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - (Available in 1.60.0+) A list of Rds instance class codes.
* `instance_classes_infos` - A list of Rds available resource. Each element contains the following attributes:
  * `class_code` - The code of the instance type.
  * `class_group` - The instance family of the instance.
  * `cpu` - The number of cores that are supported by the instance type. Unit: cores.
  * `max_connections` - The maximum number of connections that are supported by the instance type. Unit: connections.
  * `max_iombps` - The maximum I/O bandwidth that is supported by the instance type. Unit: Mbit/s.
  * `max_iops` - The maximum input/output operations per second (IOPS) that is supported by the instance type. Unit: operations per second.
  * `memory_class` - The memory capacity that is supported by the instance type. Unit: GB.
  
  -> **NOTE**: <br>If you set the CommodityCode parameter to a value that indicates the pay-as-you-go billing method, the ReferencePrice parameter specifies the hourly fee that you must pay.<br>If you set the CommodityCode parameter to a value that indicates the subscription billing method, the ReferencePrice parameter specifies the monthly fee that you must pay.
  * `reference_price` - The fee that you must pay for the instance type. Unit: cent (USD).
  
  -> **NOTE**: If the architecture of the instance type is x86, an empty string is returned by default. If the architecture of the instance type is ARM, arm is returned.
  * `instruction_set_arch` - The architecture of the instance type.
    

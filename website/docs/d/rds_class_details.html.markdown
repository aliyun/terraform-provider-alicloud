---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_class_details"
description: |-
  Provide users with Rds specification details.
---

# alicloud_rds_class_details

This data source provides details of the Rds specifications of current Alibaba Cloud users.

For information on RDS class details and how to use it, please refer to [What is RDS class details](https://www.alibabacloud.com/help/zh/apsaradb-for-rds/latest/api-rds-2014-08-15-describeclassdetails).

-> **NOTE:** Available since v1.209.0+

## Example Usage

```terraform
data "alicloud_rds_class_details" "default" {
  commodity_code = "bards"
  class_code     = "mysql.n4.medium.2c"
  engine_version = "8.0"
  engine         = "MySQL"
}
```

## Argument Reference

The following arguments are supported:

* `commodity_code` - (Required, ForceNew) The commodity code of the instance. Valid values:
  * **bards**: The instance is a pay-as-you-go primary instance. This value is available on the China site (aliyun.com).
  * **rds**: The instance is a subscription primary instance. This value is available on the China site (aliyun.com).
  * **rords**: The instance is a pay-as-you-go read-only instance. This value is available on the China site (aliyun.com).
  * **rds_rordspre_public_cn**: The instance is a subscription read-only instance. This value is available on the China site (aliyun.com).
  * **bards_intl**: The instance is a pay-as-you-go primary instance. This value is available on the International site (alibabacloud.com).
  * **rds_intl**: The instance is a subscription primary instance. This value is available on the International site (alibabacloud.com).
  * **rords_intl**: The instance is a pay-as-you-go read-only instance. This value is available on the International site (alibabacloud.com).
  * **rds_rordspre_public_intl**: The instance is a subscription read-only instance. This value is available on the International site (alibabacloud.com).
* `class_code` - (Required, ForceNew) The code of the instance type.
* `engine_version` - (Required, ForceNew) Database version. Value options:
  - MySQL: [ 5.5、5.6、5.7、8.0 ]
  - SQLServer: [ 2008r2、08r2_ent_ha、2012、2012_ent_ha、2012_std_ha、2012_web、2014_std_ha、2016_ent_ha、2016_std_ha、2016_web、2017_std_ha、2017_ent、2019_std_ha、2019_ent ]
  - PostgreSQL: [ 10.0、11.0、12.0、13.0、14.0、15.0 ]
  - MariaDB: [ 10.3 ]
* `engine` - (Required, ForceNew) Database type. Value options: MySQL, SQLServer, PostgreSQL, MariaDB.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `max_iombps` - The maximum IO bandwidth corresponding to the instance specification. Unit: Mbps. 
* `max_connections` - The maximum number of connections.
* `class_group` - The specification family.
* `cpu` - The number of CPU cores corresponding to the instance specification. Unit: pieces.
* `instruction_set_arch` - The architecture of the instance type.
* `memory_class` - The memory capacity that is supported by the instance type. Unit: GB.
* `max_iops` - The maximum IOPS of the instance.
* `reference_price` - The fee that you must pay for the instance type. Unit: cent (RMB).
* `category` - 	The RDS edition of the instance. Valid values:
  * **Basic**: Basic Edition.
  * **HighAvailability**: High-availability Edition.
  * **AlwaysOn**: Cluster Edition.
  * **Finance**: Enterprise Edition.
* `db_instance_storage_type` - 	 The storage type of the instance. Valid values:
  * **local_ssd**: specifies to use local SSDs.
  * **cloud_ssd**: specifies to use standard SSDs.
  * **cloud_essd**: specifies to use enhanced SSDs (ESSDs).
  * **cloud_essd2**: specifies to use enhanced SSDs (ESSDs).
  * **cloud_essd3**: specifies to use enhanced SSDs (ESSDs).

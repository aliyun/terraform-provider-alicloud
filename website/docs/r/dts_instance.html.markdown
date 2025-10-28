---
subcategory: "Data Transmission Service (DTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dts_instance"
sidebar_current: "docs-alicloud-resource-dts-instance"
description: |-
  Provides a Alicloud Dts Instance resource.
---

# alicloud_dts_instance

Provides a Dts Instance resource.

For information about Dts Instance and how to use it, see [What is Instance](https://www.alibabacloud.com/help/en/data-transmission-service/latest/createdtsinstance).

-> **NOTE:** Available since v1.198.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_dts_instance&exampleId=d1f3d8bc-acd2-5ee3-ab9b-62388fe4a382368040c8&activeTab=example&spm=docs.r.dts_instance.0.d1f3d8bcac&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}
data "alicloud_regions" "default" {
  current = true
}
resource "alicloud_dts_instance" "default" {
  type                             = "sync"
  resource_group_id                = data.alicloud_resource_manager_resource_groups.default.ids.0
  payment_type                     = "Subscription"
  instance_class                   = "large"
  source_endpoint_engine_name      = "MySQL"
  source_region                    = data.alicloud_regions.default.regions.0.id
  destination_endpoint_engine_name = "MySQL"
  destination_region               = data.alicloud_regions.default.regions.0.id
}
```

## Argument Reference

The following arguments are supported:
* `auto_pay` - (Optional) Whether to automatically renew the fee when it expires. Valid values:
  - **false**: No, the default value.
  - **true**: Yes.
* `auto_start` - (Optional) Whether to start the task automatically after the purchase is completed. Value:
  - **false**: No, the default value.
  - **true**: Yes.
* `compute_unit` - (Optional) Specifications of ETL. The unit is compute unit (CU),1CU = 1vCPU +4GB of memory. The value range is an integer greater than or equal to 2. **NOTE:** Enter this parameter and enable ETL to clean and convert data. 
* `database_count` - (Optional) The number of private custom RDS instances in the PolarDB-X. The default value is **1**. **NOTE:** This parameter is required only when **source_endpoint_engine_name** is **DRDS**.
* `destination_endpoint_engine_name` - (Optional, ForceNew) The target database engine type.
  - **MySQL**:MySQL databases (including RDS MySQL and self-built MySQL).
  - **PolarDB**:PolarDB MySQL.
  - **polardb_o**:PolarDB O engine.
  - **polardb_pg**:PolarDB PostgreSQL.
  - **Redis**:Redis databases (including apsaradb for Redis and user-created Redis).
  - **DRDS**: cloud-native distributed database PolarDB-X 1.0 and 2.0.
  - **PostgreSQL**: User-created PostgreSQL.
  - **odps**: MaxCompute project.
  - **oracle**: self-built Oracle.
  - **mongodb**:MongoDB databases (including apsaradb for MongoDB and user-created MongoDB).
  - **tidb**:TiDB database.
  - **ADS**: Cloud native data warehouse AnalyticDB MySQL 2.0.
  - **ADB30**: Cloud native data warehouse AnalyticDB MySQL 3.0.
  - **Greenplum**: Cloud native data warehouse AnalyticDB PostgreSQL.
  - **MSSQL**:SQL Server databases (including RDS SQL Server and self-built SQL Server).
  - **kafka**:Kafka databases (including Kafka and self-built Kafka).
  - **DataHub**: DataHub, an Alibaba cloud streaming data service.
  - **clickhouse**: ClickHouse.
  - **DB2**: self-built DB2 LUW.
  - **as400**:AS/400.
  - **Tablestore**: Tablestore.
  - **NOTE:** 
    - The default value is **MySQL**.
    - For more information about the supported source and destination databases, see [Database, Synchronization Initialization Type, and Synchronization Topology](https://www.alibabacloud.com/help/en/data-transmission-service/latest/overview-of-data-synchronization-scenarios-1) and [Supported Database and Migration Type](https://www.alibabacloud.com/help/en/data-transmission-service/latest/overview-of-data-migration-scenarios).
    - This parameter or **job_id** must be passed in.
* `du` - (Optional) Assign a specified number of DU resources to DTS tasks in the DTS exclusive cluster. Valid values: **1** ~ **100**. **NOTE:** The value of this parameter must be within the range of the number of DUs available for the DTS dedicated cluster.
* `fee_type` - (Optional) Subscription billing type, Valid values: `ONLY_CONFIGURATION_FEE`: charges only configuration fees; `CONFIGURATION_FEE_AND_DATA_FEE`: charges configuration fees and data traffic fees.
* `instance_class` - (Optional, ForceNew) The type of the migration or synchronization instance.
  - The specifications of the migration instance: **xxlarge**, **xlarge**, **large**, **medium**, **small**. 
  - The types of synchronization instances: **large**, **medium**, **small**, **micro**. 
  - **NOTE:** For performance descriptions of different specifications, see [Data Migration Link Specifications](https://www.alibabacloud.com/help/en/data-transmission-service/latest/cd773b) and [Data Synchronization Link Specifications](https://www.alibabacloud.com/help/en/data-transmission-service/latest/6bce7c).
* `job_id` - (Optional) The ID of the task obtained by calling the **ConfigureDtsJob** operation (**DtsJobId**).> After you pass in this parameter, you do not need to pass the **source_region**, **destination_region**, **type**, **source_endpoint_engine_name**, or **destination_endpoint_engine_name** parameters. Even if the input is passed in, the configuration in **job_id** shall prevail.
* `payment_type` - (Optional, ForceNew) The payment type of the resource. Valid values: `Subscription`, `PayAsYouGo`.
* `period` - (Optional) The billing method of the subscription instance. Value: `Year`, `Month`. **NOTE:** This parameter is valid and must be passed in only when `payment_type` is `Subscription`.
* `destination_region` - (Optional, ForceNew) The target instance region. For more information, see [List of supported regions](https://www.alibabacloud.com/help/en/data-transmission-service/latest/list-of-supported-regions). **NOTE:** This parameter or **job_id** must be passed in.
* `resource_group_id` - (Optional) Resource Group ID.
* `source_endpoint_engine_name` - (ForceNew, Optional) Source instance database engine type.
  - **MySQL**:MySQL databases (including RDS MySQL and self-built MySQL).
  - **PolarDB**:PolarDB MySQL.
  - **polardb_o**:PolarDB O engine.
  - **polardb_pg**:PolarDB PostgreSQL.
  - **Redis**:Redis databases (including apsaradb for Redis and user-created Redis).
  - **DRDS**: cloud-native distributed database PolarDB-X 1.0 and 2.0.
  - **PostgreSQL**: User-created PostgreSQL.
  - **odps**: MaxCompute project.
  - **oracle**: self-built Oracle.
  - **mongodb**:MongoDB databases (including apsaradb for MongoDB and user-created MongoDB).
  - **tidb**:TiDB database.
  - **ADS**: Cloud native data warehouse AnalyticDB MySQL 2.0.
  - **ADB30**: Cloud native data warehouse AnalyticDB MySQL 3.0.
  - **Greenplum**: Cloud native data warehouse AnalyticDB PostgreSQL.
  - **MSSQL**:SQL Server databases (including RDS SQL Server and self-built SQL Server).
  - **kafka**:Kafka databases (including Kafka and self-built Kafka).
  - **DataHub**: DataHub, an Alibaba cloud streaming data service.
  - **clickhouse**: ClickHouse.
  - **DB2**: self-built DB2 LUW.
  - **as400**:AS/400.
  - **Tablestore**: Tablestore.
  - **NOTE:**
    - The default value is **MySQL**.
    - For more information about the supported source and destination databases, see [Database, Synchronization Initialization Type, and Synchronization Topology](https://www.alibabacloud.com/help/en/data-transmission-service/latest/overview-of-data-synchronization-scenarios-1) and [Supported Database and Migration Type](https://www.alibabacloud.com/help/en/data-transmission-service/latest/overview-of-data-migration-scenarios).
    - This parameter or **job_id** must be passed in.
* `source_region` - (Optional, ForceNew) The source instance region. For more information, see [List of supported regions](https://www.alibabacloud.com/help/en/data-transmission-service/latest/list-of-supported-regions). **NOTE:** This parameter or **job_id** must be passed in.
* `sync_architecture` - (Optional) Synchronization topology, value:
  - **oneway**: one-way synchronization, the default value.
  - **bidirectional**: two-way synchronization.
* `tags` - (Optional) The tag value corresponding to the tag key.See the following `Block Tags`.
* `type` - (Optional, ForceNew) The instance type. Valid values:
  - **migration**: MIGRATION.
  - **sync**: synchronization.
  - **subscribe**: SUBSCRIBE.
  - **NOTE:** This parameter or **job_id** must be passed in.
* `used_time` - (Optional) Prepaid instance purchase duration.
  - When **period** is **Month**, the values are: 1, 2, 3, 4, 5, 6, 7, 8, and 9.
  - When **Period** is **Year**, the values are 1, 2, 3, and 5.
  - **NOTE:** 
    - This parameter is valid and must be passed in only when **payment_type** is `Subscription`.
    - The billing method of the subscription instance. You can set the parameter `period`.
* `synchronization_direction` - (Optional) The synchronization direction. Default value: `Forward`. Valid values:
  - `Forward`: Data is synchronized from the source database to the destination database.
  - `Reverse`: Data is synchronized from the destination database to the source database.
  - **NOTE:** You can set this parameter to Reverse to delete the reverse synchronization task only if the topology is two-way synchronization.


## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `create_time` - Instance creation time
* `dts_instance_id` - The ID of the subscription instance.
* `instance_name` - The name of Dts instance.
* `status` - Instance status.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Instance.
* `delete` - (Defaults to 5 mins) Used when delete the Instance.
* `update` - (Defaults to 5 mins) Used when update the Instance.

## Import

Dts Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_dts_instance.example <id>
```
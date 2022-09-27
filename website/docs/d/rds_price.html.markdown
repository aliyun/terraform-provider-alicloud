---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_price"
sidebar_current: "docs-alicloud-datasource-rds-price"
description: |-
  Provides a list of Rds Price to the user.
---

# alicloud\_rds\_price

This data source provides the Rds Price of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.175.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_rds_price" "example" {
  engine_version           = "13.0"
  db_instance_class        = "pg.n2.2c.2m"
  db_instance_storage      = "20"
  quantity                 = "10"
  engine                   = "PostgreSQL"
  commodity_code           = "bards"
  pay_type                 = "Prepaid"
  used_time                = "1"
  time_type                = "Year"
  instance_used_type       = "0"
  order_type               = "BUY"
  db_instance_storage_type = "cloud_essd"
  output_file              = "./price.txt"
}
```

## Argument Reference

The following arguments are supported:

* `commodity_code` - (Optional, ForceNew) The commodity code of the instances. Valid values:
  * **bards**: The instances are pay-as-you-go primary instances. This value is available at the China site (aliyun.com).
  * **rds**: The instances are subscription primary instances. This is the default value. This value is available at the China site (aliyun.com).
  * **rords**: The instances are pay-as-you-go read-only instances. This value is available at the China site (aliyun.com).
  * **rds_rordspre_public_cn**: The instances are subscription read-only instances. This value is available at the China site (aliyun.com).
  * **bards_intl**: The instances are pay-as-you-go primary instances. This value is available at the International site (alibabacloud.com).
  * **rds_intl**: The instances are subscription primary instances. This value is available at the International site (alibabacloud.com).
  * **rords_intl**: The instances are pay-as-you-go read-only instances. This value is available at the International site (alibabacloud.com).
  * **rds_rordspre_public_intl**: The instances are subscription read-only instances. This value is available at the International site (alibabacloud.com).
  
-> **NOTE:** If you want to query the price of read-only instances, you must specify this parameter.
* `engine` - (Required, ForceNew) The database engine that is run on the instances. Valid values: **MySQL**, **SQLServer**, **PostgreSQL** and **MariaDB TX**.
* `engine_version` - (Required, ForceNew) The database engine version that is run on the instances. Valid values:
  * MySQL: **5.5**, **5.6**, **5.7**, and **8.0**.
  * SQL Server: **2008r2**, **2012**, **2012_ent_ha**, **2012_std_ha**, **2012_web**, **2014_std_ha**, **2016_ent_ha**, **2016_std_ha**, **2016_web**, **2017_std_ha**, **2017_ent**, **2019_std_ha**, and **2019_ent**.
  * PostgreSQL: **10.0**, **11.0**, **12.0**, **13.0** and **14.0**.
  * MariaDB TX: **10.3**.
* `db_instance_class` - (Required, ForceNew) The instance type of the instances. For more information, see [Primary ApsaraDB RDS instance types](https://www.alibabacloud.com/help/en/apsaradb-for-rds/latest/primary-apsaradb-rds-instance-types).
* `db_instance_storage` - (Required, ForceNew) The storage capacity of the instances. Unit: GB. The storage capacity increases in increments of 5 GB. For more information, see [Primary ApsaraDB RDS instance types](https://www.alibabacloud.com/help/en/apsaradb-for-rds/latest/primary-apsaradb-rds-instance-types).
* `pay_type` - (Optional, ForceNew) The billing method of the instances. Valid values:
  * **Prepaid**: subscription.
  * **Postpaid**: pay-as-you-go.
* `zone_id` - (Optional, ForceNew) The zone ID of the primary instances. You can call the [DescribeRegions](https://www.alibabacloud.com/help/en/apsaradb-for-rds/latest/query-regions) operation to query the most recent zone list.
  
-> **NOTE:** If you specify a VPC and a vSwitch, you must specify this parameter.
* `used_time` - (Optional, ForceNew) The subscription period of the instances. Default value: **1**. Valid values:
  * If you set the **time_type** parameter to **Year**, the value of the used_time parameter ranges from **1 to 100**.
  * If you set the **time_type** parameter to **Month**, the value of the used_time parameter ranges from **1 to 999**.
* `time_type` - (Optional, ForceNew) The unit that is used to calculate the billing cycle of the instances. If the value of the **commodity_code** parameter is **rds**, **rds_rordspre_public_cn**, **rds_intl**, or **rds_rordspre_public_intl**, you must specify this parameter. Valid values: **Year** and **Month**.
* `quantity` - (Required, ForceNew) The number of instances that you want to purchase. Valid values: **0 to 30**.
* `instance_used_type` - (Optional, ForceNew) The role of the instances. Valid values:
  * **0**: primary.
  * **3**: read-only. 
* `order_type` - (Optional, ForceNew) The type of the order. Valid values:
  * **BUY**: purchase order.
  * **UPGRADE**: specification change order.
  * **RENEW**: renewal order.
* `db_instance_storage_type` - (Optional, ForceNew) The type of storage media that is used for the instances. Valid values:
  * **local_ssd**: local SSDs.
  * **cloud_ssd**: standard SSDs.
  * **cloud_essd**: enhanced SSDs (ESSDs) of performance level 1 (PL1).
  * **cloud_essd2**: ESSDs of PL2
  * **cloud_essd3**: ESSDs of PL3
* `db_instance_id` - (Optional, ForceNew) The IDs of the instances whose specifications you want to change or those that you want to renew.
  
-> **NOTE:** If you want to query the price that is required to change the specifications of specific instances or renew specific instances, you must specify this parameter.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `price` - A list of Rds Price. Each element contains the following attributes:
  * `price_info` - The price that is specified in the promotion rule.
    * `currency` - The currency that is used to measure the price.
    * `original_price` - The original price that you need to pay.
    * `discount_price` - The discount that is applied based on the promotion rule.
    * `trade_price` - The discounted price that you need to pay.
    * `rule_ids` - The ID of the promotion rule.
    * `coupons` - An array that consists of information about coupons.
      * `name` - The name of the coupon.
      * `description` - The description of the coupon.
      * `coupon_no` - The ID of the coupon.
      * `is_selected` - Indicates whether the coupon is selected.
  * `rules` - An array that consists of promotion rules.
    * `rule_id` - The ID of the promotion rule.
    * `name` - The name of the promotion rule.
    * `description` - The description of the promotion rule.
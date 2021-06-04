---
subcategory: "Cassandra"
layout: "alicloud"
page_title: "Alicloud: alicloud_cassandra_data_centers"
sidebar_current: "docs-alicloud-datasource-cassandra-data-centers"
description: |-
    Provides a collection of Cassandra Data Centers according to the specified filters.
---

# alicloud\_cassandra\_data\_centers

The `alicloud_cassandra_data_centers` data source provides a collection of Cassandra Data Centers available in Alicloud account.
Filters support regular expression for the cluster name or ids.

-> **NOTE:**  Available in 1.88.0+.

## Example Usage

```
data "alicloud_cassandra_data_centers" "cassandra" {
  name_regex        = "tf_testAccCassandra_dc"
  cluster_id        = "cds-xxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to apply to the cluster name.
* `ids` - (Optional) The list of Cassandra data center ids.
* `names` - (Optional) The name list of Cassandra data centers.
* `cluster_id` - (Required) The cluster id of dataCenters belongs to.
* `output_file` - (Optional) The name of file that can save the collection of data centers after running `terraform plan`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - The list of Cassandra data center ids.
* `names` - The name list of Cassandra data centers.
* `centers` - A list of Cassandra data centers. Its every element contains the following attributes:
  * `cluster_id` - The ID of the Cassandra cluster.
  * `commodity_instance` - The commodity ID of the Cassandra dataCenter.
  * `data_center_id` - The id of the Cassandra dataCenter.
  * `data_center_name` - The name of the Cassandra dataCenter.
  * `disk_size` - One node disk size, unit:GB.
  * `disk_type` - Cloud_ssd or cloud_efficiency.
  * `instance_type` - The instance type of the Cassandra dataCenter, eg: cassandra.c.large.
  * `lock_mode` - The lock mode of the dataCenter.
  * `node_count` - The node count of dataCenter.
  * `status` - Status of the dataCenter.
  * `create_time` - The create time of the dataCenter.
  * `expire_time` - The expire time of the dataCenter.
  * `zone_id` - Zone ID the dataCenter belongs to.
  * `vpc_id` - VPC ID the dataCenter belongs to.
  * `vswitch_id` - VSwitch ID the dataCenter belongs to.
  * `pay_type` - Billing method. Value options are `Subscription` for Pay-As-You-Go and `PayAsYouGo` for yearly or monthly subscription.


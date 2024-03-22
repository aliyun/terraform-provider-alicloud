---
subcategory: "Cassandra"
layout: "alicloud"
page_title: "Alicloud: alicloud_cassandra_clusters"
sidebar_current: "docs-alicloud-datasource-cassandra-clusters"
description: |-
    Provides a collection of Cassandra clusters according to the specified filters.
---

# alicloud\_cassandra\_clusters

The `alicloud_cassandra_clusters` data source provides a collection of Cassandra clusters available in Alicloud account.
Filters support regular expression for the cluster name, ids or tags.

-> **NOTE:**  Available in 1.88.0+.

-> **DEPRECATED:**  This data source has been [deprecated](https://www.alibabacloud.com/help/en/apsaradb-for-cassandra/latest/cassandra-delisting-notice) from version `1.219.1`.

## Example Usage

```
data "alicloud_cassandra_clusters" "cassandra" {
  name_regex        = "tf_testAccCassandra"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to apply to the cluster name.
* `ids` - (Optional) The list of Cassandra cluster ids.
* `tags` - (Optional, Available in 1.73.0) A mapping of tags to assign to the resource.
* `output_file` - (Optional) The name of file that can save the collection of clusters after running `terraform plan`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - The list of Cassandra cluster ids.
* `names` - The name list of Cassandra clusters.
* `clusters` - A list of Cassandra clusters. Its every element contains the following attributes:
  * `id` - The ID of the Cassandra cluster.
  * `cluster_id` - The ID of the Cassandra cluster.
  * `cluster_name` - The name of the Cassandra cluster.
  * `major_version` - The major version of the cluster.
  * `minor_version` - The minor version of the cluster.
  * `lock_mode` - The lock mode of the cluster.
  * `data_center_count` - The count of data centers
  * `pay_type` - Billing method. Value options are `Subscription` for Pay-As-You-Go and `PayAsYouGo` for yearly or monthly subscription.
  * `status` - Status of the cluster.
  * `create_time` - The create time of the cluster.
  * `expire_time` - The expire time of the cluster.
  * `tags` - A mapping of tags to assign to the resource.

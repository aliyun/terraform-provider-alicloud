---
layout: "alicloud"
page_title: "Alicloud: alicloud_elasticsearch"
sidebar_current: "docs-alicloud-datasource-elasticsearch"
description: |-
    Provides a collection of Elasticsearch instances according to the specified filters.
---

# alicloud\_elasticsearch

The `alicloud_elasticsearch` data source provides a collection of Elasticsearch instances available in Alicloud account.
Filters support description and other filters which are listed below.

## Example Usage

```
data "alicloud_elasticsearch" "instances" {
  description = "myes"
  es_version  = "5.5.3_with_X-Pack"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The instance's description. Fuzzy matching is supported.
* `es_version` - (Optional) Elasticsearch version. Options are `5.5.3_with_X-Pack`, and `6.3.2_with_X-Pack`. If no value is specified, all versions are returned.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `instances` - A list of Elasticsearch instances. Its every element contains the following attributes:
  * `id` - The ID of the Elasticsearch instance.
  * `description` - The description of the Elasticsearch instance.
  * `instance_charge_type` - Billing method. Value options: `PostPaid` for  Pay-As-You-Go and `PrePaid` for subscription.
  * `data_node_amount` - The Elasticsearch cluster's data node quantity, between 2 and 50.
  * `data_node_spec` - The data node specifications of the elasticsearch instance.
  * `data_node_disk` - The single data node storage space. Unit: GB.
  * `data_node_disk_type` - The data node disk type. Supported values: cloud_ssd, cloud_efficiency.
  * `vswitch_id` - VSwitch ID the instance belongs to.
  * `es_version` - Elasticsearch version.
  * `cerated_at` - The creation time of the instance. It's a GTM format, such as: "2019-01-08T15:50:50.623Z".
  * `updated_at` - The last modified time of the instance. It's a GMT format, such as: "2019-01-08T15:50:50.623Z".
  * `status` - Status of the instance. It includes `active`, `activating`, `inactive`

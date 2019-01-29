---
layout: "alicloud"
page_title: "Alicloud: alicloud_elasticsearch_instances"
sidebar_current: "docs-alicloud-datasource-elasticsearch-instances"
description: |-
  Provides a collection of Elasticsearch instances according to the specified filters.
---

# alicloud\_elasticsearch\_instances

The `alicloud_elasticsearch_instances` data source provides a collection of Elasticsearch instances available in Alicloud account.
Filters support description regex and other filters which are listed below.

## Example Usage

```
data "alicloud_elasticsearch_instances" "instances" {
  description_regex = "myes"
  version           = "5.5.3_with_X-Pack"
}
```

## Argument Reference

The following arguments are supported:

* `description_regex` - (Optional) A regex string to apply to the instance description.
* `version` - (Optional) Elasticsearch version. Options are `5.5.3_with_X-Pack`, and `6.3.2_with_X-Pack`. If no value is specified, all versions are returned.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Elasticsearch instance IDs.
* `instances` - A list of Elasticsearch instances. Its every element contains the following attributes:
  * `id` - The ID of the Elasticsearch instance.
  * `description` - The description of the Elasticsearch instance.
  * `instance_charge_type` - Billing method. Value options: `PostPaid` for  Pay-As-You-Go and `PrePaid` for subscription.
  * `data_node_amount` - The Elasticsearch cluster's data node quantity, between 2 and 50.
  * `data_node_spec` - The data node specifications of the elasticsearch instance.
  * `data_node_disk_size` - The single data node storage space. Unit: GB.
  * `data_node_disk_type` - The data node disk type. Included values: `cloud_ssd` and `cloud_efficiency`.
  * `vswitch_id` - VSwitch ID the instance belongs to.
  * `version` - Elasticsearch version includes `5.5.3_with_X-Pack` and `6.3.2_with_X-Pack`.
  * `cerated_at` - The creation time of the instance. It's a GTM format, such as: "2019-01-08T15:50:50.623Z".
  * `updated_at` - The last modified time of the instance. It's a GMT format, such as: "2019-01-08T15:50:50.623Z".
  * `status` - Status of the instance. It includes `active`, `activating`, `inactive`

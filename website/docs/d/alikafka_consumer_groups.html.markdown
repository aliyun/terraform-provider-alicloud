---
subcategory: "Alikafka"
layout: "alicloud"
page_title: "Alicloud: alicloud_alikafka_consumer_groups"
sidebar_current: "docs-alicloud-datasource-alikafka-consumer-groups"
description: |-
    Provides a list of alikafka consumer groups available to the user.
---

# alicloud\_alikakfa\_consumer\_groups

This data source provides a list of ALIKAFKA Consumer Groups in an Alibaba Cloud account according to the specified filters.

-> **NOTE:** Available in 1.56.0+

## Example Usage

```
data "alicloud_alikafka_consumer_groups" "consumer_groups_ds" {
  instance_id = "xxx"
  consumer_id_regex = "CID-alikafkaGroupDatasourceName"
  output_file = "consumerGroups.txt"
}

output "first_group_name" {
  value = "${data.alicloud_alikafka_consumer_groups.consumer_groups_ds.consumer_ids.0}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of ALIKAFKA Consumer Groups IDs, It is formatted to `<instance_id>:<consumer_id>`.
* `instance_id` - (Required) ID of the ALIKAFKA Instance that owns the consumer groups.
* `consumer_id_regex` - (Optional) A regex string to filter results by the consumer group id. 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of consumer group names.
* `groups` - A list of consumer group. Each element contains the following attributes:
    * `id` - The ID of the consumer group, It is formatted to `<instance_id>:<consumer_id>`.
    * `consumer_id` - The name of the consumer group.
    * `remark` - The remark of the consumer group.
    * `instance_id` - The instance_id of the instance.
    * `tags` - A mapping of tags to assign to the consumer group.

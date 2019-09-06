---
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
  value = "${data.alicloud_alikafka_consumer_groups.consumer_groups_ds.consumer_groups.0.consumer_id}"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) ID of the ALIKAFKA Instance that owns the consumer groups.
* `consumer_id_regex` - (Optional) A regex string to filter results by the consumer group id. 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of consumer group ids.
* `consumer_groups` - A list of consumer groups. Each element contains the following attributes:
  * `consumer_id` - The consumer id of the group.
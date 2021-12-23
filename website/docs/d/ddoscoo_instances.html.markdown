---
subcategory: "Anti-DDoS Pro"
layout: "alicloud"
page_title: "Alicloud: alicloud_ddoscoo_instances"
sidebar_current: "docs-alicloud-datasource-ddoscoo-instances"
description: |-
  Provides a list of BGP-Line Anti-DDoS Pro instances available to the user.
---

# alicloud\_ddoscoo\_instances

This data source provides a list of BGP-Line Anti-DDoS Pro instances in an Alibaba Cloud account according to the specified filters.

## Example Usage

```
data "alicloud_ddoscoo_instances" "instance" {
  name_regex = "^ddoscoo"
}

output "instance" {
  value = "${alicloud_ddoscoo_instances.instance.*.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter results by the instance name.
* `ids` - (Optional) A list of instance IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of instance IDs.
* `names` - A list of instance names.
* `instances` - A list of apis. Each element contains the following attributes:
  * `id` - The instance's id.
  * `name` - The instance's remark.
  * `base_bandwidth` - The instance's base defend bandwidth.
  * `bandwidth` - The instance's elastic defend bandwidth.
  * `service_bandwidth` - The instance's business bandwidth.
  * `port_count` - The instance's count of port retransmission config.
  * `domain_count` - The instance's count of domain retransmission config.
  * `remark` - The remark of the instance.
  * `ip_mode` - The ip mode of the instance. The Valid Values : `fnat`, `v6tov4`.
  * `debt_status` - The debt status of the instance.
  * `edition` - The edition of the instance. The Valid Values : `0`, `1`, `2`, `9`.
  * `ip_version` - The ip version of the instance. The Valid Values : `Ipv4`, `Ipv6`.
  * `status` - The status of the instance. The Valid Values : `1`, `2`.
  * `enabled` - The enabled of the instance. The Valid Values : `0`, `1`.
  * `expire_time` - The expiry time of the instance.
  * `create_time` - The creation time of the instance.
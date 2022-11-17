---
subcategory: "Anti-DDoS Pro (DdosBgp)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ddoscoo_instances"
sidebar_current: "docs-alicloud-datasource-ddoscoo-instances"
description: |-
  Provides a list of Anti-DDoS Advanced(Ddosbgp) instances available to the user.
---

# alicloud\_ddosbgp\_instances

This data source provides a list of Anti-DDoS Advanced instances in an Alibaba Cloud account according to the specified filters.

-> **NOTE:** Available in 1.183.0+ .

## Example Usage

```
data "alicloud_ddosbgp_instances" "instance" {
  name_regex = "^ddosbgp"
}

output "instance" {
  value = "${alicloud_ddosbgp_instances.instance.*.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter results by the instance name.
* `region` - (Optional) A region of instance.
* `ids` - (Optional) A list of instance IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of instance IDs.
* `names` - A list of instance names.
* `instances` - A list of apis. Each element contains the following attributes:
  * `id` - The instance's id.
  * `name` - The instance's remark.
  * `type` - The instance's type.
  * `region` - The instance's region.
  * `base_bandwidth` - The instance's base defend bandwidth.
  * `bandwidth` - The instance's elastic defend bandwidth.
  * `ip_type` - The instance's IP version.
  * `ip_count` - The instance's count of ip config.
  * `normal_bandwidth` - Normal defend bandwidth of the instance. The unit is Gbps.

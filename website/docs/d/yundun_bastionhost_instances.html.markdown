---
subcategory: "Cloud Bastionhost"
layout: "alicloud"
page_title: "Alicloud: alicloud_yundun_bastionhost_instances"
sidebar_current: "docs-alicloud-yundun_bastionhost-instances"
description: |-
  Provides a list of cloud Bastionhost(yundun_bastionhost) instances available to the user.
---

# alicloud\_yundun_bastionhost\_instances

This data source provides a list of cloud Bastionhost instances in an Alibaba Cloud account according to the specified filters.

-> **NOTE:** Available in 1.63.0+ .

## Example Usage

```
data "alicloud_yundun_bastionhost_instances" "instance" {
  name_regex = "^bastionhost"
}

output "instance" {
  value = "${alicloud_yundun_bastionhost_instances.instance.*.id}"
}
```

## Argument Reference

The following arguments are supported:

* `description_regex` - (Optional) A regex string to filter results by the instance description.
* `ids` - (Optional) Matched instance IDs to filter data source result.
* `output_file` - (Optional) File name to persist data source output.
* `descriptions` - (Optional) Descriptions to filter data source result.


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `instances` - A list of apis. Each element contains the following attributes:
  * `id` - The instance's id.
  * `description` - The instance's remark.
  * `user_vswitch_id` - The instance's vSwitch ID.
  * `private_domain` - The instance's private domain name.
  * `public_domain` - The instance's public domain name.
  * `instance_status` - The instance's status.
  * `public_network_access` - The instance's public network access configuration.
  * `security_group_ids` - The instance's security group configuration.
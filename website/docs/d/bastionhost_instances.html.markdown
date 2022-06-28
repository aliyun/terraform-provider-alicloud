---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_instances"
sidebar_current: "docs-alicloud-bastionhost-instances"
description: |-
  Provides a list of cloud Bastionhost instances available to the user.
---

# alicloud\_yundun_bastionhost\_instances

-> **NOTE:** From the version 1.132.0, the data source has been renamed to `alicloud_bastionhost_instances`.

This data source provides a list of cloud Bastionhost instances in an Alibaba Cloud account according to the specified filters.

-> **NOTE:** Available in 1.63.0+ .

## Example Usage

```
data "alicloud_bastionhost_instances" "instance" {
  description_regex = "^bastionhost"
}

output "instance" {
  value = "${alicloud_bastionhost_instances.instance.*.id}"
}
```

## Argument Reference

The following arguments are supported:

* `description_regex` - (Optional) A regex string to filter results by the instance description.
* `ids` - (Optional) Matched instance IDs to filter data source result.
* `output_file` - (Optional) File name to persist data source output.
* `descriptions` - (Optional) Descriptions to filter data source result.
* `tags` - (Optional, Available in v1.67.0+) A map of tags assigned to the bastionhost instance. It must be in the format:
  ```
  data "alicloud_bastionhost_instances" "instance" {
    tags = {
      tagKey1 = "tagValue1"
    }
  }
  ```


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
  * `tags` - A map of tags assigned to the bastionhost instance.

---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_instances"
sidebar_current: "docs-alicloud-bastionhost-instances"
description: |-
  Provides a list of cloud Bastionhost instances available to the user.
---

# alicloud_bastionhost_instances

-> **NOTE:** From the version 1.132.0, the data source has been renamed to `alicloud_bastionhost_instances`.

This data source provides a list of cloud Bastionhost instances in an Alibaba Cloud account according to the specified filters.

-> **NOTE:** Available since v1.63.0.

## Example Usage

```terraform
data "alicloud_bastionhost_instances" "instance" {
  description_regex = "^bastionhost"
}

output "instance" {
  value = data.alicloud_bastionhost_instances.instance.*.id
}
```

## Argument Reference

The following arguments are supported:

* `description_regex` - (Optional) A regex string to filter results by the instance description.
* `ids` - (Optional, ForceNew) Matched instance IDs to filter data source result.
* `output_file` - (Optional) File name to persist data source output.
* `descriptions` - (Optional) Descriptions to filter data source result.
* `tags` - (Optional) A map of tags assigned to the bastionhost instance.

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
  * `license_code` - The instance's license code.
  * `tags` - A map of tags assigned to the bastionhost instance.
  * `storage` - The storage of Cloud Bastionhost instance in TB.
  * `bandwidth` - The bandwidth of Cloud Bastionhost instance.

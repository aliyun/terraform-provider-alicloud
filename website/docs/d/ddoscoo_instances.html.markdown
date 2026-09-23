---
subcategory: "Anti-DDoS Pro (DdosCoo)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ddoscoo_instances"
sidebar_current: "docs-alicloud-datasource-ddoscoo-instances"
description: |-
  Provides a list of BGP-Line Anti-DDoS Pro(DdosCoo) instances to the user.
---

# alicloud_ddoscoo_instances

This data source provides the BGP-Line Anti-DDoS Pro(DdosCoo) instances of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.39.0.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_ddoscoo_instances" "default" {
  name_regex = "tf"
}

output "instance" {
  value = data.alicloud_ddoscoo_instances.default.instances.*.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, List) A list of instance IDs.
* `name_regex` - (Optional) A regex string to filter results by the instance name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

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
  * `ip_mode` - The ip mode of the instance.
  * `debt_status` - The debt status of the instance.
  * `edition` - The edition of the instance.
  * `ip_version` - The ip version of the instance.
  * `status` - The status of the instance.
  * `enabled` - The enabled of the instance.
  * `expire_time` - The expiry time of the instance.
  * `create_time` - The creation time of the instance.

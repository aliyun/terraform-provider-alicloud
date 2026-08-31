---
subcategory: "RabbitMQ (AMQP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_amqp_virtual_hosts"
sidebar_current: "docs-alicloud-datasource-amqp-virtual-hosts"
description: |-
  Provides a list of Amqp Virtual Hosts to the user.
---

# alicloud\_amqp\_virtual\_hosts

This data source provides the Amqp Virtual Hosts of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.126.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_amqp_virtual_hosts" "ids" {
  instance_id = "amqp-abc12345"
  ids         = ["my-VirtualHost-1", "my-VirtualHost-2"]
}
output "amqp_virtual_host_id_1" {
  value = data.alicloud_amqp_virtual_hosts.ids.hosts.0.id
}

data "alicloud_amqp_virtual_hosts" "nameRegex" {
  instance_id = "amqp-abc12345"
  name_regex  = "^my-VirtualHost"
}
output "amqp_virtual_host_id_2" {
  value = data.alicloud_amqp_virtual_hosts.nameRegex.hosts.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Virtual Host IDs. Its element value is same as Virtual Host Name.
* `instance_id` - (Required, ForceNew) InstanceId.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Virtual Host name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Virtual Host names.
* `hosts` - A list of Amqp Virtual Hosts. Each element contains the following attributes:
	* `id` - The ID of the Virtual Host.
	* `instance_id` - InstanceId.
	* `virtual_host_name` - VirtualHostName.

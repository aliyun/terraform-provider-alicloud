---
subcategory: "RabbitMQ (AMQP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_amqp_exchanges"
sidebar_current: "docs-alicloud-datasource-amqp-exchanges"
description: |-
  Provides a list of Amqp Exchanges to the user.
---

# alicloud\_amqp\_exchanges

This data source provides the Amqp Exchanges of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.128.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_amqp_exchanges" "ids" {
  instance_id       = "amqp-abc12345"
  virtual_host_name = "my-VirtualHost"
  ids               = ["my-Exchange-1", "my-Exchange-2"]
}
output "amqp_exchange_id_1" {
  value = data.alicloud_amqp_exchanges.ids.exchanges.0.id
}

data "alicloud_amqp_exchanges" "nameRegex" {
  instance_id       = "amqp-abc12345"
  virtual_host_name = "my-VirtualHost"
  name_regex        = "^my-Exchange"
}
output "amqp_exchange_id_2" {
  value = data.alicloud_amqp_exchanges.nameRegex.exchanges.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Exchange IDs. Its element value is same as Exchange Name.
* `instance_id` - (Required, ForceNew) The ID of the instance.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Exchange name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `virtual_host_name` - (Required, ForceNew) The name of virtual host where an exchange resides.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Exchange names.
* `exchanges` - A list of Amqp Exchanges. Each element contains the following attributes:
	* `attributes` - The attributes.
	* `auto_delete_state` - Indicates whether the Auto Delete attribute is configured.
	* `create_time` - The creation time.
	* `exchange_name` - The name of the exchange.
	* `exchange_type` - The type of the exchange.
	* `id` - The ID of the Exchange. Its value is same as Queue Name.
	* `instance_id` - The ID of the instance.
	* `virtual_host_name` - The name of virtual host where an exchange resides.

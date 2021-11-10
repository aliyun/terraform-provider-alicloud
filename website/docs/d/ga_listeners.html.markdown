---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_listeners"
sidebar_current: "docs-alicloud-datasource-ga-listeners"
description: |-
  Provides a list of Global Accelerator (GA) Listeners to the user.
---

# alicloud\_ga\_listeners

This data source provides the Global Accelerator (GA) Listeners of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.111.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ga_listeners" "example" {
  accelerator_id = "example_value"
  ids            = ["example_value"]
  name_regex     = "the_resource_name"
}

output "first_ga_listener_id" {
  value = data.alicloud_ga_listeners.example.listeners.0.id
}
```

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Required, ForceNew) The accelerator id.
* `ids` - (Optional, ForceNew, Computed)  A list of Listener IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Listener name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the listener. Valid values: `active`, `configuring`, `creating`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Listener names.
* `listeners` - A list of Ga Listeners. Each element contains the following attributes:
	* `certificates` - The certificates of the listener.
		* `id` - The id of the certificate.
		* `type` - The type of the certificate.
	* `client_affinity` - The clientAffinity of the listener.
	* `description` - The description of the listener.
	* `id` - The ID of the Listener.
	* `listener_id` - The listenerId of the listener.
	* `name` - The name of the listener. The length of the name is 2-128 characters. It starts with uppercase and lowercase letters or Chinese characters. It can contain numbers and underscores and dashes.
	* `port_ranges` - The portRanges of the listener.
		* `from_port` - The initial listening port used to receive requests and forward them to terminal nodes.
		* `to_port` - The end listening port used to receive requests and forward them to terminal nodes.
	* `protocol` - Type of network transport protocol monitored.
	* `status` - The status of the listener.

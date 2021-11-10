---
subcategory: "Simple Application Server"
layout: "alicloud"
page_title: "Alicloud: alicloud_simple_application_server_instances"
sidebar_current: "docs-alicloud-datasource-simple-application-server-instances"
description: |-
  Provides a list of Simple Application Server Instances to the user.
---

# alicloud\_simple\_application\_server\_instances

This data source provides the Simple Application Server Instances of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.135.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_simple_application_server_instances" "ids" {
  ids = ["example_id"]
}
output "simple_application_server_instance_id_1" {
  value = data.alicloud_simple_application_server_instances.ids.instances.0.id
}

data "alicloud_simple_application_server_instances" "nameRegex" {
  name_regex = "^my-Instance"
}
output "simple_application_server_instance_id_2" {
  value = data.alicloud_simple_application_server_instances.nameRegex.instances.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Instance IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Instance name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `payment_type` - (Optional, ForceNew) The paymen type of the resource. Valid values: `Subscription`.
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `Resetting`, `Running`, `Stopped`, `Upgrading`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Instance names.
* `instances` - A list of Simple Application Server Instances. Each element contains the following attributes:
	* `business_status` - The billing status of the simple application server. Valid values: `Normal`, `Expired` and `Overdue`.
	* `create_time` - The time when the simple application server was created.
	* `ddos_status` - The DDoS protection status. Valid values: `Normal`, `BlackHole`, and `Defense`.
	* `expired_time` - The time when the simple application server expires.
	* `id` - The ID of the Instance.
	* `image_id` - The ID of the simple application server Image.
	* `inner_ip_address` - The internal IP address of the simple application server.
	* `instance_id` - The ID of the simple application server.
	* `instance_name` - The name of the resource.
	* `payment_type` - The billing method of the simple application server.
	* `plan_id` - The ID of the simple application server plan.
	* `public_ip_address` - The public IP address of the simple application server.
	* `status` - The status of the resource.

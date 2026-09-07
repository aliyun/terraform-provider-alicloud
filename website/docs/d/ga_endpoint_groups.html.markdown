---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_endpoint_groups"
sidebar_current: "docs-alicloud-datasource-ga-endpoint-groups"
description: |-
  Provides a list of Global Accelerator (GA) Endpoint Groups to the user.
---

# alicloud_ga_endpoint_groups

This data source provides the Global Accelerator (GA) Endpoint Groups of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.113.0.

## Example Usage

Basic Usage

```terraform
variable "region" {
  default = "cn-hangzhou"
}

variable "name" {
  default = "tf-example"
}

provider "alicloud" {
  region = var.region
}

data "alicloud_ga_accelerators" "default" {
  status = "active"
}

resource "alicloud_ga_bandwidth_package" "default" {
  bandwidth              = 100
  type                   = "Basic"
  bandwidth_type         = "Basic"
  payment_type           = "PayAsYouGo"
  billing_type           = "PayBy95"
  ratio                  = 30
  bandwidth_package_name = var.name
  auto_pay               = true
  auto_use_coupon        = true
}

resource "alicloud_ga_bandwidth_package_attachment" "default" {
  accelerator_id       = data.alicloud_ga_accelerators.default.ids.0
  bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
}

resource "alicloud_ga_listener" "default" {
  accelerator_id  = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  client_affinity = "SOURCE_IP"
  protocol        = "UDP"
  name            = var.name
  port_ranges {
    from_port = "60"
    to_port   = "70"
  }
}

resource "alicloud_eip_address" "default" {
  bandwidth            = "10"
  internet_charge_type = "PayByBandwidth"
  address_name         = var.name
}

resource "alicloud_ga_endpoint_group" "default" {
  accelerator_id                = alicloud_ga_listener.default.accelerator_id
  listener_id                   = alicloud_ga_listener.default.id
  description                   = var.name
  name                          = var.name
  threshold_count               = 4
  traffic_percentage            = 20
  endpoint_group_region         = "cn-hangzhou"
  health_check_interval_seconds = "3"
  health_check_path             = "/healthcheck"
  health_check_port             = "9999"
  health_check_protocol         = "http"
  port_overrides {
    endpoint_port = "10"
    listener_port = "60"
  }
  endpoint_configurations {
    endpoint = alicloud_eip_address.default.ip_address
    type     = "PublicIp"
    weight   = "20"
  }
}

data "alicloud_ga_endpoint_groups" "default" {
  accelerator_id = alicloud_ga_endpoint_group.default.accelerator_id
  ids            = [alicloud_ga_endpoint_group.default.id]
}

output "first_ga_endpoint_group_id" {
  value = data.alicloud_ga_endpoint_groups.default.groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List) A list of Endpoint Group IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Endpoint Group name.
* `accelerator_id` - (Required, ForceNew) The ID of the Global Accelerator instance to which the endpoint group will be added.
* `listener_id` - (Optional, ForceNew) The ID of the listener that is associated with the endpoint group.
* `endpoint_group_type` - (Optional, ForceNew) The endpoint group type. Default value: `default`. Valid values: `default`, `virtual`.
* `status` - (Optional, ForceNew) The status of the endpoint group. Valid values: `active`, `configuring`, `creating`, `init`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Endpoint Group names.
* `groups` - A list of Ga Endpoint Groups. Each element contains the following attributes:
  * `id` - The ID of the Endpoint Group.
  * `endpoint_group_id` - The endpoint_group_id of the Endpoint Group.
  * `listener_id` - The ID of the listener that is associated with the endpoint group.
  * `endpoint_group_region` - The ID of the region where the endpoint group is deployed.
  * `name` - The name of the endpoint group.
  * `description` - The description of the endpoint group.
  * `health_check_interval_seconds` - The interval between two consecutive health checks. Unit: seconds.
  * `health_check_path` - The path specified as the destination of the targets for health checks.
  * `health_check_port` - The port that is used for health checks.
  * `health_check_protocol` - The protocol that is used to connect to the targets for health checks.
  * `threshold_count` - The number of consecutive failed heath checks that must occur before the endpoint is deemed unhealthy.
  * `traffic_percentage` - The weight of the endpoint group when the corresponding listener is associated with multiple endpoint groups.
  * `endpoint_group_ip_list` - (Available since v1.213.1) The list of endpoint group IP addresses.
  * `status` - The status of the endpoint group.
  * `port_overrides` - Mapping between listening port and forwarding port of boarding point.
    * `endpoint_port` - Forwarding port.
    * `listener_port` - Listener port.
  * `endpoint_configurations` - The endpointConfigurations of the endpoint group.
    * `endpoint` - The IP address or domain name of Endpoint N in the endpoint group.
	* `probe_protocol` - Probe Protocol.
	* `probe_port` - Probe Port.
    * `type` - The type of Endpoint N in the endpoint group.
    * `weight` - The weight of Endpoint N in the endpoint group.
    * `enable_clientip_preservation` - Indicates whether client IP addresses are reserved.

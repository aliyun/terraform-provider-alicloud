---
subcategory: "Anti-DDoS Pro"
layout: "alicloud"
page_title: "Alicloud: alicloud_ddoscoo_ports"
sidebar_current: "docs-alicloud-datasource-ddoscoo-ports"
description: |-
  Provides a list of Ddoscoo Ports to the user.
---

# alicloud\_ddoscoo\_ports

This data source provides the Ddoscoo Ports of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.123.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ddoscoo_ports" "example" {
  instance_id = "ddoscoo-cn-6ja1rl4j****"
  ids         = ["ddoscoo-cn-6ja1rl4j****:7001:tcp"]
}

output "first_ddoscoo_port_id" {
  value = data.alicloud_ddoscoo_ports.example.ports.0.id
}
```

## Argument Reference

The following arguments are supported:

* `frontend_port` - (Optional, ForceNew) The forwarding port.
* `frontend_protocol` - (Optional, ForceNew) The forwarding protocol. Valid values `tcp` and `udp`.
* `ids` - (Optional, ForceNew, Computed)  A list of Port IDs.
* `instance_id` - (Required, ForceNew) The Ddoscoo instance ID.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `ports` - A list of Ddoscoo Ports. Each element contains the following attributes:
	* `backend_port` - The source station port.
	* `frontend_port` - The forwarding port.
    * `frontend_protocol` - The forwarding protocol. 
    * `instance_id` - The Ddoscoo instance ID.
	* `id` - The ID of the Port.
	* `real_servers` - List of source IP addresses.

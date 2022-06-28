---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_forward_entries"
sidebar_current: "docs-alicloud-datasource-forward-entries"
description: |-
    Provides a list of Forward Entries owned by an Alibaba Cloud account.
---

# alicloud\_forward\_entries

This data source provides a list of Forward Entries owned by an Alibaba Cloud account.

-> **NOTE:** Available in 1.37.0+.

## Example Usage

```terraform
variable "name" {
  default = "forward-entry-config-example-name"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = "${alicloud_vpc.default.id}"
  cidr_block   = "172.16.0.0/21"
  zone_id      = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_nat_gateway" "default" {
  vpc_id        = "${alicloud_vpc.default.id}"
  specification = "Small"
  name          = "${var.name}"
}

resource "alicloud_eip_address" "default" {
  address_name = "${var.name}"
}

resource "alicloud_eip_association" "default" {
  allocation_id = "${alicloud_eip_address.default.id}"
  instance_id   = "${alicloud_nat_gateway.default.id}"
}

resource "alicloud_forward_entry" "default" {
  forward_table_id = "${alicloud_nat_gateway.default.forward_table_ids}"
  external_ip      = "${alicloud_eip_address.default.ip_address}"
  external_port    = "80"
  ip_protocol      = "tcp"
  internal_ip      = "172.16.0.3"
  internal_port    = "8080"
}

data "alicloud_forward_entries" "default" {
  forward_table_id = "${alicloud_forward_entry.default.forward_table_id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Forward Entries IDs.
* `name_regex` - (Optional, Available in 1.44.0+) A regex string to filter results by forward entry name.
* `external_ip` - (Optional) The public IP address.
* `internal_ip` - (Optional) The private IP address.
* `forward_table_id` - (Required) The ID of the Forward table.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `forward_entry_name` - (Optional, ForceNew, Available in 1.119.1+) The name of forward entry.
* `internal_port` - (Optional, ForceNew, Available in 1.119.1+) The internal port.
* `ip_protocol` - (Optional, ForceNew, Available in 1.119.1+) The ip protocol. Valid values: `any`,`tcp` and `udp`. 
* `status` - (Optional, ForceNew, Available in 1.119.1+) The status of farward entry. Valid value `Available`, `Deleting` and `Pending`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Forward Entries IDs.
* `names` - A list of Forward Entries names.
* `entries` - A list of Forward Entries. Each element contains the following attributes:
  * `id` - The ID of the Forward Entry.
  * `external_ip` - The public IP address.
  * `external_port` - The public port.
  * `ip_protocol` - The protocol type.
  * `internal_ip` - The private IP address.
  * `internal_port` - The private port.
  * `name` - The forward entry name.
  * `status` - The status of the Forward Entry.
  * `forward_entry_id` - The forward entry ID.
  * `forward_entry_name` - The name of forward entry.
  * `status` - The status of forward entry.


---
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

```
variable "name" {
	default = "tf-testAccForwardEntryConfig"
}

resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_nat_gateway" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	specification = "Small"
	name = "${var.name}"
}

resource "alicloud_eip" "foo" {
	name = "${var.name}"
}

resource "alicloud_eip_association" "foo" {
	allocation_id = "${alicloud_eip.foo.id}"
	instance_id = "${alicloud_nat_gateway.foo.id}"
}

resource "alicloud_forward_entry" "foo" {
	forward_table_id = "${alicloud_nat_gateway.foo.forward_table_ids}"
	external_ip = "${alicloud_eip.foo.ip_address}"
	external_port = "80"
	ip_protocol = "tcp"
	internal_ip = "172.16.0.3"
	internal_port = "8080"
}

data "alicloud_forward_entries" "foo" {
    forward_table_id = "${alicloud_forward_entry.foo.forward_table_id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Forward Entries IDs.
* `external_ip` - (Optional) The public IP address.
* `internal_ip` - (Optional) The private IP address.
* `forward_table_id` - (Required, ForceNew) The ID of the Forward table.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Forward Entries IDs.
* `entries` - A list of Forward Entries. Each element contains the following attributes:
  * `id` - The ID of the Forward Entry.
  * `external_ip` - The public IP address.
  * `external_port` - The public port.
  * `ip_protocol` - The protocol type.
  * `internal_ip` - The private IP address.
  * `internal_port` - The private port.
  * `status` - The status of the Forward Entry.


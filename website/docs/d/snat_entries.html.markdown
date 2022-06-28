---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_snat_entries"
sidebar_current: "docs-alicloud-datasource-snat-entries"
description: |-
    Provides a list of Snat Entries owned by an Alibaba Cloud account.
---

# alicloud\_snat\_entries

This data source provides a list of Snat Entries owned by an Alibaba Cloud account.

-> **NOTE:** Available in 1.37.0+.

## Example Usage

```terraform
variable "name" {
  default = "snat-entry-example-name"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "foo" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
  vpc_id            = "${alicloud_vpc.foo.id}"
  cidr_block        = "172.16.0.0/21"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name      = var.name
}

resource "alicloud_nat_gateway" "foo" {
  vpc_id        = "${alicloud_vpc.foo.id}"
  specification = "Small"
  name          = "${var.name}"
}

resource "alicloud_eip_address" "foo" {
  address_name = "${var.name}"
}

resource "alicloud_eip_association" "foo" {
  allocation_id = "${alicloud_eip_address.foo.id}"
  instance_id   = "${alicloud_nat_gateway.foo.id}"
}

resource "alicloud_snat_entry" "foo" {
  snat_table_id     = "${alicloud_nat_gateway.foo.snat_table_ids}"
  source_vswitch_id = "${alicloud_vswitch.foo.id}"
  snat_ip           = "${alicloud_eip_address.foo.ip_address}"
}

data "alicloud_snat_entries" "foo" {
  snat_table_id = "${alicloud_snat_entry.foo.snat_table_id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Snat Entries IDs.
* `snat_ip` - (Optional) The public IP of the Snat Entry.
* `source_cidr` - (Optional) The source CIDR block of the Snat Entry.
* `snat_table_id` - (Required) The ID of the Snat table.
* `name_regex` - (Optional, ForceNew, Available in 1.119.1+) A regex string to filter results by the resource name. 
* `snat_entry_name` - (Optional, ForceNew, Available in 1.119.1+) The name of snat entry.
* `source_vswitch_id` - (Optional, ForceNew, Available in 1.119.1+) The source vswitch ID.
* `status` - (Optional, ForceNew, Available in 1.119.1+) The status of the Snat Entry. Valid values: `Available`, `Deleting` and `Pending`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - (Optional) A list of Snat Entries IDs.
* `entries` - A list of Snat Entries. Each element contains the following attributes:
  * `id` - The ID of the Snat Entry.
  * `snat_ip` - The public IP of the Snat Entry.
  * `source_cidr` - The source CIDR block of the Snat Entry.
  * `status` - The status of the Snat Entry.
  * `snat_entry_id` - The ID of snat entry.
  * `snat_entry_name` - The name of snat entry.
  * `source_vswitch_id` - The source vswitch ID.


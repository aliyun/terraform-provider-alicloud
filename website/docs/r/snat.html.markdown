---
layout: "alicloud"
page_title: "Alicloud: alicloud_snat_entry"
sidebar_current: "docs-alicloud-resource-vpc"
description: |-
  Provides a Alicloud snat resource.
---

# alicloud\_snat

Provides a snat resource.

## Example Usage

Basic Usage

```
resource "alicloud_vpc" "foo" {
  ...
}

resource "alicloud_vswitch" "foo" {
  ...
}

resource "alicloud_nat_gateway" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  spec   = "Small"
  name   = "test_foo"

  bandwidth_packages = [
    {
      ip_count  = 2
      bandwidth = 5
      zone      = ""
    },
    {
      ip_count  = 1
      bandwidth = 6
      zone      = "cn-beijing-b"
    }
  ]

  depends_on = [
    "alicloud_vswitch.foo"
  ]
}

resource "alicloud_snat_entry" "foo" {
  snat_table_id     = "${alicloud_nat_gateway.foo.snat_table_ids}"
  source_vswitch_id = "${alicloud_vswitch.foo.id}"
  snat_ip           = "${alicloud_nat_gateway.foo.bandwidth_packages.0.public_ip_addresses}"
}
```

## Argument Reference

The following arguments are supported:

* `snat_table_id` - (Required, ForcesNew) The value can get from `alicloud_nat_gateway` Attributes "snat_table_ids".
* `source_vswitch_id` - (Required, ForcesNew) The vswitch ID.
* `snat_ip` - (Required) The SNAT ip address, the ip must along bandwidth package public ip which `alicloud_nat_gateway` argument `bandwidth_packages`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the snat entry. The value formats as `<snat_table_id>:<snat entry id>`

## Import

Snat Entry can be imported using the id, e.g.

```
$ terraform import alicloud_snat_entry.foo stb-1aece3:snat-232ce2
```

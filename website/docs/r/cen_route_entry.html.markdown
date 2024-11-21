---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_route_entry"
sidebar_current: "docs-alicloud-resource-cen-route-entry"
description: |-
  Provides a Alicloud CEN manage route entried resource.
---

# alicloud_cen_route_entry

Provides a CEN route entry resource. Cloud Enterprise Network (CEN) supports publishing and withdrawing route entries of attached networks. You can publish a route entry of an attached VPC or VBR to a CEN instance, then other attached networks can learn the route if there is no route conflict. You can withdraw a published route entry when CEN does not need it any more.

For information about CEN route entries publishment and how to use it, see [Manage network routes](https://www.alibabacloud.com/help/doc-detail/86980.htm).

-> **NOTE:** Available since v1.20.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cen_route_entry&exampleId=f3b90f92-d9e5-eaf6-e762-dadb6831e77417590e99&activeTab=example&spm=docs.r.cen_route_entry.0.f3b90f92d9&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_regions" "default" {
  current = true
}
data "alicloud_zones" "example" {
  available_resource_creation = "Instance"
}
data "alicloud_instance_types" "example" {
  availability_zone = data.alicloud_zones.example.zones.0.id
  cpu_core_count    = 1
  memory_size       = 2
}
data "alicloud_images" "example" {
  name_regex = "^ubuntu_18.*64"
  owners     = "system"
}

resource "alicloud_vpc" "example" {
  vpc_name   = "terraform-example"
  cidr_block = "172.17.3.0/24"
}
resource "alicloud_vswitch" "example" {
  vswitch_name = "terraform-example"
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = data.alicloud_zones.example.zones.0.id
}
resource "alicloud_security_group" "example" {
  name   = "terraform-example"
  vpc_id = alicloud_vpc.example.id
}

resource "alicloud_instance" "example" {
  availability_zone          = data.alicloud_zones.example.zones.0.id
  instance_name              = "terraform-example"
  image_id                   = data.alicloud_images.example.images.0.id
  instance_type              = data.alicloud_instance_types.example.instance_types.0.id
  security_groups            = [alicloud_security_group.example.id]
  vswitch_id                 = alicloud_vswitch.example.id
  internet_max_bandwidth_out = 5
}

resource "alicloud_cen_instance" "example" {
  cen_instance_name = "tf_example"
  description       = "an example for cen"
}
resource "alicloud_cen_instance_attachment" "example" {
  instance_id              = alicloud_cen_instance.example.id
  child_instance_id        = alicloud_vpc.example.id
  child_instance_type      = "VPC"
  child_instance_region_id = data.alicloud_regions.default.regions.0.id
}

resource "alicloud_route_entry" "example" {
  route_table_id        = alicloud_vpc.example.route_table_id
  destination_cidrblock = "11.0.0.0/16"
  nexthop_type          = "Instance"
  nexthop_id            = alicloud_instance.example.id
}

resource "alicloud_cen_route_entry" "example" {
  instance_id    = alicloud_cen_instance_attachment.example.instance_id
  route_table_id = alicloud_vpc.example.route_table_id
  cidr_block     = alicloud_route_entry.example.destination_cidrblock
}
```
## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the CEN.
* `route_table_id` - (Required, ForceNew) The route table of the attached VBR or VPC.
* `cidr_block` - (Required, ForceNew) The destination CIDR block of the route entry to publish.

->**NOTE:** The "alicloud_cen_instance_route_entries" resource depends on the related "alicloud_cen_instance_attachment" resource.

->**NOTE:** The "alicloud_cen_instance_attachment" resource should depend on the related "alicloud_vswitch" resource.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the resource, formatted as `<instance_id>:<route_table_id>:<cidr_block>`.

## Import

CEN instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_route_entry.example cen-abc123456:vtb-abc123:192.168.0.0/24
```


---
subcategory: "GWLB"
layout: "alicloud"
page_title: "Alicloud: alicloud_gwlb_listener"
description: |-
  Provides a Alicloud GWLB Listener resource.
---

# alicloud_gwlb_listener

Provides a GWLB Listener resource.



For information about GWLB Listener and how to use it, see [What is Listener](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.234.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-wulanchabu"
}

variable "region_id" {
  default = "cn-wulanchabu"
}

variable "zone_id2" {
  default = "cn-wulanchabu-c"
}

variable "zone_id1" {
  default = "cn-wulanchabu-b"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaulti9Axhl" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "default9NaKmL" {
  vpc_id       = alicloud_vpc.defaulti9Axhl.id
  zone_id      = var.zone_id1
  cidr_block   = "10.0.0.0/24"
  vswitch_name = format("%s1", var.name)
}

resource "alicloud_vswitch" "defaultH4pKT4" {
  vpc_id       = alicloud_vpc.defaulti9Axhl.id
  zone_id      = var.zone_id2
  cidr_block   = "10.0.1.0/24"
  vswitch_name = format("%s2", var.name)
}

resource "alicloud_gwlb_load_balancer" "defaultQ5setL" {
  vpc_id             = alicloud_vpc.defaulti9Axhl.id
  load_balancer_name = format("%s3", var.name)
  zone_mappings {
    vswitch_id = alicloud_vswitch.default9NaKmL.id
    zone_id    = var.zone_id1
  }
  address_ip_version = "Ipv4"
}

resource "alicloud_gwlb_server_group" "defaultoAkLbr" {
  scheduler = "5TCH"
  health_check_config {
    health_check_protocol        = "TCP"
    health_check_connect_port    = "80"
    health_check_connect_timeout = "5"
    health_check_domain          = "www.domain.com"
    health_check_enabled         = true
    health_check_http_code       = ["http_2xx", "http_4xx", "http_3xx"]
    health_check_interval        = "10"
    health_check_path            = "/health-check"
    healthy_threshold            = "2"
    unhealthy_threshold          = "2"
  }
  protocol          = "GENEVE"
  server_group_type = "Ip"
  connection_drain_config {
    connection_drain_enabled = true
    connection_drain_timeout = "1"
  }
  vpc_id = alicloud_vpc.defaulti9Axhl.id
  servers {
    server_id   = "10.0.0.1"
    server_ip   = "10.0.0.1"
    server_type = "Ip"
  }
  servers {
    server_id   = "10.0.0.2"
    server_ip   = "10.0.0.2"
    server_type = "Ip"
  }
  servers {
    server_id   = "10.0.0.3"
    server_ip   = "10.0.0.3"
    server_type = "Ip"
  }
  server_group_name = format("%s4", var.name)
}

resource "alicloud_gwlb_server_group" "defaultN4DOzm" {
  scheduler = "5TCH"
  health_check_config {
    health_check_protocol        = "TCP"
    health_check_connect_port    = "80"
    health_check_connect_timeout = "5"
    health_check_domain          = "www.domain.com"
    health_check_enabled         = true
    health_check_http_code       = ["http_2xx", "http_4xx", "http_3xx"]
    health_check_interval        = "10"
    health_check_path            = "/health-check"
    healthy_threshold            = "2"
    unhealthy_threshold          = "2"
  }
  protocol          = "GENEVE"
  server_group_type = "Ip"
  connection_drain_config {
    connection_drain_enabled = true
    connection_drain_timeout = "1"
  }
  vpc_id = alicloud_vpc.defaulti9Axhl.id
  servers {
    server_id   = "10.0.0.1"
    server_ip   = "10.0.0.1"
    server_type = "Ip"
  }
  servers {
    server_id   = "10.0.0.2"
    server_ip   = "10.0.0.2"
    server_type = "Ip"
  }
  servers {
    server_id   = "10.0.0.3"
    server_ip   = "10.0.0.3"
    server_type = "Ip"
  }
  server_group_name = format("%s5", var.name)
}


resource "alicloud_gwlb_listener" "default" {
  listener_description = "example-tf-lsn"
  server_group_id      = alicloud_gwlb_server_group.defaultoAkLbr.id
  load_balancer_id     = alicloud_gwlb_load_balancer.defaultQ5setL.id
}
```

## Argument Reference

The following arguments are supported:
* `dry_run` - (Optional) Specifies whether to perform only a dry run, without performing the actual request. 
* `listener_description` - (Optional) The custom listener description.

  The length is limited to 2 to 256 characters. It supports Chinese and English letters and can contain numbers, half-width commas (,), half-width periods (.), half-width semicolons (;), forward slashes (/), at(@), underscores (_), and dashes (-).
* `load_balancer_id` - (Required, ForceNew) The ID of the gateway load balancer instance.
* `server_group_id` - (Required) The ID of the server group.
* `tags` - (Optional, Map) The list of tags.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `region_id` - The ID of the region.
* `status` - The current status of the listener. Value:

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Listener.
* `delete` - (Defaults to 5 mins) Used when delete the Listener.
* `update` - (Defaults to 5 mins) Used when update the Listener.

## Import

GWLB Listener can be imported using the id, e.g.

```shell
$ terraform import alicloud_gwlb_listener.example <id>
```
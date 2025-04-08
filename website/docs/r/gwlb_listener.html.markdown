---
subcategory: "GWLB"
layout: "alicloud"
page_title: "Alicloud: alicloud_gwlb_listener"
description: |-
  Provides a Alicloud GWLB Listener resource.
---

# alicloud_gwlb_listener

Provides a GWLB Listener resource.



For information about GWLB Listener and how to use it, see [What is Listener](https://www.alibabacloud.com/help/en/slb/gateway-based-load-balancing-gwlb/developer-reference/api-gwlb-2024-04-15-createlistener).

-> **NOTE:** Available since v1.234.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_gwlb_listener&exampleId=dfc896fb-6d71-7093-5afe-2039bd5e16c34d9028b2&activeTab=example&spm=docs.r.gwlb_listener.0.dfc896fb6d&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
* `dry_run` - (Optional) Specifies whether to perform a dry run, without performing the actual request. Valid values:

  - `true`: performs only a dry run. The system checks the request for potential issues, including missing parameter values, incorrect request syntax, and service limits. If the request fails the dry run, an error code is returned. If the request passes the dry run, the `DryRunOperation` error code is returned.
  - `false` (default): performs a dry run and performs the actual request. If the request passes the dry run, a 2xx HTTP status code is returned and the operation is performed.
* `listener_description` - (Optional) The description of the listener.

  The description must be 2 to 256 characters in length, and can contain letters, digits, commas (,), periods (.), semicolons (;), forward slashes (/), at signs (@), underscores (\_), and hyphens (-).
* `load_balancer_id` - (Required, ForceNew) The GWLB instance ID.
* `server_group_id` - (Required) The server group ID.
* `tags` - (Optional, Map) The tags. You can specify at most 20 tags in each call.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `region_id` - The region ID.
* `status` - The status of the listener. 


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
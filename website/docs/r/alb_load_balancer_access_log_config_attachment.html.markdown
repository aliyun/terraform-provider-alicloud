---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_load_balancer_access_log_config_attachment"
description: |-
  Provides a Alicloud Application Load Balancer (ALB) Load Balancer Access Log Config Attachment resource.
---

# alicloud_alb_load_balancer_access_log_config_attachment

Provides a Application Load Balancer (ALB) Load Balancer Access Log Config Attachment resource.

Attachment between ALB and AccessLog.

For information about Application Load Balancer (ALB) Load Balancer Access Log Config Attachment and how to use it, see [What is Load Balancer Access Log Config Attachment](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.241.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
}

resource "random_uuid" "default" {
}

resource "alicloud_vpc" "alb_example_tf_vpc" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "alb_example_tf_j" {
  vpc_id       = alicloud_vpc.alb_example_tf_vpc.id
  zone_id      = "cn-beijing-j"
  cidr_block   = "192.168.1.0/24"
  vswitch_name = format("%s1", var.name)
}

resource "alicloud_vswitch" "alb_example_tf_k" {
  vpc_id       = alicloud_vpc.alb_example_tf_vpc.id
  zone_id      = "cn-beijing-k"
  cidr_block   = "192.168.2.0/24"
  vswitch_name = format("%s2", var.name)
}

resource "alicloud_vswitch" "defaultDSY0JJ" {
  vpc_id       = alicloud_vpc.alb_example_tf_vpc.id
  zone_id      = "cn-beijing-f"
  cidr_block   = "192.168.3.0/24"
  vswitch_name = format("%s3", var.name)
}

resource "alicloud_alb_load_balancer" "defaultDYswYo" {
  load_balancer_name    = format("%s4", var.name)
  load_balancer_edition = "Standard"
  vpc_id                = alicloud_vpc.alb_example_tf_vpc.id
  load_balancer_billing_config {
    pay_type = "PayAsYouGo"
  }
  address_type           = "Intranet"
  address_allocated_mode = "Fixed"
  zone_mappings {
    vswitch_id = alicloud_vswitch.defaultDSY0JJ.id
    zone_id    = alicloud_vswitch.defaultDSY0JJ.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.alb_example_tf_j.id
    zone_id    = alicloud_vswitch.alb_example_tf_j.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.alb_example_tf_k.id
    zone_id    = alicloud_vswitch.alb_example_tf_k.zone_id
  }
  lifecycle {
    ignore_changes = [access_log_config]
  }
}


resource "alicloud_alb_load_balancer_access_log_config_attachment" "default" {
  log_store        = "${var.name}-${random_uuid.default.result}"
  load_balancer_id = alicloud_alb_load_balancer.defaultDYswYo.id
  log_project      = "${var.name}-${random_uuid.default.result}"
}
```

## Argument Reference

The following arguments are supported:
* `load_balancer_id` - (Required, ForceNew) The ID of the load balancing instance.
* `log_project` - (Required, ForceNew) The log items shipped by the access log.
* `log_store` - (Required, ForceNew) Logstore for log delivery.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Load Balancer Access Log Config Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Load Balancer Access Log Config Attachment.

## Import

Application Load Balancer (ALB) Load Balancer Access Log Config Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_alb_load_balancer_access_log_config_attachment.example <id>
```
---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_load_balancer_security_group_attachment"
description: |-
  Provides a Alicloud ALB Load Balancer Security Group Attachment resource.
---

# alicloud_alb_load_balancer_security_group_attachment

Provides a ALB Load Balancer Security Group Attachment resource.

Bind a security group to an application-type Server Load Balancer instance.

For information about ALB Load Balancer Security Group Attachment and how to use it, see [What is Load Balancer Security Group Attachment](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.226.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_alb_load_balancer_security_group_attachment&exampleId=a97a449a-a99f-dac1-ca95-2377ca00ea72ab361494&activeTab=example&spm=docs.r.alb_load_balancer_security_group_attachment.0.a97a449aa9&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "eu-central-1"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "create_vpc" {
  cidr_block = "192.168.0.0/16"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "create_vsw_1" {
  vpc_id       = alicloud_vpc.create_vpc.id
  zone_id      = data.alicloud_zones.default.zones.0.id
  cidr_block   = "192.168.1.0/24"
  vswitch_name = var.name
}

resource "alicloud_vswitch" "create_vsw_2" {
  vpc_id       = alicloud_vpc.create_vpc.id
  zone_id      = data.alicloud_zones.default.zones.1.id
  cidr_block   = "192.168.2.0/24"
  vswitch_name = var.name
}

resource "alicloud_security_group" "create_security_group" {
  name   = var.name
  vpc_id = alicloud_vpc.create_vpc.id
}

resource "alicloud_alb_load_balancer" "create_alb" {
  load_balancer_name    = var.name
  load_balancer_edition = "Standard"
  vpc_id                = alicloud_vpc.create_vpc.id
  load_balancer_billing_config {
    pay_type = "PayAsYouGo"
  }
  address_type           = "Intranet"
  address_allocated_mode = "Fixed"
  zone_mappings {
    vswitch_id = alicloud_vswitch.create_vsw_2.id
    zone_id    = alicloud_vswitch.create_vsw_2.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.create_vsw_1.id
    zone_id    = alicloud_vswitch.create_vsw_1.zone_id
  }
}

resource "alicloud_alb_load_balancer_security_group_attachment" "default" {
  security_group_id = alicloud_security_group.create_security_group.id
  load_balancer_id  = alicloud_alb_load_balancer.create_alb.id
}
```

## Argument Reference

The following arguments are supported:
* `load_balancer_id` - (Required, ForceNew) The ID of the load balancing instance.
* `security_group_id` - (Optional, ForceNew, Computed) Security group ID collection.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<load_balancer_id>:<security_group_id>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Load Balancer Security Group Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Load Balancer Security Group Attachment.

## Import

ALB Load Balancer Security Group Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_alb_load_balancer_security_group_attachment.example <load_balancer_id>:<security_group_id>
```
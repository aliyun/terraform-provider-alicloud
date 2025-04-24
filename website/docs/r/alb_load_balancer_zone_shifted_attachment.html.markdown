---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_load_balancer_zone_shifted_attachment"
description: |-
  Provides a Alicloud Application Load Balancer (ALB) Load Balancer Zone Shifted Attachment resource.
---

# alicloud_alb_load_balancer_zone_shifted_attachment

Provides a Application Load Balancer (ALB) Load Balancer Zone Shifted Attachment resource.

Application load balancer start-stop zone.

For information about Application Load Balancer (ALB) Load Balancer Zone Shifted Attachment and how to use it, see [What is Load Balancer Zone Shifted Attachment](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.242.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_alb_load_balancer_zone_shifted_attachment&exampleId=2a10401b-8c0c-2f46-8529-1b110f8136b6bb30783d&activeTab=example&spm=docs.r.alb_load_balancer_zone_shifted_attachment.0.2a10401b8c&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
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

resource "alicloud_alb_load_balancer" "default78TIYG" {
  load_balancer_edition = "Standard"
  vpc_id                = alicloud_vpc.alb_example_tf_vpc.id
  load_balancer_billing_config {
    pay_type = "PayAsYouGo"
  }
  address_type           = "Intranet"
  address_allocated_mode = "Fixed"
  zone_mappings {
    vswitch_id = alicloud_vswitch.alb_example_tf_j.id
    zone_id    = alicloud_vswitch.alb_example_tf_j.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.alb_example_tf_k.id
    zone_id    = alicloud_vswitch.alb_example_tf_k.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.defaultDSY0JJ.id
    zone_id    = alicloud_vswitch.defaultDSY0JJ.zone_id
  }
}


resource "alicloud_alb_load_balancer_zone_shifted_attachment" "default" {
  zone_id          = alicloud_vswitch.defaultDSY0JJ.zone_id
  vswitch_id       = alicloud_vswitch.defaultDSY0JJ.id
  load_balancer_id = alicloud_alb_load_balancer.default78TIYG.id
}
```

## Argument Reference

The following arguments are supported:
* `load_balancer_id` - (Required, ForceNew) The ID of the load balancing instance.
* `vswitch_id` - (Required, ForceNew) The VSwitch corresponding to the zone. By default, each zone uses one VSwitch and one subnet.
* `zone_id` - (Required, ForceNew) The ID of the zone.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<load_balancer_id>:<vswitch_id>:<zone_id>`.
* `status` - Availability zone status. Value:

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Load Balancer Zone Shifted Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Load Balancer Zone Shifted Attachment.

## Import

Application Load Balancer (ALB) Load Balancer Zone Shifted Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_alb_load_balancer_zone_shifted_attachment.example <load_balancer_id>:<vswitch_id>:<zone_id>
```
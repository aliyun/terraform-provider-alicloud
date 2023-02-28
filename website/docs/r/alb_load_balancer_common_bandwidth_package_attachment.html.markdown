---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_load_balancer_common_bandwidth_package_attachment"
sidebar_current: "docs-alicloud-resource-alb-load-balancer-common-bandwidth-package-attachment"
description: |-
  Provides a Alicloud Alb Load Balancer Common Bandwidth Package Attachment resource.
---

# alicloud_alb_load_balancer_common_bandwidth_package_attachment

Provides a Alb Load Balancer Common Bandwidth Package Attachment resource.

For information about Alb Load Balancer Common Bandwidth Package Attachment and how to use it, see [What is Load Balancer Common Bandwidth Package Attachment](https://www.alibabacloud.com/help/en/server-load-balancer/latest/attachcommonbandwidthpackagetoloadbalancer).

-> **NOTE:** Available in v1.200.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "example_name"
}

data "alicloud_alb_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default_1" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_alb_zones.default.zones.0.id
}
resource "alicloud_vswitch" "vswitch_1" {
  count        = length(data.alicloud_vswitches.default_1.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 2)
  zone_id      = data.alicloud_alb_zones.default.zones.0.id
  vswitch_name = var.name
}

data "alicloud_vswitches" "default_2" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_alb_zones.default.zones.1.id
}
resource "alicloud_vswitch" "vswitch_2" {
  count        = length(data.alicloud_vswitches.default_2.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 4)
  zone_id      = data.alicloud_alb_zones.default.zones.1.id
  vswitch_name = var.name
}

resource "alicloud_alb_load_balancer" "default" {
  vpc_id                 = data.alicloud_vpcs.default.ids.0
  address_type           = "Internet"
  address_allocated_mode = "Fixed"
  load_balancer_name     = var.name
  load_balancer_edition  = "Standard"
  load_balancer_billing_config {
    pay_type = "PayAsYouGo"
  }
  zone_mappings {
    vswitch_id = length(data.alicloud_vswitches.default_1.ids) > 0 ? data.alicloud_vswitches.default_1.ids[0] : concat(alicloud_vswitch.vswitch_1.*.id, [""])[0]
    zone_id    = data.alicloud_alb_zones.default.zones.0.id
  }
  zone_mappings {
    vswitch_id = length(data.alicloud_vswitches.default_2.ids) > 0 ? data.alicloud_vswitches.default_2.ids[0] : concat(alicloud_vswitch.vswitch_2.*.id, [""])[0]
    zone_id    = data.alicloud_alb_zones.default.zones.1.id
  }
}

resource "alicloud_common_bandwidth_package" "default" {
  bandwidth            = 3
  internet_charge_type = "PayByBandwidth"
}

resource "alicloud_alb_load_balancer_common_bandwidth_package_attachment" "default" {
  bandwidth_package_id = alicloud_common_bandwidth_package.default.id
  load_balancer_id     = alicloud_alb_load_balancer.default.id
}
```

## Argument Reference

The following arguments are supported:
* `bandwidth_package_id` - (Required,ForceNew) The ID of the bound shared bandwidth package.
* `dry_run` - (Optional) Whether to PreCheck this request only. Value:-**true**: sends a check request and does not bind the shared bandwidth package to the load balancing instance. Check items include whether required parameters, request format, and business restrictions have been filled in. If the check fails, the corresponding error is returned. If the check passes, the error code 'DryRunOperation' is returned '.-**false** (default): Sends a normal request, returns the HTTP 2xx status code after the check, and directly performs the operation.
* `load_balancer_id` - (Required,ForceNew) The ID of the applied server load balancer instance.



## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.The value is formulated as `<load_balancer_id>:<bandwidth_package_id>`.
* `status` - The status of the Application Load balancing instance. Value:-**Inactive**: Stopped, indicating that the instance listener will no longer forward traffic.-**Active**: running.-**Provisioning**: The project is being created.-**Configuring**: The configuration is being changed.-**CreateFailed**: The instance cannot be deleted without any charge.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Load Balancer Common Bandwidth Package Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Load Balancer Common Bandwidth Package Attachment.

## Import

Alb Load Balancer Common Bandwidth Package Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_alb_load_balancer_common_bandwidth_package_attachment.example <load_balancer_id>:<bandwidth_package_id>
```
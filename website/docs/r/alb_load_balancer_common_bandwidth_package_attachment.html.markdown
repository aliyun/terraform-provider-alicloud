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

For information about Alb Load Balancer Common Bandwidth Package Attachment and how to use it, see [What is Load Balancer Common Bandwidth Package Attachment](https://www.alibabacloud.com/help/en/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-attachcommonbandwidthpackagetoloadbalancer).

-> **NOTE:** Available since v1.200.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_alb_load_balancer_common_bandwidth_package_attachment&exampleId=f3ebb6c7-192e-e24a-a73f-f705e591c4a566d0db00&activeTab=example&spm=docs.r.alb_load_balancer_common_bandwidth_package_attachment.0.f3ebb6c719&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}
data "alicloud_alb_zones" "default" {}
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_vswitch" "default" {
  count        = 2
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = format("10.4.%d.0/24", count.index + 1)
  zone_id      = data.alicloud_alb_zones.default.zones[count.index].id
  vswitch_name = format("${var.name}_%d", count.index + 1)
}

resource "alicloud_alb_load_balancer" "default" {
  vpc_id                 = alicloud_vpc.default.id
  address_type           = "Internet"
  address_allocated_mode = "Fixed"
  load_balancer_name     = var.name
  load_balancer_edition  = "Basic"
  resource_group_id      = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  load_balancer_billing_config {
    pay_type = "PayAsYouGo"
  }
  tags = {
    Created = "TF"
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.default.0.id
    zone_id    = data.alicloud_alb_zones.default.zones.0.id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.default.1.id
    zone_id    = data.alicloud_alb_zones.default.zones.1.id
  }
  modification_protection_config {
    status = "NonProtection"
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
* `bandwidth_package_id` - (Required, ForceNew) The ID of the bound shared bandwidth package.
* `dry_run` - (Optional) Whether to PreCheck this request only. Value:-**true**: sends a check request and does not bind the shared bandwidth package to the load balancing instance. Check items include whether required parameters, request format, and business restrictions have been filled in. If the check fails, the corresponding error is returned. If the check passes, the error code 'DryRunOperation' is returned '.-**false** (default): Sends a normal request, returns the HTTP 2xx status code after the check, and directly performs the operation.
* `load_balancer_id` - (Required, ForceNew) The ID of the applied server load balancer instance.



## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.The value is formulated as `<load_balancer_id>:<bandwidth_package_id>`.
* `status` - The status of the Application Load balancing instance. Value:-**Inactive**: Stopped, indicating that the instance listener will no longer forward traffic.-**Active**: running.-**Provisioning**: The project is being created.-**Configuring**: The configuration is being changed.-**CreateFailed**: The instance cannot be deleted without any charge.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Load Balancer Common Bandwidth Package Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Load Balancer Common Bandwidth Package Attachment.

## Import

Alb Load Balancer Common Bandwidth Package Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_alb_load_balancer_common_bandwidth_package_attachment.example <load_balancer_id>:<bandwidth_package_id>
```
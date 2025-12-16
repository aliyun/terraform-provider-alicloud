---
subcategory: "Network Load Balancer (NLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nlb_loadbalancer_common_bandwidth_package_attachment"
description: |-
  Provides a Alicloud NLB Loadbalancer Common Bandwidth Package Attachment resource.
---

# alicloud_nlb_loadbalancer_common_bandwidth_package_attachment

Provides a NLB Loadbalancer Common Bandwidth Package Attachment resource.

Bandwidth Package Operation.

For information about NLB Loadbalancer Common Bandwidth Package Attachment and how to use it, see [What is Loadbalancer Common Bandwidth Package Attachment](https://www.alibabacloud.com/help/en/server-load-balancer/latest/nlb-instances-change).

-> **NOTE:** Available since v1.209.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_nlb_loadbalancer_common_bandwidth_package_attachment&exampleId=29c3f48e-678b-ba39-c49e-4c5ab6ee67ba3b12027f&activeTab=example&spm=docs.r.nlb_loadbalancer_common_bandwidth_package_attachment.0.29c3f48e67&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_resource_manager_resource_groups" "default" {}
data "alicloud_nlb_zones" "default" {}
resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_nlb_zones.default.zones.0.id
}
resource "alicloud_vswitch" "default1" {
  vswitch_name = var.name
  cidr_block   = "10.4.1.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_nlb_zones.default.zones.1.id
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_nlb_load_balancer" "default" {
  load_balancer_name = var.name
  resource_group_id  = data.alicloud_resource_manager_resource_groups.default.ids.0
  load_balancer_type = "Network"
  address_type       = "Internet"
  address_ip_version = "Ipv4"
  vpc_id             = alicloud_vpc.default.id
  tags = {
    Created = "TF",
    For     = "example",
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.default.id
    zone_id    = data.alicloud_nlb_zones.default.zones.0.id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.default1.id
    zone_id    = data.alicloud_nlb_zones.default.zones.1.id
  }
}

resource "alicloud_common_bandwidth_package" "default" {
  bandwidth              = 2
  internet_charge_type   = "PayByTraffic"
  bandwidth_package_name = var.name
  description            = var.name
}

resource "alicloud_nlb_loadbalancer_common_bandwidth_package_attachment" "default" {
  bandwidth_package_id = alicloud_common_bandwidth_package.default.id
  load_balancer_id     = alicloud_nlb_load_balancer.default.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_nlb_loadbalancer_common_bandwidth_package_attachment&spm=docs.r.nlb_loadbalancer_common_bandwidth_package_attachment.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `bandwidth_package_id` - (Required, ForceNew) Specifies whether only to precheck the request. Valid values:
  - `true`: prechecks the request but does not disassociate the NLB instance from the EIP bandwidth plan. The system prechecks the required parameters, request syntax, and limits. If the request fails the precheck, an error message is returned. If the request passes the precheck, the `DryRunOperation` error code is returned.
  - `false` (default): sends the request. If the request passes the precheck, an HTTP 2xx status code is returned and the operation is performed.
* `load_balancer_id` - (Required, ForceNew) The ID of the EIP bandwidth plan.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<load_balancer_id>:<bandwidth_package_id>`.
* `status` - Network-based load balancing instance status. Value:, indicating that the instance listener will no longer forward traffic.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Loadbalancer Common Bandwidth Package Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Loadbalancer Common Bandwidth Package Attachment.

## Import

NLB Loadbalancer Common Bandwidth Package Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_nlb_loadbalancer_common_bandwidth_package_attachment.example <load_balancer_id>:<bandwidth_package_id>
```
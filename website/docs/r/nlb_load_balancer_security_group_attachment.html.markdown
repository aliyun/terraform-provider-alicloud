---
subcategory: "Network Load Balancer (NLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nlb_load_balancer_security_group_attachment"
description: |-
  Provides a Alicloud NLB Load Balancer Security Group Attachment resource.
---

# alicloud_nlb_load_balancer_security_group_attachment

Provides a NLB Load Balancer Security Group Attachment resource.

Security Group mount.

For information about NLB Load Balancer Security Group Attachment and how to use it, see [What is Load Balancer Security Group Attachment](https://www.alibabacloud.com/help/en/server-load-balancer/latest/loadbalancerjoinsecuritygroup).

-> **NOTE:** Available since v1.198.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_nlb_load_balancer_security_group_attachment&exampleId=860fb2cb-c63d-b129-d52d-288058f23b0c82ea7e11&activeTab=example&spm=docs.r.nlb_load_balancer_security_group_attachment.0.860fb2cbc6&intl_lang=EN_US" target="_blank">
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

resource "alicloud_nlb_load_balancer_security_group_attachment" "default" {
  security_group_id = alicloud_security_group.default.id
  load_balancer_id  = alicloud_nlb_load_balancer.default.id
}
```

## Argument Reference

The following arguments are supported:
* `dry_run` - (Optional) Specifies whether to perform a dry run, without performing the actual request. Valid values:

  - `true`: checks the request without performing the operation. The system checks the request for potential issues, including missing parameter values, incorrect request syntax, and service limits. If the request fails the dry run, an error message is returned. If the request passes the dry run, the `DryRunOperation` error code is returned.
  - `false` (default): performs a dry run and performs the actual request. If the request passes the dry run, a 2xx HTTP status code is returned and the operation is performed.
* `load_balancer_id` - (Required, ForceNew) The ID of the NLB instance to be associated with the security group.
* `security_group_id` - (Required, ForceNew, Computed) The ID of the security group to be disassociated.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<load_balancer_id>:<security_group_id>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Load Balancer Security Group Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Load Balancer Security Group Attachment.

## Import

NLB Load Balancer Security Group Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_nlb_load_balancer_security_group_attachment.example <load_balancer_id>:<security_group_id>
```
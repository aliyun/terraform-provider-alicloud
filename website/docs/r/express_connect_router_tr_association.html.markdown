---
subcategory: "Express Connect Router"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_router_tr_association"
description: |-
  Provides a Alicloud Express Connect Router Express Connect Router Tr Association resource.
---

# alicloud_express_connect_router_tr_association

Provides a Express Connect Router Express Connect Router Tr Association resource. Leased line gateway and TR binding relationship object.

For information about Express Connect Router Express Connect Router Tr Association and how to use it, see [What is Express Connect Router Tr Association](https://next.api.alibabacloud.com/api/ExpressConnectRouter/2023-09-01/CreateExpressConnectRouterAssociation).

-> **NOTE:** Available since v1.224.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_express_connect_router_tr_association&exampleId=dd472839-6d9b-59e6-6edb-30df6adee7ca28a7b887&activeTab=example&spm=docs.r.express_connect_router_tr_association.0.dd4728396d&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

variable "alowprefix1" {
  default = "10.0.0.0/24"
}

variable "allowprefix2" {
  default = "10.0.1.0/24"
}

variable "allowprefix3" {
  default = "10.0.2.0/24"
}

variable "allowprefix4" {
  default = "10.0.3.0/24"
}

variable "asn" {
  default = "4200001003"
}

resource "alicloud_express_connect_router_express_connect_router" "defaultpX0KlC" {
  alibaba_side_asn = var.asn
}

resource "alicloud_cen_instance" "default418DC9" {
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "defaultRYcjsc" {
  cen_id = alicloud_cen_instance.default418DC9.id
}

data "alicloud_account" "current" {
}

resource "alicloud_express_connect_router_tr_association" "default" {
  ecr_id                  = alicloud_express_connect_router_express_connect_router.defaultpX0KlC.id
  cen_id                  = alicloud_cen_instance.default418DC9.id
  transit_router_owner_id = data.alicloud_account.current.id
  allowed_prefixes = [
    "${var.alowprefix1}",
    "${var.allowprefix3}",
    "${var.allowprefix2}"
  ]
  transit_router_id     = alicloud_cen_transit_router.defaultRYcjsc.transit_router_id
  association_region_id = "cn-hangzhou"
}
```

## Argument Reference

The following arguments are supported:
* `allowed_prefixes` - (Optional) List of allowed route prefixes.
* `association_region_id` - (Required, ForceNew) The region to which the VPC or TR belongs.
* `cen_id` - (Optional, ForceNew) The ID of the CEN instance.
* `ecr_id` - (Required, ForceNew) The ID of the leased line gateway instance.
* `transit_router_id` - (Optional, ForceNew, Computed) The ID of the forwarding router instance.
* `transit_router_owner_id` - (Optional, ForceNew) The ID of the Alibaba Cloud account to which the forwarding router belongs.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<ecr_id>:<association_id>:<transit_router_id>`.
* `association_id` - The first ID of the resource.
* `create_time` - The creation time of the resource.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Express Connect Router Tr Association.
* `delete` - (Defaults to 5 mins) Used when delete the Express Connect Router Tr Association.
* `update` - (Defaults to 5 mins) Used when update the Express Connect Router Tr Association.

## Import

Express Connect Router Express Connect Router Tr Association can be imported using the id, e.g.

```shell
$ terraform import alicloud_express_connect_router_tr_association.example <ecr_id>:<association_id>:<transit_router_id>
```
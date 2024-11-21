---
subcategory: "Express Connect Router"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_router_vpc_association"
description: |-
  Provides a Alicloud Express Connect Router Express Connect Router Vpc Association resource.
---

# alicloud_express_connect_router_vpc_association

Provides a Express Connect Router Express Connect Router Vpc Association resource. Bind relationship object between leased line gateway and VPC.

For information about Express Connect Router Express Connect Router Vpc Association and how to use it, see [What is Express Connect Router Vpc Association](https://next.api.alibabacloud.com/api/ExpressConnectRouter/2023-09-01/CreateExpressConnectRouterAssociation).

-> **NOTE:** Available since v1.224.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_express_connect_router_vpc_association&exampleId=2453c592-867e-6ad0-2a20-e59a8e4f6a8b3fb1912c&activeTab=example&spm=docs.r.express_connect_router_vpc_association.0.2453c59286&intl_lang=EN_US" target="_blank">
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

resource "alicloud_vpc" "default8qAtD6" {
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_express_connect_router_express_connect_router" "defaultM9YxGW" {
  alibaba_side_asn = "65533"
}

data "alicloud_account" "current" {
}

resource "alicloud_express_connect_router_vpc_association" "default" {
  ecr_id = alicloud_express_connect_router_express_connect_router.defaultM9YxGW.id
  allowed_prefixes = [
    "172.16.4.0/24",
    "172.16.3.0/24",
    "172.16.2.0/24",
    "172.16.1.0/24"
  ]
  vpc_owner_id          = data.alicloud_account.current.id
  association_region_id = "cn-hangzhou"
  vpc_id                = alicloud_vpc.default8qAtD6.id
}
```

## Argument Reference

The following arguments are supported:
* `allowed_prefixes` - (Optional) List of allowed route prefixes.
* `association_region_id` - (Required, ForceNew) The region to which the VPC or TR belongs.
* `ecr_id` - (Required, ForceNew) The ID of the leased line gateway instance.
* `vpc_id` - (Required, ForceNew) The ID of the VPC instance.
* `vpc_owner_id` - (Optional, ForceNew) The ID of the Alibaba Cloud account to which the VPC belongs.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<ecr_id>:<association_id>:<vpc_id>`.
* `association_id` - The first ID of the resource.
* `create_time` - The creation time of the resource.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Express Connect Router Vpc Association.
* `delete` - (Defaults to 5 mins) Used when delete the Express Connect Router Vpc Association.
* `update` - (Defaults to 5 mins) Used when update the Express Connect Router Vpc Association.

## Import

Express Connect Router Express Connect Router Vpc Association can be imported using the id, e.g.

```shell
$ terraform import alicloud_express_connect_router_vpc_association.example <ecr_id>:<association_id>:<vpc_id>
```
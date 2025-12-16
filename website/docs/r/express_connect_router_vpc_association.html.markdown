---
subcategory: "Express Connect Router"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_router_vpc_association"
description: |-
  Provides a Alicloud Express Connect Router Express Connect Router Vpc Association resource.
---

# alicloud_express_connect_router_vpc_association

Provides a Express Connect Router Express Connect Router Vpc Association resource. Bind relationship object between leased line gateway and VPC.

For information about Express Connect Router Express Connect Router Vpc Association and how to use it, see [What is Express Connect Router Vpc Association](https://www.alibabacloud.com/help/en/express-connect/developer-reference/api-expressconnectrouter-2023-09-01-createexpressconnectrouterassociation).

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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_express_connect_router_vpc_association&spm=docs.r.express_connect_router_vpc_association.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `allowed_prefixes` - (Optional) The list of allowed route prefixes.
* `association_region_id` - (Required, ForceNew) The region ID of the resource to be associated.
* `ecr_id` - (Required, ForceNew) The ECR ID.
* `vpc_id` - (Required, ForceNew) The VPC ID.
* `vpc_owner_id` - (Optional, ForceNew) The ID of the Alibaba Cloud account that owns the VPC.
-> **NOTE:** If you want to connect to a network instance that belongs to a different account, `vpc_owner_id` is required.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<ecr_id>:<association_id>:<vpc_id>`.
* `association_id` - The ID of the association between the ECR and the VPC.
* `create_time` - The time when the association was created.
* `status` - The deployment state of the associated resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Express Connect Router Vpc Association.
* `delete` - (Defaults to 5 mins) Used when delete the Express Connect Router Vpc Association.
* `update` - (Defaults to 5 mins) Used when update the Express Connect Router Vpc Association.

## Import

Express Connect Router Express Connect Router Vpc Association can be imported using the id, e.g.

```shell
$ terraform import alicloud_express_connect_router_vpc_association.example <ecr_id>:<association_id>:<vpc_id>
```
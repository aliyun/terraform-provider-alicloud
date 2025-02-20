---
subcategory: "Vpc Ipam"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_ipam_ipam"
description: |-
  Provides a Alicloud Vpc Ipam Ipam resource.
---

# alicloud_vpc_ipam_ipam

Provides a Vpc Ipam Ipam resource.

IP Address Management.

For information about Vpc Ipam Ipam and how to use it, see [What is Ipam](https://next.api.alibabacloud.com/document/VpcIpam/2023-02-28/CreateIpam).

-> **NOTE:** Available since v1.234.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_ipam_ipam&exampleId=b459cbe2-7f5c-1a5b-7669-f4ec31cc9b1a233f7262&activeTab=example&spm=docs.r.vpc_ipam_ipam.0.b459cbe27f&intl_lang=EN_US" target="_blank">
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

data "alicloud_resource_manager_resource_groups" "default" {}


resource "alicloud_vpc_ipam_ipam" "default" {
  ipam_description      = "This is my first Ipam."
  ipam_name             = var.name
  operating_region_list = ["cn-hangzhou"]
}
```

## Argument Reference

The following arguments are supported:
* `ipam_description` - (Optional) The description of IPAM.
It must be 2 to 256 characters in length and must start with an uppercase letter or a Chinese character, but cannot start with 'http: // 'or 'https. If the description is not filled in, it is blank. The default value is blank.
* `ipam_name` - (Optional) The name of the resource.
* `operating_region_list` - (Required, Set) List of IPAM effective regions.
* `resource_group_id` - (Optional, Computed) The ID of the resource group.
* `tags` - (Optional, Map) The tag of the resource.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.
* `private_default_scope_id` - After an IPAM is created, the scope of the private network IPAM created by the system by default.
* `region_id` - The region ID of the resource.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Ipam.
* `delete` - (Defaults to 5 mins) Used when delete the Ipam.
* `update` - (Defaults to 5 mins) Used when update the Ipam.

## Import

Vpc Ipam Ipam can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_ipam_ipam.example <id>
```
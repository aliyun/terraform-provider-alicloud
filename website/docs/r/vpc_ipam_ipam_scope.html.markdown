---
subcategory: "Vpc Ipam"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_ipam_ipam_scope"
description: |-
  Provides a Alicloud Vpc Ipam Ipam Scope resource.
---

# alicloud_vpc_ipam_ipam_scope

Provides a Vpc Ipam Ipam Scope resource.

IP Address Management Scope.

For information about Vpc Ipam Ipam Scope and how to use it, see [What is Ipam Scope](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.234.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_ipam_ipam_scope&exampleId=7e1ba349-4006-2f7d-96e2-446df2dc83daa7a36e52&activeTab=example&spm=docs.r.vpc_ipam_ipam_scope.0.7e1ba34940&intl_lang=EN_US" target="_blank">
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

resource "alicloud_vpc_ipam_ipam" "defaultIpam" {
  operating_region_list = ["cn-hangzhou"]
  ipam_name             = var.name
}


resource "alicloud_vpc_ipam_ipam_scope" "default" {
  ipam_scope_name        = var.name
  ipam_id                = alicloud_vpc_ipam_ipam.defaultIpam.id
  ipam_scope_description = "This is a ipam scope."
  ipam_scope_type        = "private"
}
```

## Argument Reference

The following arguments are supported:
* `ipam_id` - (Required, ForceNew) The id of the Ipam instance.
* `ipam_scope_description` - (Optional) The description of the IPAM's scope of action.

  It must be 2 to 256 characters in length and must start with a lowercase letter, but cannot start with 'http:// 'or 'https. If it is not filled in, it is empty. The default value is empty.
* `ipam_scope_name` - (Optional) The name of the resource.
* `ipam_scope_type` - (Optional, ForceNew) IPAM scope of action type:
`private`.


-> **NOTE:**  Currently, only the role scope of the private network is supported.

* `resource_group_id` - (Optional, Computed, Available since v1.242.0) The ID of the resource group.
* `tags` - (Optional, Map) The tag of the resource.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.
* `region_id` - The region ID of the resource.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Ipam Scope.
* `delete` - (Defaults to 5 mins) Used when delete the Ipam Scope.
* `update` - (Defaults to 5 mins) Used when update the Ipam Scope.

## Import

Vpc Ipam Ipam Scope can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_ipam_ipam_scope.example <id>
```
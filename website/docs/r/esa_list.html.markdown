---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_list"
description: |-
  Provides a Alicloud ESA List resource.
---

# alicloud_esa_list

Provides a ESA List resource.



For information about ESA List and how to use it, see [What is List](https://www.alibabacloud.com/help/en/edge-security-acceleration/esa/user-guide/grouping).

-> **NOTE:** Available since v1.242.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_list&exampleId=6a2ccce9-8083-ef69-942d-1a157840249fa05a3fed&activeTab=example&spm=docs.r.esa_list.0.6a2ccce980&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "terraform-example"
}

resource "alicloud_esa_list" "default" {
  description = "resource example ip list"
  kind        = "ip"
  items = [
    "10.1.1.1",
    "10.1.1.2",
    "10.1.1.3"
  ]
  name = "resource_example_ip_list"
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) The description of the custom list.
* `items` - (Optional, List) The items in the custom list, which are displayed as an array.
* `kind` - (Optional, ForceNew) The type of the custom list.
* `name` - (Required) The name of the custom list.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the List.
* `delete` - (Defaults to 5 mins) Used when delete the List.
* `update` - (Defaults to 5 mins) Used when update the List.

## Import

ESA List can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_list.example <id>
```
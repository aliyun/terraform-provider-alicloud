---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_site"
description: |-
  Provides a Alicloud ESA Site resource.
---

# alicloud_esa_site

Provides a ESA Site resource.



For information about ESA Site and how to use it, see [What is Site](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.234.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_esa_site&exampleId=526a6ad2-d747-8099-80fe-4adba67e465727294d3a&activeTab=example&spm=docs.r.esa_site.0.526a6ad2d7&intl_lang=EN_US" target="_blank">
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

data "alicloud_resource_manager_resource_groups" "default" {
}

resource "alicloud_esa_rate_plan_instance" "defaultIEoDfU" {
  type         = "NS"
  auto_renew   = true
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = true
  plan_name    = "basic"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_esa_site" "default" {
  site_name         = "bcd${random_integer.default.result}.com"
  coverage          = "overseas"
  access_type       = "NS"
  instance_id       = alicloud_esa_rate_plan_instance.defaultIEoDfU.id
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}
```

## Argument Reference

The following arguments are supported:
* `access_type` - (Optional, ForceNew) Site Access Type
* `coverage` - (Optional) Acceleration area
* `instance_id` - (Required, ForceNew) The ID of the associated package instance.
* `resource_group_id` - (Optional, Computed) The ID of the resource group
* `site_name` - (Required, ForceNew) Site Name
* `tags` - (Optional, Map) Resource tags

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Creation time
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Site.
* `delete` - (Defaults to 5 mins) Used when delete the Site.
* `update` - (Defaults to 5 mins) Used when update the Site.

## Import

ESA Site can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_site.example <id>
```
---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_site"
description: |-
  Provides a Alicloud ESA Site resource.
---

# alicloud_esa_site

Provides a ESA Site resource.



For information about ESA Site and how to use it, see [What is Site](https://www.alibabacloud.com/help/en/edge-security-acceleration/esa/user-guide/site-management).

-> **NOTE:** Available since v1.234.0.

## Example Usage

Basic Usage

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
* `add_client_geolocation_header` - (Optional, Available since v1.244.0) Add the Visitor geolocation header. Value range:
  - on
  - off
* `add_real_client_ip_header` - (Optional, Available since v1.244.0) Add the "ali-real-client-ip" header containing the real client IP. Value range:
  - on
  - off
* `cache_architecture_mode` - (Optional, Available since v1.244.0) Multi-level cache architecture pattern. Value range:
edge: edge caching layer.
edge_smart: Edge Cache layer + Smart Cache layer.
edge_regional: Edge Cache layer + regional cache layer.
edge_regional_smart: Edge Cache layer + regional cache layer + intelligent cache layer.
* `coverage` - (Optional) Acceleration area
* `instance_id` - (Required, ForceNew) The ID of the associated package instance.
* `ipv6_enable` - (Optional, Available since v1.244.0) IPv6 switch. Value:
  - on
  - off
* `resource_group_id` - (Optional, ForceNew, Computed) The ID of the resource group
* `site_name` - (Required, ForceNew) Site Name
* `site_version` - (Optional, Int, Available since v1.244.0) The version number of the site. For a site with version management enabled, you can use this parameter to specify the effective site version. The default version is 0.
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
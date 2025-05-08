---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_edge_container_app_record"
description: |-
  Provides a Alicloud ESA Edge Container App Record resource.
---

# alicloud_esa_edge_container_app_record

Provides a ESA Edge Container App Record resource.



For information about ESA Edge Container App Record and how to use it, see [What is Edge Container App Record](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateEdgeContainerAppRecord).

-> **NOTE:** Available since v1.247.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_edge_container_app_record&exampleId=558722a7-c1d5-adae-7862-dfad05280eb8595c1763&activeTab=example&spm=docs.r.esa_edge_container_app_record.0.558722a7c1&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform.com"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_esa_site" "resource_Site_OriginPool_test" {
  site_name   = "gositecdn-${random_integer.default.result}.cn"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

resource "alicloud_esa_edge_container_app" "default" {
  health_check_host       = "example.com"
  health_check_type       = "l7"
  service_port            = "80"
  health_check_interval   = "5"
  edge_container_app_name = "terraform-app"
  health_check_http_code  = "http_2xx"
  health_check_uri        = "/"
  health_check_timeout    = "3"
  health_check_succ_times = "2"
  remarks                 = var.name
  health_check_method     = "HEAD"
  health_check_port       = "80"
  health_check_fail_times = "5"
  target_port             = "3000"
}

resource "alicloud_esa_edge_container_app_record" "default" {
  record_name = "tf.terraform.com"
  site_id     = alicloud_esa_site.resource_Site_OriginPool_test.id
  app_id      = alicloud_esa_edge_container_app.default.id
}
```

## Argument Reference

The following arguments are supported:
* `app_id` - (Required, ForceNew) The application ID
* `record_name` - (Required, ForceNew) The associated domain name.
* `site_id` - (Optional, ForceNew, Computed, Int) The website ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<app_id>:<record_name>`.
* `create_time` - The time when the domain name was added. The time follows the ISO 8601 standard in the YYYY-MM-DDThh:mm:ss format. The time is displayed in UTC.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Edge Container App Record.
* `delete` - (Defaults to 5 mins) Used when delete the Edge Container App Record.

## Import

ESA Edge Container App Record can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_edge_container_app_record.example <site_id>:<app_id>:<record_name>
```
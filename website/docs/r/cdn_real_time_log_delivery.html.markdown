---
subcategory: "CDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_cdn_real_time_log_delivery"
description: |-
  Provides a Alicloud CDN Real Time Log Delivery resource.
---

# alicloud_cdn_real_time_log_delivery

Provides a CDN Real Time Log Delivery resource.

Accelerate domain name real-time log push.

For information about CDN Real Time Log Delivery and how to use it, see [What is Real Time Log Delivery](https://www.alibabacloud.com/help/en/cdn/developer-reference/api-cdn-2018-05-10-createrealtimelogdelivery).

-> **NOTE:** Available since v1.134.0.

## Example Usage

Basic Usage

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_cdn_domain_new" "default" {
  scope       = "overseas"
  domain_name = "mycdndomain-${random_integer.default.result}.alicloud-provider.cn"
  cdn_type    = "web"
  sources {
    type     = "ipaddr"
    content  = "1.1.3.1"
    priority = 20
    port     = 80
    weight   = 15
  }
}


resource "alicloud_log_project" "default" {
  project_name = "terraform-example-${random_integer.default.result}"
  description  = "terraform-example"
}

resource "alicloud_log_store" "default" {
  project_name          = alicloud_log_project.default.project_name
  logstore_name         = "example-store"
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}

data "alicloud_regions" "default" {
  current = true
}

resource "alicloud_cdn_real_time_log_delivery" "default" {
  domain     = alicloud_cdn_domain_new.default.domain_name
  logstore   = alicloud_log_store.default.logstore_name
  project    = alicloud_log_project.default.project_name
  sls_region = data.alicloud_regions.default.regions.0.id
}
```

## Argument Reference

The following arguments are supported:
* `domain` - (Required, ForceNew) The accelerated domain name for which you want to disable real-time log delivery. You can specify multiple domain names and separate them with commas (,).
* `logstore` - (Required) The ID of the region where the Log Service project is deployed. You can specify multiple region IDs and separate them with commas (,).

  For more information about regions, see [Regions that support real-time log delivery](https://www.alibabacloud.com/help/en/doc-detail/144883.html).
* `project` - (Required) The name of the Logstore that collects log data from Alibaba Cloud CDN in real time. You can specify multiple Logstore names and separate them with commas (,).
* `sls_region` - (Required) The ID of the region where the Log Service project is deployed. For more information, see [Regions that support real-time log delivery](https://www.alibabacloud.com/help/en/doc-detail/144883.html).
* `status` - (Optional, Computed) Resource attribute fields that represent the status of the resource.

  Value:
  - offline
  - online

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Real Time Log Delivery.
* `delete` - (Defaults to 5 mins) Used when delete the Real Time Log Delivery.
* `update` - (Defaults to 5 mins) Used when update the Real Time Log Delivery.

## Import

CDN Real Time Log Delivery can be imported using the id, e.g.

```shell
$ terraform import alicloud_cdn_real_time_log_delivery.example <id>
```
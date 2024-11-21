---
subcategory: "CDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_cdn_real_time_log_delivery"
sidebar_current: "docs-alicloud-resource-cdn-real-time-log-delivery"
description: |-
  Provides a Alicloud CDN Real Time Log Delivery resource.
---

# alicloud_cdn_real_time_log_delivery

Provides a CDN Real Time Log Delivery resource.

For information about CDN Real Time Log Delivery and how to use it, see [What is Real Time Log Delivery](https://www.alibabacloud.com/help/en/cdn/developer-reference/api-cdn-2018-05-10-createrealtimelogdelivery).

-> **NOTE:** Available since v1.134.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cdn_real_time_log_delivery&exampleId=63f84f4b-707f-6281-1841-733021627c67951707c1&activeTab=example&spm=docs.r.cdn_real_time_log_delivery.0.63f84f4b70&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

* `domain` - (Required, ForceNew) The accelerated domain name for which you want to configure real-time log delivery. You can specify multiple domain names and separate them with commas (,).
* `logstore` - (Required, ForceNew) The name of the Logstore that collects log data from Alibaba Cloud Content Delivery Network (CDN) in real time.
* `project` - (Required, ForceNew) The name of the Log Service project that is used for real-time log delivery.
* `sls_region` - (Required, ForceNew) The region where the Log Service project is deployed.

-> **NOTE:** If your Project and Logstore services already exist, if you continue to create existing content, the created content will overwrite your existing indexes and custom reports. Please be careful to create your existing services to avoid affecting your online services after coverage.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Real Time Log Delivery. Its value is same as `domain`.
* `status` - The status of the real-time log delivery feature. Valid Values: `online` and `offline`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the Real Time Log Delivery.

## Import

CDN Real Time Log Delivery can be imported using the id, e.g.

```shell
$ terraform import alicloud_cdn_real_time_log_delivery.example <domain>
```

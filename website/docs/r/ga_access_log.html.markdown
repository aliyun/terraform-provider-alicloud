---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_access_log"
sidebar_current: "docs-alicloud-resource-ga-access-log"
description: |-
  Provides a Alicloud Global Accelerator (GA) Access Log resource.
---

# alicloud_ga_access_log

Provides a Global Accelerator (GA) Access Log resource.

For information about Global Accelerator (GA) Access Log and how to use it, see [What is Access Log](https://www.alibabacloud.com/help/en/global-accelerator/latest/api-ga-2019-11-20-attachlogstoretoendpointgroup).

-> **NOTE:** Available since v1.187.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ga_access_log&exampleId=f7919f03-fdb1-9468-ac0d-a0aebddd77cdda6a8a8c&activeTab=example&spm=docs.r.ga_access_log.0.f7919f03fd&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "region" {
  default = "cn-hangzhou"
}

provider "alicloud" {
  region  = var.region
  profile = "default"
}

resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_log_project" "default" {
  name = "terraform-example-${random_integer.default.result}"
}

resource "alicloud_log_store" "default" {
  project = alicloud_log_project.default.name
  name    = "terraform-example"
}

resource "alicloud_ga_accelerator" "default" {
  duration        = 1
  auto_use_coupon = true
  spec            = "2"
}

resource "alicloud_ga_bandwidth_package" "default" {
  bandwidth      = 100
  type           = "Basic"
  bandwidth_type = "Basic"
  payment_type   = "PayAsYouGo"
  billing_type   = "PayBy95"
  ratio          = 30
}

resource "alicloud_ga_bandwidth_package_attachment" "default" {
  accelerator_id       = alicloud_ga_accelerator.default.id
  bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
}

resource "alicloud_ga_listener" "default" {
  accelerator_id  = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  client_affinity = "SOURCE_IP"
  protocol        = "HTTP"
  name            = "terraform-example"

  port_ranges {
    from_port = 70
    to_port   = 70
  }
}

resource "alicloud_eip_address" "default" {
  bandwidth            = "10"
  internet_charge_type = "PayByBandwidth"
  address_name         = "terraform-example"
}

resource "alicloud_ga_endpoint_group" "default" {
  accelerator_id = alicloud_ga_listener.default.accelerator_id
  endpoint_configurations {
    endpoint = alicloud_eip_address.default.ip_address
    type     = "PublicIp"
    weight   = 20
  }
  endpoint_group_region = var.region
  listener_id           = alicloud_ga_listener.default.id
}

resource "alicloud_ga_access_log" "default" {
  accelerator_id     = alicloud_ga_accelerator.default.id
  listener_id        = alicloud_ga_listener.default.id
  endpoint_group_id  = alicloud_ga_endpoint_group.default.id
  sls_project_name   = alicloud_log_project.default.name
  sls_log_store_name = alicloud_log_store.default.name
  sls_region_id      = var.region
}
```

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Required, ForceNew) The ID of the global acceleration instance.
* `listener_id` - (Required, ForceNew) The ID of the listener.
* `endpoint_group_id` - (Required, ForceNew) The ID of the endpoint group instance.
* `sls_project_name` - (Required, ForceNew) The name of the Log Service project.
* `sls_log_store_name` - (Required, ForceNew) The name of the Log Store.
* `sls_region_id` - (Required, ForceNew) The region ID of the Log Service project.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Access Log. The value formats as `<accelerator_id>:<listener_id>:<endpoint_group_id>`.
* `status` - Whether access log is enabled.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Access Log.
* `delete` - (Defaults to 3 mins) Used when delete the Access Log.

## Import

Global Accelerator (GA) Access Log can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_access_log.example <accelerator_id>:<listener_id>:<endpoint_group_id>
```
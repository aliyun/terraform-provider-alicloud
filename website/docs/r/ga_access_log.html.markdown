---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_access_log"
sidebar_current: "docs-alicloud-resource-ga-access-log"
description: |-
  Provides a Alicloud Global Accelerator (GA) Access Log resource.
---

# alicloud\_ga\_access\_log

Provides a Global Accelerator (GA) Access Log resource.

For information about Global Accelerator (GA) Access Log and how to use it, see [What is Access Log](https://www.alibabacloud.com/help/en/global-accelerator/latest/attachlogstoretoendpointgroup).

-> **NOTE:** Available in v1.187.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ga_accelerators" "default" {
  status = "active"
}

resource "alicloud_log_project" "default" {
  name = "tf-testAcc-log-project"
}

resource "alicloud_log_store" "default" {
  project = alicloud_log_project.default.name
  name    = "tf-testAcc-log-store"
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
  accelerator_id       = data.alicloud_ga_accelerators.default.accelerators.0.id
  bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
}

resource "alicloud_ga_listener" "default" {
  accelerator_id = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  port_ranges {
    from_port = 80
    to_port   = 80
  }
}

resource "alicloud_eip_address" "default" {
  payment_type = "PayAsYouGo"
}

resource "alicloud_ga_endpoint_group" "default" {
  accelerator_id = alicloud_ga_listener.default.accelerator_id
  endpoint_configurations {
    endpoint = alicloud_eip_address.default.ip_address
    type     = "PublicIp"
    weight   = 20
  }
  endpoint_group_region = "cn-hangzhou"
  listener_id           = alicloud_ga_listener.default.id
}

resource alicloud_ga_access_log "default" {
  accelerator_id     = data.alicloud_ga_accelerators.default.accelerators.0.id
  listener_id        = alicloud_ga_listener.default.id
  endpoint_group_id  = alicloud_ga_endpoint_group.default.id
  sls_project_name   = alicloud_log_project.default.name
  sls_log_store_name = alicloud_log_store.default.name
  sls_region_id      = "cn-hangzhou"
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

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Access Log.
* `delete` - (Defaults to 3 mins) Used when delete the Access Log.

## Import

Global Accelerator (GA) Access Log can be imported using the id, e.g.

```
$ terraform import alicloud_ga_access_log.example <accelerator_id>:<listener_id>:<endpoint_group_id>
```
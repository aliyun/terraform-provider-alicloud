---
subcategory: "Cloud Storage Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_storage_gateway_gateway_logging"
sidebar_current: "docs-alicloud-resource-cloud-storage-gateway-gateway-logging"
description: |-
  Provides a Alicloud Cloud Storage Gateway Gateway Logging resource.
---

# alicloud\_cloud\_storage\_gateway\_gateway\_logging

Provides a Cloud Storage Gateway Gateway Logging resource.

For information about Cloud Storage Gateway Gateway Logging and how to use it, see [What is Gateway Logging](https://www.alibabacloud.com/help/en/doc-detail/108299.htm).

-> **NOTE:** Available in v1.144.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "example"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_cloud_storage_gateway_storage_bundle" "default" {
  storage_bundle_name = var.name
}

resource "alicloud_cloud_storage_gateway_gateway" "default" {
  description              = "tf-acctestDesalone"
  gateway_class            = "Basic"
  type                     = "File"
  payment_type             = "PayAsYouGo"
  vswitch_id               = alicloud_vswitch.default.id
  release_after_expiration = true
  public_network_bandwidth = 10
  storage_bundle_id        = alicloud_cloud_storage_gateway_storage_bundle.default.id
  location                 = "Cloud"
  gateway_name             = var.name
}

resource "alicloud_log_project" "default" {
  name        = var.name
  description = "created by terraform"
}

resource "alicloud_log_store" "default" {
  project               = alicloud_log_project.default.name
  name                  = var.name
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}

resource "alicloud_cloud_storage_gateway_gateway_logging" "default" {
  gateway_id   = alicloud_cloud_storage_gateway_gateway.default.id
  sls_logstore = alicloud_log_store.default.name
  sls_project  = alicloud_log_project.default.name
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, ForceNew) The ID of the Gateway.
* `sls_logstore` - (Required, ForceNew) The name of the Log Store.
* `sls_project` - (Required, ForceNew) The name of the Project.
* `status` - (Optional, Computed) The status of the resource. Valid values: `Enabled`, `Disable`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Gateway Logging. Its value is same as `gateway_id`.

## Import

Cloud Storage Gateway Gateway Logging can be imported using the id, e.g.

```
$ terraform import alicloud_cloud_storage_gateway_gateway_logging.example <gateway_id>
```
---
subcategory: "Cloud Storage Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_storage_gateway_gateway_logging"
sidebar_current: "docs-alicloud-resource-cloud-storage-gateway-gateway-logging"
description: |-
  Provides a Alicloud Cloud Storage Gateway Gateway Logging resource.
---

# alicloud_cloud_storage_gateway_gateway_logging

Provides a Cloud Storage Gateway Gateway Logging resource.

For information about Cloud Storage Gateway Gateway Logging and how to use it, see [What is Gateway Logging](https://www.alibabacloud.com/help/en/cloud-storage-gateway/latest/creategatewaylogging).

-> **NOTE:** Available since v1.144.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_storage_gateway_gateway_logging&exampleId=40fe5e0a-1a92-1909-6369-756c7fdc7690b7cc8fe7&activeTab=example&spm=docs.r.cloud_storage_gateway_gateway_logging.0.40fe5e0a1a&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

resource "random_uuid" "default" {
}
resource "alicloud_cloud_storage_gateway_storage_bundle" "default" {
  storage_bundle_name = substr("tf-example-${replace(random_uuid.default.result, "-", "")}", 0, 16)
}

resource "alicloud_log_project" "default" {
  project_name = substr("tf-example-${replace(random_uuid.default.result, "-", "")}", 0, 16)
  description  = "terraform-example"
}
resource "alicloud_log_store" "default" {
  project_name          = alicloud_log_project.default.project_name
  logstore_name         = var.name
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
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

resource "alicloud_cloud_storage_gateway_gateway" "default" {
  gateway_name             = var.name
  description              = var.name
  gateway_class            = "Standard"
  type                     = "File"
  payment_type             = "PayAsYouGo"
  vswitch_id               = alicloud_vswitch.default.id
  release_after_expiration = false
  public_network_bandwidth = 40
  storage_bundle_id        = alicloud_cloud_storage_gateway_storage_bundle.default.id
  location                 = "Cloud"
}

resource "alicloud_cloud_storage_gateway_gateway_logging" "default" {
  gateway_id   = alicloud_cloud_storage_gateway_gateway.default.id
  sls_logstore = alicloud_log_store.default.logstore_name
  sls_project  = alicloud_log_project.default.project_name
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, ForceNew) The ID of the Gateway.
* `sls_logstore` - (Required, ForceNew) The name of the Log Store.
* `sls_project` - (Required, ForceNew) The name of the Project.
* `status` - (Optional) The status of the resource. Valid values: `Enabled`, `Disable`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Gateway Logging. Its value is same as `gateway_id`.

## Import

Cloud Storage Gateway Gateway Logging can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_storage_gateway_gateway_logging.example <gateway_id>
```
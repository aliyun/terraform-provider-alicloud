---
subcategory: "Cloud Storage Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_storage_gateway_gateway_smb_user"
sidebar_current: "docs-alicloud-resource-cloud-storage-gateway-gateway-smb-user"
description: |-
  Provides a Alicloud Cloud Storage Gateway Gateway SMB User resource.
---

# alicloud_cloud_storage_gateway_gateway_smb_user

Provides a Cloud Storage Gateway Gateway SMB User resource.

For information about Cloud Storage Gateway Gateway SMB User and how to use it, see [What is Gateway SMB User](https://www.alibabacloud.com/help/en/cloud-storage-gateway/latest/creategatewaysmbuser).

-> **NOTE:** Available since v1.142.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_storage_gateway_gateway_smb_user&exampleId=3b90bf74-208b-37d5-fae4-782fc4ff49075cfd49ab&activeTab=example&spm=docs.r.cloud_storage_gateway_gateway_smb_user.0.3b90bf7420&intl_lang=EN_US" target="_blank">
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

resource "alicloud_cloud_storage_gateway_gateway_smb_user" "default" {
  username   = "example_username"
  password   = "password"
  gateway_id = alicloud_cloud_storage_gateway_gateway.default.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cloud_storage_gateway_gateway_smb_user&spm=docs.r.cloud_storage_gateway_gateway_smb_user.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, ForceNew) The Gateway ID of the Gateway SMB User.
* `password` - (Required, ForceNew) The password of the Gateway SMB User.
* `username` - (Required, ForceNew) The username of the Gateway SMB User.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Gateway SMB User. The value formats as `<gateway_id>:<username>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Gateway SMB User.
* `delete` - (Defaults to 1 mins) Used when delete the Gateway SMB User.

## Import

Cloud Storage Gateway Gateway SMB User can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_storage_gateway_gateway_smb_user.example <gateway_id>:<username>
```

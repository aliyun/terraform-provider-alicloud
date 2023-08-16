---
subcategory: "Smart Access Gateway (Smartag)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_connect_network_attachment"
sidebar_current: "docs-alicloud-resource-cloud-connect-network-attachment"
description: |-
  Provides a Alicloud Cloud Connect Network Attachment resource.
---

# alicloud_cloud_connect_network_attachment

Provides a Cloud Connect Network Attachment resource. This topic describes how to associate a Smart Access Gateway (SAG) instance with a network instance. You must associate an SAG instance with a network instance if you want to connect the SAG to Alibaba Cloud. You can connect an SAG to Alibaba Cloud through a leased line, the Internet, or the active and standby links.

For information about Cloud Connect Network Attachment and how to use it, see [What is Cloud Connect Network Attachment](https://www.alibabacloud.com/help/en/smart-access-gateway/latest/bindsmartaccessgateway).

-> **NOTE:** Available since v1.64.0.

-> **NOTE:** Only the following regions support. [`cn-shanghai`, `cn-shanghai-finance-1`, `cn-hongkong`, `ap-southeast-1`, `ap-southeast-2`, `ap-southeast-3`, `ap-southeast-5`, `ap-northeast-1`, `eu-central-1`]

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}
variable "sag_id" {
  default = "sag-9bifkf***"
}
provider "alicloud" {
  region = "cn-shanghai"
}
resource "alicloud_cloud_connect_network" "default" {
  name        = var.name
  description = var.name
  cidr_block  = "192.168.0.0/24"
  is_default  = true
}

resource "alicloud_cloud_connect_network_attachment" "default" {
  ccn_id = alicloud_cloud_connect_network.default.id
  sag_id = var.sag_id
}
```
## Argument Reference

The following arguments are supported:

* `ccn_id` - (Required, ForceNew) The ID of the CCN instance.
* `sag_id` - (Required, ForceNew) The ID of the Smart Access Gateway instance.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Cloud Connect Network Attachment Id and formates as `<ccn_id>:<sag_id>`.

## Import

The Cloud Connect Network Attachment can be imported using the instance_id, e.g.

```shell
$ terraform import alicloud_cloud_connect_network_attachment.example ccn-abc123456:sag-abc123456
```

---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_private_zone"
sidebar_current: "docs-alicloud-resource-cen-private-zone"
description: |-
  Provides a Alicloud Cloud Enterprise Network (CEN) Private Zone resource.
---

# alicloud_cen_private_zone

Provides a Cloud Enterprise Network (CEN) Private Zone resource.

For information about Cloud Enterprise Network (CEN) Private Zone and how to use it, see [What is Private Zone](https://www.alibabacloud.com/help/en/cloud-enterprise-network/latest/api-cbn-2017-09-12-routeprivatezoneincentovpc).

-> **NOTE:** Available since v1.83.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cen_private_zone&exampleId=2b1fc921-bfcc-5b75-b8b8-a667c2b265adbd372825&activeTab=example&spm=docs.r.cen_private_zone.0.2b1fc921bf&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_regions" "default" {
  current = true
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  description       = var.name
}

resource "alicloud_cen_instance_attachment" "default" {
  instance_id              = alicloud_cen_instance.default.id
  child_instance_id        = alicloud_vpc.default.id
  child_instance_type      = "VPC"
  child_instance_region_id = data.alicloud_regions.default.regions.0.id
}

resource "alicloud_cen_private_zone" "default" {
  cen_id           = alicloud_cen_instance_attachment.default.instance_id
  access_region_id = data.alicloud_regions.default.regions.0.id
  host_vpc_id      = alicloud_vpc.default.id
  host_region_id   = data.alicloud_regions.default.regions.0.id
}
```

## Argument Reference

The following arguments are supported:

* `cen_id` - (Required, ForceNew) The ID of the CEN instance.
* `access_region_id` - (Required, ForceNew) The ID of the region where PrivateZone is accessed. This region refers to the region in which PrivateZone is accessed by clients.
* `host_vpc_id` - (Required, ForceNew) The ID of the VPC that is associated with PrivateZone.
* `host_region_id` - (Required, ForceNew) The ID of the region where PrivateZone is deployed.

->**NOTE:** The resource `alicloud_cen_private_zone` depends on the resource `alicloud_cen_instance_attachment`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Private Zone. It formats as `<cen_id>:<access_region_id>`.
* `status` - The status of the Private Zone.

## Timeouts

-> **NOTE:** Available since v1.238.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 6 mins) Used when create the Private Zone.
* `delete` - (Defaults to 5 mins) Used when delete the Private Zone.

## Import

Cloud Enterprise Network (CEN) Private Zone can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_private_zone.example <cen_id>:<access_region_id>
```

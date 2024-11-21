---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_private_zone"
sidebar_current: "docs-alicloud-resource-cen-private-zone"
description: |-
  Provides a Alicloud CEN private zone resource.
---

# alicloud_cen_private_zone

This topic describes how to configure PrivateZone access. 
PrivateZone is a VPC-based resolution and management service for private domain names. 
After you set a PrivateZone access, the Cloud Connect Network (CCN) and Virtual Border Router (VBR) attached to a CEN instance can access the PrivateZone service through CEN.

For information about CEN Private Zone and how to use it, see [Manage CEN Private Zone](https://www.alibabacloud.com/help/en/cloud-enterprise-network/latest/api-cbn-2017-09-12-routeprivatezoneincentovpc).

-> **NOTE:** Available since v1.83.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cen_private_zone&exampleId=600c9794-2481-c94b-b4d8-cf6cd6889b0f67cab137&activeTab=example&spm=docs.r.cen_private_zone.0.600c979424&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_regions" "default" {
  current = true
}

resource "alicloud_vpc" "example" {
  vpc_name   = "tf_example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_cen_instance" "example" {
  cen_instance_name = "tf_example"
  description       = "an example for cen"
}

resource "alicloud_cen_instance_attachment" "example" {
  instance_id              = alicloud_cen_instance.example.id
  child_instance_id        = alicloud_vpc.example.id
  child_instance_type      = "VPC"
  child_instance_region_id = data.alicloud_regions.default.regions.0.id
}

resource "alicloud_cen_private_zone" "default" {
  access_region_id = data.alicloud_regions.default.regions.0.id
  cen_id           = alicloud_cen_instance_attachment.example.instance_id
  host_region_id   = data.alicloud_regions.default.regions.0.id
  host_vpc_id      = alicloud_vpc.example.id
}
```
## Argument Reference

The following arguments are supported:

* `cen_id` - (Required, ForceNew) The ID of the CEN instance.
* `access_region_id` - (Required, ForceNew) The access region. The access region is the region of the cloud resource that accesses the PrivateZone service through CEN.
* `host_region_id` - (Required, ForceNew) The service region. The service region is the target region of the PrivateZone service to be accessed through CEN. 
* `host_vpc_id` - (Required, ForceNew) The VPC that belongs to the service region.

->**NOTE:** The "alicloud_cen_private_zone" resource depends on the related "alicloud_cen_instance_attachment" resource.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the resource, formatted as `<cen_id>:<access_region_id>`.
* `status` - The status of the PrivateZone service. Valid values: ["Creating", "Active", "Deleting"].

## Import

CEN Private Zone can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_private_zone.example cen-abc123456:cn-hangzhou
```

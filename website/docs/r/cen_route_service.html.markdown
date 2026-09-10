---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_route_service"
sidebar_current: "docs-alicloud-resource-cen-route-service"
description: |-
  Provides a Alicloud CEN Route Service resource.
---

# alicloud_cen_route_service

Provides a CEN Route Service resource. The virtual border routers (VBRs) and Cloud Connect Network (CCN) instances attached to Cloud Enterprise Network (CEN) instances can access the cloud services deployed in VPCs through the CEN instances.

For information about CEN Route Service and how to use it, see [What is Route Service](https://www.alibabacloud.com/help/en/cen/developer-reference/api-cbn-2017-09-12-resolveandrouteserviceincen).

-> **NOTE:** Available since v1.99.0.

-> **NOTE:** Ensure that at least one VPC in the selected region is attached to the CEN instance.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cen_route_service&exampleId=f8f33617-0a72-1655-cc33-1c8f9ee4e1fe1d454f3d&activeTab=example&spm=docs.r.cen_route_service.0.f8f336170a&intl_lang=EN_US" target="_blank">
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

resource "alicloud_cen_route_service" "example" {
  access_region_id = data.alicloud_regions.default.regions.0.id
  host_region_id   = data.alicloud_regions.default.regions.0.id
  host_vpc_id      = alicloud_vpc.example.id
  cen_id           = alicloud_cen_instance_attachment.example.instance_id
  host             = "100.118.28.52/32"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cen_route_service&spm=docs.r.cen_route_service.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `access_region_id` - (Required, ForceNew) The region of the network instances that access the cloud services.
* `cen_id` - (Required, ForceNew) The ID of the CEN instance.
* `description` - (Optional, ForceNew) The description of the cloud service.
* `host` - (Required, ForceNew) The domain name or IP address of the cloud service.
* `host_region_id` - (Required, ForceNew) The region of the cloud service.
* `host_vpc_id` - (Required, ForceNew) The VPC associated with the cloud service.

-> **NOTE:** The values of `host_region_id` and `access_region_id` must be consistent.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the cloud service. It is formatted to `<cen_id>:<host_region_id>:<host>:<access_region_id>`.
* `status` - The status of the cloud service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 6 mins) Used when creating the cen route service (until it reaches the initial `Active` status). 
* `delete` - (Defaults to 6 mins) Used when delete the cen route service. 

## Import

CEN Route Service can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_route_service.example cen-ahixm0efqh********:cn-shanghai:100.118.28.52/32:cn-shanghai
```


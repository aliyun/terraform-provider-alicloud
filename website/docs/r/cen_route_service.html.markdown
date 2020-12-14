---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_route_service"
sidebar_current: "docs-alicloud-resource-cen-route-service"
description: |-
  Provides a Alicloud CEN Route Service resource.
---

# alicloud\_cen\_route\_service

Provides a CEN Route Service resource. The virtual border routers (VBRs) and Cloud Connect Network (CCN) instances attached to Cloud Enterprise Network (CEN) instances can access the cloud services deployed in VPCs through the CEN instances.

For information about CEN Route Service and how to use it, see [What is Route Service](https://www.alibabacloud.com/help/en/doc-detail/106671.htm).

-> **NOTE:** Available in v1.99.0+.

-> **NOTE:** Ensure that at least one VPC in the selected region is attached to the CEN instance.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-test"
}

data "alicloud_vpcs" "example" {
  is_default = true
}

resource "alicloud_cen_instance" "example" {
  name = var.name
}

resource "alicloud_cen_instance_attachment" "vpc" {
  instance_id              = alicloud_cen_instance.example.id
  child_instance_id        = data.alicloud_vpcs.example.vpcs.0.id
  child_instance_type      = "VPC"
  child_instance_region_id = data.alicloud_vpcs.example.vpcs.0.region_id
}

resource "alicloud_cen_route_service" "this" {
  access_region_id = data.alicloud_vpcs.example.vpcs.0.region_id
  host_region_id   = data.alicloud_vpcs.example.vpcs.0.region_id
  host_vpc_id      = data.alicloud_vpcs.example.vpcs.0.id
  cen_id           = alicloud_cen_instance_attachment.vpc.instance_id
  host             = "100.118.28.52/32"
}
```

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

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 6 mins) Used when creating the cen route service (until it reaches the initial `Active` status). 
* `delete` - (Defaults to 6 mins) Used when delete the cen route service. 

## Import

CEN Route Service can be imported using the id, e.g.

```
$ terraform import alicloud_cen_route_service.example cen-ahixm0efqh********:cn-shanghai:100.118.28.52/32:cn-shanghai
```


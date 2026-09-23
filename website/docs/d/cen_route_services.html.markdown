---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_route_services"
sidebar_current: "docs-alicloud-datasource-cen-route-services"
description: |-
    Provides a list of CEN Route Service owned by an Alibaba Cloud account.
---

# alicloud\_cen_\_route\_services

This data source provides CEN Route Service available to the user.

-> **NOTE:** Available in v1.102.0+

## Example Usage

Basic Usage

```terraform
data "alicloud_cen_route_services" "example" {
  cen_id = "cen-7qthudw0ll6jmc****"
}

output "first_cen_route_service_id" {
  value = data.alicloud_cen_route_services.example.services.0.id
}
```

## Argument Reference

The following arguments are supported:

* `access_region_id` - (Optional, ForceNew) The region of the network instances that access the cloud services.
* `cen_id` -(Required, ForceNew) The ID of the CEN instance.
* `host` -(Optional, ForceNew) The domain name or IP address of the cloud service.
* `host_region_id` - (Optional, ForceNew) The region of the cloud service.
* `host_vpc_id` - (Optional, ForceNew) The VPC associated with the cloud service.
* `status` - (Optional, ForceNew) The status of the cloud service. Valid values: `Active`, `Creating` and `Deleting`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of CEN Route Service IDs.
* `services` - A list of CEN Route Services. Each element contains the following attributes:
  * `id` - The ID of the route service.
  * `access_region_id` - The region of the network instances that access the cloud services.
  * `cen_id` - The ID of the CEN instance.
  * `cidrs` - The IP address of the cloud service.
  * `description` - The description of the cloud service.
  * `host` - The domain name or IP address of the cloud service.
  * `host_region_id` - The region of the cloud service.
  * `host_vpc_id` - The VPC associated with the cloud service.
  * `status` - The status of the cloud service.
  * `update_interval` - The update interval. Default value: 5. The value cannot be modified.

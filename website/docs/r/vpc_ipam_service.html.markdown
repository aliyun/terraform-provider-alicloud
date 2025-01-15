---
subcategory: "Vpc Ipam"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_ipam_service"
description: |-
  Provides a Alicloud Vpc Ipam Service resource.
---

# alicloud_vpc_ipam_service

Provides a Vpc Ipam Service resource.

Ipam service, used to support automatic provisioning of Terraform.

For information about Vpc Ipam Service and how to use it, see [What is Service](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.242.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_vpc_ipam_service" "default" {
}
```

### Deleting `alicloud_vpc_ipam_service` or removing it from your configuration

Terraform cannot destroy resource `alicloud_vpc_ipam_service`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as ``.
* `enabled` - Whether the IPAM service has been activated.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Service.

## Import

Vpc Ipam Service can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_ipam_service.example 
```
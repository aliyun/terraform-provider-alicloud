---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_public_ip_address_pool_service"
description: |-
  Provides a Alicloud VPC Public Ip Address Pool Service resource.
---

# alicloud_vpc_public_ip_address_pool_service

Provides a VPC Public Ip Address Pool Service resource. IP address pool service to support automatic provisioning of Terrafrom.

For information about VPC Public Ip Address Pool Service and how to use it, see [What is Public Ip Address Pool Service](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.225.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_vpc_public_ip_address_pool_service" "default" {
  enabled = true
}
```

### Deleting `alicloud_vpc_public_ip_address_pool_service` or removing it from your configuration

Terraform cannot destroy resource `alicloud_vpc_public_ip_address_pool_service`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `enabled` - Service Status on: Opened off: Not opened.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Public Ip Address Pool Service.

## Import

VPC Public Ip Address Pool Service can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_public_ip_address_pool_service.example <id>
```
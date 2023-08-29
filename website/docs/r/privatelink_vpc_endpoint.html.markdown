---
subcategory: "Private Link"
layout: "alicloud"
page_title: "Alicloud: alicloud_privatelink_vpc_endpoint"
sidebar_current: "docs-alicloud-resource-privatelink-vpc-endpoint"
description: |-
  Provides a Alicloud Private Link Vpc Endpoint resource.
---

# alicloud_privatelink_vpc_endpoint

Provides a Private Link Vpc Endpoint resource.

For information about Private Link Vpc Endpoint and how to use it, see [What is Vpc Endpoint](https://www.alibabacloud.com/help/en/privatelink/latest/api-privatelink-2020-04-15-createvpcendpoint).

-> **NOTE:** Available since v1.109.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}

resource "alicloud_privatelink_vpc_endpoint_service" "example" {
  service_description    = var.name
  connect_bandwidth      = 103
  auto_accept_connection = false
}

resource "alicloud_vpc" "example" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_security_group" "example" {
  name   = var.name
  vpc_id = alicloud_vpc.example.id
}

resource "alicloud_privatelink_vpc_endpoint" "example" {
  service_id         = alicloud_privatelink_vpc_endpoint_service.example.id
  security_group_ids = [alicloud_security_group.example.id]
  vpc_id             = alicloud_vpc.example.id
  vpc_endpoint_name  = var.name
}
```

## Argument Reference

The following arguments are supported:

* `dry_run` - (Optional) The dry run. Default to: `false`.
* `endpoint_description` - (Optional) The description of Vpc Endpoint. The length is 2~256 characters and cannot start with `http://` and `https://`.
* `security_group_ids` - (Required) The security group associated with the terminal node network card.
* `service_id` - (Optional, ForceNew) The terminal node service associated with the terminal node.
* `service_name` - (Optional, ForceNew) The name of the terminal node service associated with the terminal node.
* `vpc_endpoint_name` - (Optional) The name of Vpc Endpoint. The length is between 2 and 128 characters, starting with English letters or Chinese, and can include numbers, hyphens (-) and underscores (_).
* `vpc_id` - (Required, ForceNew) The private network to which the terminal node belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Vpc Endpoint. Value as `endpoint_id`.
* `bandwidth` - The Bandwidth.
* `connection_status` - The status of Connection.
* `endpoint_business_status` - The status of Endpoint Business.
* `endpoint_domain` - The Endpoint Domain.
* `status` - The status of Vpc Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 6 mins) Used when create the Vpc Endpoint.
* `update` - (Defaults to 4 mins) Used when update the Vpc Endpoint.

## Import

Private Link Vpc Endpoint can be imported using the id, e.g.

```shell
$ terraform import alicloud_privatelink_vpc_endpoint.example <endpoint_id>
```

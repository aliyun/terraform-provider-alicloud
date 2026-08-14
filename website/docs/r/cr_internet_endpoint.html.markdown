---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_internet_endpoint"
sidebar_current: "docs-alicloud-resource-cr-internet-endpoint"
description: |-
  Provides a Alicloud CR Internet Endpoint resource.
---

# alicloud_cr_internet_endpoint

Provides a CR Internet Endpoint resource.

For information about CR Internet Endpoint and how to use it, see [GetInstanceEndpoint](https://www.alibabacloud.com/help/en/acr/developer-reference/api-cr-2018-12-01-getinstanceendpoint).

-> **NOTE:** Available since v1.287.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}

resource "alicloud_cr_ee_instance" "default" {
  payment_type   = "Subscription"
  period         = 1
  renewal_status = "ManualRenewal"
  instance_type  = "Advanced"
  instance_name  = "${var.name}"
}

resource "alicloud_cr_internet_endpoint" "default" {
  instance_id = alicloud_cr_ee_instance.default.id

  entries {
    entry   = "192.168.1.0/24"
    comment = "entry-1"
  }

  entries {
    entry   = "10.0.0.0/8"
    comment = "entry-2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the Container Registry Enterprise Edition instance.
* `entries` - (Optional) The ACL entries of the Internet endpoint. See [`entries`](#entries) below.

### `entries`

The `entries` supports the following:

* `entry` - (Optional) The CIDR-formatted IP address range that is allowed to access the instance over the Internet.
* `comment` - (Optional) The comment of the entry.

-> **NOTE:** When the Internet endpoint is enabled, the CIDR block `127.0.0.1/32` with comment `default` is automatically added to the whitelist as a system-managed loopback ACL policy. It cannot be created or deleted through this resource's `entries` and is filtered out of state on Read, so adding `entry = "127.0.0.1/32"` with `comment = "default"` to `entries` causes a perpetual plan diff. Removing all user-managed entries exposes the instance to the Internet.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of the Internet Endpoint, which equals to the instance ID.
* `status` - The status of the Internet endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Internet Endpoint.
* `update` - (Defaults to 5 mins) Used when update the Internet Endpoint.
* `delete` - (Defaults to 5 mins) Used when delete the Internet Endpoint.

## Import

CR Internet Endpoint can be imported using the id, e.g.

```shell
$ terraform import alicloud_cr_internet_endpoint.example <instance_id>
```

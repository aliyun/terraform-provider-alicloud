---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_internet_endpoint"
sidebar_current: "docs-alicloud-datasource-cr-internet-endpoint"
description: |-
  Provides the CR Internet Endpoint of the current Alibaba Cloud user.
---

# alicloud_cr_internet_endpoint

This data source provides the CR Internet Endpoint of the current Alibaba Cloud user.

For information about CR Internet Endpoint and how to use it, see [GetInstanceEndpoint](https://www.alibabacloud.com/help/en/acr/developer-reference/api-cr-2018-12-01-getinstanceendpoint).

-> **NOTE:** Available since v1.287.0.

## Example Usage

Basic Usage

```terraform
variable "instance_id" {
  description = "The ID of a Container Registry Enterprise Edition instance with an enabled Internet endpoint."
  type        = string
}

data "alicloud_cr_internet_endpoint" "default" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The ID of the Container Registry Enterprise Edition instance.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the Internet Endpoint, which equals to the instance ID.
* `status` - The status of the Internet endpoint.
* `entries` - The ACL entries of the Internet endpoint. Each element contains:
  * `entry` - The CIDR-formatted IP address range that is allowed to access.
  * `comment` - The comment of the entry.

-> **NOTE:** The system-managed loopback ACL policy (entry `127.0.0.1/32`, comment `default`) is filtered out of `entries` on Read; the data source only returns user-managed entries.

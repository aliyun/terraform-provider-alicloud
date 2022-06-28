---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_endpoint_acl_service"
sidebar_current: "docs-alicloud-datasource-cr-endpoint-acl-service"
description: |-
  Provides a list of Cr Endpoint Acl Service to the user.
---

# alicloud\_cr\_endpoint\_acl\_service

This data source provides the CR Endpoint Acl Service of the current Alibaba Cloud user.

For information about Event Bridge and how to use it, see [What is CR Endpoint Acl](https://www.alibabacloud.com/help/en/doc-detail/142246.htm).

-> **NOTE:** Available in v1.139.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cr_endpoint_acl_service" "example" {
  endpoint_type = "internet"
  enable        = true
  instance_id   = "example_id"
  module_name   = "Registry"
}
```

## Argument Reference

The following arguments are supported:

* `endpoint_type` - (Required, ForceNew)  The type of endpoint. Valid values: `internet`.
* `enable` - (Optional) Whether to enable Acl Service, Setting the value to `true` to enable the acl service. Valid values: `true` and `false`.
* `instance_id` - (Required, ForceNew) The ID of the CR Instance.
* `module_name` - (Optional, ForceNew) The ModuleName. Valid values: `Registry`.

-> **NOTE:** After You enable access over the Internet, the Classless Inter-Domain Routing (CIDR) block `127.0.0.1/32` is automatically added to the whitelist.

-> **NOTE:** You may want to allow all ECS instances to access the Container Registry Enterprise Edition instance over the Internet. To achieve this purpose, you can enable access over the Internet and delete all IP addresses from the whitelist for Internet access. After you perform the preceding operation, the Container Registry Enterprise Edition instance is completely exposed to the Internet and may be attacked.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `service` - A list of Cr Endpoint Acl Service. Each element contains the following attributes:
	* `enable` - Whether to enable Acl Service.  Valid values: `true` and `false`.
	* `endpoint_type` - The type of endpoint. Valid values: `internet`.
	* `id` - The ID of the Endpoint Acl Service.
	* `instance_id` - The ID of the CR Instance.
	* `status` - The status of the resource.

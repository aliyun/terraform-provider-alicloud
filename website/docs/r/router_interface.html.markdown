---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_router_interface"
sidebar_current: "docs-alicloud-resource-router-interface"
description: |-
  Provides a VPC router interface resource to connect two VPCs.
---

# alicloud\_router\_interface

Provides a VPC router interface resource aim to build a connection between two VPCs.

-> **NOTE:** Only one pair of connected router interfaces can exist between two routers. Up to 5 router interfaces can be created for each router and each account.

-> **NOTE:** The router interface is not connected when it is created. It can be connected by means of resource [alicloud_router_interface_connection](https://www.terraform.io/docs/providers/alicloud/r/router_interface_connection).


## Example Usage

```
resource "alicloud_vpc" "foo" {
  name       = "tf_test_foo12345"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_router_interface" "interface" {
  opposite_region = "cn-beijing"
  router_type     = "VRouter"
  router_id       = alicloud_vpc.foo.router_id
  role            = "InitiatingSide"
  specification   = "Large.2"
  name            = "test1"
  description     = "test1"
}
```
## Argument Reference

The following arguments are supported:

* `opposite_region` - (Required, ForceNew) The Region of peer side.
* `router_type` - (Required, ForceNew) Router Type. Optional value: VRouter, VBR. Accepting side router interface type only be VRouter.
* `opposite_router_type` - (Deprecated) It has been deprecated from version 1.11.0. resource alicloud_router_interface_connection's 'opposite_router_type' instead.
* `router_id` - (Required, ForceNew) The Router ID.
* `opposite_router_id` - (Deprecated) It has been deprecated from version 1.11.0. Use resource alicloud_router_interface_connection's 'opposite_router_id' instead.
* `role` - (Required, ForceNew) The role the router interface plays. Optional value: `InitiatingSide`, `AcceptingSide`.
* `specification` - (Optional) Specification of router interfaces. It is valid when `role` is `InitiatingSide`. Accepting side's role is default to set as 'Negative'. For more about the specification, refer to [Router interface specification](https://www.alibabacloud.com/help/doc-detail/36037.htm).
* `access_point_id` - (Deprecated) It has been deprecated from version 1.11.0.
* `opposite_access_point_id` - (Deprecated) It has been deprecated from version 1.11.0.
* `opposite_interface_id` - (Deprecated) It has been deprecated from version 1.11.0. Use resource alicloud_router_interface_connection's 'opposite_router_id' instead.
* `opposite_interface_owner_id` - (Deprecated) It has been deprecated from version 1.11.0. Use resource alicloud_router_interface_connection's 'opposite_interface_id' instead.
* `name` - (Optional) Name of the router interface. Length must be 2-80 characters long. Only Chinese characters, English letters, numbers, period (.), underline (_), or dash (-) are permitted.
                                                    If it is not specified, the default value is interface ID. The name cannot start with http:// and https://.
* `description` - (Optional) Description of the router interface. It can be 2-256 characters long or left blank. It cannot start with http:// and https://.
* `health_check_source_ip` - (Optional) Used as the Packet Source IP of health check for disaster recovery or ECMP. It is only valid when `router_type` is `VBR`. The IP must be an unused IP in the local VPC. It and `health_check_target_ip` must be specified at the same time.
* `health_check_target_ip` - (Optional) Used as the Packet Target IP of health check for disaster recovery or ECMP. It is only valid when `router_type` is `VBR`. The IP must be an unused IP in the local VPC. It and `health_check_source_ip` must be specified at the same time.
* `instance_charge_type` - (Optional, ForceNew) The billing method of the router interface. Valid values are "PrePaid" and "PostPaid". Default to "PostPaid". Router Interface doesn't support "PrePaid" when region and opposite_region are the same.
* `period` - (Optional) The duration that you will buy the resource, in month. It is valid when `instance_charge_type` is `PrePaid`. Valid values: [1-9, 12, 24, 36]. At present, the provider does not support modify "period" and you can do that via web console.
-> **NOTE:** The attribute `period` is only used to create Subscription instance or modify the PayAsYouGo instance to Subscription. Once effect, it will not be modified that means running `terraform apply` will not effect the resource.


## Attributes Reference

The following attributes are exported:

* `id` - Router interface ID.
* `router_id` - Router ID.
* `router_type` - Router type.
* `role` - Router interface role.
* `name` - Router interface name.
* `description` - Router interface description.
* `specification` - Router nterface specification.
* `access_point_id` - Access point of the router interface.
* `opposite_access_point_id` - (Deprecated) It has been deprecated from version 1.11.0.
* `opposite_router_type` - Peer router type.
* `opposite_router_id` - Peer router ID.
* `opposite_interface_id` - Peer router interface ID.
* `opposite_interface_owner_id` - Peer account ID.
* `health_check_source_ip` - Source IP of Packet of Line HealthCheck.
* `health_check_target_ip` - Target IP of Packet of Line HealthCheck.

## Import

The router interface can be imported using the id, e.g.

```
$ terraform import alicloud_router_interface.interface ri-abc123456
```


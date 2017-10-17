---
layout: "alicloud"
page_title: "Alicloud: alicloud_router_interface"
sidebar_current: "docs-alicloud-resource-router-interface"
description: |-
  Provides a VPC router interface resource to connect two VPCs.
---

# alicloud\_router\_interface

Provides a VPC router interface resource to connect two VPCs by connecting the router interfaces .

~> **NOTE:** Only one pair of connected router interfaces can exist between two routers. Up to 5 router interfaces can be created for each router and each account.


## Example Usage

```
resource "alicloud_vpc" "foo" {
  name = "tf_test_foo12345"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_router_interface" "interface" {
  opposite_region = "cn-beijing"
  router_type = "VRouter"
  router_id = "${alicloud_vpc.foo.router_id}"
  role = "InitiatingSide"
  specification = "Large.2"
  name = "test1"
  description = "test1"
}
```
## Argument Reference

The following arguments are supported:

* `opposite_region` - (Required, Force New) The Region of peer side. At present, optional value: `cn-beijing`, `cn-hangzhou`, `cn-shanghai`, `cn-shenzhen`, `cn-hongkong`, `ap-southeast-1`, `us-east-1`, `us-west-1`.
* `router_type` - (Required, Forces New) Router Type. Optional value: VRouter, VBR.
* `opposite_router_type` - (Optional, Force New) Peer router type. Optional value: `VRouter`, `VBR`. Default to `VRouter`.
* `router_id` - (Required, Force New) Router ID. When `router_type` is VBR, the VBR specified by the `router_id` must be in the access point specified by `access_point_id`.
* `opposite_router_id` - (Optional) Peer router ID. When `opposite_router_type` is VBR, the `opposite_router_id` must be in the access point specified by `opposite_access_point_id`.
* `role` - (Required, Force New) The role the router interface plays. Optional value: `InitiatingSide`, `AcceptingSide`.
* `specification` - (Optional) Specification of router interfaces. If `role` is `AcceptingSide`, the value can be ignore or must be `Negative`. For more about the specification, refer to [Router interface specification](https://www.alibabacloud.com/help/doc-detail/52415.htm?spm=a3c0i.o52412zh.b99.10.698e566fdVCfKD).
* `access_point_id` - (Optional, Force New) Access point ID. Required when `router_type` is `VBR`. Prohibited when `router_type` is `VRouter`.
* `opposite_access_point_id` - (Optional, Force New) Access point ID of peer side. Required when `opposite_router_type` is `VBR`. Prohibited when `opposite_router_type` is `VRouter`.
* `opposite_interface_id` - (Optional) Peer router interface ID.
* `opposite_interface_owner_id` - (Optional) Peer account ID. Log on to the Alibaba Cloud console, select User Info > Account Management to check your account ID.
* `name` - (Optional) Name of the router interface. Length must be 2-80 characters long. Only Chinese characters, English letters, numbers, period (.), underline (_), or dash (-) are permitted.
                                                    If it is not specified, the default value is interface ID. The name cannot start with http:// and https://.
* `description` - (Optional) Description of the router interface. It can be 2-256 characters long or left blank. It cannot start with http:// and https://.
* `health_check_source_ip` - (Optional) Used as the Packet Source IP of health check for disaster recovery or ECMP. It is only valid when `router_type` is `VRouter` and `opposite_router_type` is `VBR`. The IP must be an unused IP in the local VPC. It and `health_check_target_ip` must be specified at the same time.
* `health_check_target_ip` - (Optional) Used as the Packet Target IP of health check for disaster recovery or ECMP. It is only valid when `router_type` is `VRouter` and `opposite_router_type` is `VBR`. The IP must be an unused IP in the local VPC. It and `health_check_source_ip` must be specified at the same time.

~> **NOTE:**
* If `router_type` is `VBR`, the `role` must be `InitiatingSide` and `opposite_router_type` must be `VRouter`.
* If `opposite_router_type` is `VBR`, the `role` must be `AcceptingSide` and `router_type` must be `VRouter`.


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
* `opposite_access_point_id` - Access point of the opposite router interface.
* `opposite_router_type` - Peer router type.
* `opposite_router_id` - Peer router ID.
* `opposite_interface_id` - Peer router interface ID.
* `opposite_interface_owner_id` - Peer account ID.
* `health_check_source_ip` - Source IP of Packet of Line HealthCheck.
* `health_check_target_ip` - Target IP of Packet of Line HealthCheck.

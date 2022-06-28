---
subcategory: "Elastic Accelerated Computing Instances (EAIS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_eais_instance"
sidebar_current: "docs-alicloud-resource-eais-instance"
description: |-
  Provides a Alicloud EAIS Instance resource.
---

# alicloud\_eais\_instance

Provides a EAIS Instance resource.

For information about EAIS Instance and how to use it, see [What is Instance](https://help.aliyun.com/document_detail/185066.html).

-> **NOTE:** Available in v1.137.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "%v"
}
data "alicloud_vpcs" "default" {
  cidr_block = "172.16.0.0/12"
}
resource "alicloud_vpc" "default" {
  count      = length(data.alicloud_vpcs.default.ids) > 0 ? 0 : 1
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}
data "alicloud_vswitches" "default" {
  vpc_id  = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
  zone_id = "cn-hangzhou-h"
}
resource "alicloud_vswitch" "default" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 2)
  zone_id      = "cn-hangzhou-h"
  vswitch_name = var.name
}
resource "alicloud_security_group" "default" {
  name        = var.name
  description = "tf test"
  vpc_id      = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
}
resource "alicloud_eais_instance" "default" {
  instance_type     = "eais.ei-a6.4xlarge"
  instance_name     = var.name
  security_group_id = alicloud_security_group.default.id
  vswitch_id        = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : alicloud_vswitch.default[0].id
}
```

## Argument Reference

The following arguments are supported:

* `force` - (Optional) Whether to force deletion when the instance status does not meet the deletion conditions. Valid values: `true` and `false`.
* `instance_name` - (Optional, ForceNew) The name of the instance.
* `instance_type` - (Required, ForceNew) The type of the resource. Valid values: `eais.ei-a6.4xlarge`, `eais.ei-a6.2xlarge`, `eais.ei-a6.xlarge`, `eais.ei-a6.large`, `eais.ei-a6.medium`.
* `security_group_id` - (Required) The ID of the security group.
* `vswitch_id` - (Required) The ID of the vswitch.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Instance.
* `status` - The status of the resource. Valid values: `Attaching`, `Available`, `Detaching`, `InUse`, `Starting`, `Unavailable`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Instance.

## Import

EAIS Instance can be imported using the id, e.g.

```
$ terraform import alicloud_eais_instance.example <id>
```

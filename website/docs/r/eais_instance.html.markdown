---
subcategory: "Elastic Accelerated Computing Instances (EAIS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_eais_instance"
sidebar_current: "docs-alicloud-resource-eais-instance"
description: |-
  Provides a Alicloud EAIS Instance resource.
---

# alicloud_eais_instance

Provides a EAIS Instance resource.

For information about EAIS Instance and how to use it, see [What is Instance](https://www.alibabacloud.com/help/en/resource-orchestration-service/latest/aliyun-eais-instance).

-> **NOTE:** Available since v1.137.0.

## Example Usage
<div class="oics-button" style="float: right;margin: 0 0 -40px 0;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_eais_instance&exampleId=e15e2f2c-80e2-839e-bd33-da0e97ae0db0a5757284&activeTab=example&spm=docs.r.eais_instance.0.e15e2f2c80" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}
locals {
  zone_id = "cn-hangzhou-h"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}
resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.1.0.0/16"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = local.zone_id
}


resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}
resource "alicloud_eais_instance" "default" {
  instance_type     = "eais.ei-a6.2xlarge"
  instance_name     = var.name
  security_group_id = alicloud_security_group.default.id
  vswitch_id        = alicloud_vswitch.default.id
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

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Instance.

## Import

EAIS Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_eais_instance.example <id>
```

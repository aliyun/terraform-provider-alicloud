---
subcategory: "Elastic Accelerated Computing Instances (EAIS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_eais_instance"
description: |-
  Provides a Alicloud EAIS Instance resource.
---

# alicloud_eais_instance

Provides a EAIS Instance resource.

Instance resource definition.

For information about EAIS Instance and how to use it, see [What is Instance](https://www.alibabacloud.com/help/en/resource-orchestration-service/latest/aliyun-eais-instance).

-> **NOTE:** Available since v1.137.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_eais_instance&exampleId=b4b42fb3-673c-12f3-9a98-6cdf7e94e29fccb36fa6&activeTab=example&spm=docs.r.eais_instance.0.b4b42fb367&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

locals {
  zone_id = "cn-hangzhou-h"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.192.0/24"
  zone_id      = local.zone_id
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_eais_instance" "default" {
  instance_type     = "eais.ei-a6.2xlarge"
  vswitch_id        = alicloud_vswitch.default.id
  security_group_id = alicloud_security_group.default.id
  instance_name     = var.name
}
```

### Deleting `alicloud_eais_instance` or removing it from your configuration

The `alicloud_eais_instance` resource allows you to manage  `category = "ei"`  instance, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration will remove it from your state file and management, but will not destroy the Instance.
You can resume managing the subscription instance via the AlibabaCloud Console.

## Argument Reference

The following arguments are supported:
* `category` - (Optional) EAIS instance category, valid values: `eais`, `jupyter`, `ei`, default is `eais`.
* `environment_var` - (Optional, List, Available since v1.246.0) Setting environment variables in eais instance on Initialization See [`environment_var`](#environment_var) below.
* `force` - (Optional, Deprecated since v1.246.0) Whether to force the deletion when the instance status does not meet the deletion conditions.
* `image` - (Optional, Available since v1.246.0) EAIS instance image.
* `instance_name` - (Optional, ForceNew) Name of the instance
* `instance_type` - (Required, ForceNew) EAIS instance type
* `resource_group_id` - (Optional, Computed, Available since v1.246.0) The ID of the resource group
* `security_group_id` - (Required, ForceNew) Security group ID
* `status` - (Optional, Computed) The status of the resource
* `tags` - (Optional, Map, Available since v1.246.0) The tags.
* `vswitch_id` - (Required, ForceNew) Switch ID.

### `environment_var`

The environment_var supports the following:
* `key` - (Optional, Available since v1.246.0) Keys for environment variables
* `value` - (Optional, Available since v1.246.0) Values of environment variables

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource
* `region_id` - Region ID

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Instance.
* `delete` - (Defaults to 5 mins) Used when delete the Instance.
* `update` - (Defaults to 5 mins) Used when update the Instance.

## Import

EAIS Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_eais_instance.example <id>
```
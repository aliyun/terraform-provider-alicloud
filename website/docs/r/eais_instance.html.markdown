---
subcategory: "Elastic Accelerated Computing Instances (EAIS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_eais_instance"
sidebar_current: "docs-alicloud-resource-eais-instance"
description: |-
  Provides a Alicloud Elastic Accelerated Computing Instances (EAIS) Instance resource.
---

# alicloud_eais_instance

Provides a Elastic Accelerated Computing Instances (EAIS) Instance resource.

For information about Elastic Accelerated Computing Instances (EAIS) Instance and how to use it, see [What is Instance](https://www.alibabacloud.com/help/en/resource-orchestration-service/latest/aliyun-eais-instance).

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

## Argument Reference

The following arguments are supported:

* `instance_type` - (Required, ForceNew) The type of the Instance. Valid values: `eais.ei-a6.4xlarge`, `eais.ei-a6.2xlarge`, `eais.ei-a6.xlarge`, `eais.ei-a6.large`, `eais.ei-a6.medium`.
* `vswitch_id` - (Required, ForceNew) The ID of the vSwitch.
* `security_group_id` - (Required, ForceNew) The ID of the security group.
* `instance_name` - (Optional, ForceNew) The name of the Instance.
* `force` - (Optional, Bool) Specifies whether to force delete the Instance. Default value: `false`. Valid values:
  - `true`: Enable.
  - `false`: Disable.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Instance.
* `status` - The status of the Instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Instance.
* `delete` - (Defaults to 5 mins) Used when delete the Instance.

## Import

Elastic Accelerated Computing Instances (EAIS) Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_eais_instance.example <id>
```

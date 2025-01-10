---
subcategory: "APIG"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_environment"
description: |-
  Provides a Alicloud APIG Environment resource.
---

# alicloud_apig_environment

Provides a APIG Environment resource.



For information about APIG Environment and how to use it, see [What is Environment](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.240.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_apig_environment&exampleId=4aacdd1a-aacf-c947-083d-b3386f109a6255aa57c2&activeTab=example&spm=docs.r.apig_environment.0.4aacdd1aaa&intl_lang=EN_US" target="_blank">
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

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_apig_gateway" "defaultgateway" {
  network_access_config {
    type = "Intranet"
  }
  vswitch {
    vswitch_id = data.alicloud_vswitches.default.ids.0
  }
  zone_config {
    select_option = "Auto"
  }
  vpc {
    vpc_id = data.alicloud_vpcs.default.ids.0
  }
  payment_type = "PayAsYouGo"
  gateway_name = format("%s2", var.name)
  spec         = "apigw.small.x1"
  log_config {
    sls {
    }
  }
}

resource "alicloud_apig_environment" "default" {
  description       = var.name
  environment_name  = var.name
  gateway_id        = alicloud_apig_gateway.defaultgateway.id
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.1
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) Description
* `environment_name` - (Required, ForceNew) The name of the resource
* `gateway_id` - (Required, ForceNew) Gateway id
* `resource_group_id` - (Optional, Computed) The ID of the resource group

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Environment.
* `delete` - (Defaults to 5 mins) Used when delete the Environment.
* `update` - (Defaults to 5 mins) Used when update the Environment.

## Import

APIG Environment can be imported using the id, e.g.

```shell
$ terraform import alicloud_apig_environment.example <id>
```
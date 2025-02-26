---
subcategory: "Api Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_instance"
description: |-
  Provides a Alicloud Api Gateway Instance resource.
---

# alicloud_api_gateway_instance

Provides a Api Gateway Instance resource. 

For information about Api Gateway Instance and how to use it, see [What is Instance](https://www.alibabacloud.com/help/en/api-gateway/product-overview/dedicated-instances).

-> **NOTE:** Available since v1.218.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_api_gateway_instance&exampleId=d73215ff-a60f-46bb-9711-20b8002b1ac5eca8eb10&activeTab=example&spm=docs.r.api_gateway_instance.0.d73215ffa6&intl_lang=EN_US" target="_blank">
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


resource "alicloud_api_gateway_instance" "default" {
  instance_name = var.name

  instance_spec = "api.s1.small"
  https_policy  = "HTTPS2_TLS1_0"
  zone_id       = "cn-hangzhou-MAZ6(i,j,k)"
  payment_type  = "PayAsYouGo"
  instance_type = "normal"
}
```

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_api_gateway_instance&exampleId=75a35892-c685-2e89-d544-e3af9724a154f64da7af&activeTab=example&spm=docs.r.api_gateway_instance.1.75a35892c6&intl_lang=EN_US" target="_blank">
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

resource "alicloud_vpc" "vpc" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "vswitch_1" {
  vpc_id       = alicloud_vpc.vpc.id
  cidr_block   = "172.16.0.0/16"
  zone_id      = "cn-hangzhou-j"
  vswitch_name = "${var.name}_1"
}

resource "alicloud_vswitch" "vswitch_2" {
  vpc_id       = alicloud_vpc.vpc.id
  cidr_block   = "172.17.0.0/16"
  zone_id      = "cn-hangzhou-k"
  vswitch_name = "${var.name}_2"
}

resource "alicloud_security_group" "security_group" {
  vpc_id              = alicloud_vpc.vpc.id
  security_group_name = var.name
}

resource "alicloud_api_gateway_instance" "vpc_integration_instance" {
  instance_name = var.name
  https_policy  = "HTTPS2_TLS1_0"
  instance_spec = "api.s1.small"
  instance_type = "vpc_connect"
  payment_type  = "PayAsYouGo"
  user_vpc_id   = alicloud_vpc.vpc.id
  instance_cidr = "192.168.0.0/16"
  zone_vswitch_security_group {
    zone_id        = alicloud_vswitch.vswitch_1.zone_id
    vswitch_id     = alicloud_vswitch.vswitch_1.id
    cidr_block     = alicloud_vswitch.vswitch_1.cidr_block
    security_group = alicloud_security_group.security_group.id
  }
  zone_vswitch_security_group {
    zone_id        = alicloud_vswitch.vswitch_2.zone_id
    vswitch_id     = alicloud_vswitch.vswitch_2.id
    cidr_block     = alicloud_vswitch.vswitch_2.cidr_block
    security_group = alicloud_security_group.security_group.id
  }
}
```

## Argument Reference

The following arguments are supported:
* `duration` - (Optional) The time of the instance package. Valid values:
  - PricingCycle is **Month**, indicating monthly payment. The value range is **1** to **9**.
  - PricingCycle is **Year**, indicating annual payment. The value range is **1** to **3**.

When the value of> ChargeType is **PrePaid**, this parameter is available and must be passed in.
* `egress_ipv6_enable` - (Optional) Specifies whether IPv6 egress capability is enabled.
* `https_policy` - (Required) Https policy.
* `instance_name` - (Required) Instance name.
* `instance_spec` - (Required, ForceNew) Instance type.
* `instance_type` - (Optional, ForceNew, Computed) The type of the instance. Valid values are:
  - `normal`: traditional dedicated instance (default).
  - `vpc_connect`: Vpc integration instance. When this type is selected, `instance_cidr`, `user_vpc_id` and `zone_vswitch_security_group` must be specified.
* `payment_type` - (Required, ForceNew) The payment type of the resource.
* `pricing_cycle` - (Optional) The subscription instance is of the subscription year or month type. This parameter is required when the Payment type is PrePaid. The value range is as follows:
  - `year`: year.
  - `month`: month.
* `ipv6_enabled` - (Optional, Available since v1.228.0) Specifies whether IPv6 ingress capability is enabled.
* `user_vpc_id` - (Optional, ForceNew) User's VpcID.
* `vpc_slb_intranet_enable` - (Optional) Whether the slb of the Vpc supports.
* `zone_id` - (Optional, ForceNew) The zone where the instance is deployed.
* `instance_cidr` - (Optional, ForceNew, Available since v1.228.0) The CIDR block for the instance deployment. Valid values are:
  - `192.168.0.0/16`.
  - `172.16.0.0/12`.
* `zone_vswitch_security_group` - (Optional, ForceNew, Available since v1.228.0) Network configuration details for Vpc integration instance which includes the availability zone, VSwitch, and security group information. See [`zone_vswitch_security_group`](#zone_vswitch_security_group) below.
* `to_connect_vpc_ip_block` - (Optional, Available since v1.228.0) The additional IP block that the VPC integration instance can access, conflict with `delete_vpc_ip_block`. See [`to_connect_vpc_ip_block`](#to_connect_vpc_ip_block) below.
* `delete_vpc_ip_block` - (Optional, Available since v1.228.0) Indicates whether to delete the IP block that the VPC can access, conflict with `to_connect_vpc_ip_block`.
* `skip_wait_switch` - (Optional, Available since v1.244.0) Specifies whether to skip the WAIT_SWITCH status of instance when modifying instance spec. Works only when instance spec change.

### `zone_vswitch_security_group`

The zone_vswitch_security_group supports the following:
* `zone_id` - (Required) The zone ID.
* `vswitch_id` - (Required) The VSwitch ID.
* `cidr_block` - (Required) The CIDR block of the VSwitch.
* `security_group` - (Required) The ID of the security group.

### `to_connect_vpc_ip_block`

The to_connect_vpc_ip_block supports the following:

* `zone_id` - (Optional) The zone ID.
* `vswitch_id` - (Optional) The VSwitch ID.
* `cidr_block` - (Required) The CIDR block of the VSwitch.
* `customized` - (Optional) Specifies whether the IP block is customized.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Creation time.
* `status` - The status of the resource.
* `support_ipv6` - Does ipv6 support.
* `connect_cidr_blocks` - (Available since v1.228.0) The CIDR blocks that can be accessed by the Vpc integration instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 25 mins) Used when create the Instance.
* `delete` - (Defaults to 25 mins) Used when delete the Instance.
* `update` - (Defaults to 25 mins) Used when update the Instance.

## Import

Api Gateway Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_api_gateway_instance.example <id>
```
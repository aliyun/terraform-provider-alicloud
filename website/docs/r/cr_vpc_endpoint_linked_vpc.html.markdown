---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_vpc_endpoint_linked_vpc"
sidebar_current: "docs-alicloud-resource-cr-vpc-endpoint-linked-vpc"
description: |-
  Provides a Alicloud CR Vpc Endpoint Linked Vpc resource.
---

# alicloud_cr_vpc_endpoint_linked_vpc

Provides a CR Vpc Endpoint Linked Vpc resource.

For information about CR Vpc Endpoint Linked Vpc and how to use it, see [What is Vpc Endpoint Linked Vpc](https://www.alibabacloud.com/help/en/acr/developer-reference/api-cr-2018-12-01-createinstancevpcendpointlinkedvpc).

-> **NOTE:** Available since v1.199.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_cr_vpc_endpoint_linked_vpc&exampleId=3bfb01d8-fad7-5567-3c5e-03e4e32c83115d54780d&activeTab=example&spm=docs.r.cr_vpc_endpoint_linked_vpc.0.3bfb01d8fa&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

resource "random_integer" "default" {
  min = 100000
  max = 999999
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
resource "alicloud_vpc" "default" {
  vpc_name   = "${var.name}-${random_integer.default.result}"
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = "${var.name}-${random_integer.default.result}"
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_cr_ee_instance" "default" {
  payment_type   = "Subscription"
  period         = 1
  renew_period   = 0
  renewal_status = "ManualRenewal"
  instance_type  = "Advanced"
  instance_name  = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_cr_vpc_endpoint_linked_vpc" "default" {
  instance_id                      = alicloud_cr_ee_instance.default.id
  vpc_id                           = alicloud_vpc.default.id
  vswitch_id                       = alicloud_vswitch.default.id
  module_name                      = "Registry"
  enable_create_dns_record_in_pvzt = true
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the instance.
* `vpc_id` - (Required, ForceNew) The ID of the VPC.
* `vswitch_id` - (Required, ForceNew) The ID of the vSwitch.
* `module_name` - (Required, ForceNew) The name of the module that you want to access. Valid Values:
  - `Registry`: the image repository.
  - `Chart`: a Helm chart.
* `enable_create_dns_record_in_pvzt` - (Optional) Specifies whether to automatically create an Alibaba Cloud DNS PrivateZone record. Valid Values:
  - `true`: automatically creates a PrivateZone record.
  - `false`: does not automatically create a PrivateZone record.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Vpc Endpoint Linked Vpc. It formats as `<instance_id>:<vpc_id>:<vswitch_id>:<module_name>`.
* `status` - The status of the Vpc Endpoint Linked Vpc.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Vpc Endpoint Linked Vpc.
* `delete` - (Defaults to 3 mins) Used when delete the Vpc Endpoint Linked Vpc.

## Import

CR Vpc Endpoint Linked Vpc can be imported using the id, e.g.

```shell
$ terraform import alicloud_cr_vpc_endpoint_linked_vpc.example <instance_id>:<vpc_id>:<vswitch_id>:<module_name>
```

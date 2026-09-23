---
subcategory: "Serverless App Engine (SAE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sae_load_balancer_internet"
sidebar_current: "docs-alicloud-resource-sae-load-balancer-internet"
description: |-
  Provides an Alicloud Serverless App Engine (SAE) Application Load Balancer Attachment resource.
---

# alicloud_sae_load_balancer_internet

Provides an Alicloud Serverless App Engine (SAE) Application Load Balancer Attachment resource.

For information about Serverless App Engine (SAE) Load Balancer Internet Attachment and how to use it, see [alicloud_sae_load_balancer_internet](https://www.alibabacloud.com/help/en/sae/latest/bindslb).

-> **NOTE:** Available since v1.164.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_sae_load_balancer_internet&exampleId=8f97dfeb-c53f-f371-89ea-575446075bb992b5dbf6&activeTab=example&spm=docs.r.sae_load_balancer_internet.0.8f97dfebc5&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tf-example"
}
data "alicloud_regions" "default" {
  current = true
}
resource "random_integer" "default" {
  max = 99999
  min = 10000
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}
resource "alicloud_security_group" "default" {
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_sae_namespace" "default" {
  namespace_id              = "${data.alicloud_regions.default.regions.0.id}:example${random_integer.default.result}"
  namespace_name            = var.name
  namespace_description     = var.name
  enable_micro_registration = false
}

resource "alicloud_sae_application" "default" {
  app_description    = var.name
  app_name           = "${var.name}-${random_integer.default.result}"
  namespace_id       = alicloud_sae_namespace.default.id
  image_url          = "registry-vpc.cn-hangzhou.aliyuncs.com/lxepoo/apache-php5"
  package_type       = "Image"
  jdk                = "Open JDK 8"
  security_group_id  = alicloud_security_group.default.id
  vpc_id             = alicloud_vpc.default.id
  vswitch_id         = alicloud_vswitch.default.id
  timezone           = "Asia/Beijing"
  replicas           = "5"
  cpu                = "500"
  memory             = "2048"
  micro_registration = "0"
}

resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name = var.name
  vswitch_id         = alicloud_vswitch.default.id
  load_balancer_spec = "slb.s2.small"
  address_type       = "internet"
}

resource "alicloud_sae_load_balancer_internet" "default" {
  app_id          = alicloud_sae_application.default.id
  internet_slb_id = alicloud_slb_load_balancer.default.id
  internet {
    protocol    = "TCP"
    port        = 80
    target_port = 8080
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_sae_load_balancer_internet&spm=docs.r.sae_load_balancer_internet.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `app_id` - (Required) The target application ID that needs to be bound to the SLB.
* `internet_slb_id` - (Optional) The internet SLB ID.
* `internet` - (Required) The bound private network SLB. See [`internet`](#internet) below.

### `internet`

The internet supports the following:

* `protocol` - (Optional) The Network protocol. Valid values: `TCP` ,`HTTP`,`HTTPS`.
* `https_cert_id` - (Optional) The SSL certificate. `https_cert_id` is required when HTTPS is selected
* `target_port` - (Optional) The Container port.
* `port` - (Optional) The SLB Port.


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID. The value is the same as the application ID.
* `internet_ip` - Use designated public network SLBs that have been purchased to support non-shared instances.

## Import

The resource can be imported using the id, e.g.

```shell
$ terraform import alicloud_sae_load_balancer_internet.example <id>
```

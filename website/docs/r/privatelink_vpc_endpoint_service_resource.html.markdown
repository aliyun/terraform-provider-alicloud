---
subcategory: "Private Link"
layout: "alicloud"
page_title: "Alicloud: alicloud_privatelink_vpc_endpoint_service_resource"
description: |-
  Provides a Alicloud Private Link Vpc Endpoint Service Resource resource.
---

# alicloud_privatelink_vpc_endpoint_service_resource

Provides a Private Link Vpc Endpoint Service Resource resource.

Endpoint service resource.

For information about Private Link Vpc Endpoint Service Resource and how to use it, see [What is Vpc Endpoint Service Resource](https://www.alibabacloud.com/help/en/privatelink/latest/api-privatelink-2020-04-15-attachresourcetovpcendpointservice).

-> **NOTE:** Available since v1.110.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_privatelink_vpc_endpoint_service_resource&exampleId=23cc3b5a-f6fc-adf3-dbfa-17aab8657574c6be7cac&activeTab=example&spm=docs.r.privatelink_vpc_endpoint_service_resource.0.23cc3b5af6&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}

data "alicloud_zones" "example" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_privatelink_vpc_endpoint_service" "example" {
  service_description    = var.name
  connect_bandwidth      = 103
  auto_accept_connection = false
}

resource "alicloud_vpc" "example" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}
resource "alicloud_vswitch" "example" {
  vswitch_name = var.name
  cidr_block   = "10.1.0.0/16"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = data.alicloud_zones.example.zones.0.id
}

resource "alicloud_security_group" "example" {
  security_group_name = var.name
  vpc_id              = alicloud_vpc.example.id
}

resource "alicloud_slb_load_balancer" "example" {
  load_balancer_name = var.name
  vswitch_id         = alicloud_vswitch.example.id
  load_balancer_spec = "slb.s2.small"
  address_type       = "intranet"
}

resource "alicloud_privatelink_vpc_endpoint" "example" {
  service_id         = alicloud_privatelink_vpc_endpoint_service.example.id
  security_group_ids = [alicloud_security_group.example.id]
  vpc_id             = alicloud_vpc.example.id
  vpc_endpoint_name  = var.name
}
resource "alicloud_privatelink_vpc_endpoint_service_resource" "example" {
  service_id    = alicloud_privatelink_vpc_endpoint_service.example.id
  resource_id   = alicloud_slb_load_balancer.example.id
  resource_type = "slb"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_privatelink_vpc_endpoint_service_resource&spm=docs.r.privatelink_vpc_endpoint_service_resource.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `dry_run` - (Optional) Specifies whether to perform only a dry run, without performing the actual request. Valid values:
  - `true`: performs only a dry run. The system checks the request for potential issues, including missing parameter values, incorrect request syntax, and service limits. If the request fails the dry run, an error code is returned. If the request passes the dry run, the DryRunOperation error code is returned.
  - **false (default)**: performs a dry run and performs the actual request. If the request passes the dry run, a 2xx HTTP status code is returned and the operation is performed.
* `resource_id` - (Required, ForceNew) The service resource ID.
* `resource_type` - (Required, ForceNew) Service resource type, value:
  - `slb`: indicates that the service resource type is Classic Load Balancer (CLB).
  - `alb`: indicates that the service resource type is Application Load Balancer (ALB).
  - `nlb`: indicates that the service resource type is Network Load Balancer (NLB).
* `service_id` - (Required, ForceNew) The endpoint service ID.
* `zone_id` - (Optional, ForceNew, Computed, Available since v1.212.0) The ID of the zone to which the service resource belongs. (valid when the resource type is nlb/alb).

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<service_id>:<resource_id>:<zone_id>`.
* `region_id` - (Available since v1.235.0) The ID of the region where the service resource is deployed.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vpc Endpoint Service Resource.
* `delete` - (Defaults to 5 mins) Used when delete the Vpc Endpoint Service Resource.

## Import

Private Link Vpc Endpoint Service Resource can be imported using the id, e.g.

```shell
$ terraform import alicloud_privatelink_vpc_endpoint_service_resource.example <service_id>:<resource_id>:<zone_id>
```
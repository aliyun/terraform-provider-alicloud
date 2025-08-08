---
subcategory: "Private Link"
layout: "alicloud"
page_title: "Alicloud: alicloud_privatelink_vpc_endpoint_service"
description: |-
  Provides a Alicloud Private Link Vpc Endpoint Service resource.
---

# alicloud_privatelink_vpc_endpoint_service

Provides a Private Link Vpc Endpoint Service resource.



For information about Private Link Vpc Endpoint Service and how to use it, see [What is Vpc Endpoint Service](https://www.alibabacloud.com/help/en/privatelink/latest/api-privatelink-2020-04-15-createvpcendpointservice).

-> **NOTE:** Available since v1.109.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_privatelink_vpc_endpoint_service&exampleId=3a7ae46a-09f6-78a6-a810-69e7b0912f50c5749917&activeTab=example&spm=docs.r.privatelink_vpc_endpoint_service.0.3a7ae46a09&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}

resource "alicloud_privatelink_vpc_endpoint_service" "example" {
  service_description    = var.name
  connect_bandwidth      = 103
  auto_accept_connection = false
}
```

## Argument Reference

The following arguments are supported:
* `address_ip_version` - (Optional, Computed, Available since v1.239.0) The IP address version.
* `auto_accept_connection` - (Optional) Indicates whether the endpoint service automatically accepts endpoint connection requests. Valid values:
  - `true`
  - `false`
* `connect_bandwidth` - (Optional, Computed, Int) The default bandwidth of the endpoint connection. Valid values: 100 to 10240. Unit: Mbit/s.
* `dry_run` - (Optional) Specifies whether to perform only a dry run, without performing the actual request.
  - `true`: performs only a dry run. The system checks the request for potential issues, including missing parameter values, incorrect request syntax, and service limits. If the request fails the dry run, an error message is returned. If the request passes the dry run, the DryRunOperation error code is returned.
  - **false (default)**: performs a dry run and performs the actual request. If the request passes the dry run, a 2xx HTTP status code is returned and the operation is performed.
* `payer` - (Optional, ForceNew, Computed) The payer of the endpoint service. Valid values:
  - `Endpoint`: the service consumer.
  - `EndpointService`: the service provider.
* `resource_group_id` - (Optional, Computed) The resource group ID.
* `service_description` - (Optional) The description of the endpoint service.
* `service_resource_type` - (Optional, ForceNew, Computed) The service resource type. Value:
  - `slb`: indicates that the service resource type is Classic Load Balancer (CLB).
  - `alb`: indicates that the service resource type is Application Load Balancer (ALB).
  - `nlb`: indicates that the service resource type is Network Load Balancer (NLB).
  - `gwlb`: indicates that the service resource type is Gateway Load Balancer (GWLB).
* `service_support_ipv6` - (Optional, Computed) Specifies whether to enable IPv6 for the endpoint service. Valid values:
  - `true`
  - **false (default)**
* `tags` - (Optional, Map) The list of tags.
* `zone_affinity_enabled` - (Optional, Computed) Specifies whether to first resolve the domain name of the nearest endpoint that is associated with the endpoint service. Valid values:
  - `true`
  - **false (default)**

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the endpoint service was created.
* `region_id` - The ID of the region to which the endpoint service belongs.
* `service_business_status` - The service state of the endpoint service. 
* `service_domain` - The domain name of the endpoint service.
* `status` - The state of the endpoint service. 
* `vpc_endpoint_service_name` - The name of the endpoint service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vpc Endpoint Service.
* `delete` - (Defaults to 5 mins) Used when delete the Vpc Endpoint Service.
* `update` - (Defaults to 5 mins) Used when update the Vpc Endpoint Service.

## Import

Private Link Vpc Endpoint Service can be imported using the id, e.g.

```shell
$ terraform import alicloud_privatelink_vpc_endpoint_service.example <id>
```
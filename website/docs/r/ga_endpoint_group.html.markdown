---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_endpoint_group"
sidebar_current: "docs-alicloud-resource-ga-endpoint-group"
description: |-
  Provides a Alicloud Global Accelerator (GA) Endpoint Group resource.
---

# alicloud_ga_endpoint_group

Provides a Global Accelerator (GA) Endpoint Group resource.

For information about Global Accelerator (GA) Endpoint Group and how to use it, see [What is Endpoint Group](https://www.alibabacloud.com/help/en/global-accelerator/latest/api-ga-2019-11-20-createendpointgroup).

-> **NOTE:** Available since v1.113.0.

-> **NOTE:** Listeners that use different protocols support different types of endpoint groups:

* For a TCP or UDP listener, you can create only one default endpoint group. 
* For an HTTP or HTTPS listener, you can create one default endpoint group and one virtual endpoint group. By default, you can create only one virtual endpoint group. 
  * A default endpoint group refers to the endpoint group that you configure when you create an HTTP or HTTPS listener. 
  * A virtual endpoint group refers to the endpoint group that you can create on the Endpoint Group page after you create a listener.
* After you create a virtual endpoint group for an HTTP or HTTPS listener, you can create a forwarding rule and associate the forwarding rule with the virtual endpoint group. Then, the HTTP or HTTPS listener forwards requests with different destination domain names or paths to the default or virtual endpoint group based on the forwarding rule. This way, you can use one Global Accelerator (GA) instance to accelerate access to multiple domain names or paths. For more information about how to create a forwarding rule, see [Manage forwarding rules](https://www.alibabacloud.com/help/en/doc-detail/204224.htm).

## Example Usage

Basic Usage

```terraform
variable "region" {
  default = "cn-hangzhou"
}

provider "alicloud" {
  region  = var.region
  profile = "default"
}

resource "alicloud_ga_accelerator" "default" {
  duration        = 1
  auto_use_coupon = true
  spec            = "1"
}

resource "alicloud_ga_bandwidth_package" "default" {
  bandwidth      = 100
  type           = "Basic"
  bandwidth_type = "Basic"
  payment_type   = "PayAsYouGo"
  billing_type   = "PayBy95"
  ratio          = 30
}

resource "alicloud_ga_bandwidth_package_attachment" "default" {
  accelerator_id       = alicloud_ga_accelerator.default.id
  bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
}

resource "alicloud_ga_listener" "default" {
  accelerator_id = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  port_ranges {
    from_port = 60
    to_port   = 70
  }
  client_affinity = "SOURCE_IP"
  protocol        = "UDP"
  name            = "terraform-example"
}

resource "alicloud_eip_address" "default" {
  count                = 2
  bandwidth            = "10"
  internet_charge_type = "PayByBandwidth"
  address_name         = "terraform-example"
}

resource "alicloud_ga_endpoint_group" "default" {
  accelerator_id = alicloud_ga_accelerator.default.id
  endpoint_configurations {
    endpoint = alicloud_eip_address.default.0.ip_address
    type     = "PublicIp"
    weight   = "20"
  }
  endpoint_configurations {
    endpoint = alicloud_eip_address.default.1.ip_address
    type     = "PublicIp"
    weight   = "20"
  }
  endpoint_group_region = var.region
  listener_id           = alicloud_ga_listener.default.id
}
```

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Required, ForceNew) The ID of the Global Accelerator instance to which the endpoint group will be added.
* `listener_id` - (Required, ForceNew) The ID of the listener that is associated with the endpoint group.
* `endpoint_group_region` - (Required, ForceNew) The ID of the region where the endpoint group is deployed.
* `endpoint_group_type` - (Optional, ForceNew) The endpoint group type. Default value: `default`. Valid values: `default`, `virtual`.
-> **NOTE:** Only the listening instance of HTTP or HTTPS protocol supports the creation of virtual terminal node group.
* `endpoint_request_protocol` - (Optional) The endpoint request protocol. Valid values: `HTTP`, `HTTPS`.
-> **NOTE:** This item is only supported when creating terminal node group for listening instance of HTTP or HTTPS protocol. For the listening instance of HTTP protocol, the back-end service protocol supports and only supports HTTP.
* `health_check_interval_seconds` - (Optional, Int) The interval between two consecutive health checks. Unit: seconds.
* `health_check_path` - (Optional) The path specified as the destination of the targets for health checks.
* `health_check_port` - (Optional, Int) The port that is used for health checks.
* `health_check_protocol` - (Optional) The protocol that is used to connect to the targets for health checks. Valid values: `http`, `https`, `tcp`.
* `threshold_count` - (Optional, Int) The number of consecutive failed heath checks that must occur before the endpoint is deemed unhealthy. Default value: `3`.
* `traffic_percentage` - (Optional, Int) The weight of the endpoint group when the corresponding listener is associated with multiple endpoint groups.
* `name` - (Optional) The name of the endpoint group.
* `description` - (Optional) The description of the endpoint group.
* `endpoint_configurations` - (Required, Set) The endpointConfigurations of the endpoint group. See [`endpoint_configurations`](#endpoint_configurations) below.
* `port_overrides` - (Optional, Set) Mapping between listening port and forwarding port of boarding point. See [`port_overrides`](#port_overrides) below.
-> **NOTE:** Port mapping is only supported when creating terminal node group for listening instance of HTTP or HTTPS protocol. The listening port in the port map must be consistent with the listening port of the current listening instance.
* `tags` - (Optional, Available since v1.207.1) A mapping of tags to assign to the resource.

### `port_overrides`

The port_overrides supports the following: 

* `endpoint_port` - (Optional, Int) Forwarding port.
* `listener_port` - (Optional, Int) Listener port.

### `endpoint_configurations`

The endpoint_configurations supports the following: 

* `endpoint` - (Required) The IP address or domain name of Endpoint N in the endpoint group.
* `type` - (Required) The type of Endpoint N in the endpoint group. Valid values:
  - `Domain`: a custom domain name.
  - `Ip`: a custom IP address.
  - `PublicIp`: an Alibaba Cloud public IP address.
  - `ECS`: an Alibaba Cloud Elastic Compute Service (ECS) instance.
  - `SLB`: an Alibaba Cloud Server Load Balancer (SLB) instance.
-> **NOTE:** When the terminal node type is ECS or SLB, if the service association role does not exist, the system will automatically create a service association role named aliyunserviceroleforgavpcndpoint.
* `weight` - (Required, Int) The weight of Endpoint N in the endpoint group. Valid values: `0` to `255`.
-> **NOTE:** If the weight of a terminal node is set to 0, global acceleration will terminate the distribution of traffic to the terminal node. Please be careful.
* `enable_proxy_protocol` - (Optional, Bool, Available since v1.207.1) Specifies whether to preserve client IP addresses by using the ProxyProtocol module. Default Value: `false`. Valid values:
  - `true`: preserves client IP addresses by using the ProxyProtocol module.
  - `false`: does not preserve client IP addresses by using the ProxyProtocol module.
* `enable_clientip_preservation` - (Optional, Bool) Indicates whether client IP addresses are reserved. Default Value: `false`. Valid values:
  - `true`: Client IP addresses are reserved.
  - `false`: Client IP addresses are not reserved.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Endpoint Group.
* `endpoint_group_ip_list` - (Available since v1.213.0) The active endpoint IP addresses of the endpoint group.
* `status` - The status of the endpoint group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 15 mins) Used when create the Endpoint Group.
* `update` - (Defaults to 2 mins) Used when update the Endpoint Group.
* `delete` - (Defaults to 10 mins) Used when delete the Endpoint Group.

## Import

Ga Endpoint Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_endpoint_group.example <id>
```

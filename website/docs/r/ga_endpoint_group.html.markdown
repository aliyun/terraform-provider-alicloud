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

-> **NOTE:** Listeners that use different protocols support different types of Endpoint Groups:

* For a `TCP` listener, if you want to create a `virtual` Endpoint Group, please ensure that the `default` Endpoint Group of the same region has been created.
* For a `UDP` listener, you can only create `default` Endpoint Group.
* For an `HTTP` or `HTTPS` listener, you can create one `default` Endpoint Group and multiple `virtual` Endpoint Group.
* After you create a `virtual` endpoint group for an `HTTP` or `HTTPS` listener, you can create a forwarding rule and associate the forwarding rule with the `virtual` endpoint group. Then, the `HTTP` or `HTTPS` listener forwards requests with different destination domain names or paths to the `default` or `virtual` Endpoint Group based on the forwarding rule. This way, you can use one Global Accelerator (GA) instance to accelerate access to multiple domain names or paths. For more information about how to create a forwarding rule, see [Manage forwarding rules](https://www.alibabacloud.com/help/en/doc-detail/204224.htm).

-> **WARN:** There is a serious bug in the `traffic_percentage` of the `alicloud_ga_endpoint_group` before version 1.211.1, while the value of `traffic_percentage` has not been explicitly specified in the Terraform code, Terraform will set `traffic_percentage` to `0`. This behavior will cause your instance traffic to drop to zero. So, please use provider greater than or equal to version `1.211.1`.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ga_endpoint_group&exampleId=698d9c61-9f7d-dbd4-38cb-127f54151ac5cf3814d3&activeTab=example&spm=docs.r.ga_endpoint_group.0.698d9c619f&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
-> **NOTE:** Currently, only `HTTP` or `HTTPS` protocol listener can directly create a `virtual` Endpoint Group. If it is `TCP` protocol listener, and you want to create a `virtual` Endpoint Group, please ensure that the `default` Endpoint Group has been created.
* `endpoint_request_protocol` - (Optional) The protocol that is used by the backend server. Valid values: `HTTP`, `HTTPS`.
-> **NOTE:** `endpoint_request_protocol` can be specified only if the listener that is associated with the endpoint group uses `HTTP` or `HTTPS`. For the listener of `HTTP` protocol, `endpoint_request_protocol` can only be set to `HTTP`.
* `endpoint_protocol_version` - (Optional, Available since v1.230.1) The backend service protocol of the endpoint that is associated with the intelligent routing listener. Valid values: `HTTP1.1`, `HTTP2`.
-> **NOTE:** `endpoint_protocol_version` is valid only when `endpoint_request_protocol` is set to `HTTPS`.
* `health_check_enabled` - (Optional, Bool, Available since v1.215.0) Specifies whether to enable the health check feature. Valid values:
  - `true`: Enables the health check feature.
  - `false`: Disables the health check feature.
* `health_check_path` - (Optional) The path specified as the destination of the targets for health checks.
* `health_check_port` - (Optional, Int) The port that is used for health checks.
* `health_check_protocol` - (Optional) The protocol that is used to connect to the targets for health checks. Valid values:
  - `TCP` or `tcp`: TCP protocol.
  - `HTTP` or `http`: HTTP protocol.
  - `HTTPS` or `https`: HTTPS protocol.
-> **NOTE:** From version 1.223.0, `health_check_protocol` can be set to `TCP`, `HTTP`, `HTTPS`.
* `health_check_interval_seconds` - (Optional, Int) The interval between two consecutive health checks. Unit: seconds.
* `threshold_count` - (Optional, Int) The number of consecutive failed heath checks that must occur before the endpoint is deemed unhealthy. Default value: `3`.
* `traffic_percentage` - (Optional, Int) The weight of the endpoint group when the corresponding listener is associated with multiple endpoint groups.
* `name` - (Optional) The name of the endpoint group.
* `description` - (Optional) The description of the endpoint group.
* `endpoint_configurations` - (Required, Set) The endpointConfigurations of the endpoint group. See [`endpoint_configurations`](#endpoint_configurations) below.
* `port_overrides` - (Optional, Set) Mapping between listening port and forwarding port of boarding point. See [`port_overrides`](#port_overrides) below.
-> **NOTE:** Port mapping is only supported when creating terminal node group for listening instance of HTTP or HTTPS protocol. The listening port in the port map must be consistent with the listening port of the current listening instance.
* `tags` - (Optional, Available since v1.207.1) A mapping of tags to assign to the resource.

### `endpoint_configurations`

The endpoint_configurations supports the following:

* `endpoint` - (Required) The IP address or domain name of Endpoint N in the endpoint group.
* `type` - (Required) The type of Endpoint N in the endpoint group. Valid values:
  - `Domain`: A custom domain name.
  - `Ip`: A custom IP address.
  - `IpTarget`: (Available since v1.262.0) An Alibaba Cloud public IP address.
  - `PublicIp`: An Alibaba Cloud public IP address.
  - `ECS`: An Elastic Compute Service (ECS) instance.
  - `SLB`: A Classic Load Balancer (CLB) instance.
  - `ALB`: (Available since v1.232.0) An Application Load Balancer (ALB) instance.
  - `NLB`: (Available since v1.232.0) A Network Load Balancer (NLB) instance.
  - `ENI`: (Available since v1.232.0) An Elastic Network Interface (ENI).
  - `OSS`: (Available since v1.232.0) An Object Storage Service (OSS) bucket.
* `weight` - (Required, Int) The weight of Endpoint N in the endpoint group. Valid values: `0` to `255`.
-> **NOTE:** If the weight of a terminal node is set to `0`, global acceleration will terminate the distribution of traffic to the terminal node. Please be careful.
* `sub_address` - (Optional, Available since v1.232.0) The private IP address of the ENI.
-> **NOTE:** `sub_address` is valid only when `type` is set to `ENI`.
* `enable_proxy_protocol` - (Optional, Bool, Available since v1.207.1) Specifies whether to preserve client IP addresses by using the ProxyProtocol module. Default Value: `false`. Valid values:
  - `true`: preserves client IP addresses by using the ProxyProtocol module.
  - `false`: does not preserve client IP addresses by using the ProxyProtocol module.
* `enable_clientip_preservation` - (Optional, Bool) Indicates whether client IP addresses are reserved. Default Value: `false`. Valid values:
  - `true`: Client IP addresses are reserved.
  - `false`: Client IP addresses are not reserved.
* `vpc_id` - (Optional, Available since v1.262.0) The ID of the VPC.
* `vswitch_ids` - (Optional, List, Available since v1.262.0) The IDs of vSwitches that are deployed in the VPC.

### `port_overrides`

The port_overrides supports the following: 

* `endpoint_port` - (Optional, Int) Forwarding port.
* `listener_port` - (Optional, Int) Listener port.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Endpoint Group.
* `endpoint_group_ip_list` - (Available since v1.213.0) The active endpoint IP addresses of the endpoint group. `endpoint_group_ip_list` will change with the growth of network traffic. You can run `terraform apply` to query the latest CIDR blocks and IP addresses.
* `status` - The status of the endpoint group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 15 mins) Used when create the Endpoint Group.
* `update` - (Defaults to 3 mins) Used when update the Endpoint Group.
* `delete` - (Defaults to 10 mins) Used when delete the Endpoint Group.

## Import

Ga Endpoint Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_endpoint_group.example <id>
```

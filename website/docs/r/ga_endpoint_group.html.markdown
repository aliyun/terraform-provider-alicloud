---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_endpoint_group"
sidebar_current: "docs-alicloud-resource-ga-endpoint-group"
description: |-
  Provides a Alicloud Global Accelerator (GA) Endpoint Group resource.
---

# alicloud\_ga\_endpoint\_group

Provides a Global Accelerator (GA) Endpoint Group resource.

For information about Global Accelerator (GA) Endpoint Group and how to use it, see [What is Endpoint Group](https://www.alibabacloud.com/help/en/doc-detail/153259.htm).

-> **NOTE:** Available in v1.113.0+.

-> **NOTE:** Listeners that use different protocols support different types of endpoint groups:
* For a TCP or UDP listener, you can create only one default endpoint group. 
* For an HTTP or HTTPS listener, you can create one default endpoint group and one virtual endpoint group. By default, you can create only one virtual endpoint group. 
  * A default endpoint group refers to the endpoint group that you configure when you create an HTTP or HTTPS listener. 
  * A virtual endpoint group refers to the endpoint group that you can create on the Endpoint Group page after you create a listener.
* After you create a virtual endpoint group for an HTTP or HTTPS listener, you can create a forwarding rule and associate the forwarding rule with the virtual endpoint group. Then, the HTTP or HTTPS listener forwards requests with different destination domain names or paths to the default or virtual endpoint group based on the forwarding rule. This way, you can use one Global Accelerator (GA) instance to accelerate access to multiple domain names or paths. For more information about how to create a forwarding rule, see [Manage forwarding rules](https://www.alibabacloud.com/help/en/doc-detail/204224.htm).

## Example Usage

Basic Usage

```terraform
resource "alicloud_ga_accelerator" "example" {
  duration        = 1
  auto_use_coupon = true
  spec            = "1"
}
resource "alicloud_ga_bandwidth_package" "de" {
  bandwidth      = "100"
  type           = "Basic"
  bandwidth_type = "Basic"
  payment_type   = "PayAsYouGo"
  billing_type   = "PayBy95"
  ratio          = 30
}
resource "alicloud_ga_bandwidth_package_attachment" "de" {
  accelerator_id       = alicloud_ga_accelerator.example.id
  bandwidth_package_id = alicloud_ga_bandwidth_package.de.id
}
resource "alicloud_ga_listener" "example" {
  depends_on     = [alicloud_ga_bandwidth_package_attachment.de]
  accelerator_id = alicloud_ga_accelerator.example.id
  port_ranges {
    from_port = 60
    to_port   = 70
  }
}
resource "alicloud_eip_address" "example" {
  bandwidth            = "10"
  internet_charge_type = "PayByBandwidth"
}
resource "alicloud_ga_endpoint_group" "example" {
  accelerator_id = alicloud_ga_accelerator.example.id
  endpoint_configurations {
    endpoint = alicloud_eip_address.example.ip_address
    type     = "PublicIp"
    weight   = "20"
  }
  endpoint_group_region = "cn-hangzhou"
  listener_id           = alicloud_ga_listener.example.id
}

```

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Required) The ID of the Global Accelerator instance to which the endpoint group will be added.
* `description` - (Optional) The description of the endpoint group.
* `endpoint_configurations` - (Required) The endpointConfigurations of the endpoint group.
* `endpoint_group_region` - (ForceNew, Required) The ID of the region where the endpoint group is deployed.
* `endpoint_group_type` - (Optional, ForceNew) The endpoint group type. Valid values: `default`, `virtual`. Default value is `default`.

-> **NOTE:** Only the listening instance of HTTP or HTTPS protocol supports the creation of virtual terminal node group.
    
* `endpoint_request_protocol` - (Optional) The endpoint request protocol. Valid value: `HTTP`, `HTTPS`.

-> **NOTE:** This item is only supported when creating terminal node group for listening instance of HTTP or HTTPS protocol. For the listening instance of HTTP protocol, the back-end service protocol supports and only supports HTTP.

* `health_check_interval_seconds` - (Optional) The interval between two consecutive health checks. Unit: seconds.
* `health_check_path` - (Optional) The path specified as the destination of the targets for health checks.
* `health_check_port` - (Optional) The port that is used for health checks.
* `health_check_protocol` - (Optional) The protocol that is used to connect to the targets for health checks. Valid values: `http`, `https`, `tcp`.
* `listener_id` - (Required, ForceNew) The ID of the listener that is associated with the endpoint group.
* `name` - (Optional) The name of the endpoint group.
* `port_overrides` - (Optional) Mapping between listening port and forwarding port of boarding point.

-> **NOTE:** Port mapping is only supported when creating terminal node group for listening instance of HTTP or HTTPS protocol. The listening port in the port map must be consistent with the listening port of the current listening instance.

* `threshold_count` - (Optional) The number of consecutive failed heath checks that must occur before the endpoint is deemed unhealthy. Default value is `3`.
* `traffic_percentage` - (Optional) The weight of the endpoint group when the corresponding listener is associated with multiple endpoint groups.

#### Block port_overrides

The port_overrides supports the following: 

* `endpoint_port` - (Optional) Forwarding port.
* `listener_port` - (Optional) Listener port.

#### Block endpoint_configurations

The endpoint_configurations supports the following: 

* `enable_clientip_preservation` - (Optional) Indicates whether client IP addresses are reserved. Valid values: `true`: Client IP addresses are reserved, `false`: Client IP addresses are not reserved. Default value is `false`.
* `endpoint` - (Required) The IP address or domain name of Endpoint N in the endpoint group.
* `probe_port` - (Computed) Probe Port.
* `probe_protocol` - (Computed) Probe Protocol.
* `type` - (Required) The type of Endpoint N in the endpoint group. Valid values: `Domain`: a custom domain name, `Ip`: a custom IP address, `PublicIp`: an Alibaba Cloud public IP address, `ECS`: an Alibaba Cloud Elastic Compute Service (ECS) instance, `SLB`: an Alibaba Cloud Server Load Balancer (SLB) instance.

-> **NOTE:** When the terminal node type is ECS or SLB, if the service association role does not exist, the system will automatically create a service association role named aliyunserviceroleforgavpcndpoint.

* `weight` - (Required) The weight of Endpoint N in the endpoint group. Valid value is 0 to 255.

-> **NOTE:** If the weight of a terminal node is set to 0, global acceleration will terminate the distribution of traffic to the terminal node. Please be careful.
             
## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Endpoint Group.
* `status` - The status of the endpoint group.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the Endpoint Group.
* `delete` - (Defaults to 6 mins) Used when delete the Endpoint Group.
* `update` - (Defaults to 2 mins) Used when update the Endpoint Group.

## Import

Ga Endpoint Group can be imported using the id, e.g.

```
$ terraform import alicloud_ga_endpoint_group.example <id>
```

---
subcategory: "Private Link"
layout: "alicloud"
page_title: "Alicloud: alicloud_privatelink_vpc_endpoint"
description: |-
  Provides a Alicloud Private Link Vpc Endpoint resource.
---

# alicloud_privatelink_vpc_endpoint

Provides a Private Link Vpc Endpoint resource. 

For information about Private Link Vpc Endpoint and how to use it, see [What is Vpc Endpoint](https://www.alibabacloud.com/help/en/privatelink/latest/api-privatelink-2020-04-15-createvpcendpoint).

-> **NOTE:** Available since v1.109.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
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

resource "alicloud_security_group" "example" {
  name   = var.name
  vpc_id = alicloud_vpc.example.id
}

resource "alicloud_privatelink_vpc_endpoint" "example" {
  service_id         = alicloud_privatelink_vpc_endpoint_service.example.id
  security_group_ids = [alicloud_security_group.example.id]
  vpc_id             = alicloud_vpc.example.id
  vpc_endpoint_name  = var.name
}
```

## Argument Reference

The following arguments are supported:
* `dry_run` - (Optional) Specifies whether to perform only a dry run, without performing the actual request. Valid values:
  - **true**: performs only a dry run. The system checks the request for potential issues, including missing parameter values, incorrect request syntax, and service limits. If the request fails the dry run, an error message is returned. If the request passes the dry run, the DryRunOperation error code is returned.
  - **false (default)**: performs a dry run and performs the actual request. If the request passes the dry run, a 2xx HTTP status code is returned and the operation is performed.
* `endpoint_description` - (Optional) The description of the endpoint.
* `endpoint_type` - (Optional, ForceNew, Computed, Available since v1.212.0) The endpoint type.Only the value: Interface, indicating the Interface endpoint. You can add the service resource types of Application Load Balancer (ALB), Classic Load Balancer (CLB), and Network Load Balancer (NLB).
* `protected_enabled` - (Optional, Available since v1.212.0) Specifies whether to enable user authentication. This parameter is available in Security Token Service (STS) mode. Valid values:
  - **true**: enables user authentication. After user authentication is enabled, only the user who creates the endpoint can modify or delete the endpoint in STS mode.
  - **false (default)**: disables user authentication.
* `resource_group_id` - (Optional, Computed, Available since v1.212.0) The resource group ID.
* `security_group_ids` - (Required) The ID of the security group that is associated with the endpoint ENI. The security group can be used to control data transfer between the VPC and the endpoint ENI.The endpoint can be associated with up to 10 security groups.
* `service_id` - (Optional, ForceNew) The ID of the endpoint service with which the endpoint is associated.
* `service_name` - (Optional, ForceNew, Computed) The name of the endpoint service with which the endpoint is associated.
* `tags` - (Optional, Map, Available since v1.212.0) The list of tags.
* `vpc_endpoint_name` - (Optional) The name of the endpoint.
* `vpc_id` - (Required, ForceNew) The ID of the VPC to which the endpoint belongs.
* `zone_private_ip_address_count` - (Optional, ForceNew, Computed, Available since v1.212.0) The number of private IP addresses that are assigned to an elastic network interface (ENI) in each zone. Only 1 is returned.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `bandwidth` - The bandwidth of the endpoint connection.  1024 to 10240. Unit: Mbit/s.Note: The bandwidth of an endpoint connection is in the range of 100 to 10,240 Mbit/s. The default bandwidth is 1,024 Mbit/s. When the endpoint is connected to the endpoint service, the default bandwidth is the minimum bandwidth. In this case, the connection bandwidth range is 1,024 to 10,240 Mbit/s.
* `connection_status` - The state of the endpoint connection. 
* `create_time` - The time when the endpoint was created.
* `endpoint_business_status` - The service state of the endpoint. 
* `endpoint_domain` - The domain name of the endpoint.
* `status` - The state of the endpoint. 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vpc Endpoint.
* `delete` - (Defaults to 5 mins) Used when delete the Vpc Endpoint.
* `update` - (Defaults to 5 mins) Used when update the Vpc Endpoint.

## Import

Private Link Vpc Endpoint can be imported using the id, e.g.

```shell
$ terraform import alicloud_privatelink_vpc_endpoint.example <id>
```
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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_privatelink_vpc_endpoint&exampleId=5272b15f-709e-789e-0a7b-61727e7a83a1e65c68d6&activeTab=example&spm=docs.r.privatelink_vpc_endpoint.0.5272b15f70&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "ap-southeast-5"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultbFzA4a" {
  description = "example-terraform"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name
}

resource "alicloud_security_group" "default1FTFrP" {
  name   = var.name
  vpc_id = alicloud_vpc.defaultbFzA4a.id
}

resource "alicloud_security_group" "defaultjljY5S" {
  name   = var.name
  vpc_id = alicloud_vpc.defaultbFzA4a.id
}

resource "alicloud_privatelink_vpc_endpoint" "default" {
  endpoint_description          = var.name
  vpc_endpoint_name             = var.name
  resource_group_id             = data.alicloud_resource_manager_resource_groups.default.ids.0
  endpoint_type                 = "Interface"
  vpc_id                        = alicloud_vpc.defaultbFzA4a.id
  service_name                  = "com.aliyuncs.privatelink.ap-southeast-5.oss"
  dry_run                       = "false"
  zone_private_ip_address_count = "1"
  policy_document               = jsonencode({ "Version" : "1", "Statement" : [{ "Effect" : "Allow", "Action" : ["*"], "Resource" : ["*"], "Principal" : "*" }] })
  security_group_ids = [
    "${alicloud_security_group.default1FTFrP.id}"
  ]
  service_id        = "epsrv-k1apjysze8u1l9t6uyg9"
  protected_enabled = "false"
}
```

## Argument Reference

The following arguments are supported:
* `address_ip_version` - (Optional, Computed, Available since v1.239.0) The IP address version. Valid values:
  - `IPv4` (default): IPv4.
  - `DualStack`: dual-stack.
* `dry_run` - (Optional) Specifies whether to perform only a dry run, without performing the actual request. Valid values:
  - `true`: performs only a dry run. The system checks the request for potential issues, including missing parameter values, incorrect request syntax, and service limits. If the request fails the dry run, an error message is returned. If the request passes the dry run, the DryRunOperation error code is returned.
  - **false (default)**: performs a dry run and performs the actual request. If the request passes the dry run, a 2xx HTTP status code is returned and the operation is performed.
* `endpoint_description` - (Optional) The description of the endpoint.
* `endpoint_type` - (Optional, ForceNew, Computed, Available since v1.212.0) The endpoint type.

  Only the value: Interface, indicating the Interface endpoint. You can add the service resource types of Application Load Balancer (ALB), Classic Load Balancer (CLB), and Network Load Balancer (NLB).
* `policy_document` - (Optional, Available since v1.223.2) RAM access policies. For more information about policy definitions, see Alibaba Cloud-access control (RAM) official guidance.
* `protected_enabled` - (Optional, Available since v1.212.0) Specifies whether to enable user authentication. This parameter is available in Security Token Service (STS) mode. Valid values:
  - `true`: enables user authentication. After user authentication is enabled, only the user who creates the endpoint can modify or delete the endpoint in STS mode.
  - **false (default)**: disables user authentication.
* `resource_group_id` - (Optional, Computed, Available since v1.212.0) The resource group ID.
* `security_group_ids` - (Optional, Set) The ID of the security group that is associated with the endpoint ENI. The security group can be used to control data transfer between the VPC and the endpoint ENI.

  The endpoint can be associated with up to 10 security groups.
* `service_id` - (Optional, ForceNew, Computed) The ID of the endpoint service with which the endpoint is associated.
* `service_name` - (Optional, ForceNew, Computed) The name of the endpoint service with which the endpoint is associated.
* `tags` - (Optional, Map, Available since v1.212.0) The list of tags.
* `vpc_endpoint_name` - (Optional) The name of the endpoint.
* `vpc_id` - (Required, ForceNew) The ID of the VPC to which the endpoint belongs.
* `zone_private_ip_address_count` - (Optional, ForceNew, Computed, Int, Available since v1.212.0) The number of private IP addresses that are assigned to an elastic network interface (ENI) in each zone. Only 1 is returned.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `bandwidth` - The bandwidth of the endpoint connection.  1024 to 10240. Unit: Mbit/s.

  Note: The bandwidth of an endpoint connection is in the range of 100 to 10,240 Mbit/s. The default bandwidth is 1,024 Mbit/s. When the endpoint is connected to the endpoint service, the default bandwidth is the minimum bandwidth. In this case, the connection bandwidth range is 1,024 to 10,240 Mbit/s.
* `connection_status` - The state of the endpoint connection. 
* `create_time` - The time when the endpoint was created.
* `endpoint_business_status` - The service state of the endpoint. 
* `endpoint_domain` - The domain name of the endpoint.
* `region_id` - (Available since v1.239.0) The region ID of the endpoint.
* `status` - The state of the endpoint. 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vpc Endpoint.
* `delete` - (Defaults to 5 mins) Used when delete the Vpc Endpoint.
* `update` - (Defaults to 5 mins) Used when update the Vpc Endpoint.

## Import

Private Link Vpc Endpoint can be imported using the id, e.g.

```shell
$ terraform import alicloud_privatelink_vpc_endpoint.example <id>
```
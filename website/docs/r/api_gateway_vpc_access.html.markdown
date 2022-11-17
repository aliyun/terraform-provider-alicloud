---
subcategory: "Api Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_vpc_access"
sidebar_current: "docs-alicloud-resource-api-gateway-vpc-access"
description: |-
  Provides a Alicloud Api Gateway vpc authorization Resource.
---

# alicloud_api_gateway_app

Provides an vpc authorization resource.This authorizes the API gateway to access your VPC instances.

For information about Api Gateway vpc and how to use it, see [Set Vpc Access](https://www.alibabacloud.com/help/doc-detail/51608.htm)

-> **NOTE:** Terraform will auto build vpc authorization while it uses `alicloud_api_gateway_vpc_access` to build vpc.

## Example Usage

Basic Usage

```terraform
resource "alicloud_api_gateway_vpc_access" "foo" {
  name        = "ApiGatewayVpc"
  vpc_id      = "vpc-awkcj192ka9zalz"
  instance_id = "i-kai2ks92kzkw92ka"
  port        = 8080
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required，ForceNew) The name of the vpc authorization. 
* `vpc_id` - (Required，ForceNew) The vpc id of the vpc authorization. 
* `instance_id` - (Required，ForceNew) ID of the instance in VPC (ECS/Server Load Balance).
* `port` - (Required，ForceNew) ID of the port corresponding to the instance.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the vpc authorization of api gateway.

## Import

Api gateway app can be imported using the id, e.g.

```shell
$ terraform import alicloud_api_gateway_vpc_access.example "APiGatewayVpc:vpc-aswcj19ajsz:i-ajdjfsdlf:8080"
```

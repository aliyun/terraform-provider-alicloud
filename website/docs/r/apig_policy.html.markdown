---
subcategory: "Cloud Native API Gateway (APIG)"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_policy"
description: |-
  Provides a Alicloud APIG Policy resource.
---

# alicloud_apig_policy

Provides a APIG Policy resource.



For information about APIG Policy and how to use it, see [What is Policy](https://next.api.alibabacloud.com/document/APIG/2024-03-27/CreateAndAttachPolicy).

-> **NOTE:** Available since v1.292.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "${var.name}-vpc"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = "${var.name}-vsw"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-hangzhou-b"
  cidr_block   = "192.168.15.0/24"
}

resource "alicloud_apig_gateway" "default" {
  gateway_name = "${var.name}-gw"
  gateway_type = "API"
  payment_type = "PayAsYouGo"
  spec         = "apigw.small.x1"
  network_access_config {
    type = "Intranet"
  }
  vpc {
    vpc_id = alicloud_vpc.default.id
  }
  vswitch {
    vswitch_id = alicloud_vswitch.default.id
  }
  zone_config {
    select_option = "Auto"
  }
}

resource "alicloud_apig_service" "default" {
  service_name = "${var.name}-svc"
  source_type  = "DNS"
  gateway_id   = alicloud_apig_gateway.default.id
  addresses    = ["httpbin.org:8080"]
}

resource "alicloud_apig_policy" "default" {
  policy_name          = var.name
  policy_class_name    = "ServiceTls"
  policy_config        = "{\"mode\":\"SIMPLE\",\"sni\":\"aaaa\",\"enable\":true}"
  gateway_id           = alicloud_apig_gateway.default.id
  environment_id       = alicloud_apig_gateway.default.environments.0.environment_id
  attach_resource_type = "GatewayService"
  attach_resource_ids  = [alicloud_apig_service.default.id]
}
```

## Argument Reference

The following arguments are supported:
* `attach_resource_ids` - (Required, ForceNew, List) The Mount point id list.
* `attach_resource_type` - (Required, ForceNew) Policies support mount point types.
  - HttpApi:HttpApi.
  - Operation: the Operation of the HttpApi.
  - GatewayRoute: Gateway route.
  - GatewayService: Gateway service.
  - GatewayServicePort: The Gateway service port.
  - Domain: The Gateway Domain name.
  - Gateway: Gateway.
* `environment_id` - (Optional, ForceNew) Environment id
* `gateway_id` - (Optional, ForceNew) Gateway id
* `policy_class_name` - (Required, ForceNew) policy class name
* `policy_config` - (Required) Policy Configuration
* `policy_name` - (Optional) Policy Name

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `policy_class_id` - policy class id.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Policy.
* `update` - (Defaults to 5 mins) Used when update the Policy.

## Import

APIG Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_apig_policy.example <policy_id>
```
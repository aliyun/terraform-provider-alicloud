---
subcategory: "Cloud Native API Gateway (APIG)"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_service"
description: |-
  Provides a Alicloud APIG Service resource.
---

# alicloud_apig_service

Provides a APIG Service resource.



For information about APIG Service and how to use it, see [What is Service](https://next.api.alibabacloud.com/document/APIG/2024-03-27/CreateService).

-> **NOTE:** Available since v1.286.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

variable "address" {
  default = "127.0.0.1:8080"
}

variable "address_1" {
  default = "127.0.0.1:7891"
}

variable "address_2" {
  default = "127.0.0.1:7890"
}

resource "alicloud_vpc" "defaultvpc" {
  cidr_block = "172.32.0.0/12"
  vpc_name   = "zhenyuan-example"
}

resource "alicloud_vswitch" "defaultvswitch" {
  vpc_id       = alicloud_vpc.defaultvpc.id
  zone_id      = "cn-hangzhou-g"
  cidr_block   = "172.32.100.0/24"
  vswitch_name = "zhenyuan-example"
}

resource "alicloud_apig_gateway" "defaultFsRKYn" {
  network_access_config {
    type = "Intranet"
  }
  vswitch {
    vswitch_id = alicloud_vswitch.defaultvswitch.id
    name       = alicloud_vswitch.defaultvswitch.vswitch_name
  }
  zone_config {
    select_option = "Auto"
  }
  vpc {
    vpc_id = alicloud_vpc.defaultvpc.id
  }
  payment_type = "PayAsYouGo"
  gateway_name = "zhenyuanexample"
  spec         = "apigw.small.x1"
  log_config {
    sls {
      enable = false
    }
  }
}


resource "alicloud_apig_service" "default" {
  addresses    = ["${var.address}"]
  service_name = "1784511687"
  source_type  = "VIP"
  gateway_id   = alicloud_apig_gateway.defaultFsRKYn.id
  namespace    = "default"
}
```

## Argument Reference

The following arguments are supported:
* `addresses` - (Optional, List) A list of domain names or fixed addresses.
* `agent_service_config` - (Optional, Set) Agent service configuration See [`agent_service_config`](#agent_service_config) below.
* `ai_service_config` - (Optional, ForceNew, Set) ai service configuration when sourceType equals AI. See [`ai_service_config`](#ai_service_config) below.
* `dns_servers` - (Optional, List) DNS servers
* `express_type` - (Optional, ForceNew) Express type
* `gateway_id` - (Optional, ForceNew) The ID of the Cloud Native API Gateway.
* `group_name` - (Optional, ForceNew) The service group name.
Required when sourceType is MSE_NACOS.
* `health_check_config` - (Optional, Set) Health check configuration See [`health_check_config`](#health_check_config) below.
* `healthy_panic_threshold` - (Optional, Float) Healthy panic threshold
* `model_provider_id` - (Optional) Model provider ID
* `namespace` - (Optional, ForceNew) The namespace of the service:
  - sourceType is K8S, indicating the namespace of the K8S service.
When-sourceType is set to MSE_NACOS, it indicates the namespace in Nacos.

When the sourceType is K8S and MSE_NACOS, it needs to be specified.
* `outlier_detection_config` - (Optional, Set) Outlier detection configuration See [`outlier_detection_config`](#outlier_detection_config) below.
* `ports` - (Optional, List) Port information See [`ports`](#ports) below.

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `protocol` - (Optional) Service protocol

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `qualifier` - (Optional, ForceNew) The function version or alias.
* `resource_group_id` - (Optional, Computed) The ID of the resource group
* `service_name` - (Optional, ForceNew) Service Name, need to fill in manually when sourceType is VIP/DNS/AI.
* `source_type` - (Optional, ForceNew) service source type, optional value is K8S/MSE_NACOS/FC3/SAE_K8S_SERVICE/VIP/DNS/AI

### `agent_service_config`

The agent_service_config supports the following:
* `address` - (Optional) Address
* `enable_health_check` - (Optional) Enable health check
* `enable_outlier_detection` - (Optional) Enable outlier detection
* `protocols` - (Optional, List) Protocols
* `provider` - (Optional) Provider

### `ai_service_config`

The ai_service_config supports the following:
* `address` - (Optional) ai provider address
* `api_keys` - (Optional, List) api key list
* `enable_health_check` - (Optional) whether enable health check
* `protocols` - (Optional, List) model protocol list
* `provider` - (Optional) ai model provider

### `health_check_config`

The health_check_config supports the following:
* `enable` - (Optional) Whether to enable health check
* `expected_statuses` - (Optional, List) Expected HTTP status codes
* `healthy_threshold` - (Optional, Int) Healthy threshold
* `http_host` - (Optional) Health check host (optional when protocol is HTTP)
* `http_path` - (Optional) Health check path (required when protocol is HTTP)
* `interval` - (Optional, Int) Health check interval
* `protocol` - (Optional) Health check protocol TCP|HTTP|GRPC
* `timeout` - (Optional, Int) Health check response timeout
* `unhealthy_threshold` - (Optional, Int) Unhealthy threshold

### `outlier_detection_config`

The outlier_detection_config supports the following:
* `base_ejection_time` - (Optional, Int) Base ejection time
* `enable` - (Optional) Whether to enable outlier detection
* `failure_percentage_minimum_hosts` - (Optional, Int) Failure percentage minimum hosts
* `failure_percentage_threshold` - (Optional, Int) Failure percentage threshold
* `interval` - (Optional, Int) Detection interval

### `ports`

The ports supports the following:
* `name` - (Optional) Port name
* `port` - (Optional, Int) Port number
* `protocol` - (Optional) Protocol TCP|UDP

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 
* `create_timestamp` - Creation timestamp.
* `health_status` - Health status.
* `outlier_endpoints` - Outlier endpoints.
* `runtime_detail_error_code` - Runtime detail error code.
* `runtime_detail_status` - Runtime detail status.
* `unhealthy_endpoints` - Unhealthy endpoints.
* `update_timestamp` - Update timestamp.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Service.
* `delete` - (Defaults to 5 mins) Used when delete the Service.
* `update` - (Defaults to 6 mins) Used when update the Service.

## Import

APIG Service can be imported using the id, e.g.

```shell
$ terraform import alicloud_apig_service.example <service_id>
```
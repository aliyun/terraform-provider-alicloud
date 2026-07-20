---
subcategory: "Cloud Native API Gateway (APIG)"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_services"
sidebar_current: "docs-alicloud-datasource-apig-services"
description: |-
  Provides a list of Apig Service owned by an Alibaba Cloud account.
---

# alicloud_apig_services

This data source provides Apig Service available to the user.[What is Service](https://next.api.alibabacloud.com/document/APIG/2024-03-27/CreateService)

-> **NOTE:** Available since v1.286.0.

## Example Usage

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
  service_name = "1784511693"
  source_type  = "VIP"
  gateway_id   = alicloud_apig_gateway.defaultFsRKYn.id
  namespace    = "default"
}

data "alicloud_apig_services" "default" {
  ids         = ["${alicloud_apig_service.default.id}"]
  name_regex  = alicloud_apig_service.default.service_name
  gateway_id  = alicloud_apig_gateway.defaultFsRKYn.id
  source_type = "VIP"
}

output "alicloud_apig_service_example_id" {
  value = data.alicloud_apig_services.default.services.0.id
}
```

## Argument Reference

The following arguments are supported:
* `gateway_id` - (ForceNew, Optional) The ID of the Cloud Native API Gateway.
* `resource_group_id` - (ForceNew, Optional) The ID of the resource group
* `source_type` - (ForceNew, Optional) service source type, optional value is K8S/MSE_NACOS/FC3/SAE_K8S_SERVICE/VIP/DNS/AI
* `ids` - (Optional, Computed) A list of Service IDs. 
* `name_regex` - (Optional) A regex string to filter results by Group Metric Rule name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Service IDs.
* `names` - A list of name of Services.
* `services` - A list of Service Entries. Each element contains the following attributes:
    * `addresses` - A list of domain names or fixed addresses.
    * `agent_service_config` - Agent service configuration.
        * `address` - Address.
        * `enable_health_check` - Enable health check.
        * `enable_outlier_detection` - Enable outlier detection.
        * `protocols` - Protocols.
        * `provider` - Provider.
    * `ai_service_config` - ai service configuration when sourceType equals AI.
        * `address` - ai provider address.
        * `api_keys` - api key list.
        * `enable_health_check` - whether enable health check.
        * `protocols` - model protocol list.
        * `provider` - ai model provider.
    * `create_timestamp` - Creation timestamp.
    * `dns_servers` - DNS servers.
    * `express_type` - Express type.
    * `gateway_id` - The ID of the Cloud Native API Gateway.
    * `group_name` - The service group name.
    * `health_check_config` - Health check configuration.
        * `enable` - Whether to enable health check.
        * `expected_statuses` - Expected HTTP status codes.
        * `healthy_threshold` - Healthy threshold.
        * `http_host` - Health check host (optional when protocol is HTTP).
        * `http_path` - Health check path (required when protocol is HTTP).
        * `interval` - Health check interval.
        * `protocol` - Health check protocol TCP|HTTP|GRPC.
        * `timeout` - Health check response timeout.
        * `unhealthy_threshold` - Unhealthy threshold.
    * `health_status` - Health status.
    * `healthy_panic_threshold` - Healthy panic threshold.
    * `model_provider_id` - Model provider ID.
    * `namespace` - The namespace of the service:.
    * `outlier_detection_config` - Outlier detection configuration.
        * `base_ejection_time` - Base ejection time.
        * `enable` - Whether to enable outlier detection.
        * `failure_percentage_minimum_hosts` - Failure percentage minimum hosts.
        * `failure_percentage_threshold` - Failure percentage threshold.
        * `interval` - Detection interval.
    * `outlier_endpoints` - Outlier endpoints.
    * `qualifier` - The function version or alias.
    * `resource_group_id` - The ID of the resource group.
    * `runtime_detail_error_code` - Runtime detail error code.
    * `runtime_detail_status` - Runtime detail status.
    * `service_id` - service id.
    * `service_name` - Service Name, need to fill in manually when sourceType is VIP/DNS/AI.
    * `source_type` - service source type, optional value is K8S/MSE_NACOS/FC3/SAE_K8S_SERVICE/VIP/DNS/AI.
    * `unhealthy_endpoints` - Unhealthy endpoints.
    * `update_timestamp` - Update timestamp.
    * `id` - The ID of the resource supplied above.

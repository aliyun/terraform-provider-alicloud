---
subcategory: "Function Compute Service (FC)"
layout: "alicloud"
page_title: "Alicloud: alicloud_fcv2_function_on_demand_config"
description: |-
  Provides a Alicloud FCV2 Function On Demand Config resource.
---

# alicloud_fcv2_function_on_demand_config

Provides a FCV2 Function On Demand Config resource. The on-demand config caps the maximum number of on-demand instances a function (under a given service alias or `LATEST`) may scale to.

For information about FCV2 Function On Demand Config and how to use it, see [What is PutFunctionOnDemandConfig](https://www.alibabacloud.com/help/en/function-compute/latest/putfunctionondemandconfig).

-> **NOTE:** Available since v2.3.0.

-> **NOTE:** The `qualifier` must be a service alias or `LATEST`; the FC-Open API does not accept bare version numbers for the on-demand config.

## Example Usage

Basic Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

variable "name" {
  default = "terraform-example"
}

resource "alicloud_fc_service" "default" {
  name        = var.name
  description = var.name
}

resource "alicloud_fcv2_function" "default" {
  service_name  = alicloud_fc_service.default.name
  function_name = var.name
  runtime       = "python3.9"
  handler       = "index.handler"
  code {
    zip_file = "UEsDBBQAAAAIABAfZl27wpNPTgAAAFMAAAAIABwAaW5kZXgucHlVVAkAA2/f7Gpv3+xqdXgLAAEE6AMAAAToAwAAHcpRCoAgDADQf08x/FKIDtBtIhcWc5M5K4juXvS+31aqqMHehF3CFfLMiVADHsg2wCJseFmcHHwUrSv/eUy91BZun5FI/AT+FKXkn+heUEsBAh4DFAAAAAgAEB9mXbvCk09OAAAAUwAAAAgAGAAAAAAAAQAAAKSBAAAAAGluZGV4LnB5VVQFAANv3+xqdXgLAAEE6AMAAAToAwAAUEsFBgAAAAABAAEATgAAAJAAAAAAAA=="
  }
}

resource "alicloud_fcv2_function_on_demand_config" "default" {
  service_name           = alicloud_fc_service.default.name
  function_name          = alicloud_fcv2_function.default.function_name
  qualifier              = "LATEST"
  maximum_instance_count = 10
}
```

## Argument Reference

The following arguments are supported:

* `function_name` - (Required, ForceNew) The name of the function that the on-demand config attaches to.
* `maximum_instance_count` - (Required) The maximum number of on-demand instances the function may scale to. PutFunctionOnDemandConfig is an upsert, so changing this value updates the existing config in place.
* `qualifier` - (Required, ForceNew) The service alias or `LATEST` that identifies the function version the on-demand config applies to. A bare version number is not accepted.
* `service_name` - (Required, ForceNew) The name of the service that owns the function.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, formatted as `<service_name>:<function_name>:<qualifier>`.
* `resource` - The FC-internal resource descriptor, formatted as `<accountId>#<serviceName>#<qualifier>#<functionName>`.

## Import

FCV2 Function On Demand Config can be imported using the id, which consists of service_name, function_name and qualifier, e.g.

```shell
$ terraform import alicloud_fcv2_function_on_demand_config.example <service_name>:<function_name>:<qualifier>
```

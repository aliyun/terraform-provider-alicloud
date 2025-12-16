---
subcategory: "Function Compute Service (FC)"
layout: "alicloud"
page_title: "Alicloud: alicloud_fc_service"
sidebar_current: "docs-alicloud-resource-fc"
description: |-
  Provides a Alicloud Function Compute Service resource. The resource is the base of launching Function and Trigger configuration.
---

# alicloud_fc_service

Provides a Alicloud Function Compute Service resource. The resource is the base of launching Function and Trigger configuration.
 For information about Service and how to use it, see [What is Function Compute](https://www.alibabacloud.com/help/en/fc/developer-reference/api-fc-open-2021-04-06-createservice).

-> **NOTE:** The resource requires a provider field 'account_id'. [See account_id](https://www.terraform.io/docs/providers/alicloud/index.html#account_id).

-> **NOTE:** If you happen the error "Argument 'internetAccess' is not supported", you need to log on web console and click button "Apply VPC Function"
which is in the upper of [Function Service Web Console](https://fc.console.aliyun.com/) page.

-> **NOTE:** Currently not all regions support Function Compute Service.
For more details supported regions, see [Service endpoints](https://www.alibabacloud.com/help/doc-detail/52984.htm)

-> **NOTE:** Available since v1.93.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_fc_service&exampleId=35d31535-c004-be8f-b8b6-1bedc610f346729e2f9e&activeTab=example&spm=docs.r.fc_service.0.35d31535c0&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_log_project" "default" {
  project_name = "example-value-${random_integer.default.result}"
}

resource "alicloud_log_store" "default" {
  project_name  = alicloud_log_project.default.project_name
  logstore_name = "example-value"
}

# add index for logstore, which is used to query logs
locals {
  sls_default_token = ", '\";=()[]{}?@&<>/:\n\t\r"
}

resource "alicloud_log_store_index" "example" {
  project  = alicloud_log_project.default.project_name
  logstore = alicloud_log_store.default.logstore_name
  full_text {
    case_sensitive = false
    token          = local.sls_default_token
  }
  field_search {
    name             = "aggPeriodSeconds"
    enable_analytics = true
    type             = "long"
    token            = local.sls_default_token
  }
  field_search {
    name             = "concurrentRequests"
    enable_analytics = true
    type             = "long"
    token            = local.sls_default_token
  }
  field_search {
    name             = "cpuPercent"
    enable_analytics = true
    type             = "double"
    token            = local.sls_default_token
  }
  field_search {
    name             = "cpuQuotaPercent"
    enable_analytics = true
    type             = "double"
    token            = local.sls_default_token
  }
  field_search {
    name             = "functionName"
    enable_analytics = true
    type             = "text"
    token            = local.sls_default_token
    case_sensitive   = true
  }
  field_search {
    name             = "hostname"
    enable_analytics = true
    type             = "text"
    token            = local.sls_default_token
  }
  field_search {
    name             = "instanceID"
    enable_analytics = true
    type             = "text"
    token            = local.sls_default_token
  }
  field_search {
    name             = "ipAddress"
    enable_analytics = true
    type             = "text"
    token            = local.sls_default_token
  }
  field_search {
    name             = "memoryLimitMB"
    enable_analytics = true
    type             = "double"
    token            = local.sls_default_token
  }
  field_search {
    name             = "memoryUsageMB"
    enable_analytics = true
    type             = "double"
    token            = local.sls_default_token
  }
  field_search {
    name             = "memoryUsagePercent"
    enable_analytics = true
    type             = "double"
    token            = local.sls_default_token
  }
  field_search {
    name             = "operation"
    enable_analytics = true
    type             = "text"
    token            = local.sls_default_token
  }
  field_search {
    name             = "qualifier"
    enable_analytics = true
    type             = "text"
    token            = local.sls_default_token
    case_sensitive   = true
  }
  field_search {
    name             = "rxBytes"
    enable_analytics = true
    type             = "long"
    token            = local.sls_default_token
  }
  field_search {
    name             = "rxTotalBytes"
    enable_analytics = true
    type             = "long"
    token            = local.sls_default_token
  }
  field_search {
    name             = "serviceName"
    enable_analytics = true
    type             = "text"
    token            = local.sls_default_token
    case_sensitive   = true
  }
  field_search {
    name             = "txBytes"
    enable_analytics = true
    type             = "long"
    token            = local.sls_default_token
  }
  field_search {
    name             = "txTotalBytes"
    enable_analytics = true
    type             = "long"
    token            = local.sls_default_token
  }
  field_search {
    name             = "versionId"
    enable_analytics = true
    type             = "text"
    token            = local.sls_default_token
  }
  field_search {
    name             = "events"
    enable_analytics = true
    type             = "json"
    token            = local.sls_default_token
  }
  field_search {
    name             = "isColdStart"
    enable_analytics = true
    type             = "text"
    token            = local.sls_default_token
  }
  field_search {
    name             = "hasFunctionError"
    enable_analytics = true
    type             = "text"
    token            = local.sls_default_token
  }
  field_search {
    name             = "errorType"
    enable_analytics = true
    type             = "text"
    token            = local.sls_default_token
  }
  field_search {
    name             = "triggerType"
    enable_analytics = true
    type             = "text"
    token            = local.sls_default_token
  }
  field_search {
    name             = "durationMs"
    enable_analytics = true
    type             = "double"
    token            = local.sls_default_token
  }
  field_search {
    name             = "statusCode"
    enable_analytics = true
    type             = "long"
    token            = local.sls_default_token
  }
}

resource "alicloud_ram_role" "default" {
  name        = "fcservicerole-${random_integer.default.result}"
  document    = <<EOF
  {
      "Statement": [
        {
          "Action": "sts:AssumeRole",
          "Effect": "Allow",
          "Principal": {
            "Service": [
              "fc.aliyuncs.com"
            ]
          }
        }
      ],
      "Version": "1"
  }
  EOF
  description = "this is a example"
  force       = true
}

resource "alicloud_ram_role_policy_attachment" "default" {
  role_name   = alicloud_ram_role.default.name
  policy_name = "AliyunLogFullAccess"
  policy_type = "System"
}

resource "alicloud_fc_service" "default" {
  name        = "example-value-${random_integer.default.result}"
  description = "example-value"
  role        = alicloud_ram_role.default.arn
  log_config {
    project                 = alicloud_log_project.default.project_name
    logstore                = alicloud_log_store.default.logstore_name
    enable_instance_metrics = true
    enable_request_metrics  = true
  }
  tags = {
    "ExampleKey" = "example-value"
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_fc_service&spm=docs.r.fc_service.example&intl_lang=EN_US)

## Module Support

You can use to the existing [fc module](https://registry.terraform.io/modules/terraform-alicloud-modules/fc/alicloud) to create a service and a function quickly and then set several triggers for it.

## Argument Reference

The following arguments are supported:

* `name` - (Optional, ForceNew) The Function Compute Service name. It is the only in one Alicloud account and is conflict with `name_prefix`.
* `name_prefix` - (Optional, ForceNew) Setting a prefix to get a only name. It is conflict with `name`.
* `description` - (Optional) The Function Compute Service description.
* `internet_access` - (Optional) Whether to allow the Service to access Internet. Default to "true".
* `role` - (Optional) RAM role arn attached to the Function Compute Service. This governs both who / what can invoke your Function, as well as what resources our Function has access to. See [User Permissions](https://www.alibabacloud.com/help/doc-detail/52885.htm) for more details.
* `log_config` - (Optional) Provide this to store your Function Compute Service logs. Fields documented below. See [Create a Service](https://www.alibabacloud.com/help/doc-detail/51924.htm). `log_config` requires the following: (**NOTE:** If both `project` and `logstore` are empty, log_config is considered to be empty or unset.). See [`log_config`](#log_config) below.
* `vpc_config` - (Optional) Provide this to allow your Function Compute Service to access your VPC. Fields documented below. See [Function Compute Service in VPC](https://www.alibabacloud.com/help/faq-detail/72959.htm). `vpc_config` requires the following: (**NOTE:** If both `vswitch_ids` and `security_group_id` are empty, vpc_config is considered to be empty or unset.). See [`vpc_config`](#vpc_config) below.
* `nas_config` - (Optional, available in 1.96.0+) Provide [NAS configuration](https://www.alibabacloud.com/help/doc-detail/87401.htm) to allow Function Compute Service to access your NAS resources. See [`nas_config`](#nas_config) below.
* `tracing_config` - (Optional, available in 1.183.0+) Provide this to allow your Function Compute to report tracing information. Fields documented below. See [Function Compute Tracing Config](https://help.aliyun.com/document_detail/189805.html). `tracing_config` requires the following: (**NOTE:** If both `type` and `params` are empty, tracing_config is considered to be empty or unset.). See [`tracing_config`](#tracing_config) below.
* `publish` - (Optional, available in 1.101.0+) Whether to publish creation/change as new Function Compute Service Version. Defaults to `false`.
* `tags` - (Optional, available in 1.212.0+) Map for tagging resources.

### `log_config`

The log_config supports the following: 

* `project` - (Required) The project name of the Alicloud Simple Log Service.
* `logstore` - (Required) The log store name of Alicloud Simple Log Service.
* `enable_request_metrics` - (Optional, available in 1.183.0+) Enable request level metrics.
* `enable_instance_metrics` - (Optional, available in 1.183.0+) Enable instance level metrics.

### `vpc_config`

The vpc_config supports the following: 

* `vpc_id` - (Optional) A vpc ID associated with the Function Compute Service.
* `vswitch_ids` - (Required) A list of vswitch IDs associated with the Function Compute Service.
* `security_group_id` - (Required) A security group ID associated with the Function Compute Service.

### `nas_config`

The nas_config supports the following: 

* `user_id` - (Required) The user id of your NAS file system.
* `group_id` - (Required) The group id of your NAS file system.
* `mount_points` - (Required) Config the NAS mount points.See [`mount_points`](#nas_config-mount_points) below.

### `nas_config-mount_points`

The nas_config-mount_points supports the following: 

* `server_addr` - (Required) The address of the remote NAS directory.
* `mount_dir` - (Required) The local address where to mount your remote NAS directory.

### `tracing_config`

The tracing_config supports the following: 

* `type` - (Required) Tracing protocol type. Currently, only Jaeger is supported.
* `params` - (Required) Tracing parameters, which type is map[string]string. When the protocol type is Jaeger, the key is "endpoint" and the value is your tracing intranet endpoint. For example endpoint: http://tracing-analysis-dc-hz.aliyuncs.com/adapt_xxx/api/traces.

## Attributes Reference

The following arguments are exported:

* `id` - The ID of the FC Service. The value is the same as name.
* `service_id` - The Function Compute Service ID.
* `last_modified` - The date this resource was last modified.
* `version` - The latest published version of your Function Compute Service.

## Import

Function Compute Service can be imported using the id or name, e.g.

```shell
$ terraform import alicloud_fc_service.foo my-fc-service
```

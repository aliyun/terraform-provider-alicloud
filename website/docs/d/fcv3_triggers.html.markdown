---
subcategory: "Function Compute Service V3 (FCV3)"
layout: "alicloud"
page_title: "Alicloud: alicloud_fcv3_triggers"
sidebar_current: "docs-alicloud-datasource-fcv3-triggers"
description: |-
  Provides a list of Fcv3 Trigger owned by an Alibaba Cloud account.
---

# alicloud_fcv3_triggers

This data source provides Fcv3 Trigger available to the user.[What is Trigger](https://next.api.alibabacloud.com/document/FC/2023-03-30/CreateTrigger)

-> **NOTE:** Available since v1.250.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-exampleTriggerResourceAPI"
}

provider "alicloud" {
  region = "cn-shanghai"
}

variable "function_name" {
  default = "terraform-exampleTriggerResourceAPI"
}

variable "trigger_name" {
  default = "exampleTrigger_HTTP"
}

resource "alicloud_fcv3_function" "function" {
  memory_size = "512"
  cpu         = 0.5
  handler     = "index.Handler"
  code {
    zip_file = "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAIAAAAaW5kZXgucHmEkEFKxEAQRfd9ig9ZTCJOooIwDMwNXLqXnnQlaalUhU5lRj2KZ/FOXkESGR114bJ/P/7jV4b1xRq1hijtFpM1682cuNgPmgysbRulPT0fRxXnMtwrSPyeCdYRokSLnuMLJTTkbUqEvDMbxm1VdcRD6Tk+T1LW2ldB66knsYdA5iNX17ebm6tN2VnPhcswMPmREPuBacb+CiapLarAj9gT6/H97dVlCNScY3mtYvRkxdZlwDKDEnanPWVLdrdkeXEGlFEazVdfPVHaVeHc3N15CUwppwOJXeK7HshAB8NuOU7J6sP4SRXuH/EvbUfMiqMmDqv5M5FNSfAj/wgAAP//UEsHCPl//NYAAQAArwEAAFBLAQIUABQACAAIAAAAAAD5f/zWAAEAAK8BAAAIAAAAAAAAAAAAAAAAAAAAAABpbmRleC5weVBLBQYAAAAAAQABADYAAAA2AQAAAAA="
  }
  function_name = var.name
  runtime       = "python3.9"
  disk_size     = "512"
  log_config {
    log_begin_rule = "None"
  }
}


resource "alicloud_fcv3_trigger" "default" {
  function_name  = alicloud_fcv3_function.function.function_name
  trigger_type   = "http"
  trigger_name   = "tf-exampleacceu-central-1fcv3trigger28547"
  description    = "create"
  qualifier      = "LATEST"
  trigger_config = jsonencode({ "authType" : "anonymous", "methods" : ["GET", "POST"] })
}

data "alicloud_fcv3_triggers" "default" {
  ids           = ["${alicloud_fcv3_trigger.default.id}"]
  name_regex    = alicloud_fcv3_trigger.default.trigger_name
  function_name = var.function_name
}

output "alicloud_fcv3_trigger_example_id" {
  value = data.alicloud_fcv3_triggers.default.triggers.0.id
}
```

## Argument Reference

The following arguments are supported:
* `function_name` - (Required, ForceNew) Function Name
* `ids` - (Optional, ForceNew, Computed) A list of Trigger IDs. The value is formulated as `<function_name>:<trigger_name>`.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group Metric Rule name.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Trigger IDs.
* `names` - A list of name of Triggers.
* `triggers` - A list of Trigger Entries. Each element contains the following attributes:
  * `create_time` - Creation time
  * `description` - Description of the trigger
  * `http_trigger` - HTTP trigger information
    * `url_internet` - The public domain name address. On the Internet, you can access the HTTP Trigger through the HTTP protocol or HTTPS protocol.
    * `url_intranet` - The private domain name address. In a VPC, you can access the HTTP Trigger through HTTP or HTTPS.
  * `invocation_role` - The role required by the event source (such as OSS) to call the function.
  * `last_modified_time` - The last modified time of the trigger
  * `qualifier` - The version or alias of the function
  * `source_arn` - Trigger Event source ARN
  * `status` - The state of the trigger
  * `target_arn` - Resource identity of the function
  * `trigger_config` - Trigger configuration. The configuration varies for different types of triggers.
  * `trigger_id` - Trigger ID
  * `trigger_name` - Trigger Name
  * `trigger_type` - The type of the trigger. Currently, the supported types include oss, log, http, timer, tablestore, cdn_events, mns_topic and eventbridge.
  * `id` - The ID of the resource supplied above.

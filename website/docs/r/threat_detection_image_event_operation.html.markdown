---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_image_event_operation"
description: |-
  Provides a Alicloud Threat Detection Image Event Operation resource.
---

# alicloud_threat_detection_image_event_operation

Provides a Threat Detection Image Event Operation resource.

Image Event Operation.

For information about Threat Detection Image Event Operation and how to use it, see [What is Image Event Operation](https://www.alibabacloud.com/help/zh/security-center/developer-reference/api-sas-2018-12-03-addimageeventoperation).

-> **NOTE:** Available since v1.212.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_threat_detection_image_event_operation&exampleId=4bdcbc07-86f9-3aa8-6b42-ded5ef3e19a8bd3c7c06&activeTab=example&spm=docs.r.threat_detection_image_event_operation.0.4bdcbc0786&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_threat_detection_image_event_operation" "default" {
  event_type     = "maliciousFile"
  operation_code = "whitelist"
  event_key      = "alibabacloud_ak"
  scenarios      = <<EOF
{
  "type":"default",
  "value":""
}
EOF
  event_name     = "阿里云AK"
  conditions     = <<EOF
[
  {
      "condition":"MD5",
      "type":"equals",
      "value":"0083a31cc0083a31ccf7c10367a6e783e"
  }
]
EOF
}
```

## Argument Reference

The following arguments are supported:
* `conditions` - (Required) The rule conditions. The value is in the JSON format. For more information, see [How to use it](https://www.alibabacloud.com/help/en/security-center/developer-reference/api-sas-2018-12-03-addimageeventoperation). **NOTE:** From version 1.255.0, `conditions` can be modified.
* `event_key` - (Optional, ForceNew) The keyword of the alert item.
* `event_name` - (Optional, ForceNew) The name of the alert item.
* `event_type` - (Required, ForceNew) The alert type.
* `note` - (Optional, Available since v1.255.0) The remarks.
* `operation_code` - (Required, ForceNew) The operation code.
* `scenarios` - (Optional) The application scope of the rule.
* `source` - (Optional, ForceNew, Available since v1.255.0) The source of the whitelist. Valid values:

  - `default`: image.
  - `agentless`: agentless detection.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Image Event Operation.
* `delete` - (Defaults to 5 mins) Used when delete the Image Event Operation.
* `update` - (Defaults to 5 mins) Used when update the Image Event Operation.

## Import

Threat Detection Image Event Operation can be imported using the id, e.g.

```shell
$ terraform import alicloud_threat_detection_image_event_operation.example <id>
```

---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_malicious_file_whitelist_config"
description: |-
  Provides a Alicloud Threat Detection Malicious File Whitelist Config resource.
---

# alicloud_threat_detection_malicious_file_whitelist_config

Provides a Threat Detection Malicious File Whitelist Config resource. malicious file add whitelist config.

For information about Threat Detection Malicious File Whitelist Config and how to use it, see [What is Malicious File Whitelist Config](https://www.alibabacloud.com/help/zh/security-center/developer-reference/api-sas-2018-12-03-createmaliciousfilewhitelistconfig/).

-> **NOTE:** Available since v1.214.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_threat_detection_malicious_file_whitelist_config&exampleId=8fae7ac9-a290-7ccf-e83f-0fe7b2b77946b2270f66&activeTab=example&spm=docs.r.threat_detection_malicious_file_whitelist_config.0.8fae7ac9a2&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_threat_detection_malicious_file_whitelist_config" "default" {
  operator     = "strEquals"
  field        = "fileMd6"
  target_value = "123"
  target_type  = "SELECTION_KEY"
  event_name   = "123"
  source       = "agentless"
  field_value  = "sadfas"
}
```

## Argument Reference

The following arguments are supported:
* `event_name` - (Optional) The name of the security alert associated with the representative rule.
* `field` - (Optional) Represents the alarm associated with the resource and the white field.
* `field_value` - (Optional) Represents the whiteout target value in effect for the resource.
* `operator` - (Optional) The decision operator in effect on behalf of the resource.
* `source` - (Optional, ForceNew) Business Source:
  - agentless: agentless detection.
* `target_type` - (Optional) The type of target in effect on behalf of the resource.
* `target_value` - (Optional) Represents the specific value of the target type in effect for the resource.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Malicious File Whitelist Config.
* `delete` - (Defaults to 5 mins) Used when delete the Malicious File Whitelist Config.
* `update` - (Defaults to 5 mins) Used when update the Malicious File Whitelist Config.

## Import

Threat Detection Malicious File Whitelist Config can be imported using the id, e.g.

```shell
$ terraform import alicloud_threat_detection_malicious_file_whitelist_config.example <id>
```
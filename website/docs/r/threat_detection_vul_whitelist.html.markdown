---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_vul_whitelist"
sidebar_current: "docs-alicloud-resource-threat-detection-vul-whitelist"
description: |-
  Provides a Alicloud Threat Detection Vul Whitelist resource.
---

# alicloud_threat_detection_vul_whitelist

Provides a Threat Detection Vul Whitelist resource.

For information about Threat Detection Vul Whitelist and how to use it, see [What is Vul Whitelist](https://www.alibabacloud.com/help/en/security-center/developer-reference/api-sas-2018-12-03-modifycreatevulwhitelist).

-> **NOTE:** Available in v1.195.0+.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_threat_detection_vul_whitelist&exampleId=227991f1-987d-57e3-0c07-4a0a39d3ec68d6674627&activeTab=example&spm=docs.r.threat_detection_vul_whitelist.0.227991f198&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_threat_detection_vul_whitelist" "default" {
  whitelist   = "[{\"aliasName\":\"RHSA-2021:2260: libwebp 安全更新\",\"name\":\"RHSA-2021:2260: libwebp 安全更新\",\"type\":\"cve\"}]"
  target_info = "{\"type\":\"GroupId\",\"uuids\":[],\"groupIds\":[10782678]}"
  reason      = "tf-example-reason"
}
```

## Argument Reference

The following arguments are supported:

* `whitelist` - (Required,ForceNew) Information about the vulnerability to be added to the whitelist. see [how to use it](https://www.alibabacloud.com/help/en/security-center/developer-reference/api-sas-2018-12-03-modifycreatevulwhitelist).
* `target_info` - (Optional) Set the effective range of the whitelist. see [how to use it](https://www.alibabacloud.com/help/en/security-center/developer-reference/api-sas-2018-12-03-modifycreatevulwhitelist).
* `reason` - (Optional) Reason for adding whitelist.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Vul Whitelist.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 3 mins) Used when create the Vul Whitelist.
* `update` - (Defaults to 3 mins) Used when update the Vul Whitelist.
* `delete` - (Defaults to 3 mins) Used when delete the Vul Whitelist.

## Import

Threat Detection Vul Whitelist can be imported using the id, e.g.

```shell
$ terraform import alicloud_threat_detection_vul_whitelist.example <id>
```

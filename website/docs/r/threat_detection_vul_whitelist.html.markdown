---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_vul_whitelist"
sidebar_current: "docs-alicloud-resource-threat-detection-vul-whitelist"
description: |-
  Provides a Alicloud Threat Detection Vul Whitelist resource.
---

# alicloud\_threat\_detection\_vul\_whitelist

Provides a Threat Detection Vul Whitelist resource.

For information about Threat Detection Vul Whitelist and how to use it, see [What is Vul Whitelist](https://www.alibabacloud.com/help/en/security-center/latest/api-doc-sas-2018-12-03-api-doc-modifycreatevulwhitelist).

-> **NOTE:** Available in v1.195.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_threat_detection_vul_whitelist" "default" {
  whitelist   = "[{\"aliasName\":\"RHSA-2021:2260: libwebp 安全更新\",\"name\":\"RHSA-2021:2260: libwebp 安全更新\",\"type\":\"cve\"}]"
  target_info = "{\"type\":\"GroupId\",\"uuids\":[],\"groupIds\":[10782678]}"
  reason      = "tf-example-reason"
}
```

## Argument Reference

The following arguments are supported:

* `whitelist` - (Required,ForceNew) Information about the vulnerability to be added to the whitelist. see [how to use it](https://www.alibabacloud.com/help/en/security-center/latest/api-doc-sas-2018-12-03-api-doc-modifycreatevulwhitelist).
* `target_info` - (Optional) Set the effective range of the whitelist. see [how to use it](https://www.alibabacloud.com/help/en/security-center/latest/api-doc-sas-2018-12-03-api-doc-modifycreatevulwhitelist).
* `reason` - (Optional) Reason for adding whitelist.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Vul Whitelist.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 3 mins) Used when create the Vul Whitelist.
* `update` - (Defaults to 3 mins) Used when update the Vul Whitelist.
* `delete` - (Defaults to 3 mins) Used when delete the Vul Whitelist.

## Import

Threat Detection Vul Whitelist can be imported using the id, e.g.

```shell
$ terraform import alicloud_threat_detection_vul_whitelist.example <id>
```

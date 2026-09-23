---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_whitelist_template"
description: |-
  Provide a whitelist template.
---

# alicloud_rds_whitelist_template

Provide a whitelist template resource.


For information about Resource AliCloudWhitelistTemplate and how to use it, see [What is Whitelist Template](https://www.alibabacloud.com/help/en/rds/developer-reference/api-rds-2014-08-15-modifywhitelisttemplate?).

-> **NOTE:** Available since v1.254.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_rds_whitelist_template&exampleId=993979d4-e07e-50c6-b17b-c8ee842dd39d0e4dc7a9&activeTab=example&spm=docs.r.rds_whitelist_template.0.993979d4e0&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_rds_whitelist_template" "example" {
  ip_white_list = "172.16.0.0"
  template_name = "example-whitelist"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_rds_whitelist_template&spm=docs.r.rds_whitelist_template.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

-> **NOTE:** Each instance can add up to 1000 IPs or IP segments, meaning that the total number of IPs or IP segments in all IP whitelist groups cannot exceed 1000. When there are many IPs, it is recommended to merge them into IP segments and fill them in, such as 10.23.XX.XX/24.
* `ip_white_list` - (Required) IP whitelist, multiple IP addresses should be separated by commas (,) and cannot be duplicated.Supports the following two formats:
  - IP address format, for example: 10.23.XX.XX.
  - CIDR format, for example: 10.23.XX.XX/24 (no inter domain routing, 24 represents the length of the prefix in the address, ranging from 1 to 32).
* `template_name` - (Required) Whitelist template name. Passed in when creating a template, and cannot have the same name under the same account, starting with a letter.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 60 mins) Used when creating whitelist templates.
* `delete` - (Defaults to 20 mins) Used when delete whitelist templates.
* `update` - (Defaults to 30 mins) Used when update whitelist templates.

## Import

You can use the whitelist template ID to import whitelist templates, e.g.

```shell
$ terraform import alicloud_rds_whitelist_template.example <id>
```
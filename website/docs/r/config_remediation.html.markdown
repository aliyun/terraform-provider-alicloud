---
subcategory: "Cloud Config (Config)"
layout: "alicloud"
page_title: "Alicloud: alicloud_config_remediation"
sidebar_current: "docs-alicloud-resource-config-remediation"
description: |-
  Provides a Alicloud Config Remediation resource.
---

# alicloud_config_remediation

Provides a Config Remediation resource.

For information about Config Remediation and how to use it, see [What is Remediation](https://www.alibabacloud.com/help/en/cloud-config/latest/api-config-2020-09-07-createremediation).

-> **NOTE:** Available since v1.204.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_config_remediation&exampleId=081a4610-7ed5-b544-d0f2-b74c7296a64f3eeac0da&activeTab=example&spm=docs.r.config_remediation.0.081a46107e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example-oss"
}
data "alicloud_regions" "default" {
  current = true
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_oss_bucket" "default" {
  bucket = "${var.name}-${random_integer.default.result}"
  tags = {
    For = "example"
  }
}

resource "alicloud_oss_bucket_acl" "name" {
  bucket = alicloud_oss_bucket.default.bucket
  acl    = "public-read"
}

resource "alicloud_config_rule" "default" {
  description               = "If the ACL policy of the OSS bucket denies read access from the Internet, the configuration is considered compliant."
  source_owner              = "ALIYUN"
  source_identifier         = "oss-bucket-public-read-prohibited"
  risk_level                = 1
  tag_key_scope             = "For"
  tag_value_scope           = "example"
  region_ids_scope          = data.alicloud_regions.default.regions.0.id
  config_rule_trigger_types = "ConfigurationItemChangeNotification"
  resource_types_scope      = ["ACS::OSS::Bucket"]
  rule_name                 = "oss-bucket-public-read-prohibited"
}

resource "alicloud_config_remediation" "default" {
  config_rule_id          = alicloud_config_rule.default.config_rule_id
  remediation_template_id = "ACS-OSS-PutBucketAcl"
  remediation_source_type = "ALIYUN"
  invoke_type             = "MANUAL_EXECUTION"
  params                  = "{\"bucketName\": \"${alicloud_oss_bucket.default.bucket}\", \"regionId\": \"${data.alicloud_regions.default.regions.0.id}\", \"permissionName\": \"private\"}"
  remediation_type        = "OOS"
}
```

## Argument Reference

The following arguments are supported:
* `config_rule_id` - (Required, ForceNew) Rule ID.
* `invoke_type` - (Required) Execution type, valid values: `Manual`, `Automatic`.
* `params` - (Required, JsonString) Remediation parameter.
* `remediation_source_type` - (Optional, ForceNew) Remediation resource type, valid values: `ALIYUN` , `CUSTOMER`.
* `remediation_template_id` - (Required) Remediation template ID.
* `remediation_type` - (Required, ForceNew) Remediation type, valid values: `OOS`, `FC`.

The following arguments will be discarded. Please use new fields as soon as possible:



## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `remediation_id` - Remediation ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Remediation.
* `delete` - (Defaults to 5 mins) Used when delete the Remediation.
* `update` - (Defaults to 5 mins) Used when update the Remediation.

## Import

Config Remediation can be imported using the id, e.g.

```shell
$ terraform import alicloud_config_remediation.example <id>
```
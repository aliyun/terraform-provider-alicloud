---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sls_collection_policy"
description: |-
  Provides a Alicloud SLS Collection Policy resource.
---

# alicloud_sls_collection_policy

Provides a SLS Collection Policy resource.

Orchestration policies for cloud product log collection.

For information about SLS Collection Policy and how to use it, see [What is Collection Policy](https://www.alibabacloud.com/help/zh/sls/developer-reference/api-sls-2020-12-30-upsertcollectionpolicy).

-> **NOTE:** Available since v1.232.0.

## Example Usage

Enable real-time log query for all of OSS buckets.

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_sls_collection_policy&exampleId=260f3fed-9de1-c582-2eee-b0c0657d8292b81d0b6e&activeTab=example&spm=docs.r.sls_collection_policy.0.260f3fed9d&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-shanghai"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_log_project" "project_create_01" {
  description  = var.name
  project_name = format("%s1%s", var.name, random_integer.default.result)
}

resource "alicloud_log_store" "logstore_create_01" {
  retention_period = "30"
  shard_count      = "2"
  project_name     = alicloud_log_project.project_create_01.project_name
  logstore_name    = format("%s1%s", var.name, random_integer.default.result)
}

resource "alicloud_log_project" "update_01" {
  description  = var.name
  project_name = format("%s2%s", var.name, random_integer.default.result)
}

resource "alicloud_log_store" "logstore002" {
  retention_period = "30"
  shard_count      = "2"
  project_name     = alicloud_log_project.update_01.project_name
  logstore_name    = format("%s2%s", var.name, random_integer.default.result)
}


resource "alicloud_sls_collection_policy" "default" {
  policy_config {
    resource_mode = "all"
    regions       = ["cn-hangzhou"]
  }
  data_code          = "metering_log"
  centralize_enabled = true
  product_code       = "oss"
  policy_name        = "xc-example-oss-01"
  enabled            = true
  data_config {
    data_region = "cn-hangzhou"
  }
  centralize_config {
    dest_ttl      = "3"
    dest_region   = "cn-shanghai"
    dest_project  = alicloud_log_project.project_create_01.project_name
    dest_logstore = alicloud_log_store.logstore_create_01.logstore_name
  }
  resource_directory {
    account_group_type = "custom"
    members            = ["1936728897040477"]
  }
}
```

Enable real-time log query for one or more specific OSS buckets
```terraform
variable "name" {
  default = "terraform-example-on-single-bucket"
}

provider "alicloud" {
  region = "cn-shanghai"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_log_project" "project_create_01" {
  description  = var.name
  project_name = format("%s1%s", var.name, random_integer.default.result)
}

resource "alicloud_log_store" "logstore_create_01" {
  retention_period = "30"
  shard_count      = "2"
  project_name     = alicloud_log_project.project_create_01.project_name
  logstore_name    = format("%s1%s", var.name, random_integer.default.result)
}

resource "alicloud_log_project" "update_01" {
  description  = var.name
  project_name = format("%s2%s", var.name, random_integer.default.result)
}

resource "alicloud_log_store" "logstore002" {
  retention_period = "30"
  shard_count      = "2"
  project_name     = alicloud_log_project.update_01.project_name
  logstore_name    = format("%s2%s", var.name, random_integer.default.result)
}

resource "alicloud_oss_bucket" "bucket" {
  bucket = format("%s1%s", var.name, random_integer.default.result)
}
resource "alicloud_sls_collection_policy" "default" {
  policy_config {
    resource_mode = "instanceMode"
    instance_ids  = [alicloud_oss_bucket.bucket.id]
  }
  data_code          = "access_log"
  centralize_enabled = false
  product_code       = "oss"
  policy_name        = "xc-example-oss-01"
  enabled            = true
}
```

## Argument Reference

The following arguments are supported:
* `centralize_config` - (Optional, List) Centralized transfer configuration. See [`centralize_config`](#centralize_config) below.
* `centralize_enabled` - (Optional) Whether to enable centralized Conversion. The default value is false.
* `data_code` - (Required, ForceNew) Log type encoding.
* `data_config` - (Optional, ForceNew, List) The configuration is supported only when the log type is global. For example, if the productCode is sls, global logs will be collected to the corresponding region during the first configuration. See [`data_config`](#data_config) below.
* `enabled` - (Required) Whether to open.
* `policy_config` - (Required, List) Collection rule configuration. See [`policy_config`](#policy_config) below.
* `policy_name` - (Required, ForceNew) The name of the rule, with a minimum of 3 characters and a maximum of 63 characters, must start with a letter.
* `product_code` - (Required, ForceNew) Product code.
* `resource_directory` - (Optional, List) For Resource Directory configuration, the account must have opened the resource directory and be an administrator or a delegated administrator. See [`resource_directory`](#resource_directory) below.

### `centralize_config`

The centralize_config supports the following:
* `dest_logstore` - (Optional) When the central logstore is transferred to the destination logstore, its geographical attribute should be consistent with the destRegion and belong to the destProject.
* `dest_project` - (Optional) The geographical attributes of the centralized transfer project should be consistent with the destRegion.
* `dest_region` - (Optional) Centralized transfer destination area.
* `dest_ttl` - (Optional, Int) The number of days for the central transfer destination. This is valid only if the central transfer destination log store is not created for the first time.

### `data_config`

The data_config supports the following:
* `data_region` - (Optional, ForceNew) If and only if the log type is global log type, for example, if productCode is sls, global logs will be collected to the corresponding region during the first configuration.

### `policy_config`

The policy_config supports the following:
* `instance_ids` - (Optional, List) A collection of instance IDs, valid only if resourceMode is instanceMode. Only instances whose instance ID is in the instance ID collection are collected.
* `regions` - (Optional, List) The region collection to which the instance belongs. Valid only when resourceMode is set to attributeMode. Wildcard characters are supported. If the region collection filter item is an empty array, it means that you do not need to filter by region, and all instances meet the filtering condition of the region collection. Otherwise, only instances with region attributes in the region collection are collected. The region collection and resource label of the instance. The instance objects are collected only when all of them are met.
* `resource_mode` - (Required) Resource collection mode. If all is configured, all instances under the account will be collected to the default logstore. If attributeMode is configured, filtering will be performed according to the region attribute and resource label of the instance. If instanceMode is configured, filtering will be performed according to the instance ID.
* `resource_tags` - (Optional, Map) Resource label, valid if and only if resourceMode is attributeMode.

  If the resource label filter item is empty, it means that you do not need to filter by resource label, and all instances meet the resource label filter condition. Otherwise, only instances whose resource label attributes meet the resource label configuration are collected.

  The resource tag and the region collection to which the instance belongs work together. The instance objects are collected only when all of them are met.

### `resource_directory`

The resource_directory supports the following:
* `account_group_type` - (Optional) Support all mode all and custom mode custom under this resource directory
* `members` - (Optional, List) When the resource directory is configured in the custom mode, the corresponding member account list

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `data_config` - The configuration is supported only when the log type is global. For example, if the productCode is sls, global logs will be collected to the corresponding region during the first configuration.
  * `data_project` - Valid only when the log type is global. For example, if the productCode is sls, the log is collected to the default dedicated Project of the account in a specific dataRegion.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Collection Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Collection Policy.
* `update` - (Defaults to 5 mins) Used when update the Collection Policy.

## Import

SLS Collection Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_sls_collection_policy.example <id>
```
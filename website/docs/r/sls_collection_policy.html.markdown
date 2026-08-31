---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sls_collection_policy"
description: |-
  Provides a Alicloud Log Service (SLS) Collection Policy resource.
---

# alicloud_sls_collection_policy

Provides a Log Service (SLS) Collection Policy resource.

Orchestration policies for cloud product log collection.

For information about Log Service (SLS) Collection Policy and how to use it, see [What is Collection Policy](https://www.alibabacloud.com/help/zh/sls/developer-reference/api-sls-2020-12-30-upsertcollectionpolicy).

-> **NOTE:** Available since v1.232.0.

## Example Usage

Basic Usage

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
<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_sls_collection_policy&exampleId=e25b6303-35ac-56a7-2958-d62c50ced91418e2ad80&activeTab=example&spm=docs.r.sls_collection_policy.1.e25b630335&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_sls_collection_policy&spm=docs.r.sls_collection_policy.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `centralize_config` - (Optional, Computed, Set) Centralized forwarding configuration. See [`centralize_config`](#centralize_config) below.
* `centralize_enabled` - (Optional) Specifies whether to enable centralized forwarding. Default value: false.
* `data_code` - (Required, ForceNew) Log type code.
* `data_config` - (Optional, ForceNew, Computed, Set) This parameter can be configured only when the log type is a global log type—for example, when productCode is sls. It indicates that global logs will be collected to the specified region upon initial configuration. See [`data_config`](#data_config) below.
* `enabled` - (Required) Whether enabled.
* `policy_config` - (Required, Set) Collection rule configuration. See [`policy_config`](#policy_config) below.
* `policy_name` - (Required, ForceNew) The naming rules are as follows:
  - It can contain only lowercase letters, digits, hyphens (-), and underscores (_).
  - It must start with a letter.
  - Its length must be between 3 and 63 characters.
* `product_code` - (Required, ForceNew) Product code.
* `resource_directory` - (Optional, Computed, Set) Resource Directory configuration. The account must have Resource Directory enabled and be either a management account or a delegated administrator. See [`resource_directory`](#resource_directory) below.

### `centralize_config`

The centralize_config supports the following:
* `dest_logstore` - (Optional) Destination Logstore for centralized forwarding. Its region must match destRegion and it must belong to destProject.
* `dest_project` - (Optional) Destination project for centralized forwarding. Its region must match destRegion.
* `dest_region` - (Optional) Destination region for centralized forwarding.
* `dest_ttl` - (Optional, Int) Retention period (in days) for the destination Logstore in centralized forwarding. This setting takes effect only when the destination Logstore is created for the first time.

### `data_config`

The data_config supports the following:
* `data_region` - (Optional, ForceNew) This parameter can be configured only when the log type is a global log type—for example, when productCode is sls. It indicates that global logs will be collected to the specified region upon initial configuration.

### `policy_config`

The policy_config supports the following:
* `instance_ids` - (Optional, List) The set of instance IDs. This parameter is valid only when resourceMode is set to instanceMode. Only instances whose IDs are included in this set are collected.
* `regions` - (Optional, List) The set of regions to which instances belong. This parameter is valid only when resourceMode is set to attributeMode and supports wildcards. If the region set filter is an empty array, no region-based filtering is applied, and all instances satisfy the region condition. Otherwise, only instances whose region attribute is included in this region set are collected. The region set and resource tags work together. An instance is collected only if it satisfies both conditions.
* `resource_mode` - (Required) Resource collection mode. If set to all, all instances under the account are collected into the default Logstore. If set to attributeMode, instances are filtered based on their region attributes and resource tags. If set to instanceMode, instances are filtered by their instance IDs.
* `resource_tags` - (Optional, Map) Resource tags. This parameter is valid only when resourceMode is set to attributeMode.  
If the resource tag filter is empty, no filtering by resource tags is applied, and all instances satisfy the resource tag condition. Otherwise, only instances whose resource tag attributes fully match the specified resource tag configuration are collected.  
Resource tags and the region set of the instance work together. An instance is collected only if it satisfies both conditions.

### `resource_directory`

The resource_directory supports the following:
* `account_group_type` - (Optional) Supports the all (select all) mode and custom mode under this Resource Directory.
* `members` - (Optional, List) The list of member accounts when the Resource Directory is configured in custom mode.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 
* `data_config` - This parameter can be configured only when the log type is a global log type—for example, when productCode is sls.
  * `data_project` - This setting is valid only when the log type is a global log type—for example, when productCode is sls.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Collection Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Collection Policy.
* `update` - (Defaults to 5 mins) Used when update the Collection Policy.

## Import

Log Service (SLS) Collection Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_sls_collection_policy.example <policy_name>
```
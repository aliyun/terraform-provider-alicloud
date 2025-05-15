---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_data_source_shared_rule"
description: |-
  Provides a Alicloud Data Works Data Source Shared Rule resource.
---

# alicloud_data_works_data_source_shared_rule

Provides a Data Works Data Source Shared Rule resource.

Data source sharing rule, which expresses A data source, from space A to space B (A user).

For information about Data Works Data Source Shared Rule and how to use it, see [What is Data Source Shared Rule](https://www.alibabacloud.com/help/en/dataworks/developer-reference/api-dataworks-public-2024-05-18-createdatasourcesharedrule).

-> **NOTE:** Available since v1.237.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_data_works_data_source_shared_rule&exampleId=4f30bc09-79e6-9c39-d9b8-c82b3fd6a306017da3b0&activeTab=example&spm=docs.r.data_works_data_source_shared_rule.0.4f30bc0979&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}

provider "alicloud" {
  region = "cn-chengdu"
}

resource "random_integer" "randint" {
  max = 999
  min = 1
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_data_works_project" "defaultQeRfvU" {
  description      = "源项目"
  project_name     = var.name
  display_name     = "shared_source2"
  pai_task_enabled = true
}
resource "alicloud_data_works_project" "defaultasjsH5" {
  description      = "目标空间"
  project_name     = format("%s1", var.name)
  display_name     = "shared_target2"
  pai_task_enabled = true
}

resource "alicloud_data_works_data_source" "defaultvzu0wG" {
  type                       = "hive"
  data_source_name           = format("%s2", var.name)
  connection_properties      = jsonencode({ "address" : [{ "host" : "127.0.0.1", "port" : "1234" }], "database" : "hive_database", "metaType" : "HiveMetastore", "metastoreUris" : "thrift://123:123", "version" : "2.3.9", "loginMode" : "Anonymous", "securityProtocol" : "authTypeNone", "envType" : "Prod", "properties" : { "key1" : "value1" } })
  project_id                 = alicloud_data_works_project.defaultQeRfvU.id
  connection_properties_mode = "UrlMode"
}


resource "alicloud_data_works_data_source_shared_rule" "default" {
  target_project_id = alicloud_data_works_project.defaultasjsH5.id
  data_source_id    = alicloud_data_works_data_source.defaultvzu0wG.data_source_id
  env_type          = "Prod"
}
```

## Argument Reference

The following arguments are supported:
* `data_source_id` - (Required, ForceNew, Int) The ID of the data source, that is, the unique identifier of the data source.
* `env_type` - (Required, ForceNew) The environment type of the data source shared to the target project, such as Dev (Development Environment) and Prod (production environment).
* `shared_user` - (Optional, ForceNew) The target user of the data source permission policy, which is null to share to the project.
* `target_project_id` - (Required, ForceNew, Int) The ID of the project to which the data source is shared.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<data_source_id>:<data_source_shared_rule_id>`.
* `create_time` - The creation time of the data source sharing rule.
* `data_source_shared_rule_id` - The data source sharing rule ID, that is, the unique identifier of the data source sharing rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Data Source Shared Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Data Source Shared Rule.

## Import

Data Works Data Source Shared Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_data_works_data_source_shared_rule.example <data_source_id>:<data_source_shared_rule_id>
```
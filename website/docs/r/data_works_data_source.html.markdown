---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_data_source"
description: |-
  Provides a Alicloud Data Works Data Source resource.
---

# alicloud_data_works_data_source

Provides a Data Works Data Source resource.



For information about Data Works Data Source and how to use it, see [What is Data Source](https://www.alibabacloud.com/help/en/dataworks/developer-reference/api-dataworks-public-2024-05-18-createdatasource).

-> **NOTE:** Available since v1.237.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_data_works_data_source&exampleId=88d802af-c4dc-6f5f-a816-b2ad347199691912c1b5&activeTab=example&spm=docs.r.data_works_data_source.0.88d802afc4&intl_lang=EN_US" target="_blank">
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

resource "alicloud_data_works_project" "defaultkguw4R" {
  status                  = "Available"
  description             = "tf_desc"
  project_name            = "${var.name}${random_integer.randint.id}"
  pai_task_enabled        = "false"
  display_name            = "tf_new_api_display"
  dev_role_disabled       = "true"
  dev_environment_enabled = "false"
  resource_group_id       = data.alicloud_resource_manager_resource_groups.default.ids.0
}

resource "alicloud_data_works_data_source" "default" {
  type                       = "hive"
  data_source_name           = var.name
  connection_properties      = jsonencode({ "address" : [{ "host" : "127.0.0.1", "port" : "1234" }], "database" : "hive_database", "metaType" : "HiveMetastore", "metastoreUris" : "thrift://123:123", "version" : "2.3.9", "loginMode" : "Anonymous", "securityProtocol" : "authTypeNone", "envType" : "Prod", "properties" : { "key1" : "value1" } })
  connection_properties_mode = "UrlMode"
  project_id                 = alicloud_data_works_project.defaultkguw4R.id
  description                = var.name
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_data_works_data_source&spm=docs.r.data_works_data_source.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `connection_properties` - (Required, JsonString) Data source connection configuration information, including the connection address, access identity, and environment information. The data source environment EnvType information is a member property of this object, including DEV (Development Environment) and PROD (production environment). The value of EnvType is not case-sensitive.
* `connection_properties_mode` - (Required) The configuration mode of the data source. Different types of data sources have different configuration modes. For example, MySQL data sources support UrlMode and InstanceMode.
* `data_source_name` - (Required, ForceNew) The data source name. The name of a data source in a specific environment (development environment or production environment) is unique in a project.
* `description` - (Optional) Description of the data source
* `project_id` - (Required, ForceNew, Int) The ID of the project to which the data source belongs.
* `type` - (Required, ForceNew) The type of data source. For a list of data source types, see the values listed in the API documentation.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<project_id>:<data_source_id>`.
* `create_time` - The creation time of the resource
* `create_user` - Creator of the data source
* `data_source_id` - The first ID of the resource
* `modify_time` - Modification time
* `modify_user` - Modifier of the data source
* `qualified_name` - Business Unique Key of Data Source

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Data Source.
* `delete` - (Defaults to 5 mins) Used when delete the Data Source.
* `update` - (Defaults to 5 mins) Used when update the Data Source.

## Import

Data Works Data Source can be imported using the id, e.g.

```shell
$ terraform import alicloud_data_works_data_source.example <project_id>:<data_source_id>
```
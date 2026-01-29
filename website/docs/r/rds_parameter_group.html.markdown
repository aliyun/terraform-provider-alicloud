---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_parameter_group"
description: |-
  Provides a Alicloud RDS Parameter Group resource.
---

# alicloud_rds_parameter_group

Provides a RDS Parameter Group resource.

Used to batch manage database parameters.

For information about RDS Parameter Group and how to use it, see [What is Parameter Group](https://next.api.alibabacloud.com/document/Rds/2014-08-15/CreateParameterGroup).

-> **NOTE:** Available since v1.270.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_rds_parameter_group&exampleId=ecba00a3-4705-c617-8963-76bd9ed57d6d16e9a7fb&activeTab=example&spm=docs.r.rds_parameter_group.0.ecba00a347&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}

resource "alicloud_rds_parameter_group" "default" {
  engine         = "mysql"
  engine_version = "5.7"
  parameter_detail {
    param_name  = "back_log"
    param_value = "4000"
  }
  parameter_detail {
    param_name  = "wait_timeout"
    param_value = "86460"
  }
  parameter_group_desc = var.name
  parameter_group_name = var.name
}
```

## Argument Reference

The following arguments are supported:
* `engine` - (Required, ForceNew) The database engine. Valid values:
  - `mysql`
  - `PostgreSQL`.
* `engine_version` - (Required, ForceNew) The database version. Valid values:
  - MySQL:
    * **5.6**
    * **5.7**
    * **8.0**
  - PostgreSQL:
    * **10.0**
    * **11.0**
    * **12.0**
    * **13.0**
    * **14.0**
    * **15.0**.
* `modify_mode` - (Optional) The modification mode of the parameter template. Valid values:
* `Collectivity` (default): Add or update.

-> **NOTE:**  Adds the parameters specified in the `Parameters` property to the current parameter template or updates existing parameters in the template. Other parameters in the current template remain unaffected.

* `Individual`: Overwrite.

-> **NOTE:**  Replaces all parameters in the current parameter template with the parameters specified in the `Parameters` property.


-> **NOTE:** This parameter only applies during resource update. If modified in isolation without other property changes, Terraform will not trigger any action.

* `parameter_detail` - (Optional, List) A JSON string that specifies parameters and their values, in the format: {"Parameter1":"Value1","Parameter2":"Value2",...}. For modifiable parameters, see [Configure parameters for an ApsaraDB RDS for MySQL instance](https://help.aliyun.com/document_detail/96063.html) or [Configure parameters for an ApsaraDB RDS for PostgreSQL instance](https://help.aliyun.com/document_detail/96751.html).

-> **NOTE:**  * If the `ModifyMode` parameter is set to `Individual`, the specified parameters overwrite the existing parameter template.

-> **NOTE:**  * If the `ModifyMode` parameter is set to `Collectivity`, the specified parameters are added to or modify the existing parameter template.

-> **NOTE:**  * If you do not specify this parameter, the original parameter settings remain unchanged.
 See [`parameter_detail`](#parameter_detail) below.
* `parameter_group_desc` - (Optional) The description of the parameter template. It can be 0 to 200 characters in length.
* `parameter_group_name` - (Required) The parameter template name.
* It must start with a letter and can contain letters, digits, periods (.), or underscores (_).
* It must be 8 to 64 characters in length.

-> **NOTE:** If you do not specify this parameter, the original parameter template name remains unchanged.

* `resource_group_id` - (Optional, Computed) The resource group ID. You can call DescribeDBInstanceAttribute to obtain it.

-> **NOTE:** This parameter only applies during resource creation, update. If modified in isolation without other property changes, Terraform will not trigger any action.

* `param_detail` - (Optional, Deprecated from v1.270.0) The attribute has been deprecated from 1.270.0 and using `parameter_detail` instead.


### `parameter_detail`

The parameter_detail supports the following:
* `param_name` - (Optional) The name of the parameter.  
* `param_value` - (Optional) The parameter value.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 
* `create_time` - The creation time.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Parameter Group.
* `delete` - (Defaults to 5 mins) Used when delete the Parameter Group.
* `update` - (Defaults to 5 mins) Used when update the Parameter Group.

## Import

RDS Parameter Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_rds_parameter_group.example <id>
```
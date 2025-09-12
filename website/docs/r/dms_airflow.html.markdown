---
subcategory: "Dms"
layout: "alicloud"
page_title: "Alicloud: alicloud_dms_airflow"
description: |-
  Provides a Alicloud Dms Airflow resource.
---

# alicloud_dms_airflow

Provides a Dms Airflow resource.

Airflow instance, used to schedule jobs.

For information about Dms Airflow and how to use it, see [What is Airflow](https://next.api.alibabacloud.com/document/Dms/2025-04-14/CreateAirflow).

-> **NOTE:** Available since v1.260.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING-dms"
}

resource "alicloud_security_group" "security_group" {
  description         = "terraform_example_group"
  security_group_name = "terraform_example_group"
  vpc_id              = data.alicloud_vpcs.default.ids.0
  security_group_type = "normal"
  inner_access_policy = "Accept"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = "cn-hangzhou-h"
}

resource "alicloud_ram_role" "dms_processing_data_role" {
  role_name                   = "AliyunDMSProcessingDataRole"
  assume_role_policy_document = <<EOF
  {
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Effect": "Allow",
        "Principal": {
          "Service": [
            "dms.aliyuncs.com"
          ]
        }
      }
    ],
    "Version": "1"
  }
  EOF
  description                 = "Role for DMS to access cloud resources for data processing"
}

resource "alicloud_ram_role_policy_attachment" "attach" {
  policy_name = "AliyunDMSProcessingDataRolePolicy"
  policy_type = "System"
  role_name   = alicloud_ram_role.dms_processing_data_role.role_name
}

resource "alicloud_dms_enterprise_workspace" "workspace" {
  description    = "terraformn-example"
  vpc_id         = data.alicloud_vpcs.default.ids.0
  workspace_name = "terraformn-example"
}


resource "alicloud_dms_airflow" "default" {
  security_group_id          = alicloud_security_group.security_group.id
  description                = "terraform-example"
  workspace_id               = alicloud_dms_enterprise_workspace.workspace.id
  zone_id                    = "cn-hangzhou-h"
  vpc_id                     = data.alicloud_vpcs.default.ids.0
  worker_serverless_replicas = "0"
  oss_path                   = "/"
  oss_bucket_name            = "hansheng"
  startup_file               = "default/startup.sh"
  app_spec                   = "SMALL"
  requirement_file           = "default/requirements.txt"
  dags_dir                   = "default/dags"
  airflow_name               = "tfaccdms48198"
  vswitch_id                 = data.alicloud_vswitches.default.ids.0
  plugins_dir                = "default/plugins"
}
```

## Argument Reference

The following arguments are supported:
* `airflow_name` - (Required) Name of the Airflow instance
* `app_spec` - (Required) Airflow instance specifications
* `dags_dir` - (Optional) Dag scan path
* `description` - (Required) airflow实例的描述
* `oss_bucket_name` - (Required, ForceNew) OSS bucket name
* `oss_path` - (Required, ForceNew) OSS path
* `plugins_dir` - (Optional) The path of the plugin scanned by the airflow instance.
* `requirement_file` - (Optional) Path to installable package
* `security_group_id` - (Required, ForceNew) Security group ID
* `startup_file` - (Optional) Launch script for the airflow container
* `vswitch_id` - (Required, ForceNew) Switch ID
* `vpc_id` - (Required, ForceNew) Vpc id
* `worker_serverless_replicas` - (Required, Int) Worker Node extension
* `workspace_id` - (Required, ForceNew) DMS workspace ID
* `zone_id` - (Required, ForceNew) Zone ID in the region

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<workspace_id>:<airflow_id>`.
* `airflow_id` - AirflowId
* `region_id` - The region ID of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Airflow.
* `delete` - (Defaults to 5 mins) Used when delete the Airflow.
* `update` - (Defaults to 5 mins) Used when update the Airflow.

## Import

Dms Airflow can be imported using the id, e.g.

```shell
$ terraform import alicloud_dms_airflow.example <workspace_id>:<airflow_id>
```
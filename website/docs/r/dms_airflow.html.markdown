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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_dms_airflow&exampleId=3e409bf2-0e83-ca64-4360-8e176a11479ba7e81400&activeTab=example&spm=docs.r.dms_airflow.0.3e409bf20e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = "cn-hangzhou-h"
}

resource "alicloud_security_group" "security_group" {
  description         = "terraform_example_group"
  security_group_name = "terraform_example_group"
  vpc_id              = data.alicloud_vpcs.default.ids.0
  security_group_type = "normal"
  inner_access_policy = "Accept"
}

resource "alicloud_dms_enterprise_workspace" "workspace" {
  description    = "terraform-example"
  vpc_id         = data.alicloud_vpcs.default.ids.0
  workspace_name = "terraform-example"
}


resource "alicloud_dms_airflow" "default" {
  vpc_id                     = data.alicloud_vpcs.default.ids.0
  oss_path                   = "/"
  dags_dir                   = "default/dags"
  zone_id                    = "cn-hangzhou-h"
  worker_serverless_replicas = "0"
  description                = "terraform-example"
  security_group_id          = alicloud_security_group.security_group.id
  requirement_file           = "default/requirements.txt"
  airflow_name               = "tfaccdms6513"
  plugins_dir                = "default/plugins"
  startup_file               = "default/startup.sh"
  app_spec                   = "SMALL"
  oss_bucket_name            = "hansheng"
  vswitch_id                 = data.alicloud_vswitches.default.ids.0
  workspace_id               = alicloud_dms_enterprise_workspace.workspace.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_dms_airflow&spm=docs.r.dms_airflow.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `airflow_name` - (Required) Name of the Airflow instance
* `app_spec` - (Required) Airflow instance specifications
* `dags_dir` - (Optional) Dag scan path
* `description` - (Required) Description of the airflow instance
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
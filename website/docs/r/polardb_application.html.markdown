---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_application"
sidebar_current: "docs-alicloud-resource-polardb-application"
description: |-
  Provides a PolarDB Application resource.
---

# alicloud_polardb_application

Provides a PolarDB Application resource. A PolarDB Application is an AI-driven database application that integrates with large language models (LLMs) to provide intelligent data processing capabilities.

-> **NOTE:** Available since v1.278.0.

## Example Usage

Create a PolarDB Application

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_polardb_application&exampleId=example-app-01&activeTab=example&spm=docs.r.polardb_application.0.example&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_vpc" "default" {
  vpc_name   = "terraform-example"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = "cn-beijing-k"
  vswitch_name = "terraform-example"
}

resource "alicloud_polardb_cluster" "default" {
  db_type     = "MySQL"
  db_version  = "8.0"
  pay_type    = "PostPaid"
  category    = "Normal"
  description = "terraform-example-cluster"
}

resource "alicloud_polardb_application" "default" {
  description      = "terraform-example-app"
  application_type = "polarclaw"
  architecture     = "x86"
  db_cluster_id    = alicloud_polardb_cluster.default.id
  vswitch_id       = alicloud_vswitch.default.id
  vpc_id           = alicloud_vpc.default.id
  region_id        = "cn-beijing"
  zone_id          = "cn-beijing-k"
  pay_type         = "PostPaid"
  model_from       = "bailian"
  model_base_url   = "https://dashscope.aliyuncs.com/compatible-mode/v1"
  model_name       = "qwen3.6-plus"

  components {
    component_type    = "polarclaw_comp"
    component_class   = "polar.app.g2.medium"
    component_replica = 1
  }

  parameters {
    parameter_name  = "secret.dashscope.apiKey"
    parameter_value = "ap-xxx"
  }
}
```

### Removing alicloud_polardb_application from your configuration

The `alicloud_polardb_application` resource allows you to manage your PolarDB application. If the application type is PrePaid, Terraform cannot destroy it directly. Removing this resource from your configuration will remove it from your statefile and management, but will not destroy the application. You can resume managing the application via the PolarDB Console.

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_polardb_application&spm=docs.r.polardb_application.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `application_type` - (Required, ForceNew) The type of the application. Valid value `polarclaw`.
* `architecture` - (Required, ForceNew) The architecture of the application. Valid value `x86`.
* `description` - (Optional, ForceNew) The description of the application. It must be 2 to 256 characters in length.
* `db_cluster_id` - (Optional, ForceNew) The ID of the associated PolarDB cluster.
* `region_id` - (Optional, ForceNew) The region ID of the application.
* `zone_id` - (Optional, Computed, ForceNew) The zone ID of the application.
* `vswitch_id` - (Optional, Computed, ForceNew) The ID of the VSwitch.
* `vpc_id` - (Optional, Computed, ForceNew) The ID of the VPC.
* `components` - (Optional, ForceNew, Type: list) The components of the application. See [`components`](#components) below.
* `pay_type` - (Optional, ForceNew) The billing method. Valid values are `PrePaid`, `PostPaid`. Default to `PostPaid`.
* `auto_renew` - (Optional, ForceNew) Whether to enable auto-renewal. Valid values are `true`, `false`.
* `period` - (Optional, ForceNew) The subscription duration in months. It is valid when `pay_type` is `PrePaid`. Valid values: `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `12`, `24`, `36`.
  -> **NOTE:** The attribute `period` is only used to create Subscription instance. Once effect, it will not be modified that means running `terraform apply` will not effect the resource.
* `used_time` - (Optional, ForceNew) The unit of the period. Valid values are `Month`, `Year`.
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.
* `model_from` - (Optional, ForceNew) The source of the model. Valid values are `bailian`,`custom`,`maas`.
* `ai_db_cluster_id` - (Optional, ForceNew) The ID of the AI DB cluster.
* `model_api_key` - (Optional, ForceNew) The API key for the model.
* `model_base_url` - (Optional, ForceNew) The base URL for the model API.
* `model_api` - (Optional, ForceNew) The API endpoint for the model.
* `model_name` - (Optional, ForceNew) The name of the model.
* `upgrade_version` - (Optional, Computed) Whether to upgrade the version. Valid values are `true`, `false`.
* `parameters` - (Optional, Computed, Type: list) The parameters of the application. See [`parameters`](#parameters) below.
* `component_id` - (Optional, Computed) The ID of the component.
* `security_ip_list` - (Optional, Computed) The list of IP addresses allowed to access the application.
* `security_ip_array_name` - (Optional, Computed) The name of the IP whitelist group.
* `modify_mode` - (Optional) The method for modifying the IP whitelist. Valid values are `Cover`, `Append`, `Delete`.

### `components`

The `components` block supports the following:

* `component_type` - (Optional, ForceNew) The type of the component.
* `component_class` - (Optional, ForceNew) The class/specification of the component.
* `component_replica` - (Optional, ForceNew) The number of replicas for the component.

### `parameters`

The `parameters` block supports the following:

* `parameter_name` - (Optional) The name of the parameter.
* `parameter_value` - (Optional) The value of the parameter.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the PolarDB Application.
* `status` - The status of the application.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 50 mins) Used when creating the polardb application.
* `update` - (Defaults to 50 mins) Used when updating the polardb application.
* `delete` - (Defaults to 10 mins) Used when deleting the polardb application.

## Import

PolarDB Application can be imported using the id, e.g.

```shell
$ terraform import alicloud_polardb_application.example pa-abc12345678
```
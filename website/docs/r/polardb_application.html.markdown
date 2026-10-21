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

-> **NOTE:** Available since v1.279.0.

-> **NOTE:** `vpc_nat_gateway_id`, `security_ip_type`, `security_ip_list`, `component_id`, the `storage_type`/`storage_instance_id`/`storage_capacity`/`storage_performance_level` members of `storages` and the database/model members of `mem_application_spec` are read back from the application, so `terraform import` restores them. `storages`, `endpoints`, `security_ip_list`, `security_ip_array_name` and `component_id` are `Computed` because the query API reports them even when the configuration declares none. Two values have an environment boundary: the query API may report only the default whitelist entry, in which case an imported custom `security_ip_array_name` is reported as `default`, and the endpoint description is only restored when the query reports it (the default endpoint may not). The remaining create-only inputs — `parameters`, `modify_mode`, `model_api_key`, `model_from`, `model_base_url`, `model_name`, `model_api`, `ai_db_cluster_id`, `resource_group_id`, `target_version`, `skill_template_id`, `dnat_ip_address`, `dnat_entries`, `auth_provider`, `auth_provider_config`, `auto_allocate_public_eip`, `agentic_db_cluster_id`, `agentic_db_branch_spec`, `knowledge_application_spec`, `storages.storage_endpoint`, `storages.endpoint_id`, `storages.mount_path`, `storages.container_mount_path`, `mem_application_spec.db_password`, `mem_application_spec.embedder_model_dimension`, `mem_application_spec.shard`, `architecture` and `upgrade_version` — are not returned by the query API (or are operation triggers) and keep their configured value; spell them out in the configuration after an import. `dnat_entries` is a special case: `DescribeApplicationAttribute` reports `DnatMappings`, but the front port is allocated by the service when it is omitted and the member has no computed attribute, so the collection is left at the configured value.

## Example Usage

Create a PolarDB Application

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_polardb_application&exampleId=451283f7-e202-aa12-7d85-67a92030844dfe3851a8&activeTab=example&spm=docs.r.polardb_application.0.451283f7e2&intl_lang=EN_US" target="_blank">
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

  modify_mode            = "Append"
  security_ip_array_name = "white_list"
  security_ip_list       = ["127.0.0.1", "192.168.1.2"]
}
```

### Removing alicloud_polardb_application from your configuration

The `alicloud_polardb_application` resource allows you to manage your PolarDB application. If the application type is PrePaid, Terraform cannot destroy it directly. Removing this resource from your configuration will remove it from your statefile and management, but will not destroy the application. You can resume managing the application via the PolarDB Console.

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_polardb_application&spm=docs.r.polardb_application.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `application_type` - (Required, ForceNew) The type of the application. Valid value `polarclaw`.

  -> **NOTE:** When `application_type` is `mem0`, `mem_application_spec` is required together with `embedder_model`, `llm_model` and `reranker_model`. When `application_type` is `knowledge`, `knowledge_application_spec` with `llm_model` is required. Other application types may omit both blocks.
* `architecture` - (Required, ForceNew) The architecture of the application. Valid value `x86`.
* `description` - (Optional, ForceNew) The description of the application. It must be 2 to 256 characters in length.
* `db_cluster_id` - (Optional, ForceNew) The ID of the associated PolarDB cluster.
* `region_id` - (Required, ForceNew) The region ID of the application.
* `zone_id` - (Optional, Computed, ForceNew) The zone ID of the application.
* `vswitch_id` - (Optional, Computed, ForceNew) The ID of the VSwitch.
* `vpc_id` - (Optional, Computed, ForceNew) The ID of the VPC.
* `components` - (Optional, ForceNew, Type: list) The components of the application. See [`components`](#components) below.
* `pay_type` - (Optional, ForceNew) The billing method. Valid values are `PrePaid`, `PostPaid`. Default to `PostPaid`.
* `auto_renew` - (Optional, ForceNew) Whether to enable auto-renewal. Valid values are `true`, `false`.
* `period` - (Optional, ForceNew) The subscription duration in months. It is valid when `pay_type` is `PrePaid`. Valid values: `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `12`, `24`, `36`.

  -> **NOTE:** The attribute `period` is only used to create Subscription instance. Once effect, it will not be modified that means running `terraform apply` will not effect the resource.
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.
* `model_from` - (Optional, ForceNew) The source of the model. Valid values are `bailian`,`custom`,`maas`.
* `ai_db_cluster_id` - (Optional, ForceNew) The ID of the AI DB cluster.
* `model_api_key` - (Optional, ForceNew) The API key for the model.
* `model_base_url` - (Optional, ForceNew) The base URL for the model API.
* `model_api` - (Optional, ForceNew) The API endpoint for the model.
* `model_name` - (Optional, ForceNew) The name of the model.
* `upgrade_version` - (Optional, Computed) Whether to upgrade the version. Valid values are `true`, `false`.
* `parameters` - (Optional, Type: list) The parameters of the application. See [`parameters`](#parameters) below.
* `component_id` - (Optional, Computed) The ID of the component.
* `security_ip_list` - (Optional, Computed, Type: list) The list of IP addresses allowed to access the application. Read back from the application whitelist entry, so an import restores it.
* `security_ip_array_name` - (Optional, Computed) The name of the IP whitelist group. Read back from the application whitelist entry.
* `modify_mode` - (Optional) The method for modifying the IP whitelist. Valid values are `Cover`, `Append`, `Delete`.
* `target_version` - (Optional, ForceNew, String) The target version of the application. Available since v2.2.0.
* `skill_template_id` - (Optional, ForceNew, String) The ID of the skill template. Available since v2.2.0.
* `dnat_ip_address` - (Optional, ForceNew, String) The DNAT IP address of the application. Available since v2.2.0.
* `auth_provider` - (Optional, ForceNew, String) The authentication provider. Available since v2.2.0.
* `auth_provider_config` - (Optional, ForceNew, String) The configuration of the authentication provider. Available since v2.2.0.
* `vpc_nat_gateway_id` - (Optional, Computed, ForceNew, String) The ID of the VPC NAT gateway. Available since v2.2.0.
* `security_ip_type` - (Optional, Computed, ForceNew, String) The type of the security IP. Read back from the application whitelist entry, so an import restores it. Available since v2.2.0.
* `auto_allocate_public_eip` - (Optional, ForceNew, Bool) Whether to automatically allocate a public EIP for the application. Available since v2.2.0.
* `agentic_db_cluster_id` - (Optional, ForceNew, String) The ID of the agentic database cluster. Available since v2.2.0.
* `storages` - (Optional, Computed, ForceNew, Set) The storage mounts of the application. See [`storages`](#storages) below. Available since v2.2.0.
* `mem_application_spec` - (Optional, Computed, ForceNew, List) The memory application specification. See [`mem_application_spec`](#mem_application_spec) below. Available since v2.2.0.
* `agentic_db_branch_spec` - (Optional, ForceNew, List) The agentic database branch specification. See [`agentic_db_branch_spec`](#agentic_db_branch_spec) below. Available since v2.2.0.
* `endpoints` - (Optional, Computed, ForceNew, Set) The endpoints of the application. See [`endpoints`](#endpoints) below. Available since v2.2.0.
* `knowledge_application_spec` - (Optional, ForceNew, List) The knowledge application specification. See [`knowledge_application_spec`](#knowledge_application_spec) below. Available since v2.2.0.
* `dnat_entries` - (Optional, ForceNew, Set) The DNAT entries of the application. See [`dnat_entries`](#dnat_entries) below. Available since v2.2.0.
* `tags` - (Optional, Map) The tags of the application. Available since v2.2.0.

### `components`

The `components` block supports the following:

* `component_type` - (Optional, ForceNew) The type of the component.
* `component_class` - (Optional, ForceNew) The class/specification of the component.
* `component_replica` - (Optional, ForceNew) The number of replicas for the component.

### `parameters`

The `parameters` block supports the following:

* `parameter_name` - (Optional) The name of the parameter.
* `parameter_value` - (Optional) The value of the parameter.

### `storages`

The `storages` block supports the following:

When a storage member is declared, the attributes the configuration omits keep the values the application reports: leaving `storage_type`, `storage_instance_id`, `storage_capacity` or `storage_performance_level` out of the block does not plan a replacement, and `storage_capacity` is compared without its unit (`20` and `20Gi` are the same capacity). Declaring one of them and changing it later plans a replacement, because every `storages` member attribute is `ForceNew`.

* `storage_type` - (Optional, Computed, ForceNew, String) The type of the storage. Available since v2.2.0.
* `storage_instance_id` - (Optional, Computed, ForceNew, String) The ID of the storage instance. Available since v2.2.0.
* `endpoint_id` - (Optional, ForceNew, String) The ID of the endpoint. Available since v2.2.0.
* `mount_path` - (Optional, ForceNew, String) The mount path of the storage. Available since v2.2.0.
* `container_mount_path` - (Optional, ForceNew, String) The mount path inside the container. Available since v2.2.0.
* `storage_capacity` - (Optional, Computed, ForceNew, String) The capacity of the storage. Available since v2.2.0.
* `storage_endpoint` - (Optional, ForceNew, String) The endpoint of the storage. Available since v2.2.0.
* `storage_performance_level` - (Optional, Computed, ForceNew, String) The performance level of the storage. Available since v2.2.0.

### `mem_application_spec`

The `mem_application_spec` block supports the following:

* `db_name` - (Optional, Computed, ForceNew, String) The name of the database. Available since v2.2.0.
* `db_password` - (Optional, ForceNew, String) The password of the database. Available since v2.2.0.
* `db_user` - (Optional, Computed, ForceNew, String) The user of the database. Available since v2.2.0.
* `embedder_model` - (Optional, Computed, ForceNew, String) The embedder model. Available since v2.2.0.
* `embedder_model_dimension` - (Optional, ForceNew, Int) The dimension of the embedder model. Available since v2.2.0.
* `graph_llm_model` - (Optional, Computed, ForceNew, String) The graph LLM model. Available since v2.2.0.
* `llm_model` - (Optional, Computed, ForceNew, String) The LLM model. Available since v2.2.0.
* `memory_type` - (Optional, Computed, ForceNew, String) The type of the memory. Available since v2.2.0.
* `project_name` - (Optional, Computed, ForceNew, String) The name of the project. Available since v2.2.0.
* `reranker_model` - (Optional, Computed, ForceNew, String) The reranker model. Available since v2.2.0.
* `shard` - (Optional, ForceNew, Int) The shard count. Available since v2.2.0.
* `visual_embedder_model` - (Optional, Computed, ForceNew, String) The visual embedder model. Available since v2.2.0.
* `visual_embedder_model_dimension` - (Optional, Computed, ForceNew, Int) The dimension of the visual embedder model. Available since v2.2.0.
* `visual_llm_modalities` - (Optional, Computed, ForceNew, Set) The modalities of the visual LLM. Available since v2.2.0.
* `visual_llm_model` - (Optional, Computed, ForceNew, String) The visual LLM model. Available since v2.2.0.

### `agentic_db_branch_spec`

The `agentic_db_branch_spec` block supports the following:

* `branch_id` - (Optional, ForceNew, String) The ID of the branch. Available since v2.2.0.
* `db_cluster_id` - (Optional, ForceNew, String) The ID of the database cluster. Available since v2.2.0.
* `fork_from_application_id` - (Optional, ForceNew, String) The application ID to fork from. Available since v2.2.0.
* `fork_from_branch` - (Optional, ForceNew, Bool) Whether to fork from a branch. Available since v2.2.0.
* `project_id` - (Optional, ForceNew, String) The ID of the project. Available since v2.2.0.
* `tenant_id` - (Optional, ForceNew, String) The ID of the tenant. Available since v2.2.0.

### `endpoints`

The `endpoints` block supports the following:

* `endpoint_type` - (Optional, ForceNew, String) The type of the endpoint. Available since v2.2.0.
* `description` - (Optional, ForceNew, String) The description of the endpoint. Available since v2.2.0.

### `knowledge_application_spec`

The `knowledge_application_spec` block supports the following:

* `dashboard_password` - (Optional, ForceNew, String) The password of the dashboard. Available since v2.2.0.
* `db_password` - (Optional, ForceNew, String) The password of the database. Available since v2.2.0.
* `llm_model` - (Optional, ForceNew, String) The LLM model. Available since v2.2.0.

### `dnat_entries`

The `dnat_entries` block supports the following:

* `port_name` - (Optional, ForceNew, String) The name of the port. Available since v2.2.0.
* `front_port` - (Optional, ForceNew, Int) The front port. Available since v2.2.0.

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
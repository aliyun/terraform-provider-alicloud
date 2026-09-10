---
subcategory: "Rds Ai"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_ai_instance"
description: |-
  Provides a Alicloud Rds Ai Instance resource.
---

# alicloud_rds_ai_instance

Provides a Rds Ai Instance resource.



For information about Rds Ai Instance and how to use it, see [What is Instance](https://next.api.alibabacloud.com/document/RdsAi/2025-05-07/CreateAppInstance).

-> **NOTE:** Available since v1.268.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_rds_ai_instance&exampleId=4f5fcca8-ae3f-d3c2-3208-68f0bc96bbcdc55f48c5&activeTab=example&spm=docs.r.rds_ai_instance.0.4f5fcca8ae&intl_lang=EN_US" target="_blank">
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

data "alicloud_vswitches" "default" {
  zone_id = "cn-hangzhou-i"
}

resource "alicloud_db_instance" "default" {
  engine                   = "PostgreSQL"
  engine_version           = "17.0"
  db_instance_storage_type = "general_essd"
  instance_type            = "pg.n2.1c.1m"
  instance_storage         = 100
  vswitch_id               = data.alicloud_vswitches.default.ids.0
  instance_name            = var.name
}

resource "alicloud_rds_ai_instance" "default" {
  app_name         = var.name
  app_type         = "supabase"
  db_instance_name = alicloud_db_instance.default.id
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_rds_ai_instance&spm=docs.r.rds_ai_instance.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `app_name` - (Required, ForceNew) The name of the new AI application.
* `app_type` - (Required, ForceNew) Application type. Currently, only `supabase` is supported.
* `auth_config_list` - (Optional, Set) Authentication information list. See [`auth_config_list`](#auth_config_list) below.
* `ca_type` - (Optional) The type of the certificate. Currently, only `custom` is supported. A custom certificate is used.

-> **NOTE:**  When `ssl_enabled` is set to `1`, this parameter must be configured.

* `db_instance_name` - (Optional, ForceNew) The ID of the RDS PostgreSQL database instance accessed by the AI application.
  supports only **newly purchased empty RDS PostgreSQL instances**. The major version is `17`, and the minor version is **20250630 or later**.>
* `dashboard_password` - (Optional) Supabase Dashboard password.
  The password must be 8 to 32 characters in length and contain three or more characters: uppercase letters, lowercase letters, numbers, and underscores (_).
* `database_password` - (Optional) The RDS Database access password.
  The password must be 8 to 32 characters in length and contain three or more characters: uppercase letters, lowercase letters, numbers, and underscores (_).
* `initialize_with_existing_data` - (Optional, Bool) Whether to recover from existing PG data. Valid values:
  - `true`: Yes.
  - `false` (default): No.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `public_endpoint_enabled` - (Optional, Bool) Whether to enable the public network connection address. Valid values:
  - `true` (default): Yes.
  - `false`: No.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `public_network_access_enabled` - (Optional, Bool) Whether to enable the public network NAT gateway. Valid values:
  - `true`: Yes.
  - `false` (default): No.

* `ssl_enabled` - (Optional, Int) Enable or disable SSL. Valid values:
  - `1`: Open
  - `0`: Closed
* `server_cert` - (Optional) Customize the certificate content.

-> **NOTE:**  When `ca_type` is set to `custom`, this parameter must be configured.

* `server_key` - (Optional) The certificate private key.

-> **NOTE:**  When `ca_type` is set to `custom`, this parameter must be configured.

* `status` - (Optional) The status of the instance. Valid values: `Running`, `Stopped`.
* `storage_config_list` - (Optional, Set) A list of storage configurations. See [`storage_config_list`](#storage_config_list) below.

### `auth_config_list`

The auth_config_list supports the following:
* `name` - (Optional) The configuration item name. For more information, see [How to use it](https://www.alibabacloud.com/help/en/rds/apsaradb-rds-for-postgresql/authentication).
* `value` - (Optional) The value of the configuration item.

### `storage_config_list`

The storage_config_list supports the following:
* `name` - (Optional) The configuration item name. For more information, see [How to use it](https://www.alibabacloud.com/help/en/rds/apsaradb-rds-for-postgresql/storage).
* `value` - (Optional) The value of the configuration item.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 20 mins) Used when create the Instance.
* `delete` - (Defaults to 5 mins) Used when delete the Instance.
* `update` - (Defaults to 5 mins) Used when update the Instance.

## Import

Rds Ai Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_rds_ai_instance.example <id>
```

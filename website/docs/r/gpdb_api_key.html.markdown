---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_api_key"
description: |-
  Provides a Alicloud AnalyticDB for PostgreSQL (GPDB) Api Key resource.
---

# alicloud_gpdb_api_key

Provides a AnalyticDB for PostgreSQL (GPDB) Api Key resource.

The API key under a GPDB SaaS workspace.

For information about AnalyticDB for PostgreSQL (GPDB) Api Key and how to use it, see [What is Api Key](https://next.api.alibabacloud.com/document/gpdb/2016-05-03/CreateApiKey).

-> **NOTE:** Available since v1.286.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_gpdb_api_key&exampleId=b8b9f256-a83a-5b90-8ec6-f0fb59b996f6d2460d20&activeTab=example&spm=docs.r.gpdb_api_key.0.b8b9f256a8&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

variable "workspace_id" {
  default = "ws-xxxxxxx"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_gpdb_api_key" "default" {
  workspace_id = var.workspace_id
  key_name     = var.name
  description  = "terraform example"
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_gpdb_api_key&spm=docs.r.gpdb_api_key.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:
* `description` - (Optional, ForceNew) The description of the API key.
* `key_name` - (Required, ForceNew) The name of the API key.
* `service_ids` - (Optional, ForceNew) The list of SaaS service IDs that the API key is authorized to access.
* `workspace_id` - (Required, ForceNew) The ID of the workspace.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<workspace_id>:<key_id>`.
* `key_id` - The ID of the API key.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Api Key.
* `delete` - (Defaults to 5 mins) Used when delete the Api Key.

## Import

AnalyticDB for PostgreSQL (GPDB) Api Key can be imported using the id, e.g.

```shell
$ terraform import alicloud_gpdb_api_key.example <workspace_id>:<key_id>
```

---
subcategory: "Function Compute Service V3 (FCV3)"
layout: "alicloud"
page_title: "Alicloud: alicloud_fcv3_alias"
description: |-
  Provides a Alicloud FCV3 Alias resource.
---

# alicloud_fcv3_alias

Provides a FCV3 Alias resource.

Alias for functions.

For information about FCV3 Alias and how to use it, see [What is Alias](https://www.alibabacloud.com/help/en/functioncompute/developer-reference/api-fc-2023-03-30-createalias).

-> **NOTE:** Available since v1.228.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_fcv3_alias&exampleId=964a205d-f3d9-3bda-b718-c780deee06a81557d045&activeTab=example&spm=docs.r.fcv3_alias.0.964a205df3&intl_lang=EN_US" target="_blank">
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

variable "function_name" {
  default = "flask-3xdg"
}


resource "alicloud_fcv3_alias" "default" {
  version_id    = "1"
  function_name = var.function_name
  description   = "create alias"
  alias_name    = var.name
  additional_version_weight = {
    "2" = "0.5"
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_fcv3_alias&spm=docs.r.fcv3_alias.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `additional_version_weight` - (Optional, Map) Grayscale version
* `alias_name` - (Optional, ForceNew, Computed) Function Alias
* `description` - (Optional) Description
* `function_name` - (Required, ForceNew) Function Name
* `version_id` - (Optional) The version that the alias points

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<function_name>:<alias_name>`.
* `create_time` - The creation time of the resource
* `last_modified_time` - (Available since v1.234.0) Last modification time

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Alias.
* `delete` - (Defaults to 5 mins) Used when delete the Alias.
* `update` - (Defaults to 5 mins) Used when update the Alias.

## Import

FCV3 Alias can be imported using the id, e.g.

```shell
$ terraform import alicloud_fcv3_alias.example <function_name>:<alias_name>
```
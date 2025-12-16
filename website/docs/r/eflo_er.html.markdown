---
subcategory: "Eflo"
layout: "alicloud"
page_title: "Alicloud: alicloud_eflo_er"
description: |-
  Provides a Alicloud Eflo Er resource.
---

# alicloud_eflo_er

Provides a Eflo Er resource.



For information about Eflo Er and how to use it, see [What is Er](https://next.api.alibabacloud.com/document/eflo/2022-05-30/CreateEr).

-> **NOTE:** Available since v1.258.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_eflo_er&exampleId=94843f9f-5e11-4c5f-3240-d11c2ac3acd23cb9cffe&activeTab=example&spm=docs.r.eflo_er.0.94843f9f5e&intl_lang=EN_US" target="_blank">
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

data "alicloud_resource_manager_resource_groups" "default" {}


resource "alicloud_eflo_er" "default" {
  er_name        = "er-example-tf"
  master_zone_id = "cn-hangzhou-a"
  description    = "example"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_eflo_er&spm=docs.r.eflo_er.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `description` - (Optional) Description.
* `er_name` - (Required) Lingjun HUB name
* `master_zone_id` - (Required, ForceNew) Primary zone
* `resource_group_id` - (Optional, Computed) The ID of the resource group instance.
* `tags` - (Optional, Map) Label List

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource
* `region_id` - region information
* `status` - Status

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Er.
* `delete` - (Defaults to 5 mins) Used when delete the Er.
* `update` - (Defaults to 5 mins) Used when update the Er.

## Import

Eflo Er can be imported using the id, e.g.

```shell
$ terraform import alicloud_eflo_er.example <id>
```
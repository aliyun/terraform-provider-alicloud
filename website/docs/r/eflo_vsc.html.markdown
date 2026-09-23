---
subcategory: "Eflo"
layout: "alicloud"
page_title: "Alicloud: alicloud_eflo_vsc"
description: |-
  Provides a Alicloud Eflo Vsc resource.
---

# alicloud_eflo_vsc

Provides a Eflo Vsc resource.

Virtual Storage Channel.

For information about Eflo Vsc and how to use it, see [What is Vsc](https://www.alibabacloud.com/help/en/pai/developer-reference/api-eflo-controller-2022-12-15-createvsc).

-> **NOTE:** Available since v1.250.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_eflo_vsc&exampleId=1d0b20df-8116-170f-16d9-ad4a83d1ec1e78fdc7c9&activeTab=example&spm=docs.r.eflo_vsc.0.1d0b20df81&intl_lang=EN_US" target="_blank">
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

resource "alicloud_eflo_vsc" "default" {
  vsc_type = "primary"
  node_id  = "e01-cn-9me49omda01"
  vsc_name = var.name
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_eflo_vsc&spm=docs.r.eflo_vsc.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `node_id` - (Required, ForceNew) The ID of the Node.
* `resource_group_id` - (Optional) The ID of the resource group.
* `tags` - (Optional, Map) The tag of the resource.
* `vsc_name` - (Optional, ForceNew) The name of the Vsc.
* `vsc_type` - (Optional, ForceNew) The type of the Vsc. Default value: `primary`. Valid values: `primary`.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `status` - The status of the Vsc.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vsc.
* `delete` - (Defaults to 5 mins) Used when delete the Vsc.
* `update` - (Defaults to 5 mins) Used when update the Vsc.

## Import

Eflo Vsc can be imported using the id, e.g.

```shell
$ terraform import alicloud_eflo_vsc.example <id>
```

---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_asset_bind"
description: |-
  Provides a Alicloud Threat Detection Asset Bind resource.
---

# alicloud_threat_detection_asset_bind

Provides a Threat Detection Asset Bind resource.

Asset Binding Information.

For information about Threat Detection Asset Bind and how to use it, see [What is Asset Bind](https://next.api.alibabacloud.com/document/Sas/2018-12-03/UpdatePostPaidBindRel).

-> **NOTE:** Available since v1.247.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_threat_detection_asset_bind&exampleId=c4693264-7d1a-a90c-b23a-da27fc117fd786c79cba&activeTab=example&spm=docs.r.threat_detection_asset_bind.0.c46932647d&intl_lang=EN_US" target="_blank">
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

data "alicloud_threat_detection_assets" "default" {
  machine_types = "ecs"
}

resource "alicloud_threat_detection_asset_bind" "default" {
  uuid         = data.alicloud_threat_detection_assets.default.assets.0.uuid
  auth_version = "5"
}
```

### Deleting `alicloud_threat_detection_asset_bind` or removing it from your configuration

Terraform cannot destroy resource `alicloud_threat_detection_asset_bind`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `auth_version` - (Optional, Int) Bind version.
* `uuid` - (Optional, ForceNew, Computed) The first ID of the resource

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Asset Bind.
* `update` - (Defaults to 5 mins) Used when update the Asset Bind.

## Import

Threat Detection Asset Bind can be imported using the id, e.g.

```shell
$ terraform import alicloud_threat_detection_asset_bind.example <id>
```
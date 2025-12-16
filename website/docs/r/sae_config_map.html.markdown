---
subcategory: "Serverless App Engine (SAE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sae_config_map"
sidebar_current: "docs-alicloud-resource-sae-config-map"
description: |-
  Provides a Alicloud Serverless App Engine (SAE) Config Map resource.
---

# alicloud_sae_config_map

Provides a Serverless App Engine (SAE) Config Map resource.

For information about Serverless App Engine (SAE) Config Map and how to use it, see [What is Config Map](https://www.alibabacloud.com/help/en/sae/latest/create-configmap).

-> **NOTE:** Available since v1.130.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_sae_config_map&exampleId=ed81114a-59d0-8efb-2187-c4f5adac6c2f36bbabd8&activeTab=example&spm=docs.r.sae_config_map.0.ed81114a59&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tf-example"
}
data "alicloud_regions" "default" {
  current = true
}
resource "random_integer" "default" {
  max = 99999
  min = 10000
}
resource "alicloud_sae_namespace" "default" {
  namespace_id              = "${data.alicloud_regions.default.regions.0.id}:example${random_integer.default.result}"
  namespace_name            = var.name
  namespace_description     = var.name
  enable_micro_registration = false
}

resource "alicloud_sae_config_map" "default" {
  data         = jsonencode({ "env.home" : "/root", "env.shell" : "/bin/sh" })
  name         = var.name
  namespace_id = alicloud_sae_namespace.default.namespace_id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_sae_config_map&spm=docs.r.sae_config_map.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `data` - (Required) ConfigMap instance data.
* `description` - (Optional) The Description of ConfigMap.
* `name` - (Required, ForceNew) ConfigMap instance name.
* `namespace_id` - (Required, ForceNew) The NamespaceId of ConfigMap.It can contain 2 to 32 lowercase characters.The value is in format `{RegionId}:{namespace}`

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Config Map.

## Import

Serverless App Engine (SAE) Config Map can be imported using the id, e.g.

```shell
$ terraform import alicloud_sae_config_map.example <id>
```
